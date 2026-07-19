<template>
  <div class="library-root">
    <div class="grain" aria-hidden="true" />

    <!-- Loading State -->
    <div v-if="loading" class="library-loading">
      <div class="loading-ring" />
      <span class="loading-text">{{ $t("common.loading") }}</span>
    </div>

    <!-- Main Layout -->
    <div v-else class="library-layout">
      <!-- Back to Profile Link -->
      <router-link :to="`/users/${userId}`" class="back-link">
        <ChevronLeft :size="16" />
        <span>{{ $t("navigation.profile") }}</span>
      </router-link>

      <!-- Profile Header Summary Card -->
      <header class="profile-summary-card">
        <div class="user-info-section">
          <div class="avatar-wrapper">
            <div class="avatar-ring">
              <img
                v-if="profileUser?.avatar_url"
                :src="profileUser.avatar_url"
                :alt="profileUser.username"
                class="avatar-img"
              />
              <div v-else class="avatar-fallback">
                <UserIcon :size="28" />
              </div>
            </div>
            <span class="level-chip">Lv.{{ profileUser?.level ?? 1 }}</span>
          </div>
          <div class="user-meta">
            <h1 class="user-name">{{ profileUser?.display_name || profileUser?.username }}</h1>
            <p class="user-tag">@{{ profileUser?.username }}</p>
            <p class="user-bio">{{ profileUser?.bio || $t("profile.noBio") }}</p>
          </div>
        </div>
        <div class="user-counts-section">
          <div class="count-box">
            <span class="count-val">{{ profileUser?.review_count ?? 0 }}</span>
            <span class="count-lbl">{{ $t("library.tabs.reviews") }}</span>
          </div>
          <div class="count-box">
            <span class="count-val">{{ profileUser?.follower_count ?? 0 }}</span>
            <span class="count-lbl">{{ $t("profile.stats.followers") }}</span>
          </div>
          <div class="count-box">
            <span class="count-val">{{ profileUser?.following_count ?? 0 }}</span>
            <span class="count-lbl">{{ $t("profile.stats.following") }}</span>
          </div>
        </div>
      </header>

      <!-- Navigation Tabs -->
      <nav class="library-tabs">
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'dashboard' }"
          @click="activeTab = 'dashboard'"
        >
          <Flame :size="16" />
          <span>{{ $t("library.dashboard.title") }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'watchlist' }"
          @click="activeTab = 'watchlist'"
        >
          <Bookmark :size="16" />
          <span>{{ $t("library.tabs.watchlist") }}</span>
          <span class="tab-count">{{ watchlistItems.length }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'watched' }"
          @click="activeTab = 'watched'"
        >
          <Eye :size="16" />
          <span>{{ $t("library.tabs.watched") }}</span>
          <span class="tab-count">{{ watchedItems.length }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'likes' }"
          @click="activeTab = 'likes'"
        >
          <Heart :size="16" />
          <span>{{ $t("library.tabs.likes") }}</span>
          <span class="tab-count">{{ likedItems.length }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === 'reviews' }"
          @click="activeTab = 'reviews'"
        >
          <Star :size="16" />
          <span>{{ $t("library.tabs.reviews") }}</span>
          <span class="tab-count">{{ userReviews.length }}</span>
        </button>
      </nav>

      <!-- TAB CONTENT: Dashboard -->
      <section v-if="activeTab === 'dashboard'" class="dashboard-section fade-in-up">
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

      <!-- TAB CONTENT: Lists (Watchlist, Watched, Likes, Reviews) -->
      <section v-else class="list-section fade-in-up">
        <!-- Advanced Filter Bar -->
        <div class="filter-bar">
          <!-- Title Search -->
          <div class="filter-input-wrapper search-wrapper">
            <Search :size="16" class="filter-icon" />
            <input
              v-model="filters.search"
              type="text"
              :placeholder="$t('library.filters.search')"
              class="filter-input"
            />
          </div>

          <!-- Media Type Select -->
          <div class="filter-select-wrapper">
            <Film :size="14" class="select-icon" />
            <select v-model="filters.mediaType" class="filter-select">
              <option value="all">{{ $t("library.filters.allMedia") }}</option>
              <option value="movie">{{ $t("library.filters.movies") }}</option>
              <option value="tv">{{ $t("library.filters.tvSeries") }}</option>
            </select>
          </div>

          <!-- Genre Filter -->
          <div class="filter-select-wrapper">
            <ArrowUpDown :size="14" class="select-icon" />
            <select v-model="filters.genre" class="filter-select">
              <option value="all">{{ $t("library.filters.allGenres") }}</option>
              <option
                v-for="genre in availableGenres"
                :key="genre"
                :value="genre"
              >{{ genre }}</option>
            </select>
          </div>

          <!-- Sort Order -->
          <div class="filter-select-wrapper">
            <ArrowUpDown :size="14" class="select-icon" />
            <select v-model="filters.sortBy" class="filter-select">
              <option value="newest">{{ $t("library.filters.sort.newest") }}</option>
              <option value="oldest">{{ $t("library.filters.sort.oldest") }}</option>
              <option v-if="activeTab === 'reviews'" value="rating_high">{{ $t("library.filters.sort.ratingHigh") }}</option>
              <option v-if="activeTab === 'reviews'" value="rating_low">{{ $t("library.filters.sort.ratingLow") }}</option>
              <option value="title_az">{{ $t("library.filters.sort.titleAZ") }}</option>
              <option value="title_za">{{ $t("library.filters.sort.titleZA") }}</option>
            </select>
          </div>
        </div>

        <!-- Render List Items -->
        <!-- Tab: Reviews -->
        <div v-if="activeTab === 'reviews'" class="reviews-list">
          <div
            v-for="review in filteredReviews"
            :key="review.id"
            class="review-library-card"
          >
            <router-link
              :to="review.media.media_type === 'tv' ? `/tv/${review.media.id}` : `/movies/${review.media.id}`"
              class="review-media-poster"
            >
              <img
                v-if="review.media.poster_url"
                :src="`${TMDB_IMG}${review.media.poster_url}`"
                :alt="review.media.title"
                loading="lazy"
              />
              <div v-else class="poster-fallback">
                <Film :size="18" />
              </div>
            </router-link>
            <div class="review-card-content">
              <div class="review-card-header">
                <router-link
                  :to="review.media.media_type === 'tv' ? `/tv/${review.media.id}` : `/movies/${review.media.id}`"
                  class="review-media-title"
                >
                  {{ review.media.title }}
                </router-link>
                <div class="review-card-stars">
                  <span class="rating-num">{{ review.rating.toFixed(1) }}</span>
                  <span class="stars-icon">★</span>
                </div>
              </div>
              <p class="review-card-body">{{ review.body }}</p>
              <div class="review-card-footer">
                <span class="review-date">{{ new Date(review.created_at).toLocaleDateString() }}</span>
                <div class="review-meta-actions">
                  <span class="action-count"><Heart :size="12" class="liked-heart" /> {{ review.like_count }}</span>
                </div>
              </div>
            </div>
          </div>
          <div v-if="filteredReviews.length === 0" class="empty-state-card">
            <Star :size="32" class="empty-icon" />
            <p>{{ $t("library.empty_reviews") }}</p>
          </div>
        </div>

        <!-- Tab: Library Items (Watchlist, Watched, Likes) -->
        <div v-else class="media-grid">
          <div
            v-for="item in filteredLibraryItems"
            :key="item.id"
            class="media-library-card"
          >
            <router-link
              :to="item.media.media_type === 'tv' ? `/tv/${item.media.id}` : `/movies/${item.media.id}`"
              class="media-poster-frame"
            >
              <img
                v-if="item.media.poster_url"
                :src="`${TMDB_IMG}${item.media.poster_url}`"
                :alt="item.media.title"
                loading="lazy"
              />
              <div v-else class="poster-fallback">
                <Film :size="24" />
              </div>
              <span class="media-type-chip">{{ item.media.media_type === 'tv' ? 'TV' : 'Movie' }}</span>
            </router-link>
            <div class="media-card-meta">
              <router-link
                :to="item.media.media_type === 'tv' ? `/tv/${item.media.id}` : `/movies/${item.media.id}`"
                class="media-title"
              >
                {{ item.media.title }}
              </router-link>
              <div class="media-details-row">
                <span class="media-rating-badge">★ {{ item.media.vote_average.toFixed(1) }}</span>
                <span class="save-date">{{ new Date(item.created_at).toLocaleDateString() }}</span>
              </div>
              <!-- User tags if present -->
              <div v-if="item.tags && item.tags.length > 0" class="tags-wrapper">
                <span v-for="t in item.tags" :key="t" class="tag-badge">{{ t }}</span>
              </div>
            </div>
          </div>
          <div v-if="filteredLibraryItems.length === 0" class="empty-state-card">
            <Bookmark :size="32" class="empty-icon" />
            <p>{{ $t("library.empty") }}</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRoute } from "vue-router"
