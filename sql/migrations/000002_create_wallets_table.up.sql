CREATE TABLE IF NOT EXISTS wallets(
   wallet_id UUID NOT NULL DEFAULT gen_random_uuid(),
   user_id UUID NOT NULL REFERENCES users(user_id),
   balance NUMERIC(12, 2) DEFAULT 0,
   created_at timestamp DEFAULT now()
);

CREATE INDEX wallet_user_id ON wallets(user_id);