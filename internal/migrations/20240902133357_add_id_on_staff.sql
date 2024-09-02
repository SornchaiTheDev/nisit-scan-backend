-- +goose Up
-- +goose StatementBegin
ALTER TABLE staffs
ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE staffs
DROP COLUMN id;
-- +goose StatementEnd
