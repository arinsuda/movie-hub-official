<template>
  <div class="home">
    <section class="hero" v-if="heroMovies.length">
      <div class="hero__slides">
        <TransitionGroup name="hero-slide">
          <div
            v-for="(movie, i) in heroMovies"
            :key="movie.id"
            v-show="heroIndex === i"
            class="hero__slide"
          >
            <img
              v-if="movie.backdrop_path && !failedBackdrops.has(movie.id)"
              :src="`https://image.tmdb.org/t/p/original${movie.backdrop_path}`"
              :alt="movie.title"
              class="hero__backdrop"
              @error="failedBackdrops.add(movie.id)"
            />
            <div v-else class="hero__backdrop hero__backdrop--fallback">
              <span class="hero__fallback-title">{{ movie.title }}</span>
            </div>
            <div class="hero__overlay" />
          </div>
        </TransitionGroup>
      </div>

      <div class="hero__content" v-if="currentHero">
        <h2 class="hero__movie-title">{{ currentHero.title }}</h2>
        <p class="hero__meta">
          {{ currentHero.release_date?.slice(0, 4) }}
          <span class="hero__dot-sep">·</span>
          ⭐ {{ currentHero.vote_average?.toFixed(1) }}
        </p>
      </div>

      <div class="hero__dots">
        <button
          v-for="(_, i) in heroMovies"
          :key="i"
          class="hero__dot"
          :class="{ 'hero__dot--active': heroIndex === i }"
          @click="handleDotClick(i)"
          :aria-label="`Slide ${i + 1}`"
        />
      </div>
    </section>

    <div class="content-area">
      <section class="movie-section" v-if="hasGenres">
        <div class="section-header">
          <h2 class="section-title">
            <Sparkles :size="16" class="section-icon" />
            แนะนำสำหรับคุณ
          </h2>
        </div>

        <div class="movie-grid" v-if="recommendedMovies.length">
          <div
            v-for="(movie, i) in recommendedMovies"
            :key="`rec-${movie.id}`"
            class="poster-wrap"
            @mouseenter="onCardEnter(movie.id, i, 'rec', $event)"
            @mouseleave="onCardLeave(movie.id, 'rec')"
          >
            <RouterLink
              :to="{ name: 'movie-detail', params: { id: movie.id } }"
              class="poster-card"
            >
              <img
                v-if="movie.poster_path && !failedPosters.has(movie.id)"
                :src="`https://image.tmdb.org/t/p/w342${movie.poster_path}`"
                :alt="movie.title"
                loading="lazy"
                @error="failedPosters.add(movie.id)"
              />
              <div v-else class="poster-fallback">
                <span class="poster-fallback__title">{{ movie.title }}</span>
              </div>
            </RouterLink>

            <Transition name="popup">
              <div
                v-if="hoveredMovieKey === `rec-${movie.id}`"
                class="hover-popup"
                :class="getPopupPosition(i)"
                @mouseenter="onPopupEnter"
                @mouseleave="onPopupLeave(movie.id, 'rec')"
              >
                <PopupCard
                  :movie="movie"
                  :trailer="getTrailer(movie.id)"
                  :is-iframe-mounted="getState(movie.id).isIframeMounted.value"
                  :is-iframe-loaded="getState(movie.id).isIframeLoaded.value"
                  :show-skeleton="getState(movie.id).showSkeleton.value"
                  :show-fallback="getState(movie.id).showFallback.value"
                  @iframe-load="getState(movie.id).onIframeLoad()"
                />
              </div>
            </Transition>
          </div>
        </div>

        <div class="movie-grid" v-else-if="isLoadingRecommended">
          <div v-for="i in 10" :key="i" class="poster-skeleton" />
        </div>
      </section>

      <section class="movie-section">
        <div class="section-header">
          <h2 class="section-title">
            <Clapperboard :size="16" class="section-icon" />
            หนังมาใหม่
          </h2>
          <RouterLink :to="{ name: 'movies' }" class="see-all-btn">
            ดูทั้งหมด <ChevronRight :size="14" />
          </RouterLink>
        </div>

        <div class="movie-grid" v-if="nowPlayingMovies.length">
          <div
            v-for="(movie, i) in nowPlayingMovies"
            :key="`now-${movie.id}`"
            class="poster-wrap"
            @mouseenter="onCardEnter(movie.id, i, 'now', $event)"
            @mouseleave="onCardLeave(movie.id, 'now')"
          >
            <RouterLink :to="{ name: 'movie-detail', params: { id: movie.id } }" class="poster-card">
              <img
                v-if="movie.poster_path && !failedPosters.has(movie.id)"
                :src="`https://image.tmdb.org/t/p/w342${movie.poster_path}`"
                :alt="movie.title"
                loading="lazy"
                @error="failedPosters.add(movie.id)"
              />
              <div v-else class="poster-fallback">
                <span class="poster-fallback__title">{{ movie.title }}</span>
              </div>
            </RouterLink>

            <Transition name="popup">
              <div
                v-if="hoveredMovieKey === `now-${movie.id}`"
                class="hover-popup"
                :class="getPopupPosition(i)"
                @mouseenter="onPopupEnter"
                @mouseleave="onPopupLeave(movie.id, 'now')"
              >
                <PopupCard
                  :movie="movie"
                  :trailer="getTrailer(movie.id)"
                  :is-iframe-mounted="getState(movie.id).isIframeMounted.value"
                  :is-iframe-loaded="getState(movie.id).isIframeLoaded.value"
                  :show-skeleton="getState(movie.id).showSkeleton.value"
                  :show-fallback="getState(movie.id).showFallback.value"
                  @iframe-load="getState(movie.id).onIframeLoad()"
                />
              </div>
            </Transition>
          </div>
        </div>

        <div class="movie-grid" v-else-if="isLoadingNowPlaying">
          <div v-for="i in 10" :key="i" class="poster-skeleton" />
        </div>
      </section>
    </div>

    <Transition name="fade">
      <div
        class="loading-overlay"
        v-if="isLoadingNowPlaying && !nowPlayingMovies.length"
      >
        <div class="spinner" />
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, onUnmounted } from "vue"
  import { useQuery } from "@tanstack/vue-query"
  import { movieApi } from "@/api/api"
  import { useAuthStore } from "@/stores/auth"
  import { Sparkles, Clapperboard, ChevronRight } from "lucide-vue-next"
  import type { Movie } from "@/types"
  import PopupCard from "@/components/movie/PopupCard.vue"
  import {
    resolveTrailer,
    useTrailerPreview,
    type ResolvedTrailer,
  } from "@/composables/useTrailerPreview"

  // ── Auth / genre ───────────────────────────────────────────
  const authStore = useAuthStore()

  const hasGenres = computed(() => {
    const g = authStore.user?.favorite_genres
    if (!g || g === "null" || g === "skip") return false
    try {
      const parsed = JSON.parse(g)
      return Array.isArray(parsed) && parsed.length > 0
    } catch {
      return g.length > 0
    }
  })

  const genreQuery = computed(() => {
    if (!hasGenres.value) return ""
    const raw = authStore.user?.favorite_genres
    if (!raw) return ""
    try {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed)) return parsed.join(",")
    } catch {}
    return raw
  })

  // ── State ──────────────────────────────────────────────────
  const heroIndex = ref(0)
  const currentHero = computed(() => heroMovies.value[heroIndex.value] ?? null)
  const hoveredMovieKey = ref<string | null>(null)
  const hoveredMovieId = ref<number | null>(null)
  const failedPosters = ref(new Set<number>())
  const failedBackdrops = ref(new Set<number>())
  let heroTimer: ReturnType<typeof setInterval> | null = null

  // ── Per-card trailer state store ──────────────────────────
  const cardStates = new Map<number, ReturnType<typeof useTrailerPreview>>()
  const cardTrailers = new Map<number, ResolvedTrailer | null>()
  const videoCache = new Map<number, any[]>()

  let insidePopup = false
  let showTimer: ReturnType<typeof setTimeout> | null = null
  const SHOW_DELAY = 200

  // ── Queries ────────────────────────────────────────────────
  const { data: recommendedData, isLoading: isLoadingRecommended } = useQuery({
    queryKey: computed(() => ["movies-recommended", genreQuery.value]),
    queryFn: () => movieApi.getRecommended(genreQuery.value).then(r => r.data),
    enabled: hasGenres,
  })

  const { data: nowPlayingData, isLoading: isLoadingNowPlaying } = useQuery({
    queryKey: ["movies-now-playing-home"],
    queryFn: () => movieApi.getNowPlaying(1).then(r => r.data),
  })

  // ── Computed lists ─────────────────────────────────────────
  const recommendedMovies = computed<Movie[]>(
    () => recommendedData.value?.results?.slice(0, 10) ?? [],
  )

  const nowPlayingMovies = computed<Movie[]>(
    () => nowPlayingData.value?.results?.slice(0, 10) ?? [],
  )

  const heroMovies = computed<Movie[]>(() => {
    const source = recommendedMovies.value.length
      ? recommendedMovies.value
      : nowPlayingMovies.value
    return source.slice(0, 5)
  })

  // ── Per-card helpers ───────────────────────────────────────
  function getState(movieId: number): ReturnType<typeof useTrailerPreview> {
    if (!cardStates.has(movieId)) {
      cardStates.set(movieId, useTrailerPreview({ mountDelay: 500 }))
    }
    return cardStates.get(movieId)!
  }

  function getTrailer(movieId: number): ResolvedTrailer | null {
    return cardTrailers.get(movieId) ?? null
  }

  async function fetchAndCacheTrailer(movie: Movie) {
    if (cardTrailers.has(movie.id)) return
    cardTrailers.set(movie.id, null)

    try {
      const res = await movieApi.getVideos(movie.id)
      const videos = res.data?.results ?? []
      videoCache.set(movie.id, videos)
      const trailer = resolveTrailer(videos)
      cardTrailers.set(movie.id, trailer)

      if (hoveredMovieId.value === movie.id && trailer) {
        getState(movie.id).scheduleMount()
      }
    } catch {
      cardTrailers.set(movie.id, null)
    }
  }

  // ── Hover handlers ─────────────────────────────────────────
  function onCardEnter(
    movieId: number,
    index: number,
    section: "rec" | "now",
    event: MouseEvent,
  ) {
    clearTimeout(showTimer ?? undefined)
    insidePopup = false

    hoveredMovieId.value = movieId

    const movie = [...recommendedMovies.value, ...nowPlayingMovies.value].find(
      m => m.id === movieId,
    )
    if (movie) fetchAndCacheTrailer(movie)

    showTimer = setTimeout(() => {
      hoveredMovieKey.value = `${section}-${movieId}`

      const trailer = getTrailer(movieId)
      if (trailer) {
        getState(movieId).scheduleMount()
      }
    }, SHOW_DELAY)
  }

  function onCardLeave(movieId: number, section: "rec" | "now") {
    clearTimeout(showTimer ?? undefined)
    setTimeout(() => {
      if (!insidePopup) closeCard(movieId, section)
    }, 80)
  }

  function onPopupEnter() {
    insidePopup = true
  }

  function onPopupLeave(movieId: number, section: "rec" | "now") {
    insidePopup = false
    closeCard(movieId, section)
  }

  function closeCard(movieId: number, section: "rec" | "now") {
    const targetKey = `${section}-${movieId}`
    if (hoveredMovieKey.value !== targetKey) return

    hoveredMovieKey.value = null
    hoveredMovieId.value = null

    const state = cardStates.get(movieId)
    if (state) state.unmount()
  }

  // ── Popup position (same logic as original) ────────────────
  function getPopupPosition(index: number): string {
    const col = index % 5
    if (col === 0) return "popup--right"
    if (col === 4) return "popup--left"
    return "popup--center"
  }

  // ── Hero actions & timer ───────────────────────────────────
  function handleDotClick(index: number) {
    heroIndex.value = index
    resetTimer()
  }

  function resetTimer() {
    if (heroTimer) clearInterval(heroTimer)
    heroTimer = setInterval(() => {
      heroIndex.value =
        (heroIndex.value + 1) % Math.max(heroMovies.value.length, 1)
    }, 4000)
  }

  onMounted(() => resetTimer())
  onUnmounted(() => {
    if (heroTimer) clearInterval(heroTimer)
    clearTimeout(showTimer ?? undefined)
    cardStates.forEach(s => s.unmount())
    cardStates.clear()
  })
