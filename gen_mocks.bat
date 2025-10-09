mockgen -source=internal/services/jwt_service.go -destination=mocks/mock_jwt_service.go -package=mocks
mockgen -source=internal/services/order_service.go -destination=mocks/mock_order_service.go -package=mocks
mockgen -source=internal/services/user_service.go -destination=mocks/mock_user_service.go -package=mocks
mockgen -source=internal/repositories/order_repository.go -destination=mocks/mock_order_repository.go -package=mocks
mockgen -source=internal/repositories/user_repository.go -destination=mocks/mock_user_repository.go -package=mocks