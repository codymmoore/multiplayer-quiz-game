-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_verification_code ON verification_code (verification_code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_username;
-- +goose StatementEnd
