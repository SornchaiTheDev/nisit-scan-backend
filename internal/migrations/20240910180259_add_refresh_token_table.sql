-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
	token VARCHAR(255),
	email VARCHAR(255),
	PRIMARY KEY(token, email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
