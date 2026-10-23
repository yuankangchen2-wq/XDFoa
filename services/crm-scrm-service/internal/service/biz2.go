package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/model"
	"github.com/oa-portal/crm-scrm-service/internal/repo"
)

// ==================== 线索工作台 ====================

func (s *ScrmService) RecycleLeads(ctx context.Context, ids []int64, reason string) error {
	if len(ids) == 0 {
		return fmt.Errorf("请先选择线索")
	}
	if reason == "" {
		return fmt.Errorf("请选择回收原因")
	}
	if err := s.leadRepo.Recycle(ctx, CurrentTenant, ids, reason, CurrentUserName, time.Now()); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "lead", BizIDs: joinIDs(ids),
		Action: "线索回收", Operator: CurrentUserName, Detail: "回收原因：" + reason,
	})
}

func (s *ScrmService) ArchiveLeads(ctx context.Context, ids []int64, way, rule, reason string) error {
	if len(ids) == 0 {
		return fmt.Errorf("请先选择线索")
	}
	return s.leadRepo.BatchArchiveWithReason(ctx, CurrentTenant, ids, way, rule, reason, CurrentUserName, time.Now())
}

func (s *ScrmService) LeadTaskSnapshot(ctx context.Context, scope string) ([]repo.TaskSnapshotRow, error) {
	return s.leadRepo.TaskSnapshot(ctx, CurrentTenant, scope, CurrentUserID)
}

// ==================== 商机评分引擎 ====================

// GradeByScore 按得分评级：S≥85 / A≥70 / B≥55 / C 其余
func GradeByScore(score int) string {
	switch {
	case score >= 85:
		return "S"
	case score >= 70:
		return "A"
	case score >= 55:
		return "B"
	default:
		return "C"
	}
}

// ComputeScore 初始分+行为分的确定性评分（演示口径，规则见 ScoreRuleText）
func ComputeScore(o *model.Opportunity) (initScore, score int32, diagnosis string) {
	base := 40
	switch o.Level {
	case "高":
		base += 22
	case "中":
		base += 12
	case "低":
		base += 4
	}
	if o.ChannelL1 == "市场地推" {
		base += 8
	}
	if o.AddWechat == 1 {
		base += 6
	}
	if o.NewOldType == "老学员" {
		base += 5
	}
	initScore = int32(math.Max(20, math.Min(95, float64(base))))
	score = initScore + o.FollowCount*2
	if o.LastFollowContent != "" {
		score += 3
	}
	if o.PromiseVisitAt != nil {
		score += 6
	}
	if o.VisitConverted == 1 {
		score += 8
	}
	if o.ContactPoint == "呼入" {
		score += 4
	}
	if score > 100 {
		score = 100
	}
	switch {
	case score >= 85:
		diagnosis = "高意向高活跃，建议优先推进到访与成单"
	case score >= 70:
		diagnosis = "意向明确，保持高频跟进争取本周到访"
	case score >= 55:
		diagnosis = "意向一般，建议通过体验课激活"
	default:
		diagnosis = "暂无诊断"
	}
	return
}

// ScoreRuleText 表头 ② 悬浮说明
func ScoreRuleText() map[string]string {
	return map[string]string{
		"level":     "商机等级：按得分自动评级 S(≥85) / A(≥70) / B(≥55) / C(<55)",
		"score":     "商机分数：初始分 + 跟进×2 + 跟进记录3 + 诺访6 + 到访转化8 + 呼入4 + 加微6",
		"init":      "初始分数：创建时按意向等级、渠道、新老生、加微情况计算（20~95）",
		"max":       "历史最高分：得分历次变化的峰值，只升不降",
		"ai":        "AI诊断：根据分数区间生成质量诊断与推进建议",
	}
}

// ==================== 商机动作 ====================

