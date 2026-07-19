<template>
  <section class="dashboard-section fade-in-up">
    <!-- Highlight stats grid -->
    <div class="stats-grid">
      <!-- Total Watched -->
      <div class="stat-card">
        <div class="stat-icon-box watched-theme">
          <Eye :size="22" />
        </div>
        <div class="stat-data">
          <h3>{{ $t("library.dashboard.totalWatched") }}</h3>
          <p class="stat-value">{{ watchedItems.length }}</p>
        </div>
      </div>

      <!-- Total Reviews & Trend -->
      <div class="stat-card">
        <div class="stat-icon-box review-theme">
          <Star :size="22" />
        </div>
        <div class="stat-data">
          <h3>{{ $t("library.dashboard.totalReviews") }}</h3>
          <p class="stat-value">{{ userReviews.length }}</p>
          <div v-if="userReviews.length > 0" class="trend-badge" :class="trendClass">
            <component :is="trendIcon" :size="12" />
            <span class="trend-text">
              {{
                reviewsStats.diff > 0
                  ? $t("library.dashboard.moreThanLastMonth", { count: reviewsStats.diff })
                  : reviewsStats.diff < 0
                  ? $t("library.dashboard.lessThanLastMonth", { count: Math.abs(reviewsStats.diff) })
                  : $t("library.dashboard.equalLastMonth")
              }}
            </span>
          </div>
        </div>
      </div>

      <!-- Average Rating -->
      <div class="stat-card">
        <div class="stat-icon-box rating-theme">
          <Star :size="22" />
        </div>
        <div class="stat-data">
          <h3>{{ $t("library.dashboard.averageRating") }}</h3>
          <div class="rating-display">
            <span class="stat-value">{{ averageRating > 0 ? averageRating : "—" }}</span>
            <span class="max-val">/ 5.0</span>
          </div>
          <div class="stars-row">
            <span
              v-for="s in 5"
              :key="s"
              class="star-dot"
              :class="{ 'star-dot--active': s <= Math.round(averageRating) }"
            >★</span>
          </div>
        </div>
      </div>

      <!-- Dynamic Top Genre -->
      <div class="stat-card">
        <div class="stat-icon-box genre-theme">
          <Trophy :size="22" />
        </div>
        <div class="stat-data">
          <h3>{{ $t("library.dashboard.topGenre") }}</h3>
          <p class="stat-value genre-text">{{ topGenre || $t("library.dashboard.noWatchedData") }}</p>
        </div>
      </div>
    </div>

    <!-- Charts Section (Rating Distribution + Genres Breakdown) -->
    <div class="charts-row">
      <!-- Rating Distribution CSS Bar Chart -->
      <div class="chart-card">
        <h4 class="chart-title">{{ $t("library.dashboard.averageRating") }}</h4>
        <div class="distribution-list">
          <div
            v-for="dist in ratingDistribution"
            :key="dist.rating"
            class="dist-row"
          >
            <span class="dist-label">{{ dist.rating.toFixed(1) }} ★</span>
            <div class="dist-bar-wrapper">
              <div
                class="dist-bar-fill"
                :style="{ width: `${dist.percentage}%` }"
                aria-hidden="true"
              />
            </div>
            <span class="dist-count">{{ dist.count }}</span>
          </div>
        </div>
      </div>

      <!-- Top Genres Frequency List -->
      <div class="chart-card">
        <h4 class="chart-title">{{ $t("library.dashboard.topGenre") }}</h4>
        <div class="genres-ranking">
          <div
            v-for="(genre, index) in topGenresList"
            :key="genre.name"
            class="genre-rank-row"
          >
            <div class="rank-badge" :class="`rank-${index + 1}`">{{ index + 1 }}</div>
            <span class="genre-rank-name">{{ genre.name }}</span>
            <span class="genre-rank-count">{{ genre.count }} {{ $t("library.tabs.watched") }}</span>
          </div>
          <div v-if="topGenresList.length === 0" class="genres-empty">
            {{ $t("library.dashboard.noWatchedData") }}
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { Eye, Star, Trophy, TrendingUp, TrendingDown, Minus } from "lucide-vue-next"
import type { LibraryItemResponse, ReviewResponse } from "@/types"

const props = defineProps<{
  watchedItems: LibraryItemResponse[]
  userReviews: ReviewResponse[]
  likedItems: LibraryItemResponse[]
}>()

// Dashboard Computations
const averageRating = computed(() => {
  if (props.userReviews.length === 0) return 0
  const sum = props.userReviews.reduce((acc, r) => acc + r.rating, 0)
  return Number((sum / props.userReviews.length).toFixed(1))
})

const topGenresList = computed(() => {
  const counts: Record<string, number> = {}
  const allMediaItems = [...props.watchedItems, ...props.likedItems]

  allMediaItems.forEach(item => {
    if (item.media && item.media.genres) {
      item.media.genres.forEach(g => {
        counts[g.name] = (counts[g.name] || 0) + 1
      })
    }
  })

  return Object.entries(counts)
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 5)
})

