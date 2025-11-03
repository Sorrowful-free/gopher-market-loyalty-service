CREATE TABLE balance(
    id SERIAL PRIMARY KEY,
    user_id INT,
    order_id INT,
    sum DECIMAL(10, 2),
    processed_at TIMESTAMP
);

CREATE INDEX idx_balance_id ON balance(id);
ALTER TABLE balance ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE balance ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id);
