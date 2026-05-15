<template>
  <div class="home">
    <!-- ── Hero Banner ─────────────────────────────────────── -->
    <section class="hero" v-if="featuredMovies.length">
      <div class="hero__slides">
        <TransitionGroup name="hero-slide">
          <div
            v-for="(movie, i) in featuredMovies"
            :key="movie.id"
            v-show="heroIndex === i"
            class="hero__slide"
          >
            <img
              :src="`https://image.tmdb.org/t/p/original${movie.backdrop_path}`"
              :alt="movie.title"
              class="hero__backdrop"
            />
            <div class="hero__overlay" />
          </div>
        </TransitionGroup>
      </div>
      <div class="hero__dots">
        <button
          v-for="(_, i) in featuredMovies"
          :key="i"
          class="hero__dot"
          :class="{ 'hero__dot--active': heroIndex === i }"
          @click="
            heroIndex = i;
            resetTimer();
          "
        />
      </div>
    </section>

    <!-- ── Main content ────────────────────────────────────── -->
    <div class="content-area">
      <!-- Row 1: poster grid + hover popup -->
      <div class="movie-grid-section" v-if="popularMovies.length">
        <div class="movie-grid">
          <div
            v-for="movie in popularMovies.slice(0, 10)"
            :key="movie.id"
            class="poster-wrap"
            @mouseenter="hoveredMovie = movie"
            @mouseleave="hoveredMovie = null"
          >
            <RouterLink :to="`/movies/${movie.id}`" class="poster-card">
              <img
                :src="`https://image.tmdb.org/t/p/w342${movie.poster_path}`"
                :alt="movie.title"
                loading="lazy"
              />
            </RouterLink>

            <!-- Hover popup card -->
            <Transition name="popup">
              <div
                v-if="hoveredMovie?.id === movie.id"
                class="hover-popup"
                :class="
                  getPopupPosition(popularMovies.slice(0, 10).indexOf(movie))
                "
              >
                <div class="popup__trailer">
                  <Play :size="28" class="popup__play-icon" />
                  <span class="popup__trailer-label">trailer auto play</span>
                </div>
                <div class="popup__info">
                  <div class="popup__title-row">
                    <h3 class="popup__title">{{ movie.title }}</h3>
                    <Star :size="13" class="popup__star" />
                    <span class="popup__rating">{{
                      movie.vote_average?.toFixed(1)
                    }}</span>
                  </div>
                  <p class="popup__overview">
                    {{ truncate(movie.overview, 140) }}
                  </p>
                  <div class="popup__actions">
                    <button class="action-btn action-btn--watched">
                      <Eye :size="16" />
                      <span>{{ fmtCount(movie.vote_count) }}</span>
                    </button>
                    <button
                      class="action-btn action-btn--review"
                      @click.prevent="() => {}"
                    >
                      <PenLine :size="16" />
                      <span>Review</span>
                    </button>
                    <button class="action-btn action-btn--favorite">
                      <Heart :size="16" />
                      <span>{{
                        fmtCount(Math.floor(movie.vote_count * 0.6))
                      }}</span>
                    </button>
                    <button class="action-btn action-btn--watchlist">
                      <BookmarkPlus :size="16" />
                      <span>{{
                        fmtCount(Math.floor(movie.vote_count * 0.4))
                      }}</span>
                    </button>
                  </div>
                </div>
              </div>
            </Transition>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div class="pagination" v-if="totalPages > 1">
        <button
          v-for="p in paginationPages"
          :key="p"
          class="page-btn"
          :class="{
            'page-btn--active': p === currentPage,
            'page-btn--ellipsis': p === '...',
          }"
          :disabled="p === '...'"
          @click="typeof p === 'number' && goToPage(p)"
        >
          {{ p }}
        </button>

        <div class="page-size-select" ref="pageSizeRef">
          <button
            class="page-size-trigger"
            @click="pageSizeOpen = !pageSizeOpen"
          >
            {{ pageSize }} <ChevronDown :size="12" />
          </button>
          <div class="page-size-dropdown" v-if="pageSizeOpen">
            <button
              v-for="s in [20, 50, 100, 200, 500]"
              :key="s"
              class="page-size-option"
              :class="{ active: pageSize === s }"
              @click="changePageSize(s)"
            >
              {{ s }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div class="loading-overlay" v-if="isLoading">
      <div class="spinner" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { useQuery } from "@tanstack/vue-query";
import { movieApi } from "@/api/api";
import type { Movie } from "@/types";
import {
  Play,
  Star,
  Eye,
  PenLine,
  Heart,
  BookmarkPlus,
  ChevronDown,
} from "lucide-vue-next";

const router = useRouter();

const heroIndex = ref(0);
const hoveredMovie = ref<Movie | null>(null);
const currentPage = ref(1);
const pageSize = ref(20);
const pageSizeOpen = ref(false);
const pageSizeRef = ref<HTMLElement | null>(null);

let heroTimer: ReturnType<typeof setInterval> | null = null;

const { data: popularData, isLoading } = useQuery({
  queryKey: ["movies-popular", currentPage],
  queryFn: () => movieApi.getPopular(currentPage.value).then((r) => r.data),
});

const popularMovies = computed<Movie[]>(() => popularData.value?.results ?? []);
const totalPages = computed(() => popularData.value?.total_pages ?? 1);
const featuredMovies = computed(() => popularMovies.value.slice(0, 5));

const paginationPages = computed(() => {
  const total = Math.min(totalPages.value, 500);
  const cur = currentPage.value;
  const pages: (number | string)[] = [1, 2, 3];
  if (cur > 4) pages.push("...");
  if (cur > 3 && cur < total - 1) pages.push(cur);
  pages.push("...");
  pages.push(total);
  return [...new Set(pages)];
});

// กำหนดทิศทาง popup ไม่ให้ล้นออกนอกหน้าจอ
// grid 5 คอลัมน์: index 0,5 = ซ้ายสุด → popup ชิดขวา, index 4,9 = ขวาสุด → popup ชิดซ้าย
function getPopupPosition(index: number): string {
  const col = index % 5;
  if (col === 0) return "popup--right";
  if (col === 4) return "popup--left";
  return "popup--center";
}

function resetTimer() {
  if (heroTimer) clearInterval(heroTimer);
  heroTimer = setInterval(() => {
    heroIndex.value =
      (heroIndex.value + 1) % Math.max(featuredMovies.value.length, 1);
  }, 4000);
}

function goToPage(p: number) {
  currentPage.value = p;
  window.scrollTo({ top: 0, behavior: "smooth" });
}

function changePageSize(s: number) {
  pageSize.value = s;
  pageSizeOpen.value = false;
  currentPage.value = 1;
}

function truncate(str: string, len: number) {
  if (!str) return "";
  return str.length > len ? str.slice(0, len) + "…" : str;
}

function fmtCount(n: number) {
  if (n >= 1000) return (n / 1000).toFixed(0) + "k";
  return String(n);
}

function onClickOutside(e: MouseEvent) {
  if (pageSizeRef.value && !pageSizeRef.value.contains(e.target as Node)) {
    pageSizeOpen.value = false;
  }
}

onMounted(() => {
  resetTimer();
  document.addEventListener("click", onClickOutside);
});
onUnmounted(() => {
  if (heroTimer) clearInterval(heroTimer);
  document.removeEventListener("click", onClickOutside);
});
</script>

<style scoped>
.home {
  background: #141414;
  min-height: 100vh;
}

/* ── Hero ─────────────────────────────────────────────── */
.hero {
  position: relative;
  width: 100%;
  height: 300px;
  overflow: hidden;
}
.hero__slides {
  position: absolute;
  inset: 0;
}
.hero__slide {
  position: absolute;
  inset: 0;
}
.hero__backdrop {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 30%;
}
.hero__overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    rgba(0, 0, 0, 0.05) 0%,
    rgba(20, 20, 20, 0.75) 100%
  );
}
.hero__dots {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 6px;
  z-index: 2;
}
.hero__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.35);
  border: none;
  cursor: pointer;
  transition:
    background 0.2s,
    transform 0.2s;
}
.hero__dot--active {
  background: #e50914;
  transform: scale(1.2);
}
.hero-slide-enter-active,
.hero-slide-leave-active {
  transition: opacity 0.6s;
}
.hero-slide-enter-from,
.hero-slide-leave-to {
  opacity: 0;
}

