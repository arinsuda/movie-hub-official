<template>
  <div class="popup__trailer">
    <Play :size="28" class="popup__play-icon" />
    <span class="popup__trailer-label">trailer auto play</span>
  </div>
  <div class="popup__info">
    <div class="popup__title-row">
      <h3 class="popup__title">{{ movie.title }}</h3>
      <Star :size="13" class="popup__star" />
      <span class="popup__rating">{{ movie.vote_average?.toFixed(1) }}</span>
    </div>
    <p class="popup__overview">{{ truncate(movie.overview, 140) }}</p>
    <div class="popup__actions">
      <button class="action-btn action-btn--watched">
        <Eye :size="16" />
        <span>{{ fmtCount(movie.vote_count) }}</span>
      </button>
      <button class="action-btn action-btn--review">
        <PenLine :size="16" />
        <span>Review</span>
      </button>
      <button class="action-btn action-btn--favorite">
        <Heart :size="16" />
        <span>{{ fmtCount(Math.floor(movie.vote_count * 0.6)) }}</span>
      </button>
      <button class="action-btn action-btn--watchlist">
        <BookmarkPlus :size="16" />
        <span>{{ fmtCount(Math.floor(movie.vote_count * 0.4)) }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Play, Star, Eye, PenLine, Heart, BookmarkPlus } from "lucide-vue-next";
import type { Movie } from "@/types";

defineProps<{ movie: Movie }>();

function truncate(str: string, len: number) {
  if (!str) return "";
  return str.length > len ? str.slice(0, len) + "…" : str;
}

function fmtCount(n: number) {
  return n >= 1000 ? (n / 1000).toFixed(0) + "k" : String(n);
}
</script>

<style scoped>
.popup__trailer {
  aspect-ratio: 16/9;
  background: #000;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}
.popup__play-icon {
  color: #555;
}
.popup__trailer-label {
  font-size: 0.72rem;
  color: #555;
}

.popup__info {
  padding: 0.875rem;
}

.popup__title-row {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  margin-bottom: 0.5rem;
}
.popup__title {
  font-size: 0.9rem;
  font-weight: 700;
  color: #fff;
  margin: 0;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.popup__star {
  color: #f59e0b;
  flex-shrink: 0;
}
.popup__rating {
  font-size: 0.82rem;
  font-weight: 700;
  color: #f59e0b;
  flex-shrink: 0;
}

.popup__overview {
  font-size: 0.72rem;
  color: #888;
  line-height: 1.55;
  margin: 0 0 0.75rem;
}

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
  padding: 0.5rem 0.2rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.62rem;
  font-weight: 700;
  color: #fff;
  transition:
    filter 0.15s,
    transform 0.1s;
}
.action-btn:hover {
  filter: brightness(1.15);
  transform: translateY(-1px);
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
