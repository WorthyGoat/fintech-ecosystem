-- Add scopes column to api_keys
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS scopes TEXT DEFAULT '*';
