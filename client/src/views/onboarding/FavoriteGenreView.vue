<template>
  <div class="onboarding-page">
    <div class="onboarding-bg"><div class="onboarding-bg__overlay" /></div>

    <div class="onboarding-container">
      <!-- Header -->
      <div class="onboarding-header">
        <div class="brand-logo">
          <span class="brand-movie">MOVIE</span
          ><span class="brand-hub">HUB</span>
        </div>
        <h1 class="onboarding-title">เลือกหนังที่คุณชอบ</h1>
        <p class="onboarding-sub">
          เลือกอย่างน้อย {{ MIN_PICKS }} เรื่องที่คุณชอบหรือเคยดูแล้ว
          เพื่อให้เราแนะนำหนังที่ใช่สำหรับคุณ
        </p>

        <!-- Progress bar -->
        <div class="progress-bar-wrap">
          <div
            class="progress-bar-fill"
            :style="{
              width: `${Math.min((selectedIds.length / MIN_PICKS) * 100, 100)}%`,
            }"
            :class="{
              'progress-bar-fill--done': selectedIds.length >= MIN_PICKS,
            }"
          />
        </div>
        <p class="progress-label">
          <span :class="{ 'count--done': selectedIds.length >= MIN_PICKS }">
            {{ selectedIds.length }}
          </span>
          / {{ MIN_PICKS }} เรื่องขึ้นไป
        </p>
      </div>

      <!-- Movie Grid -->
      <div class="movie-grid" v-if="movies.length">
        <button
          v-for="movie in movies"
          :key="movie.id"
          class="movie-card"
          :class="{ 'movie-card--selected': selectedIds.includes(movie.id) }"
          @click="toggleMovie(movie)"
        >
          <img
            v-if="movie.poster_path && !failedIds.has(movie.id)"
            :src="`https://image.tmdb.org/t/p/w185${movie.poster_path}`"
            :alt="movie.title"
            loading="lazy"
            class="movie-poster"
            @error="failedIds.add(movie.id)"
          />
          <div v-else class="movie-poster movie-poster--placeholder">
            <Film :size="28" />
            <span class="placeholder-title">{{ movie.title }}</span>
          </div>
          <!-- Overlay on selected -->
          <Transition name="check-fade">
            <div v-if="selectedIds.includes(movie.id)" class="movie-check">
              <Check :size="22" />
            </div>
          </Transition>
          <!-- Title tooltip on hover -->
          <div class="movie-title-bar">
            <span class="movie-title">{{ movie.title }}</span>
          </div>
        </button>
      </div>

      <!-- Skeleton -->
      <div class="movie-grid" v-else>
        <div
          v-for="i in TOTAL_MOVIES"
          :key="i"
          class="movie-card movie-card--skeleton"
        />
      </div>

      <!-- Footer -->
      <div class="onboarding-footer">
        <button
          class="btn-confirm"
          :disabled="selectedIds.length < MIN_PICKS || isSaving"
          @click="handleConfirm"
        >
          <span v-if="isSaving" class="spinner" />
          <span v-else>เริ่มต้นใช้งาน</span>
        </button>
        <button class="btn-skip" @click="handleSkip" :disabled="isSaving">
          ข้ามไปก่อน
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { movieApi, userApi } from "@/api/api";
import { useAuthStore } from "@/stores/auth";
import { Film, Check } from "lucide-vue-next";
import type { Movie } from "@/types";

const router = useRouter();
const authStore = useAuthStore();

const MIN_PICKS = 5;
const TOTAL_MOVIES = 20;

const movies = ref<Movie[]>([]);
const selectedIds = ref<number[]>([]);
const selectedMovies = ref<Movie[]>([]);
const isSaving = ref(false);
const failedIds = ref(new Set<number>());

function toggleMovie(movie: Movie) {
  const idx = selectedIds.value.indexOf(movie.id);
  if (idx !== -1) {
    selectedIds.value.splice(idx, 1);
    selectedMovies.value = selectedMovies.value.filter(
      (m) => m.id !== movie.id,
    );
  } else {
    selectedIds.value.push(movie.id);
    selectedMovies.value.push(movie);
  }
}

