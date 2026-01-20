CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- asset, liability, equity, revenue, expense
    balance BIGINT NOT NULL DEFAULT 0,
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
    amount BIGINT NOT NULL, -- Positive for Debit, Negative for Credit
    direction VARCHAR(10) NOT NULL, -- 'debit' or 'credit' (redundant but useful for validation/readability)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_entries_transaction_id ON entries(transaction_id);
CREATE INDEX IF NOT EXISTS idx_entries_account_id ON entries(account_id);
