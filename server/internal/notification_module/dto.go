package notification_module

import "time"

// ─── Request DTOs ─────────────────────────────────────────────────────────────

// CreateNotificationRequest ใช้ภายใน service-to-service เท่านั้น (ไม่ expose ผ่าน HTTP)
type CreateNotificationRequest struct {
	UserID    uint  // เจ้าของ notification
	ActorID   *uint // ผู้กระทำ (nil = system)
	Type      NotificationType
	TargetID  *uint   // entity id ที่เกี่ยวข้อง
	TargetRef *string // "movie" | "review" | "user"
	Message   string
}

// MarkReadRequest ใช้กับ endpoint PATCH /notifications/read
type MarkReadRequest struct {
	// IDs ที่ต้องการ mark as read; ถ้า empty = mark ทั้งหมด
	IDs []uint `json:"ids"`
}

// ListNotificationsQuery query params สำหรับ GET /notifications
type ListNotificationsQuery struct {
	Unread   *bool  `query:"unread"`    // filter เฉพาะที่ยังไม่อ่าน
	Page     int    `query:"page"`      // default 1
	PageSize int    `query:"page_size"` // default 20, max 100
	Type     string `query:"type"`      // filter by NotificationType
}

// ─── Response DTOs ────────────────────────────────────────────────────────────

// ActorSummary ข้อมูลย่อของผู้กระทำ (embed ใน response)
type ActorSummary struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}

// NotificationResponse รูปแบบที่ client รับ
type NotificationResponse struct {
	ID        uint             `json:"id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	IsRead    bool             `json:"is_read"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	Actor     *ActorSummary    `json:"actor,omitempty"`     // nil ถ้าเป็น system notification
	TargetID  *uint            `json:"target_id,omitempty"` // id สำหรับ deep-link
	TargetRef *string          `json:"target_ref,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

// NotificationListResponse ผลลัพธ์รายการพร้อม pagination
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                  `json:"unread_count"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
}
