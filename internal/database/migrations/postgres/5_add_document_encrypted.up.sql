ALTER TABLE documents ADD COLUMN IF NOT EXISTS encrypted boolean NOT NULL DEFAULT false;
