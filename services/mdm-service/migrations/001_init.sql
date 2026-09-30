-- MDM 主数据服务 - 初始化表结构
-- 数据库: PostgreSQL

-- 客户表
CREATE TABLE IF NOT EXISTS mdm_customer (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    customer_code VARCHAR(64),
    customer_name VARCHAR(256) NOT NULL,
    customer_type VARCHAR(32) DEFAULT 'enterprise',
    industry VARCHAR(64),
    region VARCHAR(64),
    contact_name VARCHAR(64),
    contact_phone VARCHAR(32),
    contact_email VARCHAR(128),
    address VARCHAR(512),
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_mdm_customer_tenant ON mdm_customer(tenant_id);
CREATE INDEX IF NOT EXISTS idx_mdm_customer_code ON mdm_customer(tenant_id, customer_code);

-- 商品表
CREATE TABLE IF NOT EXISTS mdm_product (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    sku_code VARCHAR(64),
    product_name VARCHAR(256) NOT NULL,
    category VARCHAR(64),
    unit VARCHAR(16),
    spec VARCHAR(128),
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_mdm_product_tenant ON mdm_product(tenant_id);
CREATE INDEX IF NOT EXISTS idx_mdm_product_sku ON mdm_product(tenant_id, sku_code);

-- 供应商表
CREATE TABLE IF NOT EXISTS mdm_supplier (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    supplier_code VARCHAR(64),
    supplier_name VARCHAR(256) NOT NULL,
    contact_name VARCHAR(64),
    contact_phone VARCHAR(32),
    contact_email VARCHAR(128),
    address VARCHAR(512),
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_mdm_supplier_tenant ON mdm_supplier(tenant_id);
CREATE INDEX IF NOT EXISTS idx_mdm_supplier_code ON mdm_supplier(tenant_id, supplier_code);

-- 组织表（树形）
CREATE TABLE IF NOT EXISTS mdm_org (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    parent_id BIGINT DEFAULT 0,
    org_code VARCHAR(64) NOT NULL,
    org_name VARCHAR(256) NOT NULL,
    org_path VARCHAR(512),
    org_level INT DEFAULT 1,
    sort_order INT DEFAULT 0,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_mdm_org_tenant ON mdm_org(tenant_id);
CREATE INDEX IF NOT EXISTS idx_mdm_org_parent ON mdm_org(tenant_id, parent_id);
CREATE UNIQUE INDEX IF NOT EXISTS uk_mdm_org_code ON mdm_org(tenant_id, org_code);
