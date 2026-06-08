<template>
  <div class="watchlist-wrap">
    <!-- ── หัวข้อรายการเฝ้าดู ── -->
    <div class="section-head">
      <span class="section-label">Watchlist</span>
      <span class="count-badge">{{ watchlist.length }}</span>
      <div class="section-rule" />
    </div>

    <!-- ── Loading State ── -->
    <div v-if="loading" class="state-loading">
      <div class="loading-bar"><div class="loading-fill" /></div>
    </div>

    <!-- ── Empty State ── -->
    <div v-else-if="watchlist.length === 0" class="state-empty">
      <Bookmark :size="32" :stroke-width="1.2" aria-hidden="true" />
      <p>Your watchlist is empty</p>
    </div>

    <!-- ── โปสเตอร์กริด (Poster Grid) ── -->
    <div v-else class="poster-grid">
      <div v-for="item in watchlist" :key="item.id" class="poster-card">
        <div class="poster-wrap" :aria-label="item.title">
          <!-- รูปภาพหน้าปก -->
          <img
            v-if="item.coverUrl"
            :src="item.coverUrl"
            :alt="item.title"
            class="poster-img"
          />
          <div v-else class="poster-fallback">
            <Film :size="20" :stroke-width="1.5" aria-hidden="true" />
          </div>

          <!-- แท็กหมวดหมู่มุมซ้ายบน -->
          <span class="category-tag">{{ item.category }}</span>

          <!-- ปุ่มลบออกจากรายการ (Hover แล้วขึ้น) -->
          <button
            class="remove-btn"
            :aria-label="`Remove ${item.title} from watchlist`"
            @click="handleRemove(item.id)"
          >
            <X :size="12" aria-hidden="true" />
          </button>

          <!-- ตัวเคลือบเงาการ์ดขณะโฮเวอร์ -->
          <div class="poster-overlay" aria-hidden="true">
            <p class="overlay-title">{{ item.title }}</p>
          </div>
        </div>

        <!-- รายละเอียดด้านล่างโปสเตอร์ -->
        <div class="poster-meta">
          <h4 class="poster-title">{{ item.title }}</h4>
          <span class="poster-date">
            <Clock :size="10" aria-hidden="true" />
            {{ item.addedAt }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref } from "vue"
  import { Bookmark, Film, X, Clock } from "lucide-vue-next"

  const props = defineProps<{ userId: number }>()

  const loading = ref(false)

  interface WatchlistItem {
    id: number
    title: string
    category: string
    coverUrl?: string
    addedAt: string
  }

  const watchlist = ref<WatchlistItem[]>([])

  onMounted(async () => {
    try {
      loading.value = true
      await new Promise(resolve => setTimeout(resolve, 400))
      watchlist.value = [
        {
          id: 101,
          title: "Interstellar",
          category: "SCI-FI",
          coverUrl:
            "https://image.tmdb.org/t/p/w300/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg",
          addedAt: "1 Jun 2026",
        },
        {
          id: 102,
          title: "Oppenheimer",
          category: "DRAMA",
          coverUrl:
            "https://image.tmdb.org/t/p/w300/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg",
          addedAt: "28 May 2026",
        },
        {
          id: 103,
          title: "Dune: Part Two",
          category: "SCI-FI",
          coverUrl:
            "https://image.tmdb.org/t/p/w300/1pdfLvkbY9ohJlCjQH2CZjjYVvJ.jpg",
          addedAt: "20 May 2026",
        },
      ]
    } catch (err) {
      console.error("Fetch watchlist failed:", err)
    } finally {
      loading.value = false
    }
  })

  function handleRemove(id: number) {
    watchlist.value = watchlist.value.filter(i => i.id !== id)
  }
</script>

