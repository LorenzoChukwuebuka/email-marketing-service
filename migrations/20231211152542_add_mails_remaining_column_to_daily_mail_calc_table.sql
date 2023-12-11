-- +goose Up
-- +goose StatementBegin

ALTER TABLE daily_mail_calc
ADD COLUMN mails_remaining INT NOT NULL DEFAULT 0;
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
