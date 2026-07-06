import type { AppNotification, NotificationListResponse } from "@/types/notification"
import api from "../index"


export interface ListNotificationParams {
  page?: number
  limit?: number
  unread_only?: boolean
}

export const notiApi = {
  list: (params: ListNotificationParams = {}) =>
    api.get<NotificationListResponse>("/notifications", {
      params: { page: 1, limit: 20, ...params },
    }),

  getUnreadCount: () =>
    api.get<{ unread_count: number }>("/notifications/unread-count"),

  markAsRead: (notificationId: number) =>
    api.patch<{ notification: AppNotification }>(
      `/notifications/${notificationId}/read`,
    ),

  markAllAsRead: () =>
    api.patch<{ updated: number }>("/notifications/read-all"),

  deleteNotification: (notificationId: number) =>
    api.delete(`/notifications/${notificationId}`),

  clearAll: () => api.delete("/notifications"),
}
