-- +migrate Up

ALTER TABLE contracts ADD COLUMN chain_id bigint NOT NULL DEFAULT 0;

-- +migrate Down

