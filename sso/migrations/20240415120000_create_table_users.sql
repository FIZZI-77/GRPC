-- +goose Up
-- +goose StatementBegin

CREATE TABLE(
    id INTEGER PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE ,
    pass_hash BYTEA
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE top_up CASCADE;
-- +goose StatementEnd