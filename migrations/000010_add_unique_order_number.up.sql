-- Добавляем UNIQUE ограничение на order_number в таблице orders
-- Сначала удаляем возможные дубликаты (если есть)
-- Затем добавляем UNIQUE индекс
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_order_number_unique ON orders(order_number);

