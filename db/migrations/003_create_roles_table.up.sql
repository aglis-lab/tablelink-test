CREATE TABLE roles (
    id bigserial PRIMARY KEY,
    user_id bigint,
    role_right_id bigint,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