func (s *ScrmService) CreateOpportunity(ctx context.Context, o *model.Opportunity) error {
	o.TenantID = CurrentTenant
	if o.Code == "" {
		o.Code = fmt.Sprintf("XSJ%s%06d", time.Now().Format("060102"), time.Now().UnixNano()%1e6)
	}
	if o.Creator == "" {
		o.Creator = CurrentUserName
	}
	if o.CreateMethod == "" {
		o.CreateMethod = "搜索手工创建"
	}
	if o.EnterType == "" {
		o.EnterType = "直接创建商机"
	}
	if o.Status == "" {
		o.Status = model.OppStatusPendingFollow
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	initScore, score, diag := ComputeScore(o)
	o.InitScore, o.Score, o.AiDiagnosis = initScore, score, diag
	o.MaxScore = score
	o.Grade2 = GradeByScore(int(score))
	return s.oppRepo.Create(ctx, o)
}

func (s *ScrmService) ListOpportunities(ctx context.Context, scope string, f repo.OppFilter, p repo.Page) (*repo.PageResult[model.Opportunity], error) {
	return s.oppRepo.List(ctx, CurrentTenant, scope, CurrentUserID, f, p)
}

// UpdateOpportunityFields 跟进计划等局部字段更新（白名单）
func (s *ScrmService) UpdateOpportunityFields(ctx context.Context, id int64, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return fmt.Errorf("没有可更新的字段")
	}
	allowed := map[string]bool{
		"next_follow_at": true, "promise_visit_at": true, "remark": true,
		"last_follow_content": true, "owner_name": true, "owner_id": true, "owner_dept": true,
		"level": true, "add_wechat": true,
	}
	updates := map[string]interface{}{}
	for k, v := range fields {
		if allowed[k] {
			updates[k] = v
		}
	}
	if len(updates) == 0 {
		return fmt.Errorf("字段不允许更新")
	}
	updates["updated_at"] = time.Now()
	return s.db.WithContext(ctx).Model(&model.Opportunity{}).
		Where("tenant_id = ? AND id = ?", CurrentTenant, id).Updates(updates).Error
}

func (s *ScrmService) AssignOpportunities(ctx context.Context, ids []int64, receiverID, receiverName, receiverDept string, auto int32) error {
	if len(ids) == 0 || receiverName == "" {
		return fmt.Errorf("请选择商机与接收人")
	}
	now := time.Now()
	if err := s.oppRepo.Assign(ctx, CurrentTenant, ids, receiverID, receiverName, receiverDept, auto, now); err != nil {
		return err
	}
	// 写分配流水
	rows := []model.OppAssignRecord{}
	for _, id := range ids {
		o, err := s.oppRepo.Get(ctx, id)
		if err != nil {
			continue
		}
		rows = append(rows, model.OppAssignRecord{
			TenantID: CurrentTenant, OppID: o.ID, Code: o.Code, Name: o.Name, StudentNo: o.StudentNo,
			Grade: o.Grade, School: o.School, Remark: o.Remark, Operator: CurrentUserName,
			AutoAssign: auto, OpAt: now, Receiver: receiverName, ReceiverDept: receiverDept,
			Status: o.Status, LastFollowAt: o.LastFollowAt,
		})
	}
	return s.assignRepo.CreateBatch(ctx, rows)
}

func (s *ScrmService) RecycleOpportunities(ctx context.Context, ids []int64, reason string) error {
	if len(ids) == 0 {
		return fmt.Errorf("请先选择商机")
	}
	return s.oppRepo.Recycle(ctx, CurrentTenant, ids, reason, time.Now())
}

func (s *ScrmService) ArchiveOpportunities(ctx context.Context, ids []int64, rule string) error {
	if len(ids) == 0 {
		return fmt.Errorf("请先选择商机")
	}
	return s.oppRepo.BatchArchive(ctx, CurrentTenant, ids, rule, time.Now())
}

func (s *ScrmService) ActivateOpportunities(ctx context.Context, ids []int64) error {
	return s.oppRepo.Activate(ctx, CurrentTenant, ids)
}

func (s *ScrmService) Push91(ctx context.Context, ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("请先选择可连线商机")
	}
	now := time.Now()
	n := 0
	for _, id := range ids {
		resNo := fmt.Sprintf("91GY%s%05d", now.Format("060102"), id)
		if err := s.oppRepo.MarkPush91(ctx, CurrentTenant, []int64{id}, "推送成功", resNo, "91跟进顾问", "91推进人", now); err != nil {
			continue
		}
		n++
	}
	return n, nil
}

func (s *ScrmService) PushHedgehog(ctx context.Context, ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("请先选择商机")
	}
	return len(ids), s.oppRepo.MarkPushHedgehog(ctx, CurrentTenant, ids)
}

func (s *ScrmService) ListAssignRecords(ctx context.Context, scope string, f repo.AssignFilter, p repo.Page) (*repo.PageResult[model.OppAssignRecord], error) {
	return s.assignRepo.List(ctx, CurrentTenant, scope, CurrentUserID, f, p)
}

