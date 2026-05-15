export interface UserSummary {
  id: number
  username: string
  avatar_url: string
}

export interface UserProfile {
  id: number
  username: string
  display_name: string | null
  bio: string | null
  avatar_url: string | null
  gender: string
  favorite_genres: string | null
  review_count: number
  follower_count: number
  following_count: number
  is_private: boolean
  role: string
}

export interface AuthUser {
  id: number
  username: string
  email: string
  display_name: string | null
  avatar_url: string | null
  is_verified: boolean
  role: string
}

export interface LoginRequest {
  identifier: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  confirm_password: string
}

export interface AuthResponse {
  user: AuthUser
}
