-- +goose Up

CREATE TABLE IF NOT EXISTS wallets (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    last_operation_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

INSERT INTO wallets (wallet_id, balance, last_operation_type, created_at) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1000.50, 'deposit', NOW() - INTERVAL '10 days'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 250.75, 'withdraw', NOW() - INTERVAL '8 days'),
('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 5000.00, 'deposit', NOW() - INTERVAL '6 days'),
('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 150.30, 'deposit', NOW() - INTERVAL '4 days'),
('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 75.20, 'withdraw', NOW() - INTERVAL '2 days'),
('f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 3000.00, 'deposit', NOW() - INTERVAL '1 day'),
(gen_random_uuid(), 10.00, 'deposit', NOW());

CREATE INDEX idx_wallet_id ON wallets(wallet_id);

-- +goose Down
DROP TABLE IF EXISTS wallets;
