package repo

import (
	"time"

	"gorm.io/gorm"

	"github.com/oa-portal/portal-bff/internal/model"
)

// PortalRepo 门户业务数据
type PortalRepo struct {
	db *gorm.DB
}

func NewPortalRepo(db *gorm.DB) *PortalRepo {
	return &PortalRepo{db: db}
}

// ---------- 审批 ----------

func (r *PortalRepo) ListApprovals(typ, status, keyword, currentUserID string) ([]model.PortalApproval, error) {
	q := r.db.Model(&model.PortalApproval{}).Order("created_at DESC")

	if status == "mine" {
		q = q.Where("applicant_id = ?", currentUserID)
	} else if status == "done" {
		q = q.Where("status = ?", 2)
	} else {
		q = q.Where("status = ?", 1)
	}

	if typ != "" && typ != "全部" {
		q = q.Where("type = ?", typ)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR serial_no LIKE ?", like, like)
	}

	var list []model.PortalApproval
	err := q.Limit(50).Find(&list).Error
	return list, err
}

// ApprovalCounts 各状态审批数量
func (r *PortalRepo) ApprovalCounts(typ, currentUserID string) (map[string]int64, error) {
	base := r.db.Model(&model.PortalApproval{})
	if typ != "" && typ != "全部" {
		base = base.Where("type = ?", typ)
	}

	var todo, done, mine int64
	if err := base.Session(&gorm.Session{}).Where("status = ?", 1).Count(&todo).Error; err != nil {
		return nil, err
	}
	if err := base.Session(&gorm.Session{}).Where("status = ?", 2).Count(&done).Error; err != nil {
		return nil, err
	}
	if err := base.Session(&gorm.Session{}).Where("applicant_id = ?", currentUserID).Count(&mine).Error; err != nil {
		return nil, err
	}
	return map[string]int64{"todo": todo, "done": done, "mine": mine, "all": todo + done}, nil
}

func (r *PortalRepo) CreateApproval(a *model.PortalApproval) error {
	return r.db.Create(a).Error
}

// ---------- 日程 ----------

func (r *PortalRepo) ListSchedules(userID, month string) ([]model.PortalSchedule, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, err
	}
	end := start.AddDate(0, 1, 0)
	var list []model.PortalSchedule
	err = r.db.Where("user_id = ? AND event_date >= ? AND event_date < ?", userID, start, end).
		Order("event_date ASC").Find(&list).Error
	return list, err
}

// ---------- 常用 ----------

func (r *PortalRepo) ListFavorites(userID string) ([]model.PortalFavorite, error) {
	var list []model.PortalFavorite
	err := r.db.Where("user_id = ?", userID).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *PortalRepo) AddFavorite(f *model.PortalFavorite) error {
	return r.db.Create(f).Error
}

func (r *PortalRepo) DeleteFavorite(id int64, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.PortalFavorite{}).Error
}

// EmployeeRepo 员工只读查询（iam 库）
type EmployeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) *EmployeeRepo {
	return &EmployeeRepo{db: db}
}

func (r *EmployeeRepo) Search(keyword string, limit int) ([]model.IAMUser, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	q := r.db.Model(&model.IAMUser{}).Where("status = ?", 1)
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("real_name LIKE ? OR username LIKE ? OR email LIKE ?", like, like, like)
	}
	var list []model.IAMUser
	err := q.Limit(limit).Find(&list).Error
	return list, err
}
