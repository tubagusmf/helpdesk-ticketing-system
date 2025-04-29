
-- +migrate Up
create table attachments (
    "id" SERIAL PRIMARY KEY,
    "ticket_id" INT NOT NULL REFERENCES tickets("id") ON DELETE CASCADE,
    "file_path" TEXT NOT NULL,
    "uploaded_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
drop table attachments;
