-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS client_scope (
    client_id UUID NOT NULL,
    scope_id UUID NOT NULL,
    PRIMARY KEY (client_id, scope_id),
    FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE,
    FOREIGN KEY (scope_id) REFERENCES scope(scope_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS client_scope;
-- +goose StatementEnd
