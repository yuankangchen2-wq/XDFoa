package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/model"
	"gorm.io/gorm"
)

// SeedDefaults 在空库时植入一期演示数据
func SeedDefaults(db *gorm.DB) error {
	ctx := context.Background()
	var actN int64
	db.Model(&model.Activity{}).Count(&actN)
	if actN > 0 {
		return nil
	}

	// ---------- 活动 ----------
	day := func(y int, m time.Month, d int) time.Time {
		return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	}
	acts := []model.Activity{
		{TenantID: CurrentTenant, Name: "20261014-高中-迂回转学", SourceSystem: "SCRM", Category: "地推活动", Type: "高中", TargetObject: "新高二", IntentionDept: "高中部", IntentionProject: "一对一", StartTime: ptime(day(2026, 10, 15)), EndTime: ptime(day(2026, 10, 20)), Remark: "针对重点高中转学生源地推", ExternalID: "3829506", Code: "HDPC28101463829506", Status: 0, Creator: "陈元素", CreatorDept: "市场部"},
		{TenantID: CurrentTenant, Name: "20261014-高中-中职宣讲", SourceSystem: "SCRM", Category: "讲座沙龙", Type: "高中", TargetObject: "新高三", IntentionDept: "高中部", IntentionProject: "大班课", StartTime: ptime(day(2026, 10, 12)), EndTime: ptime(day(2026, 10, 30)), Remark: "中职院校升学宣讲系列", ExternalID: "3829507", Code: "HDPC28101463829507", Status: 1, Creator: "陈元素", CreatorDept: "市场部"},
		{TenantID: CurrentTenant, Name: "武汉一中秋季地推", SourceSystem: "SCRM", Category: "地推活动", Type: "高中", TargetObject: "高中在读", IntentionDept: "高中部", IntentionProject: "小班课", StartTime: ptime(day(2026, 9, 1)), EndTime: ptime(day(2026, 9, 30)), ExternalID: "3811201", Code: "HDPC28090163811201", Status: 2, Creator: "王主管", CreatorDept: "市场部"},
		{TenantID: CurrentTenant, Name: "江汉区高三体验营", SourceSystem: "SCRM", Category: "体验课", Type: "高中", TargetObject: "新高三", IntentionDept: "高中部", IntentionProject: "集训营", StartTime: ptime(day(2026, 10, 10)), EndTime: ptime(day(2026, 11, 10)), ExternalID: "3827777", Code: "HDPC28101063827777", Status: 1, Creator: "王主管", CreatorDept: "市场部"},
		{TenantID: CurrentTenant, Name: "武昌实验中学家长会", SourceSystem: "SCRM", Category: "讲座沙龙", Type: "高中", TargetObject: "高二", IntentionDept: "高中部", IntentionProject: "一对一", StartTime: ptime(day(2026, 8, 20)), EndTime: ptime(day(2026, 8, 21)), ExternalID: "3808866", Code: "HDPC28082063808866", Status: 2, Creator: "陈元素", CreatorDept: "市场部"},
		{TenantID: CurrentTenant, Name: "初二期末冲刺讲座", SourceSystem: "SCRM", Category: "讲座沙龙", Type: "初中", TargetObject: "初中在读", IntentionDept: "初中部", IntentionProject: "大班课", StartTime: ptime(day(2026, 10, 25)), EndTime: ptime(day(2026, 11, 5)), ExternalID: "3829999", Code: "HDPC28102563829999", Status: 0, Creator: "王主管", CreatorDept: "市场部"},
	}
	if err := db.CreateInBatches(&acts, 50).Error; err != nil {
		return err
	}

	// ---------- 线索 ----------
	names := []string{"魏进青", "都析艺", "李明轩", "王思远", "张子涵", "刘宇航", "陈雨桐", "杨欣怡",
		"赵子豪", "黄嘉怡", "周俊杰", "吴雅静", "徐浩然", "孙若曦", "胡一鸣", "朱梓萱",
		"高志强", "林诗涵", "何俊熙", "郭梦琪", "马天佑", "罗晓萌", "梁文杰", "宋佳怡",
		"郑泽宇", "谢雨泽", "韩雪莹", "唐博文", "冯乐瑶", "于浩宇", "董思远", "萧雅芙",
		"程一鸣", "曹可欣", "袁志远", "邓紫函", "蒋昊然", "沈雨桐", "钱俊熙", "曾俊杰",
		"彭雨萱", "吕子墨", "苏晓彤", "卢伟宸", "范佳琪", "石浩然", "史静姝", "孟子轩", "熊若彤",
		"秦俊杰", "侯梦瑶", "白敬宇", "雷佳音", "龙泽宇", "段艺璇", "章若楠", "薛家宁",
		"岳奇峰", "贺子轩", "赖雨桐", "文卓群", "庞博", "樊振宇", "蓝天野", "施诗"}
	grades := []string{"高三", "高三", "高二", "高一", "初三", "初二"}
	campuses := []string{"武汉校区", "武昌校区", "汉口校区", "江汉校区"}
	schools := []string{"武汉一中", "武汉二中", "华师一附中", "武昌实验中学", "武汉外国语学校"}
	districts := []string{"江汉区", "武昌区", "洪山区", "江岸区"}
	tags := []string{"未知", "高意向", "待跟进", "无意向", "竞品学员"}
	levels := []string{"A", "B", "C", "D"}
	employees := []struct{ id, name string }{
		{"U1001", "美到810"}, {"U1002", "qinz"}, {"U1003", "阳光顾问"},
		{"U1004", "小鹿老师"}, {"U1005", "阿星"}, {"U1007", "陈元康"},
	}

	// 状态配额（名称池长度 49，按配额分配）
	statusPlan := []struct {
		status string
		n      int
		dup    bool
		owner  bool
	}{
		{model.StatusPendingHandover, 12, false, false},
		{model.StatusPendingHandover, 4, true, false},
		{model.StatusPendingFollow, 12, false, true},
		{model.StatusPendingFollow, 3, true, true},
		{model.StatusFollowing, 8, false, true},
		{model.StatusPendingAssign, 6, false, false},
		{model.StatusConverted, 9, false, true},
		{model.StatusStoredArchive, 4, false, true},
		{model.StatusInvalidArchive, 3, false, false},
		{model.StatusRecycled, 3, false, false},
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.Local)

	leads := make([]model.Lead, 0, 64)
	nameIdx := 0
	for _, plan := range statusPlan {
		for i := 0; i < plan.n && nameIdx < len(names); i++ {
			n := names[nameIdx]
			nameIdx++
			act := acts[nameIdx%len(acts)]
			created := day(2026, 8, 14).AddDate(0, 0, (nameIdx*7)%60)
			if plan.status == model.StatusConverted {
				created = day(2026, 10, 1).AddDate(0, 0, (nameIdx*2)%13)
			}
			l := model.Lead{
				TenantID: CurrentTenant,
				Code:     fmt.Sprintf("XXL2026%08d", 10146325+nameIdx),
				Name:     n,
				Phone:    fmt.Sprintf("13%09d", 800000000+nameIdx*137),
				StudentNo: fmt.Sprintf("XDF2026%04d", 1000+nameIdx),
				StudentTag: tags[nameIdx%len(tags)],
				NewOldType: "新生",
				School:   schools[nameIdx%len(schools)],
				Grade:    grades[nameIdx%len(grades)],
				District: districts[nameIdx%len(districts)],
				LeadType: "普通线索", EnterType: "纯新线索",
				Status: plan.status,
				IntentionLevel:   levels[nameIdx%len(levels)],
				IntentionGrade:   grades[nameIdx%len(grades)],
				IntentionDept:    "高中部",
				IntentionProject: []string{"一对一", "小班课", "大班课", "集训营"}[nameIdx%4],
				IntentionCampus:  campuses[nameIdx%len(campuses)],
				IntentionSubject: []string{"数学", "英语", "物理", "全科"}[nameIdx%4],
				ChannelL1: "市场地推", ChannelL2: "学校地推", ChannelL3: []string{"校门口摆点", "班级宣讲", "家长会"}[nameIdx%3],
				SystemSource: "SCRM", CreateMethod: "Excel上传",
				ActivityID: act.ID, ActivityName: act.Name,
				TaskName: act.Name + "导入任务",
				FollowCount: int32((nameIdx*3 + 1) % 9),
				Creator:  "陈元素", CreatorDept: "市场部",
				Remark:      []string{"", "家长晚上有空可电话", "关注冲刺提分", "预算待确认"}[nameIdx%4],
				ImportType:  "市场地推", EnterWay: "Excel上传",
				CreatedAt: created, UpdatedAt: created,
			}
			if plan.status == model.StatusConverted {
				l.EnterType = "激活线索"
				if nameIdx%3 == 0 {
					l.EnterType = "纯新线索"
				}
			}
			if plan.dup {
				l.IsDuplicate = 1
				l.StudentTag = "待跟进"
			}
			if plan.owner {
				emp := employees[nameIdx%len(employees)]
				l.OwnerID = emp.id
				l.OwnerName = emp.name
				t := created.AddDate(0, 0, 1)
				l.AssignTime = &t
				if plan.status == model.StatusHandedOver || plan.status == model.StatusConverted {
					ht := created.AddDate(0, 0, 2)
					l.HandoverTime = &ht
				}
			}
			if plan.status == model.StatusStoredArchive || plan.status == model.StatusInvalidArchive || plan.status == model.StatusRecycled {
				at := created.AddDate(0, 0, 20)
				l.ArchiveTime = &at
				l.ArchiveRule = "30天未跟进归档"
				if plan.status == model.StatusStoredArchive {
					l.ArchiveWay = "入库归档"
				}
				if plan.status == model.StatusInvalidArchive {
					l.ArchiveWay = "无效归档"
				}
				l.EnterType = "激活线索"
			}
			if plan.status == model.StatusFollowing || plan.status == model.StatusPendingFollow {
				ft := today.AddDate(0, 0, -((nameIdx % 10) + 1))
				l.LastFollowAt = &ft
				l.LastFollowContent = []string{"电话未接，短信已发", "家长有意向，约周末到访", "加微信沟通课程方案", "孩子住校，晚上再联系"}[nameIdx%4]
			}
			// 今日待回访 3 条、逾期 6 条
			if (nameIdx%17 == 0) || (nameIdx%19 == 0) || (nameIdx%23 == 0) {
				nf := today
				l.NextFollowAt = &nf
			} else if nameIdx%9 == 0 || nameIdx%11 == 0 || nameIdx%13 == 0 {
				nf := today.AddDate(0, 0, -(nameIdx%5 + 1))
				l.NextFollowAt = &nf
			}
			leads = append(leads, l)
		}
	}
	if err := db.CreateInBatches(&leads, 100).Error; err != nil {
		return err
	}

	// ---------- 导入任务/日志 ----------
	deadline := today.AddDate(0, 0, 3)
	failRows := [][]string{
		{"姓名", "手机号", "学员号", "失败原因"},
		{"钱多多", "13911112222", "XDF20260901", "学员【钱多多】已存在线索（编码XXL202610146331，状态跟进中），不符合线索唯一性"},
		{"孙丽", "13833334444", "XDF20260912", "无效数据：手机号格式不正确"},
		{"周伟", "13755556666", "XDF20260933", "学员【周伟】命中活动商机去重规则"},
	}
	fb, _ := json.Marshal(failRows)
	dd, _ := json.Marshal([]map[string]interface{}{
		{"rule": "手机号/学员号与库内线索重复", "count": 98},
		{"rule": "老学生去重-归档学员重新激活", "count": 18},
	},
	)
	cfg, _ := json.Marshal(map[string]interface{}{"enable": true, "old_student_checked": true})
	finished1 := day(2026, 10, 14).Add(9 * time.Hour + 30*time.Minute)
	tasks := []model.ImportTask{
		{TenantID: CurrentTenant, TaskName: "20261014-高中-中职宣讲导入", ActivityID: acts[1].ID, ActivityName: acts[1].Name,
			ImportMethod: "phone", DataSource: "市场地推", ChannelL1: "市场地推", ChannelL2: "学校地推", ChannelL3: "班级宣讲",
			SubmitType: "supervisor", DedupConfig: string(cfg), Operator: "陈元素", SerialNo: "LS20261014093001",
			BatchNo: "PC20261014001", Status: "已完成", FileName: "高中宣讲名单_1014.xlsx",
			TotalCount: 581, SuccessCount: 403, NewCount: 7, ActiveCount: 396, FailCount: 178,
			DuplicateCount: 62, InvalidCount: 12, UniquenessFailCount: 3, DedupHitCount: 116,
			DedupConditionCount: 2, DedupDetails: string(dd), FailRows: string(fb),
			StartedAt: &finished1, FinishedAt: &finished1, CreatedAt: finished1, UpdatedAt: finished1},
		{TenantID: CurrentTenant, TaskName: "20261010-高一-秋季地推", ActivityID: acts[2].ID, ActivityName: acts[2].Name,
			ImportMethod: "phone", DataSource: "市场地推", ChannelL1: "市场地推", ChannelL2: "学校地推", ChannelL3: "校门口摆点",
			SubmitType: "employee", Operator: "美到810", SerialNo: "LS20261010150022", BatchNo: "PC20261010002",
			Status: "已完成", TotalCount: 220, SuccessCount: 205, NewCount: 30, ActiveCount: 175, FailCount: 15,
			DuplicateCount: 9, InvalidCount: 6, DedupHitCount: 9, DedupConditionCount: 1,
			DedupDetails: `[{"rule":"手机号与库内线索重复","count":9}]`,
			CreatedAt: day(2026, 10, 10), UpdatedAt: day(2026, 10, 10)},
		{TenantID: CurrentTenant, TaskName: "高中地推标准任务", IsTemplate: 1, ActivityID: acts[1].ID, ActivityName: acts[1].Name,
			ImportMethod: "phone", DataSource: "市场地推", ChannelL1: "市场地推", ChannelL2: "学校地推",
			SubmitType: "supervisor", DedupConfig: string(cfg), Operator: "陈元素", Status: "模板",
			CreatedAt: day(2026, 9, 20), UpdatedAt: day(2026, 9, 20)},
		{TenantID: CurrentTenant, TaskName: "20261014-高中-晚间批次", ActivityID: acts[3].ID, ActivityName: acts[3].Name,
			ImportMethod: "phone", DataSource: "大数据图文", ChannelL1: "大数据图文", ChannelL2: "信息流投放", ChannelL3: "抖音投放",
			SubmitType: "none", Operator: "陈元素", SerialNo: "LS20261014200003", BatchNo: "PC20261014003",
			Status: "进行中", FileName: "晚间投放线索.xlsx", Deadline: &deadline,
			CreatedAt: today, UpdatedAt: today},
	}
	if err := db.CreateInBatches(&tasks, 50).Error; err != nil {
		return err
	}

	// ---------- 数据进入记录 ----------
	recStatus := []struct {
		name, phone, enter, fail, l1, l2 string
	}{
		{"都析艺", "13912340001", "成功_新建线索", "", "市场地推", "学校地推"},
		{"李明轩", "13912340002", "成功_激活线索", "", "市场地推", "学校地推"},
		{"钱多多", "13912340003", "失败_线索重复", "学员【钱多多】命中重复，库内已有跟进中线索", "市场地推", "学校地推"},
		{"孙丽", "13912340004", "失败_无效数据", "手机号格式不正确", "大数据图文", "信息流投放"},
		{"周伟", "13912340005", "成功_新建线索", "", "大数据图文", "信息流投放"},
		{"吴琼", "13912340006", "成功_激活线索", "", "转介绍", "学员转介绍"},
		{"郑凯", "13912340007", "失败_线索重复", "学员【郑凯】手机号命中活动商机去重", "市场地推", "社区地推"},
		{"王芳", "13912340008", "成功_新建线索", "", "线上渠道", "官网留资"},
		{"冯雷", "13912340009", "成功_激活线索", "", "市场地推", "学校地推"},
		{"陈晨", "13912340010", "失败_线索重复", "不符合线索唯一性：手机号与学员号指向不同学员", "市场地推", "学校地推"},
		{"褚静", "13912340011", "成功_新建线索", "", "转介绍", "员工推荐"},
		{"卫蓝", "13912340012", "成功_激活线索", "", "市场地推", "异业合作"},
		{"蒋欣", "13912340013", "成功_新建线索", "", "线上渠道", "小程序"},
		{"沈月", "13912340014", "失败_无效数据", "缺少学员姓名", "大数据图文", "短信触达"},
	}
	records := make([]model.DataRecord, 0, len(recStatus))
	for i, r := range recStatus {
		grade := grades[i%len(grades)]
		records = append(records, model.DataRecord{
			TenantID: CurrentTenant, ImportTaskID: tasks[0].ID,
			RelatedName: r.name, RelatedPhone: r.phone, RelatedStudentNo: fmt.Sprintf("XDF2026%04d", 2000+i),
			ActivityStudentName: r.name, ActivityPhone: r.phone, ImportType: "大数据图文",
			ActivityGrade: grade, RelatedGrade: grade, EnterWay: "Excel上传",
			LeadEnterType: r.enter, FailReason: r.fail,
			ChannelL1: r.l1, ChannelL2: r.l2, ChannelL3: "校门口摆点",
			ActivityName: acts[1].Name, IntentionDept: "高中部", IntentionProject: "一对一",
			OwnerName: "", CreatedAt: day(2026, 10, 14).Add(time.Duration(i) * time.Hour),
		})
	}
	if err := db.CreateInBatches(&records, 100).Error; err != nil {
		return err
	}

	// ---------- 异常线索 ----------
	excs := []model.ExceptionLead{
		{TenantID: CurrentTenant, RelatedName: "钱多多", RelatedPhone: "13912340003", RelatedStudentNo: "XDF20260901", ActivityStudentName: "钱多", ActivityPhone: "13912340003", ChannelL1: "市场地推", ChannelL2: "学校地推", ChannelL3: "班级宣讲", SystemSourceL1: "SCRM", SystemSourceL2: "Excel导入", ActivityName: acts[1].Name, DataRemark: "活动数据姓名与库内学员不一致", IntentionDept: "高中部", IntentionProject: "一对一", Status: "待处理"},
		{TenantID: CurrentTenant, RelatedName: "陈晨", RelatedPhone: "13912340010", RelatedStudentNo: "XDF20260945", ActivityStudentName: "陈辰", ActivityPhone: "13912340010", ChannelL1: "市场地推", ChannelL2: "学校地推", SystemSourceL1: "SCRM", SystemSourceL2: "Excel导入", ActivityName: acts[1].Name, DataRemark: "手机号与学员号指向不同学员", IntentionDept: "高中部", IntentionProject: "大班课", Status: "待处理"},
		{TenantID: CurrentTenant, RelatedName: "郑凯", RelatedPhone: "13912340007", ActivityStudentName: "郑恺", ActivityPhone: "13912340007", ChannelL1: "大数据图文", ChannelL2: "信息流投放", ChannelL3: "抖音投放", SystemSourceL1: "SCRM", SystemSourceL2: "投放回传", ActivityName: acts[3].Name, DataRemark: "同一条数据命中两个活动", IntentionDept: "高中部", IntentionProject: "集训营", Status: "待处理"},
		{TenantID: CurrentTenant, RelatedName: "沈月", RelatedPhone: "13912340014", ActivityStudentName: "", ActivityPhone: "13912340014", ChannelL1: "大数据图文", ChannelL2: "短信触达", SystemSourceL1: "SCRM", SystemSourceL2: "短信平台", ActivityName: acts[3].Name, DataRemark: "缺少学员姓名，信息不完整", IntentionDept: "高中部", IntentionProject: "小班课", Status: "待处理"},
		{TenantID: CurrentTenant, RelatedName: "刘洋", RelatedPhone: "13866667777", ActivityStudentName: "刘洋", ActivityPhone: "13866660000", ChannelL1: "市场地推", ChannelL2: "社区地推", SystemSourceL1: "SCRM", SystemSourceL2: "Excel导入", ActivityName: acts[2].Name, DataRemark: "活动手机号与库内手机号不一致", IntentionDept: "高中部", Status: "已忽略"},
		{TenantID: CurrentTenant, RelatedName: "杨帆", RelatedPhone: "13788889999", ActivityStudentName: "杨凡", ActivityPhone: "13788889999", ChannelL1: "转介绍", ChannelL2: "学员转介绍", SystemSourceL1: "SCRM", SystemSourceL2: "手工录入", ActivityName: acts[3].Name, DataRemark: "学员标签缺失", IntentionDept: "高中部", Status: "待处理"},
		{TenantID: CurrentTenant, RelatedName: "杜甫", RelatedPhone: "13612123434", ActivityStudentName: "杜甫", ActivityPhone: "13612123434", ChannelL1: "市场地推", ChannelL2: "异业合作", SystemSourceL1: "SCRM", SystemSourceL2: "合作方回传", ActivityName: acts[0].Name, DataRemark: "意向部门为空", IntentionDept: "", Status: "已忽略"},
		{TenantID: CurrentTenant, RelatedName: "许可", RelatedPhone: "13501010101", ActivityStudentName: "许可", ActivityPhone: "13501010100", ChannelL1: "线上渠道", ChannelL2: "公众号", SystemSourceL1: "SCRM", SystemSourceL2: "公众号留资", ActivityName: acts[4].Name, DataRemark: "手机号尾号冲突待核实", IntentionDept: "高中部", Status: "待处理"},
	}
	if err := db.CreateInBatches(&excs, 50).Error; err != nil {
		return err
	}

	// ---------- 日程 ----------
	sched := func(d int, title, tm, typ, status string) model.Schedule {
		return model.Schedule{
			TenantID: CurrentTenant, UserID: CurrentUserID, Title: title, Type: typ, Status: status,
			EventDate: day(2026, 10, d), EventTime: tm, CreatedAt: day(2026, 9, 28),
		}
	}
	schedules := []model.Schedule{
		sched(12, "部门周例会", "09:30", "日程", "已确访"),
		{TenantID: CurrentTenant, UserID: CurrentUserID, Title: "王同学体验课到访确认", Type: "体验课", Status: "待确访", EventDate: today, EventTime: "10:00", CreatedAt: day(2026, 10, 10)},
		sched(15, "李同学家长电话回访", "14:00", "日程", "待确访"),
		sched(16, "体验课到访未转化跟进-张同学", "16:00", "体验课", "到访未转化"),
		sched(20, "季度经营复盘会", "14:00", "日程", "已确访"),
		sched(22, "到访未转化家长二次邀约", "11:00", "日程", "到访未转化"),
		sched(28, "11月活动排期讨论", "15:30", "日程", "待确访"),
		{TenantID: CurrentTenant, UserID: CurrentUserID, Title: "11月部门启动会", Type: "日程", Status: "已确访", EventDate: day(2026, 11, 2), EventTime: "09:00", CreatedAt: day(2026, 10, 8)},
	}
	if err := db.CreateInBatches(&schedules, 50).Error; err != nil {
		return err
	}

	// ---------- 顾问业绩 ----------
	perf := func(month, campus, uid, name, amount string, persons, visits int32) model.ConsultantPerf {
		return model.ConsultantPerf{TenantID: CurrentTenant, Month: month, Campus: campus, UserID: uid, Name: name, DealAmount: amount, DealPersons: persons, DealVisits: visits}
	}
	perfs := []model.ConsultantPerf{
		perf("2026-10", "武汉校区", "U1001", "美到810", "134500.00", 12, 14),
		perf("2026-10", "武汉校区", "U1002", "qinz", "118800.00", 10, 11),
		perf("2026-10", "武汉校区", "U1003", "阳光顾问", "96500.00", 9, 10),
		perf("2026-10", "江汉校区", "U1004", "小鹿老师", "84200.00", 7, 9),
		perf("2026-10", "汉口校区", "U1005", "阿星", "62000.00", 6, 7),
		perf("2026-10", "武昌校区", "U1006", "王老师", "57800.00", 5, 6),
		perf("2026-10", "武汉校区", "U1007", "陈元康", "43200.00", 4, 5),
		perf("2026-09", "武汉校区", "U1001", "美到810", "152000.00", 14, 16),
		perf("2026-09", "武汉校区", "U1002", "qinz", "101000.00", 9, 12),
		perf("2026-09", "武昌校区", "U1006", "王老师", "66000.00", 6, 8),
		perf("2026-09", "江汉校区", "U1004", "小鹿老师", "59000.00", 5, 7),
	}
	if err := db.CreateInBatches(&perfs, 50).Error; err != nil {
		return err
	}

	// ---------- 待办数量种子 ----------
	seed := func(role, key string, cnt int32) model.TodoSeed {
		return model.TodoSeed{TenantID: CurrentTenant, Role: role, TodoKey: key, Count: cnt}
	}
	employeeSeeds := []model.TodoSeed{
		seed("employee", "opp_follow", 18), seed("employee", "opp_active_revisit", 12),
		seed("employee", "opp_today_revisit", 5), seed("employee", "opp_overdue", 7),
		seed("employee", "opp_appeal_pending", 2), seed("employee", "opp_3day_recycle", 4),
		seed("employee", "opp_3day_archive", 3),
		seed("employee", "trial_confirm", 6), seed("employee", "trial_visit_notconv", 4),
		seed("employee", "task_expiring_3d", 2), seed("employee", "ticket_refund_follow", 3),
		seed("employee", "new_sign_schedule", 5), seed("employee", "first_lesson", 2),
		seed("employee", "student_handover", 8),
	}
	superSeeds := []model.TodoSeed{
		seed("supervisor", "schedule_confirm", 0), seed("supervisor", "schedule_visit_notconv", 0),
		seed("supervisor", "lead_handover", 0),
	}
	if err := db.CreateInBatches(&employeeSeeds, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&superSeeds, 10).Error; err != nil {
		return err
	}
	_ = ctx
	return nil
}

func ptime(t time.Time) *time.Time { return &t }
