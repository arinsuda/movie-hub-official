// client/src/api/endpoints/library.ts
import api from "../index"
import type {
  AddItemRequest,
  LibraryItemResponse,
  MediaStatusResponse,
} from "@/types/movie"

export const libraryApi = {
  /**
   * ตรวจสอบว่าหนัง/ซีรีส์เรื่องนี้อยู่ในลิสต์ไหนของยูสเซอร์บ้าง (ได้ item_id กลับมาด้วย)
   * GET /api/v3/users/:userId/library/status?media_id=xxx&media_type=xxx
   */
  getMediaStatus: (
    userId: number,
    mediaId: number,
    mediaType: "movie" | "tv",
  ) =>
    api.get<MediaStatusResponse>(`/users/${userId}/library/status`, {
      params: { media_id: mediaId, media_type: mediaType },
      withCredentials: true,
    }),

  /**
   * เพิ่มไอเทมเข้าคลังภาพยนตร์ (Watchlist / Watched / Favorite)
   * POST /api/v3/users/:userId/library
   */
  addItem: (userId: number, data: AddItemRequest) =>
    api.post<{ item: LibraryItemResponse }>(`/users/${userId}/library`, data),

  /**
   * ลบไอเทมออกจากคลังภาพยนตร์ (Unwatchlist / Unwatched)
   * DELETE /api/v3/users/:userId/library/:itemId
   */
  removeItem: (userId: number, itemId: number) =>
    api.delete(`/users/${userId}/library/${itemId}`),

  /**
   * ดึงรายการหนังทั้งหมดในคลังแยกตามประเภทลิสต์
   * GET /api/v3/users/:userId/library?list_type=watchlist
   */
  getLibrary: (
    userId: number,
    params?: { list_type?: string; media_type?: string },
  ) =>
    api.get<{ items: LibraryItemResponse[] }>(`/users/${userId}/library`, {
      params,
    }),
}
