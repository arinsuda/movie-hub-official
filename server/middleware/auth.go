package middleware

import (
	"errors"

	"github.com/arinsuda/movie-hub/config"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type AuthMiddleware struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (m *AuthMiddleware) RequireAuth(c fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	claims, err := ParseAccess(token, m.cfg.JWT.AccessSecret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	c.Locals("claims", claims)
	return c.Next()
}

func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*Claims)
		if !ok || claims.Role != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}

func (m *AuthMiddleware) RequireCurrentAdmin(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*Claims)
		if !ok || claims.UserID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		var count int64
		err := db.Table("users").
			Joins("JOIN roles ON roles.id = users.role_id").
			Where("users.id = ? AND users.deleted_at IS NULL AND users.is_active = true AND roles.role_name = 'admin'", claims.UserID).
			Count(&count).Error

		if err != nil || count == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin authorization lost"})
		}

		return c.Next()
	}
}
