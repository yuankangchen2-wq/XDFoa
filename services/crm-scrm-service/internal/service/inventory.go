package service

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/model"
	"github.com/oa-portal/crm-scrm-service/internal/repo"
)

// 商机盘点：预聚合宽表（运行时内存聚合，演示口径）

type InvCard struct {
	Key    string  `json:"key"`
	Title  string  `json:"title"`
	Value  float64 `json:"value"`
	Unit   string  `json:"unit"`
}

type InvMonth struct {
	Month string `json:"month"`
	S     int    `json:"s"`
	A     int    `json:"a"`
	B     int    `json:"b"`
	C     int    `json:"c"`
	Deals int    `json:"deals"`
	Conv  string `json:"conv"`
}

type InvMetric struct {
	Key       string  `json:"key"`
	Title     string  `json:"title"`
	AvgScore  float64 `json:"avg_score"`
	HitRate   string  `json:"hit_rate"`
	HitCount  int     `json:"hit_count"`
	Total     int     `json:"total"`
}

type InvRow struct {
	Name      string     `json:"name"`
	Total     int        `json:"total"`
	S         int        `json:"s"`
	A         int        `json:"a"`
	B         int        `json:"b"`
	C         int        `json:"c"`
	Deals     int        `json:"deals"`
	Conv      string     `json:"conv"`
	Months    []InvMonth `json:"months"`
	Metrics   []InvMetric `json:"metrics"`
}

type InvResp struct {
	Rows       []InvRow  `json:"rows"`
	Yesterday  []InvCard `json:"yesterday"`
	Circle     []InvCard `json:"circle"`
	GlobalMetric []InvMetric `json:"global_metric"`
}

// metricScores 计算四个指标分类的命中率分数（0~100）
func metricScores(o model.Opportunity) map[string]int {
	client := 40
	if o.NewOldType != "" { client += 15 }
	if o.School != "" { client += 15 }
	if o.Grade != "" { client += 10 }
	if o.Phone != "" { client += 10 }
	if o.AddWechat == 1 { client += 10 }

	activity := 35
	if o.ActivityName != "" { activity += 20 }
	if o.ChannelL1 != "" { activity += 15 }
	if o.PromiseVisitAt != nil { activity += 15 }
	if o.VisitConverted == 1 { activity += 15 }

	comm := 30
	comm += int(o.FollowCount) * 8
	if comm > 90 { comm = 90 }
	if o.LastFollowContent != "" { comm += 8 }
	if o.ContactPoint == "呼入" { comm += 5 }
	if comm > 100 { comm = 100 }

	pay := 25
	if o.Status == model.OppStatusDealt { pay = 95 } else if o.Grade2 == "S" { pay += 35 } else if o.Grade2 == "A" { pay += 20 }
	if o.VisitConverted == 1 { pay += 15 }
	if o.Status == model.OppStatusRecycled || o.Status == model.OppStatusDeadArchive { pay -= 20 }
	if pay < 0 { pay = 0 }
	if pay > 100 { pay = 100 }
	return map[string]int{"client": client, "activity": activity, "communication": comm, "payment": pay}
}

var metricMeta = []struct{ Key, Title string }{
	{"client", "客户信息"},
	{"activity", "活动行为"},
	{"communication", "沟通质量"},
	{"payment", "付费行为"},
}

func pct(n, d int) string {
	if d == 0 { return "0%" }
	return formatFloat(float64(n)*100/float64(d)) + "%"
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 1, 64)
}

