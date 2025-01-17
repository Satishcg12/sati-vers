-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_token (
    access_token UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    refresh_token UUID NOT NULL,
    client_id UUID NOT NULL,
    user_id UUID NOT NULL,
    scopes TEXT[] NOT NULL,
    expires_in TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE,
    FOREIGN KEY (refresh_token) REFERENCES refresh_token(refresh_token) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_token;
-- +goose StatementEnd
