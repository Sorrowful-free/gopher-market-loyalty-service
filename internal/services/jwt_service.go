package services

import (
	"strings"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthorizationKey = "Authorization"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.Claims
}

type JWTService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (*JWTClaims, error)
	ExtractToken(c *fiber.Ctx) (string, error)
}

type JWTServiceImpl struct {
	jwtSecret string
	logger    logger.Logger
}

func NewJWTService(jwtSecret string, logger logger.Logger) JWTService {
	return &JWTServiceImpl{jwtSecret: jwtSecret, logger: logger}
}

func (s *JWTServiceImpl) ExtractToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get(AuthorizationKey)
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

func (s *JWTServiceImpl) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Error("Invalid signing method")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		s.logger.Error("Invalid token")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		s.logger.Error("Invalid token claims")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	return claims, nil
}

func (s *JWTServiceImpl) GenerateToken(userID string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
