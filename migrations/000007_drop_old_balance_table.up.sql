-- Удаляем старую таблицу balance
ALTER TABLE balance DROP CONSTRAINT IF EXISTS fk_order_id;
ALTER TABLE balance DROP CONSTRAINT IF EXISTS fk_user_id;
DROP INDEX IF EXISTS idx_balance_id;
DROP TABLE IF EXISTS balance;

