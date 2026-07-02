import { reviewApi } from "@/api/api"
import { useAuthStore } from "@/stores/auth"
import { useToast } from "@/composables/useToast"
import { computed, ref } from "vue"
import type { ReviewResponse } from "@/types/review"

type MediaType = "movie" | "tv"

function toDateString(d: Date): string {
  // BE parses with time.Parse("2006-01-02", ...), so this must stay zero-padded
  // local-time YYYY-MM-DD (not toISOString, which shifts to UTC).
  const yyyy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, "0")
  const dd = String(d.getDate()).padStart(2, "0")
  return `${yyyy}-${mm}-${dd}`
}

export function useReviewForm(mediaId: number, mediaType: MediaType) {
  const auth = useAuthStore()
  const { error: toastError } = useToast()

  const body = ref("")
  const rating = ref(5)
  const isPublic = ref(true)
  const watchedAt = ref<Date | null>(null)
  const isSubmitting = ref(false)

  const isAuthenticated = computed(() => Boolean(auth.user?.id))
  const canSubmit = computed(
    () =>
      isAuthenticated.value &&
      body.value.trim().length > 0 &&
      !isSubmitting.value,
  )

  function reset() {
    body.value = ""
    rating.value = 5
    isPublic.value = true
    watchedAt.value = null
  }

  async function submit(): Promise<ReviewResponse | null> {
    if (!isAuthenticated.value) {
      toastError("กรุณาเข้าสู่ระบบก่อนเขียนรีวิว")
      return null
    }
    if (!body.value.trim() || isSubmitting.value) return null

    isSubmitting.value = true
    try {
      const res = await reviewApi.createReview(auth.user!.id, {
        media_id: mediaId,
        media_type: mediaType,
        rating: rating.value,
        body: body.value,
        is_public: isPublic.value,
        watched_at: watchedAt.value ? toDateString(watchedAt.value) : null,
      })
      reset()
      return res.data.review
    } catch (err) {
      console.error("submitReview failed:", err)
      toastError("ส่งรีวิวไม่สำเร็จ กรุณาลองใหม่")
      return null
    } finally {
      isSubmitting.value = false
    }
  }

  return {
    body,
    rating,
    isPublic,
    watchedAt,
    isSubmitting,
    isAuthenticated,
    canSubmit,
    submit,
  }
}
