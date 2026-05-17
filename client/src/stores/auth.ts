import { defineStore } from "pinia";
import { ref, computed } from "vue";

import { authApi } from "@/api/api";

import type { AuthUser, LoginRequest, RegisterRequest } from "@/types";

export const useAuthStore = defineStore("auth", () => {
  const user = ref<AuthUser | null>(null);

  const isLoading = ref(false);

  const isInitialized = ref(false);

  const isLoggedIn = computed(() => !!user.value);

  const isAdmin = computed(() => user.value?.role === "admin");

  function setUser(u: AuthUser) {
    user.value = u;
  }

  function clearUser() {
    user.value = null;
  }

  async function login(data: LoginRequest) {
    isLoading.value = true;

    try {
      const res = await authApi.login(data);

      setUser(res.data.user);

      return res.data.user;
    } finally {
      isLoading.value = false;
    }
  }

  async function register(data: RegisterRequest) {
    isLoading.value = true;

    try {
      const res = await authApi.register(data);

      setUser(res.data.user);

      return res.data.user;
    } finally {
      isLoading.value = false;
    }
  }

  async function logout() {
    try {
      await authApi.logout();
    } finally {
      clearUser();
    }
  }

  async function fetchMe() {
    try {
      // ใช้ me endpoint
      // interceptor จะ refresh ให้อัตโนมัติถ้า access token หมด
      const res = await authApi.refresh(); 

      setUser(res.data.user);
    } catch {
      clearUser();
    } finally {
      isInitialized.value = true;
    }
  }

  return {
    user,
    isLoading,
    isInitialized,
    isLoggedIn,
    isAdmin,
    setUser,
    clearUser,
    login,
    register,
    logout,
    fetchMe,
  };
});
