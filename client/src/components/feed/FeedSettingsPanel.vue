<template>
  <div
    class="w-72 rounded-xl border border-white/10 bg-neutral-900 p-4 shadow-xl"
  >
    <h3 class="mb-1 text-sm font-semibold text-white">การแสดงกิจกรรมในฟีด</h3>
    <p class="mb-3 text-xs text-white/40">
      เลือกว่ากิจกรรมประเภทไหนของคุณจะไปโผล่ในฟีดของคนที่ติดตามคุณอยู่
    </p>

    <div v-if="loading" class="py-6 text-center text-xs text-white/40">
      กำลังโหลด...
    </div>

    <p v-else-if="settingsError" class="py-6 text-center text-xs text-red-400">
      {{ settingsError }}
    </p>

    <ul v-else-if="settings" class="space-y-2.5">
      <li
        v-for="opt in options"
        :key="opt.key"
        class="flex items-center justify-between gap-3"
      >
        <span class="text-sm text-white/80">{{ opt.label }}</span>
        <button
          type="button"
          role="switch"
          :aria-checked="settings[opt.key]"
          class="relative h-5 w-9 shrink-0 rounded-full transition-colors"
          :class="settings[opt.key] ? 'bg-emerald-500' : 'bg-white/15'"
          @click="toggle(opt.key)"
        >
          <span
            class="absolute top-0.5 h-4 w-4 rounded-full bg-white transition-transform"
            :class="settings[opt.key] ? 'translate-x-4' : 'translate-x-0.5'"
          />
        </button>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useFeedStore } from "@/stores/feed";
import type { ActivitySettingsResponse } from "@/types/feed";

const store = useFeedStore();
const {
  settings,
  settingsLoading: loading,
  settingsError,
} = storeToRefs(store);

const options: { key: keyof ActivitySettingsResponse; label: string }[] = [
  { key: "review_created", label: "สร้างรีวิวใหม่" },
  { key: "review_commented", label: "แสดงความคิดเห็นในรีวิว" },
  { key: "review_liked", label: "ถูกใจรีวิว" },
  { key: "media_liked", label: "ถูกใจหนัง/ซีรีส์" },
  { key: "watchlist_added", label: "เพิ่มลงวอทช์ลิสต์" },
  { key: "watched_added", label: "ทำเครื่องหมายว่าดูแล้ว" },
  { key: "achievement_unlocked", label: "ปลดล็อกความสำเร็จ" },
];

function toggle(key: keyof ActivitySettingsResponse) {
  if (!settings.value) return;
  store.updateSettings({ [key]: !settings.value[key] });
}

onMounted(() => {
  if (!settings.value) store.fetchSettings();
});
</script>
