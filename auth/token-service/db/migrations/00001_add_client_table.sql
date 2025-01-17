-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Enable UUID generation

CREATE TABLE IF NOT EXISTS client (
    client_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_secret TEXT NOT NULL,
    client_name varchar(100) NOT NULL,
    client_description TEXT,
    redirect_uris TEXT[] NOT NULL,
    grant_types TEXT[] NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS client;
-- +goose StatementEnd
