-- CRM 销售服务（oa_crm 库）
USE oa_crm;

-- 销售订单
CREATE TABLE IF NOT EXISTS crm_sales_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    order_no VARCHAR(64) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    status INT DEFAULT 1,
    total_amount DECIMAL(18,2) DEFAULT 0,
    remark VARCHAR(512),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id),
    INDEX idx_customer (customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 合同
CREATE TABLE IF NOT EXISTS crm_contract (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    contract_no VARCHAR(64) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    amount DECIMAL(18,2) DEFAULT 0,
    paid_amount DECIMAL(18,2) DEFAULT 0,
    status INT DEFAULT 1,
    sign_date VARCHAR(32),
    remark VARCHAR(512),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id),
    INDEX idx_customer (customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 回款
CREATE TABLE IF NOT EXISTS crm_payment (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    contract_id BIGINT NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    pay_date VARCHAR(32),
    pay_method VARCHAR(32),
    remark VARCHAR(512),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_contract (contract_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
