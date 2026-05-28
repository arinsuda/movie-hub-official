package analytics_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// RegisterRoutes ลงทะเบียน analytics routes
//
// Routes (ทั้งหมดเป็น public — ไม่ต้อง auth):
//
//	GET /analytics/:mediaType/:mediaId  → stats ของ media นั้น
//	GET /analytics/:mediaType/trending  → top 20 trending
//
// :mediaType = "movies" | "series"
func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	analytics := router.Group("/analytics")

	// ⚠️  route จำเพาะ ("trending") ต้องขึ้นก่อน wildcard (":mediaId")
	// เพื่อกัน fiber match "trending" เป็น mediaId
	analytics.Get("/:mediaType/trending", h.GetTrending)
	analytics.Get("/:mediaType/:mediaId", h.GetMediaAnalytics)
}
