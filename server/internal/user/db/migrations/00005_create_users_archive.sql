-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_archive (
   id SERIAL PRIMARY KEY,
   users_id INTEGER NOT NULL,
   username VARCHAR(15) NOT NULL,
   email VARCHAR(255) NOT NULL,
   password_hash TEXT NOT NULL,
   is_verified BOOLEAN NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE NOT NULL,
   updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
   archived_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_archive;
-- +goose StatementEnd