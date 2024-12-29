-- SQL for the 'up' migration
-- Add your 'up' migration SQL here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE rooms (
    id uuid DEFAULT uuid_generate_v4(),
    name VARCHAR(250) NOT NULL,
    users VARCHAR(254)[] NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);

CREATE INDEX idx_users ON rooms(users);