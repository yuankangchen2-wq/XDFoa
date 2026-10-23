package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/model"
	"gorm.io/gorm"
)

type seedEmp struct{ id, name, dept string }

var seedEmps = []seedEmp{
	{"U1001", "美到810", "高中部-江汉组"},
	{"U1002", "qinz", "高中部-武昌组"},
	{"U1003", "阳光顾问", "高中部-汉口组"},
	{"U1004", "小鹿老师", "初中部-江汉组"},
	{"U1005", "阿星", "高中部-武昌组"},
	{"U1007", "陈元康", "高中部-江汉组"},
}

func countAny(db *gorm.DB, m interface{}, where string, args ...interface{}) int64 {
	var n int64
	q := db.Model(m)
	if where != "" {
		q = q.Where(where, args...)
	}
	q.Count(&n)
	return n
}

// SeedP2P3 二三期演示数据（各表独立幂等）
func SeedP2P3(db *gorm.DB) error {
	if err := seedLeadP2Fields(db); err != nil {
		return err
	}
	if err := seedOpportunities(db); err != nil {
		return err
	}
	// 存量数据修复：保证本周诺访/到访未转化池有演示数据（幂等）
	if err := db.Exec(`UPDATE scrm_opportunity SET status = '跟进中', promise_visit_at = DATE_ADD(CURDATE(), INTERVAL 3 DAY),
		visit_converted = 0, owner_id = 1, owner_name = '陈元素'
		WHERE tenant_id = ? AND status = '跟进中' AND promise_visit_at IS NULL
		ORDER BY id LIMIT 3`, CurrentTenant).Error; err != nil {
		return err
	}
	// 今日高意向演示数据（幂等，每次重启固定落在最早两条跟进商机上）
	if err := db.Exec(`UPDATE scrm_opportunity SET level = '高', last_follow_at = NOW(), owner_id = 1, owner_name = '陈元素'
		WHERE tenant_id = ? AND status IN ('跟进中','待跟进')
		ORDER BY id LIMIT 2`, CurrentTenant).Error; err != nil {
		return err
	}
	// 保证当前登录顾问（陈元素/ID=1）名下有商机，员工工作台非空（幂等）
	if err := db.Exec(`UPDATE scrm_opportunity SET owner_id = 1, owner_name = '陈元素'
		WHERE tenant_id = ? AND (
			id % 4 = 0
			OR (next_follow_at IS NOT NULL AND DATE(next_follow_at) = CURDATE())
			OR (planned_recycle_at IS NOT NULL AND planned_recycle_at BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 3 DAY))
			OR (promise_visit_at IS NOT NULL AND YEARWEEK(promise_visit_at, 1) = YEARWEEK(CURDATE(), 1))
		)`, CurrentTenant).Error; err != nil {
		return err
	}
	if err := seedAssignRecords(db); err != nil {
		return err
	}
	if err := seedEmployeeTasks(db); err != nil {
		return err
	}
	if err := seedAppeals(db); err != nil {
		return err
	}
	if err := seedDealApprovals(db); err != nil {
		return err
	}
	if err := seedDealLogs(db); err != nil {
		return err
	}
	if err := seedTickets(db); err != nil {
		return err
	}
	if err := seedMallOrders(db); err != nil {
		return err
	}
	return nil
}

