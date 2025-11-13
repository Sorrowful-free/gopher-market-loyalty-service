ALTER TABLE balance DROP CONSTRAINT IF EXISTS fk_order_id;

CREATE SEQUENCE IF NOT EXISTS orders_id_seq;

SELECT setval('orders_id_seq', COALESCE((SELECT MAX(id) FROM orders), 0) + 1, false);

ALTER TABLE orders ALTER COLUMN id TYPE INTEGER;
ALTER TABLE orders ALTER COLUMN id SET DEFAULT nextval('orders_id_seq');
ALTER TABLE orders ALTER COLUMN id SET NOT NULL;

ALTER SEQUENCE orders_id_seq OWNED BY orders.id;

ALTER TABLE orders ADD COLUMN order_number VARCHAR(255);

ALTER TABLE balance ALTER COLUMN order_id TYPE INTEGER;
ALTER TABLE balance ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id);