/* ── Content ──────────────────────────────────────────── */
.content-area {
  padding: 1.5rem 1.5rem 3rem;
  max-width: 1100px;
  margin: 0 auto;
}

/* ── Movie grid ───────────────────────────────────────── */
.movie-grid-section {
  margin-bottom: 1.5rem;
}

.movie-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
}

/* Each poster cell needs position:relative for popup anchor */
.poster-wrap {
  position: relative;
}

.poster-card {
  display: block;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition:
    transform 0.2s,
    box-shadow 0.2s;
}
.poster-card:hover {
  transform: scale(1.03);
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.5);
}
.poster-card img {
  width: 100%;
  aspect-ratio: 2/3;
  object-fit: cover;
  display: block;
}

/* ── Hover popup ──────────────────────────────────────── */
.hover-popup {
  position: absolute;
  top: 0;
  z-index: 50;
  width: 280px;
  background: #1c1c1c;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.8);
  pointer-events: auto;
}

/* Position variants */
.popup--right {
  left: calc(100% + 8px);
}
.popup--left {
  right: calc(100% + 8px);
}
.popup--center {
  left: 50%;
  transform: translateX(-50%);
}

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

/* Popup transition */
.popup-enter-active {
  transition:
    opacity 0.15s,
    transform 0.15s;
}
.popup-leave-active {
  transition: opacity 0.1s;
}
.popup-enter-from {
  opacity: 0;
  transform: translateY(6px) scale(0.97);
}
.popup-leave-to {
  opacity: 0;
}