import { libraryApi, reviewApi, userApi } from "@/api/api"
import type { UserProfile } from "@/types/user"
import type { LibraryItemResponse, ReviewResponse } from "@/types"
import {
  Search,
  Bookmark,
  Eye,
  Heart,
  Star,
  ArrowUpDown,
  Flame,
  Trophy,
  ChevronLeft,
  User as UserIcon,
  Film,
  TrendingUp,
  TrendingDown,
  Minus
} from "lucide-vue-next"

const route = useRoute()

const TMDB_IMG = "https://image.tmdb.org/t/p/w342"

const userId = computed(() => {
  const raw = route.params.userId
  return Number(Array.isArray(raw) ? raw[0] : raw)
})

const loading = ref(true)
const activeTab = ref<"dashboard" | "watchlist" | "watched" | "likes" | "reviews">("dashboard")

const profileUser = ref<UserProfile | null>(null)
const watchlistItems = ref<LibraryItemResponse[]>([])
const watchedItems = ref<LibraryItemResponse[]>([])
const likedItems = ref<LibraryItemResponse[]>([])
const userReviews = ref<ReviewResponse[]>([])

// Filter States
const filters = ref({
  search: "",
  mediaType: "all",
  genre: "all",
  sortBy: "newest"
})

// Load everything on mount
const loadData = async () => {
  try {
    loading.value = true
    const id = userId.value

    // Fetch profile info
    const profileRes = await userApi.getProfile(id)
    profileUser.value = profileRes.data.user

    // Fetch lists in parallel
    const [watchlistRes, watchedRes, likedRes, reviewsRes] = await Promise.all([
      libraryApi.getVisibleUserLibrary(id, { list_type: "watchlist" }),
      libraryApi.getVisibleUserLibrary(id, { list_type: "watched" }),
      libraryApi.getVisibleUserLibrary(id, { list_type: "likes" }),
      reviewApi.getUserReviews(id)
    ])

    watchlistItems.value = watchlistRes.data.items
    watchedItems.value = watchedRes.data.items
    likedItems.value = likedRes.data.items
    userReviews.value = reviewsRes.data.reviews
  } catch (err) {
    console.error("Failed to load library data:", err)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
watch(userId, loadData)

// Dashboard Computations
const averageRating = computed(() => {
  if (userReviews.value.length === 0) return 0
  const sum = userReviews.value.reduce((acc, r) => acc + r.rating, 0)
  return Number((sum / userReviews.value.length).toFixed(1))
})

const topGenresList = computed(() => {
  const counts: Record<string, number> = {}
  const allMediaItems = [...watchedItems.value, ...likedItems.value]

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

  userReviews.value.forEach(r => {
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
  userReviews.value.forEach(r => {
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

// Dynamic Genres for Dropdown Filter
const availableGenres = computed(() => {
  const genresMap = new Set<string>()
  const items =
    activeTab.value === "reviews"
      ? userReviews.value.map(r => r.media)
      : activeTab.value === "watchlist"
      ? watchlistItems.value.map(i => i.media)
      : activeTab.value === "watched"
      ? watchedItems.value.map(i => i.media)
      : likedItems.value.map(i => i.media)

  items.forEach(media => {
    if (media && media.genres) {
      media.genres.forEach(g => genresMap.add(g.name))
    }
  })

  return Array.from(genresMap).sort()
})

// Filter & Sort Logic for Library Items
const filteredLibraryItems = computed(() => {
  let list =
    activeTab.value === "watchlist"
      ? watchlistItems.value
      : activeTab.value === "watched"
      ? watchedItems.value
      : likedItems.value

  // 1. Search Query
  if (filters.value.search.trim()) {
    const query = filters.value.search.toLowerCase()
    list = list.filter(item => item.media.title.toLowerCase().includes(query))
  }

  // 2. Media Type
  if (filters.value.mediaType !== "all") {
    list = list.filter(item => item.media.media_type === filters.value.mediaType)
  }

  // 3. Genre
  if (filters.value.genre !== "all") {
    list = list.filter(
      item => item.media.genres && item.media.genres.some(g => g.name === filters.value.genre)
    )
  }

  // 4. Sort
  const sorted = [...list]
  if (filters.value.sortBy === "newest") {
    sorted.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  } else if (filters.value.sortBy === "oldest") {
    sorted.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
  } else if (filters.value.sortBy === "title_az") {
    sorted.sort((a, b) => a.media.title.localeCompare(b.media.title))
  } else if (filters.value.sortBy === "title_za") {
    sorted.sort((a, b) => b.media.title.localeCompare(a.media.title))
  }

  return sorted
})

// Filter & Sort Logic for Reviews
const filteredReviews = computed(() => {
  let list = userReviews.value

  // 1. Search Query
  if (filters.value.search.trim()) {
    const query = filters.value.search.toLowerCase()
    list = list.filter(r => r.media.title.toLowerCase().includes(query))
  }

  // 2. Media Type
  if (filters.value.mediaType !== "all") {
    list = list.filter(r => r.media.media_type === filters.value.mediaType)
  }

  // 3. Genre
  if (filters.value.genre !== "all") {
    list = list.filter(
      r => r.media.genres && r.media.genres.some(g => g.name === filters.value.genre)
    )
  }

  // 4. Sort
  const sorted = [...list]
  if (filters.value.sortBy === "newest") {
    sorted.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  } else if (filters.value.sortBy === "oldest") {
    sorted.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
  } else if (filters.value.sortBy === "rating_high") {
    sorted.sort((a, b) => b.rating - a.rating)
  } else if (filters.value.sortBy === "rating_low") {
    sorted.sort((a, b) => a.rating - b.rating)
  } else if (filters.value.sortBy === "title_az") {
    sorted.sort((a, b) => a.media.title.localeCompare(b.media.title))
  } else if (filters.value.sortBy === "title_za") {
    sorted.sort((a, b) => b.media.title.localeCompare(a.media.title))
  }

  return sorted
})

// Reset Filters on Tab Change
watch(activeTab, () => {
  filters.value.search = ""
  filters.value.mediaType = "all"
  filters.value.genre = "all"
  filters.value.sortBy = "newest"
})
</script>

<style scoped>
.library-root {
  font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background-color: #0c0c0e;
  color: #f3f4f6;
  min-height: 100vh;
  padding: 2.5rem 2rem;
  position: relative;
  overflow: hidden;
}

.grain {
  position: fixed;
  inset: 0;
  background-image: radial-gradient(rgba(255, 255, 255, 0.015) 1px, transparent 0);
  background-size: 24px 24px;
  pointer-events: none;
  z-index: 1;
}

/* Loading styling */
.library-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 70vh;
  gap: 1rem;
}

.loading-ring {
  width: 48px;
  height: 48px;
  border: 3px solid rgba(225, 37, 27, 0.1);
  border-top-color: #e1251b;
  border-radius: 50%;
  animation: spin 1s infinite linear;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: 0.95rem;
  color: #9ca3af;
}

.library-layout {
  position: relative;
  z-index: 2;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Back Link */
.back-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: #9ca3af;
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
  transition: color 0.2s;
  align-self: flex-start;
}

.back-link:hover {
  color: #e1251b;
}

/* Summary Card Header */
.profile-summary-card {
  background: rgba(25, 25, 28, 0.65);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

@media (min-width: 768px) {
  .profile-summary-card {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
  }
}

.user-info-section {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.avatar-wrapper {
  position: relative;
}

.avatar-ring {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e1251b, #8a1612);
  padding: 2px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.avatar-fallback {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: #27272a;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
}

.level-chip {
  position: absolute;
  bottom: -4px;
  right: -4px;
  background: #e1251b;
  color: #ffffff;
  font-size: 0.675rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 8px;
  border: 2px solid #161618;
}

.user-meta {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
  letter-spacing: -0.01em;
}

.user-tag {
  font-size: 0.875rem;
  color: #9ca3af;
  margin: 0.125rem 0 0.5rem 0;
}

.user-bio {
  font-size: 0.875rem;
  color: #d1d5db;
  margin: 0;
  max-width: 500px;
  line-height: 1.5;
}

.user-counts-section {
  display: flex;
  gap: 2rem;
}

.count-box {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.count-val {
  font-size: 1.5rem;
  font-weight: 800;
  color: #ffffff;
}

.count-lbl {
  font-size: 0.75rem;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-top: 0.25rem;
}

/* Tabs */
.library-tabs {
  display: flex;
  gap: 0.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  padding-bottom: 0.5rem;
  overflow-x: auto;
}

.tab-btn {
  background: transparent;
  border: none;
  color: #9ca3af;
  padding: 0.75rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-radius: 10px;
  transition: background-color 0.2s, color 0.2s;
  white-space: nowrap;
}

.tab-btn:hover {
  background: rgba(255, 255, 255, 0.03);
  color: #ffffff;
}

.tab-btn--active {
  background: rgba(225, 37, 27, 0.1);
  color: #e1251b;
}

.tab-count {
  font-size: 0.75rem;
  background: rgba(255, 255, 255, 0.08);
  color: #d1d5db;
  padding: 2px 6px;
  border-radius: 6px;
  margin-left: 0.25rem;
}

.tab-btn--active .tab-count {
  background: rgba(225, 37, 27, 0.2);
  color: #e1251b;
}

/* Dashboard styling */
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

/* Filters styling */
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  background: rgba(20, 20, 22, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.03);
  border-radius: 16px;
  padding: 1rem;
  margin-bottom: 2rem;
  align-items: center;
}

.filter-input-wrapper {
  display: flex;
  align-items: center;
  background: #161618;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 0 1rem;
  height: 42px;
}

.search-wrapper {
  flex: 1;
  min-width: 200px;
}

.filter-icon, .select-icon {
  color: #9ca3af;
}

.filter-input {
  background: transparent;
  border: none;
  color: #ffffff;
  padding-left: 0.75rem;
  font-size: 0.875rem;
  width: 100%;
  outline: none;
}

.filter-select-wrapper {
  display: flex;
  align-items: center;
  background: #161618;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 0 1rem;
  height: 42px;
  min-width: 150px;
  position: relative;
}

.filter-select {
  background: transparent;
  border: none;
  color: #ffffff;
  padding-left: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  width: 100%;
  outline: none;
  cursor: pointer;
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
}

.filter-select option {
  background: #161618;
  color: #ffffff;
}

/* Grid lists and cards */
.media-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

@media (min-width: 640px) {
  .media-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 768px) {
  .media-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (min-width: 1024px) {
  .media-grid {
    grid-template-columns: repeat(5, 1fr);
  }
}

.media-library-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 16px;
  overflow: hidden;
  transition: transform 0.25s, border-color 0.25s;
  display: flex;
  flex-direction: column;
}

.media-library-card:hover {
  transform: translateY(-4px);
  border-color: rgba(255, 255, 255, 0.12);
}

.media-poster-frame {
  position: relative;
  aspect-ratio: 2/3;
  background: #1c1c1e;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.media-poster-frame img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.media-library-card:hover .media-poster-frame img {
  transform: scale(1.05);
}

.media-type-chip {
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(0, 0, 0, 0.75);
  color: #ffffff;
  font-size: 0.625rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.media-card-meta {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1;
}

.media-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: #ffffff;
  text-decoration: none;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.media-title:hover {
  color: #e1251b;
}

.media-details-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.75rem;
}

.media-rating-badge {
  color: #ffb800;
  font-weight: 600;
}

.save-date {
  color: #71717a;
}

.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.tag-badge {
  font-size: 0.675rem;
  background: rgba(255, 255, 255, 0.05);
  color: #a1a1aa;
  padding: 2px 6px;
  border-radius: 4px;
}

/* Reviews List card styling */
.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.review-library-card {
  background: rgba(20, 20, 22, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 20px;
  padding: 1.5rem;
  display: flex;
  gap: 1.5rem;
}

.review-media-poster {
  width: 80px;
  height: 120px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
  background: #1c1c1e;
  display: flex;
  align-items: center;
  justify-content: center;
}

.review-media-poster img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.review-card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.review-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.review-media-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: #ffffff;
  text-decoration: none;
  transition: color 0.2s;
}

.review-media-title:hover {
  color: #e1251b;
}

.review-card-stars {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  background: rgba(255, 184, 0, 0.08);
  border: 1px solid rgba(255, 184, 0, 0.15);
  padding: 4px 10px;
  border-radius: 8px;
}

.rating-num {
  font-size: 0.8125rem;
  font-weight: 700;
  color: #ffb800;
}

.stars-icon {
  color: #ffb800;
  font-size: 0.8125rem;
}

.review-card-body {
  font-size: 0.875rem;
  color: #d1d5db;
  line-height: 1.6;
  margin: 0;
  white-space: pre-wrap;
}

.review-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 0.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.review-date {
  font-size: 0.75rem;
  color: #71717a;
}

.review-meta-actions {
  display: flex;
  gap: 1rem;
}

.action-count {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.75rem;
  color: #a1a1aa;
}

.liked-heart {
  color: #ef4444;
}

/* Empty state styling */
.empty-state-card {
  grid-column: 1 / -1;
  text-align: center;
  padding: 4rem 2rem;
  background: rgba(20, 20, 22, 0.3);
  border: 1px dashed rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  color: #71717a;
}

.empty-icon {
  margin-bottom: 1rem;
  color: #3f3f46;
}

/* Poster fallback */
.poster-fallback {
  color: #3f3f46;
}

/* entrance animation */
.fade-in-up {
  opacity: 0;
  transform: translateY(15px);
  animation: fadeInUp 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
