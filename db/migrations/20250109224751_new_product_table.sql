-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products(
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cost FLOAT8 NOT NULL,
    quantity_stock INT NOT NULL DEFAULT 0,
    guarantees TIMESTAMP NOT NULL,
    country VARCHAR(255) NOT NULL,
    likes INT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
