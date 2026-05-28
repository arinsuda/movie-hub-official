<template>
  <!-- ── Trailer / Media zone ─────────────────────────────────────────────── -->
  <div class="popup__media">
    <!-- 1. Skeleton shimmer while iframe loads -->
    <div v-if="showSkeleton" class="popup__skeleton" aria-hidden="true">
      <div class="popup__skeleton-bar popup__skeleton-bar--1" />
      <div class="popup__skeleton-bar popup__skeleton-bar--2" />
    </div>

    <!-- 2. Fallback: backdrop image (or icon) before iframe mounts -->
    <template v-if="showFallback">
      <img
        v-if="backdropUrl"
        :src="backdropUrl"
        :alt="movie.title"
        class="popup__backdrop"
      />
      <div v-else class="popup__no-media">
        <Film :size="28" class="popup__no-media-icon" />
      </div>
    </template>

    <!-- 3. YouTube iframe — only rendered after hover delay -->
    <iframe
      v-if="isIframeMounted && trailer"
      :src="trailer.embedUrl"
      class="popup__iframe"
      :class="{ 'popup__iframe--visible': isIframeLoaded }"
      frameborder="0"
      allow="autoplay; encrypted-media"
      referrerpolicy="strict-origin-when-cross-origin"
      @load="onIframeLoad"
    />

    <!-- No trailer notice (only shown while fallback is active) -->
    <div v-if="showFallback && !trailer" class="popup__no-trailer">
      <VolumeX :size="12" />
      <span>No trailer</span>
    </div>

    <!-- Live badge when trailer is playing -->
    <div v-if="isIframeMounted && trailer" class="popup__badge">
      <span class="popup__badge-dot" />
      {{ trailer.type }}
    </div>

    <!-- Gradient overlay fading into info section -->
    <div class="popup__media-fade" />
  </div>

  <!-- ── Info section ──────────────────────────────────────────────────────── -->
  <div class="popup__info">
    <!-- Title + rating -->
    <div class="popup__title-row">
      <h3 class="popup__title">{{ movie.title }}</h3>
      <div class="popup__rating">
        <Star :size="11" class="popup__star" />
        <span>{{ movie.vote_average?.toFixed(1) }}</span>
      </div>
    </div>

    <!-- Meta: year + runtime -->
    <div class="popup__meta" v-if="releaseYear">
      <span class="popup__chip">{{ releaseYear }}</span>
      <span
        v-if="'runtime' in movie && (movie as any).runtime"
        class="popup__chip"
      >
        {{ fmtRuntime((movie as any).runtime) }}
      </span>
    </div>

    <!-- Overview -->
    <p class="popup__overview">{{ truncate(movie.overview, 130) }}</p>

    <!-- Action buttons -->
    <div class="popup__actions">
      <button class="action-btn action-btn--watched" title="Watched">
        <Eye :size="15" />
        <span>{{ fmtCount(movie.vote_count) }}</span>
      </button>
      <button class="action-btn action-btn--review" title="Review">
        <PenLine :size="15" />
        <span>Review</span>
      </button>
      <button class="action-btn action-btn--favorite" title="Favourite">
        <Heart :size="15" />
        <span>{{ fmtCount(Math.floor(movie.vote_count * 0.6)) }}</span>
      </button>
      <button class="action-btn action-btn--watchlist" title="Watchlist">
        <BookmarkPlus :size="15" />
        <span>Watchlist</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
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
import type { Movie } from "@/types";
import type { ResolvedTrailer } from "@/composables/useTrailerPreview";

// ─── Props ────────────────────────────────────────────────────────────────────

const props = defineProps<{
  movie: Movie;
  trailer: ResolvedTrailer | null;
  isIframeMounted: boolean;
  isIframeLoaded: boolean;
  showSkeleton: boolean;
  showFallback: boolean;
}>();

const emit = defineEmits<{
  (e: "iframe-load"): void;
}>();

// ─── Composables ──────────────────────────────────────────────────────────────

const { backdropImage } = useImageUrl();

// ─── Computed ──────────────────────────────────────────────────────────────────

const backdropUrl = computed(() => {
  const path = props.movie.backdrop_path;
  return path ? backdropImage(path) : null;
});

const releaseYear = computed(() =>
  props.movie.release_date ? props.movie.release_date.slice(0, 4) : null,
);

// ─── Handlers ────────────────────────────────────────────────────────────────

function onIframeLoad() {
  emit("iframe-load");
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

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
/* ── Media zone ──────────────────────────────────────────────────────────── */
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

/* iframe */
.popup__iframe {
  position: absolute;
  inset: -1px;
  width: calc(100% + 2px);
  height: calc(100% + 2px);
  border: none;
  opacity: 0;
  transition: opacity 0.45s ease;
}

.popup__iframe--visible {
  opacity: 1;
}

/* Gradient fade into info */
.popup__media-fade {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 55%;
  background: linear-gradient(to bottom, transparent, #1c1c1c);
  pointer-events: none;
}

/* No-trailer notice */
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

/* Live trailer badge */
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

/* ── Skeleton shimmer ──────────────────────────────────────────────────── */
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

/* ── Info section ──────────────────────────────────────────────────────── */
.popup__info {
  padding: 0.8rem 0.875rem 0.875rem;
}

.popup__title-row {
  display: flex;
  align-items: flex-start;
  gap: 0.4rem;
  margin-bottom: 0.4rem;
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

/* Meta chips */
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

/* Overview */
.popup__overview {
  font-size: 0.7rem;
  color: #666;
  line-height: 1.6;
  margin: 0 0 0.75rem;
}

/* Action buttons */
.popup__actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 5px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
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
    transform 0.12s;
}

.action-btn:hover {
  filter: brightness(1.2);
  transform: translateY(-1px);
}

.action-btn:active {
  transform: translateY(0);
  filter: brightness(0.9);
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
