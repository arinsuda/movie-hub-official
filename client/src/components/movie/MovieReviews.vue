<template>
  <div class="movie-reviews-component">
    <div class="write-review-box">
      <textarea
        v-model="form.body.value"
        placeholder="แบ่งปันความคิดเห็นของคุณต่อภาพยนตร์เรื่องนี้..."
        rows="3"
      ></textarea>

      <div class="extra-fields-row">
        <div class="watched-at-field">
          <label for="watched-date">วันที่คุณดู (ไม่บังคับ)</label>
          <VueDatePicker
            id="watched-date"
            v-model="form.watchedAt.value"
            :max-date="new Date()"
            :enable-time-picker="false"
            format="dd/MM/yyyy"
            :locale="th"
            dark
            auto-apply
            placeholder="เลือกวันที่"
            teleport
            menu-class-name="movie-dp-menu"
          />
        </div>

        <label class="visibility-toggle">
          <input type="checkbox" v-model="form.isPublic.value" />
          <span>{{
            form.isPublic.value ? "รีวิวสาธารณะ" : "รีวิวส่วนตัว"
          }}</span>
          <i :class="form.isPublic.value ? 'pi pi-globe' : 'pi pi-lock'"></i>
        </label>
      </div>

      <div class="form-footer">
        <div class="rating-input">
          <span>ให้คะแนนของคุณ: </span>
          <StarRatingInput v-model="form.rating.value" />
        </div>
        <button
          class="btn-submit-review"
          :disabled="!form.canSubmit.value"
          @click="handleSubmit"
        >
          {{ form.isSubmitting.value ? "กำลังส่ง..." : "ส่งรีวิว" }}
        </button>
      </div>

      <p v-if="!form.isAuthenticated.value" class="auth-hint">
        เข้าสู่ระบบเพื่อเขียนรีวิว
      </p>
    </div>

    <div v-if="isLoadingReviews" class="reviews-status">กำลังโหลดรีวิว...</div>

    <div v-else-if="reviews.length" class="reviews-list">
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
  import { computed, onMounted, ref, watch } from "vue"
  import type { ReviewResponse } from "@/types/review"
  import { VueDatePicker } from "@vuepic/vue-datepicker"
  import "@vuepic/vue-datepicker/dist/main.css"
  import { th } from "date-fns/locale"
  import StarRatingInput from "@/components/movie/Starratinginput.vue"
  import { useReviewForm } from "@/composables/useReviewForm"

  const props = withDefaults(
    defineProps<{
      movieId: number
      mediaType?: "movie" | "tv"
    }>(),
    { mediaType: "movie" },
  )

  const emit = defineEmits<{
    (e: "review-submitted"): void
  }>()

  const reviews = ref<ReviewResponse[]>([])
  const isLoadingReviews = ref(false)

  const reviewSegment = computed(() =>
    props.mediaType === "tv" ? "tv" : "movies",
  )
  const form = useReviewForm(props.movieId, props.mediaType)

  const dateFormatter = new Intl.DateTimeFormat("th-TH", {
    year: "numeric",
    month: "short",
    day: "numeric",
  })
  function formatDate(iso: string) {
    return dateFormatter.format(new Date(iso))
  }

  async function loadReviews() {
    if (!props.movieId) return
    isLoadingReviews.value = true
    try {
      const res = await reviewApi.getMediaReviews(
        reviewSegment.value,
        props.movieId,
      )
      reviews.value = res.data.reviews
    } catch (err) {
      console.error("getMediaReviews failed:", err)
    } finally {
      isLoadingReviews.value = false
    }
  }

  async function handleSubmit() {
    const created = await form.submit()
    if (!created) return
    reviews.value.unshift(created)
    emit("review-submitted")
  }

  // Re-fetch if this component instance gets reused for a different movie
  // (e.g. router reuses the view between two movie detail pages).
  watch(() => props.movieId, loadReviews)
  onMounted(loadReviews)
</script>

