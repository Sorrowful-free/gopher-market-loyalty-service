-- Создаем таблицу для истории всех операций с балансом
CREATE TABLE balance_transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    order_id INTEGER, -- NULL для withdrawals
    transaction_type VARCHAR(16) NOT NULL, -- 'ACCRUAL' или 'WITHDRAWAL'
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_balance_transactions_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_balance_transactions_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL,
    CONSTRAINT chk_transaction_type CHECK (transaction_type IN ('ACCRUAL', 'WITHDRAWAL'))
);

CREATE INDEX idx_balance_transactions_user_id ON balance_transactions(user_id);
CREATE INDEX idx_balance_transactions_order_id ON balance_transactions(order_id);
CREATE INDEX idx_balance_transactions_created_at ON balance_transactions(created_at DESC);
CREATE INDEX idx_balance_transactions_type ON balance_transactions(transaction_type);
-- Уникальный индекс для предотвращения дубликатов начислений за один заказ
CREATE UNIQUE INDEX idx_balance_transactions_user_order_accrual ON balance_transactions(user_id, order_id, transaction_type) 
WHERE transaction_type = 'ACCRUAL' AND order_id IS NOT NULL;

