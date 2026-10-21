CREATE TABLE IF NOT EXISTS collections (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collection_documents (
    collection_id TEXT NOT NULL,
    document_id TEXT NOT NULL,
    position INTEGER NOT NULL,
    PRIMARY KEY (collection_id, position)
);
