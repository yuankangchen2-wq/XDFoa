package model

import "time"

// 仓库
type Warehouse struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      string    `gorm:"index;size:64" json:"tenant_id"`
	WarehouseCode string    `gorm:"size:64;uniqueIndex:uk_tenant_wh" json:"warehouse_code"`
	WarehouseName string    `gorm:"size:128" json:"warehouse_name"`
	Address       string    `gorm:"size:512" json:"address"`
	Status        int32     `gorm:"default:1" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Warehouse) TableName() string { return "inv_warehouse" }

// 库存账户（sku + warehouse，可改余额）
type StockAccount struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	SkuCode     string    `gorm:"size:64;index" json:"sku_code"`
	WarehouseID int64     `gorm:"index" json:"warehouse_id"`
	Quantity    int64     `gorm:"default:0" json:"quantity"`   // 可用库存
	Allocated   int64     `gorm:"default:0" json:"allocated"`  // 已预扣锁定
	SafetyStock int64     `gorm:"default:0" json:"safety_stock"`
	BatchNo     string    `gorm:"size:64" json:"batch_no"`
	ExpireAt    int64     `json:"expire_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (StockAccount) TableName() string { return "inv_stock_account" }

// 库存流水（不可变，所有变动溯源）
type StockJournal struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	SkuCode      string    `gorm:"size:64;index" json:"sku_code"`
	WarehouseID  int64     `gorm:"index" json:"warehouse_id"`
	ChangeType   string    `gorm:"size:32" json:"change_type"` // in/out/allocate/release/transfer_in/transfer_out/adjust
	Quantity     int64     `json:"quantity"`     // 变动数量（正负）
	BeforeQty    int64     `json:"before_qty"`  // 变动前可用量
	AfterQty     int64     `json:"after_qty"`   // 变动后可用量
	ReferenceID  string    `gorm:"size:128" json:"reference_id"` // 关联单据
	BatchNo      string    `gorm:"size:64" json:"batch_no"`
	CreatedAt    time.Time `json:"created_at"`
}

func (StockJournal) TableName() string { return "inv_stock_journal" }

// 调拨单
type StockTransfer struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	TenantID        string    `gorm:"index;size:64" json:"tenant_id"`
	TransferNo      string    `gorm:"size:64;uniqueIndex" json:"transfer_no"`
	FromWarehouseID int64     `json:"from_warehouse_id"`
	ToWarehouseID   int64     `json:"to_warehouse_id"`
	SkuCode         string    `gorm:"size:64" json:"sku_code"`
	Quantity        int64     `json:"quantity"`
	Status          int32     `gorm:"default:1" json:"status"` // 1=待出库 2=已出库 3=已入库 4=已取消
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (StockTransfer) TableName() string { return "inv_stock_transfer" }

// 盘点单
type Stocktake struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	StocktakeNo string    `gorm:"size:64;uniqueIndex" json:"stocktake_no"`
	WarehouseID int64     `json:"warehouse_id"`
	Status      int32     `gorm:"default:1" json:"status"` // 1=进行中 2=已完成
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Stocktake) TableName() string { return "inv_stocktake" }

// 盘点明细
type StocktakeItem struct {
	ID          int64 `gorm:"primaryKey" json:"id"`
	StocktakeID int64 `gorm:"index" json:"stocktake_id"`
	SkuCode     string    `gorm:"size:64" json:"sku_code"`
	SystemQty   int64     `json:"system_qty"`
	ActualQty   int64     `json:"actual_qty"`
	DiffQty     int64     `json:"diff_qty"`
}

func (StocktakeItem) TableName() string { return "inv_stocktake_item" }
