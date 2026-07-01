<template>
  <div class="reviews-root">
    <div class="section-head">
      <span class="eyebrow">My Reviews</span>
      <span class="count-chip">{{ reviews.length }}</span>
      <div class="rule" />
    </div>

    <!-- Filters -->
    <div class="filter-rail" role="tablist">
      <button
        v-for="f in filters"
        :key="f.key"
        class="filter-pill"
        :class="{ 'filter-pill--active': activeFilter === f.key }"
        :aria-selected="activeFilter === f.key"
        role="tab"
        @click="activeFilter = f.key"
      >
        {{ f.label }}
      </button>
    </div>

    <!-- States -->
    <div v-if="loading" class="state-loading">
      <div class="loader-bar"><div class="loader-fill" /></div>
    </div>

    <div v-else-if="filteredReviews.length === 0" class="state-empty">
      <Star :size="28" :stroke-width="1.2" />
      <p>No reviews found</p>
    </div>

    <!-- List -->
    <div v-else class="reviews-list">
      <article
        v-for="(review, i) in displayReviews"
        :key="review.id"
        class="review-card"
        :style="{ '--i': i }"
      >
        <!-- Poster -->
        <div
          class="poster"
          @click="goToDetail(review.mediaId, review.mediaType)"
        >
          <img
            v-if="review.posterUrl"
            :src="review.posterUrl"
            :alt="review.targetName"
            class="poster-img"
          />
          <div v-else class="poster-fallback">
            <Film :size="16" :stroke-width="1.4" />
          </div>
        </div>

        <!-- Body -->
        <div class="review-body">
          <div class="title-rating-row">
            <h4
              class="review-title"
              @click="goToDetail(review.mediaId, review.mediaType)"
            >
              {{ review.targetName }}
            </h4>
            <div class="inline-rating">
              <span class="inline-rating-num">{{
                review.rating.toFixed(1)
              }}</span>
              <div class="stars">
                <Star
                  v-for="n in 5"
                  :key="n"
                  :size="9"
                  :class="n <= review.rating ? 'star-on' : 'star-off'"
                />
              </div>
            </div>
          </div>
          <time class="review-date">{{ review.createdAt }}</time>
          <p class="review-content">{{ review.content }}</p>
          <div v-if="review.tags?.length" class="tag-row">
            <span v-for="tag in review.tags" :key="tag" class="tag">{{
              tag
            }}</span>
          </div>
        </div>

        <!-- Actions Menu -->
        <div class="menu-wrap" @click.stop>
          <button
            class="menu-trigger"
            :class="{ 'menu-trigger--open': openMenuId === review.id }"
            @click="toggleMenu(review.id)"
            :aria-label="'Options for ' + review.targetName"
          >
            <MoreHorizontal :size="14" />
          </button>
          <Transition name="dropdown">
            <div v-if="openMenuId === review.id" class="dropdown">
              <button class="dropdown-item" @click="openEdit(review)">
                <Pencil :size="12" /> Edit
              </button>
              <button
                class="dropdown-item dropdown-item--del"
                @click="openDelete(review)"
              >
                <Trash2 :size="12" /> Delete
              </button>
            </div>
          </Transition>
        </div>
      </article>
    </div>

    <!-- Modals -->
    <EditReview
      v-if="editTarget"
      :userId="userId"
      :review="editTarget"
      @close="editTarget = null"
      @saved="onReviewSaved"
    />
    <DeleteReview
      v-if="deleteTarget"
      :userId="userId"
      :review="deleteTarget"
      @close="deleteTarget = null"
      @deleted="onReviewDeleted"
    />
  </div>
</template>

