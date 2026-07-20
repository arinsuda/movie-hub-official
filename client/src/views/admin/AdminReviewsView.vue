<template>
  <div class="admin-reviews-view">
    <div class="view-header">
      <h1 class="page-title">{{ $t("admin.reviews.title") }}</h1>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="search-input-wrapper">
        <Search :size="16" class="search-icon" />
        <input
          v-model="search"
          type="text"
          class="search-input"
          :placeholder="$t('admin.reviews.searchPlaceholder')"
          @input="onSearchInput"
        />
      </div>

      <div class="filter-controls">
        <select v-model="mediaTypeFilter" class="filter-select" @change="fetchData(1)">
          <option value="all">{{ $t("admin.reviews.typeAll") }}</option>
          <option value="movie">{{ $t("admin.reviews.typeMovie") }}</option>
          <option value="tv">{{ $t("admin.reviews.typeTv") }}</option>
        </select>

        <select v-model="visibilityFilter" class="filter-select" @change="fetchData(1)">
          <option value="all">{{ $t("admin.reviews.visibilityAll") }}</option>
          <option value="public">{{ $t("admin.reviews.visibilityPublic") }}</option>
          <option value="private">{{ $t("admin.reviews.visibilityPrivate") }}</option>
        </select>

        <select v-model="sortBy" class="filter-select" @change="fetchData(1)">
          <option value="created_at">Sort by Date</option>
          <option value="rating">Sort by Rating</option>
          <option value="like_count">Sort by Likes</option>
        </select>

        <button class="sort-order-btn" @click="toggleSortOrder">
          <ArrowUpDown :size="16" />
          <span>{{ sortOrder.toUpperCase() }}</span>
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="table-container">
      <table class="admin-table">
        <thead>
          <tr>
            <th>{{ $t("admin.reviews.colId") }}</th>
            <th>{{ $t("admin.reviews.colAuthor") }}</th>
            <th>{{ $t("admin.reviews.colType") }}</th>
            <th>{{ $t("admin.reviews.colRating") }}</th>
            <th>{{ $t("admin.reviews.colBody") }}</th>
            <th>{{ $t("admin.reviews.colLikes") }}</th>
            <th>{{ $t("admin.reviews.colCreated") }}</th>
            <th class="text-right">{{ $t("admin.reviews.colActions") }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="adminStore.isLoadingReviews">
            <td colspan="8" class="text-center py-8">
              <div class="spinner inline" />
            </td>
          </tr>
          <tr v-else-if="reviews.length === 0">
            <td colspan="8" class="text-center py-8 text-muted">
              No reviews found.
            </td>
          </tr>
          <tr v-for="r in reviews" :key="r.id">
            <td>#{{ r.id }}</td>
            <td class="font-semibold">@{{ r.username }}</td>
            <td>
              <span class="media-type-tag" :class="r.media_type">
                {{ r.media_type === "movie" ? "Movie" : "TV" }} #{{ r.media_id }}
              </span>
            </td>
            <td>
              <div class="rating-cell">
                <Star :size="14" fill="#f5c518" color="#f5c518" />
                <span>{{ r.rating.toFixed(1) }}</span>
              </div>
            </td>
            <td class="review-body-cell" :title="r.body">
              {{ r.body }}
            </td>
            <td>{{ r.like_count }}</td>
            <td class="text-secondary">{{ formatDate(r.created_at) }}</td>
            <td class="text-right">
              <button
                class="action-btn danger"
                title="Delete Review"
                @click="handleDelete(r)"
              >
                <Trash2 :size="14" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination Bar -->
    <div v-if="adminStore.reviewsTotalPages > 1" class="pagination-bar">
      <span class="pagination-info">
        {{ $t("admin.pagination.showing", {
          from: (adminStore.reviewsPage - 1) * adminStore.reviewsLimit + 1,
          to: Math.min(adminStore.reviewsPage * adminStore.reviewsLimit, adminStore.reviewsTotal),
          total: adminStore.reviewsTotal
        }) }}
      </span>

      <div class="pagination-buttons">
        <button
          class="page-btn"
          :disabled="adminStore.reviewsPage <= 1"
          @click="fetchData(adminStore.reviewsPage - 1)"
        >
          {{ $t("admin.pagination.previous") }}
        </button>
        <span class="page-num">{{ adminStore.reviewsPage }} / {{ adminStore.reviewsTotalPages }}</span>
        <button
          class="page-btn"
          :disabled="adminStore.reviewsPage >= adminStore.reviewsTotalPages"
          @click="fetchData(adminStore.reviewsPage + 1)"
        >
          {{ $t("admin.pagination.next") }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useAdminStore } from "@/stores/admin"
import type { AdminReviewRow } from "@/types"
import { Search, ArrowUpDown, Star, Trash2 } from "lucide-vue-next"

const adminStore = useAdminStore()
const reviews = computed(() => adminStore.reviews)

const search = ref("")
const mediaTypeFilter = ref("all")
const visibilityFilter = ref("all")
const sortBy = ref("created_at")
const sortOrder = ref<"asc" | "desc">("desc")

let debounceTimer: ReturnType<typeof setTimeout> | null = null

function onSearchInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => fetchData(1), 300)
}

