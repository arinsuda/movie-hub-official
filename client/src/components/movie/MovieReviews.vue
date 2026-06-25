<template>
  <div class="movie-reviews-component">
    <div class="write-review-box">
      <textarea
        v-model="newReview"
        placeholder="แบ่งปันความคิดเห็นของคุณต่อภาพยนตร์เรื่องนี้..."
        rows="3"
      ></textarea>
      <div class="form-footer">
        <div class="rating-input">
          <span>ให้คะแนนของคุณ: </span>
          <div class="stars-wrap">
            <div
              v-for="n in 5"
              :key="n"
              class="star-slot"
              @mouseleave="hoverRating = 0"
            >
              <div
                class="half-zone left"
                @click="userRating = n - 0.5"
                @mouseenter="hoverRating = n - 0.5"
              />
              <div
                class="half-zone right"
                @click="userRating = n"
                @mouseenter="hoverRating = n"
              />
              <svg
                class="star-base"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
                  stroke="currentColor"
                  stroke-width="1.5"
                  stroke-linejoin="round"
                  fill="none"
                />
              </svg>
              <svg
                class="star-fill"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <defs>
                  <clipPath :id="`clip-movie-${n}-half`">
                    <rect x="0" y="0" width="12" height="24" />
                  </clipPath>
                  <clipPath :id="`clip-movie-${n}-full`">
                    <rect x="0" y="0" width="24" height="24" />
                  </clipPath>
                </defs>
                <path
                  d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
                  fill="currentColor"
                  :clip-path="starClip(n, displayRating)"
                  :class="n - 0.5 <= displayRating ? 'fill--on' : 'fill--off'"
                />
              </svg>
            </div>
          </div>
          <span class="rating-label">{{ displayRating }} / 5</span>
        </div>
        <button class="btn-submit-review" @click="submitReview">
          ส่งรีวิว
        </button>
      </div>
    </div>

    <div class="reviews-list" v-if="reviews.length">
      <div v-for="item in reviews" :key="item.id" class="review-card">
        <div class="review-header">
          <div class="user-info">
            <div class="user-avatar">
              <img
                v-if="item.user.avatar_url"
                :src="item.user.avatar_url"
                :alt="item.user.username"
              />
              <i v-else class="pi pi-user"></i>
            </div>
            <div class="user-meta">
              <span class="username">
                {{ item.user.display_name || item.user.username }}
              </span>
              <span class="review-date">{{ formatDate(item.created_at) }}</span>
            </div>
          </div>
          <div class="user-score">
            <i class="pi pi-star-fill"></i>
            <span>{{ item.rating }}/5</span>
          </div>
        </div>
        <p class="review-comment">{{ item.body }}</p>
      </div>
    </div>

    <div v-else class="empty-reviews">
      <i class="pi pi-inbox"></i>
      <p>ยังไม่มีรีวิวสำหรับเรื่องนี้ มารีวิวเป็นคนแรกกันเถอะ!</p>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { reviewApi } from "@/api/api"
  import { useAuthStore } from "@/stores/auth"
  import { computed, onMounted, ref } from "vue"
  import type { ReviewResponse } from "@/types/review"

  const props = defineProps<{
    movieId: number
    mediaType?: "movie" | "tv"
  }>()

  const auth = useAuthStore()
  const newReview = ref("")
  const userRating = ref(5)
  const hoverRating = ref<number>(0)
  const reviews = ref<ReviewResponse[]>([])
  const userId = auth.user?.id || 0
  const displayRating = computed(() => hoverRating.value || userRating.value)

  // ─── Derived values ───────────────────────────────────────────────────────
  // media_type ส่งไป API: "movie" | "tv"
  const mt = computed(() => props.mediaType ?? "movie")
  // path segment สำหรับ getMediaReviews: "movies" | "tv"
  const reviewSegment = computed(() => (mt.value === "tv" ? "tv" : "movies"))

  function formatDate(iso: string) {
    return new Date(iso).toLocaleDateString("th-TH", {
      year: "numeric",
      month: "short",
      day: "numeric",
    })
  }

  async function submitReview() {
    if (!newReview.value.trim()) return

    try {
      const res = await reviewApi.createReview(userId, {
        media_id: props.movieId,
        media_type: mt.value,
        rating: userRating.value,
        body: newReview.value,
        is_public: true,
      })
      reviews.value.unshift(res.data.review)
      newReview.value = ""
      userRating.value = 5
    } catch (err) {
      console.error("submitReview failed:", err)
    }
  }

  function starClip(n: number, rating: number) {
    if (rating >= n) return `url(#clip-movie-${n}-full)`
    if (rating >= n - 0.5) return `url(#clip-movie-${n}-half)`
    return "none"
  }

  onMounted(async () => {
    if (props.movieId !== 0) {
      try {
        const res = await reviewApi.getMediaReviews(
          reviewSegment.value,
          props.movieId,
        )
        reviews.value = res.data.reviews
      } catch (err) {
        console.error("getMediaReviews failed:", err)
      }
    }
  })
