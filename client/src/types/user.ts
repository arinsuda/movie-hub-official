export interface UserProfile {
  id: number
  username: string
  display_name: string | null
  bio: string | null
  avatar_url: string | null
  date_of_birth: string
  gender: string
  favorite_genres: string | null
  review_count: number
  follower_count: number
  following_count: number
  is_private: boolean
  level: number
  role: string
  created_at: string
}

export interface UserSummary {
  id: number
  username: string
  display_name: string | null
  avatar_url: string | null
}
