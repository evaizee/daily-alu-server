-- Drop triggers first
DROP TRIGGER IF EXISTS update_api_keys_updated_at ON api_keys;

-- Drop functions
DROP FUNCTION IF EXISTS update_api_keys_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_api_keys_key_value;
DROP INDEX IF EXISTS idx_api_keys_status;
DROP INDEX IF EXISTS idx_api_keys_expires_at;

-- Drop the table
DROP TABLE IF EXISTS api_keys;