<style scoped>
  .movie-reviews-component {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    margin-top: 0.5rem;
    font-family: "Noto Sans Thai", sans-serif;
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

  .btn-submit-review {
    background: var(--red, #e50914);
    color: #fff;
    border: none;
    padding: 0.5rem 1.4rem;
    border-radius: 8px;
    font-weight: 600;
    font-size: 0.85rem;
    cursor: pointer;
    transition:
      background 0.2s,
      opacity 0.2s;
  }
  .btn-submit-review:hover:not(:disabled) {
    background: #f40612;
  }
  .btn-submit-review:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .auth-hint {
    margin: 0.5rem 0 0;
    font-size: 0.75rem;
    color: #8a8a8a;
  }

  .reviews-status {
    text-align: center;
    padding: 1.5rem;
    color: #8a8a8a;
    font-size: 0.875rem;
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
    border: 1px solid rgba(245, 197, 24, 0.25);
    color: var(--gold, #f5c518);
    padding: 0.3rem 0.7rem;
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

  .extra-fields-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-top: 0.85rem;
    padding-top: 0.85rem;
    border-top: 1px solid var(--border, rgba(255, 255, 255, 0.08));
    flex-wrap: wrap;
  }

  .watched-at-field {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    font-size: 0.78rem;
    color: var(--muted, #8a8a8a);
  }

  .watched-at-field label {
    white-space: nowrap;
  }

  .watched-at-field :deep(.dp__input) {
    background: rgba(0, 0, 0, 0.35);
    border: 1px solid var(--border, rgba(255, 255, 255, 0.08));
    border-radius: 6px;
    padding: 0.4rem 0.7rem 0.4rem 2.2rem;
    color: #fff;
    font-family: inherit;
    font-size: 0.8rem;
    min-height: unset;
    width: 150px;
  }

  .watched-at-field :deep(.dp__input:hover) {
    border-color: rgba(255, 255, 255, 0.2);
  }

  .watched-at-field :deep(.dp__input_focus) {
    border-color: var(--red, #e50914);
  }

  .watched-at-field :deep(.dp__icon) {
    color: var(--muted, #8a8a8a);
  }

  .visibility-toggle {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.4rem 0.85rem;
    border-radius: 9999px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid var(--border, rgba(255, 255, 255, 0.08));
    font-size: 0.78rem;
    color: var(--muted, #8a8a8a);
    cursor: pointer;
    user-select: none;
    transition:
      background 0.2s,
      border-color 0.2s,
      color 0.2s;
  }
  .visibility-toggle:hover {
    background: rgba(255, 255, 255, 0.09);
    border-color: rgba(255, 255, 255, 0.2);
  }
  .visibility-toggle input[type="checkbox"] {
    accent-color: var(--red, #e50914);
    cursor: pointer;
  }
  .visibility-toggle i {
    font-size: 0.8rem;
  }
  .visibility-toggle:has(input:checked) {
    background: rgba(245, 197, 24, 0.12);
    border-color: rgba(245, 197, 24, 0.3);
    color: var(--gold, #f5c518);
  }
  .visibility-toggle:has(input:checked) i {
    color: var(--gold, #f5c518);
  }

  /* ============================================================
     Minimal calendar redesign for VueDatePicker
     Scoped via the menu-class-name="movie-dp-menu" passed to the
     component. :deep() keeps working after teleport because the
     scoped data-v attribute stays on the teleported nodes.
     ============================================================ */

  :deep(.movie-dp-menu.dp__menu) {
    background: #111113;
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 16px;
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.45);
    padding: 0.5rem;
    font-family: inherit;
    --dp-primary-color: #e50914;
    --dp-primary-text-color: #fff;
    --dp-border-color: transparent;
    --dp-border-color-hover: transparent;
    --dp-menu-border-color: transparent;
  }

  :deep(.movie-dp-menu .dp__arrow_top) {
    display: none;
  }

  :deep(.movie-dp-menu .dp__calendar_header) {
    color: #6b6b6b;
    font-size: 0.7rem;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  :deep(.movie-dp-menu .dp__calendar_header_separator) {
    display: none;
  }

  :deep(.movie-dp-menu .dp__month_year_row) {
    padding: 0.35rem 0.4rem 0.75rem;
  }

  :deep(.movie-dp-menu .dp__month_year_select) {
    color: #fff;
    font-size: 0.9rem;
    font-weight: 600;
    border-radius: 8px;
    transition: background 0.15s;
  }
  :deep(.movie-dp-menu .dp__month_year_select:hover) {
    background: rgba(255, 255, 255, 0.06);
  }

  :deep(.movie-dp-menu .dp__inner_nav) {
    color: #8a8a8a;
    border-radius: 999px;
    transition:
      background 0.15s,
      color 0.15s;
  }
  :deep(.movie-dp-menu .dp__inner_nav:hover) {
    background: rgba(255, 255, 255, 0.08);
    color: #fff;
  }
  :deep(.movie-dp-menu .dp__inner_nav svg) {
    width: 14px;
    height: 14px;
  }

  :deep(.movie-dp-menu .dp__calendar_item) {
    padding: 1px;
  }

  :deep(.movie-dp-menu .dp__cell_inner) {
    height: 34px;
    width: 34px;
    font-size: 0.82rem;
    font-weight: 400;
    color: #d4d4d4;
    border-radius: 999px;
    border: none;
    transition:
      background 0.15s,
      color 0.15s;
  }

  :deep(.movie-dp-menu .dp__cell_inner:hover) {
    background: rgba(255, 255, 255, 0.08);
  }

  :deep(.movie-dp-menu .dp__cell_offset) {
    color: #3d3d3d;
  }

  :deep(.movie-dp-menu .dp__today) {
    border: 1px solid rgba(245, 197, 24, 0.5);
    color: #f5c518;
    font-weight: 600;
  }

  :deep(.movie-dp-menu .dp__active_date) {
    background: var(--red, #e50914);
    color: #fff;
    font-weight: 600;
  }
  :deep(.movie-dp-menu .dp__active_date:hover) {
    background: #f40612;
  }

  :deep(.movie-dp-menu .dp__cell_disabled) {
    color: #3d3d3d;
    opacity: 0.5;
  }

  :deep(.movie-dp-menu .dp__action_row) {
    border-top: 1px solid rgba(255, 255, 255, 0.06);
    padding: 0.6rem 0.4rem 0.2rem;
    margin-top: 0.35rem;
  }

  :deep(.movie-dp-menu .dp__action_button) {
    font-size: 0.78rem;
    font-weight: 600;
    border-radius: 8px;
    padding: 0.4rem 0.9rem;
  }

  :deep(.movie-dp-menu .dp__select) {
    background: var(--red, #e50914);
    color: #fff;
  }
  :deep(.movie-dp-menu .dp__select:hover) {
    background: #f40612;
  }
</style>
