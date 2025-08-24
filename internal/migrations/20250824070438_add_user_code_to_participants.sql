-- +goose Up
-- +goose StatementBegin
ALTER TABLE participants ADD COLUMN user_code VARCHAR(10);
ALTER TABLE participants ADD CONSTRAINT fk_user_code FOREIGN KEY (user_code) REFERENCES users(code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE participants DROP CONSTRAINT fk_user_code;
ALTER TABLE participants DROP COLUMN user_code;
-- +goose StatementEnd
