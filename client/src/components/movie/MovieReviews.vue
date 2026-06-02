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
          <i
            v-for="star in 5"
            :key="star"
            class="pi"
            :class="star <= userRating ? 'pi-star-fill active' : 'pi-star'"
            @click="userRating = star"
          ></i>
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
              <i class="pi pi-user"></i>
            </div>
            <div class="user-meta">
              <span class="username">{{ item.username }}</span>
              <span class="review-date">{{ item.date }}</span>
            </div>
          </div>
          <div class="user-score">
            <i class="pi pi-star-fill"></i>
            <span>{{ item.score }}/5</span>
          </div>
        </div>
        <p class="review-comment">{{ item.comment }}</p>
      </div>
    </div>

    <div v-else class="empty-reviews">
      <i class="pi pi-inbox"></i>
      <p>ยังไม่มีรีวิวสำหรับเรื่องนี้ มารีวิวเป็นคนแรกกันเถอะ!</p>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from "vue"

  const props = defineProps<{
    movieId: number
  }>()

  const newReview = ref("")
  const userRating = ref(5)

  // Mock Data สำหรับตัวอย่างรีวิว
  const reviews = ref([
    {
      id: 1,
      username: "MovieLover99",
      date: "1 วันที่แล้ว",
      score: 5,
      comment:
        "งานภาพสวยมาก เนื้อเรื่องดำเนินเร็วสไตล์ผู้กำกับคนนี้ ไม่ผิดหวังเลยครับ แนะนำให้ดูในระบบ IMAX!",
    },
    {
      id: 2,
      username: "Reviewer_X",
      date: "3 วันที่แล้ว",
      score: 4,
      comment:
        "ครึ่งแรกทำออกมาได้น่าติดตามมาก เสียดายช่วงท้ายรวบรัดไปนิดนึง แต่โดยรวมเป็นหนังที่ดีและคุ้มค่าตั๋ว",
    },
  ])

  function submitReview() {
    if (!newReview.value.trim()) return

    reviews.value.unshift({
      id: Date.now(),
      username: "คุณ (ผู้ใช้งาน)",
      date: "เมื่อสักครู่",
      score: userRating.value,
      comment: newReview.value,
    })

    newReview.value = ""
    userRating.value = 5
  }
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
</style>