</script>

<style scoped>
  .home {
    background: #141414;
    min-height: 100vh;
  }

  /* ── Hero ─────────────────────────────────────────────────── */
  .hero {
    position: relative;
    width: 100%;
    height: 320px;
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
  .hero__backdrop--fallback {
    background: #1a1a1a;
    display: flex;
    align-items: flex-end;
    padding: 2rem;
  }
  .hero__fallback-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: #444;
  }
  .hero__overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(
      to bottom,
      rgba(0, 0, 0, 0.05) 0%,
      rgba(20, 20, 20, 0.85) 100%
    );
  }
  .hero__content {
    position: absolute;
    bottom: 48px;
    left: 50%;
    transform: translateX(-50%);
    width: 100%;
    max-width: 1100px;
    padding: 0 1.5rem;
    z-index: 2;
    pointer-events: none;
  }
  .hero__movie-title {
    font-size: 1.4rem;
    font-weight: 800;
    color: #fff;
    margin: 0 0 0.25rem;
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.6);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 480px;
  }
  .hero__meta {
    font-size: 0.8rem;
    color: rgba(255, 255, 255, 0.6);
    margin: 0;
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }
  .hero__dot-sep {
    opacity: 0.4;
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
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.3);
    border: none;
    cursor: pointer;
    transition:
      background 0.2s,
      transform 0.2s;
    padding: 0;
  }
  .hero__dot--active {
    background: #e50914;
    transform: scale(1.25);
  }
  .hero-slide-enter-active,
  .hero-slide-leave-active {
    transition: opacity 0.6s;
  }
  .hero-slide-enter-from,
  .hero-slide-leave-to {
    opacity: 0;
  }

  /* ── Content area ─────────────────────────────────────────── */
  .content-area {
    padding: 1.75rem 1.5rem 4rem;
    max-width: 1100px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 2.5rem;
  }

  /* ── Section header ───────────────────────────────────────── */
  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.875rem;
  }
  .section-title {
    display: flex;
    align-items: center;
    gap: 0.45rem;
    font-size: 1rem;
    font-weight: 700;
    color: #fff;
    margin: 0;
    letter-spacing: 0.2px;
  }
  .section-icon {
    color: #e50914;
  }
  .see-all-btn {
    display: flex;
    align-items: center;
    gap: 0.15rem;
    font-size: 0.8rem;
    color: #666;
    text-decoration: none;
    transition: color 0.2s;
  }
  .see-all-btn:hover {
    color: #fff;
  }

  /* ── Movie grid ───────────────────────────────────────────── */
  .movie-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 10px;
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
    background: #1a1a1a;
  }
  .poster-card:hover {
    transform: scale(1.04);
    box-shadow: 0 8px 28px rgba(0, 0, 0, 0.55);
  }
  .poster-card img {
    width: 100%;
    aspect-ratio: 2/3;
    object-fit: cover;
    display: block;
  }

  /* ── Fallback card ────────────────────────────────────────── */
  .poster-fallback {
    width: 100%;
    aspect-ratio: 2/3;
    background: #1a1a1a;
    border: 1px solid #252525;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.75rem;
  }
  .poster-fallback__title {
    font-size: 0.72rem;
    color: #484848;
    text-align: center;
    line-height: 1.45;
    display: -webkit-box;
    -webkit-line-clamp: 4;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  /* ── Skeleton ─────────────────────────────────────────────── */
  .poster-skeleton {
    aspect-ratio: 2/3;
    border-radius: 8px;
    background: linear-gradient(90deg, #1a1a1a 25%, #242424 50%, #1a1a1a 75%);
    background-size: 200% 100%;
    animation: shimmer 1.3s infinite;
  }
  @keyframes shimmer {
    0% {
      background-position: 200% 0;
    }
    100% {
      background-position: -200% 0;
    }
  }

  /* ── Hover popup ──────────────────────────────────────────── */
  .hover-popup {
    position: absolute;
    top: 0;
    z-index: 50;
    width: 300px;
    background: #1c1c1c;
    border: 1px solid rgba(255, 255, 255, 0.07);
    border-radius: 12px;
    overflow: hidden;
    box-shadow:
      0 20px 60px rgba(0, 0, 0, 0.8),
      0 0 0 1px rgba(255, 255, 255, 0.04);
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

  .popup-enter-active {
    transition:
      opacity 0.18s,
      transform 0.18s cubic-bezier(0.2, 0, 0.13, 1.35);
  }
  .popup-leave-active {
    transition: opacity 0.12s;
  }
  .popup-enter-from {
    opacity: 0;
    transform: scale(0.93) translateY(4px);
  }
  .popup--center.popup-enter-from {
    transform: translateX(-50%) scale(0.93) translateY(4px);
  }
  .popup-leave-to {
    opacity: 0;
  }

  /* ── Loading overlay ──────────────────────────────────────── */
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
    width: 38px;
    height: 38px;
    border: 2.5px solid rgba(255, 255, 255, 0.12);
    border-top-color: #e50914;
    border-radius: 50%;
    animation: spin 0.75s linear infinite;
  }
  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.25s;
  }
  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
  }
</style>
