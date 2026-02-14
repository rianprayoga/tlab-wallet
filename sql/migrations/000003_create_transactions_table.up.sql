CREATE TYPE transaction_status AS ENUM ('SUCCESS', 'FAILED');

CREATE TABLE IF NOT EXISTS transactions(
   transaction_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
   amount BIGINT NOT NULL,
   sender_id UUID NOT NULL REFERENCES users(user_id),
   receiver_id UUID NOT NULL REFERENCES users(user_id),
   status transaction_status NOT NULL,
   created_at timestamp DEFAULT now()
);

CREATE INDEX transactions_sender_id ON transactions(sender_id);