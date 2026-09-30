package repo

import (
	"context"

	"github.com/oa-portal/inventory-service/internal/model"
	"gorm.io/gorm"
)

// ---------- 通用分页 ----------
type Page struct {
	Page     int
	PageSize int
}

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func Paginate(p Page) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.Page <= 0 {
			p.Page = 1
		}
		if p.PageSize <= 0 || p.PageSize > 100 {
			p.PageSize = 20
		}
		return db.Offset((p.Page - 1) * p.PageSize).Limit(p.PageSize)
	}
}

// ---------- 仓库 ----------
type WarehouseRepo struct{ db *gorm.DB }

func NewWarehouseRepo(db *gorm.DB) *WarehouseRepo { return &WarehouseRepo{db: db} }

func (r *WarehouseRepo) Create(ctx context.Context, w *model.Warehouse) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *WarehouseRepo) List(ctx context.Context, tenantID string, p Page) (*PageResult[model.Warehouse], error) {
	var result PageResult[model.Warehouse]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if err := q.Model(&model.Warehouse{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *WarehouseRepo) Get(ctx context.Context, id int64) (*model.Warehouse, error) {
	var w model.Warehouse
	err := r.db.WithContext(ctx).First(&w, id).Error
	return &w, err
}

// ---------- 库存账户 ----------
type StockAccountRepo struct{ db *gorm.DB }

func NewStockAccountRepo(db *gorm.DB) *StockAccountRepo { return &StockAccountRepo{db: db} }

func (r *StockAccountRepo) DB() *gorm.DB { return r.db }

// GetOrCreate 获取或创建库存账户（按 sku + warehouse）
func (r *StockAccountRepo) GetOrCreate(ctx context.Context, tenantID, skuCode string, warehouseID int64) (*model.StockAccount, error) {
	var acc model.StockAccount
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND sku_code = ? AND warehouse_id = ?", tenantID, skuCode, warehouseID).
		First(&acc).Error
	if err == nil {
		return &acc, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	acc = model.StockAccount{
		TenantID:    tenantID,
		SkuCode:     skuCode,
		WarehouseID: warehouseID,
		Quantity:    0,
		Allocated:   0,
	}
	if err := r.db.WithContext(ctx).Create(&acc).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *StockAccountRepo) List(ctx context.Context, tenantID string, skuCode string, warehouseID int64, p Page) (*PageResult[model.StockAccount], error) {
	var result PageResult[model.StockAccount]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if skuCode != "" {
		q = q.Where("sku_code LIKE ?", "%"+skuCode+"%")
	}
	if warehouseID > 0 {
		q = q.Where("warehouse_id = ?", warehouseID)
	}
	if err := q.Model(&model.StockAccount{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

// ---------- 库存流水 ----------
type StockJournalRepo struct{ db *gorm.DB }

func NewStockJournalRepo(db *gorm.DB) *StockJournalRepo { return &StockJournalRepo{db: db} }

func (r *StockJournalRepo) Create(ctx context.Context, j *model.StockJournal) error {
	return r.db.WithContext(ctx).Create(j).Error
}

func (r *StockJournalRepo) List(ctx context.Context, tenantID, skuCode string, warehouseID int64, p Page) (*PageResult[model.StockJournal], error) {
	var result PageResult[model.StockJournal]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if skuCode != "" {
		q = q.Where("sku_code = ?", skuCode)
	}
	if warehouseID > 0 {
		q = q.Where("warehouse_id = ?", warehouseID)
	}
	if err := q.Model(&model.StockJournal{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

// ---------- 调拨单 ----------
type StockTransferRepo struct{ db *gorm.DB }

func NewStockTransferRepo(db *gorm.DB) *StockTransferRepo { return &StockTransferRepo{db: db} }

func (r *StockTransferRepo) Create(ctx context.Context, t *model.StockTransfer) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *StockTransferRepo) Get(ctx context.Context, id int64) (*model.StockTransfer, error) {
	var t model.StockTransfer
	err := r.db.WithContext(ctx).First(&t, id).Error
	return &t, err
}

func (r *StockTransferRepo) List(ctx context.Context, tenantID string, p Page) (*PageResult[model.StockTransfer], error) {
	var result PageResult[model.StockTransfer]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if err := q.Model(&model.StockTransfer{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *StockTransferRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.StockTransfer{}).Where("id = ?", id).Update("status", status).Error
}

// ---------- 盘点单 ----------
type StocktakeRepo struct{ db *gorm.DB }

func NewStocktakeRepo(db *gorm.DB) *StocktakeRepo { return &StocktakeRepo{db: db} }

func (r *StocktakeRepo) Create(ctx context.Context, s *model.Stocktake) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *StocktakeRepo) List(ctx context.Context, tenantID string, p Page) (*PageResult[model.Stocktake], error) {
	var result PageResult[model.Stocktake]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if err := q.Model(&model.Stocktake{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *StocktakeRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.Stocktake{}).Where("id = ?", id).Update("status", status).Error
}
