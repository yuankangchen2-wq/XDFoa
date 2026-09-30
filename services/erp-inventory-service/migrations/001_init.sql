-- 库存服务数据库初始化
CREATE DATABASE IF NOT EXISTS oa_inventory;
\c oa_inventory;

-- 仓库表
CREATE TABLE IF NOT EXISTS inv_warehouse (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    warehouse_code VARCHAR(64) NOT NULL,
    warehouse_name VARCHAR(128) NOT NULL,
    address VARCHAR(512),
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, warehouse_code)
);

-- 库存账户表
CREATE TABLE IF NOT EXISTS inv_stock_account (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    sku_code VARCHAR(64) NOT NULL,
    warehouse_id BIGINT NOT NULL,
    quantity BIGINT DEFAULT 0,
    allocated BIGINT DEFAULT 0,
    safety_stock BIGINT DEFAULT 0,
    batch_no VARCHAR(64),
    expire_at BIGINT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, sku_code, warehouse_id)
);
CREATE INDEX IF NOT EXISTS idx_stock_account_sku ON inv_stock_account(sku_code);

-- 库存流水表（不可变）
CREATE TABLE IF NOT EXISTS inv_stock_journal (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    sku_code VARCHAR(64) NOT NULL,
    warehouse_id BIGINT NOT NULL,
    change_type VARCHAR(32) NOT NULL,
    quantity BIGINT NOT NULL,
    before_qty BIGINT NOT NULL,
    after_qty BIGINT NOT NULL,
    reference_id VARCHAR(128),
    batch_no VARCHAR(64),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_journal_sku ON inv_stock_journal(sku_code);
CREATE INDEX IF NOT EXISTS idx_journal_ref ON inv_stock_journal(reference_id);

-- 调拨单
CREATE TABLE IF NOT EXISTS inv_stock_transfer (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    transfer_no VARCHAR(64) UNIQUE NOT NULL,
    from_warehouse_id BIGINT NOT NULL,
    to_warehouse_id BIGINT NOT NULL,
    sku_code VARCHAR(64) NOT NULL,
    quantity BIGINT NOT NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 盘点单
CREATE TABLE IF NOT EXISTS inv_stocktake (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    stocktake_no VARCHAR(64) UNIQUE NOT NULL,
    warehouse_id BIGINT NOT NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 盘点明细
CREATE TABLE IF NOT EXISTS inv_stocktake_item (
    id BIGSERIAL PRIMARY KEY,
    stocktake_id BIGINT NOT NULL,
    sku_code VARCHAR(64) NOT NULL,
    system_qty BIGINT NOT NULL,
    actual_qty BIGINT NOT NULL,
    diff_qty BIGINT NOT NULL
);