// seedLeadP2Fields 回填二期新增列演示值，并保证回收/归档池有演示线索（幂等）
func seedLeadP2Fields(db *gorm.DB) error {
	// 新增列的演示值
	db.Exec(`UPDATE scrm_lead SET total_follow_count = follow_count + 1 WHERE tenant_id = ? AND total_follow_count = 0`, CurrentTenant)
	db.Exec(`UPDATE scrm_lead SET ai_call_count = id % 4 WHERE tenant_id = ? AND ai_call_count = 0`, CurrentTenant)
	db.Exec(`UPDATE scrm_lead SET ai_call_at = updated_at WHERE tenant_id = ? AND ai_call_count > 0 AND ai_call_at IS NULL`, CurrentTenant)
	db.Exec(`UPDATE scrm_lead SET deduct_wechat = 1 WHERE tenant_id = ? AND deduct_wechat = 0 AND id % 3 = 0`, CurrentTenant)
	db.Exec(`UPDATE scrm_lead SET business_channel_l2 = '线下地推', business_channel_l3 = '校门口摆点'
		WHERE tenant_id = ? AND (business_channel_l2 IS NULL OR business_channel_l2 = '')`, CurrentTenant)

	now := time.Now()
	ensureStatus := func(target string, want int, from []string, cols map[string]interface{}) error {
		var n int64
		db.Model(&model.Lead{}).Where("tenant_id = ? AND status = ?", CurrentTenant, target).Count(&n)
		diff := want - int(n)
		if diff <= 0 {
			return nil
		}
		var ids []int64
		if err := db.Model(&model.Lead{}).Where("tenant_id = ? AND status IN ?", CurrentTenant, from).
			Order("id DESC").Limit(diff).Pluck("id", &ids).Error; err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}
		updates := map[string]interface{}{"status": target}
		for k, v := range cols {
			updates[k] = v
		}
		return db.Model(&model.Lead{}).Where("id IN ?", ids).Updates(updates).Error
	}

	if err := ensureStatus(model.StatusRecycled, 3,
		[]string{model.StatusPendingAssign, model.StatusPendingFollow},
		map[string]interface{}{
			"recycle_reason": "30天未跟进自动回收", "recycle_at": now.AddDate(0, 0, -2),
			"archive_time": nil, "archive_way": "",
		}); err != nil {
		return err
	}
	if err := ensureStatus(model.StatusInvalidArchive, 2,
		[]string{model.StatusPendingFollow, model.StatusFollowing},
		map[string]interface{}{
			"archive_way": "无效归档", "archive_rule": "空号停机", "archive_reason": "空号停机",
			"archive_operator": "王主管", "archive_time": now.AddDate(0, 0, -1),
		}); err != nil {
		return err
	}
	if err := ensureStatus(model.StatusStoredArchive, 2,
		[]string{model.StatusFollowing, model.StatusPendingFollow},
		map[string]interface{}{
			"archive_way": "入库归档", "archive_rule": "重复数据入库", "archive_reason": "重复数据入库",
			"archive_operator": "王主管", "archive_time": now.AddDate(0, 0, -3),
		}); err != nil {
		return err
	}
	return nil
}

