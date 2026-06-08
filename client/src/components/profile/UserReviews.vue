<template>
  <div class="reviews-wrap">
    <div class="section-head">
      <span class="section-label">My Reviews</span>
      <span class="count-badge">{{ reviews.length }}</span>
      <div class="section-rule" />
    </div>

    <div class="filter-rail" role="tablist">
      <button
        v-for="f in filters"
        :key="f.key"
        role="tab"
        :aria-selected="activeFilter === f.key"
        class="filter-btn"
        :class="{ 'filter-btn--active': activeFilter === f.key }"
        @click="activeFilter = f.key"
      >
        {{ f.label }}
      </button>
    </div>

    <div v-if="loading" class="state-loading">
      <div class="loading-bar"><div class="loading-fill" /></div>
    </div>

    <div v-else-if="filteredReviews.length === 0" class="state-empty">
      <Star :size="32" :stroke-width="1.2" aria-hidden="true" />
      <p>No reviews found</p>
    </div>

    <div v-else class="reviews-list">
      <article
        v-for="review in filteredReviews"
        :key="review.id"
        class="review-card"
      >
        <div class="poster">
          <img
            v-if="review.posterUrl"
            :src="review.posterUrl"
            :alt="review.targetName"
            class="poster-img"
          />
          <div v-else class="poster-fallback">
            <Film :size="18" :stroke-width="1.5" aria-hidden="true" />
          </div>
        </div>

        <div class="review-body">
          <div class="review-meta-top">
            <h4 class="review-title">{{ review.targetName }}</h4>
            <span class="review-date">{{ review.createdAt }}</span>
          </div>

          <p class="review-content">{{ review.content }}</p>

          <div v-if="review.tags?.length" class="review-tags">
            <span v-for="tag in review.tags" :key="tag" class="tag">
              {{ tag }}
            </span>
          </div>

          <div class="review-actions">
            <button class="action-btn" @click="handleEdit(review.id)">
              <Pencil :size="12" aria-hidden="true" />
              Edit
            </button>
            <button
              class="action-btn action-btn--del"
              @click="handleDelete(review.id)"
            >
              <Trash2 :size="12" aria-hidden="true" />
              Delete
            </button>
          </div>
        </div>

        <div class="review-rating">
          <span class="rating-score">{{ review.rating.toFixed(1) }}</span>
          <div
            class="rating-stars"
            :aria-label="`Rating ${review.rating} out of 5`"
          >
            <Star
              v-for="i in 5"
              :key="i"
              :size="10"
              :class="i <= review.rating ? 'star-fill' : 'star-empty'"
              aria-hidden="true"
            />
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, onMounted, ref } from "vue"
  import { Star, Film, Pencil, Trash2 } from "lucide-vue-next"

  const props = defineProps<{ userId: number }>()

  const loading = ref(false)
  type FilterKey = "all" | "high" | "low"
  const activeFilter = ref<FilterKey>("all")

  const filters: { key: FilterKey; label: string }[] = [
    { key: "all", label: "All Reviews" },
    { key: "high", label: "Highly Rated" },
    { key: "low", label: "Critical" },
  ]

  interface ReviewItem {
    id: number
    targetName: string
    rating: number
    content: string
    createdAt: string
    posterUrl?: string
    tags?: string[]
  }

  const reviews = ref<ReviewItem[]>([])

  const filteredReviews = computed(() => {
    if (activeFilter.value === "high")
      return reviews.value.filter(r => r.rating >= 4)
    if (activeFilter.value === "low")
      return reviews.value.filter(r => r.rating < 4)
    return reviews.value
  })

  onMounted(async () => {
    try {
      loading.value = true
      await new Promise(resolve => setTimeout(resolve, 400))
      reviews.value = [
        {
          id: 1,
          targetName: "The Shawshank Redemption",
          rating: 5,
          content:
            "หนังที่ดูกี่ครั้งก็ยังซึ้ง เรื่องราวของความหวังในสถานที่ที่มืดที่สุด Morgan Freeman บรรยายได้ดีมาก",
          createdAt: "15 May 2026",
          posterUrl:
            "https://image.tmdb.org/t/p/w200/q6y0Go1tsGEsmtFryDOJo3dEmqu.jpg",
          tags: ["DRAMA", "CLASSIC"],
        },
        {
          id: 2,
          targetName: "The Godfather",
          rating: 4.5,
          content:
            "Masterpiece ที่ไม่มีใครเถียงได้ การแสดงของ Marlon Brando ยังทรงพลังทุกครั้งที่ดูซ้ำ",
          createdAt: "20 Apr 2026",
          posterUrl:
            "https://image.tmdb.org/t/p/w200/3bhkrj58Vtu7enYsLe1rLtozaGT.jpg",
          tags: ["CRIME", "EPIC"],
        },
        {
          id: 3,
          targetName: "Advanced Go Backend — Course",
          rating: 3.5,
          content:
            "เนื้อหาแน่นมาก ได้เข้าใจ Concurrency กับ Clean Architecture แบบจริงจัง แต่หักคะแนนเสียงไมค์",
          createdAt: "12 Apr 2026",
          tags: ["COURSE", "BACKEND"],
        },
      ]
    } catch (err) {
      console.error("Fetch reviews failed:", err)
    } finally {
      loading.value = false
    }
  })

  function handleEdit(id: number) {
    console.log("Edit review ID:", id)
  }

  function handleDelete(id: number) {
    if (!confirm("Are you sure you want to delete this review?")) return
    reviews.value = reviews.value.filter(r => r.id !== id)
  }
