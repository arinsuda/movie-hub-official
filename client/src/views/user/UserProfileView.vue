<template>
  <div class="profile">
    <div class="profile-hero">
      <div class="hero-gradient" aria-hidden="true" />
    </div>

    <div class="profile-root">
      <div class="identity-row">
        <div class="avatar-container">
          <div class="avatar-frame">
            <img
              v-if="user?.avatar_url"
              :src="user.avatar_url"
              :alt="user.display_name || user.username"
              class="avatar-img"
            />
            <div v-else class="avatar-placeholder">
              <UserIcon :size="28" />
            </div>
          </div>
          <span class="user-minimal-level">{{ userLevel }}</span>
        </div>

        <div class="identity-meta">
          <div class="meta-top-row">
            <div class="name-zone">
              <h1 class="display-name">
                {{ user?.display_name || user?.username }}
              </h1>
            </div>

            <div class="header-actions">
              <button v-if="notMe" class="btn-primary">
                <UserPlus :size="13" />
                Follow
              </button>
              <button v-else class="btn-ghost">
                <Settings :size="13" />
                Edit Profile
              </button>
            </div>
          </div>

          <p class="username-handle">@{{ user?.username }}</p>
          <p v-if="user?.bio" class="bio-text">{{ user.bio }}</p>
        </div>
      </div>

      <div v-if="user" class="stats-bar">
        <div class="stat-card">
          <span class="stat-num">{{ user.review_count }}</span>
          <span class="stat-lbl">Reviews</span>
        </div>
        <div class="stat-card">
          <span class="stat-num">{{ user.follower_count }}</span>
          <span class="stat-lbl">Followers</span>
        </div>
        <div class="stat-card">
          <span class="stat-num">{{ user.following_count }}</span>
          <span class="stat-lbl">Following</span>
        </div>
        <div class="stat-card">
          <span class="stat-num">{{ avgRating }}</span>
          <span class="stat-lbl">Avg Rating</span>
        </div>
        <div class="stat-card">
          <span class="stat-num">{{ joinYear }}</span>
          <span class="stat-lbl">Joined</span>
        </div>
      </div>

      <div v-if="loading" class="global-loading">
        <div class="loading-bar"><div class="loading-fill" /></div>
        <span>Loading profile data...</span>
      </div>

      <div v-else-if="!user" class="empty-state">
        <Monitor :size="32" :stroke-width="1" />
        <p>User not found</p>
        <span>Please try logging in again.</span>
      </div>

      <div v-else class="page-body">
        <aside class="sidebar-left">
          <div class="sidebar-section">
            <span class="sidebar-label">Menu</span>
            <ul class="sidebar-nav" role="tablist">
              <li
                v-for="tab in tabs"
                :key="tab.key"
                role="tab"
                :aria-selected="activeTab === tab.key"
                class="nav-item"
                :class="{ 'nav-item--active': activeTab === tab.key }"
                @click="activeTab = tab.key"
              >
                <div class="nav-link-content">
                  <component :is="tab.icon" :size="15" aria-hidden="true" />
                  <span class="nav-text">{{ tab.label }}</span>
                </div>
                <span v-if="tab.count !== undefined" class="nav-count">{{
                  tab.count
                }}</span>
              </li>
            </ul>
          </div>
        </aside>

        <main class="main-content">
          <component :is="activeComponent" v-bind="activeProps" />
        </main>

        <aside class="sidebar-right">
          <div class="sidebar-section">
            <span class="sidebar-label">Badges</span>
            <slot name="badges">
              <p class="badge-empty">No badges unlocked</p>
            </slot>
          </div>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { authApi } from "@/api/api"
  import { useAuthStore } from "@/stores/auth"
  import type { UserProfile } from "@/types"
  import { computed, onMounted, ref } from "vue"
  import {
    Star,
    Bookmark,
    Heart,
    Trophy,
    User as UserIcon,
    Settings,
    UserPlus,
    Monitor,
    TvMinimalPlay
  } from "lucide-vue-next"

  import ProfileInfo from "@/components/profile/ProfileInfo.vue"
  import UserReviews from "@/components/profile/UserReviews.vue"
  import UserWatchlist from "@/components/profile/UserWatchlist.vue"
  import UserLikes from "@/components/profile/UserLikes.vue"
  import UserAchievements from "@/components/profile/UserAchievements.vue"

  const auth = useAuthStore()
  const userId = auth.user?.id ?? 0
  const user = ref<UserProfile | null>(null)
  const loading = ref(true)

  const notMe = ref(false)

  type TabKey = "profile" | "reviews" | "watchlist" | "likes" | "achievements"
  const activeTab = ref<TabKey>("reviews")

  const tabs = computed(() => [
    {
      key: "reviews" as TabKey,
      label: "My Reviews",
      icon: Star,
      count: user.value?.review_count,
    },
    {
      key: "watchlist" as TabKey,
      label: "Watchlist",
      icon: Bookmark,
      count: undefined,
    },
    {
      key: "likes" as TabKey,
      label: "Likes",
      icon: Heart,
      count: undefined,
    },
    {
      key: "watched" as TabKey,
      label: "Watched",
      icon: TvMinimalPlay,
      count: undefined,
    },
    {
      key: "achievements" as TabKey,
      label: "Achievements",
      icon: Trophy,
      count: undefined,
    },
  ])

  const componentMap: Record<TabKey, unknown> = {
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

  const avgRating = computed(() => "4.2")
  const userLevel = computed(() => 18)

  onMounted(async () => {
    if (auth.user) {
      try {
        loading.value = true
        const res = await authApi.me(userId)
        user.value = res.data.user
      } catch (err) {
        console.error("fetchUserProfile failed:", err)
      } finally {
        loading.value = false
      }
    } else {
      loading.value = false
    }
  })
</script>

<style scoped>
  .profile {
    width: 100%;
    min-height: 100vh;
    background: #0a0a0a;
  }

  .profile-root {
    --surface: #121212;
    --s2: #1a1a1a;
    --s3: #242424;
    --border: rgba(255, 255, 255, 0.05);
    --border2: rgba(255, 255, 255, 0.08);
    --red: #e1251b;
    --red-b: #ff3b30;
    --red-glow: rgba(225, 37, 27, 0.05);
    --text: #f5f5f7;
    --sub: #8e8e93;
    --muted: #48484a;
    --ui:
      -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue",
      system-ui, sans-serif;
    --radius-sm: 8px;
    --radius-md: 12px;
    --radius-lg: 16px;

    color: var(--text);
    font-family: var(--ui);

    max-width: 1100px;
    margin: 0 auto;
    padding: 0 1.5rem;
    width: 100%;
  }

  .profile-hero {
    position: relative;
    width: 100%;
    height: 160px;
    overflow: hidden;
    background: #0f0f11;
  }

  .hero-gradient {
    position: absolute;
    inset: 0;
    background: radial-gradient(
      circle at 80% 20%,
      rgba(225, 37, 27, 0.06) 0%,
      transparent 60%
    );
  }

  .identity-row {
    display: flex;
    align-items: flex-end;
    gap: 2rem;
    padding: 0 1rem;
    margin-top: -40px;
    position: relative;
    z-index: 2;
  }

  .avatar-container {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    flex-shrink: 0;
  }

  .avatar-frame {
    position: relative;
  }

  .avatar-img,
  .avatar-placeholder {
    width: 130px;
    height: 130px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 4px solid #0a0a0a;
    border-radius: 50%;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
  }

  .avatar-img {
    object-fit: cover;
  }

  .avatar-placeholder {
    background: var(--s2);
    color: var(--sub);
  }

  .user-minimal-level {
    position: absolute;
    bottom: -12px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text);
    border: 2px solid #0a0a0a;
    background: var(--s3);
    border-radius: 50%;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);
  }

  .identity-meta {
    flex: 1;
    padding-bottom: 0.25rem;
  }

  .meta-top-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1.5rem;
  }

  .name-zone {
    flex: 1;
  }

  .display-name {
    font-size: 1.8rem;
    font-weight: 600;
    color: #ffffff;
    line-height: 1.2;
    letter-spacing: -0.01em;
    margin: 0;
  }

  .username-handle {
    font-size: 0.85rem;
    color: var(--sub);
    margin: 6px 0 12px 0;
  }

  .bio-text {
    font-size: 0.9rem;
    color: var(--sub);
    line-height: 1.5;
    margin: 0;
    max-width: 550px;
  }

  .header-actions {
    display: flex;
    align-items: center;
    flex-shrink: 0;
  }

  .btn-primary,
  .btn-ghost {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-family: var(--ui);
    font-size: 0.75rem;
    font-weight: 500;
    padding: 8px 16px;
    cursor: pointer;
    border: 1px solid transparent;
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
    white-space: nowrap;
  }

  .btn-primary {
    background: var(--text);
    color: #0a0a0a;
  }
  .btn-primary:hover {
    background: #ffffff;
    transform: translateY(-1px);
  }

  .btn-ghost {
    background: var(--surface);
    border-color: var(--border);
    color: var(--text);
  }
  .btn-ghost:hover {
    background: var(--s2);
    border-color: var(--border2);
  }

  /* ── Stats Bar ── */
  .stats-bar {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 1rem;
    padding: 0 1rem;
    margin-top: 2.5rem;
    margin-bottom: 2rem;
  }

  .stat-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    padding: 1.25rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 6px;
    transition:
      transform 0.2s ease,
      border-color 0.2s ease;
  }

  .stat-card:hover {
    border-color: var(--border2);
    transform: translateY(-2px);
  }

  .stat-num {
    font-size: 1.5rem;
    font-weight: 600;
    color: #ffffff;
    line-height: 1;
  }

  .stat-lbl {
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--sub);
  }

  /* ── Page Body ── */
  .page-body {
    display: grid;
    grid-template-columns: 220px 1fr 220px;
    gap: 1rem;
    min-height: 50vh;
    padding: 0 1rem;
    margin-bottom: 4rem;
  }

  /* ── sidebar left & right ── */
  .sidebar-left {
    border-right: 1px solid var(--border);
    padding: 1rem 1.5rem 2rem 0; /* ถอยชิดซ้ายสุดแนวคอนเทนต์ */
  }

  .sidebar-right {
    border-left: 1px solid var(--border);
    padding: 1rem 0 2rem 1.5rem; /* ถอยชิดขวาสุดแนวคอนเทนต์ */
  }

  .sidebar-section {
    display: flex;
    flex-direction: column;
  }

  .sidebar-label {
    display: block;
    font-size: 0.7rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    color: var(--muted);
    margin-bottom: 1rem;
    text-transform: uppercase;
  }

  .sidebar-nav {
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 0;
    margin: 0;
  }

  .nav-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 12px;
    font-size: 0.85rem;
    color: var(--sub);
    cursor: pointer;
    border-radius: var(--radius-sm);
    transition: all 0.15s ease;
    user-select: none;
  }

  .nav-link-content {
    display: flex;
    align-items: center;
    gap: 12px;
    line-height: 1;
  }

  .nav-text {
    margin-top: 1px;
  }

  .nav-item:hover {
    color: var(--text);
    background: var(--surface);
  }
  .nav-item--active {
    color: #ffffff;
    background: var(--surface);
    border: 1px solid var(--border);
    font-weight: 500;
  }

  .nav-count {
    font-size: 0.7rem;
    color: var(--sub);
    background: var(--s2);
    padding: 2px 8px;
    border-radius: 10px;
  }

  .badge-empty {
    font-size: 0.8rem;
    color: var(--muted);
  }

  /* ── main content ── */
  .main-content {
    padding: 1rem 1rem 2rem;
  }

  /* ── global loading ── */
  .global-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 5rem 0;
    color: var(--sub);
    font-size: 0.85rem;
  }

  .loading-bar {
    width: 100px;
    height: 2px;
    background: var(--border);
    border-radius: var(--radius-sm);
    overflow: hidden;
  }

  .loading-fill {
    height: 100%;
    background: var(--red);
    border-radius: var(--radius-sm);
    animation: sweep 1.4s ease-in-out infinite;
  }

  @keyframes sweep {
    0% {
      transform: translateX(-100%);
      width: 50%;
    }
    50% {
      transform: translateX(100%);
      width: 50%;
    }
    100% {
      transform: translateX(-100%);
      width: 50%;
    }
  }

  /* ── empty ── */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    padding: 5rem 0;
    color: var(--muted);
  }
  .empty-state p {
    font-size: 0.9rem;
    color: var(--sub);
    margin: 0;
  }
  .empty-state span {
    font-size: 0.8rem;
    color: var(--muted);
  }

  /* ── responsive ── */
  @media (max-width: 1100px) {
    .page-body {
      grid-template-columns: 200px 1fr;
      padding: 0;
    }
    .sidebar-right {
      grid-column: span 2;
      border-left: none;
      border-top: 1px solid var(--border);
      padding: 1.5rem 0;
    }
    .stats-bar {
      grid-template-columns: repeat(3, 1fr);
      padding: 0;
    }
    .identity-row {
      padding: 0;
    }
  }

  @media (max-width: 768px) {
    .page-body {
      grid-template-columns: 1fr;
    }
    .sidebar-left {
      border-right: none;
      border-bottom: 1px solid var(--border);
      padding: 1.5rem 0;
    }
    .sidebar-right {
      grid-column: span 1;
      padding: 1.5rem 0;
    }
    .meta-top-row {
      flex-direction: column;
      align-items: flex-start;
      gap: 1rem;
    }
    .identity-row {
      flex-direction: column;
      align-items: center;
      text-align: center;
      gap: 1.5rem;
    }
    .meta-top-row {
      align-items: center;
    }
    .bio-text {
      margin: 0 auto;
    }
    .stats-bar {
      grid-template-columns: repeat(2, 1fr);
    }
    .main-content {
      padding: 1.5rem 0;
    }
  }
</style>
