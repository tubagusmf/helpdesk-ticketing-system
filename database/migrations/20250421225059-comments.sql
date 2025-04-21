
-- +migrate Up

CREATE TABLE comments (
    "id" SERIAL PRIMARY KEY,
    "ticket_id" INT NOT NULL REFERENCES tickets("id") ON DELETE CASCADE,
    "user_id" INT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "content" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS comments;
