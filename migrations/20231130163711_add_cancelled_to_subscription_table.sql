-- +goose Up
-- +goose StatementBegin

ALTER TABLE subscriptions
ADD COLUMN cancelled BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN date_cancelled TIMESTAMP WITH TIME ZONE;
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
