-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY,
    login VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    hashed_password VARCHAR(512) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    wallet_usdt FLOAT8 NOT NULL DEFAULT 0.0,
    number_purchases INT NOT NULL DEFAULT 0,
    salt VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
