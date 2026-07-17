<template>
  <div class="popup__media">
    <RouterLink
      :to="{ name: detailRouteName, params: { id: movie.id } }"
      class="popup__media-link"
    >
      <div v-if="showSkeleton" class="popup__skeleton" aria-hidden="true">
        <div class="popup__skeleton-bar popup__skeleton-bar--1" />
        <div class="popup__skeleton-bar popup__skeleton-bar--2" />
      </div>

      <template v-if="showFallback">
        <img
          v-if="backdropUrl"
          :src="backdropUrl"
          :alt="displayTitle"
          class="popup__backdrop"
        />
        <div v-else class="popup__no-media">
          <Film :size="28" class="popup__no-media-icon" />
        </div>
      </template>

      <!-- ── YouTube Player mount target ── -->
      <!-- ใช้ div ธรรมดาแทน <iframe :src="..."> เพราะต้องให้ YT.Player
           คุม element นี้เองถึงจะดัก onError ได้จริง -->
      <div
        v-if="isIframeMounted"
        :id="playerElId"
        class="popup__iframe"
        :class="{ 'popup__iframe--visible': isIframeLoaded }"
      />

      <div v-if="trailerUnavailable" class="popup__no-trailer">
        <VolumeX :size="12" />
        <span>No trailer</span>
      </div>

      <div v-if="isIframeMounted && currentTrailer" class="popup__badge">
        <span class="popup__badge-dot" />
        {{ currentTrailer.type }}
      </div>

      <div class="popup__media-fade" />
    </RouterLink>
  </div>

  <div class="popup__info">
    <div class="popup__title-row">
      <RouterLink
        :to="{ name: detailRouteName, params: { id: movie.id } }"
        class="popup__title-link"
      >
        <h3 class="popup__title">{{ displayTitle }}</h3>
      </RouterLink>
      <div class="popup__rating">
        <Star :size="11" class="popup__star" />
        <span>{{ movie.vote_average?.toFixed(1) }}</span>
      </div>
    </div>

    <div class="popup__meta" v-if="releaseYear">
      <span class="popup__chip">{{ releaseYear }}</span>
      <span v-if="displayRuntime" class="popup__chip">
        {{ fmtRuntime(displayRuntime) }}
      </span>
    </div>

    <p class="popup__overview">{{ truncate(movie.overview, 130) }}</p>

    <div class="popup__actions">
      <button class="action-btn action-btn--watched" title="Views Static">
        <Eye :size="15" />
        <span>{{ fmtCount(stats.view_count) }}</span>
      </button>

      <RouterLink
        :to="{
          name: detailRouteName,
          params: { id: movie.id },
          query: { action: 'review' },
        }"
        class="action-btn action-btn--review"
        title="Review"
      >
        <PenLine :size="15" />
        <span>{{ fmtCount(stats.review_count) }}</span>
      </RouterLink>

      <button
        class="action-btn action-btn--favorite"
        :class="{ 'action-btn--active': isLiked }"
        title="Favourite"
        @click.stop="handleLikeToggle"
      >
        <Heart :size="15" :fill="isLiked ? 'currentColor' : 'none'" />
        <span>{{ fmtCount(stats.like_count) }}</span>
      </button>

      <button
        class="action-btn action-btn--watchlist"
        :class="{ 'action-btn--active': isInWatchlist }"
        title="Watchlist"
        @click.stop="handleWatchlistToggle"
      >
        <BookmarkPlus
          :size="15"
          :fill="isInWatchlist ? 'currentColor' : 'none'"
        />
        <span>{{ fmtCount(stats.watchlist_count) }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  computed,
  ref,
  onMounted,
  onBeforeUnmount,
  watch,
  nextTick,
} from "vue";
import {
  Film,
  Star,
  Eye,
  PenLine,
  Heart,
  BookmarkPlus,
  VolumeX,
} from "lucide-vue-next";
import { useImageUrl } from "@/composables/useImageUrl";
import type { Movie, TVSeries } from "@/types";
import type { ResolvedTrailer } from "@/composables/useTrailerPreview";
import { libraryApi } from "@/api/endpoints/library";
import { useAuthStore } from "@/stores/auth";
import { mediaApi } from "@/api/endpoints/media";

const props = defineProps<{
  movie: Movie | TVSeries;
  // ── Trailer/player state: มาจาก useTrailerPreview instance ของการ์ดนี้ ──
  currentTrailer: ResolvedTrailer | null;
  trailerUnavailable: boolean;
  isIframeMounted: boolean;
  isIframeLoaded: boolean;
  showSkeleton: boolean;
  showFallback: boolean;
  attachPlayer: (elementId: string) => void;
  mediaType?: "movie" | "tv";
}>();

const authStore = useAuthStore();
const currentUserId = computed(() => authStore.user?.id ?? null);

// ─── Derived media type (default: "movie") ────────────────────────────────
const mt = computed(() => props.mediaType ?? "movie");

