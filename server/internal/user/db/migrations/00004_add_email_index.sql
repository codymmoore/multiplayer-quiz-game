-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_email ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_email;
-- +goose StatementEnd
