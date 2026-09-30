-- 商机表（CRM 库）
USE oa_crm;

CREATE TABLE IF NOT EXISTS crm_opportunity (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    customer_id BIGINT NOT NULL,
    stage INT DEFAULT 1,
    amount DECIMAL(18,2) DEFAULT 0,
    win_rate INT DEFAULT 10,
    expected_close_date VARCHAR(32),
    owner_id VARCHAR(64),
    remark VARCHAR(512),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id),
    INDEX idx_customer (customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