/** Extract unique genre IDs from selected movies, sorted by frequency */
function extractGenres(): number[] {
  const freq: Record<number, number> = {};
  for (const movie of selectedMovies.value) {
    for (const gid of movie.genre_ids ?? []) {
      freq[gid] = (freq[gid] ?? 0) + 1;
    }
  }
  return Object.entries(freq)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 5)
    .map(([id]) => Number(id));
}

async function handleConfirm() {
  if (selectedIds.value.length < MIN_PICKS || isSaving.value) return;
  isSaving.value = true;
  try {
    const genres = extractGenres(); // number[]
    await userApi.updateFavoriteGenres(authStore.user!.id, genres);
    // store ต้องเก็บเป็น JSON string ให้ hasGenres parse ได้
    authStore.setUser({
      ...authStore.user!,
      favorite_genres: JSON.stringify(genres),
    });
    router.push({ name: "home" });
  } catch {
    // handle error
  } finally {
    isSaving.value = false;
  }
}

async function handleSkip() {
  try {
    await userApi.updateProfile(authStore.user!.id, {
      favorite_genres: "skip",
    });
  } catch {
    /* ไม่ block user ถ้า api fail */
  }
  authStore.setUser({ ...authStore.user!, favorite_genres: "skip" });
  router.push({ name: "home" });
}

onMounted(async () => {
  try {
    // Fetch 2 pages of popular movies, then pick 20 randomly for variety
    const [p1, p2] = await Promise.all([
      movieApi.getPopular(1).then((r) => r.data.results),
      movieApi.getPopular(2).then((r) => r.data.results),
    ]);
    const pool: Movie[] = [...p1, ...p2].filter(
      (m) => !!m.poster_path,
    ) as Movie[];
    for (let i = pool.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      const temp = pool[i]!;
      pool[i] = pool[j]!;
      pool[j] = temp;
    }
    movies.value = pool.slice(0, TOTAL_MOVIES);
  } catch {
    // fallback: stay empty, skeleton will show
  }
});
</script>

<style scoped>
/* ── Page shell ───────────────────────────────────────── */
.onboarding-page {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  background: #0a0a0a;
  padding: 2rem 1rem 3rem;
}

.onboarding-bg {
  position: fixed;
  inset: 0;
  background:
    url("https://image.tmdb.org/t/p/w300/qJ2tW6WMUDux911r6m7haRef0WH.jpg") 0 0 /
      20% auto repeat-y,
    url("https://image.tmdb.org/t/p/w300/udDclJoHjfjb8Ekgsd4FDteOkCU.jpg") 20%
      0 / 20% auto repeat-y,
    url("https://image.tmdb.org/t/p/w300/rAiYTfKGqDCRIIqo664sY9XZIvQ.jpg") 40%
      0 / 20% auto repeat-y,
    url("https://image.tmdb.org/t/p/w300/b0PlSFdDwbyK0cf5RxwDpaOJQvQ.jpg") 60%
      0 / 20% auto repeat-y,
    url("https://image.tmdb.org/t/p/w300/74xTEgt7R36Fpooo50r9T25onhq.jpg") 80%
      0 / 20% auto repeat-y;
  filter: brightness(0.18) saturate(0.4);
  pointer-events: none;
}

.onboarding-bg__overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    rgba(10, 10, 10, 0.5) 0%,
    rgba(10, 10, 10, 0.8) 100%
  );
}

/* ── Container ────────────────────────────────────────── */
.onboarding-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 900px;
  display: flex;
  flex-direction: column;
  gap: 1.75rem;
}

/* ── Header ───────────────────────────────────────────── */
.onboarding-header {
  text-align: center;
}

.brand-logo {
  font-family: "Impact", "Arial Black", sans-serif;
  font-size: 2rem;
  font-weight: 900;
  margin-bottom: 1rem;
}
.brand-movie {
  color: #fff;
}
.brand-hub {
  color: #e50914;
}

