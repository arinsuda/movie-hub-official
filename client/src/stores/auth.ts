import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/api'
import type { AuthUser, LoginRequest, RegisterRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUser | null>(null)
  const isLoading = ref(false)

  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setUser(u: AuthUser) {
    user.value = u
  }

  function clearUser() {
    user.value = null
  }

  async function login(data: LoginRequest) {
    isLoading.value = true
    try {
      const res = await authApi.login(data)
      setUser(res.data.user)
      return res.data.user
    } finally {
      isLoading.value = false
    }
  }

  async function register(data: RegisterRequest) {
    isLoading.value = true
    try {
      const res = await authApi.register(data)
      setUser(res.data.user)
      return res.data.user
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    await authApi.logout()
    clearUser()
  }

  async function fetchMe() {
    // เรียก refresh เพื่อเช็คว่า cookie ยังใช้ได้ไหม
    try {
      const res = await authApi.refresh()
      setUser(res.data.user)
    } catch {
      clearUser()
    }
  }

  return {
    user,
    isLoading,
    isLoggedIn,
    isAdmin,
    setUser,
    clearUser,
    login,
    register,
    logout,
    fetchMe,
  }
})
