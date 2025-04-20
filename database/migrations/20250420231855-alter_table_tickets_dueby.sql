
-- +migrate Up
ALTER TABLE tickets ADD COLUMN "due_by" TIMESTAMP;

-- +migrate Down

ALTER TABLE tickets DROP COLUMN "due_by";
