import api from "../index"

import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
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
    api.post("/auth/resend-verification", {
      email,
    }),

  logoutAll: () => api.post("auth/logout-all"),
}
