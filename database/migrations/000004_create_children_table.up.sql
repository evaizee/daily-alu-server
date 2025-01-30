-- Create users table
CREATE TABLE IF NOT EXISTS children (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    details JSONB,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create index on user_id
CREATE INDEX IF NOT EXISTS idx_user_id ON children(user_id);