func seedOpportunities(db *gorm.DB) error {
	if countAny(db, &model.Opportunity{}, "tenant_id = ?", CurrentTenant) > 0 {
		return nil
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.Local)
	mon := func() time.Time {
		wd := int(now.Weekday())
		if wd == 0 { wd = 7 }
		return time.Date(now.Year(), now.Month(), now.Day()-wd+1, 0, 0, 0, 0, time.Local)
	}()
	names := []string{"魏进青", "都析艺", "李明轩", "王思远", "张子涵", "刘宇航", "陈雨桐", "杨欣怡",
		"赵子豪", "黄嘉怡", "周俊杰", "吴雅静", "徐浩然", "孙若曦", "胡一鸣", "朱梓萱",
		"高志强", "林诗涵", "何俊熙", "郭梦琪", "马天佑", "罗晓萌", "梁文杰", "宋佳怡",
		"郑泽宇", "谢雨泽", "韩雪莹", "唐博文", "冯乐瑶", "于浩宇", "董思远", "萧雅芙",
		"程一鸣", "曹可欣", "袁志远", "邓紫函", "蒋昊然", "沈雨桐", "钱俊熙", "曾俊杰"}
	statuses := []string{
		model.OppStatusPendingFollow, model.OppStatusFollowing, model.OppStatusFollowing,
		model.OppStatusPendingFollow, model.OppStatusFollowing, model.OppStatusRecycled,
		model.OppStatusPendingAssign, model.OppStatusDeadArchive, model.OppStatusDealt,
	}
	levels := []string{"高", "中", "低"}
	grades := []string{"高三", "高二", "高一", "初三"}
	schools := []string{"武汉一中", "武汉二中", "华师一附中", "武昌实验中学", "武汉外国语学校"}
	districts := []string{"江汉区", "武昌区", "洪山区", "江岸区"}
	channels := [][3]string{
		{"市场地推", "学校地推", "校门口摆点"},
		{"市场地推", "学校地推", "班级宣讲"},
		{"线上推广", "抖音投放", "信息流广告"},
		{"转介绍", "老带新", "家长推荐"},
	}
	projects := []string{"一对一", "小班课", "大班课", "集训营"}
	contactPoints := []string{"呼入", "呼出"}
	list := make([]model.Opportunity, 0, 48)
	for i := 0; i < 48; i++ {
		emp := seedEmps[i%len(seedEmps)]
		st := statuses[i%len(statuses)]
		created := now.AddDate(0, 0, -(i * 3 % 150))
		ch := channels[i%len(channels)]
		o := model.Opportunity{
			TenantID: CurrentTenant,
			Code:     fmt.Sprintf("XSJ%s%06d", now.Format("060102"), 100001+i),
			Name:     names[i%len(names)],
			Phone:    fmt.Sprintf("13%09d", 900000000+i*173),
			StudentNo: fmt.Sprintf("XDF2026%04d", 3000+i),
			NewOldType: []string{"新生", "老学员"}[i%2],
			Grade: grades[i%len(grades)],
			School: schools[i%len(schools)],
			District: districts[i%len(districts)],
			AddWechat: int32(i % 2),
			Status: st,
			EnterType: []string{"直接创建商机", "激活商机", "优秀跟进创建"}[i%3],
			OppType: []string{"常规商机", "转介绍商机", "活动商机"}[i%3],
			Level: levels[i%len(levels)],
			IntentionDept: "高中部",
			IntentionProject: projects[i%4],
			IntentionCampus: []string{"江汉校区", "武昌校区", "汉口校区"}[i%3],
			IntentionGrade: grades[i%len(grades)],
			IntentionSubject: []string{"数学", "英语", "物理"}[i%3],
			ContactPoint: contactPoints[i%2],
			ActivityName: "20261014-高中-迂回转学",
			ChannelL1: ch[0], ChannelL2: ch[1], ChannelL3: ch[2],
			BusinessChannel: ch[0], SystemSource: "SCRM",
			CreateMethod: []string{"搜索手工创建", "导入"}[i%2],
			LeadCode: fmt.Sprintf("XXL2026%08d", 10146325+i),
			OwnerID: emp.id, OwnerName: emp.name, OwnerDept: emp.dept,
			FollowCount: int32(i%8),
			Creator: "陈元素", CreatorDept: "市场部",
			CreatedAt: created, UpdatedAt: created,
		}
		assignT := created.Add(2 * time.Hour)
		o.AssignAt = &assignT
		if i%3 != 0 {
			lf := today.AddDate(0, 0, -(i % 12))
			o.LastFollowAt = &lf
			o.LastFollowContent = []string{"家长在对比竞品，强调师资", "约了本周六下午到访", "孩子想补数学，预算OK", "微信发了课程表"}[i%4]
			o.LastFollowBy = emp.name
		} else if i%12 == 0 {
			// 今日高意向（高等级 + 今日跟进）
			lf := today
			o.LastFollowAt = &lf
			o.LastFollowContent = "今日到店咨询，意向高，家长在比价"
			o.LastFollowBy = emp.name
		}
		if i%5 == 0 {
			nf := today.AddDate(0, 0, i%4)
			o.NextFollowAt = &nf
		}
		// 本周诺访 / 到访未转化（诺访行必须为跟进中才会进入备单池）
		if i%6 == 0 {
			pv := mon.AddDate(0, 0, 3+i%4)
			o.PromiseVisitAt = &pv
			o.Status = model.OppStatusFollowing
			st = model.OppStatusFollowing
		}
		// 3天内即将回收
		if i%7 == 0 {
			pr := now.AddDate(0, 0, 1+i%3)
			o.PlannedRecycleAt = &pr
		}
		if st == model.OppStatusRecycled {
			rt := created.AddDate(0, 0, 30)
			o.RecycleAt = &rt
			o.RecycleReason = []string{"30天未跟进自动回收", "主管手动回收-线索质量差", "撞单回收"}[i%3]
		}
		if st == model.OppStatusDeadArchive {
			at := created.AddDate(0, 0, 45)
			o.ArchiveAt = &at
			o.ArchiveRule = "死单归档-家长明确拒绝"
		}
		// 国际游学商机
		if i >= 42 {
			o.IsIntlTravel = 1
			o.Category = "国际游学"
			o.OppType = "国际游学"
			o.IntentionProject = "游学营"
		}
		// 91 国游推送
		if i >= 36 && i < 42 {
			o.Category = "国内游学"
			o.Pushed91 = 1
			o.Push91Result = "推送成功"
			o.ResourceNo91 = fmt.Sprintf("91GY2610%05d", 2000+i)
			pt := created.AddDate(0, 0, 2)
			o.Push91At = &pt
			o.Follow91 = "91跟进顾问"
			o.Advancer91 = "91推进人"
		}
		initScore, score, diag := ComputeScore(&o)
		o.InitScore, o.Score, o.AiDiagnosis = initScore, score, diag
		o.MaxScore = score
		if i%9 == 0 && score > 50 {
			o.MaxScore = score + 6
		}
		o.Grade2 = GradeByScore(int(o.MaxScore))
		list = append(list, o)
	}
	return db.CreateInBatches(&list, 100).Error
}

