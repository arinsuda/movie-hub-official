package user_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	users := router.Group("/users")

	users.Get("/:userId", h.GetProfile)
	users.Patch("/:userId", h.UpdateProfile)
	users.Delete("/:userId", h.DeleteUser)
}
