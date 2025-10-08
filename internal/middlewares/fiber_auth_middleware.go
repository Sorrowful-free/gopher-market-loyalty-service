package middlewares

import (
	"strings"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey = "user_id"
	LoginKey  = "login"

	AuthorizationKey = "Authorization"
)

type JWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.Claims
}

type FiberAuthMiddleware struct {
	jwtSecret []byte
	logger    logger.Logger
}

func NewFiberAuthMiddleware(jwtSecret string, logger logger.Logger) *FiberAuthMiddleware {
	return &FiberAuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

func (m *FiberAuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	token, err := m.extractToken(c)
	if err != nil {
		m.logger.Error("Failed to extract token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	claims, err := m.validateToken(token)
	if err != nil {
		m.logger.Error("Failed to validate token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	c.Locals(UserIDKey, claims.UserID)

	return c.Next()
}

func (m *FiberAuthMiddleware) extractToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get(AuthorizationKey)
	if authHeader == "" {
		m.logger.Error("Authorization header required")
		return "", fiber.NewError(fiber.StatusUnauthorized, "Authorization header required")
	}

	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		m.logger.Error("Invalid authorization format")
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
	}

	return tokenParts[1], nil
}

func (m *FiberAuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			m.logger.Error("Invalid signing method")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return m.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		m.logger.Error("Invalid token")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		m.logger.Error("Invalid token claims")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	return claims, nil
}

func (m *FiberAuthMiddleware) generateToken(userID int64) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.jwtSecret)
}

func GetUserIDFromContext(c *fiber.Ctx) (int64, error) {
	userID, ok := c.Locals(UserIDKey).(int64)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}
	return userID, nil
}

func GetLoginFromContext(c *fiber.Ctx) (string, error) {
	login, ok := c.Locals(LoginKey).(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}
	return login, nil
}
