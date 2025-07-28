-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_token (
    id SERIAL PRIMARY KEY,
    token_hash TEXT NOT NULL,
    issued_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    active BOOLEAN DEFAULT true NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_token;
-- +goose StatementEnd
