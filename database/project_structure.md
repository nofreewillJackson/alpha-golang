# Project Directory Structure

```
database/
├── migrations/
    ├── 001_create_messages_table.up.sql
    ├── 002_create_summaries_table.up.sql
    ├── 003_create_digests_table.up.sql
    ├── 004_add_summarized_column_to_messages.up.sql
    ├── 005_add_digested.up.sql
    ├── 006_add_synthesized.up.sql
    ├── 007_add_synthesizedTEXT.up.sql
    └── README.md
├── printProject.py
└── project_structure.md
```

# File Contents

## `C:\\Dev\\alpha-golang\\database\\migrations\\001_create_messages_table.up.sql`
```sql
CREATE TABLE IF NOT EXISTS messages (
                                        id SERIAL PRIMARY KEY,
                                        content TEXT NOT NULL,
                                        author_id VARCHAR(50) NOT NULL,
    channel_id VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    );

```

## `C:\\Dev\\alpha-golang\\database\\migrations\\002_create_summaries_table.up.sql`
```sql
CREATE TABLE IF NOT EXISTS summaries (
                                         id SERIAL PRIMARY KEY,
                                         summary TEXT NOT NULL,
                                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

```

## `C:\\Dev\\alpha-golang\\database\\migrations\\003_create_digests_table.up.sql`
```sql
CREATE TABLE IF NOT EXISTS digests (
                                       id SERIAL PRIMARY KEY,
                                       digest TEXT NOT NULL,
                                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

```

## `C:\\Dev\\alpha-golang\\database\\migrations\\004_add_summarized_column_to_messages.up.sql`
```sql
ALTER TABLE messages ADD COLUMN summarized BOOLEAN DEFAULT FALSE;

```

## `C:\\Dev\\alpha-golang\\database\\migrations\\005_add_digested.up.sql`
```sql
-- 004_add_digested_column_to_messages.up.sql
ALTER TABLE messages ADD COLUMN digested BOOLEAN DEFAULT FALSE;

```

## `C:\\Dev\\alpha-golang\\database\\migrations\\006_add_synthesized.up.sql`
```sql
-- 006
ALTER TABLE messages ADD COLUMN synthesized BOOLEAN DEFAULT FALSE;
```

## `C:\\Dev\\alpha-golang\\database\\migrations\\007_add_synthesizedTEXT.up.sql`
```sql
-- 007
ALTER TABLE messages ADD COLUMN synthesis TEXT;
```

