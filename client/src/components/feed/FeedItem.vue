<template>
  <article
    class="group flex gap-3 border-b border-white/5 px-4 py-4 transition-colors hover:bg-white/[0.03]"
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

        <button
          v-if="canManage"
          type="button"
          class="shrink-0 rounded p-1 text-white/30 opacity-0 transition-opacity hover:bg-white/10 hover:text-white/70 group-hover:opacity-100"
          title="ซ่อนกิจกรรมนี้จากฟีด"
          @click="emit('hide', item.id)"
        >
          <svg
            viewBox="0 0 20 20"
            class="h-4 w-4"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path d="M5 5l10 10M15 5L5 15" stroke-linecap="round" />
          </svg>
        </button>
      </div>

      <div class="mt-1 flex items-center gap-1.5 text-xs text-white/40">
        <span>{{ activityMeta.icon }}</span>
        <span>{{ activityMeta.label }}</span>
        <span aria-hidden="true">·</span>
        <time :datetime="item.created_at">{{ relativeTime }}</time>
      </div>

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
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { FeedItemResponse } from "@/types/feed";
import { getActivityMeta, mediaRoutePrefix } from "@/utils/feedActivity";
import { formatRelativeTime } from "@/utils/relativeTime";

const props = defineProps<{
  item: FeedItemResponse;
  canManage?: boolean;
}>();

const emit = defineEmits<{ hide: [id: number] }>();

// TODO: เปลี่ยน path ให้ตรงกับ default avatar จริงของโปรเจกต์
const defaultAvatar = "/images/default-avatar.png";

const activityMeta = computed(() => getActivityMeta(props.item.type));
const relativeTime = computed(() => formatRelativeTime(props.item.created_at));
const actorLink = computed(() => `/users/${props.item.actor.id}`);

const mediaLink = computed(() => {
  if (!props.item.media) return "/";
  return `/${mediaRoutePrefix(props.item.media.media_type)}/${props.item.media.id}`;
});
</script>
