CREATE TABLE IF NOT EXISTS documents (
	id varchar(255) PRIMARY KEY,
	content text NOT NULL,
	created_at timestamp with time zone DEFAULT now(),
	updated_at timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts (
	id SERIAL PRIMARY KEY,
	username varchar(255) NOT NULL,
	password varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
	public varchar(255) PRIMARY KEY,
	token varchar(255) NOT NULL,
	secret varchar
);
