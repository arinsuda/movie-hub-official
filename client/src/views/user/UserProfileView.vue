<template>
  <div class="profile-root">
    <!-- Ambient background grain -->
    <div class="grain" aria-hidden="true" />

    <div v-if="loading" class="profile-loading">
      <div class="loading-ring" />
      <span class="loading-text">Loading profile…</span>
    </div>

    <div v-else-if="!user" class="profile-empty">
      <Monitor class="empty-icon" />
      <p class="empty-title">User not found</p>
      <p class="empty-sub">Please try logging in again.</p>
    </div>

    <div v-else class="profile-layout">
      <!-- ── Hero Section ── -->
      <header class="profile-hero">
        <div class="hero-inner">
          <!-- Avatar -->
          <div class="avatar-shell">
            <div class="avatar-ring" />
            <Avatar
              v-if="user?.avatar_url"
              :image="user.avatar_url"
              shape="circle"
              class="avatar-img"
            />
            <div v-else class="avatar-fallback">
              <UserIcon :size="28" />
            </div>
            <span class="level-chip">Lv.{{ user?.level }}</span>
          </div>

          <!-- Identity -->
          <div class="hero-identity">
            <div class="identity-top">
              <div class="name-block">
                <h1 class="display-name">
                  {{ user?.display_name || user?.username }}
                </h1>
                <p class="username">@{{ user?.username }}</p>
              </div>
              <div class="hero-actions">
                <button v-if="notMe" class="btn btn-primary">
                  <UserPlus :size="14" />
                  <span>Follow</span>
                </button>
                <button v-else class="btn btn-ghost" @click="showEdit = true">
                  <Settings :size="14" />
                  <span>Edit Profile</span>
                </button>
              </div>
            </div>
            <p v-if="user?.bio" class="bio">{{ user.bio }}</p>
            <p v-else class="bio bio--empty">No bio yet.</p>
          </div>
        </div>
      </header>

      <!-- ── Stats Row ── -->
      <div class="stats-row">
        <div v-for="(s, i) in statItems" :key="i" class="stat-cell">
          <span class="stat-value">{{ s.value }}</span>
          <span class="stat-label">{{ s.label }}</span>
        </div>
      </div>

      <div class="divider" />

      <!-- ── Main Content Grid ── -->
      <div class="content-grid">
        <!-- Left: Nav -->
        <aside class="nav-aside">
          <p class="nav-eyebrow">Menu</p>
          <nav>
            <button
              v-for="tab in tabs"
              :key="tab.key"
              class="nav-item"
              :class="{ 'nav-item--active': activeTab === tab.key }"
              @click="activeTab = tab.key"
            >
              <component :is="tab.icon" :size="14" />
              <span>{{ tab.label }}</span>
              <span v-if="tab.count !== undefined" class="nav-count">{{
                tab.count
              }}</span>
            </button>
          </nav>
        </aside>

        <!-- Center: Tab Content -->
        <main class="content-main">
          <Transition name="fade-up" mode="out-in">
            <component
              :is="activeComponent"
              v-bind="activeProps"
              :key="activeTab"
            />
          </Transition>
        </main>

        <!-- Right: Badges -->
        <aside class="badges-aside">
          <p class="nav-eyebrow">Badges</p>
          <slot name="badges">
            <p class="badges-empty">No badges unlocked yet.</p>
          </slot>
        </aside>
      </div>
    </div>

    <!-- Edit Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showEdit"
          class="modal-backdrop"
          @click.self="showEdit = false"
        >
          <div class="modal-panel">
            <EditProfile :user="user" @close="handleEditClose" />
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
  import { authApi } from "@/api/api"
  import { useAuthStore } from "@/stores/auth"
  import type { UserProfile } from "@/types/user"
  import { computed, onMounted, ref } from "vue"

  import Avatar from "primevue/avatar"

  import {
    Star,
    Bookmark,
    Heart,
    Trophy,
    User as UserIcon,
    Settings,
    UserPlus,
    Monitor,
    TvMinimalPlay,
  } from "lucide-vue-next"

  import ProfileInfo from "@/components/profile/ProfileInfo.vue"
  import UserReviews from "@/components/profile/UserReviews.vue"
  import UserWatchlist from "@/components/profile/UserWatchlist.vue"
  import UserLikes from "@/components/profile/UserLikes.vue"
  import UserAchievements from "@/components/profile/UserAchievements.vue"
  import EditProfile from "@/components/profile/EditProfile.vue"

  const auth = useAuthStore()
  const userId = auth.user?.id ?? 0
  const user = ref<UserProfile | null>(null)
  const loading = ref(true)
  const notMe = ref(false)
  const showEdit = ref(false)

  type TabKey = "profile" | "reviews" | "watchlist" | "likes" | "achievements"
  const activeTab = ref<TabKey>("reviews")

  const tabs = computed(() => [
    {
      key: "reviews" as TabKey,
      label: "Reviews",
      icon: Star,
    },
    {
      key: "watchlist" as TabKey,
      label: "Watchlist",
      icon: Bookmark,
    },
    { key: "likes" as TabKey, label: "Likes", icon: Heart, count: undefined },
    {
      key: "watched" as TabKey,
      label: "Watched",
      icon: TvMinimalPlay,
    },
    {
      key: "achievements" as TabKey,
      label: "Achievements",
      icon: Trophy,
    },
  ])

  const statItems = computed(() => [
    { label: "Reviews", value: user.value?.review_count ?? 0 },
    { label: "Followers", value: user.value?.follower_count ?? 0 },
    { label: "Following", value: user.value?.following_count ?? 0 },
    { label: "Joined", value: joinYear.value },
  ])

  const componentMap: Record<string, unknown> = {
    profile: ProfileInfo,
    reviews: UserReviews,
    watchlist: UserWatchlist,
    likes: UserLikes,
    achievements: UserAchievements,
  }

  const activeComponent = computed(() => componentMap[activeTab.value])
  const activeProps = computed(() =>
    activeTab.value === "profile" ? { user: user.value } : { userId },
  )

  const joinYear = computed(() =>
    user.value?.created_at
      ? new Date(user.value.created_at).getFullYear()
      : "—",
  )
  const userLevel = computed(() => 18)
  const xpCurrent = computed(() => 3_400)
  const xpNext = computed(() => 5_000)
  const xpPercent = computed(() =>
    Math.round((xpCurrent.value / xpNext.value) * 100),
  )

  async function handleEditClose() {
    showEdit.value = false
    await loadProfile()
  }

  async function loadProfile() {
    if (!auth.user) return
    try {
      loading.value = true
      const res = await authApi.me(userId)
      user.value = res.data.user
    } catch (err) {
      console.error("fetchUserProfile failed:", err)
    } finally {
      loading.value = false
    }
  }

  onMounted(loadProfile)
