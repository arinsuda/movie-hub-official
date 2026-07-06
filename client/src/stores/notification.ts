import { defineStore } from "pinia"
import { ref, computed } from "vue"
import { notiApi } from "@/api/endpoints/notification"
import { useSocket } from "@/composables/useSocket"
import type {
  AppNotification,
  NotificationPaginationMeta,
  SocketNotificationDeletedPayload,
  SocketNotificationNewPayload,
  SocketNotificationReadPayload,
} from "@/types/notification"

export const useNotificationStore = defineStore("notification", () => {
  const DEFAULT_PAGINATION: NotificationPaginationMeta = {
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0,
  }
  const notifications = ref<AppNotification[]>([])
  const unreadCount = ref(0)
  const loading = ref(false)
  const loadingMore = ref(false)
  const error = ref<string | null>(null)
  const isPanelOpen = ref(false)

  const pagination = ref<NotificationPaginationMeta>({ ...DEFAULT_PAGINATION })

  const hasMore = computed(
    () => pagination.value.page < pagination.value.total_pages,
  )

  let socketBound = false

  async function fetchNotifications(page = 1) {
    if (loading.value || loadingMore.value) return

    const limit = pagination.value?.limit ?? DEFAULT_PAGINATION.limit

    if (page === 1) loading.value = true
    else loadingMore.value = true

    error.value = null

    try {
      const res = await notiApi.list({ page, limit })

      const nextNotifications = res.data.notifications ?? []

      notifications.value =
        page === 1
          ? nextNotifications
          : [...notifications.value, ...nextNotifications]

      unreadCount.value = res.data.unread_count ?? 0

      pagination.value = {
        ...DEFAULT_PAGINATION,
        ...(res.data.pagination ?? {}),
        page,
        limit,
      }
    } catch (err) {
      console.error("fetchNotifications failed:", err)
      error.value = "โหลดการแจ้งเตือนไม่สำเร็จ"
    } finally {
      loading.value = false
      loadingMore.value = false
    }
  }

  async function loadMore() {
    if (!hasMore.value || loadingMore.value) return
    await fetchNotifications(pagination.value.page + 1)
  }

  async function fetchUnreadCount() {
    try {
      const res = await notiApi.getUnreadCount()
      unreadCount.value = res.data.unread_count
    } catch (err) {
      console.error("fetchUnreadCount failed:", err)
    }
  }

  async function refresh() {
    await fetchNotifications(1)
  }

  async function markAsRead(id: number) {
    const target = notifications.value.find(n => n.id === id)
    if (!target || target.is_read) return

    target.is_read = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)

    try {
      await notiApi.markAsRead(id)
    } catch (err) {
      target.is_read = false
      unreadCount.value++
      console.error("markAsRead failed:", err)
    }
  }

  async function markAllAsRead() {
    if (unreadCount.value === 0) return
    const prevState = notifications.value.map(n => n.is_read)
    const prevCount = unreadCount.value

    notifications.value.forEach(n => (n.is_read = true))
    unreadCount.value = 0

    try {
      await notiApi.markAllAsRead()
    } catch (err) {
      notifications.value.forEach((n, i) => (n.is_read = prevState[i] ?? false))
      unreadCount.value = prevCount
      console.error("markAllAsRead failed:", err)
    }
  }

  async function removeNotification(id: number) {
    const idx = notifications.value.findIndex(n => n.id === id)
    if (idx === -1) return

    const [removed] = notifications.value.splice(idx, 1)
    if (!removed) return

    if (!removed.is_read) {
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
    pagination.value.total = Math.max(0, pagination.value.total - 1)

    try {
      await notiApi.deleteNotification(id)
    } catch (err) {
      notifications.value.splice(idx, 0, removed)
      if (!removed.is_read) unreadCount.value++
      pagination.value.total++
      console.error("removeNotification failed:", err)
    }
  }

  async function clearAll() {
    const prev = notifications.value
    const prevUnread = unreadCount.value
    notifications.value = []
    unreadCount.value = 0

    try {
      await notiApi.clearAll()
    } catch (err) {
      notifications.value = prev
      unreadCount.value = prevUnread
      console.error("clearAll failed:", err)
    }
  }

  function prependNotification(n: AppNotification) {
    if (notifications.value.some(x => x.id === n.id)) return
    notifications.value.unshift(n)
    if (!n.is_read) unreadCount.value++
    pagination.value.total++
  }

  function handleRemoteRead({ ids }: SocketNotificationReadPayload) {
    ids.forEach(id => {
      const target = notifications.value.find(n => n.id === id)
      if (target && !target.is_read) {
        target.is_read = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    })
  }

  function handleRemoteDelete({ ids }: SocketNotificationDeletedPayload) {
    ids.forEach(id => {
      const idx = notifications.value.findIndex(n => n.id === id)
      if (idx === -1) return
      const [removed] = notifications.value.splice(idx, 1)
      if (removed && !removed.is_read) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    })
  }

  /**
   * เรียกตอน user login แล้ว (เช่นใน MainLayout onMounted / watch(authStore.user))
   * ผูก listener ของ event ที่ backend จะยิงเข้ามา:
   *  - "notification:new"     -> AppNotification ก้อนใหม่
   *  - "notification:read"    -> { id }  (เผื่อ mark read จากอุปกรณ์อื่น)
   *  - "notification:deleted" -> { id }
   */
  function bindSocket() {
    if (socketBound) return
    const { connect, on } = useSocket()
    connect()

    on<SocketNotificationNewPayload>("notification:new", payload => {
      prependNotification(payload)
      ;(window as any).$toast?.info?.(payload.message, payload.title)
    })

    on<SocketNotificationReadPayload>("notification:read", handleRemoteRead)

    on<SocketNotificationDeletedPayload>(
      "notification:deleted",
      handleRemoteDelete,
    )

    on<{ unread_count: number }>(
      "notification:unread_count",
      ({ unread_count }) => {
        unreadCount.value = unread_count
      },
    )
    socketBound = true
  }

  function unbindSocket() {
    const { disconnect } = useSocket()
    disconnect()
    socketBound = false
  }

  function togglePanel() {
    const nextOpen = !isPanelOpen.value
    isPanelOpen.value = nextOpen

    if (!nextOpen) return

    if (
      notifications.value.length === 0 &&
      !loading.value &&
      !loadingMore.value
    ) {
      fetchNotifications(1)
    }
  }

  function closePanel() {
    isPanelOpen.value = false
  }

  function reset() {
    notifications.value = []
    unreadCount.value = 0
    pagination.value = { ...DEFAULT_PAGINATION }
    isPanelOpen.value = false
    unbindSocket()
  }

  return {
    notifications,
    unreadCount,
    loading,
    loadingMore,
    error,
    pagination,
    hasMore,
    isPanelOpen,

    fetchNotifications,
    loadMore,
    fetchUnreadCount,
    refresh,
    markAsRead,
    markAllAsRead,
    removeNotification,
    clearAll,

    bindSocket,
    unbindSocket,
    togglePanel,
    closePanel,
    reset,
  }
})
