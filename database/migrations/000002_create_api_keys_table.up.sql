CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    key_value VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    rate_limit INTEGER NOT NULL DEFAULT 1000,
    allowed_ips TEXT[] DEFAULT '{}',
    rotation_date TIMESTAMP WITH TIME ZONE,
    sunset_date TIMESTAMP WITH TIME ZONE,
    successor_key_id UUID REFERENCES api_keys(id),
    predecessor_key_id UUID REFERENCES api_keys(id),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_api_keys_key_value ON api_keys(key_value);
CREATE INDEX idx_api_keys_status ON api_keys(status);
CREATE INDEX idx_api_keys_expires_at ON api_keys(expires_at);

-- Add comments for better documentation
COMMENT ON TABLE api_keys IS 'Stores API keys for application authentication';
COMMENT ON COLUMN api_keys.id IS 'Unique identifier for the API key';
COMMENT ON COLUMN api_keys.name IS 'Human-readable name for the API key';
COMMENT ON COLUMN api_keys.key_value IS 'The actual API key value';
COMMENT ON COLUMN api_keys.status IS 'Current status of the API key (active, rotating, revoked)';
COMMENT ON COLUMN api_keys.rate_limit IS 'Maximum number of requests allowed per hour';
COMMENT ON COLUMN api_keys.allowed_ips IS 'List of allowed IP addresses/ranges';
COMMENT ON COLUMN api_keys.rotation_date IS 'Date when key rotation should begin';
COMMENT ON COLUMN api_keys.sunset_date IS 'Date when key will be fully revoked';
COMMENT ON COLUMN api_keys.successor_key_id IS 'Reference to the new key during rotation';
COMMENT ON COLUMN api_keys.predecessor_key_id IS 'Reference to the old key being replaced';
