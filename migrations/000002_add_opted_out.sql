-- +goose Up
ALTER TABLE users ADD COLUMN opted_out BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN opted_out;
