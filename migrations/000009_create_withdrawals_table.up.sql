-- Создаем таблицу для истории списаний средств
CREATE TABLE withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    order_number VARCHAR(255) NOT NULL,
    sum DECIMAL(10, 2) NOT NULL,
    processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_withdrawals_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_withdrawals_user_id ON withdrawals(user_id);
CREATE INDEX idx_withdrawals_processed_at ON withdrawals(processed_at DESC);
CREATE INDEX idx_withdrawals_order_number ON withdrawals(order_number);

