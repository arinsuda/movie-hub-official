package library_module

import "time"

type AddItemRequest struct {
	MediaID   int       `json:"media_id"`
	MediaType MediaType `json:"media_type"` // "movie" | "tv"
	ListType  ListType  `json:"list_type"`  // "watchlist" | "favorite" | "watched"
	WatchedAt *string   `json:"watched_at"` // "2026-05-01" ใช้เฉพาะ watched
	Tags      []string  `json:"tags"`
	Note      *string   `json:"note"`
}

type UpdateItemRequest struct {
	WatchedAt *string  `json:"watched_at"`
	Tags      []string `json:"tags"`
	Note      *string  `json:"note"`
}

type LibraryItemResponse struct {
	ID        uint       `json:"id"`
	MediaID   int        `json:"media_id"`
	MediaType MediaType  `json:"media_type"`
	ListType  ListType   `json:"list_type"`
	WatchedAt *time.Time `json:"watched_at"`
	Tags      []string   `json:"tags"`
	Note      *string    `json:"note"`
	CreatedAt time.Time  `json:"created_at"`
}

// สำหรับ query สถานะ media นึงว่าอยู่ใน list ไหนบ้าง
type MediaStatusResponse struct {
	MediaID   int        `json:"media_id"`
	MediaType MediaType  `json:"media_type"`
	InLists   []ListType `json:"in_lists"` // ["watched", "favorite"]
}
