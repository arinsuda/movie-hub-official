import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useToast } from './useToast'

export function useAuth() {
  const authStore = useAuthStore()
  const { user, isLoggedIn, isLoading } = storeToRefs(authStore)
  const router = useRouter()
  const toast = useToast()

  async function logout() {
    await authStore.logout()
    toast.success('ออกจากระบบแล้ว')
    router.push('/login')
  }

  return {
    user,
    isLoggedIn,
    isLoading,
    logout,
    login: authStore.login,
    register: authStore.register,
  }
}
