CREATE TABLE orders(
    id INT PRIMARY KEY,
    user_id INT,
    status VARCHAR(16),
    accrual DECIMAL(10, 2),
    uploaded_at TIMESTAMP
)

CREATE INDEX idx_order_id ON orders(id);
ALTER TABLE orders ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);