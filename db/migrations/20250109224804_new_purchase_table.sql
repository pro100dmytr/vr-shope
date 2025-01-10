-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS purchases (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    cost FLOAT8 NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS purchases;
-- +goose StatementEnd
