-- portal-bff 门户库（oa_portal，PostgreSQL）

CREATE TABLE IF NOT EXISTS portal_approval (
    id            BIGSERIAL PRIMARY KEY,
    title         VARCHAR(200) NOT NULL,
    type          VARCHAR(32)  NOT NULL,
    serial_no     VARCHAR(64)  DEFAULT '',
    applicant     VARCHAR(64)  DEFAULT '',
    applicant_id  VARCHAR(64)  DEFAULT '',
    status        INTEGER      DEFAULT 1,
    current_node  VARCHAR(128) DEFAULT '',
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_portal_approval_type ON portal_approval(type);
CREATE INDEX IF NOT EXISTS idx_portal_approval_status ON portal_approval(status);
CREATE INDEX IF NOT EXISTS idx_portal_approval_applicant ON portal_approval(applicant_id);

CREATE TABLE IF NOT EXISTS portal_schedule (
    id          BIGSERIAL PRIMARY KEY,
    user_id     VARCHAR(64) DEFAULT '',
    title       VARCHAR(200) NOT NULL,
    event_date  DATE,
    created_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_portal_schedule_user_date ON portal_schedule(user_id, event_date);

CREATE TABLE IF NOT EXISTS portal_favorite (
    id          BIGSERIAL PRIMARY KEY,
    user_id     VARCHAR(64) DEFAULT '',
    name        VARCHAR(64) NOT NULL,
    icon        VARCHAR(64) DEFAULT 'Star',
    path        VARCHAR(256) DEFAULT '',
    sort        INTEGER DEFAULT 0,
    created_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_portal_favorite_user ON portal_favorite(user_id);
