<template>
  <div class="movie-detail-page" ref="pageRef">
    <div class="detail-backdrop" v-if="movie">
      <img
        :src="`https://image.tmdb.org/t/p/w1280${movie.backdrop_path}`"
        :alt="movie.title"
        class="detail-backdrop__img"
      />
      <div class="detail-backdrop__overlay" />
    </div>

    <button class="btn-back" @click="goBack" ref="backBtnRef">
      <i class="pi pi-arrow-left"></i>
      <span>ย้อนกลับ</span>
    </button>

    <main class="detail-container" v-if="movie">
      <div class="detail-grid">
        <div class="detail-sidebar" ref="sidebarRef">
          <div class="poster-wrapper">
            <img
              v-if="movie.poster_path"
              :src="`https://image.tmdb.org/t/p/w500${movie.poster_path}`"
              :alt="movie.title"
              class="main-poster"
            />
            <div v-else class="main-poster main-poster--empty">
              <i class="pi pi-image" style="font-size: 3rem"></i>
            </div>
          </div>

          <div class="action-buttons">
            <button class="btn-primary-action" @click="showTrailerPopup = true">
              <i class="pi pi-play-circle"></i>
              <span>รับชมตัวอย่างหนัง</span>
            </button>

            <Teleport to="body">
              <PopupTrailer
                v-if="showTrailerPopup"
                :trailer-url="movieTrailerUrl"
                @close="showTrailerPopup = false"
              />
            </Teleport>

            <button
              :class="['btn-secondary-action', { 'active-like': isLiked }]"
              @click="toggleLike"
            >
              <i
                :class="isLiked ? 'pi pi-thumbs-up-fill' : 'pi pi-thumbs-up'"
              ></i>
              <span>{{ isLiked ? "ถูกใจแล้ว" : "ถูกใจหนังเรื่องนี้" }}</span>
            </button>

            <button
              :class="[
                'btn-secondary-action',
                { 'active-watchlist': isWatchlisted },
              ]"
              @click="toggleWatchlist"
            >
              <i
                :class="
                  isWatchlisted ? 'pi pi-bookmark-fill' : 'pi pi-bookmark'
                "
              ></i>
              <span>{{
                isWatchlisted ? "อยู่ในเพลย์ลิสต์" : "เพิ่มในรายการโปรด"
              }}</span>
            </button>

            <button
              :class="['btn-secondary-action', { 'active-watched': isWatched }]"
              @click="toggleWatched"
            >
              <i :class="isWatched ? 'pi pi-check-circle' : 'pi pi-circle'"></i>
              <span>{{
                isWatched ? "รับชมแล้ว" : "ทำเครื่องหมายว่าดูแล้ว"
              }}</span>
            </button>
          </div>
        </div>

        <div class="detail-content" ref="contentRef">
          <div class="content-header">
            <h1 class="movie-title">{{ movie.title }}</h1>
            <p class="movie-tagline" v-if="movie.tagline">
              "{{ movie.tagline }}"
            </p>
          </div>

          <div class="metadata-row">
            <span class="meta-badge" v-if="movie.release_date">
              <i class="pi pi-calendar"></i>
              {{ formatYear(movie.release_date) }}
            </span>
            <span class="meta-badge" v-if="movie.runtime">
              <i class="pi pi-clock"></i>
              {{ formatRuntime(movie.runtime) }}
            </span>
            <span class="meta-badge rating-tmdb" v-if="movie.vote_average">
              <i class="pi pi-star-fill"></i>
              TMDB: {{ movie.vote_average.toFixed(2) }}
            </span>
            <span class="meta-badge rating-remov">
              <i class="pi pi-heart-fill"></i>
              REMOV: {{ removStats.average_rating.toFixed(2) }}
            </span>
          </div>

          <div class="genres-list" v-if="movie.genres?.length">
            <span
              v-for="genre in movie.genres"
              :key="genre.id"
              class="genre-tag"
            >
              {{ genre.name }}
            </span>
          </div>

          <section class="remov-info-section">
            <div class="remov-stat-card view">
              <i class="pi pi-eye"></i>
              <div class="stat-details">
                <span class="stat-count">{{
                  removStats.view_count.toLocaleString()
                }}</span>
                <span class="stat-name">ยอดเข้าชม</span>
              </div>
            </div>
            <div class="remov-stat-card like">
              <i class="pi pi-thumbs-up-fill"></i>
              <div class="stat-details">
                <span class="stat-count">{{
                  removStats.like_count.toLocaleString()
                }}</span>
                <span class="stat-name">ถูกใจ</span>
              </div>
            </div>
            <div class="remov-stat-card review">
              <i class="pi pi-comment"></i>
              <div class="stat-details">
                <span class="stat-count">{{
                  removStats.review_count.toLocaleString()
                }}</span>
                <span class="stat-name">รีวิว</span>
              </div>
            </div>
            <div class="remov-stat-card watchlist">
              <i class="pi pi-bookmark-fill"></i>
              <div class="stat-details">
                <span class="stat-count">{{
                  removStats.watchlist_count.toLocaleString()
                }}</span>
                <span class="stat-name">เพลย์ลิสต์</span>
              </div>
            </div>
          </section>

          <hr class="section-divider" />

          <section class="info-section">
            <h2 class="section-title">
              <i class="pi pi-align-left"></i>เรื่องย่อ
            </h2>
            <p class="overview-text">
              {{
                movie.overview ||
                "ไม่มีข้อมูลเรื่องย่อภาษาไทยสำหรับภาพยนตร์เรื่องนี้"
              }}
            </p>
          </section>

          <section class="info-section grid-stats">
            <div class="stat-box" v-if="movie.status">
              <span class="stat-label">สถานะการฉาย</span>
              <span class="stat-value">{{ formatStatus(movie.status) }}</span>
            </div>
            <div class="stat-box" v-if="movie.budget">
              <span class="stat-label">ทุนสร้าง</span>
              <span class="stat-value">${{ formatMoney(movie.budget) }}</span>
            </div>
            <div class="stat-box" v-if="movie.revenue">
              <span class="stat-label">รายได้รวม</span>
              <span class="stat-value">${{ formatMoney(movie.revenue) }}</span>
            </div>
          </section>

          <section class="info-section" v-if="castList.length">
            <h2 class="section-title"><i class="pi pi-users"></i>นักแสดงนำ</h2>
            <div class="cast-scroll">
              <div v-for="actor in castList" :key="actor.id" class="cast-card">
                <div class="cast-avatar-wrap">
                  <img
                    v-if="actor.profile_path"
                    :src="`https://image.tmdb.org/t/p/w185${actor.profile_path}`"
                    :alt="actor.name"
                    class="cast-avatar"
                  />
                  <div v-else class="cast-avatar-empty">
                    <i class="pi pi-user"></i>
                  </div>
                </div>
                <div class="cast-info">
                  <span class="actor-name">{{ actor.name }}</span>
                  <span class="actor-character">{{ actor.character }}</span>
                </div>
              </div>
            </div>
          </section>

          <section class="info-section review-zone">
            <h2 class="section-title">
              <i class="pi pi-comments"></i>รีวิวจากผู้ใช้งาน REMOV
            </h2>
            <MovieReviews
              :movie-id="movieId"
              @review-submitted="onReviewSubmitted"
            />
          </section>
        </div>
      </div>
    </main>

    <div class="loading-state" v-if="isLoading">
      <i class="pi pi-spin pi-spinner spinner-icon"></i>
      <p>กำลังโหลดข้อมูลภาพยนตร์...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from "vue"
  import { useRoute, useRouter } from "vue-router"
  import { movieApi } from "@/api/api"
  import { libraryApi } from "@/api/endpoints/library"
  import { gsap } from "gsap"
  import PopupTrailer from "@/components/movie/PopupTrailer.vue"
  import MovieReviews from "@/components/movie/MovieReviews.vue"
  import { resolveTrailer } from "@/composables/useTrailerPreview"
  import { useAuthStore } from "@/stores/auth"
  import { mediaApi } from "@/api/endpoints/media"

  const route = useRoute()
  const router = useRouter()
  const movieId = computed(() => Number(route.params.id))

  const authStore = useAuthStore()
  const currentUserId = computed(() => authStore.user?.id ?? null)
  const showTrailerPopup = ref(false)
  const movieTrailerUrl = computed(() => {
    const videos = videoList.value
    const resolved = resolveTrailer(videos)
    return resolved ? resolved.embedUrl : null
  })

  // Animation Refs
  const pageRef = ref<HTMLElement | null>(null)
  const backBtnRef = ref<HTMLElement | null>(null)
  const sidebarRef = ref<HTMLElement | null>(null)
  const contentRef = ref<HTMLElement | null>(null)

  // States
  const movie = ref<any>(null)
  const castList = ref<any[]>([])
  const videoList = ref<any[]>([])
  const isLoading = ref(true)

  const removStats = ref({
    media_id: 0,
    media_type: "movie",
    like_count: 0,
    view_count: 0,
    review_count: 0,
    watchlist_count: 0,
    liked_at: null as string | null,
    average_rating: 0.0,
    has_rating: false,
  })

  const isLiked = computed(() => removStats.value.liked_at !== null)

  const isWatchlisted = ref(false)
  const isWatched = ref(false)

  const watchlist_ItemId = ref<number | null>(null)
  const watched_ItemId = ref<number | null>(null)

  function goBack() {
    router.back()
  }

  function formatYear(dateStr: string): string {
    if (!dateStr) return ""
    return dateStr.split("-")[0] || ""
  }

  function formatRuntime(minutes: number): string {
    if (!minutes) return ""
    const hours = Math.floor(minutes / 60)
    const mins = minutes % 60
    return hours > 0 ? `${hours} ชม. ${mins} นาที` : `${mins} นาที`
  }

  function formatMoney(amount: number): string {
    if (!amount || amount === 0) return "-"
    return amount.toLocaleString("en-US")
  }

  function formatStatus(status: string): string {
    const statusMap: Record<string, string> = {
      Released: "เข้าฉายแล้ว",
      "Post Production": "อยู่ในขั้นตอนหลังการถ่ายทำ",
      "In Production": "กำลังถ่ายทำ",
      Planned: "วางแผนการสร้าง",
    }
    return statusMap[status] || status
  }

  async function toggleLike() {
    try {
      if (isLiked.value) {
        await mediaApi.unlikeMedia("movie", movieId.value)
        removStats.value.liked_at = null
        removStats.value.like_count = Math.max(
          0,
          removStats.value.like_count - 1,
        )
        window.$toast?.info("ลบภาพยนตร์ออกจากรายการที่ชอบแล้ว")
      } else {
        await mediaApi.likeMedia("movie", movieId.value)
        removStats.value.liked_at = new Date().toISOString()
        removStats.value.like_count += 1
        window.$toast?.success("เพิ่มเข้าภาพยนตร์ที่คุณชื่นชอบแล้ว ❤️")
      }
    } catch (err) {
      console.error("ล้มเหลวในการบันทึกสถานะถูกใจ:", err)
      window.$toast?.error("เกิดข้อผิดพลาดในการบันทึกสถานะไลก์")
    }
  }

  async function toggleWatchlist() {
    if (!currentUserId.value) {
      window.$toast?.error("กรุณาเข้าสู่ระบบก่อนใช้งาน")
      return
    }

    try {
      if (isWatchlisted.value) {
        if (!watchlist_ItemId.value) {
          const libRes = await libraryApi.getMediaStatus(
            currentUserId.value,
            movieId.value,
            "movie",
          )
          const inLists = libRes.data?.in_lists || []
          const watchlistInfo = inLists.find(
            (l: any) => l.list_type === "watchlist",
          )

          if (watchlistInfo) {
            watchlist_ItemId.value = watchlistInfo.item_id
          } else {
            isWatchlisted.value = false
            return
          }
        }

        await libraryApi.removeItem(
          currentUserId.value,
          watchlist_ItemId.value!,
        )

        removStats.value.watchlist_count = Math.max(
          0,
          removStats.value.watchlist_count - 1,
        )
        isWatchlisted.value = false
        watchlist_ItemId.value = null
        window.$toast?.info("ลบออกจากเพลย์ลิสต์แล้ว")
      } else {
        const res = await libraryApi.addItem(currentUserId.value, {
          media_id: movieId.value,
          media_type: "movie",
          list_type: "watchlist",
        })

        if (res.data?.item) {
          watchlist_ItemId.value = res.data.item.id
        }

        removStats.value.watchlist_count += 1
        isWatchlisted.value = true
        window.$toast?.success("เพิ่มเข้าเพลย์ลิสต์สำเร็จ 🍿")
      }
    } catch (err) {
      console.error("ล้มเหลวในการแก้ไขสถานะ Watchlist:", err)
      window.$toast?.error("ไม่สามารถบันทึกข้อมูลเพลย์ลิสต์ได้")
    }
  }

  async function toggleWatched() {
    if (!currentUserId.value) {
      window.$toast?.error("กรุณาเข้าสู่ระบบก่อนใช้งาน")
      return
    }

    try {
      if (isWatched.value && watched_ItemId.value) {
        await libraryApi.removeItem(currentUserId.value, watched_ItemId.value)

        isWatched.value = false
        watched_ItemId.value = null
        window.$toast?.info("เปลี่ยนสถานะเป็นยังไม่ได้จัดส่งรับชม")
      } else {
        const res = await libraryApi.addItem(currentUserId.value, {
          media_id: movieId.value,
          media_type: "movie",
          list_type: "watched",
          watched_at: new Date().toISOString().split("T")[0],
        })

        if (res.data?.item) {
          watched_ItemId.value = res.data.item.id
        }

        isWatched.value = true
        window.$toast?.success("มาร์กสถานะว่ารับชมหนังเรื่องนี้แล้ว! ✔️")
      }
    } catch (err) {
      console.error("ล้มเหลวในการแก้ไขสถานะ Watched:", err)
      window.$toast?.error("ไม่สามารถบันทึกสถานะการรับชมได้")
    }
  }

  // Entrance GSAP Animation
  function animateEntrance() {
    if (!backBtnRef.value || !sidebarRef.value || !contentRef.value) return

    const tl = gsap.timeline({ defaults: { ease: "power3.out" } })

    tl.fromTo(
      backBtnRef.value,
      { opacity: 0, x: -20 },
      { opacity: 1, x: 0, duration: 0.4 },
    )
    tl.fromTo(
      sidebarRef.value,
      { opacity: 0, y: 30, scale: 0.95 },
      { opacity: 1, y: 0, scale: 1, duration: 0.6 },
      "-=0.2",
    )

    const contentElements = Array.from(contentRef.value.children)
    if (contentElements.length > 0) {
      tl.fromTo(
        contentElements,
        { opacity: 0, y: 20 },
        { opacity: 1, y: 0, duration: 0.5, stagger: 0.06 },
        "-=0.4",
      )
    }
  }

  async function fetchStats() {
    try {
      const statsRes = await mediaApi.getMediaStats("movie", movieId.value)

      if (statsRes.data && statsRes.data.stats) {
        const incomingStats = statsRes.data.stats

        removStats.value = {
          media_id: incomingStats.media_id ?? 0,
          media_type: incomingStats.media_type ?? "movie",
          like_count: incomingStats.like_count ?? 0,
          view_count: incomingStats.view_count ?? 0,
          review_count: incomingStats.review_count ?? 0,
          watchlist_count: incomingStats.watchlist_count ?? 0,
          liked_at:
            incomingStats.liked_at !== undefined
              ? incomingStats.liked_at
              : null,
          average_rating: incomingStats.average_rating ?? 0.0,
          has_rating: incomingStats.has_rating ?? false,
        }

        if (incomingStats.watchlisted_at) {
          isWatchlisted.value = true
        }
      }
    } catch (err) {
      console.error("ไม่สามารถดึงข้อมูลสถิติของภาพยนตร์ได้:", err)
    }
  }

  function onReviewSubmitted() {
    fetchStats()
  }

  onMounted(async () => {
    if (!movieId.value) {
      router.push({ name: "upcoming" })
      return
    }

    try {
      isLoading.value = true

      const res = await movieApi.getById(movieId.value)
      movie.value = res.data.movie
      castList.value = res.data.credits?.cast?.slice(0, 8) || []
      videoList.value = res.data.videos || []

      mediaApi.recordMediaView("movie", movieId.value).catch(err => {
        console.warn("record view failed:", err)
      })

      await fetchStats()

      if (currentUserId.value) {
        try {
          const libRes = await libraryApi.getMediaStatus(
            currentUserId.value,
            movieId.value,
            "movie",
          )
          const inLists = libRes.data?.in_lists || []

          const watchlistInfo = inLists.find(l => l.list_type === "watchlist")
          if (watchlistInfo) {
            isWatchlisted.value = true
            watchlist_ItemId.value = watchlistInfo.item_id
          }

          const watchedInfo = inLists.find(l => l.list_type === "watched")
          if (watchedInfo) {
            isWatched.value = true
            watched_ItemId.value = watchedInfo.item_id
          }
        } catch (libErr) {
          console.error(
            "ไม่สามารถดึงข้อมูลลิสต์สถานะของผู้ใช้จากห้องสมุดได้:",
            libErr,
          )
        }
      }

      isLoading.value = false
      setTimeout(() => {
        animateEntrance()
      }, 50)
    } catch (err) {
      console.error("Error fetching movie details:", err)
      isLoading.value = false
    }
  })
