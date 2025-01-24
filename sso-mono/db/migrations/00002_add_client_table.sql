-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients (
    client_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_secret_hash TEXT NOT NULL,
    client_name VARCHAR(255) UNIQUE NOT NULL,
    client_description TEXT,
    clinet_logo_url TEXT,
    client_homepage_url TEXT,
    client_redirect_uris TEXT[],
    client_scopes TEXT[] check (clent_scopes IN ('openid', 'profile', 'email', 'offline_access')),
    client_grants TEXT[] check (client_grants IN ('authorization_code', 'client_credentials', 'refresh_token', 'password')),
    is_trusted BOOLEAN DEFAULT false,
    status VARCHAR(255) DEFAULT 'active',
    created_by UUID REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE clients;
-- +goose StatementEnd
