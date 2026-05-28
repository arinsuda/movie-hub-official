<template>
  <div class="list-page">
    <div class="page-header">
      <h1 class="page-title">MOVIES</h1>
      <div class="tabs">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="tab-btn"
          :class="{ active: activeTab === tab.key }"
          @click="switchTab(tab.key)"
        >
          {{ tab.label }}
        </button>
      </div>
      <div class="filter-bar">
        <div class="filter-select" ref="genreRef">
          <button class="filter-trigger" @click="genreOpen = !genreOpen">
            <Tag :size="14" />
            {{
              selectedGenre
                ? genres.find((g) => g.id === selectedGenre)?.name
                : "All Genres"
            }}
            <ChevronDown
              :size="13"
              :class="{ 'rotate-180': genreOpen }"
              class="chevron"
            />
          </button>
          <div class="filter-dropdown" v-if="genreOpen">
            <button
              class="filter-option"
              :class="{ active: !selectedGenre }"
              @click="selectGenre(null)"
            >
              All Genres
            </button>
            <button
              v-for="g in genres"
              :key="g.id"
              class="filter-option"
              :class="{ active: selectedGenre === g.id }"
              @click="selectGenre(g.id)"
            >
              {{ g.name }}
            </button>
          </div>
        </div>
        <div class="filter-select" ref="sortRef">
          <button class="filter-trigger" @click="sortOpen = !sortOpen">
            <ArrowUpDown :size="14" />
            {{ sortOptions.find((s) => s.key === sortBy)?.label }}
            <ChevronDown
              :size="13"
              :class="{ 'rotate-180': sortOpen }"
              class="chevron"
            />
          </button>
          <div class="filter-dropdown" v-if="sortOpen">
            <button
              v-for="s in sortOptions"
              :key="s.key"
              class="filter-option"
              :class="{ active: sortBy === s.key }"
              @click="selectSort(s.key)"
            >
              {{ s.label }}
            </button>
          </div>
        </div>
        <div class="search-inline">
          <Search :size="14" class="search-ic" />
          <input
            v-model="searchQuery"
            class="search-input"
            placeholder="Search movies..."
            @keydown.enter="doSearch"
          />
          <button v-if="searchQuery" class="search-clear" @click="clearSearch">
            <X :size="13" />
          </button>
        </div>
      </div>
    </div>

    <div class="content-area">
      <div class="movie-grid" v-if="movies.length">
        <div
          v-for="(show, i) in movies"
          :key="show.id"
          class="poster-wrap"
          @mouseenter="hoveredId = show.id"
          @mouseleave="hoveredId = null"
        >
          <RouterLink :to="`/movies/${show.id}`" class="poster-card">
            <img
              :src="
                show.poster_path
                  ? `https://image.tmdb.org/t/p/w342${show.poster_path}`
                  : '/placeholder.jpg'
              "
              :alt="show.original_title"
              loading="lazy"
            />
          </RouterLink>
          <Transition name="popup">
            <div
              v-if="hoveredId === show.id"
              class="hover-popup"
              :class="getPopupPos(i)"
            >
              <div class="popup__trailer">
                <Play :size="26" class="popup__play-icon" />
                <span class="popup__trailer-label">trailer auto play</span>
              </div>
              <div class="popup__info">
                <div class="popup__title-row">
                  <h3 class="popup__title">{{ show.original_title }}</h3>
                  <Star :size="12" class="popup__star" />
                  <span class="popup__rating">{{
                    show.vote_average?.toFixed(1)
                  }}</span>
                </div>
                <p class="popup__meta">
                  {{ show.release_date?.slice(0, 4) }}
                </p>
                <p class="popup__overview">
                  {{ truncate(show.overview, 130) }}
                </p>
                <div class="popup__actions">
                  <button class="action-btn action-btn--watched">
                    <Eye :size="15" /><span>{{
                      fmtCount(show.vote_count)
                    }}</span>
                  </button>
                  <button class="action-btn action-btn--review">
                    <PenLine :size="15" /><span>Review</span>
                  </button>
                  <button class="action-btn action-btn--favorite">
                    <Heart :size="15" /><span>{{
                      fmtCount(Math.floor(show.vote_count * 0.6))
                    }}</span>
                  </button>
                  <button class="action-btn action-btn--watchlist">
                    <BookmarkPlus :size="15" /><span>{{
                      fmtCount(Math.floor(show.vote_count * 0.4))
                    }}</span>
                  </button>
                </div>
              </div>
            </div>
          </Transition>
        </div>
      </div>
      <div class="empty-state" v-else-if="!isLoading">
        <Film :size="48" class="empty-icon" />
        <p>No movies found</p>
      </div>

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
        <div class="page-size-wrap" ref="pageSizeRef">
          <button
            class="page-size-trigger"
            @click="pageSizeOpen = !pageSizeOpen"
          >
            {{ pageSize }} <ChevronDown :size="12" />
          </button>
          <div class="page-size-dropdown" v-if="pageSizeOpen">
            <button
              v-for="s in [20, 50, 100]"
              :key="s"
              class="page-size-opt"
              :class="{ active: pageSize === s }"
              @click="changeSize(s)"
            >
              {{ s }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div class="loading-overlay" v-if="isLoading"><div class="spinner" /></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { useRoute } from "vue-router";
import { useQuery } from "@tanstack/vue-query";
import { movieApi } from "@/api/api";
import type { Movie, Genre } from "@/types";
import {
  Search,
  X,
  Tag,
  ArrowUpDown,
  ChevronDown,
  Play,
  Star,
  Eye,
  PenLine,
  Heart,
  BookmarkPlus,
  Film,
} from "lucide-vue-next";
import {
  resolveTrailer,
  useTrailerPreview,
  type ResolvedTrailer,
} from "@/composables/useTrailerPreview";

const route = useRoute();

type TabKey = "popular" | "now_playing" | "top_rated";
const tabs: { key: TabKey; label: string }[] = [
  { key: "popular", label: "Popular" },
  { key: "now_playing", label: "Now Playing" },
  { key: "top_rated", label: "Top Rated" },
];
const sortOptions = [
  { key: "default", label: "Default" },
  { key: "rating", label: "Rating" },
  { key: "date", label: "First Air Date" },
  { key: "popular", label: "Popularity" },
];

const activeTab = ref<TabKey>("popular");
const sortBy = ref("default");
const selectedGenre = ref<number | null>(null);
const searchQuery = ref((route.query.q as string) ?? "");
const currentPage = ref(1);
const pageSize = ref(20);
const hoveredId = ref<number | null>(null);
const cardStates = new Map<number, ReturnType<typeof useTrailerPreview>>();
const cardTrailers = new Map<number, ResolvedTrailer | null>();
let insidePopup = false;
let showTimer: ReturnType<typeof setTimeout> | null = null;
const SHOW_DELAY = 200;
const genres = ref<Genre[]>([]);
const genreOpen = ref(false);
const sortOpen = ref(false);
const pageSizeOpen = ref(false);
const genreRef = ref<HTMLElement | null>(null);
const sortRef = ref<HTMLElement | null>(null);
const pageSizeRef = ref<HTMLElement | null>(null);

const isSearchMode = computed(() => !!searchQuery.value.trim());
const queryKey = computed(() => [
  "movies",
  activeTab.value,
  currentPage.value,
  searchQuery.value,
]);

const { data, isLoading } = useQuery({
  queryKey,
  queryFn: async () => {
    if (isSearchMode.value)
      return movieApi
        .search(searchQuery.value.trim(), currentPage.value)
        .then((r) => r.data);
    const fn = {
      popular: () => movieApi.getPopular(currentPage.value),
      now_playing: () => movieApi.getNowPlaying(currentPage.value),
      top_rated: () => movieApi.getTopRated(currentPage.value),
    }[activeTab.value];
    return fn().then((r) => r.data);
  },
});

const rawMovies = computed<Movie[]>(() => data.value?.results ?? []);
const totalPages = computed(() => data.value?.total_pages ?? 1);

const movies = computed(() => {
  let list = [...rawMovies.value];
  if (selectedGenre.value)
    list = list.filter((s) => s.genre_ids?.includes(selectedGenre.value!));
  if (sortBy.value === "rating")
    list.sort((a, b) => b.vote_average - a.vote_average);
  if (sortBy.value === "date")
    list.sort((a, b) =>
      (b.release_date ?? "").localeCompare(a.release_date ?? ""),
    );
  if (sortBy.value === "popular")
    list.sort((a, b) => b.popularity - a.popularity);
  return list;
});

const paginationPages = computed(() => {
  const total = Math.min(totalPages.value, 500);
  const cur = currentPage.value;
  const p: (number | string)[] = [];
  if (cur > 2) p.push(1);
  if (cur > 3) p.push("...");
  for (let i = Math.max(1, cur - 1); i <= Math.min(total, cur + 1); i++)
    p.push(i);
  if (cur < total - 2) p.push("...");
  if (cur < total - 1) p.push(total);
  return [...new Set(p)];
});

function getPopupPos(i: number) {
  const col = i % 5;
  if (col === 0) return "popup--right";
  if (col === 4) return "popup--left";
  return "popup--center";
}
function switchTab(key: TabKey) {
  activeTab.value = key;
  currentPage.value = 1;
  searchQuery.value = "";
}
function doSearch() {
  currentPage.value = 1;
}
function clearSearch() {
  searchQuery.value = "";
  currentPage.value = 1;
}
function selectGenre(id: number | null) {
  selectedGenre.value = id;
  genreOpen.value = false;
}
function selectSort(key: string) {
  sortBy.value = key;
  sortOpen.value = false;
}
function goToPage(p: number) {
  currentPage.value = p;
  window.scrollTo({ top: 0, behavior: "smooth" });
}
function changeSize(s: number) {
  pageSize.value = s;
  pageSizeOpen.value = false;
  currentPage.value = 1;
}
function truncate(s: string, n: number) {
  return s?.length > n ? s.slice(0, n) + "…" : (s ?? "");
}
function fmtCount(n: number) {
  return n >= 1000 ? (n / 1000).toFixed(0) + "k" : String(n);
}
function onClickOutside(e: MouseEvent) {
  const t = e.target as Node;
  if (genreRef.value && !genreRef.value.contains(t)) genreOpen.value = false;
  if (sortRef.value && !sortRef.value.contains(t)) sortOpen.value = false;
  if (pageSizeRef.value && !pageSizeRef.value.contains(t))
    pageSizeOpen.value = false;
}

function getState(movieId: number) {
  if (!cardStates.has(movieId)) {
    cardStates.set(movieId, useTrailerPreview({ mountDelay: 500 }));
  }
  return cardStates.get(movieId)!;
}

function getTrailer(movieId: number): ResolvedTrailer | null {
  return cardTrailers.get(movieId) ?? null;
}

async function fetchAndCacheTrailer(movie: Movie) {
  if (cardTrailers.has(movie.id)) return;
  cardTrailers.set(movie.id, null);
  try {
    const res = await movieApi.getVideos(movie.id);
    const videos = res.data?.results ?? [];
    const trailer = resolveTrailer(videos);
    cardTrailers.set(movie.id, trailer);
    if (hoveredId.value === movie.id && trailer) {
      getState(movie.id).scheduleMount();
    }
  } catch {
    cardTrailers.set(movie.id, null);
  }
}

function onCardEnter(movieId: number, movie: Movie) {
  clearTimeout(showTimer ?? undefined);
  insidePopup = false;
  fetchAndCacheTrailer(movie);
  showTimer = setTimeout(() => {
    hoveredId.value = movieId;
    const trailer = getTrailer(movieId);
    if (trailer) getState(movieId).scheduleMount();
  }, SHOW_DELAY);
}

function onCardLeave(movieId: number) {
  clearTimeout(showTimer ?? undefined);
  setTimeout(() => {
    if (!insidePopup) closeCard(movieId);
  }, 80);
}

function onPopupEnter() {
  insidePopup = true;
}

function onPopupLeave(movieId: number) {
  insidePopup = false;
  closeCard(movieId);
}

function closeCard(movieId: number) {
  if (hoveredId.value !== movieId) return;
  hoveredId.value = null;
  cardStates.get(movieId)?.unmount();
}

watch(
  () => route.query.q,
  (q) => {
    searchQuery.value = (q as string) ?? "";
    currentPage.value = 1;
  },
);
onMounted(async () => {
  document.addEventListener("click", onClickOutside);
  const res = await movieApi.getGenres();
  genres.value = res.data.genres;
});
onUnmounted(() => {
  document.removeEventListener("click", onClickOutside);
  clearTimeout(showTimer ?? undefined);
  cardStates.forEach((s) => s.unmount());
  cardStates.clear();
});
</script>

<style scoped>
.list-page {
  background: #141414;
  min-height: 100vh;
}
.page-header {
  padding: 5.5rem 1.5rem 0;
  max-width: 1100px;
  margin: 0 auto;
}
.page-title {
  font-family: "Impact", "Arial Black", sans-serif;
  font-size: 1.75rem;
  font-weight: 900;
  letter-spacing: 2px;
  color: #fff;
  margin: 0 0 1.25rem;
}
.tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 1rem;
}
.tab-btn {
  padding: 0.4rem 1rem;
  background: #1f1f1f;
  border: 1px solid #2a2a2a;
  border-radius: 9999px;
  font-size: 0.82rem;
  font-weight: 600;
  color: #a3a3a3;
  cursor: pointer;
  transition: all 0.2s;
}
.tab-btn:hover {
  color: #fff;
  border-color: #444;
}
.tab-btn.active {
  background: #e50914;
  border-color: #e50914;
  color: #fff;
}
.filter-bar {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding-bottom: 1.25rem;
  flex-wrap: wrap;
}
.filter-select {
  position: relative;
}
.filter-trigger {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.45rem 0.875rem;
  background: #1f1f1f;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  color: #ccc;
  font-size: 0.82rem;
  cursor: pointer;
  transition: border-color 0.2s;
  white-space: nowrap;
}
.filter-trigger:hover {
  border-color: #555;
  color: #fff;
}
.chevron {
  transition: transform 0.2s;
}
.rotate-180 {
  transform: rotate(180deg);
}
.filter-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  min-width: 160px;
  max-height: 280px;
  overflow-y: auto;
  background: #1f1f1f;
  border: 1px solid #2a2a2a;
  border-radius: 10px;
  padding: 0.35rem;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
  z-index: 30;
}
.filter-option {
  display: block;
  width: 100%;
  padding: 0.45rem 0.75rem;
  background: none;
  border: none;
  color: #ccc;
  font-size: 0.82rem;
  cursor: pointer;
  border-radius: 6px;
  text-align: left;
  transition:
    background 0.15s,
    color 0.15s;
}
.filter-option:hover {
  background: rgba(255, 255, 255, 0.07);
  color: #fff;
}
.filter-option.active {
  color: #e50914;
  font-weight: 600;
}
.search-inline {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #1f1f1f;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  padding: 0.4rem 0.75rem;
  flex: 1;
  max-width: 320px;
}
.search-inline:focus-within {
  border-color: #555;
}
.search-ic {
  color: #666;
  flex-shrink: 0;
}
.search-input {
  background: none;
  border: none;
  outline: none;
  color: #fff;
  font-size: 0.875rem;
  width: 100%;
}
.search-input::placeholder {
  color: #555;
}
.search-clear {
  background: none;
  border: none;
  cursor: pointer;
  color: #666;
  display: flex;
  padding: 0;
}
.search-clear:hover {
  color: #ccc;
}
.content-area {
  padding: 0 1.5rem 3rem;
  max-width: 1100px;
  margin: 0 auto;
}
.movie-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
  margin-bottom: 1.5rem;
}
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
}
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
  gap: 0.4rem;
}
.popup__play-icon {
  color: #555;
}
.popup__trailer-label {
  font-size: 0.7rem;
  color: #555;
}
.popup__info {
  padding: 0.875rem;
}
.popup__title-row {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  margin-bottom: 0.25rem;
}
.popup__title {
  font-size: 0.875rem;
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
  font-size: 0.8rem;
  font-weight: 700;
  color: #f59e0b;
  flex-shrink: 0;
}
.popup__meta {
  font-size: 0.7rem;
  color: #666;
  margin: 0 0 0.35rem;
}
.popup__overview {
  font-size: 0.7rem;
  color: #888;
  line-height: 1.5;
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
  padding: 0.45rem 0.2rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.6rem;
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
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 5rem 0;
  color: #555;
}
.empty-icon {
  opacity: 0.4;
}
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
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
.page-size-wrap {
  position: relative;
  margin-left: 8px;
}
.page-size-trigger {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 58px;
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
.page-size-opt {
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
.page-size-opt:hover,
.page-size-opt.active {
  background: rgba(229, 9, 20, 0.15);
  color: #fff;
}
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
