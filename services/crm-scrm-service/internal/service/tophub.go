package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/model"
	"github.com/oa-portal/crm-scrm-service/internal/repo"
)

// ================= 统一搜索（跨线索/学员/商机） =================

type SearchResponse struct {
	Keyword      string             `json:"keyword"`
	Leads        []model.Lead       `json:"leads"`
	Students     []StudentHit       `json:"students"`
	Opportunities []OpportunityStub `json:"opportunities"`
}

type StudentHit struct {
	LeadID    int64  `json:"lead_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	StudentNo string `json:"student_no"`
	School    string `json:"school"`
	Grade     string `json:"grade"`
}

// OpportunityStub 第三期商机管理上线后接入真实查询
type OpportunityStub struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
}

func (s *ScrmService) SearchAll(ctx context.Context, keyword string) (*SearchResponse, error) {
	leads, err := s.leadRepo.Search(ctx, CurrentTenant, keyword, 20)
	if err != nil {
		return nil, err
	}
	res := &SearchResponse{Keyword: keyword, Leads: leads, Opportunities: []OpportunityStub{}}
	seen := map[string]bool{}
	for _, l := range leads {
		if l.StudentNo != "" && !seen[l.StudentNo] {
			seen[l.StudentNo] = true
			res.Students = append(res.Students, StudentHit{
				LeadID: l.ID, Name: l.Name, Phone: l.Phone, StudentNo: l.StudentNo,
				School: l.School, Grade: l.Grade,
			})
		}
	}
	return res, nil
}

func (s *ScrmService) GetLead(ctx context.Context, id int64) (*model.Lead, error) {
	return s.leadRepo.Get(ctx, id)
}

// ================= 云座席双呼 =================

// Dial 点击「双呼上号」即落一条跟进日志（拨号事件），返回日志ID
func (s *ScrmService) Dial(ctx context.Context, leadID int64, phone string) (*model.FollowLog, error) {
	if leadID <= 0 {
		return nil, fmt.Errorf("请先选择要外呼的线索")
	}
	lead, err := s.leadRepo.Get(ctx, leadID)
	if err != nil {
		return nil, fmt.Errorf("线索不存在: %w", err)
	}
	if phone == "" {
		phone = lead.Phone
	}
	if !validPhone(phone) {
		return nil, fmt.Errorf("客户手机号无效，无法外呼")
	}
	log := &model.FollowLog{
		TenantID:      CurrentTenant,
		LeadID:        lead.ID,
		LeadName:      lead.Name,
		Type:          "双呼",
		Content:       fmt.Sprintf("发起双呼：系统同时呼起顾问手机与客户手机（%s），等待接通…", phone),
		CustomerPhone: phone,
		Operator:      CurrentUserName,
		OperatorID:    CurrentUserID,
	}
	if err := s.logRepo.Create(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// CompleteCall 通话结束：更新日志结果/时长，并回写线索通话计数与接线状态
func (s *ScrmService) CompleteCall(ctx context.Context, logID int64, connected bool, durationSec int32) error {
	if durationSec < 0 {
		durationSec = 0
	}
	logEntry, err := s.logRepo.Get(ctx, logID)
	if err != nil {
		return fmt.Errorf("通话记录不存在: %w", err)
	}
	now := time.Now()
	if connected {
		logEntry.Result = "已接通"
		logEntry.Content = fmt.Sprintf("双呼接通，通话时长 %s", fmtDuration(durationSec))
	} else {
		logEntry.Result = "未接通"
		logEntry.Content = "双呼发起成功，客户未接通"
	}
	logEntry.DurationSec = durationSec
	if err := s.logRepo.Update(ctx, logEntry); err != nil {
		return err
	}
	return s.leadRepo.ApplyCallResult(ctx, CurrentTenant, logEntry.LeadID, connected, now)
}

func (s *ScrmService) ListFollowLogs(ctx context.Context, leadID int64) ([]model.FollowLog, error) {
	return s.logRepo.ListByLead(ctx, CurrentTenant, leadID, 50)
}

func (s *ScrmService) RecentCalls(ctx context.Context) ([]model.FollowLog, error) {
	return s.logRepo.ListRecent(ctx, CurrentTenant, CurrentUserID, 20)
}

func fmtDuration(sec int32) string {
	m := sec / 60
	return fmt.Sprintf("%02d:%02d", m, sec%60)
}

// ================= 站内通知 =================

func (s *ScrmService) ListNotifications(ctx context.Context, unread string, p repo.Page) (*repo.PageResult[model.Notification], error) {
	return s.notifyRepo.List(ctx, CurrentTenant, repo.NotificationFilter{UserID: CurrentUserID, Unread: unread}, p)
}

func (s *ScrmService) UnreadNotifications(ctx context.Context) (int64, error) {
	return s.notifyRepo.UnreadCount(ctx, CurrentTenant, CurrentUserID)
}

func (s *ScrmService) MarkNotificationsRead(ctx context.Context, ids []int64, all bool) error {
	return s.notifyRepo.MarkRead(ctx, CurrentTenant, CurrentUserID, ids, all)
}
