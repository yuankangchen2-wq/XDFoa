-- 订单服务数据库（MySQL）
CREATE DATABASE IF NOT EXISTS oa_order DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
USE oa_order;

-- 订单主表
CREATE TABLE IF NOT EXISTS ord_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    order_no VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    status INT DEFAULT 1,
    total_amount DECIMAL(18,2) DEFAULT 0,
    warehouse_id BIGINT NOT NULL,
    remark VARCHAR(512),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id),
    INDEX idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 订单明细
CREATE TABLE IF NOT EXISTS ord_order_item (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    sku_code VARCHAR(64) NOT NULL,
    quantity BIGINT NOT NULL,
    unit_price DECIMAL(18,2) NOT NULL,
    INDEX idx_order (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 订单状态日志
CREATE TABLE IF NOT EXISTS ord_order_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    from_status INT NOT NULL,
    to_status INT NOT NULL,
    remark VARCHAR(512),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
