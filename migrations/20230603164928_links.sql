-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS link
(
    initial_url text NOT NULL,
    shorten_url text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE link;
-- +goose StatementEnd
