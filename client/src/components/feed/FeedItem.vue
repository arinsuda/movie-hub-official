<template>
  <article
    class="group flex gap-3 border-b border-white/5 px-4 py-4 transition-colors hover:bg-white/[0.03] relative"
  >
    <RouterLink :to="actorLink" class="shrink-0">
      <img
        :src="item.actor.avatar_url || defaultAvatar"
        :alt="item.actor.username"
        class="h-10 w-10 rounded-full object-cover"
      />
    </RouterLink>

    <div class="min-w-0 flex-1">
      <div class="flex items-start justify-between gap-2">
        <p class="text-sm leading-snug">
          <RouterLink
            :to="actorLink"
            class="font-semibold text-white hover:underline"
          >
            {{ item.actor.display_name || item.actor.username }}
          </RouterLink>
          <span class="ml-1 text-white/70">{{ item.message }}</span>
        </p>

        <!-- Privacy & Delete Management Dropdown -->
        <div v-if="canManage" class="relative" ref="menuRef">
          <button
            type="button"
            class="shrink-0 rounded p-1 text-white/30 hover:bg-white/10 hover:text-white/70 transition-colors"
            :title="`ความเป็นส่วนตัว: ${visibilityLabel}`"
            @click="isMenuOpen = !isMenuOpen"
          >
            <!-- Globe Icon (public) -->
            <svg
              v-if="item.visibility === 'public'"
              viewBox="0 0 20 20"
              class="h-4 w-4"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <circle cx="10" cy="10" r="7.25" stroke-linecap="round" />
              <path d="M3 10h14M10 3a15.3 15.3 0 013 7 15.3 15.3 0 01-3 7 15.3 15.3 0 01-3-7 15.3 15.3 0 013-7z" stroke-linecap="round" />
            </svg>
            <!-- Users Icon (followers) -->
            <svg
              v-else-if="item.visibility === 'followers'"
              viewBox="0 0 20 20"
              class="h-4 w-4"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <path d="M6.5 13.5a3.5 3.5 0 017 0M10 10a3 3 0 100-6 3 3 0 000 6z" stroke-linecap="round" />
            </svg>
            <!-- Lock Icon (private) -->
            <svg
              v-else-if="item.visibility === 'private'"
              viewBox="0 0 20 20"
              class="h-4 w-4"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <rect x="5.5" y="8.5" width="9" height="7.5" rx="1.5" stroke-linecap="round" />
              <path d="M7.75 8.5V6a2.25 2.25 0 114.5 0v2.5" stroke-linecap="round" />
            </svg>
            <!-- Gear/Default Icon -->
            <svg
              v-else
              viewBox="0 0 20 20"
              class="h-4 w-4"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <circle cx="10" cy="10" r="7.25" stroke-dasharray="2 2" stroke-linecap="round" />
              <circle cx="10" cy="10" r="2.5" />
            </svg>
          </button>

          <!-- Dropdown Menu -->
          <Transition name="slide-up">
            <div
              v-if="isMenuOpen"
              class="absolute right-0 top-full z-30 mt-1 w-48 rounded-lg border border-white/10 bg-neutral-900 py-1 shadow-2xl"
            >
              <div
                class="px-3 py-1.5 border-b border-white/5 text-[10px] font-semibold uppercase tracking-wider text-white/40"
              >
                ความเป็นส่วนตัว
              </div>
              <button
                v-for="opt in visibilityOptions"
                :key="opt.value"
                type="button"
                class="flex w-full items-center gap-2 px-3 py-2 text-left text-xs transition-colors hover:bg-white/5"
                :class="
                  item.visibility === opt.value
                    ? 'text-emerald-400 font-medium'
                    : 'text-white/70'
                "
                @click="changeVisibility(opt.value)"
              >
                <!-- Render Icon inside option -->
                <span v-html="opt.iconSvg" class="h-3.5 w-3.5 shrink-0" />
                <span>{{ opt.label }}</span>
              </button>

              <div class="border-t border-white/5 my-1" />

              <button
                type="button"
                class="flex w-full items-center gap-2 px-3 py-2 text-left text-xs text-rose-400 transition-colors hover:bg-rose-500/10"
                @click="confirmDelete"
              >
                <svg
                  viewBox="0 0 20 20"
                  class="h-3.5 w-3.5 shrink-0"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <path
                    d="M4 6h12M7 6v10a1 1 0 001 1h4a1 1 0 001-1V6M9 6V4.5a1 1 0 011-1h0a1 1 0 011 1V6"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
                <span>ลบกิจกรรม</span>
              </button>
            </div>
          </Transition>
        </div>
      </div>

      <div class="mt-1 flex items-center gap-1.5 text-xs text-white/40">
        <span>{{ activityMeta.icon }}</span>
        <span>{{ activityMeta.label }}</span>
        <span aria-hidden="true">·</span>
        <time :datetime="item.created_at">{{ relativeTime }}</time>
        <span v-if="canManage" class="text-white/20" aria-hidden="true">·</span>
        <span
          v-if="canManage"
          class="text-[9px] uppercase tracking-wider px-1 py-0.5 rounded bg-white/5 text-white/40"
        >
          {{ visibilityLabel }}
        </span>
      </div>

      <!-- Media Reference Card -->
      <RouterLink
        v-if="item.media"
        :to="mediaLink"
        class="mt-3 flex items-center gap-3 rounded-lg border border-white/10 bg-white/[0.02] p-2 transition-colors hover:border-white/20"
      >
        <img
          v-if="item.media.poster_url"
          :src="item.media.poster_url"
          :alt="item.media.title"
          class="h-16 w-11 shrink-0 rounded object-cover"
        />
        <div class="min-w-0">
          <p class="truncate text-sm font-medium text-white">
            {{ item.media.title }}
          </p>
          <p v-if="item.media.vote_average" class="text-xs text-white/40">
            ★ {{ item.media.vote_average.toFixed(1) }}
          </p>
        </div>
      </RouterLink>

      <!-- Target User Reference Card (For follow activities) -->
      <RouterLink
        v-if="item.type === 'user_followed' && item.target_user"
        :to="`/users/${item.target_user.id}`"
        class="mt-3 flex items-center gap-3 rounded-lg border border-white/10 bg-white/[0.02] p-2 transition-colors hover:border-white/20"
      >
        <img
          :src="item.target_user.avatar_url || defaultAvatar"
          :alt="item.target_user.username"
          class="h-10 w-10 shrink-0 rounded-full object-cover"
        />
        <div class="min-w-0">
          <p class="truncate text-sm font-medium text-white">
            {{ item.target_user.display_name || item.target_user.username }}
          </p>
          <p class="text-xs text-white/40">
            @{{ item.target_user.username }}
          </p>
        </div>
      </RouterLink>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import type { FeedItemResponse } from "@/types/feed";
