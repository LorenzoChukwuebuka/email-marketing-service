-- +goose Up
-- +goose StatementBegin
  CREATE TABLE IF NOT EXISTS public.logger (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NOT NUll,
    action TEXT NOT NULL,
    performed_by INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

SELECT 'up SQL query';
-- +goose StatementEnd
Drop Table logger;
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