/* ── Pagination ───────────────────────────────────────── */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  margin-top: 1rem;
  flex-wrap: wrap;
}
.page-btn {
  min-width: 36px;
  height: 36px;
  padding: 0 0.5rem;
  background: #1f1f1f;
  border: 1px solid #2a2a2a;
  color: #a3a3a3;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  transition:
    background 0.15s,
    color 0.15s;
}
.page-btn:hover:not(:disabled):not(.page-btn--ellipsis) {
  background: #2a2a2a;
  color: #fff;
}
.page-btn--active {
  background: #e50914;
  border-color: #e50914;
  color: #fff;
}
.page-btn--ellipsis {
  cursor: default;
}

.page-size-select {
  position: relative;
  margin-left: 8px;
}
.page-size-trigger {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 60px;
  height: 36px;
  padding: 0 0.75rem;
  background: #1f1f1f;
  border: 1px solid #2a2a2a;
  color: #a3a3a3;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
}
.page-size-dropdown {
  position: absolute;
  bottom: calc(100% + 4px);
  right: 0;
  background: #1f1f1f;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  overflow: hidden;
  min-width: 70px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
}
.page-size-option {
  display: block;
  width: 100%;
  padding: 0.5rem 0.75rem;
  background: none;
  border: none;
  color: #a3a3a3;
  font-size: 0.875rem;
  cursor: pointer;
  text-align: center;
  transition:
    background 0.15s,
    color 0.15s;
}
.page-size-option:hover,
.page-size-option.active {
  background: rgba(229, 9, 20, 0.15);
  color: #fff;
}

/* ── Loading ──────────────────────────────────────────── */
.loading-overlay {
  position: fixed;
  inset: 0;
  background: rgba(20, 20, 20, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}
.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.15);
  border-top-color: #e50914;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
