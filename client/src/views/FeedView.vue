<template>
  <div class="mx-auto max-w-2xl px-4 py-8 text-white">
    <!-- Header Section -->
    <header class="mb-8 flex items-end justify-between">
      <div>
        <h1 class="text-3xl font-extrabold tracking-tight bg-gradient-to-r from-white via-white/95 to-white/60 bg-clip-text text-transparent">
          ฟีดกิจกรรม
        </h1>
        <p class="mt-1.5 text-sm text-white/50">
          ความเคลื่อนไหวล่าสุดจากเพื่อน ๆ และผู้ใช้ที่คุณติดตาม
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          class="flex h-10 w-10 items-center justify-center rounded-xl border border-white/5 bg-white/5 text-white/60 transition-all duration-300 hover:bg-white/10 hover:text-white hover:scale-105 active:scale-95 disabled:opacity-50"
          title="รีเฟรชฟีด"
          :disabled="feedStore.loading"
          @click="feedStore.refresh"
        >
          <RefreshCw
            class="h-4.5 w-4.5"
            :class="{ 'animate-spin': feedStore.loading }"
          />
        </button>

        <div ref="settingsRef" class="relative">
          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-xl border border-white/5 bg-white/5 text-white/60 transition-all duration-300 hover:bg-white/10 hover:text-white hover:scale-105 active:scale-95"
            :class="{ 'bg-white/10 text-white border-white/20': isSettingsOpen }"
            title="ตั้งค่าฟีด"
            @click="isSettingsOpen = !isSettingsOpen"
          >
            <Settings class="h-4.5 w-4.5" />
          </button>

          <Transition name="fade-slide">
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

    <!-- Category Tabs Filter -->
    <div class="mb-6 flex gap-1 border-b border-white/5 pb-px overflow-x-auto scrollbar-none">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        type="button"
        class="relative px-4 py-3 text-sm font-medium transition-all duration-300 shrink-0 cursor-pointer"
        :class="feedStore.selectedCategory === tab.value ? 'text-red-500 font-semibold' : 'text-white/40 hover:text-white/75'"
        @click="feedStore.setCategory(tab.value)"
      >
        <span>{{ tab.label }}</span>
        <Transition name="tab-line">
          <span
            v-if="feedStore.selectedCategory === tab.value"
            class="absolute bottom-0 left-2 right-2 h-0.5 bg-red-500 rounded-full"
          />
        </Transition>
      </button>
    </div>

    <!-- Floating New Feed Items Alert -->
    <Transition name="slide-down">
      <div
        v-if="feedStore.newItemsCount > 0"
        class="sticky top-6 z-10 mb-6 flex justify-center"
      >
        <button
          type="button"
          class="flex items-center gap-1.5 rounded-full border border-red-500/30 bg-red-950/80 px-5 py-2.5 text-sm font-semibold text-red-400 backdrop-blur-md transition-all duration-300 hover:bg-red-900/90 active:scale-95 shadow-lg shadow-red-500/10 cursor-pointer"
          @click="feedStore.showNewItems()"
        >
          <ArrowUp class="h-4 w-4" />
          มี {{ feedStore.newItemsCount }} กิจกรรมใหม่ในหน้านี้
        </button>
      </div>
    </Transition>

    <!-- Feed List Component -->
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
import { RefreshCw, Settings, ArrowUp } from "lucide-vue-next";

const feedStore = useFeedStore();
const authStore = useAuthStore();
const currentUserId = computed(() => authStore.user?.id);

const isSettingsOpen = ref(false);
const settingsRef = ref<HTMLElement | null>(null);
useClickOutside(settingsRef, () => (isSettingsOpen.value = false));

const tabs = [
  { label: "ทั้งหมด", value: "" },
  { label: "รีวิวภาพยนตร์", value: "reviews" },
  { label: "ลิสต์ภาพยนตร์", value: "lists" },
  { label: "โซเชียล & ความสำเร็จ", value: "social" }
];

onMounted(() => {
  feedStore.bindSocket();
  if (feedStore.items.length === 0) feedStore.fetchFeed(1);
});

onUnmounted(() => {
  feedStore.unbindSocket();
});
</script>

<style scoped>
.fade-slide-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.fade-slide-leave-active {
  transition: all 0.2s cubic-bezier(0.36, 0.07, 0.19, 0.97);
}
.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.96);
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition:
    transform 0.3s cubic-bezier(0.16, 1, 0.3, 1),
    opacity 0.2s ease;
}
.slide-down-enter-from,
.slide-down-leave-to {
  transform: translateY(-12px);
  opacity: 0;
}

.tab-line-enter-active,
.tab-line-leave-active {
  transition: all 0.22s cubic-bezier(0.4, 0, 0.2, 1);
}
.tab-line-enter-from,
.tab-line-leave-to {
  opacity: 0;
  transform: scaleX(0.1);
}
</style>

