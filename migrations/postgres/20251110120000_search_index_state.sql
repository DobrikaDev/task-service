-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS search_index_state (
    id SMALLINT PRIMARY KEY DEFAULT 1,
    last_synced_at TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00+00'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS search_index_state;
-- +goose StatementEnd

