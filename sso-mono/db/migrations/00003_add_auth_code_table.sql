-- +goose Up
-- +goose StatementBegin
CREATE TABLE auth_codes (
    auth_code_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    client_id UUID REFERENCES clients(client_id) ON DELETE CASCADE,
    auth_code_hash TEXT NOT NULL,
    redirect_uri TEXT NOT NULL,
    scopes TEXT[],
    revoked BOOLEAN DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE auth_codes;
-- +goose StatementEnd
