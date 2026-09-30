-- crm-analytics-service init (MySQL, shared oa_crm)

CREATE TABLE IF NOT EXISTS ana_sales_snapshot (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    period VARCHAR(32) NOT NULL,
    product_id BIGINT NOT NULL DEFAULT 0,
    customer_id BIGINT NOT NULL DEFAULT 0,
    total_amount DECIMAL(18,2) NOT NULL,
    order_count INT NOT NULL DEFAULT 0,
    avg_amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ana_sales_tenant (tenant_id),
    INDEX idx_ana_sales_period (period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ana_funnel_stage (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    period VARCHAR(32) NOT NULL,
    stage VARCHAR(64) NOT NULL,
    opportunity_count INT NOT NULL DEFAULT 0,
    amount DECIMAL(18,2) NOT NULL DEFAULT 0,
    conversion_rate DECIMAL(8,4) NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ana_funnel_tenant (tenant_id),
    INDEX idx_ana_funnel_period (period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ana_customer_stat (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    period VARCHAR(32) NOT NULL,
    new_customers INT NOT NULL DEFAULT 0,
    active_customers INT NOT NULL DEFAULT 0,
    total_contribution DECIMAL(18,2) NOT NULL DEFAULT 0,
    top_customer_id BIGINT NOT NULL DEFAULT 0,
    top_contribution DECIMAL(18,2) NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ana_cust_tenant (tenant_id),
    INDEX idx_ana_cust_period (period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