func (s *ScrmService) ListEmployeeTasks(ctx context.Context, f repo.EmployeeTaskFilter, p repo.Page) (*repo.PageResult[model.EmployeeTask], error) {
	return s.task2Repo.List(ctx, CurrentTenant, f, p)
}

// ==================== 申请成单 / 申诉 / 审批 ====================

func (s *ScrmService) ApplyDeal(ctx context.Context, a *model.DealApproval) error {
	if a.OppCode == "" || a.BizType == "" {
		return fmt.Errorf("商机编码与业务类型必填")
	}
	a.TenantID = CurrentTenant
	a.Status = "待审核"
	a.Applicant = CurrentUserName
	a.CreatedAt = time.Now()
	return s.approvalRepo.Create(ctx, a)
}

func (s *ScrmService) CreateAppeal(ctx context.Context, a *model.Appeal) error {
	if a.OppCode == "" || a.BizLine == "" {
		return fmt.Errorf("商机编码必填")
	}
	a.TenantID = CurrentTenant
	a.Status = "待审批"
	a.Applicant = CurrentUserName
	a.CreatedAt = time.Now()
	return s.appealRepo.Create(ctx, a)
}

func (s *ScrmService) ListAppeals(ctx context.Context, bizLine, keyword, status string, p repo.Page) (*repo.PageResult[model.Appeal], error) {
	return s.appealRepo.List(ctx, CurrentTenant, bizLine, keyword, status, p)
}

func (s *ScrmService) ListApprovals(ctx context.Context, keyword, status, applicant, bizType string, p repo.Page) (*repo.PageResult[model.DealApproval], error) {
	return s.approvalRepo.List(ctx, CurrentTenant, keyword, status, applicant, bizType, p)
}

// HandleApprovals 审批成单申请；通过时同步生成成单日志并把商机置为已成单
func (s *ScrmService) HandleApprovals(ctx context.Context, ids []int64, pass bool) error {
	now := time.Now()
	if pass {
		for _, id := range ids {
			var a model.DealApproval
			if err := s.db.First(&a, id).Error; err != nil {
				continue
			}
			dl := &model.DealLog{
				TenantID: CurrentTenant, BizLine: a.BizType, OppCode: a.OppCode, OwnerName: a.Applicant,
				OrderPerson: a.OrderCreator, CertNo: a.CertNo, Amount: a.Amount,
				PayAt: &now, MatchAt: &now, CreatedAt: now,
				ProductName: a.BizType + "成单",
			}
			_ = s.dealRepo.Create(ctx, dl)
			s.db.Model(&model.Opportunity{}).Where("tenant_id = ? AND code = ?", CurrentTenant, a.OppCode).
				Updates(map[string]interface{}{"status": model.OppStatusDealt, "deal_at": now})
		}
	}
	return s.approvalRepo.Handle(ctx, CurrentTenant, ids, pass, CurrentUserName, now)
}

func (s *ScrmService) HandleAppeals(ctx context.Context, ids []int64, pass bool) error {
	return s.appealRepo.Handle(ctx, CurrentTenant, ids, pass, CurrentUserName, time.Now())
}

// ==================== 成单/工单/订单/列配置 ====================

func (s *ScrmService) ListDeals(ctx context.Context, f repo.DealFilter, p repo.Page) (*repo.PageResult[model.DealLog], error) {
	return s.dealRepo.List(ctx, CurrentTenant, f, p)
}

func (s *ScrmService) ListTickets(ctx context.Context, scope, status, bizType, keyword string, p repo.Page) (*repo.PageResult[model.BizTicket], error) {
	return s.ticketRepo.List(ctx, CurrentTenant, scope, status, bizType, keyword, p)
}

func (s *ScrmService) ListOrders(ctx context.Context, f repo.OrderFilter, p repo.Page) (*repo.PageResult[model.MallOrder], error) {
	return s.orderRepo.List(ctx, CurrentTenant, f, p)
}

func (s *ScrmService) GetGridConfig(ctx context.Context, pageKey, cfgType string) (string, error) {
	cfg, err := s.gridRepo.Get(ctx, CurrentTenant, CurrentUserID, pageKey, cfgType)
	if err != nil {
		return "", nil
	}
	return cfg, nil
}

func (s *ScrmService) SaveGridConfig(ctx context.Context, pageKey, cfgType, config string) error {
	return s.gridRepo.Upsert(ctx, CurrentTenant, CurrentUserID, pageKey, cfgType, config)
}
