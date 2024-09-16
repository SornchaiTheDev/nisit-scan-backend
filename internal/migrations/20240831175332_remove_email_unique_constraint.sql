-- +goose Up
-- +goose StatementBegin
ALTER TABLE admins
DROP CONSTRAINT admins_email_key;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE admins
ADD CONSTRAINT admins_email_key UNIQUE (email);
-- +goose StatementEnd
