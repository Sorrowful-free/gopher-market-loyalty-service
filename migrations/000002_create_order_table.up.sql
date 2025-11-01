CREATE TABLE order(
    order_id INT PRIMARY KEY,
    user_id INT,
    status VARCHAR(16),
    accrual DECIMAL(10, 2),
    uploaded_at TIMESTAMP
)

CREATE INDEX idx_order_id ON orders(order_id);
ALTER TABLE orders ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id);