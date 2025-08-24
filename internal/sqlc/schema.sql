CREATE TABLE admins (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email VARCHAR(255) NOT NULL,
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
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

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
	timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	event_id UUID,
	user_code VARCHAR(10) NOT NULL,

	PRIMARY KEY(barcode,event_id),
	FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
	FOREIGN KEY(user_code) REFERENCES users(code) ON DELETE CASCADE
);

CREATE TABLE refresh_tokens (
	email VARCHAR(255) PRIMARY KEY,
	token VARCHAR(255) NOT NULL
);

CREATE TABLE users (
	code VARCHAR(10) PRIMARY KEY,
	full_name TEXT NOT NULL,
	gmail TEXT NOT NULL,
	major VARCHAR(6) NOT NULL
)
