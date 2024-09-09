-- +goose Up
-- +goose StatementBegin
ALTER TABLE participants DROP CONSTRAINT IF EXISTS participants_barcode_event_id_key;

ALTER TABLE participants DROP CONSTRAINT IF EXISTS participants_pkey;

ALTER TABLE participants ADD PRIMARY KEY (barcode, event_id);

ALTER TABLE participants DROP COLUMN IF EXISTS id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE participants ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();

ALTER TABLE participants DROP CONSTRAINT IF EXISTS participants_pkey;

ALTER TABLE participants ADD CONSTRAINT participants_barcode_event_id_key UNIQUE (barcode, event_id);
-- +goose StatementEnd
