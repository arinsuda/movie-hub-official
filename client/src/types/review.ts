import type { UserSummary } from './auth'
import type { MediaSummary } from './movie'

export interface ReviewResponse {
  id: number
  user: UserSummary
  media: MediaSummary
  rating: number
  body: string
  is_public: boolean
  watched_at: string | null
  like_count: number
  comment_count: number
  is_liked: boolean
  created_at: string
  updated_at: string
}

export interface CreateReviewRequest {
  media_id: number
  media_type: 'movie' | 'tv'
  rating: number
  body: string
  is_public: boolean
  watched_at?: string
}

export interface UpdateReviewRequest {
  rating?: number
  body?: string
  is_public?: boolean
  watched_at?: string
}

export interface CommentResponse {
  id: number
  review_id: number
  user: UserSummary
  body: string
  created_at: string
  updated_at: string
}

export interface CreateCommentRequest {
  body: string
}

// Library
export type ListType = 'watchlist' | 'favorite' | 'watched'
export type MediaType = 'movie' | 'tv'

export interface LibraryItemResponse {
  id: number
  media: MediaSummary
  list_type: ListType
  watched_at: string | null
  tags: string[]
  note: string | null
  created_at: string
}

export interface AddItemRequest {
  media_id: number
  media_type: MediaType
  list_type: ListType
  watched_at?: string
  tags?: string[]
  note?: string
}

export interface MediaStatusResponse {
  media_id: number
  media_type: MediaType
  in_lists: ListType[]
}

// Pagination
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  limit: number
  total_pages: number
}