func (s *ScrmService) Inventory(ctx context.Context, dim string) (*InvResp, error) {
	opps, err := s.oppRepo.ListAll(ctx, CurrentTenant, repo.OppFilter{}, 5000)
	if err != nil {
		return nil, err
	}
	deals, err := s.dealRepo.List(ctx, CurrentTenant, repo.DealFilter{}, repo.Page{Page: 1, PageSize: 5000})
	if err != nil {
		return nil, err
	}
	dealCodes := map[string]bool{}
	for _, d := range deals.Items { dealCodes[d.OppCode] = true }

	type agg struct {
		opps   []model.Opportunity
		metric [4]struct{ sum, hit, n int }
	}
	groups := map[string]*agg{}
	for i := range opps {
		o := opps[i]
		key := ""
		switch dim {
		case "dept":
			key = o.OwnerDept
		case "consultant":
			key = o.OwnerName
		default:
			key = o.School
		}
		if key == "" { key = "未分配" }
		g, ok := groups[key]
		if !ok { g = &agg{}; groups[key] = g }
		g.opps = append(g.opps, o)
		ms := metricScores(o)
		for i, m := range metricMeta {
			v := ms[m.Key]
			g.metric[i].sum += v
			g.metric[i].n++
			if v >= 60 { g.metric[i].hit++ }
		}
	}

	// 近6个月
	now := time.Now()
	months := make([]string, 6)
	for i := 5; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		months[5-i] = t.Format("2006-01")
	}

	names := make([]string, 0, len(groups))
	for k := range groups { names = append(names, k) }
	sort.Strings(names)
	rows := make([]InvRow, 0, len(names))
	var globalMetric [4]struct{ sum, hit, n int }
	for _, name := range names {
		g := groups[name]
		r := InvRow{Name: name}
		mm := map[string]*InvMonth{}
		for _, mo := range months { mm[mo] = &InvMonth{Month: mo} }
		for _, o := range g.opps {
			r.Total++
			switch o.Grade2 {
			case "S": r.S++
			case "A": r.A++
			case "B": r.B++
			default: r.C++
			}
			if dealCodes[o.Code] { r.Deals++ }
			if mc, ok := mm[o.CreatedAt.Format("2006-01")]; ok {
				switch o.Grade2 {
				case "S": mc.S++
				case "A": mc.A++
				case "B": mc.B++
				default: mc.C++
				}
				if dealCodes[o.Code] { mc.Deals++ }
			}
		}
		r.Conv = pct(r.Deals, r.Total)
		for _, mo := range months {
			mc := mm[mo]
			tot := mc.S + mc.A + mc.B + mc.C
			mc.Conv = pct(mc.Deals, tot)
			r.Months = append(r.Months, *mc)
		}
		for i, m := range metricMeta {
			d := g.metric[i]
			avg := 0.0
			if d.n > 0 { avg = float64(d.sum) / float64(d.n) }
			r.Metrics = append(r.Metrics, InvMetric{
				Key: m.Key, Title: m.Title,
				AvgScore: float64(int(avg*10+0.5)) / 10,
				HitRate: pct(d.hit, d.n), HitCount: d.hit, Total: d.n,
			})
			globalMetric[i].sum += d.sum
			globalMetric[i].hit += d.hit
			globalMetric[i].n += d.n
		}
		rows = append(rows, r)
	}

	gm := make([]InvMetric, 0, 4)
	for i, m := range metricMeta {
		d := globalMetric[i]
		avg := 0.0
		if d.n > 0 { avg = float64(d.sum) / float64(d.n) }
		gm = append(gm, InvMetric{Key: m.Key, Title: m.Title,
			AvgScore: float64(int(avg*10+0.5)) / 10, HitRate: pct(d.hit, d.n), HitCount: d.hit, Total: d.n})
	}

	// 昨日统计 6 卡
	today := now.Format("2006-01-02")
	var sStuck, bDrop, cStuck int
	var up, down, same int
	followByGrade := map[string]struct{ n, sum int }{}
	for _, o := range opps {
		if o.LastFollowAt != nil && o.LastFollowAt.Format("2006-01-02") == today {
			switch {
			case o.Score < o.InitScore:
				down++
				if o.Grade2 == "B" { bDrop++ }
			case o.Score > o.InitScore:
				up++
			default:
				same++
				if o.Grade2 == "S" { sStuck++ }
				if o.Grade2 == "C" { cStuck++ }
			}
		}
		if o.Grade2 != "" {
			v := followByGrade[o.Grade2]
			v.n++
			v.sum += int(o.FollowCount)
			followByGrade[o.Grade2] = v
		}
	}
	circle := []InvCard{}
	for _, g := range []string{"S", "A", "B", "C"} {
		v := followByGrade[g]
		avg := 0.0
		if v.n > 0 { avg = float64(v.sum) / float64(v.n) }
		circle = append(circle, InvCard{Key: g, Title: g + "级平均跟进次数", Value: float64(int(avg*10+0.5)) / 10, Unit: "次"})
	}

	return &InvResp{
		Rows: rows,
		GlobalMetric: gm,
		Yesterday: []InvCard{
			{Key: "s_stuck", Title: "昨日S级分数不升", Value: float64(sStuck), Unit: "个"},
			{Key: "b_drop", Title: "昨日B级分数下滑", Value: float64(bDrop), Unit: "个"},
			{Key: "c_stuck", Title: "昨日C级分数不升", Value: float64(cStuck), Unit: "个"},
			{Key: "up", Title: "昨日分数上涨", Value: float64(up), Unit: "个"},
			{Key: "down", Title: "昨日分数下滑", Value: float64(down), Unit: "个"},
			{Key: "same", Title: "昨日分数持平", Value: float64(same), Unit: "个"},
		},
		Circle: circle,
	}, nil
}
