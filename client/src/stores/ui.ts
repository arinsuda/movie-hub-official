import { defineStore } from 'pinia'
import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'info' | 'warning'

export interface Toast {
  id: string
  type: ToastType
  message: string
  duration?: number
}

export const useUIStore = defineStore('ui', () => {
  const toasts = ref<Toast[]>([])
  const isSearchOpen = ref(false)

  function addToast(message: string, type: ToastType = 'info', duration = 3000) {
    const id = crypto.randomUUID()
    toasts.value.push({ id, type, message, duration })
    setTimeout(() => removeToast(id), duration)
  }

  function removeToast(id: string) {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  function toast(message: string) {
    addToast(message, 'info')
  }

  function toastSuccess(message: string) {
    addToast(message, 'success')
  }

  function toastError(message: string) {
    addToast(message, 'error')
  }

  function toastWarning(message: string) {
    addToast(message, 'warning')
  }

  function openSearch() {
    isSearchOpen.value = true
  }

  function closeSearch() {
    isSearchOpen.value = false
  }

  return {
    toasts,
    isSearchOpen,
    addToast,
    removeToast,
    toast,
    toastSuccess,
    toastError,
    toastWarning,
    openSearch,
    closeSearch,
  }
})