func seedAssignRecords(db *gorm.DB) error {
	if countAny(db, &model.OppAssignRecord{}, "tenant_id = ?", CurrentTenant) > 0 {
		return nil
	}
	var opps []model.Opportunity
	if err := db.Where("tenant_id = ?", CurrentTenant).Limit(14).Find(&opps).Error; err != nil {
		return err
	}
	rows := make([]model.OppAssignRecord, 0, len(opps))
	now := time.Now()
	for i, o := range opps {
		from := seedEmps[(i+1)%len(seedEmps)]
		to := seedEmps[i%len(seedEmps)]
		auto := int32(0)
		if i%4 == 0 { auto = 1 }
		opAt := now.AddDate(0, 0, -(i * 2))
		rows = append(rows, model.OppAssignRecord{
			TenantID: CurrentTenant, OppID: o.ID, Code: o.Code, Name: o.Name, StudentNo: o.StudentNo,
			Grade: o.Grade, School: o.School, Remark: "工作再分配", Operator: from.name, OperatorDept: from.dept,
			AutoAssign: auto, OpAt: opAt, Receiver: to.name, ReceiverDept: to.dept,
			Status: o.Status, LastFollowAt: o.LastFollowAt,
		})
	}
	return db.CreateInBatches(&rows, 100).Error
}

