-- +goose Up
-- +goose StatementBegin
 CREATE TABLE IF NOT EXISTS public.logger(
      id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NUll,
    action TEXT NULL
 )

SELECT 'up SQL query';
-- +goose StatementEnd
Drop Table logger;
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
