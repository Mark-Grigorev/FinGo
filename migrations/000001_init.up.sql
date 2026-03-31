CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL DEFAULT '',
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS accounts (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    type       VARCHAR(50) NOT NULL DEFAULT 'card',
    currency   VARCHAR(10) NOT NULL DEFAULT 'RUB',
    balance    NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);

CREATE TABLE IF NOT EXISTS categories (
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT REFERENCES users(id) ON DELETE CASCADE,
    name      VARCHAR(255) NOT NULL,
    icon      VARCHAR(20) NOT NULL DEFAULT '💳',
    type      VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(user_id, name, type)
);
CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);

CREATE TABLE IF NOT EXISTS transactions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id  BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    type        VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    amount      NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    name        VARCHAR(500) NOT NULL DEFAULT '',
    date        DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_account  ON transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date     ON transactions(date DESC);

CREATE TABLE IF NOT EXISTS budgets (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id   BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    month         DATE NOT NULL,
    limit_amount  NUMERIC(15,2) NOT NULL CHECK (limit_amount > 0),
    spent_amount  NUMERIC(15,2) NOT NULL DEFAULT 0,
    UNIQUE(user_id, category_id, month)
);

CREATE TABLE IF NOT EXISTS recurring_payments (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id        BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id       BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    name              VARCHAR(255) NOT NULL,
    amount            NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    frequency         VARCHAR(20) NOT NULL DEFAULT 'monthly',
    next_payment_date DATE NOT NULL,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
