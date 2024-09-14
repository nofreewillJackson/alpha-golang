CREATE TABLE IF NOT EXISTS summaries (
                                         id SERIAL PRIMARY KEY,
                                         summary TEXT NOT NULL,
                                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
