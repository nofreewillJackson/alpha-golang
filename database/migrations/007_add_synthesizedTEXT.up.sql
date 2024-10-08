-- 007_add_synthesis_column.up.sql
ALTER TABLE messages ADD COLUMN IF NOT EXISTS synthesis TEXT;
