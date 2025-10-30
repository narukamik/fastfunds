-- PostgreSQL schema for FastFunds

CREATE TABLE accounts (
    account_id INTEGER PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0 -- pennies
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    source_account_id INTEGER NOT NULL REFERENCES accounts(account_id) ON DELETE RESTRICT,
    destination_account_id INTEGER NOT NULL REFERENCES accounts(account_id) ON DELETE RESTRICT,
    amount BIGINT NOT NULL, -- pennies
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_source ON transactions(source_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_destination ON transactions(destination_account_id);

-- Seed data

INSERT INTO accounts (account_id, balance) VALUES
    (123,  10023),
    (456,  5000);

INSERT INTO transactions (source_account_id, destination_account_id, amount, status)
VALUES (123, 456, 1000, 'completed');
