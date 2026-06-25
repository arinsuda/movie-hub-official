import type { UserProfile } from "./user"

export interface AuthUser extends UserProfile {
  email: string
  is_verified: boolean
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

export interface ForgotPasswordRequest {
  email: string
}

export interface ResetPasswordRequest {
  token: string
  user_id: number
  new_password: string
  confirm_password: string
}