CREATE TABLE IF NOT EXISTS collections (
    id varchar(255) PRIMARY KEY,
    title text NOT NULL DEFAULT '',
    username text NOT NULL DEFAULT '',
    created_at timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS collection_documents (
    collection_id varchar(255) NOT NULL,
    document_id varchar(255) NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (collection_id, position)
);
