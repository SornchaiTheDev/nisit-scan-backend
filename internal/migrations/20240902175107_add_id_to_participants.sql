-- +goose Up
-- +goose StatementBegin
ALTER TABLE participants
ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE participants
DROP COLUMN id;
-- +goose StatementEnd
