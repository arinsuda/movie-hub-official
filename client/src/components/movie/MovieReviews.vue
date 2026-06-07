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
            <div v-for="star in 5" :key="star" class="star-slot">
              <!-- ซีกซ้าย = ครึ่งดาว -->
              <div
                class="half-zone left"
                @click="userRating = star - 0.5"
                @mouseenter="hoverRating = star - 0.5"
                @mouseleave="hoverRating = 0"
              />
              <!-- ซีกขวา = เต็มดาว -->
              <div
                class="half-zone right"
                @click="userRating = star"
                @mouseenter="hoverRating = star"
                @mouseleave="hoverRating = 0"
              />

              <!-- ดาวเต็ม -->
              <i v-if="displayRating >= star" class="pi pi-star-fill active" />
              <!-- ครึ่งดาว — ใช้ clip-path แบ่งครึ่ง -->
              <span v-else-if="displayRating >= star - 0.5" class="half-star">
                <i class="pi pi-star-fill active half" />
                <i class="pi pi-star" />
              </span>
              <!-- ดาวว่าง -->
              <i v-else class="pi pi-star" />
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
              <!-- ถ้ามี avatar ให้แสดงรูป ถ้าไม่มีแสดง icon -->
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
            <span>{{ item.rating }}/10</span>
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
  import { useAuthStore } from "@/stores/auth" // ปรับ path ตาม project
  import { computed, onMounted, ref } from "vue"
  import type { ReviewResponse } from "@/types/review"

  const props = defineProps<{
    movieId: number
  }>()

  const auth = useAuthStore()
  const newReview = ref("")
  const userRating = ref(5)
  const hoverRating = ref<number>(0)
  const reviews = ref<ReviewResponse[]>([])
  const userId = auth.user?.id || 0
  const displayRating = computed(() => hoverRating.value || userRating.value)

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
        media_type: "movie",
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

  onMounted(async () => {
    if (props.movieId !== 0) {
      try {
        const res = await reviewApi.getMediaReviews("movies", props.movieId)
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
    gap: 4px;
  }

  .star-slot {
    position: relative;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* ซ่อน i ไว้ข้างหลัง zone */
  .star-slot .pi {
    font-size: 1.2rem;
    pointer-events: none;
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

  /* ครึ่งดาว — ซ้อน 2 icon แล้ว clip ซีกซ้าย */
  .half-star {
    position: relative;
    display: inline-block;
    width: 1.2rem;
    height: 1.2rem;
  }
  .half-star .pi {
    position: absolute;
    left: 0;
    top: 0;
  }
  .half-star .pi.half {
    clip-path: inset(0 50% 0 0); /* แสดงแค่ซีกซ้าย */
  }

  .pi.active {
    color: #f5c518;
  }

  .rating-label {
    font-size: 0.85rem;
    color: #8a8a8a;
    min-width: 40px;
  }
</style>
