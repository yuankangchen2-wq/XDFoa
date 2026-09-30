package repo

import (
	"context"

	"github.com/oa-portal/mdm-service/internal/model"
	"gorm.io/gorm"
)

// ---------- 商品 ----------
type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *model.Product) error {
	return GetDB(ctx, r.db).Create(p).Error
}

func (r *ProductRepo) Update(ctx context.Context, p *model.Product) error {
	return GetDB(ctx, r.db).Save(p).Error
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64, tenantID string) (*model.Product, error) {
	var p model.Product
	err := GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) List(ctx context.Context, tenantID, keyword, category string, pq PageQuery) (*PageResult[model.Product], error) {
	q := GetDB(ctx, r.db).Model(&model.Product{}).Where("tenant_id = ?", tenantID)
	if keyword != "" {
		q = q.Where("product_name LIKE ? OR sku_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	return Paginate[model.Product](r.db, q, pq)
}

func (r *ProductRepo) Delete(ctx context.Context, id int64, tenantID string) error {
	return GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Product{}).Error
}

// ---------- 供应商 ----------
type SupplierRepo struct {
	db *gorm.DB
}

func NewSupplierRepo(db *gorm.DB) *SupplierRepo {
	return &SupplierRepo{db: db}
}

func (r *SupplierRepo) Create(ctx context.Context, s *model.Supplier) error {
	return GetDB(ctx, r.db).Create(s).Error
}

func (r *SupplierRepo) Update(ctx context.Context, s *model.Supplier) error {
	return GetDB(ctx, r.db).Save(s).Error
}

func (r *SupplierRepo) GetByID(ctx context.Context, id int64, tenantID string) (*model.Supplier, error) {
	var s model.Supplier
	err := GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SupplierRepo) List(ctx context.Context, tenantID, keyword string, pq PageQuery) (*PageResult[model.Supplier], error) {
	q := GetDB(ctx, r.db).Model(&model.Supplier{}).Where("tenant_id = ?", tenantID)
	if keyword != "" {
		q = q.Where("supplier_name LIKE ? OR supplier_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	return Paginate[model.Supplier](r.db, q, pq)
}

func (r *SupplierRepo) Delete(ctx context.Context, id int64, tenantID string) error {
	return GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Supplier{}).Error
}

// ---------- 组织 ----------
type OrganizationRepo struct {
	db *gorm.DB
}

func NewOrganizationRepo(db *gorm.DB) *OrganizationRepo {
	return &OrganizationRepo{db: db}
}

func (r *OrganizationRepo) Create(ctx context.Context, o *model.Organization) error {
	return GetDB(ctx, r.db).Create(o).Error
}

func (r *OrganizationRepo) Update(ctx context.Context, o *model.Organization) error {
	return GetDB(ctx, r.db).Save(o).Error
}

func (r *OrganizationRepo) GetByID(ctx context.Context, id int64, tenantID string) (*model.Organization, error) {
	var o model.Organization
	err := GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrganizationRepo) ListByParent(ctx context.Context, tenantID string, parentID int64) ([]model.Organization, error) {
	var list []model.Organization
	err := GetDB(ctx, r.db).Where("tenant_id = ? AND parent_id = ?", tenantID, parentID).
		Order("sort_order ASC").Find(&list).Error
	return list, err
}

func (r *OrganizationRepo) GetByCode(ctx context.Context, tenantID, orgCode string) (*model.Organization, error) {
	var o model.Organization
	err := GetDB(ctx, r.db).Where("tenant_id = ? AND org_code = ?", tenantID, orgCode).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrganizationRepo) Delete(ctx context.Context, id int64, tenantID string) error {
	return GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Organization{}).Error
}
