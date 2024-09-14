CREATE TABLE IF NOT EXISTS digests (
                                       id SERIAL PRIMARY KEY,
                                       digest TEXT NOT NULL,
                                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
