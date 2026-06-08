<template>
  <div class="ach-wrap">
    <div class="section-head">
      <span class="section-label">Achievements</span>
      <span class="count-badge">
        {{ unlockedCount }} / {{ achievements.length }}
      </span>
      <div class="section-rule" />
    </div>

    <div v-if="loading" class="state-loading">
      <div class="loading-bar"><div class="loading-fill" /></div>
    </div>

    <div v-else class="ach-list">
      <div
        v-for="item in achievements"
        :key="item.id"
        class="ach-card"
        :class="item.isUnlocked ? 'ach-card--unlocked' : 'ach-card--locked'"
        :aria-disabled="!item.isUnlocked"
      >
        <div class="ach-icon" :class="item.isUnlocked ? 'ach-icon--lit' : ''">
          <component
            :is="item.icon"
            :size="16"
            :stroke-width="1.8"
            aria-hidden="true"
          />
        </div>

        <div class="ach-info">
          <p class="ach-title">{{ item.title }}</p>
          <p class="ach-desc">{{ item.description }}</p>
          <p v-if="item.isUnlocked && item.unlockedAt" class="ach-date">
            <CheckCircle :size="10" aria-hidden="true" />
            Unlocked on {{ item.unlockedAt }}
          </p>
        </div>

        <div class="ach-status">
          <span v-if="item.isUnlocked" class="status-badge status-badge--done">
            <CheckCircle :size="10" aria-hidden="true" />
            Unlocked
          </span>
          <span v-else class="status-badge status-badge--pending">
            <Lock :size="10" aria-hidden="true" />
            Locked
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, onMounted, ref } from "vue"
  import {
    UserCheck,
    PenLine,
    Eye,
    Heart,
    Flame,
    CheckCircle,
    Lock,
    type LucideIcon,
  } from "lucide-vue-next"

  const props = defineProps<{ userId: number }>()

  const loading = ref(false)

  interface Achievement {
    id: number
    title: string
    description: string
    icon: LucideIcon
    isUnlocked: boolean
    unlockedAt?: string
  }

  const achievements = ref<Achievement[]>([])

  const unlockedCount = computed(
    () => achievements.value.filter(a => a.isUnlocked).length,
  )

  onMounted(async () => {
    try {
      loading.value = true
      // จำลอง API response เล็กน้อย
      await new Promise(resolve => setTimeout(resolve, 400))
      achievements.value = [
        {
          id: 1,
          title: "First Step",
          description: "สมัครสมาชิกและตั้งค่าข้อมูลโปรไฟล์เสร็จสมบูรณ์",
          icon: UserCheck,
          isUnlocked: true,
          unlockedAt: "10 Jan 2026",
        },
        {
          id: 2,
          title: "Pro Critic",
          description: "เขียนรีวิวภาพยนตร์หรือเนื้อหาครบ 5 ครั้ง",
          icon: PenLine,
          isUnlocked: true,
          unlockedAt: "15 May 2026",
        },
        {
          id: 3,
          title: "True Supporter",
          description: "กดถูกใจรีวิวหรือคอมเมนต์ที่ชื่นชอบครบ 20 ครั้ง",
          icon: Heart,
          isUnlocked: true,
          unlockedAt: "1 Jun 2026",
        },
        {
          id: 4,
          title: "Night Watcher",
          description: "เพิ่มรายการหนังไว้ใน Watchlist ครบ 10 เรื่อง",
          icon: Eye,
          isUnlocked: false,
        },
        {
          id: 5,
          title: "Elite Streak",
          description: "เข้าใช้งานแพลตฟอร์มติดต่อกันยาวนานครบ 30 วัน",
          icon: Flame,
          isUnlocked: false,
        },
      ]
    } catch (err) {
      console.error("Fetch achievements failed:", err)
    } finally {
      loading.value = false
    }
  })
</script>

<style scoped>
  .ach-wrap {
    /* ใช้ Token ดีไซน์ชุดเดียวกับระเบียบ UI หลัก */
    --surface: #121212;
    --surface-card: #1a1a1a;
    --brand-red: #e50914;
    --brand-red-hover: #b20710;
    --border: rgba(255, 255, 255, 0.05);
    --border-hover: rgba(255, 255, 255, 0.1);
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
    font-variant-numeric: tabular-nums;
  }

  .section-rule {
    flex: 1;
    height: 1px;
    background: var(--border);
  }

  /* ── States ── */
  .state-loading {
    padding: 3rem 0;
    display: flex;
    justify-content: center;
  }
  .loading-bar {
    width: 120px;
    height: 2px;
    background: var(--border);
    border-radius: 2px;
    overflow: hidden;
    position: relative;
  }
  .loading-fill {
    height: 100%;
    width: 40%;
    background: #ffffff;
    position: absolute;
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

  /* ── List Layout ── */
  .ach-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .ach-card {
    display: grid;
    grid-template-columns: 40px 1fr auto;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 8px;
    position: relative;
    transition: all 0.2s ease;
  }

  .ach-card:hover {
    border-color: var(--border-hover);
    background: var(--surface-card);
  }

  /* เลื่อน Accent Bar ด้านข้างให้ดูคมและมินิมอลขึ้น */
  .ach-card--unlocked::before {
    content: "";
    position: absolute;
    left: 0;
    top: 25%;
    bottom: 25%;
    width: 3px;
    background: #ffffff;
    border-radius: 0 4px 4px 0;
  }

  .ach-card--locked {
    opacity: 0.45;
    filter: grayscale(1);
  }

  /* Icon Area */
  .ach-icon {
    width: 40px;
    height: 40px;
    background: var(--surface-card);
    border: 1px solid var(--border);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--sub);
    transition: all 0.2s ease;
  }

  .ach-card--unlocked:hover .ach-icon {
    border-color: rgba(255, 255, 255, 0.2);
    color: #ffffff;
  }

  .ach-icon--lit {
    color: #ffffff;
    background: rgba(255, 255, 255, 0.03);
  }

  /* Info Block */
  .ach-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }

  .ach-title {
    font-size: 0.85rem;
    font-weight: 500;
    color: #ffffff;
    margin: 0;
  }

  .ach-desc {
    font-size: 0.75rem;
    color: var(--sub);
    line-height: 1.4;
    margin: 0;
  }

  .ach-date {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 0.68rem;
    color: var(--sub);
    margin: 4px 0 0 0;
    opacity: 0.8;
  }

  /* Status Badges */
  .ach-status {
    flex-shrink: 0;
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 0.65rem;
    font-weight: 500;
    padding: 4px 8px;
    border-radius: 4px;
    white-space: nowrap;
  }

  .status-badge--done {
    background: rgba(255, 255, 255, 0.06);
    color: #ffffff;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .status-badge--pending {
    background: transparent;
    color: var(--muted);
    border: 1px solid var(--border);
  }

  /* Mobile Responsive */
  @media (max-width: 540px) {
    .ach-card {
      grid-template-columns: 40px 1fr;
      row-gap: 0.75rem;
    }
    .ach-status {
      grid-column: 2;
      justify-self: start;
    }
  }
</style>
