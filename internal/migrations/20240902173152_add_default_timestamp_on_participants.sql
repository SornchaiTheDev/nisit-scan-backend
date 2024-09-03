-- +goose Up
-- +goose StatementBegin
ALTER TABLE participants
ALTER COLUMN timestamp SET DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE participants
ALTER COLUMN timestamp DROP DEFAULT;
-- +goose StatementEnd
