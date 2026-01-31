CREATE TABLE IF NOT EXISTS payment_intents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL,
    description TEXT,
    user_id UUID NOT NULL, -- Logical foreign key to Auth service User
    application_fee_amount BIGINT DEFAULT 0,
    on_behalf_of UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for looking up payments by user
CREATE INDEX IF NOT EXISTS idx_payment_intents_user_id ON payment_intents(user_id);

CREATE TABLE IF NOT EXISTS idempotency_keys (
    user_id UUID NOT NULL,
    key VARCHAR(255) NOT NULL,
    response_body TEXT NOT NULL,
    status_code INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, key)
);
