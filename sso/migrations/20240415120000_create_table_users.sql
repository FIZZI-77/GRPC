-- +goose Up
-- +goose StatementBegin

CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users CASCADE;
-- +goose StatementEnd