</script>

<style scoped>
  @import url("https://fonts.googleapis.com/css2?family=Noto+Sans+Thai:wght@300;400;500;600;700&display=swap");

  .movie-detail-page {
    --red: #e50914;
    --gold: #f5c518;
    --neon-pink: #ff2a74;
    --bg: #080808;
    --surface: #121212;
    --surface2: #1c1c1c;
    --border: rgba(255, 255, 255, 0.08);
    --text: #f0f0f0;
    --muted: #8a8a8a;

    font-family: "Noto Sans Thai", sans-serif;
    background: var(--bg);
    color: var(--text);
    min-height: 100vh;
    position: relative;
    overflow-x: hidden;
  }

  .btn-secondary-action.active-like {
    background: rgba(255, 42, 116, 0.2) !important;
    border-color: var(--neon-pink) !important;
    color: var(--neon-pink) !important;
  }

  .btn-secondary-action.active-watchlist {
    background: rgba(245, 197, 24, 0.15) !important;
    border-color: var(--gold) !important;
    color: var(--gold) !important;
  }

  .btn-secondary-action.active-watched {
    background: rgba(46, 213, 115, 0.15) !important;
    border-color: #2ed573 !important;
    color: #2ed573 !important;
  }

  .detail-backdrop {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 550px;
    z-index: 0;
    pointer-events: none;
  }
  .detail-backdrop__img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center 20%;
    opacity: 0.22;
    filter: blur(4px);
  }
  .detail-backdrop__overlay {
    position: absolute;
    inset: 0;
    background:
      linear-gradient(
        to top,
        var(--bg) 0%,
        rgba(8, 8, 8, 0.7) 60%,
        transparent 100%
      ),
      linear-gradient(
        to right,
        var(--bg) 0%,
        transparent 30%,
        transparent 70%,
        var(--bg) 100%
      );
  }
  .btn-back {
    position: absolute;
    top: 2rem;
    left: 2rem;
    z-index: 10;
    background: rgba(0, 0, 0, 0.5);
    border: 1px solid var(--border);
    color: #fff;
    padding: 0.55rem 1.1rem;
    border-radius: 9999px;
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.85rem;
    cursor: pointer;
    backdrop-filter: blur(8px);
    transition:
      background 0.2s,
      border-color 0.2s;
  }
  .btn-back:hover {
    background: var(--surface2);
    border-color: rgba(255, 255, 255, 0.2);
  }
  .detail-container {
    position: relative;
    z-index: 2;
    max-width: 1120px;
    margin: 0 auto;
    padding: 7.5rem 1.75rem 6rem;
  }
  .detail-grid {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 3.5rem;
    align-items: flex-start;
  }
  .poster-wrapper {
    position: relative;
    border-radius: 14px;
    overflow: hidden;
    border: 1px solid var(--border);
    box-shadow: 0 32px 64px rgba(0, 0, 0, 0.8);
    aspect-ratio: 2/3;
    background: var(--surface);
  }
  .main-poster {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }
  .main-poster--empty {
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--muted);
  }
  .action-buttons {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 1.5rem;
  }
  .btn-primary-action {
    background: var(--red);
    color: #fff;
    border: none;
    font-weight: 600;
    font-size: 0.875rem;
    padding: 0.75rem;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    transition: background 0.2s;
  }
  .btn-primary-action:hover {
    background: #f40612;
  }
  .btn-secondary-action {
    background: rgba(255, 255, 255, 0.05);
    color: #fff;
    border: 1px solid var(--border);
    font-weight: 500;
    font-size: 0.875rem;
    padding: 0.75rem;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    transition:
      background 0.2s,
      border-color 0.2s;
  }
  .btn-secondary-action:hover {
    background: rgba(255, 255, 255, 0.09);
    border-color: rgba(255, 255, 255, 0.2);
  }
  .detail-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  .movie-title {
    font-size: clamp(2.2rem, 4vw, 3.4rem);
    font-weight: 700;
    color: #fff;
    margin: 0;
    line-height: 1.15;
  }
  .movie-tagline {
    font-size: 1.05rem;
    color: var(--gold);
    font-style: italic;
    font-weight: 300;
    margin: -0.5rem 0 0;
  }
  .metadata-row {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-top: 0.25rem;
  }
  .meta-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.8rem;
    background: var(--surface);
    border: 1px solid var(--border);
    padding: 0.35rem 0.8rem;
    border-radius: 6px;
    color: #fff;
  }
  .meta-badge i {
    font-size: 0.75rem;
    color: var(--muted);
  }
  .rating-tmdb {
    background: rgba(245, 197, 24, 0.1);
    border-color: rgba(245, 197, 24, 0.25);
    color: var(--gold);
    font-weight: 600;
  }
  .rating-tmdb i {
    color: var(--gold);
  }
  .rating-remov {
    background: rgba(255, 42, 42, 0.15);
    border-color: rgba(255, 0, 0, 0.3);
    color: var(--red);
    font-weight: 600;
  }
  .rating-remov i {
    color: var(--red);
  }
  .genres-list {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }
  .genre-tag {
    background: transparent;
    border: 1px solid rgba(255, 255, 255, 0.15);
    padding: 0.2rem 0.75rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    color: #ddd;
  }
  .section-divider {
    border: none;
    border-top: 1px solid var(--border);
    margin: 0.5rem 0;
  }
  .remov-info-section {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin: 0.5rem 0;
  }
  .remov-stat-card {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 0.75rem 1rem;
    display: flex;
    align-items: center;
    gap: 0.85rem;
    transition:
      transform 0.2s,
      background 0.2s;
  }
  .remov-stat-card:hover {
    background: rgba(255, 255, 255, 0.06);
    transform: translateY(-2px);
  }
  .remov-stat-card i {
    font-size: 1.25rem;
    padding: 0.5rem;
    border-radius: 8px;
  }
  .stat-details {
    display: flex;
    flex-direction: column;
  }
  .stat-count {
    font-size: 1.05rem;
    font-weight: 700;
    color: #fff;
  }
  .stat-name {
    font-size: 0.68rem;
    color: var(--muted);
  }
  .remov-stat-card.view i {
    background: rgba(115, 164, 255, 0.15);
    color: #73a4ff;
  }
  .remov-stat-card.like i {
    background: rgba(255, 42, 116, 0.15);
    color: var(--neon-pink);
  }
  .remov-stat-card.review i {
    background: rgba(46, 213, 115, 0.15);
    color: #2ed573;
  }
  .remov-stat-card.watchlist i {
    background: rgba(245, 197, 24, 0.15);
    color: var(--gold);
  }
  .info-section {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
  }
  .section-title {
    font-size: 1.1rem;
    font-weight: 600;
    color: #fff;
    margin: 0;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .section-title i {
    color: var(--red);
    font-size: 0.95rem;
  }
  .overview-text {
    font-size: 0.925rem;
    line-height: 1.8;
    color: #ccc;
    margin: 0;
    font-weight: 300;
  }
  .grid-stats {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
    margin: 0.5rem 0;
  }
  .stat-box {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 0.8rem 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  .stat-label {
    font-size: 0.7rem;
    color: var(--muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  .stat-value {
    font-size: 0.95rem;
    font-weight: 600;
    color: #fff;
  }
  .cast-scroll {
    display: flex;
    gap: 1rem;
    overflow-x: auto;
    padding-bottom: 1rem;
    scrollbar-width: thin;
    scrollbar-color: var(--surface2) transparent;
  }
  .cast-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 8px;
    min-width: 125px;
    width: 125px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    transition: transform 0.2s;
  }
  .cast-card:hover {
    transform: translateY(-3px);
  }
  .cast-avatar-wrap {
    width: 100%;
    aspect-ratio: 1/1.2;
    background: var(--surface2);
    overflow: hidden;
  }
  .cast-avatar {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .cast-avatar-empty {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--muted);
    font-size: 1.5rem;
  }
  .cast-info {
    padding: 0.45rem 0.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }
  .actor-name {
    font-size: 0.75rem;
    font-weight: 600;
    color: #fff;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .actor-character {
    font-size: 0.65rem;
    color: var(--muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .review-zone {
    margin-top: 1rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border);
  }
  .loading-state {
    position: fixed;
    inset: 0;
    background: var(--bg);
    z-index: 100;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    color: var(--muted);
  }
  .spinner-icon {
    font-size: 2rem;
    color: var(--red);
  }

  @media (max-width: 820px) {
    .detail-grid {
      grid-template-columns: 1fr;
      gap: 2.5rem;
    }
    .detail-sidebar {
      display: grid;
      grid-template-columns: 180px 1fr;
      gap: 1.5rem;
      align-items: center;
    }
    .action-buttons {
      margin-top: 0;
    }
    .remov-info-section {
      grid-template-columns: repeat(2, 1fr);
    }
    .grid-stats {
      grid-template-columns: 1fr;
      gap: 0.75rem;
    }
    .btn-back {
      top: 1.25rem;
      left: 1.25rem;
      padding: 0.45rem 0.9rem;
    }
    .detail-container {
      padding-top: 5.5rem;
    }
  }

  @media (max-width: 480px) {
    .detail-sidebar {
      grid-template-columns: 1fr;
    }
    .poster-wrapper {
      max-width: 200px;
      margin: 0 auto;
      width: 100%;
    }
    .remov-info-section {
      grid-template-columns: 1fr;
      gap: 0.75rem;
    }
  }
</style>
