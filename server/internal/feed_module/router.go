package feed_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Module struct {
	Service Service
}

func New(db *gorm.DB) *Module {
	return &Module{Service: NewService(db)}
}

// RegisterRoutes ผูก route ทั้งหมดของ feed_module เข้ากับ router
// ทุก route ต้อง login ก่อน (เหมือน notification_module / review_module / library_module)
// ดังนั้นควรเรียกด้วย router ที่ผ่าน auth middleware มาแล้ว (เช่น `protected` group ใน router/router.go)
func (m *Module) RegisterRoutes(router fiber.Router) {
	h := newHandler(m.Service)

	router.Get("/feed", h.GetFeed)
	router.Get("/users/:userId/activities", h.GetUserActivities)
	router.Patch("/activities/:activityId/visibility", h.UpdateVisibility)
	router.Get("/me/activity-settings", h.GetSettings)
	router.Patch("/me/activity-settings", h.UpdateSettings)
}
