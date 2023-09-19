-- +goose Up
-- +goose StatementBegin

 CREATE TABLE IF NOT EXISTS public.subscriptions (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NOT NUll,
    user_id INTEGER NOT NULL,
    plan_id INTEGER NOT NULL,
    payment_id INTEGER NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    expired BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

SELECT 'up SQL query';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE subscriptions;
SELECT 'down SQL query';
-- +goose StatementEnd
