mockgen -source=internal/services/jwt_service.go -destination=internal/services/mock_jwt_service.go -package=services
mockgen -source=internal/services/order_service.go -destination=internal/services/mock_order_service.go -package=services
mockgen -source=internal/services/user_service.go -destination=internal/services/mock_user_service.go -package=services
mockgen -source=internal/services/balance_service.go -destination=internal/services/mock_balance_service.go -package=services
mockgen -source=internal/repositories/order_repository.go -destination=internal/repositories/mock_order_repository.go -package=repositories
mockgen -source=internal/repositories/user_repository.go -destination=internal/repositories/mock_user_repository.go -package=repositories