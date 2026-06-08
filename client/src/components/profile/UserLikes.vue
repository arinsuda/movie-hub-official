<template>
  <div class="likes-wrap">
    <!-- ── หัวข้อรายการที่ชอบ ── -->
    <div class="section-head">
      <span class="section-label">Liked Items</span>
      <span class="count-badge">{{ likedItems.length }}</span>
      <div class="section-rule" />
    </div>

    <!-- ── Loading State ── -->
    <div v-if="loading" class="state-loading">
      <div class="loading-bar"><div class="loading-fill" /></div>
    </div>

    <!-- ── Empty State ── -->
    <div v-else-if="likedItems.length === 0" class="state-empty">
      <Heart :size="32" :stroke-width="1.2" aria-hidden="true" />
      <p>No liked items yet</p>
    </div>

    <!-- ── รายการที่ชอบ (List Layout) ── -->
    <div v-else class="likes-list">
      <div v-for="item in likedItems" :key="item.id" class="like-card">
        <!-- กลุ่มข้อมูลซ้ายมือ -->
        <div class="like-main">
          <!-- ไอคอนหัวใจดวงเล็กสีแดงเนียนตา -->
          <div class="heart-icon-box" aria-hidden="true">
            <Heart :size="16" class="heart-icon" />
          </div>

          <!-- ข้อความและรายละเอียด -->
          <div class="like-info">
            <h4 class="like-title">{{ item.title }}</h4>
            <p class="like-meta">
              by {{ item.author }} • Liked on {{ item.likedAt }}
            </p>
          </div>
        </div>

        <!-- ปุ่ม Action ขวามือ -->
        <button
          class="unlike-btn"
          title="Remove from likes"
          @click="handleUnlike(item.id)"
        >
          Unlike
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref } from "vue"
  import { Heart } from "lucide-vue-next"

  const props = defineProps<{
    userId: number
  }>()

  const loading = ref(false)

  interface LikedItem {
    id: number
    title: string
    author: string
    likedAt: string
  }

  const likedItems = ref<LikedItem[]>([])

  onMounted(async () => {
    try {
      loading.value = true
      await new Promise(resolve => setTimeout(resolve, 400))

      likedItems.value = [
        {
          id: 201,
          title: "How to Optimize SQL Queries in MySQL for Production",
          author: "Tech Blogger",
          likedAt: "3 Jun 2026",
        },
        {
          id: 202,
          title: "State Management in Flutter: Comprehensive Guide to Riverpod",
          author: "Dev community",
          likedAt: "29 May 2026",
        },
        {
          id: 203,
          title: "Designing Scalable ERP System Database Models",
          author: "Database Architect",
          likedAt: "12 May 2026",
        },
      ]
    } catch (err) {
      console.error("Fetch liked items failed:", err)
    } finally {
      loading.value = false
    }
  })

  const handleUnlike = (id: number) => {
    likedItems.value = likedItems.value.filter(item => item.id !== id)
    console.log("Unliked item ID:", id)
  }
</script>

<style scoped>
  .likes-wrap {
    /* Token ระบบ แดง - ดำ คริมนอน */
    --bg-dark: #0a0a0a;
    --surface-card: #141414;
    --brand-red: #e50914;
    --brand-red-dim: rgba(229, 9, 20, 0.1);
    --border: rgba(255, 255, 255, 0.06);
    --border-hover: rgba(235, 9, 20, 0.15);
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
    box-shadow: 0 2px 8px rgba(229, 9, 20, 0.25);
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

  /* ── List Layout & Cards ── */
  .likes-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .like-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1.5rem;
    padding: 1rem 1.25rem;
    background: #111111;
    border: 1px solid var(--border);
    border-radius: 8px;
    transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .like-card:hover {
    background: var(--surface-card);
    border-color: var(--border-hover);
    transform: translateX(2px);
  }

  .like-main {
    display: flex;
    align-items: center;
    gap: 1rem;
    min-width: 0; /* แก้ไขปัญหา Text Truncate บั๊กใน Flexbox */
  }

  /* กล่องไอคอนรูปหัวใจ */
  .heart-icon-box {
    width: 36px;
    height: 36px;
    background: var(--brand-red-dim);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .heart-icon {
    color: var(--brand-red);
    fill: var(--brand-red); /* ถมสีให้ดูโดดเด่นสไตล์ Like */
  }

  /* ข้อมูล Content */
  .like-info {
    display: flex;
    flex-direction: column;
    gap: 3px;
    min-width: 0;
  }

  .like-title {
    font-size: 0.88rem;
    font-weight: 500;
    color: #ffffff;
    margin: 0;
    line-height: 1.3;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    transition: color 0.2s ease;
  }

  .like-card:hover .like-title {
    color: var(--brand-red);
  }

  .like-meta {
    font-size: 0.72rem;
    color: var(--sub);
    margin: 0;
  }

  /* ปุ่ม Unlike แก้วกระจกแดงมินิมอล */
  .unlike-btn {
    font-family: var(--ui);
    font-size: 0.72rem;
    font-weight: 500;
    color: var(--sub);
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid var(--border);
    padding: 0.45rem 0.85rem;
    border-radius: 6px;
    cursor: pointer;
    white-space: nowrap;
    flex-shrink: 0;
    transition: all 0.2s ease;
  }

  .unlike-btn:hover {
    color: #ffffff;
    background: var(--brand-red);
    border-color: var(--brand-red);
    box-shadow: 0 2px 10px rgba(229, 9, 20, 0.4);
  }

  /* Mobile Responsive Dropdown */
  @media (max-width: 480px) {
    .like-card {
      flex-direction: column;
      align-items: flex-start;
      gap: 1rem;
      padding: 1.1rem;
    }
    .like-main {
      width: 100%;
    }
    .unlike-btn {
      align-self: flex-end;
      width: 100%;
      text-align: center;
      padding: 0.55rem 0;
    }
  }
</style>
