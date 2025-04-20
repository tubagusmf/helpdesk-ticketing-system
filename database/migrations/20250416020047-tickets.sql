
-- +migrate Up
Create type "status" AS ENUM ('open', 'in_progress', 'pending', 'resolved', 'closed');

Create type "priority" AS ENUM ('high', 'medium', 'low', 'very_low');

Create table tickets (
    "id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "status" status NOT NULL,
    "priority" priority NOT NULL,
    "assigned_to" INT REFERENCES users("id") ON DELETE SET NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS tickets;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS priority;
