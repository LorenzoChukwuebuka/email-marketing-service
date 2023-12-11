-- +goose Up
-- +goose StatementBegin

 CREATE TABLE IF NOT EXISTS public.contact (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(255) NUll,
    user_id INT NOT Null,
    contact VARCHAR(255) NOT NULL,
    subscribed VARCHAR(255),
    blocked BOOLEAN DEFAULT FALSE,
    firstname VARCHAR(255) NULL,
    lastname VARCHAR(255) NULL, 
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
