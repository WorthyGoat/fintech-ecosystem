-- Add user_id to idempotency_keys
ALTER TABLE idempotency_keys ADD COLUMN user_id UUID;
-- We can't strictly enforce NOT NULL if there are existing rows without data, 
-- but in this project we can probably assume it's fresh or we can backfill.
-- For safety in a real migration:
-- ALTER TABLE idempotency_keys ALTER COLUMN user_id SET NOT NULL;

-- Update Primary Key to include user_id
ALTER TABLE idempotency_keys DROP CONSTRAINT idempotency_keys_pkey;
ALTER TABLE idempotency_keys ADD PRIMARY KEY (user_id, key);
