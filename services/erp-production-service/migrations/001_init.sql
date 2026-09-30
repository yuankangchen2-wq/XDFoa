CREATE DATABASE IF NOT EXISTS oa_production;
\c oa_production;

-- BOM 主表
CREATE TABLE IF NOT EXISTS prd_bom (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    product_sku VARCHAR(64) NOT NULL,
    bom_name VARCHAR(128) NOT NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- BOM 明细
CREATE TABLE IF NOT EXISTS prd_bom_item (
    id BIGSERIAL PRIMARY KEY,
    bom_id BIGINT NOT NULL,
    component_sku VARCHAR(64) NOT NULL,
    quantity BIGINT NOT NULL,
    level INT DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_bom_item_bom ON prd_bom_item(bom_id);

-- 生产工单
CREATE TABLE IF NOT EXISTS prd_work_order (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    order_no VARCHAR(64) UNIQUE NOT NULL,
    product_sku VARCHAR(64) NOT NULL,
    quantity BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    status INT DEFAULT 1,
    produced_qty BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_work_order_tenant ON prd_work_order(tenant_id);

-- 领料单
CREATE TABLE IF NOT EXISTS prd_material_issue (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    issue_no VARCHAR(64) UNIQUE NOT NULL,
    work_order_id BIGINT NOT NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_material_issue_wo ON prd_material_issue(work_order_id);

-- 领料明细
CREATE TABLE IF NOT EXISTS prd_material_issue_item (
    id BIGSERIAL PRIMARY KEY,
    issue_id BIGINT NOT NULL,
    sku_code VARCHAR(64) NOT NULL,
    quantity BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_material_issue_item ON prd_material_issue_item(issue_id);

-- 报工记录
CREATE TABLE IF NOT EXISTS prd_production_report (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    work_order_id BIGINT NOT NULL,
    produced_qty BIGINT NOT NULL,
    report_no VARCHAR(64) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
