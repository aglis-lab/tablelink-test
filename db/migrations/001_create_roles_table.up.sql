CREATE TABLE role_rights (
    id bigserial PRIMARY KEY,
    role_name varchar(50),
    right_create bool DEFAULT false,
    right_read bool DEFAULT false,
    right_update bool DEFAULT false,
    right_delete bool DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
