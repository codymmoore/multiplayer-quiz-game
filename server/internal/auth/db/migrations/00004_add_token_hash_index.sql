-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_token_hash ON refresh_token (token_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_token_hash;
-- +goose StatementEnd
