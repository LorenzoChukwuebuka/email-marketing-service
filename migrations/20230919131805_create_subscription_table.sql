-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS public.subscriptions
(
    id integer NOT NULL DEFAULT nextval('subscriptions_id_seq'::regclass),
    user_id integer NOT NULL,
    plan_id integer NOT NULL,
    payment_id integer NOT NULL,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    expired boolean NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    CONSTRAINT subscriptions_pkey PRIMARY KEY (id)
)
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
