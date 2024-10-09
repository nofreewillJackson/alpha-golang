-- 011_update_logs_table.up.sql

-- Convert existing tags to valid JSON
UPDATE logs
SET tags = '["' || REPLACE(tags, ', ', '", "') || '"]'
WHERE tags IS NOT NULL AND tags != '';

-- Alter the tags column to JSONB
ALTER TABLE logs
ALTER COLUMN tags TYPE JSONB USING tags::JSONB;
