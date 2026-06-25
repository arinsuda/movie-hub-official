export interface UserProfile {
  id: number
  username: string
  email: string 
  verified_email_at: string | null
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

export interface UpdateProfileRequest {
  display_name?: string
  bio?: string
  gender?: "male" | "female" | "other"
  gender_other?: string
  favorite_genres?: string
  date_of_birth?: string
  is_private?: boolean
}

export interface UserSummary {
  id: number
  username: string
  display_name: string | null
  avatar_url: string | null
}

export interface RequestEmailChangeRequest {
  new_email: string
}

export interface VerifyEmailChangeRequest {
  otp: string
}

export interface ChangePassword {
  old_password: string
  new_password: string
  confirm_password: string
}