import type { AppNotification, NotificationListResponse } from "@/types/notification"
import api from "../index"


export interface ListNotificationParams {
  page?: number
  limit?: number
  unread_only?: boolean
}

export const notiApi = {
  /**
   * GET /api/v3/notifications?page=&limit=&unread_only=
   */
  list: (params: ListNotificationParams = {}) =>
    api.get<NotificationListResponse>("/notifications", {
      params: { page: 1, limit: 20, ...params },
    }),

  /**
   * GET /api/v3/notifications/unread-count
   * เอาไว้ sync ตัวเลข badge ตอนเปิดแอปใหม่/socket หลุดแล้วต่อกลับ
   */
  getUnreadCount: () =>
    api.get<{ unread_count: number }>("/notifications/unread-count"),

  /**
   * PATCH /api/v3/notifications/:id/read
   */
  markAsRead: (notificationId: number) =>
    api.patch<{ notification: AppNotification }>(
      `/notifications/${notificationId}/read`,
    ),

  /**
   * PATCH /api/v3/notifications/read-all
   */
  markAllAsRead: () =>
    api.patch<{ updated: number }>("/notifications/read-all"),

  /**
   * DELETE /api/v3/notifications/:id
   */
  deleteNotification: (notificationId: number) =>
    api.delete(`/notifications/${notificationId}`),

  /**
   * DELETE /api/v3/notifications
   * ล้างการแจ้งเตือนทั้งหมดของ user
   */
  clearAll: () => api.delete("/notifications"),
}
