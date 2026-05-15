import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true, // ส่ง cookie ไปด้วยทุก request (JWT อยู่ใน cookie)
  headers: {
    'Content-Type': 'application/json',
  },
})

// Response interceptor — จัดการ 401 auto refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // ถ้าได้ 401 และยังไม่เคย retry
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        // ลอง refresh token
        await axios.post('/api/auth/refresh', {}, { withCredentials: true })
        // retry request เดิม
        return api(originalRequest)
      } catch {
        // refresh ไม่สำเร็จ → logout
        const authStore = useAuthStore()
        authStore.clearUser()
        window.location.href = '/login'
      }
    }

    return Promise.reject(error)
  },
)

export default api