</script>

<style scoped>
  /* ─────────────────────────────────────────
   Design Tokens (Updated for Minimal & Clean Visual)
───────────────────────────────────────── */
  .profile-root {
    --c-bg: #080808;
    --c-surface: #111111;
    --c-card: #161616;
    --c-border: rgba(255, 255, 255, 0.06);
    --c-border-h: rgba(255, 255, 255, 0.12);
    --c-red: #e1251b;
    --c-red-dim: rgba(225, 37, 27, 0.12);
    --c-text: #f0f0f0;
    --c-sub: #8a8a8e;
    --c-muted: #3a3a3c;

    /* ── ชุดฟอนต์ระดับพรีเมียมสไตล์ SaaS ยุคใหม่ (Inter / Geist / SF Pro) ── */
    --font-display:
      "Geist", "Inter", "Noto Sans Thai", -apple-system, BlinkMacSystemFont,
      "Segoe UI", sans-serif;
    --font-ui:
      "Inter", "Noto Sans Thai", -apple-system, BlinkMacSystemFont,
      "SF Pro Text", "Helvetica Neue", system-ui, sans-serif;

    --radius: 10px;
    --ease: cubic-bezier(0.16, 1, 0.3, 1);

    position: relative;
    min-height: 100vh;
    background: var(--c-bg);
    color: var(--c-text);
    font-family: var(--font-ui);
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  /* Grain overlay */
  .grain {
    position: fixed;
    inset: 0;
    pointer-events: none;
    z-index: 0;
    opacity: 0.025;
    background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
    background-size: 200px 200px;
  }

  /* ─────────── Loading ─────────── */
  .profile-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    min-height: 60vh;
  }

  .loading-ring {
    width: 32px;
    height: 32px;
    border: 2px solid var(--c-border);
    border-top-color: var(--c-red);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .loading-text {
    font-size: 0.8rem;
    color: var(--c-sub);
    letter-spacing: 0.04em;
    font-weight: 400;
  }

  /* ─────────── Empty ─────────── */
  .profile-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    min-height: 60vh;
  }

  .empty-icon {
    width: 28px;
    height: 28px;
    color: var(--c-muted);
  }
  .empty-title {
    font-size: 0.95rem;
    font-weight: 500;
    letter-spacing: -0.01em;
    margin: 0;
  }
  .empty-sub {
    font-size: 0.8rem;
    color: var(--c-sub);
    margin: 0;
    font-weight: 400;
  }

  /* ─────────── Layout ─────────── */
  .profile-layout {
    position: relative;
    z-index: 1;
    max-width: 1080px;
    margin: 0 auto;
    padding: 48px 24px 80px;
    animation: fadeIn 0.5s var(--ease) both;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* ─────────── Hero ─────────── */
  .profile-hero {
    margin-bottom: 32px;
  }

  .hero-inner {
    display: flex;
    align-items: flex-start;
    gap: 28px;
    margin-bottom: 20px;
  }

  /* Avatar */
  .avatar-shell {
    position: relative;
    flex-shrink: 0;
  }

  .avatar-ring {
    position: absolute;
    inset: -3px;
    border-radius: 50%;
    border: 1.5px solid var(--c-border-h);
    pointer-events: none;
  }

  .avatar-img {
    width: 88px !important;
    height: 88px !important;
  }
  .avatar-img :deep(img) {
    width: 88px;
    height: 88px;
    object-fit: cover;
    border-radius: 50%;
  }

  .avatar-fallback {
    width: 88px;
    height: 88px;
    border-radius: 50%;
    background: var(--c-card);
    border: 1px solid var(--c-border);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--c-muted);
  }

  .level-chip {
    position: absolute;
    bottom: -4px;
    left: 50%;
    transform: translateX(-50%);
    background: var(--c-surface);
    border: 1px solid var(--c-border-h);
    color: #fff;
    font-size: 0.58rem;
    font-weight: 600;
    letter-spacing: 0.04em;
    padding: 2px 7px;
    border-radius: 20px;
    white-space: nowrap;
  }

  /* Identity */
  .hero-identity {
    flex: 1;
    min-width: 0;
    padding-top: 4px;
  }

  .identity-top {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;
    margin-bottom: 10px;
    flex-wrap: wrap;
  }

  .name-block {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .display-name {
    font-family: var(--font-display);
    font-size: 1.65rem; /* บีบขนาดลงเล็กน้อยให้ดูสุขุมขึ้น */
    font-weight: 600; /* เปลี่ยนจาก 700 เพื่ออารมณ์ที่ Clean และไม่หนักเกินไป */
    color: #fff;
    margin: 0;
    line-height: 1.2;
    letter-spacing: -0.025em; /* บีบตัวอักษรชิดกันขึ้นสไตล์เว็บชั้นนำ */
  }

  .username {
    font-size: 0.8rem;
    color: var(--c-sub);
    margin: 0;
    letter-spacing: -0.01em;
  }

  .bio {
    font-size: 0.85rem;
    color: rgba(240, 240, 240, 0.75);
    line-height: 1.6;
    margin: 0;
    max-width: 560px;
    font-weight: 400;
  }
  .bio--empty {
    color: var(--c-muted);
    font-style: italic;
  }

  /* Buttons */
  .btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-family: var(--font-ui);
    font-size: 0.78rem;
    font-weight: 500;
    padding: 8px 14px;
    border-radius: 6px; /* บีบความโค้งลงมาที่ 6px ตามสไตล์โมเดิร์น */
    border: none;
    cursor: pointer;
    transition: all 0.2s var(--ease);
    white-space: nowrap;
  }

  .btn-primary {
    background: var(--c-red);
    color: #fff;
  }
  .btn-primary:hover {
    background: #ff3b30;
    transform: translateY(-0.5px);
    box-shadow: 0 4px 12px rgba(225, 37, 27, 0.25);
  }

  .btn-ghost {
    background: var(--c-card);
    color: var(--c-text);
    border: 1px solid var(--c-border);
  }
  .btn-ghost:hover {
    background: #1e1e1e;
    border-color: var(--c-border-h);
  }

  /* ─────────── Stats ─────────── */
  .stats-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1px;
    background: var(--c-border);
    border: 1px solid var(--c-border);
    border-radius: var(--radius);
    overflow: hidden;
    margin-bottom: 28px;
  }

  .stat-cell {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 16px 20px;
    background: var(--c-surface);
    transition: background 0.2s;
  }
  .stat-cell:hover {
    background: var(--c-card);
  }

  .stat-value {
    font-family: var(--font-display);
    font-size: 1.35rem;
    font-weight: 600; /* เปลี่ยนเป็น 600 เพื่อความสบายตา */
    color: #fff;
    line-height: 1;
    letter-spacing: -0.02em;
  }

  .stat-label {
    font-size: 0.65rem;
    font-weight: 500;
    letter-spacing: 0.05em; /* จัดโครงสร้างให้อ่านง่าย */
    text-transform: uppercase;
    color: var(--c-sub);
  }

  /* ─────────── Divider ─────────── */
  .divider {
    height: 1px;
    background: var(--c-border);
    margin-bottom: 24px;
  }

  /* ─────────── Content Grid ─────────── */
  .content-grid {
    display: grid;
    grid-template-columns: 180px 1fr 180px;
    gap: 0;
    min-height: 480px;
  }

  /* Nav Aside */
  .nav-aside {
    border-right: 1px solid var(--c-border);
    padding-right: 20px;
    padding-top: 4px;
  }

  .nav-eyebrow {
    font-size: 0.6rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--c-muted);
    margin: 0 0 12px 8px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 9px;
    width: 100%;
    padding: 8px 10px;
    border-radius: 6px;
    border: none;
    background: none;
    font-family: var(--font-ui);
    font-size: 0.8rem;
    color: var(--c-sub);
    cursor: pointer;
    transition: all 0.15s var(--ease);
    margin-bottom: 2px;
    text-align: left;
  }

  .nav-item:hover {
    color: var(--c-text);
    background: rgba(255, 255, 255, 0.03);
  }

  .nav-item--active {
    color: #fff;
    background: rgba(255, 255, 255, 0.05);
    font-weight: 500;
  }

  .nav-count {
    margin-left: auto;
    font-size: 0.6rem;
    font-weight: 600;
    background: var(--c-red);
    color: #fff;
    padding: 1px 5px;
    border-radius: 4px;
  }

  /* Main Content */
  .content-main {
    padding: 0 28px;
  }

  /* Badges Aside */
  .badges-aside {
    border-left: 1px solid var(--c-border);
    padding-left: 20px;
    padding-top: 4px;
  }

  .badges-empty {
    font-size: 0.78rem;
    color: var(--c-muted);
    margin: 0;
  }

  /* ─────────── Transitions ─────────── */
  .fade-up-enter-active,
  .fade-up-leave-active {
    transition: all 0.25s var(--ease);
  }
  .fade-up-enter-from {
    opacity: 0;
    transform: translateY(8px);
  }
  .fade-up-leave-to {
    opacity: 0;
    transform: translateY(-4px);
  }

  /* Modal */
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(6px);
    z-index: 999;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px;
  }

  .modal-panel {
    background: var(--c-surface);
    border: 1px solid var(--c-border-h);
    border-radius: 12px;
    width: 100%;
    max-width: 480px;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 24px 64px rgba(0, 0, 0, 0.6);
  }

  .modal-enter-active,
  .modal-leave-active {
    transition: all 0.3s var(--ease);
  }
  .modal-enter-from,
  .modal-leave-to {
    opacity: 0;
    transform: scale(0.96);
  }

  /* ─────────── Responsive ─────────── */
  @media (max-width: 900px) {
    .content-grid {
      grid-template-columns: 1fr;
    }
    .nav-aside {
      border-right: none;
      border-bottom: 1px solid var(--c-border);
      padding-right: 0;
      padding-bottom: 16px;
      margin-bottom: 20px;
    }
    nav {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
    }
    .nav-item {
      width: auto;
    }
    .badges-aside {
      border-left: none;
      border-top: 1px solid var(--c-border);
      padding-left: 0;
      padding-top: 20px;
      margin-top: 24px;
    }
    .content-main {
      padding: 0;
    }
    .stats-row {
      grid-template-columns: repeat(3, 1fr);
    }
  }

  @media (max-width: 600px) {
    .profile-layout {
      padding: 28px 16px 60px;
    }
    .hero-inner {
      flex-direction: column;
      align-items: center;
      text-align: center;
    }
    .identity-top {
      flex-direction: column;
      align-items: center;
    }
    .bio {
      text-align: center;
    }
    .stats-row {
      grid-template-columns: repeat(2, 1fr);
    }
  }
</style>