<style scoped>
  .watchlist-wrap {
    /* ชุดตัวแปรสี Premium Crimson Red & Deep Dark */
    --bg-dark: #0a0a0a;
    --surface-card: #141414;
    --brand-red: #e50914;
    --brand-red-hover: #b20710;
    --border: rgba(255, 255, 255, 0.06);
    --border-hover: rgba(235, 9, 20, 0.2);
    --text: #f5f5f7;
    --sub: #8e8e93;
    --muted: #48484a;
    --ui:
      -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue",
      system-ui, sans-serif;

    color: var(--text);
    font-family: var(--ui);
  }

  .section-head {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }

  .section-label {
    font-size: 0.7rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--muted);
    white-space: nowrap;
  }

  .count-badge {
    font-size: 0.7rem;
    font-weight: 600;
    color: #ffffff;
    background: var(--brand-red);
    padding: 2px 8px;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(229, 9, 20, 0.3);
  }

  .section-rule {
    flex: 1;
    height: 1px;
    background: var(--border);
  }

  /* ── Loader States ── */
  .state-loading {
    padding: 4rem 0;
    display: flex;
    justify-content: center;
  }

  .loading-bar {
    width: 120px;
    height: 2px;
    background: var(--border);
    overflow: hidden;
    position: relative;
    border-radius: 2px;
  }

  .loading-fill {
    height: 100%;
    width: 40%;
    background: var(--brand-red);
    position: absolute;
    box-shadow: 0 0 8px var(--brand-red);
    animation: sweep 1.5s infinite ease-in-out;
  }

  @keyframes sweep {
    0% {
      left: -40%;
    }
    100% {
      left: 100%;
    }
  }

  .state-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    padding: 5rem 0;
    color: var(--muted);
  }

  .state-empty p {
    font-size: 0.85rem;
    margin: 0;
    color: var(--sub);
  }

  /* ── Poster Grid Layout ── */
  .poster-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
    gap: 1.25rem 1rem;
  }

  .poster-card {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
  }

  .poster-wrap {
    position: relative;
    aspect-ratio: 2/3;
    background: var(--surface-card);
    border: 1px solid var(--border);
    border-radius: 8px;
    overflow: hidden;
    cursor: pointer;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
    transition:
      transform 0.25s ease,
      border-color 0.25s ease,
      box-shadow 0.25s ease;
  }

  .poster-wrap:hover {
    transform: translateY(-4px);
    border-color: var(--border-hover);
    box-shadow: 0 8px 20px rgba(229, 9, 20, 0.15);
  }

  .poster-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  .poster-wrap:hover .poster-img {
    transform: scale(1.05);
  }

  .poster-fallback {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--muted);
  }

  /* Tags บนรูปภาพ */
  .category-tag {
    position: absolute;
    top: 8px;
    left: 8px;
    font-size: 0.55rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    background: rgba(0, 0, 0, 0.75);
    color: #ffffff;
    padding: 3px 6px;
    border-radius: 4px;
    backdrop-filter: blur(4px);
    border: 1px solid rgba(255, 255, 255, 0.05);
    z-index: 2;
  }

  /* ปุ่มลบ (X) ดีไซน์สไตล์วงกลมกระจกแบบโมเดิร์น */
  .remove-btn {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 24px;
    height: 24px;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(4px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 50%;
    color: #ffffff;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    opacity: 0;
    transform: scale(0.9);
    transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
    z-index: 3;
  }

  .poster-wrap:hover .remove-btn {
    opacity: 1;
    transform: scale(1);
  }

  .remove-btn:hover {
    background: var(--brand-red);
    border-color: var(--brand-red);
    box-shadow: 0 0 8px rgba(229, 9, 20, 0.6);
  }

  /* ไล่เฉดสีมืดด้านล่างการ์ด */
  .poster-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(
      to top,
      rgba(0, 0, 0, 0.95) 0%,
      rgba(0, 0, 0, 0.4) 40%,
      transparent 80%
    );
    opacity: 0;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    padding: 10px;
    transition: opacity 0.2s ease;
    z-index: 1;
  }

  .poster-wrap:hover .poster-overlay {
    opacity: 1;
  }

  .overlay-title {
    font-size: 0.65rem;
    font-weight: 500;
    color: #ffffff;
    line-height: 1.3;
    margin: 0;
  }

  /* Metadata ส่วนข้อมูลภายนอกการ์ด */
  .poster-meta {
    padding: 0 2px;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .poster-title {
    font-size: 0.8rem;
    font-weight: 500;
    color: var(--text);
    margin: 0;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
    cursor: pointer;
    transition: color 0.2s ease;
  }

  .poster-title:hover {
    color: var(--brand-red);
  }

  .poster-date {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 0.68rem;
    color: var(--sub);
  }
</style>
