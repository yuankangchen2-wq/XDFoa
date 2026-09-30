-- 采购服务数据库
CREATE DATABASE IF NOT EXISTS oa_purchase;
\c oa_purchase;

-- 采购订单
CREATE TABLE IF NOT EXISTS pur_order (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    order_no VARCHAR(64) UNIQUE NOT NULL,
    supplier_code VARCHAR(64) NOT NULL,
    warehouse_id BIGINT NOT NULL,
    status INT DEFAULT 1,
    total_amount DECIMAL(18,2) DEFAULT 0,
    remark VARCHAR(512),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_pur_order_tenant ON pur_order(tenant_id);

-- 采购订单明细
CREATE TABLE IF NOT EXISTS pur_order_item (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    sku_code VARCHAR(64) NOT NULL,
    quantity BIGINT NOT NULL,
    unit_price DECIMAL(18,2) NOT NULL,
    received_qty BIGINT DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_pur_order_item_order ON pur_order_item(order_id);

-- 到货单
CREATE TABLE IF NOT EXISTS pur_receipt (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    receipt_no VARCHAR(64) UNIQUE NOT NULL,
    order_id BIGINT NOT NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_pur_receipt_order ON pur_receipt(order_id);

-- 到货明细
CREATE TABLE IF NOT EXISTS pur_receipt_item (
    id BIGSERIAL PRIMARY KEY,
    receipt_id BIGINT NOT NULL,
    sku_code VARCHAR(64) NOT NULL,
    quantity BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_pur_receipt_item_receipt ON pur_receipt_item(receipt_id);
