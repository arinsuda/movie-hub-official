import { computed, ref, type Ref } from "vue";
import { feedApi } from "@/api/endpoints/feed";
import type { FeedItemResponse, FeedPaginationMeta } from "@/types/feed";

const DEFAULT_PAGINATION: FeedPaginationMeta = {
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0,
};

// ใช้แสดง activity ของ user คนเดียว เช่น แท็บ "กิจกรรม" ในหน้าโปรไฟล์
// แยกจาก useFeedStore เพราะอาจเปิดดูโปรไฟล์หลายคนพร้อมกัน ไม่อยาก share state กับฟีดหลัก
export function useUserActivities(userId: Ref<number> | number) {
  const items = ref<FeedItemResponse[]>([]);
  const pagination = ref<FeedPaginationMeta>({ ...DEFAULT_PAGINATION });
  const loading = ref(false);
  const loadingMore = ref(false);
  const error = ref<string | null>(null);

  const hasMore = computed(
    () => pagination.value.page < pagination.value.total_pages,
  );
  const isEmpty = computed(() => !loading.value && items.value.length === 0);

  const resolveUserId = () =>
    typeof userId === "number" ? userId : userId.value;

  async function fetch(page = 1) {
    if (loading.value || loadingMore.value) return;
    if (page === 1) loading.value = true;
    else loadingMore.value = true;
    error.value = null;

    try {
      const res = await feedApi.getUserActivities(resolveUserId(), {
        page,
        limit: pagination.value.limit,
      });
      const next = res.data.items ?? [];
      items.value = page === 1 ? next : [...items.value, ...next];
      pagination.value = {
        ...DEFAULT_PAGINATION,
        ...(res.data.pagination ?? {}),
      };
    } catch (err) {
      console.error("useUserActivities fetch failed:", err);
      error.value = "โหลดกิจกรรมไม่สำเร็จ";
    } finally {
      loading.value = false;
      loadingMore.value = false;
    }
  }

  async function loadMore() {
    if (!hasMore.value || loadingMore.value) return;
    await fetch(pagination.value.page + 1);
  }

  function reset() {
    items.value = [];
    pagination.value = { ...DEFAULT_PAGINATION };
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
    fetch,
    loadMore,
    reset,
  };
}
