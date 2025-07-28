-- +goose Up
-- +goose StatementBegin
CREATE TABLE verification_code (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    verification_code VARCHAR(6) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE verification_code;
-- +goose StatementEnd
