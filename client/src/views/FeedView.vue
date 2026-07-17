<template>
  <div class="mx-auto max-w-2xl px-4 py-6 text-white">
    <header class="mb-4 flex items-center justify-between">
      <h1 class="text-xl font-semibold">ฟีดกิจกรรม</h1>

      <div class="flex items-center gap-1">
        <button
          type="button"
          class="rounded-md p-2 text-white/50 transition-colors hover:bg-white/10 hover:text-white"
          title="รีเฟรชฟีด"
          :disabled="feedStore.loading"
          @click="feedStore.refresh"
        >
          <svg
            viewBox="0 0 20 20"
            class="h-5 w-5"
            :class="{ 'animate-spin': feedStore.loading }"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path
              d="M16 10a6 6 0 10-1.76 4.24M16 10V6m0 4h-4"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>

        <div ref="settingsRef" class="relative">
          <button
            type="button"
            class="rounded-md p-2 text-white/50 transition-colors hover:bg-white/10 hover:text-white"
            title="ตั้งค่าฟีด"
            @click="isSettingsOpen = !isSettingsOpen"
          >
            <svg
              viewBox="0 0 20 20"
              class="h-5 w-5"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <path d="M10 12.5a2.5 2.5 0 100-5 2.5 2.5 0 000 5z" />
              <path
                d="M15.5 10a1.4 1.4 0 00.28 1.54l.05.05a1.7 1.7 0 11-2.4 2.4l-.05-.05a1.4 1.4 0 00-1.54-.28 1.4 1.4 0 00-.85 1.28v.14a1.7 1.7 0 11-3.4 0v-.07a1.4 1.4 0 00-.92-1.28 1.4 1.4 0 00-1.54.28l-.05.05a1.7 1.7 0 11-2.4-2.4l.05-.05a1.4 1.4 0 00.28-1.54 1.4 1.4 0 00-1.28-.85H4.5a1.7 1.7 0 110-3.4h.07a1.4 1.4 0 001.28-.92 1.4 1.4 0 00-.28-1.54l-.05-.05a1.7 1.7 0 112.4-2.4l.05.05a1.4 1.4 0 001.54.28h.07a1.4 1.4 0 00.85-1.28V4.5a1.7 1.7 0 113.4 0v.07a1.4 1.4 0 00.85 1.28 1.4 1.4 0 001.54-.28l.05-.05a1.7 1.7 0 112.4 2.4l-.05.05a1.4 1.4 0 00-.28 1.54v.07a1.4 1.4 0 001.28.85h.14a1.7 1.7 0 110 3.4h-.14a1.4 1.4 0 00-1.28.85z"
              />
            </svg>
          </button>

          <Transition name="fade">
            <div
              v-if="isSettingsOpen"
              class="absolute right-0 top-full z-20 mt-2"
            >
              <FeedSettingsPanel />
            </div>
          </Transition>
        </div>
      </div>
    </header>

    <!-- Floating New Feed Items Alert -->
    <Transition name="slide-down">
      <div
        v-if="feedStore.newItemsCount > 0"
        class="sticky top-4 z-10 mb-4 flex justify-center"
      >
        <button
          type="button"
          class="flex items-center gap-1.5 rounded-full border border-blue-500/30 bg-blue-500/10 px-4 py-2 text-sm font-semibold text-blue-400 backdrop-blur-md transition-all hover:bg-blue-500/20 active:scale-95 shadow-lg shadow-blue-500/5"
          @click="feedStore.showNewItems()"
        >
          <svg
            viewBox="0 0 20 20"
            class="h-4 w-4 animate-bounce"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              d="M10 3v14m0-14l-4 4m4-4l4 4"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          มี {{ feedStore.newItemsCount }} กิจกรรมใหม่
        </button>
      </div>
    </Transition>

    <FeedList
      :items="feedStore.items"
      :loading="feedStore.loading"
      :loading-more="feedStore.loadingMore"
      :error="feedStore.error"
      :has-more="feedStore.hasMore"
      :is-empty="feedStore.isEmpty"
      :current-user-id="currentUserId"
      @load-more="feedStore.loadMore"
      @retry="feedStore.refresh"
      @update-visibility="feedStore.updateVisibility"
      @delete-activity="feedStore.deleteActivity"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { useFeedStore } from "@/stores/feed";
import { useClickOutside } from "@/composables/useClickOutside";
import FeedList from "@/components/feed/FeedList.vue";
import FeedSettingsPanel from "@/components/feed/FeedSettingsPanel.vue";
import { useAuthStore } from "@/stores/auth";

const feedStore = useFeedStore();
const authStore = useAuthStore();
const currentUserId = computed(() => authStore.user?.id);

const isSettingsOpen = ref(false);
const settingsRef = ref<HTMLElement | null>(null);
useClickOutside(settingsRef, () => (isSettingsOpen.value = false));

onMounted(() => {
  feedStore.bindSocket();
  if (feedStore.items.length === 0) feedStore.fetchFeed(1);
});

onUnmounted(() => {
  feedStore.unbindSocket();
});
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition:
    transform 0.25s cubic-bezier(0.16, 1, 0.3, 1),
    opacity 0.2s ease;
}
.slide-down-enter-from,
.slide-down-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}
</style>

