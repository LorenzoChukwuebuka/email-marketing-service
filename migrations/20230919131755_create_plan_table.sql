-- +goose Up
-- +goose StatementBegin

 CREATE TABLE IF NOT EXISTS public.plans (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NOT NUll,
    planname VARCHAR(255) NOT NULL,
    duration VARCHAR(50) NOT NULL,
    price REAL NOT NULL,
    details TEXT NOT NULL,
    status VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

SELECT 'up SQL query';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE plans;
SELECT 'down SQL query';
-- +goose StatementEnd
