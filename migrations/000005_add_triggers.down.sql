DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP TRIGGER IF EXISTS set_timestamp_wallets ON wallets;

DROP FUNCTION IF EXISTS update_modified_column();