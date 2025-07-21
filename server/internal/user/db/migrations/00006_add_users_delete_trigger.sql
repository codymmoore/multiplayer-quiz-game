-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION archive_user()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO users_archive (
        users_id, username, email, password_hash, is_verified, created_at, updated_at
    ) VALUES (
        OLD.id, OLD.username, OLD.email, OLD.password_hash, OLD.is_verified, OLD.created_at, OLD.updated_at
    );

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER archive_user
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION archive_user();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS archive_user ON users;
DROP FUNCTION IF EXISTS archive_user;
-- +goose StatementEnd
