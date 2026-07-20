import { defineStore } from "pinia"
import { ref, computed } from "vue"

import { authApi, userApi } from "@/api/api"

import type { AuthUser, LoginRequest, RegisterRequest } from "@/types"

export const useAuthStore = defineStore("auth", () => {
  const user = ref<AuthUser | null>(null)

  const isLoading = ref(false)

  const isInitialized = ref(false)

  const isLoggedIn = computed(() => !!user.value)

  const isAdmin = computed(() => user.value?.role === "admin")

  function setUser(u: AuthUser) {
    user.value = u
  }

  function clearUser() {
    user.value = null
    localStorage.removeItem("access_token")
    localStorage.removeItem("refresh_token")
  }

  async function login(data: LoginRequest) {
    isLoading.value = true

    try {
      const res = await authApi.login(data)
      if (res.data.access_token) {
        localStorage.setItem("access_token", res.data.access_token)
      }
      if (res.data.refresh_token) {
        localStorage.setItem("refresh_token", res.data.refresh_token)
      }

      setUser(res.data.user)

      return res.data.user
    } finally {
      isLoading.value = false
    }
  }

  const needsOnboarding = computed(
    () => isLoggedIn.value && !user.value?.favorite_genres,
  )

  async function register(data: RegisterRequest) {
    isLoading.value = true

    try {
      const res = await authApi.register(data)
      if (res.data.access_token) {
        localStorage.setItem("access_token", res.data.access_token)
      }
      if (res.data.refresh_token) {
        localStorage.setItem("refresh_token", res.data.refresh_token)
      }

      setUser(res.data.user)

      return res.data.user
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    try {
      await authApi.logout()
    } finally {
      clearUser()
    }
  }

  async function fetchMe() {
    try {
      const urlParams = new URLSearchParams(window.location.search)
      const urlAccess = urlParams.get("access_token")
      const urlRefresh = urlParams.get("refresh_token")

      if (urlAccess && urlRefresh) {
        localStorage.setItem("access_token", urlAccess)
        localStorage.setItem("refresh_token", urlRefresh)
        window.history.replaceState({}, document.title, window.location.pathname)
      }

      const refreshRes = await authApi.refreshToken()
      const userId = refreshRes.data.user.id

      if (refreshRes.data.access_token) {
        localStorage.setItem("access_token", refreshRes.data.access_token)
      }
      if (refreshRes.data.refresh_token) {
        localStorage.setItem("refresh_token", refreshRes.data.refresh_token)
      }

      const meRes = await userApi.me(userId)
      const profile = meRes.data.user

      setUser({
        ...profile, 
        email: refreshRes.data.user.email, 
        is_verified: refreshRes.data.user.is_verified, 
      })
    } catch {
      clearUser()
    } finally {
      isInitialized.value = true
    }
  }

  return {
    user,
    isLoading,
    isInitialized,
    isLoggedIn,
    isAdmin,
    needsOnboarding,
    setUser,
    clearUser,
    login,
    register,
    logout,
    fetchMe,
  }
})
