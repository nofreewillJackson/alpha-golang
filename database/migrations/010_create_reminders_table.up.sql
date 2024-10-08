-- 010_create_reminders_table.up.sql
CREATE TABLE IF NOT EXISTS reminders (
                                         id SERIAL PRIMARY KEY,
                                         description TEXT NOT NULL,
                                         author_id VARCHAR(50) NOT NULL,
    due_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
