-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS id_email ON users (email);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_email;
-- +goose StatementEnd
