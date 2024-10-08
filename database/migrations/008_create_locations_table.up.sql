-- 008_create_locations_table.up.sql
CREATE TABLE IF NOT EXISTS locations (
                                         id SERIAL PRIMARY KEY,
                                         item TEXT UNIQUE NOT NULL,
                                         description TEXT NOT NULL
);
