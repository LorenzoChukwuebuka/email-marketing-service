-- +goose Up
-- +goose StatementBegin

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
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE otp;
SELECT 'down SQL query';
-- +goose StatementEnd
