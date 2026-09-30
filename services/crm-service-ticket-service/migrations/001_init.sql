-- CRM 工单服务（oa_crm 库）
USE oa_crm;

-- 工单表
CREATE TABLE IF NOT EXISTS crm_service_ticket (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    ticket_no VARCHAR(64) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    title VARCHAR(256) NOT NULL,
    description TEXT,
    status INT DEFAULT 1,
    priority INT DEFAULT 2,
    assignee_id VARCHAR(64),
    satisfaction INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id),
    INDEX idx_customer (customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 工单回复表
CREATE TABLE IF NOT EXISTS crm_ticket_reply (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    ticket_id BIGINT NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ticket (ticket_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
