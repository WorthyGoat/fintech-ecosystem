-- Remove scopes column from api_keys
ALTER TABLE api_keys DROP COLUMN IF EXISTS scopes;
