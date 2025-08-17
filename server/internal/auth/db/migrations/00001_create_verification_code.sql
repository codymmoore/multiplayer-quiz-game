-- +goose Up
-- +goose StatementBegin
CREATE TABLE verification_code (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    verification_code VARCHAR(32) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '15 minutes'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE verification_code;
-- +goose StatementEnd
