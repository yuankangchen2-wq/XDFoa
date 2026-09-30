-- erp-cost-service init (PostgreSQL)

CREATE TABLE IF NOT EXISTS cost_center (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    type INT NOT NULL,
    remark VARCHAR(512),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE TABLE IF NOT EXISTS product_cost (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    product_id BIGINT NOT NULL,
    material_cost DECIMAL(18,4) NOT NULL DEFAULT 0,
    labor_cost DECIMAL(18,4) NOT NULL DEFAULT 0,
    overhead_cost DECIMAL(18,4) NOT NULL DEFAULT 0,
    total_cost DECIMAL(18,4) NOT NULL DEFAULT 0,
    quantity DECIMAL(18,4) NOT NULL,
    unit_cost DECIMAL(18,4) NOT NULL DEFAULT 0,
    period VARCHAR(32),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cost_record (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    product_id BIGINT NOT NULL,
    work_order_id BIGINT,
    cost_center_id BIGINT,
    element INT NOT NULL,
    amount DECIMAL(18,4) NOT NULL,
    period VARCHAR(32),
    status INT NOT NULL DEFAULT 1,
    remark VARCHAR(512),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cost_center_tenant ON cost_center(tenant_id);
CREATE INDEX IF NOT EXISTS idx_product_cost_tenant ON product_cost(tenant_id);
CREATE INDEX IF NOT EXISTS idx_product_cost_product ON product_cost(product_id);
CREATE INDEX IF NOT EXISTS idx_cost_record_tenant ON cost_record(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cost_record_product ON cost_record(product_id);
