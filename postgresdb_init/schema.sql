CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS users (
    user_id UUID NOT NULL DEFAULT (uuid_generate_v4()),
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    slug VARCHAR(10) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL,
    UNIQUE (user_id, url),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
	        REFERENCES users(user_id)
	        ON DELETE CASCADE
);

CREATE INDEX trgm_idx_links_description ON links USING gin (description gin_trgm_ops);
