<template>
  <div class="likes-root">
    <div class="section-head">
      <span class="eyebrow">Liked Items</span>
      <span class="count-chip">{{ likedItems.length }}</span>
      <div class="rule" />
    </div>

    <div v-if="loading" class="state-loading">
      <div class="loader-bar"><div class="loader-fill" /></div>
    </div>

    <div v-else-if="likedItems.length === 0" class="state-empty">
      <Heart :size="28" :stroke-width="1.2" />
      <p>No liked items yet</p>
    </div>

    <div v-else class="likes-list">
      <div
        v-for="(item, i) in likedItems"
        :key="item.id"
        class="like-row"
        :style="{ '--i': i }"
      >
        <div class="heart-dot" aria-hidden="true">
          <Heart :size="13" class="heart-icon" />
        </div>

        <div class="like-body">
          <h4 class="like-title">{{ item.title }}</h4>
          <p class="like-meta">by {{ item.author }} · {{ item.likedAt }}</p>
        </div>

        <button
          class="unlike-btn"
          title="Unlike"
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

  const props = defineProps<{ userId: number }>()

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
      await new Promise(r => setTimeout(r, 400))
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
          author: "Dev Community",
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

  function handleUnlike(id: number) {
    likedItems.value = likedItems.value.filter(i => i.id !== id)
  }
</script>

<style scoped>
  .likes-root {
    --c-card: #161616;
    --c-border: rgba(255, 255, 255, 0.06);
    --c-border-h: rgba(255, 255, 255, 0.12);
    --c-red: #e1251b;
    --c-red-dim: rgba(225, 37, 27, 0.1);
    --c-text: #f0f0f0;
    --c-sub: #8a8a8e;
    --c-muted: #3a3a3c;
    --font:
      -apple-system, BlinkMacSystemFont, "SF Pro Text", system-ui, sans-serif;
    --ease: cubic-bezier(0.16, 1, 0.3, 1);

    font-family: var(--font);
    color: var(--c-text);
  }

  .section-head {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 20px;
  }
  .eyebrow {
    font-size: 0.6rem;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--c-muted);
    white-space: nowrap;
  }
  .count-chip {
    font-size: 0.65rem;
    font-weight: 700;
    background: var(--c-red);
    color: #fff;
    padding: 2px 7px;
    border-radius: 4px;
  }
  .rule {
    flex: 1;
    height: 1px;
    background: var(--c-border);
  }

  /* States */
  .state-loading {
    padding: 48px 0;
    display: flex;
    justify-content: center;
  }
  .loader-bar {
    width: 100px;
    height: 2px;
    background: var(--c-border);
    border-radius: 2px;
    overflow: hidden;
    position: relative;
  }
  .loader-fill {
    height: 100%;
    width: 40%;
    background: var(--c-red);
    position: absolute;
    animation: sweep 1.4s infinite ease-in-out;
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
    gap: 10px;
    padding: 56px 0;
    color: var(--c-muted);
  }
  .state-empty p {
    font-size: 0.82rem;
    margin: 0;
    color: var(--c-sub);
  }

  /* List */
  .likes-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .like-row {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 14px 16px;
    background: #101010;
    border: 1px solid var(--c-border);
    border-radius: 10px;
    transition: all 0.22s var(--ease);
    animation: rowIn 0.4s var(--ease) calc(var(--i) * 50ms) both;
  }
  .like-row:hover {
    background: var(--c-card);
    border-color: var(--c-border-h);
    transform: translateX(3px);
  }

  @keyframes rowIn {
    from {
      opacity: 0;
      transform: translateX(-8px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  /* Heart dot */
  .heart-dot {
    width: 34px;
    height: 34px;
    background: var(--c-red-dim);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    transition: background 0.2s;
  }
  .like-row:hover .heart-dot {
    background: rgba(225, 37, 27, 0.18);
  }
  .heart-icon {
    color: var(--c-red);
    fill: var(--c-red);
  }

  /* Body */
  .like-body {
    display: flex;
    flex-direction: column;
    gap: 3px;
    flex: 1;
    min-width: 0;
  }
  .like-title {
    font-size: 0.84rem;
    font-weight: 500;
    color: #fff;
    margin: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    transition: color 0.15s;
  }
  .like-row:hover .like-title {
    color: rgba(255, 255, 255, 0.85);
  }
  .like-meta {
    font-size: 0.7rem;
    color: var(--c-sub);
    margin: 0;
  }

  /* Unlike btn */
  .unlike-btn {
    font-family: var(--font);
    font-size: 0.7rem;
    font-weight: 500;
    color: var(--c-sub);
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--c-border);
    padding: 5px 12px;
    border-radius: 6px;
    cursor: pointer;
    white-space: nowrap;
    flex-shrink: 0;
    transition: all 0.2s var(--ease);
  }
  .unlike-btn:hover {
    color: #fff;
    background: var(--c-red);
    border-color: var(--c-red);
    box-shadow: 0 3px 12px rgba(225, 37, 27, 0.35);
  }

  @media (max-width: 440px) {
    .like-row {
      flex-wrap: wrap;
    }
    .unlike-btn {
      width: 100%;
      text-align: center;
    }
  }
</style>
