-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS scope (
    scope_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scope_name varchar(100) NOT NULL,
    scope_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scope;
-- +goose StatementEnd
