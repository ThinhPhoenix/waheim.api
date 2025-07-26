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

CREATE TABLE apps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    status TEXT,
    uri TEXT,
    icon TEXT,
    publisher_id UUID,
    screenshots TEXT[],
    category TEXT,
    tags TEXT[],
    rating DOUBLE PRECISION DEFAULT 0,
    downloads INT DEFAULT 0
);

CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    app_id UUID NOT NULL,
    comment TEXT,
    stars INT NOT NULL CHECK (stars >= 1 AND stars <= 5),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    status TEXT,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_app FOREIGN KEY(app_id) REFERENCES apps(id)
);