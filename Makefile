DB_URI = postgresql://postgres:postgres@localhost:5432/gophermart
MIGRATIONS_DIR = migrations
GO_BIN = go
PSQL_BIN = psql

migrate-up:
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000001_create_users_table.up.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000002_create_orders_table.up.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000003_create_balance_table.up.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000004_create_user_login_index.up.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000005_alter_orders_id_to_bigint.up.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000006_alter_orders_add_order_number_and_change_id_to_serial.up.sql

migrate-down:
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000001_create_users_table.down.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000002_create_orders_table.down.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000003_create_balance_table.down.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000004_create_user_login_index.down.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000005_alter_orders_id_to_bigint.down.sql
	$(PSQL_BIN) $(DB_URI) -f $(MIGRATIONS_DIR)/000006_alter_orders_add_order_number_and_change_id_to_serial.down.sql

migrate-reset:
	make migrate-down
	make migrate-up

gen_mocks:
	mockgen -source=internal/services/jwt_service.go -destination=internal/services/mock_jwt_service.go -package=services
	mockgen -source=internal/services/order_service.go -destination=internal/services/mock_order_service.go -package=services
	mockgen -source=internal/services/user_service.go -destination=internal/services/mock_user_service.go -package=services
	mockgen -source=internal/services/balance_service.go -destination=internal/services/mock_balance_service.go -package=services
	mockgen -source=internal/repositories/order_repository.go -destination=internal/repositories/mock_order_repository.go -package=repositories
	mockgen -source=internal/repositories/user_repository.go -destination=internal/repositories/mock_user_repository.go -package=repositories
	mockgen -source=internal/services/external_accrual_service.go -destination=internal/services/mock_external_accrual_service.go -package=services
	
test:
	make gen_mocks
	$(GO_BIN) test -v ./...

run:
	make migrate-reset
	make test
	$(GO_BIN) run cmd/gophermart/main.go

build:
	make migrate-reset
	make test
	$(GO_BIN) build -o gophermart cmd/gophermart/main.go