</script>

<style scoped>
  .reviews-wrap {
    --surface: #121212;
    --surface-card: #1a1a1a;
    --brand-red: #e50914;
    --brand-red-hover: #b20710;
    --border: rgba(255, 255, 255, 0.05);
    --border-hover: rgba(255, 255, 255, 0.1);
    --text: #f5f5f7;
    --sub: #8e8e93;
    --muted: #48484a;
    --ui:
      -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue",
      system-ui, sans-serif;

    color: var(--text);
    font-family: var(--ui);
  }

  .section-head {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }

  .section-label {
    font-size: 0.7rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--muted);
    white-space: nowrap;
  }

  .count-badge {
    font-size: 0.7rem;
    font-weight: 600;
    color: #ffffff;
    background: var(--brand-red);
    padding: 2px 8px;
    border-radius: 4px;
  }

  .section-rule {
    flex: 1;
    height: 1px;
    background: var(--border);
  }

  /* ── Filter Rail ── */
  .filter-rail {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1.75rem;
  }

  .filter-btn {
    font-family: var(--ui);
    font-size: 0.78rem;
    font-weight: 500;
    color: var(--sub);
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--border);
    padding: 0.45rem 1rem;
    border-radius: 20px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .filter-btn:hover {
    color: #ffffff;
    border-color: var(--border-hover);
    background: rgba(255, 255, 255, 0.06);
  }

  .filter-btn--active {
    color: #000000;
    background: #ffffff;
    border-color: #ffffff;
  }

  .filter-btn--active:hover {
    color: #000000;
    background: #ffffff;
    border-color: #ffffff;
  }

  /* ── Loader & Empty State ── */
  .state-loading {
    padding: 4rem 0;
    display: flex;
    justify-content: center;
  }

  .loading-bar {
    width: 120px;
    height: 2px;
    background: var(--border);
    overflow: hidden;
    position: relative;
    border-radius: 2px;
  }

  .loading-fill {
    height: 100%;
    width: 40%;
    background: #ffffff;
    position: absolute;
    animation: sweep 1.5s infinite ease-in-out;
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
    gap: 0.75rem;
    padding: 5rem 0;
    color: var(--muted);
  }

  .state-empty p {
    font-size: 0.85rem;
    margin: 0;
    color: var(--sub);
  }

  /* ── Reviews Card Layout ── */
  .reviews-list {
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .review-card {
    display: grid;
    grid-template-columns: 56px 1fr auto;
    gap: 1.25rem;
    padding: 1.5rem 0.5rem;
    border-bottom: 1px solid var(--border);
    align-items: start;
    transition: background 0.2s ease;
  }

  .review-card:last-child {
    border-bottom: none;
  }

  /* Poster Media Layout */
  .poster {
    width: 56px;
    height: 80px;
    background: var(--surface-card);
    border: 1px solid var(--border);
    border-radius: 6px;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }

  .poster-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .poster-fallback {
    color: var(--muted);
  }

  /* Content Core Layout */
  .review-body {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .review-meta-top {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .review-title {
    font-size: 0.92rem;
    font-weight: 600;
    color: #ffffff;
    margin: 0;
    line-height: 1.3;
  }

  .review-date {
    font-size: 0.7rem;
    color: var(--sub);
  }

  .review-content {
    font-size: 0.82rem;
    color: var(--text);
    line-height: 1.55;
    margin: 0.2rem 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  /* Tags rendering styling */
  .review-tags {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
    margin: 0.15rem 0;
  }

  .tag {
    font-size: 0.62rem;
    font-weight: 600;
    letter-spacing: 0.04em;
    color: var(--sub);
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border);
    padding: 2px 8px;
    border-radius: 4px;
  }

  /* Interactive action links */
  .review-actions {
    display: flex;
    gap: 1rem;
    margin-top: 0.25rem;
  }

  .action-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 0.72rem;
    font-weight: 500;
    font-family: var(--ui);
    color: var(--sub);
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    transition: color 0.2s ease;
  }

  .action-btn:hover {
    color: #ffffff;
  }

  .action-btn--del:hover {
    color: #ff453a; /* เฉพาะปุ่มลบให้ขึ้นสีแดงเมื่อโฮเวอร์ */
  }

  /* Rating Badge Component styling */
  .review-rating {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 6px;
    flex-shrink: 0;
  }

  .rating-score {
    font-size: 1.4rem;
    font-weight: 700;
    color: #ffffff;
    line-height: 1;
    letter-spacing: -0.02em;
  }

  .rating-stars {
    display: flex;
    gap: 3px;
  }

  .star-fill {
    color: #ffffff; /* เปลี่ยนเป็นดาวสีขาวสว่าง เรียบหรูสไตล์ฟิล์ม */
  }

  .star-empty {
    color: var(--muted);
  }

  /* Mobile Grid Adaptation */
  @media (max-width: 580px) {
    .review-card {
      grid-template-columns: 56px 1fr;
      row-gap: 1rem;
    }
    .review-rating {
      grid-column: 2;
      flex-direction: row;
      align-items: center;
      gap: 10px;
      justify-self: start;
    }
  }
</style>
