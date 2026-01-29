CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- asset, liability, equity, revenue, expense
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    user_id UUID, -- Optional link to auth user
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_id VARCHAR(255) UNIQUE NOT NULL, -- Idempotency key / Reference to external event
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    account_id UUID NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL CHECK (amount != 0), -- Must be positive for Debit, Negative for Credit
    direction VARCHAR(10) NOT NULL, -- 'debit' or 'credit'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE -- NULL means pending
);

-- TRIGGERS FOR IMMUTABILITY
CREATE OR REPLACE FUNCTION prevent_mutation()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Financial records are immutable. Mutation of %% is forbidden.', TG_TABLE_NAME;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION prevent_outbox_mutation()
RETURNS TRIGGER AS $$
BEGIN
    -- Allow updating processed_at from NULL to a timestamp
    IF (OLD.processed_at IS NULL AND NEW.processed_at IS NOT NULL) THEN
        -- Ensure other columns haven't changed
        IF (OLD.id = NEW.id AND OLD.event_type = NEW.event_type AND OLD.payload = NEW.payload AND OLD.created_at = NEW.created_at) THEN
            RETURN NEW;
        END IF;
    END IF;
    RAISE EXCEPTION 'Immutable outbox violation. Only processed_at can be updated exactly once.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_immutable_transactions
BEFORE UPDATE OR DELETE ON transactions
FOR EACH ROW EXECUTE FUNCTION prevent_mutation();

CREATE TRIGGER trg_immutable_entries
BEFORE UPDATE OR DELETE ON entries
FOR EACH ROW EXECUTE FUNCTION prevent_mutation();

CREATE TRIGGER trg_immutable_outbox
BEFORE UPDATE OR DELETE ON outbox
FOR EACH ROW EXECUTE FUNCTION prevent_outbox_mutation();

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_entries_transaction_id ON entries(transaction_id);
CREATE INDEX IF NOT EXISTS idx_entries_account_id ON entries(account_id);
