-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients (
    client_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    client_secret_hash TEXT NOT NULL,
    client_name VARCHAR(255) UNIQUE NOT NULL,
    client_description TEXT,
    clinet_logo_url TEXT,
    client_tos_url TEXT,
    client_policy_url TEXT,
    client_homepage_url TEXT,
    client_redirect_uris TEXT[] NOT NULL,
    client_scopes TEXT[] NOT NULL,
    client_grants TEXT[] NOT NULL,
    is_trusted BOOLEAN DEFAULT false,
    status VARCHAR(255) DEFAULT 'active',
    created_by UUID REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()

    -- scopes: [openid, profile, email, phone, address, offline_access]
    
    -- grants: [authorization_code, implicit, password, client_credentials, refresh_token]

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE clients;
-- +goose StatementEnd
