-- Authorization codes table with PKCE support
CREATE TABLE IF NOT EXISTS authorization_codes (
    code VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255) REFERENCES oauth_clients(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    redirect_uri TEXT NOT NULL,
    scope TEXT,
    code_challenge VARCHAR(255),
    code_challenge_method VARCHAR(10),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Client redirect URIs allowlist
CREATE TABLE IF NOT EXISTS client_redirect_uris (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id VARCHAR(255) REFERENCES oauth_clients(id) ON DELETE CASCADE,
    redirect_uri TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(client_id, redirect_uri)
);

-- Extend oauth_tokens with refresh token expiry
ALTER TABLE oauth_tokens ADD COLUMN IF NOT EXISTS refresh_token_expires_at TIMESTAMP WITH TIME ZONE;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_auth_codes_client_id ON authorization_codes(client_id);
CREATE INDEX IF NOT EXISTS idx_auth_codes_expires ON authorization_codes(expires_at);
CREATE INDEX IF NOT EXISTS idx_redirect_uris_client_id ON client_redirect_uris(client_id);
