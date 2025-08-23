-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
	code VARCHAR(10) PRIMARY KEY,
	full_name TEXT NOT NULL,
	gmail TEXT NOT NULL,
	major VARCHAR(6) NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd
