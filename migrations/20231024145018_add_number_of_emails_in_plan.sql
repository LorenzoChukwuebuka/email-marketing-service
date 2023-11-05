-- +goose Up
-- +goose StatementBegin
ALTER TABLE plans
ADD COLUMN number_of_emails_per_day VARCHAR(255);
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
