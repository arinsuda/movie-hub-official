<template>
  <div class="upcoming-page" ref="pageRef">
    <!-- Ambient backdrop -->
    <div class="ambient-bg">
      <div
        class="ambient-blob"
        :style="
          featuredMovie
            ? `background-image: url(https://image.tmdb.org/t/p/w780${featuredMovie.backdrop_path})`
            : ''
        "
      />
      <div class="ambient-overlay" />
    </div>

    <main class="page-main">
      <!-- Section title -->
      <div class="section-intro" ref="introRef">
        <div class="section-eyebrow">
          <span class="eyebrow-line" />
          <span class="eyebrow-text">COMING SOON</span>
          <span class="eyebrow-line" />
        </div>
        <h1 class="section-title">หนังที่กำลังจะมา</h1>
        <p class="section-sub">ติดตามภาพยนตร์ที่น่าจับตามองที่สุดแห่งปี</p>
      </div>

      <!-- HERO FEATURED FILM -->
      <section class="hero-section" ref="heroRef" v-if="featuredMovie">
        <div class="hero-backdrop">
          <img
            :src="`https://image.tmdb.org/t/p/w1280${featuredMovie.backdrop_path}`"
            :alt="featuredMovie.title"
            class="hero-backdrop-img"
          />
          <div class="hero-backdrop-grad" />
        </div>

        <div class="hero-content">
          <div class="hero-left">
            <div class="hero-poster-wrap">
              <img
                :src="`https://image.tmdb.org/t/p/w342${featuredMovie.poster_path}`"
                :alt="featuredMovie.title"
                class="hero-poster"
              />
              <div class="hero-poster-glow" />
            </div>
          </div>

          <div class="hero-right">
            <div class="hero-badge">Featured</div>
            <h2 class="hero-title">{{ featuredMovie.title }}</h2>
            <p class="hero-overview">
              {{ featuredMovie.overview?.slice(0, 200)
              }}{{ (featuredMovie.overview?.length ?? 0) > 200 ? "…" : "" }}
            </p>

            <div class="hero-meta">
              <span class="meta-chip">
                <CalendarDays :size="13" />
                {{ formatDate(featuredMovie.release_date) }}
              </span>
              <span
                class="meta-chip meta-chip--accent"
                v-if="daysUntil(featuredMovie.release_date) > 0"
              >
                <Clock :size="13" />
                อีก {{ daysUntil(featuredMovie.release_date) }} วัน
              </span>
              <span class="meta-chip meta-chip--green" v-else>
                🎬 เข้าฉายแล้ว
              </span>
            </div>

            <!-- Countdown -->
            <div class="countdown-wrap" v-if="countdown.length > 0">
              <div
                class="countdown-block"
                v-for="unit in countdown"
                :key="unit.label"
              >
                <span class="countdown-num">{{ unit.value }}</span>
                <span class="countdown-label">{{ unit.label }}</span>
              </div>
            </div>
            <div class="countdown-wrap" v-else-if="featuredMovie.release_date">
              <div class="release-now">🎬 เข้าฉายแล้ว!</div>
            </div>

            <button class="btn-detail" @click="goToDetail(featuredMovie.id)">
              <span>ดูรายละเอียด</span>
              <ArrowRight :size="16" />
            </button>
          </div>
        </div>
      </section>

      <!-- Hero Skeleton -->
      <section class="hero-section hero-skeleton" v-else-if="isLoading">
        <div class="skeleton-block hero-sk-backdrop" />
        <div class="hero-content">
          <div class="skeleton-block hero-sk-poster" />
          <div class="hero-right">
            <div class="skeleton-block sk-line sk-short" />
            <div
              class="skeleton-block sk-line sk-long"
              style="height: 2rem; margin-bottom: 0.75rem"
            />
            <div class="skeleton-block sk-line" />
            <div class="skeleton-block sk-line sk-med" />
          </div>
        </div>
      </section>

      <!-- Error state -->
      <div class="error-state" v-else-if="hasError">
        <Film :size="40" />
        <p>ไม่สามารถโหลดข้อมูลได้ กรุณาลองใหม่อีกครั้ง</p>
        <button class="btn-retry" @click="retryLoad">ลองใหม่</button>
      </div>

      <!-- MONTH TIMELINE GRID -->
      <section class="timeline-section" ref="timelineRef" v-if="movies.length">
        <template v-for="group in sortedMonthGroups" :key="group.month">
          <div class="month-group">
            <div class="month-label">
              <span class="month-line" />
              <span class="month-name">{{ group.month }}</span>
              <span class="month-count">{{ group.movies.length }} เรื่อง</span>
            </div>
            <div class="movie-rail">
              <button
                v-for="(movie, idx) in group.movies"
                :key="movie.id"
                class="rail-card"
                :class="{
                  'rail-card--featured': movie.id === featuredMovie?.id,
                }"
                :data-delay="idx"
                @click="setFeatured(movie)"
                @mouseenter="onRailHover($event, true)"
                @mouseleave="onRailHover($event, false)"
              >
                <div class="rail-poster-wrap">
                  <img
                    v-if="movie.poster_path"
                    :src="`https://image.tmdb.org/t/p/w185${movie.poster_path}`"
                    :alt="movie.title"
                    class="rail-poster"
                    loading="lazy"
                  />
                  <div v-else class="rail-poster rail-poster--placeholder">
                    <Film :size="20" />
                  </div>
                  <div class="rail-overlay">
                    <Play :size="18" fill="white" />
                  </div>
                </div>
                <div class="rail-info">
                  <span class="rail-title">{{ movie.title }}</span>
                  <span class="rail-date">{{
                    formatDateShort(movie.release_date)
                  }}</span>
                  <span
                    class="rail-days"
                    :class="{
                      'rail-days--soon':
                        daysUntil(movie.release_date) > 0 &&
                        daysUntil(movie.release_date) <= 14,
                      'rail-days--out': daysUntil(movie.release_date) <= 0,
                    }"
                  >
                    {{
                      daysUntil(movie.release_date) <= 0
                        ? "เข้าฉายแล้ว"
                        : `อีก ${daysUntil(movie.release_date)} วัน`
                    }}
                  </span>
                </div>
                <div
                  class="rail-active-dot"
                  v-if="movie.id === featuredMovie?.id"
                />
              </button>
            </div>
          </div>
        </template>

        <!-- Skeleton months -->
        <template v-if="!movies.length && isLoading">
          <div class="month-group" v-for="g in 2" :key="g">
            <div class="month-label">
              <span class="month-line" />
              <div
                class="skeleton-block sk-line sk-short"
                style="width: 100px; height: 1rem"
              />
            </div>
            <div class="movie-rail">
              <div
                class="skeleton-block rail-card-skeleton"
                v-for="i in 5"
                :key="i"
              />
            </div>
          </div>
        </template>
      </section>

      <!-- Load More -->
      <div
        class="load-more-wrap"
        ref="loadMoreRef"
        v-if="movies.length && hasMore"
      >
        <button class="btn-load-more" @click="loadMore" :disabled="isLoading">
          <span v-if="isLoading" class="spinner" />
          <span v-else>โหลดเพิ่มเติม</span>
        </button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, onUnmounted } from "vue"
  import { useRouter } from "vue-router"
  import { movieApi } from "@/api/api"
  import type { Movie } from "@/types"
  import { CalendarDays, Clock, ArrowRight, Film, Play } from "lucide-vue-next"
  import { gsap } from "gsap"
  import { ScrollTrigger } from "gsap/ScrollTrigger"

  gsap.registerPlugin(ScrollTrigger)

  const router = useRouter()

  // ── Refs ─────────────────────────────────────────────────────────
  const pageRef = ref<HTMLElement | null>(null)
  const introRef = ref<HTMLElement | null>(null)
  const heroRef = ref<HTMLElement | null>(null)
  const timelineRef = ref<HTMLElement | null>(null)
  const loadMoreRef = ref<HTMLElement | null>(null)

  // ── State ─────────────────────────────────────────────────────────
  const movies = ref<Movie[]>([])
  const currentPage = ref(1)
  const hasMore = ref(true)
  const isLoading = ref(false)
  const hasError = ref(false)
  const featuredMovie = ref<Movie | null>(null)

  // Countdown interval handle
  let countdownInterval: ReturnType<typeof setInterval> | null = null
  const now = ref(new Date())

  // ── Date helpers (timezone-safe) ──────────────────────────────────
  /**
   * Parse a YYYY-MM-DD string as a local date (not UTC).
   * Avoids the off-by-one caused by `new Date("2025-07-04")` being
   * interpreted as midnight UTC which is the previous day in UTC+7.
   */
  function parseLocalDate(dateStr: string): Date {
    if (!dateStr) return new Date(NaN)
    const [year, month, day] = dateStr.split("-").map(Number)
    return new Date(year!, month! - 1, day!)
  }

  function daysUntil(date: string): number {
    if (!date) return 0
    const release = parseLocalDate(date)
    const today = new Date(
      now.value.getFullYear(),
      now.value.getMonth(),
      now.value.getDate(),
    )
    const diff = release.getTime() - today.getTime()
    return Math.max(0, Math.ceil(diff / 86400000))
  }

  function formatDate(date: string): string {
    if (!date) return ""
    return parseLocalDate(date).toLocaleDateString("th-TH", {
      day: "numeric",
      month: "long",
      year: "numeric",
    })
  }

  function formatDateShort(date: string): string {
    if (!date) return ""
    return parseLocalDate(date).toLocaleDateString("th-TH", {
      day: "numeric",
      month: "short",
    })
  }

  // ── Computed ──────────────────────────────────────────────────────

  /** Group movies by month, keyed as "MMMM YYYY" (Thai locale) */
  const groupedByMonth = computed(() => {
    const map: Record<string, Movie[]> = {}
    for (const m of movies.value) {
      if (!m.release_date) continue
      const d = parseLocalDate(m.release_date)
      const key = d.toLocaleDateString("th-TH", {
        month: "long",
        year: "numeric",
      })
      if (!map[key]) map[key] = []
      map[key]!.push(m)
    }
    return map
  })

  /** Month groups sorted chronologically */
  const sortedMonthGroups = computed(() => {
    return Object.entries(groupedByMonth.value)
      .map(([month, moviesInMonth]) => ({ month, movies: moviesInMonth }))
      .sort((a, b) => {
        // Use the first movie's release date in each group for ordering
        const dateA = parseLocalDate(a.movies[0]?.release_date ?? "")
        const dateB = parseLocalDate(b.movies[0]?.release_date ?? "")
        return dateB.getTime() - dateB.getTime()
      })
  })

  const countdown = computed(() => {
    if (!featuredMovie.value?.release_date) return []
    const release = parseLocalDate(featuredMovie.value.release_date)
    // Use exact millisecond diff for countdown (hours/mins/secs need real time)
    const diff = release.getTime() - now.value.getTime()
    if (diff <= 0) return []
    const days = Math.floor(diff / 86400000)
    const hours = Math.floor((diff % 86400000) / 3600000)
    const mins = Math.floor((diff % 3600000) / 60000)
    const secs = Math.floor((diff % 60000) / 1000)
    return [
      { value: String(days).padStart(2, "0"), label: "วัน" },
      { value: String(hours).padStart(2, "0"), label: "ชั่วโมง" },
      { value: String(mins).padStart(2, "0"), label: "นาที" },
      { value: String(secs).padStart(2, "0"), label: "วินาที" },
    ]
  })

  // ── Helpers ───────────────────────────────────────────────────────
  function goToDetail(id: number) {
    router.push({ name: "movie-detail", params: { id } })
  }

  // ── Featured movie switch ─────────────────────────────────────────
  function setFeatured(movie: Movie) {
    if (movie.id === featuredMovie.value?.id) return

    const hero = heroRef.value
    gsap.to(hero, {
      opacity: 0,
      y: 10,
      duration: 0.2,
      ease: "power2.in",
      onComplete: () => {
        featuredMovie.value = movie
        gsap.fromTo(
          hero,
          { opacity: 0, y: -10 },
          { opacity: 1, y: 0, duration: 0.35, ease: "power3.out" },
        )
      },
    })
  }

  // ── Rail hover (GSAP) ─────────────────────────────────────────────
  function onRailHover(event: MouseEvent, entering: boolean) {
    const card = event.currentTarget as HTMLElement
    gsap.to(card, {
      y: entering ? -5 : 0,
      scale: entering ? 1.03 : 1,
      duration: entering ? 0.22 : 0.18,
      ease: "power2.out",
    })
  }

  // ── Pick best featured movie ──────────────────────────────────────
  /**
   * Pick the featured hero movie:
   * 1. Must have both backdrop_path and poster_path
   * 2. Among those, prefer the most popular upcoming movie
   *    (release date >= today, sorted by popularity desc)
   * 3. Fallback to any movie with images, then any movie
   */
  function pickFeatured(results: Movie[]): Movie | null {
    if (!results.length) return null

    const today = new Date(
      now.value.getFullYear(),
      now.value.getMonth(),
      now.value.getDate(),
    )

    const withImages = results.filter(m => m.backdrop_path && m.poster_path)

    // Prefer upcoming movies (not yet released)
    const upcoming = withImages.filter(
      m => m.release_date && parseLocalDate(m.release_date) >= today,
    )

    if (upcoming.length) {
      // Sort by popularity desc → pick most popular upcoming with images
      return (
        [...upcoming].sort(
          (a, b) => (b.popularity ?? 0) - (a.popularity ?? 0),
        )[0] ?? null
      )
    }

    // All are already released — just pick most popular with images
    if (withImages.length) {
      return (
        [...withImages].sort(
          (a, b) => (b.popularity ?? 0) - (a.popularity ?? 0),
        )[0] ?? null
      )
    }

    return results[0] ?? null
  }

  // ── Data fetching ─────────────────────────────────────────────────
  async function fetchMovies(page = 1) {
    isLoading.value = true
    hasError.value = false
    try {
      const res = await movieApi.getUpcoming(page)
      const results = res.data.results as Movie[]

      results.sort(
        (a, b) =>
          parseLocalDate(b.release_date).getTime() -
          parseLocalDate(a.release_date).getTime(),
      )

      console.log("Fetched upcoming movies:", results)

      if (page === 1) {
        movies.value = results
        featuredMovie.value = pickFeatured(results)
      } else {
        movies.value.push(...results)
      }

      hasMore.value = res.data.page < res.data.total_pages
      currentPage.value = page
    } catch {
      hasError.value = page === 1
    } finally {
      isLoading.value = false
    }
  }

  async function retryLoad() {
    await fetchMovies(1)
    await new Promise(r => setTimeout(r, 60))
    setupAnimations()
  }

  async function loadMore() {
    if (isLoading.value || !hasMore.value) return
    await fetchMovies(currentPage.value + 1)
    await new Promise(r => setTimeout(r, 50))
    const newCards = timelineRef.value?.querySelectorAll(
      ".rail-card:not(.gsap-revealed)",
    )
    if (newCards?.length) {
      gsap.fromTo(
        newCards,
        { opacity: 0, y: 30 },
        {
          opacity: 1,
          y: 0,
          duration: 0.45,
          stagger: 0.04,
          ease: "power3.out",
          onComplete: () =>
            newCards.forEach(c => c.classList.add("gsap-revealed")),
        },
      )
    }
  }

  // ── GSAP entrance + scroll triggers ──────────────────────────────
  function setupAnimations() {
    const tl = gsap.timeline({ defaults: { ease: "power3.out" } })

    tl.fromTo(
      introRef.value,
      { opacity: 0, y: 24 },
      { opacity: 1, y: 0, duration: 0.65 },
      "-=0.3",
    )

    tl.fromTo(
      heroRef.value,
      { opacity: 0, y: 40, scale: 0.97 },
      { opacity: 1, y: 0, scale: 1, duration: 0.75 },
      "-=0.3",
    )

    ScrollTrigger.batch(".month-group", {
      onEnter: batch => {
        gsap.fromTo(
          batch,
          { opacity: 0, y: 50 },
          {
            opacity: 1,
            y: 0,
            duration: 0.6,
            stagger: 0.12,
            ease: "power3.out",
            onComplete: () =>
              batch.forEach((el: Element) => el.classList.add("gsap-revealed")),
          },
        )
      },
      start: "top 88%",
      once: true,
    })

    ScrollTrigger.batch(".rail-card", {
      onEnter: batch => {
        gsap.fromTo(
          batch,
          { opacity: 0, x: -20 },
          {
            opacity: 1,
            x: 0,
            duration: 0.4,
            stagger: 0.05,
            ease: "power2.out",
            onComplete: () =>
              batch.forEach((el: Element) => el.classList.add("gsap-revealed")),
          },
        )
      },
      start: "top 92%",
      once: true,
    })
  }

  // ── Lifecycle ─────────────────────────────────────────────────────
  onMounted(async () => {
    gsap.set([introRef.value, heroRef.value], {
      opacity: 0,
    })

    // Countdown tick every second
    countdownInterval = setInterval(() => {
      now.value = new Date()
    }, 1000)

    await fetchMovies(1)
    await new Promise(r => setTimeout(r, 60))
    setupAnimations()
  })

  onUnmounted(() => {
    if (countdownInterval) clearInterval(countdownInterval)
    ScrollTrigger.getAll().forEach(t => t.kill())
  })
