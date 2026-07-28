import { reviewApi } from "@/api/api"
import { useAuthStore } from "@/stores/auth"
import { useToast } from "@/composables/useToast"
import { computed, ref } from "vue"
import type { ReviewResponse } from "@/types/review"
import type { Visibility } from "@/types"

type MediaType = "movie" | "tv"



export function useReviewForm(mediaId: number, mediaType: MediaType) {
  const auth = useAuthStore()
  const { error: toastError } = useToast()

  const body = ref("")
  const rating = ref(5)
  const visibility = ref<Visibility>("public")
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
    visibility.value = "public"
  }

  async function submit(): Promise<ReviewResponse | null> {
    if (!isAuthenticated.value) {
      toastError("กรุณาเข้าสู่ระบบก่อนเขียนรีวิว")
      return null
    }
    if (!body.value.trim() || isSubmitting.value) return null

    isSubmitting.value = true
    try {
      const res = await reviewApi.createReview({
        media_id: mediaId,
        media_type: mediaType,
        rating: rating.value,
        body: body.value,
        visibility: visibility.value,
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
    visibility,
    isSubmitting,
    isAuthenticated,
    canSubmit,
    submit,
  }
}
