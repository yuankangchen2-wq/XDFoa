package model

import "time"

// 客户主数据
type Customer struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      string    `gorm:"index;size:64" json:"tenant_id"`
	CustomerCode  string    `gorm:"size:64" json:"customer_code"`
	CustomerName  string    `gorm:"size:256" json:"customer_name"`
	CustomerType  string    `gorm:"size:32" json:"customer_type"` // enterprise/individual
	Industry      string    `gorm:"size:64" json:"industry"`
	Region        string    `gorm:"size:64" json:"region"`
	ContactName   string    `gorm:"size:64" json:"contact_name"`
	ContactPhone  string    `gorm:"size:32" json:"contact_phone"`
	ContactEmail  string    `gorm:"size:128" json:"contact_email"`
	Address       string    `gorm:"size:512" json:"address"`
	Status        int32     `gorm:"default:1" json:"status"` // 1=正常 0=禁用
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Customer) TableName() string { return "mdm_customer" }

// 商品/SKU
type Product struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	SkuCode     string    `gorm:"size:64" json:"sku_code"`
	ProductName string    `gorm:"size:256" json:"product_name"`
	Category    string    `gorm:"size:64" json:"category"`
	Unit        string    `gorm:"size:16" json:"unit"`
	Spec        string    `gorm:"size:128" json:"spec"`
	Status      int32     `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Product) TableName() string { return "mdm_product" }

// 供应商
type Supplier struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      string    `gorm:"index;size:64" json:"tenant_id"`
	SupplierCode  string    `gorm:"size:64" json:"supplier_code"`
	SupplierName  string    `gorm:"size:256" json:"supplier_name"`
	ContactName   string    `gorm:"size:64" json:"contact_name"`
	ContactPhone  string    `gorm:"size:32" json:"contact_phone"`
	ContactEmail  string    `gorm:"size:128" json:"contact_email"`
	Address       string    `gorm:"size:512" json:"address"`
	Status        int32     `gorm:"default:1" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Supplier) TableName() string { return "mdm_supplier" }

// 组织（树形）
type Organization struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;size:64" json:"tenant_id"`
	ParentID  int64     `gorm:"default:0" json:"parent_id"`
	OrgCode   string    `gorm:"size:64" json:"org_code"`
	OrgName   string    `gorm:"size:256" json:"org_name"`
	OrgPath   string    `gorm:"size:512" json:"org_path"` // /root/division/dept
	OrgLevel  int32     `gorm:"default:1" json:"org_level"`
	SortOrder int32     `gorm:"default:0" json:"sort_order"`
	Status    int32     `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Organization) TableName() string { return "mdm_org" }
