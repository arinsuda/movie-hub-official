// client/src/types/movie.ts
export interface Genre {
  id: number
  name: string
}

export interface Movie {
  id: number
  title: string
  original_title: string
  overview: string
  poster_path: string
  backdrop_path: string
  release_date: string
  vote_average: number
  vote_count: number
  popularity: number
  genre_ids: number[]
  adult: boolean
}

export interface MovieDetail extends Movie {
  runtime: number
  status: string
  tagline: string
  budget: number
  revenue: number
  homepage: string
  genres: Genre[]
  videos?: Video[] // ← populated after /movie/{id}/videos fetch
}

export interface TVSeries {
  id: number
  name: string
  original_name: string
  overview: string
  poster_path: string
  backdrop_path: string
  first_air_date: string
  vote_average: number
  vote_count: number
  popularity: number
  genre_ids: number[]
  adult: boolean
}

export interface TVSeriesDetail extends TVSeries {
  number_of_seasons: number
  number_of_episodes: number
  status: string
  tagline: string
  homepage: string
  genres: Genre[]
}

export interface Cast {
  id: number
  name: string
  character: string
  profile_path: string
  order: number
}

export interface Crew {
  id: number
  name: string
  job: string
  department: string
  profile_path: string
}

export interface Credits {
  cast: Cast[]
  crew: Crew[]
}

export interface Video {
  id: string
  key: string
  name: string
  site: string
  type: string
  official: boolean
}

export interface PaginatedResult<T> {
  page: number
  results: T[]
  total_pages: number
  total_results: number
}

// enriched media จาก backend (shared.MediaSummary)
export interface MediaSummary {
  id: number
  title: string
  poster_url: string
  backdrop_url: string
  release_year: string
  vote_average: number
  media_type: "movie" | "tv"
}

export interface MediaStats {
  media_id: number
  media_type: string
  like_count: number
  view_count: number
  review_count: number
  watchlist_count: number
  average_rating: number
  has_rating: boolean
  liked_at: string | null
  watchlisted_at: string | null
}

export interface MediaStatsResponse {
  stats: MediaStats
}

// ── เพิ่มเติมต่อท้ายไฟล์ client/src/types/movie.ts ──────────────────

export type ListType = "watchlist" | "favorite" | "watched"

export interface AddItemRequest {
  media_id: number
  media_type: "movie" | "tv"
  list_type: ListType
  watched_at?: string // รูปแบบ "YYYY-MM-DD" (ใช้เฉพาะกรณี list_type: 'watched')
  note?: string
  tags?: string[]
}

export interface LibraryItemResponse {
  id: number
  media_id: number
  media_type: "movie" | "tv"
  list_type: ListType
  watched_at: string | null
  tags: string[]
  note: string | null
  created_at: string
}

// 💡 โครงสร้างใหม่ที่ส่งไอดีแถวกลับมาให้หน้าบ้านใช้ลบข้อมูลได้
export interface MediaItemStatus {
  list_type: ListType
  item_id: number
}

export interface MediaStatusResponse {
  media_id: number
  media_type: "movie" | "tv"
  in_lists: MediaItemStatus[]
}