// ─── Route name ───────────────────────────────────────────────────────────
const detailRouteName = computed(() =>
  mt.value === "tv" ? "tv-detail" : "movie-detail",
);

// ─── Derived display fields (รองรับทั้ง Movie และ TVSeries) ───────────────
const displayTitle = computed(() => {
  const m = props.movie as any;
  return m.title ?? m.name ?? "ไม่ทราบชื่อ";
});

const displayReleaseDate = computed(() => {
  const m = props.movie as any;
  return m.release_date ?? m.first_air_date ?? "";
});

const releaseYear = computed(() =>
  displayReleaseDate.value ? displayReleaseDate.value.slice(0, 4) : null,
);

const displayRuntime = computed(() => {
  const m = props.movie as any;
  if (typeof m.runtime === "number") return m.runtime;
  if (Array.isArray(m.episode_run_time) && m.episode_run_time.length) {
    return m.episode_run_time[0];
  }
  return null;
});

// ─── YouTube Player mount ───────────────────────────────────────────────────
// unique element id ต่อการ์ด กันชนกันเวลามีหลายป็อปอัพ instance พร้อมกันใน DOM
const playerElId = `yt-player-${props.movie.id}-${Math.random()
  .toString(36)
  .slice(2)}`;

// เมื่อ isIframeMounted กลายเป็น true (ไม่ว่าจะครั้งแรก หรือ retry candidate ถัดไป)
// ให้เรียก attachPlayer เพื่อสร้าง YT.Player ผูกกับ div ของการ์ดนี้
watch(
  () => props.isIframeMounted,
  async (mounted) => {
    if (!mounted) return;
    await nextTick(); // รอให้ div :id="playerElId" render เสร็จก่อน
    props.attachPlayer(playerElId);
  },
);

onMounted(() => {
  if (props.isIframeMounted) {
    nextTick(() => props.attachPlayer(playerElId));
  }
});

// ─── Reactive Stats State ─────────────────────────────────────────────────
const stats = ref({
  view_count: 0,
  like_count: 0,
  review_count: 0,
  watchlist_count: 0,
});
const isLiked = ref(false);
const isInWatchlist = ref(false);
const watchlistItemId = ref<number | null>(null);

onMounted(async () => {
  try {
    const resStats = await mediaApi.getMediaStats(mt.value, props.movie.id);
    const s = resStats.data?.stats || (resStats.data as any)?.stats;

    if (s) {
      stats.value = {
        view_count: s.view_count || 0,
        like_count: s.like_count || 0,
        review_count: s.review_count || 0,
        watchlist_count: s.watchlist_count || 0,
      };
      isLiked.value = !!s.liked_at;
    }

    if (!currentUserId.value) return;

    const userId = currentUserId.value;
    const resLibrary = await libraryApi.getOwnMediaStatus(
      props.movie.id,
      mt.value,
    );
    const watchlistInfo = resLibrary.data?.in_lists?.find(
      (item: any) => item.list_type === "watchlist",
    );

    if (watchlistInfo) {
      isInWatchlist.value = true;
      watchlistItemId.value = watchlistInfo.item_id;
    }
  } catch (err) {
    console.error("ไม่สามารถโหลดสถิติจริงจากระบบหลังบ้านได้:", err);
  }
});

// ─── Like Toggle ──────────────────────────────────────────────────────────
async function handleLikeToggle() {
  try {
    if (isLiked.value) {
      await mediaApi.unlikeMedia(mt.value, props.movie.id);
      stats.value.like_count = Math.max(0, stats.value.like_count - 1);
      isLiked.value = false;
      window.$toast?.info(`ลบออกจากรายการที่ชอบแล้ว`, displayTitle.value);
    } else {
      await mediaApi.likeMedia(mt.value, props.movie.id);
      stats.value.like_count++;
      isLiked.value = true;
      window.$toast?.success(
        `เพิ่มไปยังรายการที่ชอบเรียบร้อย! ❤️`,
        displayTitle.value,
      );
    }
  } catch (err) {
    window.$toast?.error("เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", "REMOV HUB");
  }
}

// ─── Watchlist Toggle ─────────────────────────────────────────────────────
async function handleWatchlistToggle() {
  if (!currentUserId.value) {
    window.$toast?.warning("กรุณาเข้าสู่ระบบก่อนใช้งาน", "แจ้งเตือน");
    return;
  }

  try {
    if (isInWatchlist.value && watchlistItemId.value) {
      await libraryApi.removeItem(watchlistItemId.value);
      stats.value.watchlist_count = Math.max(
        0,
        stats.value.watchlist_count - 1,
      );
      isInWatchlist.value = false;
      watchlistItemId.value = null;
      window.$toast?.info(`ลบออกจาก Watchlist แล้ว`, displayTitle.value);
    } else {
      const res = await libraryApi.addItem({
        media_id: props.movie.id,
        media_type: mt.value,
        list_type: "watchlist",
      });

      if (res.data?.item) {
        watchlistItemId.value = res.data.item.id;
      }

      stats.value.watchlist_count++;
      isInWatchlist.value = true;
      window.$toast?.success(
        `เพิ่มเข้า Watchlist สำเร็จ! 🍿`,
        displayTitle.value,
      );
    }
  } catch (err) {
    console.error("Watchlist Error:", err);
    window.$toast?.error(
      "บันทึกข้อมูลไม่สำเร็จ กรุณาลองใหม่อีกครั้ง",
      "REMOV HUB",
    );
  }
}

