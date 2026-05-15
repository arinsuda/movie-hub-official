import type { UserSummary } from './auth'
import type { MediaSummary } from './movie'

export interface FollowStatsResponse {
  user_id: number
  followers: number
  following: number
  is_following: boolean
}

export interface FollowUserResponse {
  user: UserSummary
  is_following: boolean
  followed_at: string
}

export interface FeedReviewResponse {
  id: number
  user: UserSummary
  media: MediaSummary
  rating: number
  body: string
  like_count: number
  comment_count: number
  is_liked: boolean
  created_at: string
}

export interface FeedItemResponse {
  review: FeedReviewResponse
  created_at: string
}