.onboarding-title {
  font-size: 1.6rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 0.5rem;
}

.onboarding-sub {
  font-size: 0.875rem;
  color: #a3a3a3;
  margin: 0 0 1.25rem;
  line-height: 1.6;
}

/* Progress */
.progress-bar-wrap {
  width: 100%;
  max-width: 320px;
  margin: 0 auto 0.4rem;
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 9999px;
  overflow: hidden;
}
.progress-bar-fill {
  height: 100%;
  background: #e50914;
  border-radius: 9999px;
  transition:
    width 0.3s ease,
    background 0.3s;
}
.progress-bar-fill--done {
  background: #22c55e;
}
.progress-label {
  font-size: 0.8rem;
  color: #666;
  margin: 0;
}
.count--done {
  color: #22c55e;
  font-weight: 700;
}

/* ── Movie Grid ───────────────────────────────────────── */
.movie-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
}

@media (max-width: 680px) {
  .movie-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
@media (max-width: 480px) {
  .movie-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* ── Movie Card ───────────────────────────────────────── */
.movie-card {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  background: #1a1a1a;
  padding: 0;
  transition:
    border-color 0.2s,
    transform 0.2s,
    box-shadow 0.2s;
  aspect-ratio: 2/3;
  display: flex;
  flex-direction: column;
}

.movie-card:hover {
  transform: scale(1.04);
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.6);
  border-color: rgba(255, 255, 255, 0.25);
}

.movie-card--selected {
  border-color: #e50914 !important;
  box-shadow:
    0 0 0 1px #e50914,
    0 8px 28px rgba(229, 9, 20, 0.3) !important;
}

.movie-card--skeleton {
  background: linear-gradient(90deg, #1f1f1f 25%, #2a2a2a 50%, #1f1f1f 75%);
  background-size: 200% 100%;
  animation: shimmer 1.2s infinite;
  cursor: default;
  border-color: transparent !important;
  aspect-ratio: 2/3;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

/* Poster image */
.movie-poster {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  flex: 1;
}

.movie-poster--placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #1a1a1a;
  color: #333;
  flex: 1;
}

/* Checkmark overlay */
.movie-check {
  position: absolute;
  inset: 0;
  background: rgba(229, 9, 20, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  backdrop-filter: blur(1px);
}

.placeholder-title {
  font-size: 0.6rem;
  color: #555;
  text-align: center;
  padding: 0 0.4rem;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.check-fade-enter-active {
  transition: opacity 0.15s;
}
.check-fade-leave-active {
  transition: opacity 0.1s;
}
.check-fade-enter-from,
.check-fade-leave-to {
  opacity: 0;
}

/* Title bar (hover) */
.movie-title-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 1.5rem 0.5rem 0.4rem;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.85) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.2s;
  pointer-events: none;
}

.movie-card:hover .movie-title-bar {
  opacity: 1;
}

.movie-title {
  display: block;
  font-size: 0.65rem;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
}

/* ── Footer ───────────────────────────────────────────── */
.onboarding-footer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding-top: 0.25rem;
}

.btn-confirm {
  width: 100%;
  max-width: 320px;
  padding: 0.85rem;
  background: #e50914;
  color: #fff;
  border: none;
  border-radius: 9999px;
  font-family: "Impact", "Arial Black", sans-serif;
  font-size: 1.1rem;
  font-weight: 900;
  letter-spacing: 2px;
  cursor: pointer;
  transition:
    background 0.2s,
    transform 0.1s,
    opacity 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
}
.btn-confirm:hover:not(:disabled) {
  background: #f40612;
  transform: translateY(-1px);
}
.btn-confirm:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.btn-skip {
  background: none;
  border: none;
  color: #555;
  font-size: 0.82rem;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 3px;
  transition: color 0.2s;
}
.btn-skip:hover {
  color: #a3a3a3;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
