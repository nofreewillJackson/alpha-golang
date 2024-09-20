-- 004_add_digested_column_to_messages.up.sql
ALTER TABLE messages ADD COLUMN digested BOOLEAN DEFAULT FALSE;