import { getActivityMeta, mediaRoutePrefix } from "@/utils/feedActivity";
import { formatRelativeTime } from "@/utils/relativeTime";
import { useClickOutside } from "@/composables/useClickOutside";

const props = defineProps<{
  item: FeedItemResponse;
  canManage?: boolean;
}>();

const emit = defineEmits<{
  "update-visibility": [
    id: number,
    visibility: "default" | "public" | "followers" | "private",
  ];
  "delete-activity": [id: number];
}>();

const defaultAvatar = "/images/default-avatar.png";
const isMenuOpen = ref(false);
const menuRef = ref<HTMLElement | null>(null);

useClickOutside(menuRef, () => (isMenuOpen.value = false));

const activityMeta = computed(() => getActivityMeta(props.item.type));
const relativeTime = computed(() => formatRelativeTime(props.item.created_at));
const actorLink = computed(() => `/users/${props.item.actor.id}`);

const mediaLink = computed(() => {
  if (!props.item.media) return "/";
  return `/${mediaRoutePrefix(props.item.media.media_type)}/${props.item.media.id}`;
});

const visibilityLabel = computed(() => {
  switch (props.item.visibility) {
    case "public":
      return "สาธารณะ";
    case "followers":
      return "ผู้ติดตามเท่านั้น";
    case "private":
      return "เฉพาะฉัน";
    case "default":
    default:
      return "ค่าเริ่มต้น";
  }
});

const visibilityOptions = [
  {
    value: "default" as const,
    label: "ค่าเริ่มต้น (ตามประเภทบัญชี)",
    iconSvg: `<svg viewBox="0 0 20 20" class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="10" cy="10" r="7.25" stroke-dasharray="2 2" stroke-linecap="round"/><circle cx="10" cy="10" r="2.5"/></svg>`,
  },
  {
    value: "public" as const,
    label: "สาธารณะ",
    iconSvg: `<svg viewBox="0 0 20 20" class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="10" cy="10" r="7.25" stroke-linecap="round"/><path d="M3 10h14M10 3a15.3 15.3 0 013 7 15.3 15.3 0 01-3 7 15.3 15.3 0 01-3-7 15.3 15.3 0 013-7z" stroke-linecap="round"/></svg>`,
  },
  {
    value: "followers" as const,
    label: "ผู้ติดตามเท่านั้น",
    iconSvg: `<svg viewBox="0 0 20 20" class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6.5 13.5a3.5 3.5 0 017 0M10 10a3 3 0 100-6 3 3 0 000 6z" stroke-linecap="round"/></svg>`,
  },
  {
    value: "private" as const,
    label: "เฉพาะฉัน",
    iconSvg: `<svg viewBox="0 0 20 20" class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="5.5" y="8.5" width="9" height="7.5" rx="1.5" stroke-linecap="round"/><path d="M7.75 8.5V6a2.25 2.25 0 114.5 0v2.5" stroke-linecap="round"/></svg>`,
  },
];

function changeVisibility(
  visibility: "default" | "public" | "followers" | "private"
) {
  emit("update-visibility", props.item.id, visibility);
  isMenuOpen.value = false;
}

function confirmDelete() {
  if (confirm("คุณแน่ใจหรือไม่ที่จะลบกิจกรรมนี้ออกจากฟีดอย่างถาวร?")) {
    emit("delete-activity", props.item.id);
  }
  isMenuOpen.value = false;
}
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition:
    transform 0.15s ease,
    opacity 0.15s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(4px);
  opacity: 0;
}
</style>
