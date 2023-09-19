-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

 CREATE TABLE IF NOT EXISTS public.admin
(
    id integer NOT NULL DEFAULT nextval('admin_id_seq'::regclass),
    firstname character varying(255) COLLATE pg_catalog."default",
    middlename character varying(255) COLLATE pg_catalog."default",
    lastname character varying(255) COLLATE pg_catalog."default",
    email character varying(255) COLLATE pg_catalog."default",
    password character varying(255) COLLATE pg_catalog."default",
    type character varying(50) COLLATE pg_catalog."default",
    created_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT admin_pkey PRIMARY KEY (id)
)

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
