-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS public.payments
(
    id integer NOT NULL DEFAULT nextval('payments_id_seq'::regclass),
    user_id integer NOT NULL,
    amount_paid numeric(10,2) NOT NULL,
    plan_id integer NOT NULL,
    duration character varying(255) COLLATE pg_catalog."default" NOT NULL,
    expiry_date timestamp without time zone NOT NULL,
    reference character varying(255) COLLATE pg_catalog."default" NOT NULL,
    status character varying(255) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    CONSTRAINT payments_pkey PRIMARY KEY (id)
)
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