</script>

<style scoped>
  .movie-reviews-component {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    margin-top: 0.5rem;
  }

  .write-review-box {
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 12px;
    padding: 1rem;
  }

  textarea {
    width: 100%;
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 0.75rem;
    color: #fff;
    font-family: inherit;
    font-size: 0.9rem;
    resize: none;
    outline: none;
  }
  textarea:focus {
    border-color: #e50914;
  }

  .form-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 0.75rem;
  }

  .rating-input {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.85rem;
    color: #8a8a8a;
  }
  .rating-input i {
    cursor: pointer;
    font-size: 1rem;
    transition: color 0.1s;
  }
  .rating-input i.active {
    color: #f5c518;
  }

  .btn-submit-review {
    background: #fff;
    color: #000;
    border: none;
    padding: 0.45rem 1.25rem;
    border-radius: 6px;
    font-weight: 600;
    font-size: 0.85rem;
    cursor: pointer;
    transition: opacity 0.2s;
  }
  .btn-submit-review:hover {
    opacity: 0.9;
  }

  .reviews-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .review-card {
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 10px;
    padding: 1rem;
  }

  .review-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 0.75rem;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .user-avatar {
    background: rgba(255, 255, 255, 0.1);
    width: 36px;
    height: 36px;
    border-radius: 9999px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
  }

  .user-meta {
    display: flex;
    flex-direction: column;
  }

  .username {
    font-size: 0.875rem;
    font-weight: 600;
    color: #fff;
  }

  .review-date {
    font-size: 0.75rem;
    color: #8a8a8a;
  }

  .user-score {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    background: rgba(245, 197, 24, 0.1);
    color: #f5c518;
    padding: 0.25rem 0.6rem;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 600;
  }

  .review-comment {
    font-size: 0.9rem;
    line-height: 1.6;
    color: #ccc;
    margin: 0;
  }

  .empty-reviews {
    text-align: center;
    padding: 2.5rem;
    color: #8a8a8a;
    background: rgba(255, 255, 255, 0.01);
    border-radius: 12px;
    border: 1px dashed rgba(255, 255, 255, 0.05);
  }
  .empty-reviews i {
    font-size: 2rem;
    margin-bottom: 0.5rem;
  }
  .empty-reviews p {
    margin: 0;
    font-size: 0.875rem;
  }

  .stars-wrap {
    display: flex;
    gap: 3px;
  }
  .star-slot {
    position: relative;
    width: 22px;
    height: 22px;
    cursor: pointer;
    transition: transform 0.12s ease;
    flex-shrink: 0;
  }
  .star-slot:hover {
    transform: scale(1.18);
  }
  .star-base,
  .star-fill {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
  }
  .star-base {
    color: rgba(255, 255, 255, 0.2);
  }
  .fill--on {
    color: #f5c518;
  }
  .fill--off {
    color: transparent;
  }
  .half-zone {
    position: absolute;
    top: 0;
    height: 100%;
    width: 50%;
    cursor: pointer;
    z-index: 1;
  }
  .half-zone.left {
    left: 0;
  }
  .half-zone.right {
    right: 0;
  }
  .rating-label {
    font-size: 0.85rem;
    color: #8a8a8a;
    min-width: 40px;
  }
</style>
