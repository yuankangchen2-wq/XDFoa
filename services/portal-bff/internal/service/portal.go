package service

import (
	"time"

	"github.com/oa-portal/portal-bff/internal/model"
	"github.com/oa-portal/portal-bff/internal/repo"
)

// 当前登录用户（门户简化版：固定员工 陈元康）
const CurrentUserID = "1"

type PortalService struct {
	portal *repo.PortalRepo
	emp    *repo.EmployeeRepo
}

func NewPortalService(portal *repo.PortalRepo, emp *repo.EmployeeRepo) *PortalService {
	return &PortalService{portal: portal, emp: emp}
}

// ---------- 审批 ----------

func (s *PortalService) ListApprovals(typ, status, keyword string) ([]model.PortalApproval, error) {
	return s.portal.ListApprovals(typ, status, keyword, CurrentUserID)
}

func (s *PortalService) ApprovalCounts(typ string) (map[string]int64, error) {
	return s.portal.ApprovalCounts(typ, CurrentUserID)
}

func (s *PortalService) CreateApproval(a *model.PortalApproval) error {
	if a.Applicant == "" {
		a.Applicant = "陈元康"
	}
	if a.SerialNo == "" {
		a.SerialNo = "SP" + time.Now().Format("20060102150405")
	}
	if a.CurrentNode == "" {
		a.CurrentNode = "直属领导审批"
	}
	return s.portal.CreateApproval(a)
}

// ---------- 日程 ----------

func (s *PortalService) ListSchedules(month string) ([]model.PortalSchedule, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	return s.portal.ListSchedules(CurrentUserID, month)
}

// ---------- 常用 ----------

func (s *PortalService) ListFavorites() ([]model.PortalFavorite, error) {
	return s.portal.ListFavorites(CurrentUserID)
}

func (s *PortalService) AddFavorite(name, icon, path string) (*model.PortalFavorite, error) {
	f := &model.PortalFavorite{
		UserID: CurrentUserID,
		Name:   name,
		Icon:   icon,
		Path:   path,
	}
	if f.Icon == "" {
		f.Icon = "Star"
	}
	if err := s.portal.AddFavorite(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *PortalService) DeleteFavorite(id int64) error {
	return s.portal.DeleteFavorite(id, CurrentUserID)
}

// ---------- 员工搜索 ----------

func (s *PortalService) SearchEmployees(keyword string) ([]model.IAMUser, error) {
	return s.emp.Search(keyword, 10)
}
