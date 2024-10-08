-- 009_create_logs_table.up.sql
CREATE TABLE IF NOT EXISTS logs (
                                    id SERIAL PRIMARY KEY,
                                    entry TEXT NOT NULL,
                                    tags TEXT,
                                    author_id VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
