CREATE UNIQUE INDEX idx_wallets_user_id_active 
ON wallets (user_id) 
WHERE deleted_at IS NULL;

CREATE INDEX idx_transactions_wallet_id ON transactions (wallet_id);

CREATE INDEX idx_transactions_wallet_date 
ON transactions (wallet_id, created_at DESC);