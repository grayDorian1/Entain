-- Schema
CREATE SCHEMA IF NOT EXISTS accounts;
CREATE SCHEMA IF NOT EXISTS payments;

-- Users table
CREATE TABLE IF NOT EXISTS accounts.users (
    id      BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    balance NUMERIC(18, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0)
);

-- Transactions table
CREATE TABLE IF NOT EXISTS payments.transactions (
    id             BIGSERIAL PRIMARY KEY,
    transaction_id TEXT NOT NULL UNIQUE,
    user_id        BIGINT NOT NULL REFERENCES accounts.users(id),
    source_type    VARCHAR(50) NOT NULL,
    state          VARCHAR(10) NOT NULL,
    amount         NUMERIC(18, 2) NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW()
);

-- Predefined users
INSERT INTO accounts.users (id, balance)
OVERRIDING SYSTEM VALUE VALUES
    (1, 1000.00),
    (2, 1000.00),
    (3, 1000.00)
ON CONFLICT (id) DO NOTHING;