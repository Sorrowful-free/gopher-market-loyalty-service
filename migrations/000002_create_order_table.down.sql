ALTER TABLE orders DROP CONSTRAINT fk_user_id;
DROP INDEX idx_order_id;
DROP TABLE orders;