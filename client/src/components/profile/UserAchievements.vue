<template>
  <div class="ach-root">
    <div class="section-head">
      <span class="eyebrow">Achievements</span>
      <span class="count-chip"
        >{{ unlockedCount }} / {{ achievements.length }}</span
      >
      <div class="rule" />
    </div>

    <div v-if="loading" class="state-loading">
      <div class="loader-bar"><div class="loader-fill" /></div>
    </div>

    <div v-else class="ach-list">
      <div
        v-for="(item, i) in achievements"
        :key="item.id"
        class="ach-card"
        :class="item.isUnlocked ? 'ach-card--unlocked' : 'ach-card--locked'"
        :style="{ '--i': i }"
      >
        <div class="ach-icon-box" :class="{ lit: item.isUnlocked }">
          <component :is="item.icon" :size="15" :stroke-width="1.8" />
        </div>

        <div class="ach-info">
          <p class="ach-name">{{ item.title }}</p>
          <p class="ach-desc">{{ item.description }}</p>
          <p v-if="item.isUnlocked && item.unlockedAt" class="ach-unlocked-at">
            <CheckCircle :size="9" /> {{ item.unlockedAt }}
          </p>
        </div>

        <div class="ach-badge">
          <span v-if="item.isUnlocked" class="badge badge--done">
            <CheckCircle :size="9" /> Unlocked
          </span>
          <span v-else class="badge badge--locked">
            <Lock :size="9" /> Locked
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
      await new Promise(r => setTimeout(r, 400))
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
  .ach-root {
    --c-card: #161616;
    --c-border: rgba(255, 255, 255, 0.06);
    --c-border-h: rgba(255, 255, 255, 0.12);
    --c-red: #e1251b;
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
    background: #fff;
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

  /* List */
  .ach-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .ach-card {
    display: grid;
    grid-template-columns: 40px 1fr auto;
    align-items: center;
    gap: 14px;
    padding: 14px 16px;
    background: #101010;
    border: 1px solid var(--c-border);
    border-radius: 10px;
    position: relative;
    overflow: hidden;
    transition: all 0.22s var(--ease);
    animation: cardIn 0.4s var(--ease) calc(var(--i) * 50ms) both;
  }

  @keyframes cardIn {
    from {
      opacity: 0;
      transform: translateY(8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .ach-card:hover {
    background: var(--c-card);
    border-color: var(--c-border-h);
  }

  /* Unlocked accent line */
  .ach-card--unlocked::before {
    content: "";
    position: absolute;
    left: 0;
    top: 20%;
    bottom: 20%;
    width: 2px;
    background: #fff;
    border-radius: 0 2px 2px 0;
  }

  .ach-card--locked {
    opacity: 0.4;
    filter: grayscale(1);
  }
  .ach-card--locked:hover {
    opacity: 0.55;
  }

  /* Icon */
  .ach-icon-box {
    width: 40px;
    height: 40px;
    background: var(--c-card);
    border: 1px solid var(--c-border);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--c-sub);
    transition: all 0.2s;
  }
  .ach-card--unlocked:hover .ach-icon-box {
    border-color: rgba(255, 255, 255, 0.18);
    color: #fff;
  }
  .ach-icon-box.lit {
    color: #fff;
    background: rgba(255, 255, 255, 0.04);
  }

  /* Info */
  .ach-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }
  .ach-name {
    font-size: 0.84rem;
    font-weight: 500;
    color: #fff;
    margin: 0;
  }
  .ach-desc {
    font-size: 0.73rem;
    color: var(--c-sub);
    line-height: 1.45;
    margin: 0;
  }
  .ach-unlocked-at {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 0.65rem;
    color: var(--c-sub);
    margin: 3px 0 0;
    opacity: 0.75;
  }

  /* Badges */
  .ach-badge {
    flex-shrink: 0;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 0.62rem;
    font-weight: 600;
    padding: 4px 9px;
    border-radius: 5px;
    white-space: nowrap;
  }
  .badge--done {
    background: rgba(255, 255, 255, 0.05);
    color: #fff;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  .badge--locked {
    background: transparent;
    color: var(--c-muted);
    border: 1px solid var(--c-border);
  }

  @media (max-width: 480px) {
    .ach-card {
      grid-template-columns: 40px 1fr;
      row-gap: 8px;
    }
    .ach-badge {
      grid-column: 2;
      justify-self: start;
    }
  }
</style>
