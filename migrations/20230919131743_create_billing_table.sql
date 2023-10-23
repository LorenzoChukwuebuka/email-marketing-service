-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.billing (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    amount_paid DECIMAL(10,2) NOT NULL,
    plan_id INTEGER NOT NULL,
    duration VARCHAR(255) NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    reference VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
SELECT 'up SQL query';
-- +goose StatementEnd
 

-- +goose Down
-- +goose StatementBegin
DROP TABLE payments;
SELECT 'down SQL query';
-- +goose StatementEnd
