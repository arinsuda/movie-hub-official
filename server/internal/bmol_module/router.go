package bmol_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	bmol := router.Group("/bmol")
	bmol.Post("", h.AddItem)
	bmol.Get("/user/:userId", h.GetUserBMOL)
	bmol.Put("/:itemId", h.UpdateItem)
	bmol.Delete("/:itemId", h.RemoveItem)
}