func seedEmployeeTasks(db *gorm.DB) error {
	now := time.Now()

	// 旧线索可能只有任务名称没有 task_id，按任务名回填确定性 ID
	var taskNames []string
	db.Model(&model.Lead{}).Where("tenant_id = ? AND (task_id IS NULL OR task_id = 0) AND task_name <> ''", CurrentTenant).
		Distinct("task_name").Pluck("task_name", &taskNames)
	for i, name := range taskNames {
		db.Model(&model.Lead{}).Where("tenant_id = ? AND task_name = ? AND (task_id IS NULL OR task_id = 0)", CurrentTenant, name).
			Update("task_id", int64(5000+i))
	}

	rows := []model.EmployeeTask{}
	// 线索任务：从线索表按任务×归属人聚合（每次重建，保持快照新鲜）
	type grp struct {
		TaskID    int64
		TaskName  string
		OwnerID   string
		OwnerName string
		OwnerDept string
	}
	var grps []grp
	db.Raw(`SELECT task_id AS task_id, MAX(task_name) AS task_name, owner_id AS owner_id,
		MAX(owner_name) AS owner_name, MAX(owner_dept) AS owner_dept
		FROM scrm_lead WHERE tenant_id = ? AND task_id > 0 GROUP BY task_id, owner_id`, CurrentTenant).Scan(&grps)
	statusKeys := []string{
		model.StatusPendingAssign, model.StatusPendingFollow, model.StatusFollowing,
		model.StatusRecycled, model.StatusInvalidArchive, model.StatusStoredArchive,
		model.StatusConverted,
	}
	for i, g := range grps {
		counts := map[string]int{}
		var total int64
		db.Model(&model.Lead{}).Where("tenant_id = ? AND task_id = ? AND owner_id = ?", CurrentTenant, g.TaskID, g.OwnerID).Count(&total)
		for _, st := range statusKeys {
			var n int64
			db.Model(&model.Lead{}).Where("tenant_id = ? AND task_id = ? AND owner_id = ? AND status = ?",
				CurrentTenant, g.TaskID, g.OwnerID, st).Count(&n)
			counts[st] = int(n)
		}
		cb, _ := json.Marshal(counts)
		start := now.AddDate(0, 0, -60)
		end := now.AddDate(0, 0, 30)
		rows = append(rows, model.EmployeeTask{
			TenantID: CurrentTenant, BizKind: "lead", TaskID: g.TaskID, TaskName: g.TaskName,
			TaskType: []string{"线下地推", "线上推广", "转介绍"}[i%3],
			OwnerID: g.OwnerID, OwnerName: g.OwnerName, OwnerDept: g.OwnerDept,
			Supervisor: "王主管", StartAt: &start, EndAt: &end,
			AutoArchive: int32(i % 2), TotalCount: int32(total), StatusCounts: string(cb),
		})
	}
	// 线索任务快照每次重建
	if err := db.Where("tenant_id = ? AND biz_kind = ?", CurrentTenant, "lead").Delete(&model.EmployeeTask{}).Error; err != nil {
		return err
	}
	if len(rows) > 0 {
		if err := db.CreateInBatches(&rows, 200).Error; err != nil {
			return err
		}
	}

	// 商机任务：每个顾问 1~2 个（幂等）
	if countAny(db, &model.EmployeeTask{}, "tenant_id = ? AND biz_kind = ?", CurrentTenant, "opp") > 0 {
		return nil
	}
	oppRows := []model.EmployeeTask{}
	taskDefs := []struct{ name, typ string }{
		{"202610-高中-国际游学招生", "国际游学"},
		{"202610-高中-一对一冲刺", "线下地推"},
		{"202610-初中-体验课转化", "体验课"},
	}
	for i, emp := range seedEmps {
		for j, td := range taskDefs {
			if j > 0 && i%2 == 0 { continue }
			start := now.AddDate(0, 0, -40)
			end := now.AddDate(0, 0, 20)
			counts := map[string]int{
				model.OppStatusPendingAssign: j, model.OppStatusPendingFollow: 2 + i%3,
				model.OppStatusFollowing: 3 + j, model.OppStatusRecycled: i % 3,
				model.OppStatusDeadArchive: i % 2, model.OppStatusDealt: i % 4,
			}
			total := 0
			for _, v := range counts { total += v }
			cb, _ := json.Marshal(counts)
			oppRows = append(oppRows, model.EmployeeTask{
				TenantID: CurrentTenant, BizKind: "opp",
				TaskID: int64(9000 + i*10 + j), TaskName: td.name, TaskType: td.typ,
				OwnerID: emp.id, OwnerName: emp.name, OwnerDept: emp.dept,
				Supervisor: "王主管", StartAt: &start, EndAt: &end, TotalCount: int32(total), StatusCounts: string(cb),
			})
		}
	}
	return db.CreateInBatches(&oppRows, 200).Error
}

