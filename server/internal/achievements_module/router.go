package achievementsmodule

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Module struct {
	Service Service
}

func New(db *gorm.DB) *Module {
	repo := newRepository(db)
	svc := newService(repo)
	return &Module{Service: svc}
}

func (m *Module) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	h := newHandler(m.Service)

	router.Get("/achievements", h.listAllAchievements)

	router.Get("/users/:userID/achievements", authMiddleware, h.listUserAchievements)
}
