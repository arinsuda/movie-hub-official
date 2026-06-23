export type NotificationType =
  | "followed_you"
  | "movie_now_playing"
  | "following_reviewed"
  | "following_liked_review"
  | "following_added_watchlist"
  | "following_added_watched"

export interface MarkReadRequest {
  id: number
}

export interface ListNotifications {
  unread: boolean
  page: number
  page_size: number
  type: string
}

export interface ActorSummary {
  id: number
  username: string
  display_name: string
  avatar_url: string
}

export interface NotificationResponse {
  id: number
  type: NotificationType
  message: string
  is_read: boolean
  read_at: string
  actor: ActorSummary
  target_id: number
  target_ref: string
  created_at: string
}

export interface NotificationList {
  notifications: NotificationResponse[]
  unread_count: number
  total: number
  page: number
  page_size: number
}
