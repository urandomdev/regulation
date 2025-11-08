-- Modify "sync_cursors" table
ALTER TABLE "sync_cursors" ADD COLUMN "last_error" character varying NOT NULL DEFAULT '', ADD COLUMN "consecutive_failures" bigint NOT NULL DEFAULT 0;
