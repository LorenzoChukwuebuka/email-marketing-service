-- +goose Up
-- +goose StatementBegin

 CREATE TABLE IF NOT EXISTS public.admin (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NUll,
    firstname VARCHAR(255),
    middlename VARCHAR(255),
    lastname VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    type VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE admin;
SELECT 'down SQL query';
-- +goose StatementEnd
