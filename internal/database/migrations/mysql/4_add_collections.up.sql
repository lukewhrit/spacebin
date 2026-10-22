CREATE TABLE IF NOT EXISTS collections (
	id varchar(255) PRIMARY KEY,
	title varchar(255) NOT NULL DEFAULT '',
	username varchar(255) NOT NULL DEFAULT '',
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collection_documents (
	collection_id varchar(255) NOT NULL,
	document_id varchar(255) NOT NULL,
	position integer NOT NULL,
	PRIMARY KEY (collection_id, position)
);
