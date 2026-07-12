<template>
  <div>
    <div v-if="loading && items.length === 0" class="divide-y divide-white/5">
      <FeedItemSkeleton v-for="n in 5" :key="n" />
    </div>

    <div
      v-else-if="error && items.length === 0"
      class="flex flex-col items-center gap-3 py-16 text-center"
    >
      <p class="text-sm text-white/50">{{ error }}</p>
      <button
        type="button"
        class="rounded-md bg-white/10 px-4 py-2 text-sm text-white hover:bg-white/20"
        @click="emit('retry')"
      >
        ลองอีกครั้ง
      </button>
    </div>

    <div
      v-else-if="isEmpty"
      class="flex flex-col items-center gap-2 py-20 text-center text-white/40"
    >
      <p class="text-sm">ยังไม่มีความเคลื่อนไหวในฟีด</p>
      <p class="text-xs">ติดตามเพื่อน ๆ เพิ่ม เพื่อดูกิจกรรมของพวกเขาที่นี่</p>
    </div>

    <template v-else>
      <div class="divide-y divide-white/5">
        <FeedItem
          v-for="item in items"
          :key="item.id"
          :item="item"
          :can-manage="item.actor.id === currentUserId"
          @hide="emit('hide', $event)"
        />
      </div>

      <div ref="sentinel" class="h-1" />

      <div v-if="loadingMore" class="flex justify-center py-4">
        <span class="text-xs text-white/40">กำลังโหลดเพิ่มเติม...</span>
      </div>
      <p
        v-else-if="!hasMore && items.length > 0"
        class="py-4 text-center text-xs text-white/30"
      >
        ดูครบทุกกิจกรรมแล้ว
      </p>

      <div
        v-if="error"
        class="flex flex-col items-center gap-2 py-4 text-center"
      >
        <p class="text-xs text-red-400">{{ error }}</p>
        <button
          type="button"
          class="text-xs text-white/60 underline"
          @click="emit('retry')"
        >
          ลองอีกครั้ง
        </button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import type { FeedItemResponse } from "@/types/feed";
import { useInfiniteScroll } from "@/composables/useInfiniteScroll";
import FeedItem from "./FeedItem.vue";
import FeedItemSkeleton from "./FeedItemSkeleton.vue";

const props = defineProps<{
  items: FeedItemResponse[];
  loading: boolean;
  loadingMore: boolean;
  error: string | null;
  hasMore: boolean;
  isEmpty: boolean;
  currentUserId?: number;
}>();

const emit = defineEmits<{
  "load-more": [];
  retry: [];
  hide: [id: number];
}>();

const sentinel = ref<HTMLElement | null>(null);

useInfiniteScroll(sentinel, () => emit("load-more"), {
  disabled: computed(
    () => props.loading || props.loadingMore || !props.hasMore,
  ),
});
</script>
