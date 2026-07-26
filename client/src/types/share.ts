import type { MediaType } from "./common"
import type { Genre } from "./movie"

export interface ShareMediaContext {
  id: number
  mediaType: MediaType
  title: string
  posterPath: string | null
  posterUrl: string | null
  releaseYear: string | null
  genres: Genre[]
  voteAverage: number
  removRating: number | null
}

export interface ShareReviewContext {
  id: number
  authorDisplayName: string
  authorUsername: string
  authorAvatarUrl: string | null
  rating: number
  body: string
  isPublic: boolean
  watchedAt: string | null
  createdAt: string
}
