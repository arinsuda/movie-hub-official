import { ref, computed } from "vue"
import { defineStore } from "pinia"
import { achievementApi } from "@/api/endpoints/achievement"
import type {
  Achievement,
  AchievementFilterState,
  PaginationMeta,
  UserAchievement,
} from "@/types/achievement"

export const useAchievementStore = defineStore("achievement", () => {
  // ── Catalog: achievement ทั้งหมดในระบบ ใช้สร้าง filter action_type ──
  const catalog = ref<Achievement[]>([])
  const catalogLoaded = ref(false)

  const actionTypes = computed(() => {
    const set = new Set(catalog.value.map(a => a.action_type))
    return Array.from(set).sort()
  })

  async function fetchCatalog() {
    if (catalogLoaded.value) return
    const res = await achievementApi.listAll({ page: 1, limit: 100 })
    catalog.value = res.data
    catalogLoaded.value = true
  }

  // ── User progress ──
  const userAchievements = ref<UserAchievement[]>([])
  const pagination = ref<PaginationMeta>({
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0,
  })
  const loading = ref(false)
  const error = ref<string | null>(null)

  const filters = ref<AchievementFilterState>({
    unlockedFilter: "all",
    actionType: "",
    search: "",
    page: 1,
    limit: 12,
  })

  const unlockedCount = computed(
    () => userAchievements.value.filter(ua => ua.is_unlocked).length,
  )

  const showingRange = computed(() => {
    const { page, limit, total } = pagination.value
    if (total === 0) return { from: 0, to: 0 }
    const from = (page - 1) * limit + 1
    const to = Math.min(page * limit, total)
    return { from, to }
  })

  /** ใช้กับหน้าเต็ม /users/:userId/achievements (ทั้ง locked/unlocked + filter) */
  async function fetchUserAchievements(userId: number) {
    loading.value = true
    error.value = null
    try {
      const params: Record<string, unknown> = {
        page: filters.value.page,
        limit: filters.value.limit,
      }
      if (filters.value.unlockedFilter !== "all") {
        params.unlocked = filters.value.unlockedFilter === "unlocked"
      }
      if (filters.value.actionType) {
        params.action_type = filters.value.actionType
      }

      const res = await achievementApi.listByUser(userId, params)
      userAchievements.value = res.data
      pagination.value = res.pagination
    } catch (err) {
      console.error("fetchUserAchievements failed:", err)
      error.value = "ไม่สามารถโหลด Achievement ได้"
    } finally {
      loading.value = false
    }
  }

  /** ใช้กับ widget UserAchievements.vue — เฉพาะที่ unlocked แล้วเท่านั้น */
  async function fetchUnlockedOnly(userId: number, limit = 20) {
    loading.value = true
    error.value = null
    try {
      const res = await achievementApi.listByUser(userId, {
        unlocked: true,
        page: 1,
        limit,
      })
      userAchievements.value = res.data
      pagination.value = res.pagination
    } catch (err) {
      console.error("fetchUnlockedOnly failed:", err)
      error.value = "ไม่สามารถโหลด Achievement ได้"
    } finally {
      loading.value = false
    }
  }

  // filter จากคำค้นหาทำฝั่ง client (backend ยังไม่รองรับ search)
  const filteredBySearch = computed(() => {
    const q = filters.value.search.trim().toLowerCase()
    if (!q) return userAchievements.value
    return userAchievements.value.filter(
      ua =>
        ua.achievement.name.toLowerCase().includes(q) ||
        ua.achievement.description.toLowerCase().includes(q),
    )
  })

  function setPage(page: number) {
    filters.value.page = page
  }

  function resetFilters() {
    filters.value = {
      unlockedFilter: "all",
      actionType: "",
      search: "",
      page: 1,
      limit: 12,
    }
  }

  return {
    catalog,
    actionTypes,
    fetchCatalog,

    userAchievements,
    pagination,
    loading,
    error,
    filters,
    filteredBySearch,
    unlockedCount,
    showingRange,

    fetchUserAchievements,
    fetchUnlockedOnly,
    setPage,
    resetFilters,
  }
})
