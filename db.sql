DROP TABLE IF EXISTS users CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username    TEXT NOT NULL UNIQUE,
    password    TEXT NOT NULL,
    email       TEXT NOT NULL UNIQUE,
    phone       TEXT,
    address     TEXT,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TEXT,
    is_active   BOOLEAN DEFAULT TRUE,
    role        TEXT NOT NULL DEFAULT 'user',
    avatar      TEXT,
    first_name  TEXT,
    last_name   TEXT,
    date_of_birth TIMESTAMPTZ,
    gender      TEXT,
    status      TEXT
);