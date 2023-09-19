-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS public.plans
(
    id integer NOT NULL DEFAULT nextval('plans_id_seq'::regclass),
    planname character varying(255) COLLATE pg_catalog."default" NOT NULL,
    duration character varying(50) COLLATE pg_catalog."default" NOT NULL,
    price real NOT NULL,
    details text COLLATE pg_catalog."default" NOT NULL,
    status character varying(50) COLLATE pg_catalog."default",
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT plans_pkey PRIMARY KEY (id)
)
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
