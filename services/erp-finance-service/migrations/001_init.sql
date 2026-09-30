-- erp-finance-service init (PostgreSQL)

CREATE TABLE IF NOT EXISTS fin_receivable (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    customer_id BIGINT NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    received_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    status INT NOT NULL DEFAULT 1,
    due_date VARCHAR(32),
    remark VARCHAR(512),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fin_payable (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    supplier_id BIGINT NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    paid_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    status INT NOT NULL DEFAULT 1,
    due_date VARCHAR(32),
    remark VARCHAR(512),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fin_payment (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    type INT NOT NULL,
    ref_id BIGINT NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    pay_date VARCHAR(32),
    pay_method VARCHAR(32),
    remark VARCHAR(512),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fin_receivable_tenant ON fin_receivable(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fin_receivable_customer ON fin_receivable(customer_id);
CREATE INDEX IF NOT EXISTS idx_fin_payable_tenant ON fin_payable(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fin_payable_supplier ON fin_payable(supplier_id);
CREATE INDEX IF NOT EXISTS idx_fin_payment_tenant ON fin_payment(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fin_payment_ref ON fin_payment(ref_id);
