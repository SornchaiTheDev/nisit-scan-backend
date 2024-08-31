-- +goose Up
-- +goose StatementBegin
CREATE TABLE admins (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email VARCHAR(255) UNIQUE NOT NULL,
	full_name VARCHAR(255) NOT NULL,
	deleted_at TIMESTAMP
);

CREATE TABLE events (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	place VARCHAR(255) NOT NULL,
	date DATE NOT NULL,
	host VARCHAR(255) NOT NULL,
	admin_id UUID,

	UNIQUE(name,place,date),
	FOREIGN KEY(admin_id) REFERENCES admins(id) ON DELETE CASCADE
);

CREATE TABLE staffs (
	email VARCHAR(255) NOT NULL,
	event_id UUID,

	UNIQUE (email, event_id),
	FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE
);

CREATE TABLE participants (
	barcode VARCHAR(14) NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	event_id UUID,

	UNIQUE(barcode,event_id),
	FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE participants;
DROP TABLE staffs;
DROP TABLE events;
DROP TABLE admins;
-- +goose StatementEnd
