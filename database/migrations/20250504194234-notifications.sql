
-- +migrate Up
create table notifications (
    "id" SERIAL PRIMARY KEY,
    "ticket_id" INT NOT NULL REFERENCES tickets("id") ON DELETE CASCADE,
    "user_id" INT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "email" VARCHAR(100) NOT NULL,
    "subject" VARCHAR(255) NOT NULL,
    "message" TEXT NOT NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'pending',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
drop table if exists notifications;