package repo

import (
	"context"

	"github.com/oa-portal/mdm-service/internal/model"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) Create(ctx context.Context, c *model.Customer) error {
	return GetDB(ctx, r.db).Create(c).Error
}

func (r *CustomerRepo) Update(ctx context.Context, c *model.Customer) error {
	return GetDB(ctx, r.db).Save(c).Error
}

func (r *CustomerRepo) GetByID(ctx context.Context, id int64, tenantID string) (*model.Customer, error) {
	var c model.Customer
	err := GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) List(ctx context.Context, tenantID, keyword, customerType string, pq PageQuery) (*PageResult[model.Customer], error) {
	q := GetDB(ctx, r.db).Model(&model.Customer{}).Where("tenant_id = ?", tenantID)
	if keyword != "" {
		q = q.Where("customer_name LIKE ? OR customer_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if customerType != "" {
		q = q.Where("customer_type = ?", customerType)
	}
	return Paginate[model.Customer](r.db, q, pq)
}

func (r *CustomerRepo) Delete(ctx context.Context, id int64, tenantID string) error {
	return GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Customer{}).Error
}
