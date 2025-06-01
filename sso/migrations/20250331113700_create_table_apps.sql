-- +goose Up
-- +goose StatementBegin
CREATE TABLE apps(
    id INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE ,
    secret VARCHAR(50) NOT NULL UNIQUE
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE apps CASCADE ;
-- +goose StatementEnd
