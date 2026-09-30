-- CRM 客户服务数据库（MySQL）
CREATE DATABASE IF NOT EXISTS oa_crm DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
USE oa_crm;

-- 客户表
CREATE TABLE IF NOT EXISTS crm_customer (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    contact_person VARCHAR(64),
    phone VARCHAR(32),
    email VARCHAR(128),
    address VARCHAR(256),
    level VARCHAR(8),
    industry VARCHAR(64),
    owner_id VARCHAR(64),
    status INT DEFAULT 1,
    remark VARCHAR(512),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 联系人表
CREATE TABLE IF NOT EXISTS crm_contact (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    name VARCHAR(64) NOT NULL,
    position VARCHAR(64),
    phone VARCHAR(32),
    email VARCHAR(128),
    is_primary INT DEFAULT 0,
    INDEX idx_customer (customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 跟进记录表
CREATE TABLE IF NOT EXISTS crm_follow_up (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    type VARCHAR(32),
    content TEXT,
    user_id VARCHAR(64),
    follow_up_at DATETIME,
    INDEX idx_customer (customer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
