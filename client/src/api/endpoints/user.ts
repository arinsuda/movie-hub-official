import api from "../index"
import type { ChangePassword, UserProfile } from "@/types"

export const userApi = {
  me: (userId: number) => api.get<{ user: UserProfile }>(`/users/${userId}`),

  getProfile: (userId: number) =>
    api.get<{ user: UserProfile }>(`/users/${userId}`),

  updateProfile: (userId: number, data: FormData) =>
    api.patch<{ user: UserProfile }>(`/users/${userId}`, data, {
      headers: { "Content-Type": "multipart/form-data" },
    }),

  deleteUser: (userId: number) => api.delete(`/users/${userId}`),

  updateFavoriteGenres: (userId: number, genres: number[]) =>
    api.patch<{ user: UserProfile }>(`/users/${userId}/genres`, {
      favorite_genres: genres,
    }),

  requestEmailChange: (userId: number) => api.post(`/users/${userId}/email`),

  verifyEmailChange: (userId: number, otp: string) =>
    api.put<{ user: UserProfile }>(`/users/${userId}/email`, { otp }),

  updateEmail: (userId: number, newEmail: string) =>
    api.patch(`/users/${userId}/email`, { new_email: newEmail }),

  // รับ body เพื่อส่ง old/new/confirm password ไปด้วย
  changePassword: (userId: number, data: ChangePassword) =>
    api.patch(`/users/${userId}/password`, data),
}
