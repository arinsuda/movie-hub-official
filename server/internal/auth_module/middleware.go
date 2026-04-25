package auth_module

import (
	"github.com/arinsuda/movie-hub/internal/config"
	"github.com/gofiber/fiber/v3"
)

type Middleware struct {
	jwt *jwtManager
}

func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{jwt: newJWTManager(cfg.JWT)}
}

// RequireAuth บังคับ login — ใส่ claims ใน locals
func (m *Middleware) RequireAuth(c fiber.Ctx) error {
	claims, err := m.extractClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	c.Locals("claims", claims)
	return c.Next()
}

// RequireRole บังคับ role เฉพาะ
func (m *Middleware) RequireRole(role string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*Claims)
		if !ok || claims.Role != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}

func (m *Middleware) extractClaims(c fiber.Ctx) (*Claims, error) {
	token := c.Cookies("access_token")
	if token == "" {
		return nil, ErrInvalidToken
	}
	return m.jwt.ParseAccess(token)
}

// GetClaims ดึง claims จาก context (helper สำหรับ handler)
func GetClaims(c fiber.Ctx) *Claims {
	claims, _ := c.Locals("claims").(*Claims)
	return claims
}
