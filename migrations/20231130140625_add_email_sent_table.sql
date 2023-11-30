-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.daily_mail_calc (
    id SERIAL PRIMARY KEY,
    uuid character varying,
    subscription_id integer,
    mails_for_a_day integer,
    mails_sent integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE
);
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