const topGenre = computed(() => {
  return topGenresList.value[0]?.name || null
})

const reviewsStats = computed(() => {
  const now = new Date()
  const currentYear = now.getFullYear()
  const currentMonth = now.getMonth()

  let thisMonthCount = 0
  let lastMonthCount = 0

  props.userReviews.forEach(r => {
    const rDate = new Date(r.created_at)
    const rYear = rDate.getFullYear()
    const rMonth = rDate.getMonth()

    if (rYear === currentYear && rMonth === currentMonth) {
      thisMonthCount++
    } else if (
      (currentMonth === 0 && rYear === currentYear - 1 && rMonth === 11) ||
      (currentMonth > 0 && rYear === currentYear && rMonth === currentMonth - 1)
    ) {
      lastMonthCount++
    }
  })

  const diff = thisMonthCount - lastMonthCount
  return { thisMonth: thisMonthCount, lastMonth: lastMonthCount, diff }
})

const trendIcon = computed(() => {
  if (reviewsStats.value.diff > 0) return TrendingUp
  if (reviewsStats.value.diff < 0) return TrendingDown
  return Minus
})

const trendClass = computed(() => {
  if (reviewsStats.value.diff > 0) return "trend-green"
  if (reviewsStats.value.diff < 0) return "trend-red"
  return "trend-neutral"
})

const ratingDistribution = computed(() => {
  const distribution = Array(10).fill(0) // 0.5 to 5.0 in steps of 0.5
  props.userReviews.forEach(r => {
    const index = Math.round(r.rating * 2) - 1
    if (index >= 0 && index < 10) {
      distribution[index]++
    }
  })
  const maxCount = Math.max(...distribution, 1)
  return distribution
    .map((count, i) => ({
      rating: (i + 1) / 2,
      count,
      percentage: (count / maxCount) * 100
    }))
    .reverse()
})
</script>

<style scoped>
.dashboard-section {
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

@media (min-width: 640px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

.stat-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1.25rem;
  backdrop-filter: blur(8px);
}

.stat-icon-box {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.watched-theme { color: #3b82f6; background: rgba(59, 130, 246, 0.1); }
.review-theme { color: #e1251b; background: rgba(225, 37, 27, 0.1); }
.rating-theme { color: #ffb800; background: rgba(255, 184, 0, 0.1); }
.genre-theme { color: #a855f7; background: rgba(168, 85, 247, 0.1); }

.stat-data h3 {
  font-size: 0.8125rem;
  color: #9ca3af;
  margin: 0 0 0.25rem 0;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 800;
  color: #ffffff;
  margin: 0;
}

.genre-text {
  font-size: 1.125rem;
  font-weight: 700;
}

.trend-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 2px 6px;
  border-radius: 6px;
  margin-top: 0.375rem;
  font-size: 0.725rem;
  font-weight: 600;
}

.trend-green { background: rgba(16, 185, 129, 0.1); color: #10b981; }
.trend-red { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.trend-neutral { background: rgba(255, 255, 255, 0.05); color: #9ca3af; }

.rating-display {
  display: flex;
  align-items: baseline;
  gap: 0.125rem;
}

.max-val {
  font-size: 0.8125rem;
  color: #9ca3af;
}

.stars-row {
  display: flex;
  gap: 0.125rem;
  margin-top: 0.25rem;
}

.star-dot {
  font-size: 0.875rem;
  color: #3f3f46;
}

.star-dot--active {
  color: #ffb800;
}

/* Charts Section */
.charts-row {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}

@media (min-width: 1024px) {
  .charts-row {
    grid-template-columns: 1fr 1fr;
  }
}

.chart-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 20px;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.chart-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
  border-left: 3px solid #e1251b;
  padding-left: 0.75rem;
}

/* Rating Distribution CSS Chart */
.distribution-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.dist-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.dist-label {
  font-size: 0.8125rem;
  color: #9ca3af;
  width: 36px;
  text-align: right;
  font-weight: 600;
}

.dist-bar-wrapper {
  flex: 1;
  height: 8px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 4px;
  overflow: hidden;
}

.dist-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #8a1612, #e1251b);
  border-radius: 4px;
}

.dist-count {
  font-size: 0.8125rem;
  color: #ffffff;
  width: 24px;
  font-weight: 700;
}

/* Top Genres List ranking */
.genres-ranking {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.genre-rank-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.rank-badge {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.8125rem;
  background: rgba(255, 255, 255, 0.05);
  color: #9ca3af;
}

.rank-1 { background: rgba(255, 184, 0, 0.15); color: #ffb800; }
.rank-2 { background: rgba(156, 163, 175, 0.15); color: #d1d5db; }
.rank-3 { background: rgba(180, 83, 9, 0.15); color: #b45309; }

.genre-rank-name {
  flex: 1;
  font-size: 0.875rem;
  font-weight: 600;
  color: #f3f4f6;
}

.genre-rank-count {
  font-size: 0.8125rem;
  color: #9ca3af;
}

.genres-empty {
  text-align: center;
  color: #52525b;
  font-size: 0.875rem;
  padding: 2rem 0;
}
</style>
