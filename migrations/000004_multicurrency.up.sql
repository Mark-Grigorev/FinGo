ALTER TABLE users ADD COLUMN IF NOT EXISTS base_currency VARCHAR(10) NOT NULL DEFAULT 'RUB';

CREATE TABLE IF NOT EXISTS exchange_rates (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    currency   VARCHAR(10) NOT NULL,
    rate       NUMERIC(18,6) NOT NULL CHECK (rate > 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, currency)
);
CREATE INDEX IF NOT EXISTS idx_exchange_rates_user ON exchange_rates(user_id);
