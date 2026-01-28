-- Drop indexes first
DROP INDEX IF EXISTS idx_redirect_uris_client_id;
DROP INDEX IF EXISTS idx_auth_codes_expires;
DROP INDEX IF EXISTS idx_auth_codes_client_id;

-- Remove refresh token expiry column
ALTER TABLE oauth_tokens DROP COLUMN IF EXISTS refresh_token_expires_at;

-- Drop tables
DROP TABLE IF EXISTS client_redirect_uris;
DROP TABLE IF EXISTS authorization_codes;
