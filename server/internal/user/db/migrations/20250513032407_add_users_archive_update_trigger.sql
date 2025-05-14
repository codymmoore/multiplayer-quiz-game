-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION disable_users_archive_update()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Updates to users_archive are not allowed';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER disable_users_archive_update
BEFORE UPDATE ON users_archive
FOR EACH ROW
EXECUTE FUNCTION disable_users_archive_update();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS disable_users_archive_update ON users_archive;
DROP FUNCTION IF EXISTS disable_users_archive_update;
-- +goose StatementEnd