func seedAppeals(db *gorm.DB) error {
	if countAny(db, &model.Appeal{}, "tenant_id = ?", CurrentTenant) > 0 {
		return nil
	}
	bizLines := []struct{ line, certLabel string }{
		{"班课", "TKZ"}, {"高端", "HTB"}, {"会员卡", "HYK"}, {"课时包", "KSB"},
	}
	rows := []model.Appeal{}
	now := time.Now()
	for i := 0; i < 8; i++ {
		bl := bizLines[i%4]
		emp := seedEmps[i%len(seedEmps)]
		status := "待审批"
		if i%3 == 0 { status = "已通过" }
		if i%5 == 0 { status = "已驳回" }
		a := model.Appeal{
			TenantID: CurrentTenant, BizLine: bl.line,
			OppCode: fmt.Sprintf("XSJ%s%06d", now.Format("060102"), 100001+i*3),
			CertNo: fmt.Sprintf("%s%08d", bl.certLabel[:2], 88000000+i*7),
			OwnerName: emp.name, OwnerDept: emp.dept, Status: status,
			Reason: []string{"订单归属错误，实际为本人跟进成单", "试听后转正课业绩未计入", "续费业绩归属调整"}[i%3],
			AdjustMode: []string{"调整归属", "业绩拆分"}[i%2],
			Applicant: emp.name, ApplicantDept: emp.dept,
			CreatedAt: now.AddDate(0, 0, -(i * 3)),
		}
		if status != "待审批" {
			ht := a.CreatedAt.AddDate(0, 0, 1)
			a.HandledAt = &ht
			a.HandledBy = "王主管"
		}
		rows = append(rows, a)
	}
	return db.CreateInBatches(&rows, 100).Error
}

func seedDealApprovals(db *gorm.DB) error {
	if countAny(db, &model.DealApproval{}, "tenant_id = ?", CurrentTenant) > 0 {
		return nil
	}
	bizTypes := []string{"班课", "高端", "会员卡", "课时包"}
	now := time.Now()
	rows := []model.DealApproval{}
	for i := 0; i < 7; i++ {
		emp := seedEmps[i%len(seedEmps)]
		status := "待审核"
		if i%3 == 0 { status = "已通过" }
		a := model.DealApproval{
			TenantID: CurrentTenant,
			OppCode: fmt.Sprintf("XSJ%s%06d", now.Format("060102"), 100005+i*4),
			Status: status, Applicant: emp.name, ApplicantDept: emp.dept,
			BizType: bizTypes[i%4], CertNo: fmt.Sprintf("TKZ%08d", 77000000+i*11),
			Amount: fmt.Sprintf("%d.00", 3800+i*2600),
			Reason: "家长已缴费，申请成单审批",
			OrderCreator: emp.name, CreatedAt: now.AddDate(0, 0, -(i * 2)),
		}
		if status != "待审核" {
			ht := a.CreatedAt.AddDate(0, 0, 1)
			a.HandledAt = &ht
			a.Approver = "王主管"
		}
		rows = append(rows, a)
	}
	return db.CreateInBatches(&rows, 100).Error
}

func seedDealLogs(db *gorm.DB) error {
	if countAny(db, &model.DealLog{}, "tenant_id = ?", CurrentTenant) > 0 {
		return nil
	}
	now := time.Now()
	bizLines := []string{"班课", "高端", "会员卡", "课时包", "国游"}
	products := map[string]string{
		"班课": "高二数学冲刺班", "高端": "高中全托VIP方案", "会员卡": "学习一卡通银卡",
		"课时包": "60课时一对一包", "国游": "新加坡国立大学研学营",
	}
	rows := []model.DealLog{}
	for i := 0; i < 18; i++ {
		bl := bizLines[i%5]
		emp := seedEmps[i%len(seedEmps)]
		pay := now.AddDate(0, 0, -(i * 6))
		match := pay.Add(time.Duration(2+i*3) * time.Hour)
		rows = append(rows, model.DealLog{
			TenantID: CurrentTenant, BizLine: bl,
			OppCode: fmt.Sprintf("XSJ%s%06d", now.Format("060102"), 100001+i*2),
			OwnerName: emp.name, OwnerDept: emp.dept, OrderPerson: emp.name,
			StudentNo: fmt.Sprintf("XDF2026%04d", 3000+i),
			StudentName: []string{"李明轩", "王思远", "张子涵", "刘宇航", "陈雨桐"}[i%5],
			ProductName: products[bl],
			Amount: fmt.Sprintf("%d.00", 2800+(i%6)*3200),
			CertNo: fmt.Sprintf("TKZ%08d", 66000000+i*13),
			OrderNo: fmt.Sprintf("DD2026%07d", 9000000+i*17),
			DealType: []string{"新签", "续费", "扩科"}[i%3],
			PayAt: &pay, MatchAt: &match, CreatedAt: match,
		})
	}
	return db.CreateInBatches(&rows, 100).Error
}

