-- Create users table
CREATE TABLE IF NOT EXISTS activities (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    child_id BIGINT NOT NULL,
    type VARCHAR(255) NOT NULL,
    details JSONB,
    happens_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create index on user_id
CREATE INDEX IF NOT EXISTS idx_user_id ON activities(user_id);
CREATE INDEX IF NOT EXISTS idx_child_id ON activities(child_id);