</script>

<style scoped>
  @import url("https://fonts.googleapis.com/css2?family=Bebas+Neue&family=DM+Sans:ital,wght@0,300;0,400;0,500;0,600;1,300&display=swap");

  /* ── Base ─────────────────────────────────────────────── */
  *,
  *::before,
  *::after {
    box-sizing: border-box;
  }

  .upcoming-page {
    --red: #e50914;
    --red-dim: rgba(229, 9, 20, 0.18);
    --gold: #f5c518;
    --green: #22c55e;
    --bg: #080808;
    --surface: #111111;
    --surface2: #1a1a1a;
    --border: rgba(255, 255, 255, 0.07);
    --text: #f0f0f0;
    --muted: #6b6b6b;
    --muted2: #9a9a9a;

    font-family: "Noto Sans Thai", sans-serif;
    background: var(--bg);
    color: var(--text);
    min-height: 100vh;
    position: relative;
    overflow-x: hidden;
  }

  /* ── Ambient bg ───────────────────────────────────────── */
  .ambient-bg {
    position: fixed;
    inset: 0;
    pointer-events: none;
    z-index: 0;
  }

  .ambient-blob {
    position: absolute;
    inset: 0;
    background-size: cover;
    background-position: center top;
    filter: blur(80px) brightness(0.12) saturate(1.6);
    transform: scale(1.1);
    transition: background-image 0.8s ease;
  }

  .ambient-overlay {
    position: absolute;
    inset: 0;
    background: radial-gradient(
      ellipse at 50% 0%,
      rgba(8, 8, 8, 0.4) 0%,
      #080808 70%
    );
  }

  /* ── Main ─────────────────────────────────────────────── */
  .page-main {
    position: relative;
    z-index: 1;
    max-width: 1100px;
    margin: 0 auto;
    padding: 3rem 1.5rem 5rem;
    display: flex;
    flex-direction: column;
    gap: 3.5rem;
  }

  /* ── Section Intro ────────────────────────────────────── */
  .section-intro {
    text-align: center;
  }

  .section-eyebrow {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    margin-bottom: 0.75rem;
  }

  .eyebrow-line {
    display: block;
    width: 40px;
    height: 1px;
    background: var(--red);
    opacity: 0.6;
  }

  .eyebrow-text {
    font-family: "Noto Sans Thai", sans-serif;
    font-size: 0.75rem;
    letter-spacing: 0.4em;
    color: var(--red);
  }

  .section-title {
    font-family: "Noto Sans Thai", sans-serif;
    font-size: clamp(2.4rem, 5vw, 3.5rem);
    letter-spacing: 0.04em;
    color: #fff;
    margin: 0 0 0.4rem;
    line-height: 1;
  }

  .section-sub {
    color: var(--muted2);
    font-size: 0.9rem;
    margin: 0;
  }

  /* ── Hero Section ─────────────────────────────────────── */
  .hero-section {
    position: relative;
    border-radius: 16px;
    overflow: hidden;
    border: 1px solid var(--border);
    background: var(--surface);
    min-height: 340px;
  }

  .hero-backdrop {
    position: absolute;
    inset: 0;
  }

  .hero-backdrop-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center 30%;
    display: block;
  }

  .hero-backdrop-grad {
    position: absolute;
    inset: 0;
    background:
      linear-gradient(
        to right,
        rgba(8, 8, 8, 0.98) 35%,
        rgba(8, 8, 8, 0.5) 70%,
        rgba(8, 8, 8, 0.15) 100%
      ),
      linear-gradient(to top, rgba(8, 8, 8, 0.7) 0%, transparent 50%);
  }

  .hero-content {
    position: relative;
    z-index: 1;
    display: flex;
    gap: 2.5rem;
    padding: 2.5rem;
    align-items: flex-start;
  }

  .hero-left {
    flex-shrink: 0;
  }

  .hero-poster-wrap {
    position: relative;
    width: 160px;
  }

  .hero-poster {
    width: 160px;
    border-radius: 10px;
    display: block;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.8);
  }

  .hero-poster-glow {
    position: absolute;
    inset: -4px;
    border-radius: 14px;
    background: linear-gradient(135deg, var(--red), transparent 60%);
    opacity: 0.4;
    z-index: -1;
    filter: blur(8px);
  }

  .hero-right {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 0.9rem;
  }

  .hero-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    background: var(--red-dim);
    border: 1px solid rgba(229, 9, 20, 0.3);
    color: #ff6b6b;
    font-size: 0.7rem;
    font-weight: 600;
    padding: 0.2rem 0.65rem;
    border-radius: 9999px;
    width: fit-content;
    letter-spacing: 0.05em;
    text-transform: uppercase;
  }

  .hero-title {
    font-family: "Noto Sans Thai", sans-serif;
    font-size: clamp(2rem, 4vw, 3rem);
    letter-spacing: 0.04em;
    color: #fff;
    margin: 0;
    line-height: 1;
  }

  .hero-overview {
    color: var(--muted2);
    font-size: 0.875rem;
    line-height: 1.7;
    margin: 0;
    font-weight: 300;
    max-width: 520px;
  }

  .hero-meta {
    display: flex;
    gap: 0.6rem;
    flex-wrap: wrap;
  }

  .meta-chip {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid var(--border);
    color: var(--muted2);
    font-size: 0.75rem;
    padding: 0.25rem 0.7rem;
    border-radius: 6px;
  }

  .meta-chip--accent {
    background: rgba(245, 197, 24, 0.1);
    border-color: rgba(245, 197, 24, 0.25);
    color: var(--gold);
  }

  .meta-chip--green {
    background: rgba(34, 197, 94, 0.1);
    border-color: rgba(34, 197, 94, 0.25);
    color: var(--green);
  }

  /* Countdown */
  .countdown-wrap {
    display: flex;
    gap: 0.6rem;
    flex-wrap: wrap;
  }

  .countdown-block {
    display: flex;
    flex-direction: column;
    align-items: center;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 0.5rem 0.9rem;
    min-width: 64px;
  }

  .countdown-num {
    font-family: "Noto Sans Thai", sans-serif;
    font-size: 2rem;
    line-height: 1;
    color: #fff;
    letter-spacing: 0.05em;
  }

  .countdown-label {
    font-size: 0.6rem;
    color: var(--muted);
    letter-spacing: 0.06em;
    text-transform: uppercase;
    margin-top: 2px;
  }

  .release-now {
    font-family: "Noto Sans Thai", sans-serif;
    font-size: 1.4rem;
    color: var(--green);
    letter-spacing: 0.06em;
  }

  .btn-detail {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--red);
    color: #fff;
    border: none;
    border-radius: 9999px;
    padding: 0.65rem 1.4rem;
    font-family: "Noto Sans Thai", sans-serif;
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    width: fit-content;
    transition:
      background 0.2s,
      transform 0.15s;
  }
  .btn-detail:hover {
    background: #f40612;
    transform: translateY(-1px);
  }

  /* ── Hero Skeleton ────────────────────────────────────── */
  .hero-skeleton {
    min-height: 340px;
  }

  .hero-sk-backdrop {
    position: absolute !important;
    inset: 0;
    border-radius: 0;
  }

  .hero-sk-poster {
    width: 160px !important;
    height: 240px !important;
    border-radius: 10px;
    flex-shrink: 0;
  }

  /* ── Error state ──────────────────────────────────────── */
  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 4rem 0;
    color: var(--muted2);
    text-align: center;
  }

  .btn-retry {
    background: var(--surface2);
    border: 1px solid var(--border);
    color: var(--text);
    font-family: "Noto Sans Thai", sans-serif;
    font-size: 0.85rem;
    padding: 0.5rem 1.5rem;
    border-radius: 9999px;
    cursor: pointer;
    transition: background 0.2s;
  }
  .btn-retry:hover {
    background: #222;
  }

  /* ── Skeletons ────────────────────────────────────────── */
  .skeleton-block {
    background: linear-gradient(90deg, #181818 25%, #222 50%, #181818 75%);
    background-size: 200% 100%;
    animation: shimmer 1.4s infinite;
    border-radius: 8px;
  }

  .sk-line {
    height: 0.85rem;
    width: 100%;
    margin-bottom: 0.5rem;
  }
  .sk-short {
    width: 40%;
  }
  .sk-med {
    width: 70%;
  }
  .sk-long {
    width: 90%;
  }

  .rail-card-skeleton {
    width: 120px;
    flex-shrink: 0;
    height: 220px;
    border-radius: 10px;
  }

  @keyframes shimmer {
    0% {
      background-position: 200% 0;
    }
    100% {
      background-position: -200% 0;
    }
  }

  /* ── Timeline Section ─────────────────────────────────── */
  .timeline-section {
    display: flex;
    flex-direction: column;
    gap: 2.5rem;
  }

  .month-group {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .month-label {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .month-line {
    flex: 1;
    height: 1px;
    background: var(--border);
  }

  .month-name {
    font-family: "Noto Sans Thai", sans-serif;
    font-size: 1.1rem;
    letter-spacing: 0.2em;
    color: var(--muted2);
    white-space: nowrap;
  }

  .month-count {
    font-size: 0.7rem;
    color: var(--muted);
    background: var(--surface2);
    border: 1px solid var(--border);
    padding: 0.15rem 0.5rem;
    border-radius: 9999px;
    white-space: nowrap;
  }

  /* ── Rail ─────────────────────────────────────────────── */
  .movie-rail {
    display: flex;
    gap: 12px;
    overflow-x: auto;
    padding-bottom: 0.5rem;
    scrollbar-width: thin;
    scrollbar-color: var(--surface2) transparent;
  }

  .movie-rail::-webkit-scrollbar {
    height: 4px;
  }
  .movie-rail::-webkit-scrollbar-track {
    background: transparent;
  }
  .movie-rail::-webkit-scrollbar-thumb {
    background: var(--surface2);
    border-radius: 9999px;
  }

  /* ── Rail Card ────────────────────────────────────────── */
  .rail-card {
    flex-shrink: 0;
    width: 130px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 10px;
    overflow: hidden;
    cursor: pointer;
    padding: 0;
    text-align: left;
    position: relative;
    will-change: transform;
    transition:
      border-color 0.2s,
      box-shadow 0.2s;
  }

  .rail-card--featured {
    border-color: var(--red) !important;
    box-shadow:
      0 0 0 1px var(--red),
      0 8px 24px rgba(229, 9, 20, 0.2) !important;
  }

  .rail-poster-wrap {
    position: relative;
    width: 100%;
    aspect-ratio: 2/3;
    overflow: hidden;
  }

  .rail-poster {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.3s ease;
  }

  .rail-card:hover .rail-poster {
    transform: scale(1.06);
  }

  .rail-poster--placeholder {
    background: var(--surface2);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--muted);
    height: 100%;
  }

  .rail-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.2s;
  }

  .rail-card:hover .rail-overlay {
    opacity: 1;
  }

  .rail-info {
    padding: 0.6rem 0.65rem 0.7rem;
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }

  .rail-title {
    font-size: 0.72rem;
    font-weight: 600;
    color: var(--text);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: block;
  }

  .rail-date {
    font-size: 0.65rem;
    color: var(--muted);
  }

  .rail-days {
    font-size: 0.62rem;
    color: var(--muted);
    font-style: italic;
  }

  .rail-days--soon {
    color: var(--gold);
    font-style: normal;
    font-weight: 600;
  }

  .rail-days--out {
    color: var(--green);
    font-style: normal;
    font-weight: 600;
  }

  .rail-active-dot {
    position: absolute;
    top: 7px;
    right: 7px;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--red);
    box-shadow: 0 0 0 2px rgba(229, 9, 20, 0.3);
    animation: pulse-dot 1.5s ease-in-out infinite;
  }

  @keyframes pulse-dot {
    0%,
    100% {
      box-shadow: 0 0 0 2px rgba(229, 9, 20, 0.3);
    }
    50% {
      box-shadow: 0 0 0 5px rgba(229, 9, 20, 0);
    }
  }

  /* ── Load More ────────────────────────────────────────── */
  .load-more-wrap {
    display: flex;
    justify-content: center;
  }

  .btn-load-more {
    background: transparent;
    border: 1px solid var(--border);
    color: var(--muted2);
    font-family: "Noto Sans Thai", sans-serif;
    font-size: 0.85rem;
    padding: 0.65rem 2rem;
    border-radius: 9999px;
    cursor: pointer;
    transition:
      border-color 0.2s,
      color 0.2s,
      background 0.2s;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 140px;
    justify-content: center;
  }

  .btn-load-more:hover:not(:disabled) {
    border-color: rgba(255, 255, 255, 0.2);
    color: #fff;
    background: rgba(255, 255, 255, 0.04);
  }

  .btn-load-more:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .spinner {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.2);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
    flex-shrink: 0;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  /* ── Responsive ───────────────────────────────────────── */
  @media (max-width: 720px) {
    .hero-content {
      flex-direction: column;
      gap: 1.5rem;
      padding: 1.5rem;
    }
    .hero-poster-wrap,
    .hero-poster {
      width: 120px;
    }
    .hero-sk-poster {
      width: 120px !important;
      height: 180px !important;
    }
    .countdown-block {
      min-width: 52px;
      padding: 0.4rem 0.6rem;
    }
    .countdown-num {
      font-size: 1.5rem;
    }
  }

  @media (max-width: 480px) {
    .page-main {
      padding: 2rem 1rem 4rem;
    }
    .section-title {
      font-size: 2rem;
    }
  }
</style>
