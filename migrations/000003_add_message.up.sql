-- SQL for the 'up' migration
-- Add your 'up' migration SQL here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE messages (
    id uuid DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    room_id uuid NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_room_id FOREIGN KEY (room_id) REFERENCES rooms (id)
);