-- +goose Up
-- +goose StatementBegin
ALTER TABLE subscriptions
ADD COLUMN transaction_id VARCHAR(255);
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