func seedTickets(db *gorm.DB) error {
	if countAny(db, &model.BizTicket{}, "tenant_id = ?", CurrentTenant) > 0 {
		return nil
	}
	now := time.Now()
	bizTypes := []string{"退班", "调班", "转班"}
	statuses := []string{"待审核", "待客服处理", "已完成", "已驳回"}
	rows := []model.BizTicket{}
	for i := 0; i < 12; i++ {
		submit := now
		switch {
		case i%4 == 0:
			submit = now.Add(-26 * time.Hour) // 超24H客服未处理
		case i%4 == 1:
			submit = now.Add(-8 * 24 * time.Hour) // 超7天未完成
		default:
			submit = now.AddDate(0, 0, -(i % 5))
		}
		st := statuses[i%4]
		if i%4 == 0 { st = "待客服处理" }
		rows = append(rows, model.BizTicket{
			TenantID: CurrentTenant, Code: fmt.Sprintf("GD2026%07d", 5000000+i*19),
			SubmitAccount: "chenyuankang", SubmitName: "陈元康", SubmitAt: submit,
			BizType: bizTypes[i%3],
			StudentName: []string{"李明轩", "王思远", "张子涵"}[i%3],
			StudentNo: fmt.Sprintf("XDF2026%04d", 3000+i),
			Status: st, Source: []string{"积分商城", "POS", "线上自助"}[i%3],
			CsAccount: []string{"", "kefu01", "kefu02"}[i%3],
			CreatedAt: submit,
		})
	}
	return db.CreateInBatches(&rows, 100).Error
}

func seedMallOrders(db *gorm.DB) error {
	if countAny(db, &model.MallOrder{}, "tenant_id = ?", CurrentTenant) > 0 {
		return nil
	}
	now := time.Now()
	products := []struct{ typ, code, name, sale string }{
		{"课程商品", "KC001", "高二数学秋季系统班", "新签"},
		{"实物商品", "SW102", "高考真题礼盒", "零售"},
		{"增值服务", "FW203", "一对一学情诊断", "赠送"},
		{"会员卡", "HY304", "学习一卡通金卡", "新签"},
	}
	rows := []model.MallOrder{}
	for i := 0; i < 15; i++ {
		p := products[i%4]
		rows = append(rows, model.MallOrder{
			TenantID: CurrentTenant,
			StudentCode: fmt.Sprintf("XDF2026%04d", 3000+i),
			StudentName: []string{"李明轩", "王思远", "张子涵", "刘宇航"}[i%4],
			ProductType: p.typ, SaleType: p.sale, ProductCode: p.code, ProductName: p.name,
			OrderNo: fmt.Sprintf("JF2026%08d", 40000000+i*23),
			BatchNo: fmt.Sprintf("PC2610%03d", i+1),
			BizType: []string{"新签", "续费"}[i%2],
			Status: []string{"已完成", "已发货", "待发货", "已取消"}[i%4],
			OrderAccount: "chenyuankang", FlagDept: "高中部-江汉组",
			CreatedAt: now.AddDate(0, 0, -(i * 2)),
		})
	}
	return db.CreateInBatches(&rows, 100).Error
}
