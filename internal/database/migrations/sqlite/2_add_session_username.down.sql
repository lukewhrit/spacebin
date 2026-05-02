CREATE TABLE sessions_backup (
    public TEXT PRIMARY KEY,
    token TEXT NOT NULL,
    secret TEXT NOT NULL
);

INSERT INTO sessions_backup (public, token, secret)
SELECT public, token, secret FROM sessions;

DROP TABLE sessions;

ALTER TABLE sessions_backup RENAME TO sessions;
