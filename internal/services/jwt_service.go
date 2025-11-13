package services

import (
	"strings"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source=jwt_service.go -destination=mock_jwt_service.go -package=services
type JWTService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (models.JWTClaims, error)
	ExtractToken(c *fiber.Ctx) (string, error)
}

type JWTServiceImpl struct {
	jwtSecret []byte
	logger    logger.Logger
}

func NewJWTService(jwtSecret string, logger logger.Logger) JWTService {
	return &JWTServiceImpl{jwtSecret: []byte(jwtSecret), logger: logger}
}

func (s *JWTServiceImpl) ExtractToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get(fiber.HeaderAuthorization)
	if authHeader == "" {
		s.logger.Error("Authorization header required")
		return "", fiber.NewError(fiber.StatusUnauthorized, "Authorization header required")
	}

	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		s.logger.Error("Invalid authorization format")
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
	}

	return tokenParts[1], nil
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (models.JWTClaims, error) {
	var claims models.JWTClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Error("Invalid signing method")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return models.EmptyJWTClaims, err
	}

	if !token.Valid {
		s.logger.Error("Invalid token")
		return models.EmptyJWTClaims, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	newClaims, ok := (token.Claims).(*models.JWTClaims)
	if !ok {
		s.logger.Error("Invalid token claims")
		return models.EmptyJWTClaims, fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	return *newClaims, nil
}

func (s *JWTServiceImpl) GenerateToken(userID int) (string, error) {
	claims := &models.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
