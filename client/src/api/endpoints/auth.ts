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

  logout: () => api.post("/auth/logout"),

  refreshToken: () => api.post<AuthResponse>("/auth/refresh"),

  me: (userId: number) => api.get<{ user: UserProfile }>(`/users/${userId}`),

  verifyEmail: (token: string) => api.get(`/auth/verify-email?token=${token}`),

  resendVerification: (email: string) =>
    api.post("/auth/resend-verification", {
      email,
    }),
}
