-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS authorization_code (
    authorization_code_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID NOT NULL,
    scope_id UUID NOT NULL,
    user_id UUID NOT NULL,
    code TEXT NOT NULL,
    redirect_uri TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE,
    FOREIGN KEY (scope_id) REFERENCES scope(scope_id) ON DELETE CASCADE,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
