-- IAM 服务 - 初始化表结构
-- 数据库: PostgreSQL

-- 用户表
CREATE TABLE IF NOT EXISTS iam_user (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    username VARCHAR(64) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    real_name VARCHAR(64),
    email VARCHAR(128),
    phone VARCHAR(32),
    org_id BIGINT DEFAULT 0,
    org_path VARCHAR(512),
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_tenant_username ON iam_user(tenant_id, username);
CREATE INDEX IF NOT EXISTS idx_iam_user_tenant ON iam_user(tenant_id);

-- 角色表
CREATE TABLE IF NOT EXISTS iam_role (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    role_code VARCHAR(64) NOT NULL,
    role_name VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    data_scope VARCHAR(32) DEFAULT 'self',
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_tenant_role ON iam_role(tenant_id, role_code);

-- 用户-角色关联
CREATE TABLE IF NOT EXISTS iam_user_role (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_iam_user_role_user ON iam_user_role(user_id);

-- 权限表
CREATE TABLE IF NOT EXISTS iam_permission (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    perm_code VARCHAR(128) NOT NULL,
    perm_name VARCHAR(128) NOT NULL,
    perm_type INT DEFAULT 3,
    path VARCHAR(256),
    method VARCHAR(16),
    status INT DEFAULT 1
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_tenant_perm ON iam_permission(tenant_id, perm_code);

-- Casbin 策略表
CREATE TABLE IF NOT EXISTS iam_casbin_rule (
    id BIGSERIAL PRIMARY KEY,
    ptype VARCHAR(16) NOT NULL,
    v0 VARCHAR(128),
    v1 VARCHAR(256),
    v2 VARCHAR(128),
    v3 VARCHAR(128),
    v4 VARCHAR(128),
    v5 VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_casbin_ptype_v0 ON iam_casbin_rule(ptype, v0);

-- 审计日志表
CREATE TABLE IF NOT EXISTS iam_audit_log (
    id BIGSERIAL PRIMARY KEY,
    tenant_id VARCHAR(64),
    user_id BIGINT,
    username VARCHAR(64),
    module VARCHAR(64),
    action VARCHAR(64),
    resource VARCHAR(128),
    client_ip VARCHAR(64),
    user_agent VARCHAR(512),
    detail TEXT,
    status INT,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_tenant_time ON iam_audit_log(tenant_id, created_at);

-- 初始化数据：默认管理员用户 (admin/admin123)
INSERT INTO iam_user (tenant_id, username, password_hash, real_name, status)
VALUES ('default', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '系统管理员', 1)
ON CONFLICT DO NOTHING;

-- 初始化角色
INSERT INTO iam_role (tenant_id, role_code, role_name, description, data_scope)
VALUES ('default', 'admin', '超级管理员', '拥有所有权限', 'all')
ON CONFLICT DO NOTHING;
