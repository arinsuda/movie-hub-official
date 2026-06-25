import api from "../index"

import type {
  AuthResponse,
  ForgotPasswordRequest,
  LoginRequest,
  RegisterRequest,
  ResetPasswordRequest,
  UserProfile,
} from "@/types"

export const authApi = {
  login: (data: LoginRequest) => api.post<AuthResponse>("/auth/login", data),

  register: (data: RegisterRequest) =>
    api.post<AuthResponse>("/auth/register", data),

  refreshToken: () => api.post<AuthResponse>("/auth/refresh"),

  logout: () => api.post("/auth/logout"),

  verifyEmail: (token: string) => api.get(`/auth/verify-email?token=${token}`),

  resendVerification: (email: string) =>
    api.post("/auth/resend-verification", { email }),

  logoutAll: () => api.post("auth/logout-all"),

  // ส่ง email เพื่อรับ reset link
  forgotPassword: (data: ForgotPasswordRequest) =>
    api.post("auth/forgot-password", data),

  // ใช้ token จาก email ตั้งรหัสผ่านใหม่
  resetPassword: (data: ResetPasswordRequest) =>
    api.post("auth/reset-password", data),
}
