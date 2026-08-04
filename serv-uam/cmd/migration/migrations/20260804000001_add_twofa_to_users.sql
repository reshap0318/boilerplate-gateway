-- +goose Up
ALTER TABLE users
    ADD COLUMN twofa BOOLEAN NOT NULL DEFAULT FALSE AFTER locked_until;

-- +goose Down
ALTER TABLE users
    DROP COLUMN twofa;
