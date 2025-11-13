-- Удаляем таблицу balance_transactions
DROP INDEX IF EXISTS idx_balance_transactions_user_order_accrual;
DROP INDEX IF EXISTS idx_balance_transactions_type;
DROP INDEX IF EXISTS idx_balance_transactions_created_at;
DROP INDEX IF EXISTS idx_balance_transactions_order_id;
DROP INDEX IF EXISTS idx_balance_transactions_user_id;
DROP TABLE IF EXISTS balance_transactions;