// ─── Composables & Computed ───────────────────────────────────────────────
const { backdropImage } = useImageUrl();

const backdropUrl = computed(() => {
  const path = props.movie.backdrop_path;
  return path ? backdropImage(path) : null;
});

function truncate(str: string, len: number) {
  if (!str) return "";
  return str.length > len ? str.slice(0, len) + "…" : str;
}

function fmtCount(n: number) {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + "M";
  if (n >= 1_000) return (n / 1_000).toFixed(0) + "K";
  return String(n);
}

function fmtRuntime(mins: number) {
  const h = Math.floor(mins / 60);
  const m = mins % 60;
  return h > 0 ? `${h}h ${m}m` : `${m}m`;
}
</script>

<style scoped>
.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  gap: 3px;
  padding: 0.45rem 0.15rem;
  border: none;
  border-radius: 7px;
  cursor: pointer;
  font-size: 0.59rem;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  transition:
    filter 0.15s,
    transform 0.12s,
    background-color 0.2s;
}
.action-btn:hover {
  filter: brightness(1.2);
  transform: translateY(-1px);
}
.action-btn--active {
  box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.4);
  filter: saturate(1.4) brightness(1.1);
}

.popup__media {
  position: relative;
  aspect-ratio: 16 / 9;
  background: #080808;
  overflow: hidden;
}
.popup__backdrop {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center top;
}
.popup__no-media {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #111 0%, #1c1c1c 100%);
}
.popup__no-media-icon {
  color: #2a2a2a;
}
.popup__iframe {
  position: absolute;
  inset: -1px;
  width: calc(100% + 2px);
  height: calc(100% + 2px);
  border: none;
  opacity: 0;
  transition: opacity 0.45s ease;
  pointer-events: none;
}
.popup__iframe--visible {
  opacity: 1;
}
.popup__media-fade {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 55%;
  background: linear-gradient(to bottom, transparent, #1c1c1c);
  pointer-events: none;
}
.popup__no-trailer {
  position: absolute;
  bottom: 8px;
  left: 10px;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.62rem;
  color: #444;
  z-index: 2;
}
.popup__badge {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.58rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.5);
  background: rgba(0, 0, 0, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 4px;
  padding: 3px 6px;
  backdrop-filter: blur(6px);
  z-index: 3;
}
.popup__badge-dot {
  display: inline-block;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #e50914;
  animation: pulse-dot 1.8s ease-in-out infinite;
}
.popup__media-link {
  display: block;
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  z-index: 4;
  cursor: pointer;
}
@keyframes pulse-dot {
  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.35;
    transform: scale(0.65);
  }
}
.popup__skeleton {
  position: absolute;
  inset: 0;
  background: #101010;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  z-index: 1;
}
.popup__skeleton-bar {
  border-radius: 4px;
  background: linear-gradient(90deg, #1a1a1a 25%, #252525 50%, #1a1a1a 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s ease-in-out infinite;
}
.popup__skeleton-bar--1 {
  width: 55%;
  height: 7px;
}
.popup__skeleton-bar--2 {
  width: 38%;
  height: 5px;
}
@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
.popup__info {
  padding: 0.8rem 0.875rem 0.875rem;
}
.popup__title-row {
  display: flex;
  align-items: flex-start;
  gap: 0.4rem;
  margin-bottom: 0.4rem;
}
.popup__title-link {
  cursor: pointer;
  color: inherit;
}
.popup__title {
  flex: 1;
  font-size: 0.88rem;
  font-weight: 700;
  color: #f0f0f0;
  margin: 0;
  line-height: 1.3;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.popup__rating {
  display: flex;
  align-items: center;
  gap: 3px;
  flex-shrink: 0;
  margin-top: 2px;
  font-size: 0.76rem;
  font-weight: 700;
  color: #f59e0b;
}
.popup__star {
  color: #f59e0b;
}
.popup__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 0.5rem;
}
.popup__chip {
  font-size: 0.6rem;
  font-weight: 600;
  color: #666;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 4px;
  padding: 2px 6px;
}
.popup__overview {
  font-size: 0.7rem;
  color: #666;
  line-height: 1.6;
  margin: 0 0 0.75rem;
}
.popup__actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 5px;
}
.action-btn--watched {
  background: #1d4ed8;
}
.action-btn--review {
  background: #16a34a;
}
.action-btn--favorite {
  background: #e50914;
}
.action-btn--watchlist {
  background: #d97706;
}
</style>
