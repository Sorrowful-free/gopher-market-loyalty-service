ALTER TABLE balance DROP CONSTRAINT fk_order_id;
ALTER TABLE balance DROP CONSTRAINT fk_user_id;
DROP INDEX idx_balance_id;
DROP TABLE balance;