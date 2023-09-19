-- +goose Up
-- +goose StatementBegin
 CREATE TABLE IF NOT EXISTS public.users
(
    id SERIAL PRIMARY KEY,
    uuid VARCHAR,
    firstname VARCHAR,
    middlename VARCHAR,
    lastname VARCHAR,
    email VARCHAR,
    password VARCHAR,
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    username VARCHAR
);

SELECT 'up SQL query';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE users;
SELECT 'down SQL query';
-- +goose StatementEnd
