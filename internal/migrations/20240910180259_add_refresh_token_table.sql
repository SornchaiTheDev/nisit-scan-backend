-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
	email VARCHAR(255) PRIMARY KEY,
	token VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
