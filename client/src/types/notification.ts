export type NotificationType =
  | "followed_you"
  | "following_reviewed"
  | "following_liked_review"
  | "following_added_watchlist"
  | "following_added_watched"
  | "following_marked_helpful"
  | "following_commented"
  | "review_liked"
  | "review_commented"
  | "review_marked_helpful"
  | "achievement_unlocked"
  | "movie_now_playing"
  | "email_verified"
  | "password_changed"
  | "welcome"

export type NotificationCategory = "system" | "social" | "media" | "achievement"

export interface NotificationActor {
  id: number
  username: string
  display_name?: string | null
  avatar_url?: string | null
}

export interface AppNotification {
  id: number
  type: NotificationType
  category: NotificationCategory
  message: string
  is_read: boolean
  read_at?: string | null
  actor?: NotificationActor | null
  target_id?: number | null
  target_ref?: string | null
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

export interface SocketNotificationNewPayload extends AppNotification {}
export interface SocketNotificationReadPayload {
  ids: number[]
}
export interface SocketNotificationDeletedPayload {
  ids: number[] | null
}

const TITLE_BY_TYPE: Record<NotificationType, string> = {
  followed_you: "ผู้ติดตามใหม่",
  following_reviewed: "กิจกรรมจากคนที่คุณติดตาม",
  following_liked_review: "กิจกรรมจากคนที่คุณติดตาม",
  following_added_watchlist: "กิจกรรมจากคนที่คุณติดตาม",
  following_added_watched: "กิจกรรมจากคนที่คุณติดตาม",
  following_marked_helpful: "กิจกรรมจากคนที่คุณติดตาม",
  following_commented: "กิจกรรมจากคนที่คุณติดตาม",
  review_liked: "มีคนถูกใจรีวิวของคุณ",
  review_commented: "มีคอมเมนต์ใหม่ในรีวิวของคุณ",
  review_marked_helpful: "รีวิวของคุณถูกโหวตว่ามีประโยชน์",
  achievement_unlocked: "ปลดล็อกความสำเร็จ!",
  movie_now_playing: "หนังในวอทช์ลิสต์เข้าฉายแล้ว",
  email_verified: "ยืนยันอีเมลสำเร็จ",
  password_changed: "เปลี่ยนรหัสผ่านสำเร็จ",
  welcome: "ยินดีต้อนรับ",
}

export function getNotificationTitle(n: AppNotification): string {
  return TITLE_BY_TYPE[n.type] ?? "การแจ้งเตือน"
}

export function getNotificationActionUrl(n: AppNotification): string | null {
  if (!n.target_ref || n.target_id == null) return null
  switch (n.target_ref) {
    case "review":
      return `/reviews/${n.target_id}`
    case "user":
      return `/profile/${n.target_id}`
    case "movie":
      return `/movies/${n.target_id}`
    default:
      return null
  }
}

export const CATEGORY_LABELS: Record<NotificationCategory, string> = {
  social: "ติดตามและโซเชียล",
  media: "หนังและซีรีส์",
  achievement: "ความสำเร็จ",
  system: "ระบบ",
}

export interface NotificationCategoryTab {
  label: string
  value: NotificationCategory | null
}

export const NOTIFICATION_CATEGORY_TABS: NotificationCategoryTab[] = [
  { label: "ทั้งหมด", value: null },
  { label: CATEGORY_LABELS.social, value: "social" },
  { label: CATEGORY_LABELS.media, value: "media" },
  { label: CATEGORY_LABELS.achievement, value: "achievement" },
  { label: CATEGORY_LABELS.system, value: "system" },
]