<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from "vue"
  import { Star, Film, Pencil, Trash2, MoreHorizontal } from "lucide-vue-next"
  import { reviewApi } from "@/api/api"
  import { useRouter } from "vue-router"
  import type { ReviewResponse } from "@/types"
  import EditReview from "./components/EditReview.vue"
  import DeleteReview from "./components/DeleteReview.vue"

  const router = useRouter()

  const props = defineProps<{ userId: number }>()

  const loading = ref(false)
  type FilterKey = "all" | "high" | "low"
  const activeFilter = ref<FilterKey>("all")
  const filters: { key: FilterKey; label: string }[] = [
    { key: "all", label: "All" },
    { key: "high", label: "High Rated" },
    { key: "low", label: "Critical" },
  ]

  const reviews = ref<ReviewResponse[]>([])
  const openMenuId = ref<number | null>(null)
  const editTarget = ref<ReviewResponse | null>(null)
  const deleteTarget = ref<ReviewResponse | null>(null)

  function goToDetail(mediaId: number, mediaType: string) {
    router.push({
      name: mediaType === "tv" ? "tv-detail" : "movie-detail",
      params: { id: mediaId },
    })
  }

  function toggleMenu(id: number) {
    openMenuId.value = openMenuId.value === id ? null : id
  }
  function closeMenu() {
    openMenuId.value = null
  }
  function openEdit(review: ReviewResponse) {
    closeMenu()
    editTarget.value = review
  }
  function openDelete(review: ReviewResponse) {
    closeMenu()
    deleteTarget.value = review
  }
  function onReviewSaved(updated: ReviewResponse) {
    const idx = reviews.value.findIndex(r => r.id === updated.id)
    if (idx !== -1) reviews.value[idx] = updated
    editTarget.value = null
  }
  function onReviewDeleted(id: number) {
    reviews.value = reviews.value.filter(r => r.id !== id)
    deleteTarget.value = null
  }

  onMounted(() => document.addEventListener("click", closeMenu))
  onUnmounted(() => document.removeEventListener("click", closeMenu))

  const filteredReviews = computed(() => {
    if (activeFilter.value === "high")
      return reviews.value.filter(r => r.rating >= 4)
    if (activeFilter.value === "low")
      return reviews.value.filter(r => r.rating < 4)
    return reviews.value
  })

  const displayReviews = computed(() =>
    filteredReviews.value.map(r => ({
      ...r,
      targetName: r.media.title || `#${r.media.id}`,
      posterUrl: r.media.poster_url || null,
      mediaId: r.media.id,
      mediaType: r.media.media_type,
      createdAt: new Date(r.created_at).toLocaleDateString("th-TH", {
        year: "numeric",
        month: "short",
        day: "numeric",
      }),
      content: r.body,
      tags: [] as string[],
    })),
  )

  onMounted(async () => {
    try {
      loading.value = true
      const res = await reviewApi.getUserReviews(props.userId)
      reviews.value = res.data.reviews ?? []
    } catch (err) {
      console.error("Fetch reviews failed:", err)
    } finally {
      loading.value = false
    }
  })
</script>

