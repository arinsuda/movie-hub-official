package notification_module

import "time"

type CreateNotificationRequest struct {
	UserID    uint
	ActorID   *uint
	Type      NotificationType
	TargetID  *uint
	TargetRef *string
	Message   string
}
type MarkReadRequest struct {
	IDs []uint `json:"ids"`
}
type DeleteNotificationsRequest struct {
	IDs []uint `json:"ids"`
}
type ActorSummary struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}
type ListNotificationsQuery struct {
	Unread   *bool  `query:"unread"`
	Category string `query:"category"` // NEW: system | social | media | achievement
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Type     string `query:"type"`
}
type NotificationResponse struct {
	ID        uint                 `json:"id"`
	Type      NotificationType     `json:"type"`
	Category  NotificationCategory `json:"category"` // NEW
	Message   string               `json:"message"`
	IsRead    bool                 `json:"is_read"`
	ReadAt    *time.Time           `json:"read_at,omitempty"`
	Actor     *ActorSummary        `json:"actor,omitempty"`
	TargetID  *uint                `json:"target_id,omitempty"`
	TargetRef *string              `json:"target_ref,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
}
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                  `json:"unread_count"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
}
type UnreadByCategoryResponse struct {
	Category NotificationCategory `json:"category"`
	Count    int64                `json:"count"`
}
