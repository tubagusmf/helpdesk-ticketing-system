
-- +migrate Up

create table ticket_histories (
    "id" SERIAL PRIMARY KEY,
    "ticket_id" INT NOT NULL REFERENCES tickets("id") ON DELETE CASCADE,
    "user_id" INT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "status" status NOT NULL,
    "priority" priority NOT NULL,
    "changed_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down

DROP TABLE IF EXISTS ticket_histories;