<style scoped>
  .reviews-root {
    --c-card: #161616;
    --c-border: rgba(255, 255, 255, 0.06);
    --c-border-h: rgba(255, 255, 255, 0.12);
    --c-red: #e1251b;
    --c-text: #f0f0f0;
    --c-sub: #8a8a8e;
    --c-muted: #3a3a3c;
    --c-star: #fbbf24;
    --font:
      "Inter", -apple-system, BlinkMacSystemFont, "SF Pro Text", system-ui,
      sans-serif;
    --ease: cubic-bezier(0.16, 1, 0.3, 1);
    font-family: var(--font);
    color: var(--c-text);
  }

  .section-head {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 20px;
  }
  .eyebrow {
    font-size: 0.6rem;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--c-muted);
    white-space: nowrap;
  }
  .count-chip {
    font-size: 0.65rem;
    font-weight: 700;
    background: var(--c-red);
    color: #fff;
    padding: 2px 7px;
    border-radius: 4px;
  }
  .rule {
    flex: 1;
    height: 1px;
    background: var(--c-border);
  }

  /* Filters */
  .filter-rail {
    display: flex;
    gap: 6px;
    margin-bottom: 24px;
  }
  .filter-pill {
    font-family: var(--font);
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--c-sub);
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--c-border);
    padding: 6px 14px;
    border-radius: 20px;
    cursor: pointer;
    transition: all 0.2s var(--ease);
  }
  .filter-pill:hover {
    color: var(--c-text);
    border-color: var(--c-border-h);
  }
  .filter-pill--active {
    color: #000;
    background: #fff;
    border-color: #fff;
  }
  .filter-pill--active:hover {
    color: #000;
  }

  /* States */
  .state-loading {
    padding: 48px 0;
    display: flex;
    justify-content: center;
  }
  .loader-bar {
    width: 100px;
    height: 2px;
    background: var(--c-border);
    border-radius: 2px;
    overflow: hidden;
    position: relative;
  }
  .loader-fill {
    height: 100%;
    width: 40%;
    background: #fff;
    position: absolute;
    animation: sweep 1.4s infinite ease-in-out;
  }
  @keyframes sweep {
    0% {
      left: -40%;
    }
    100% {
      left: 100%;
    }
  }

  .state-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    padding: 56px 0;
    color: var(--c-muted);
  }
  .state-empty p {
    font-size: 0.82rem;
    margin: 0;
    color: var(--c-sub);
  }

  /* List */
  .reviews-list {
    display: flex;
    flex-direction: column;
  }

  .review-card {
    display: grid;
    grid-template-columns: 56px 1fr auto;
    gap: 0 16px;
    padding: 20px 0;
    border-bottom: 1px solid var(--c-border);
    align-items: start;
    position: relative;
    animation: cardIn 0.4s var(--ease) calc(var(--i) * 50ms) both;
    transition: background 0.15s;
  }
  .review-card:last-child {
    border-bottom: none;
  }
  .review-card:hover {
    background: rgba(255, 255, 255, 0.015);
    padding-left: 8px;
    padding-right: 8px;
    margin: 0 -8px;
    border-radius: 10px;
  }

  @keyframes cardIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .poster {
    width: 56px;
    height: 80px;
    background: var(--c-card);
    border: 1px solid var(--c-border);
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    cursor: pointer;
  }
  .poster-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .poster-fallback {
    color: var(--c-muted);
  }

  .review-body {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }

  .title-rating-row {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
  }
  .review-title {
    font-size: 0.9rem;
    font-weight: 600;
    color: #fff;
    margin: 0;
    line-height: 1.3;
    cursor: pointer;
    transition: color 0.15s;
  }
  .review-title:hover {
    color: var(--c-red);
  }

  .inline-rating {
    display: flex;
    align-items: center;
    gap: 5px;
    flex-shrink: 0;
  }
  .inline-rating-num {
    font-size: 0.82rem;
    font-weight: 700;
    color: var(--c-star);
    font-variant-numeric: tabular-nums;
    letter-spacing: -0.02em;
  }
  .stars {
    display: flex;
    gap: 2px;
  }
  .star-on {
    color: var(--c-star);
    fill: var(--c-star);
  }
  .star-off {
    color: var(--c-muted);
  }

  .review-date {
    font-size: 0.68rem;
    color: var(--c-sub);
  }
  .review-content {
    font-size: 0.8rem;
    color: rgba(240, 240, 240, 0.72);
    line-height: 1.58;
    margin: 2px 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  .tag-row {
    display: flex;
    gap: 5px;
    flex-wrap: wrap;
  }
  .tag {
    font-size: 0.6rem;
    font-weight: 600;
    letter-spacing: 0.04em;
    color: var(--c-sub);
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--c-border);
    padding: 2px 7px;
    border-radius: 4px;
  }

  /* Menu */
  .menu-wrap {
    position: relative;
    align-self: start;
  }
  .menu-trigger {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 6px;
    background: transparent;
    border: 1px solid transparent;
    color: var(--c-sub);
    cursor: pointer;
    transition: all 0.15s var(--ease);
  }
  .menu-trigger:hover,
  .menu-trigger--open {
    color: var(--c-text);
    background: rgba(255, 255, 255, 0.06);
    border-color: var(--c-border-h);
  }

  .dropdown {
    position: absolute;
    top: calc(100% + 6px);
    right: 0;
    min-width: 120px;
    background: #1e1e1e;
    border: 1px solid var(--c-border-h);
    border-radius: 8px;
    padding: 4px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
    z-index: 50;
    transform-origin: top right;
  }
  .dropdown-item {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 7px 10px;
    border-radius: 5px;
    font-family: var(--font);
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--c-sub);
    background: none;
    border: none;
    cursor: pointer;
    text-align: left;
    transition: all 0.12s;
  }
  .dropdown-item:hover {
    color: #fff;
    background: rgba(255, 255, 255, 0.06);
  }
  .dropdown-item--del:hover {
    color: #ff453a;
    background: rgba(255, 69, 58, 0.08);
  }

  .dropdown-enter-active {
    transition:
      opacity 0.15s var(--ease),
      transform 0.15s var(--ease);
  }
  .dropdown-leave-active {
    transition:
      opacity 0.1s ease,
      transform 0.1s ease;
  }
  .dropdown-enter-from,
  .dropdown-leave-to {
    opacity: 0;
    transform: scale(0.92);
  }

  @media (max-width: 520px) {
    .review-card {
      grid-template-columns: 56px 1fr auto;
    }
    .title-rating-row {
      flex-direction: column;
      align-items: flex-start;
      gap: 4px;
    }
  }
</style>
