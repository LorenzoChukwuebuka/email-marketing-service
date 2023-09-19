 -- +goose Up
-- +goose StatementBegin

-- Add a comment describing the purpose of this migration
-- This is optional but can be helpful for future reference

CREATE TABLE IF NOT EXISTS public.otp (
    id SERIAL PRIMARY KEY,
    user_id integer,
    token character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    uuid character varying,
    FOREIGN KEY (user_id) REFERENCES public.users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

-- +goose StatementEnd
