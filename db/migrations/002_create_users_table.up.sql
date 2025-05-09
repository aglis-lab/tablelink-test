CREATE TABLE users (
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL,
    email varchar(100) NOT NULL,
    password varchar(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