function toggleSortOrder() {
  sortOrder.value = sortOrder.value === "asc" ? "desc" : "asc"
  fetchData(1)
}

function fetchData(page = 1) {
  adminStore.fetchReviews({
    page,
    limit: 20,
    search: search.value.trim(),
    media_type: mediaTypeFilter.value,
    visibility: visibilityFilter.value,
    sort_by: sortBy.value,
    sort_order: sortOrder.value,
  })
}

async function handleDelete(review: AdminReviewRow) {
  if (!confirm(`Are you sure you want to delete review #${review.id} by @${review.username}?`)) return
  const reason = prompt("Enter optional reason for review deletion:") || undefined

  try {
    await adminStore.deleteReview(review.id, reason)
    fetchData(adminStore.reviewsPage)
  } catch (err: unknown) {
    const e = err as { response?: { data?: { error?: string } } }
    alert(e.response?.data?.error || "Failed to delete review")
  }
}

function formatDate(isoStr: string) {
  if (!isoStr) return "-"
  return new Date(isoStr).toLocaleDateString()
}

onMounted(() => {
  fetchData(1)
})
</script>

<style scoped>
.admin-reviews-view {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.page-title {
  font-size: 1.6rem;
  font-weight: 800;
  margin: 0;
  color: var(--color-text-primary);
}

.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: space-between;
  align-items: center;
}

.search-input-wrapper {
  position: relative;
  flex: 1;
  min-width: 260px;
}

.search-icon {
  position: absolute;
  left: 0.85rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-secondary);
}

.search-input {
  width: 100%;
  padding: 0.6rem 0.85rem 0.6rem 2.4rem;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  color: var(--color-text-primary);
  outline: none;
  font-size: 0.875rem;
  box-sizing: border-box;
}

.filter-controls {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.filter-select {
  padding: 0.6rem 0.85rem;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  color: var(--color-text-primary);
  font-size: 0.85rem;
  outline: none;
  cursor: pointer;
}

.sort-order-btn {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.6rem 0.85rem;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  color: var(--color-text-primary);
  font-size: 0.85rem;
  cursor: pointer;
}

.table-container {
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-card);
  overflow-x: auto;
}

.admin-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
}

.admin-table th {
  padding: 0.85rem 1rem;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  font-weight: 600;
}

.admin-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  color: var(--color-text-primary);
}

.font-semibold { font-weight: 600; }

.media-type-tag {
  font-size: 0.775rem;
  font-weight: 700;
  padding: 0.2rem 0.5rem;
  border-radius: 0.35rem;
  background: var(--color-surface-3);
  color: var(--color-text-secondary);
}

.rating-cell {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  font-weight: 700;
  color: #f5c518;
}

.review-body-cell {
  max-width: 360px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.action-btn {
  padding: 0.4rem;
  background: var(--color-surface-3);
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  border-radius: 0.35rem;
  cursor: pointer;
  transition: all 0.15s;
}
.action-btn.danger:hover { color: #ef4444; border-color: #ef4444; }

.text-right { text-align: right; }
.text-center { text-align: center; }
.text-secondary { color: var(--color-text-secondary); }
.text-muted { color: var(--color-text-muted); }
.py-8 { padding-top: 2rem; padding-bottom: 2rem; }

.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: var(--color-text-secondary);
}

.pagination-buttons {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.page-btn {
  padding: 0.4rem 0.85rem;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
  border-radius: 0.35rem;
  cursor: pointer;
}
.page-btn:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
