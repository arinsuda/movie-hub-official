import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { feedApi } from "@/api/endpoints/feed";
import type {
  FeedItemResponse,
  FeedPaginationMeta,
  ActivitySettingsResponse,
  UpdateActivitySettingsRequest,
} from "@/types/feed";

const DEFAULT_PAGINATION: FeedPaginationMeta = {
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0,
};

export const useFeedStore = defineStore("feed", () => {
  const items = ref<FeedItemResponse[]>([]);
  const pagination = ref<FeedPaginationMeta>({ ...DEFAULT_PAGINATION });
  const loading = ref(false);
  const loadingMore = ref(false);
  const error = ref<string | null>(null);

  const settings = ref<ActivitySettingsResponse | null>(null);
  const settingsLoading = ref(false);
  const settingsError = ref<string | null>(null);

  const hasMore = computed(
    () => pagination.value.page < pagination.value.total_pages
  );
  const isEmpty = computed(() => !loading.value && items.value.length === 0);

  async function fetchFeed(page = 1) {
    if (loading.value || loadingMore.value) return;

    const limit = pagination.value?.limit ?? DEFAULT_PAGINATION.limit;
    if (page === 1) loading.value = true;
    else loadingMore.value = true;
    error.value = null;

    try {
      const res = await feedApi.getFeed({ page, limit });
      const next = res.data.items ?? [];
      items.value = page === 1 ? next : [...items.value, ...next];
      pagination.value = {
        ...DEFAULT_PAGINATION,
        ...res.data.pagination,
      };
    } catch (err) {
      console.error("fetchFeed failed:", err);
      error.value = "โหลดฟีดไม่สำเร็จ";
    } finally {
      loading.value = false;
      loadingMore.value = false;
    }
  }

  async function loadMore() {
    if (!hasMore.value || loadingMore.value || loading.value) return;
    await fetchFeed(pagination.value.page + 1);
  }

  async function refresh() {
    await fetchFeed(1);
  }

  async function updateVisibility(
    id: number,
    visibility: "default" | "public" | "followers" | "private"
  ) {
    const item = items.value.find((i) => i.id === id);
    if (!item) return;

    const prevVisibility = item.visibility;
    item.visibility = visibility;

    try {
      await feedApi.updateVisibility(id, visibility);
    } catch (err) {
      item.visibility = prevVisibility;
      console.error("updateVisibility failed:", err);
      throw err;
    }
  }

  async function deleteActivity(id: number) {
    const idx = items.value.findIndex((i) => i.id === id);
    if (idx === -1) return;

    const [removed] = items.value.splice(idx, 1);
    if (!removed) return;
    pagination.value.total = Math.max(0, pagination.value.total - 1);

    try {
      await feedApi.deleteActivity(id);
    } catch (err) {
      items.value.splice(idx, 0, removed);
      pagination.value.total++;
      console.error("deleteActivity failed:", err);
      throw err;
    }
  }

  async function fetchSettings() {
    if (settingsLoading.value) return;
    settingsLoading.value = true;
    settingsError.value = null;
    try {
      const res = await feedApi.getSettings();
      settings.value = res.data;
    } catch (err) {
      console.error("fetchSettings failed:", err);
      settingsError.value = "โหลดการตั้งค่าไม่สำเร็จ";
    } finally {
      settingsLoading.value = false;
    }
  }

  async function updateSettings(patch: UpdateActivitySettingsRequest) {
    const prev = settings.value;
    if (settings.value) {
      settings.value = {
        ...settings.value,
        ...patch,
      } as ActivitySettingsResponse;
    }

    try {
      const res = await feedApi.updateSettings(patch);
      settings.value = res.data;
      return true;
    } catch (err) {
      settings.value = prev;
      console.error("updateSettings failed:", err);
      return false;
    }
  }

  function reset() {
    items.value = [];
    pagination.value = { ...DEFAULT_PAGINATION };
    loading.value = false;
    loadingMore.value = false;
    error.value = null;
  }

  return {
    items,
    pagination,
    loading,
    loadingMore,
    error,
    hasMore,
    isEmpty,

    settings,
    settingsLoading,
    settingsError,

    fetchFeed,
    loadMore,
    refresh,
    updateVisibility,
    deleteActivity,
    fetchSettings,
    updateSettings,
    reset,
  };
});
