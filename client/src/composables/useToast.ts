import { useUIStore } from '@/stores/ui'

export function useToast() {
  const ui = useUIStore()
  return {
    toast: ui.toast,
    success: ui.toastSuccess,
    error: ui.toastError,
    warning: ui.toastWarning,
  }
}
