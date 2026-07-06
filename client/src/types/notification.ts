

export type NotificationType =
  | "follow"
  | "review_like"
  | "review_comment"
  | "comment_reply"
  | "achievement_unlock"
  | "system"

export interface NotificationActor {
  id: number
  username: string
  avatar_url?: string | null
}

export interface AppNotification {
  id: number
  user_id: number
  type: NotificationType
  title: string
  message: string
  is_read: boolean
  action_url?: string | null
  actor?: NotificationActor | null
  metadata?: Record<string, unknown> | null
  created_at: string
}

export interface NotificationPaginationMeta {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface NotificationListResponse {
  notifications: AppNotification[]
  unread_count: number
  pagination: NotificationPaginationMeta
}

// ── Socket.IO event payloads ─────────────────────────────────────
// (สัญญากับฝั่ง Backend — ปรับชื่อ event ให้ตรงกับที่ BE ยิงจริง)

export interface SocketNotificationNewPayload extends AppNotification {}

export interface SocketNotificationReadPayload {
  id: number
}

export interface SocketNotificationDeletedPayload {
  id: number
}
