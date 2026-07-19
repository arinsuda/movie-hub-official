<template>
  <div class="profile-root">
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
        <div class="hero-backdrop" aria-hidden="true" />
        <div class="hero-inner">
          <!-- Avatar -->
          <div class="avatar-shell">
            <div class="avatar-ring">
              <div class="avatar-inner">
                <Avatar
                  v-if="user?.avatar_url"
                  :image="user.avatar_url"
                  shape="circle"
                  class="avatar-img"
                />
                <div v-else class="avatar-fallback">
                  <UserIcon :size="32" />
                </div>
              </div>
            </div>
            <span class="level-chip">Lv.{{ user?.level ?? 1 }}</span>
          </div>

          <!-- Identity -->
          <div class="hero-identity">
            <h1 class="display-name">
              {{ user?.display_name || user?.username }}
            </h1>
            <p class="username">@{{ user?.username }}</p>
            <p v-if="user?.bio" class="bio">{{ user.bio }}</p>
            <p v-else class="bio bio--empty">No bio yet.</p>

            <!-- Inline stats row (Instagram-style) -->
            <div class="hero-stats">
              <div class="hstat">
                <span class="hstat-val">{{ user?.review_count ?? 0 }}</span>
                <span class="hstat-lbl">Reviews</span>
              </div>
              <div class="hstat-sep" />
              <div class="hstat">
                <span class="hstat-val">{{ user?.follower_count ?? 0 }}</span>
                <span class="hstat-lbl">Followers</span>
              </div>
              <div class="hstat-sep" />
              <div class="hstat">
                <span class="hstat-val">{{ user?.following_count ?? 0 }}</span>
                <span class="hstat-lbl">Following</span>
              </div>
              <div class="hstat-sep" />
              <div class="hstat">
                <span class="hstat-val">{{ joinYear }}</span>
                <span class="hstat-lbl">Joined</span>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="hero-actions">
            <button
              v-if="notMe"
              class="btn"
              :class="followButtonClass"
              :disabled="followLoading || isPendingFollow"
              @click="handleFollowButtonClick"
            >
              <component :is="followButtonIcon" :size="14" />
              <span>{{ followButtonLabel }}</span>
            </button>

            <button v-else class="btn btn-ghost" @click="showEdit = true">
              <Settings :size="14" />
              <span>Edit Profile</span>
            </button>
          </div>
        </div>
      </header>

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
              :key="activeComponentKey"
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
    <ConfirmModal
      v-model="showUnfollowConfirm"
      list-type="unfollow"
      :item-name="user?.display_name || user?.username"
      :confirm-disabled="followLoading"
      @confirm="confirmUnfollow"
      @cancel="showUnfollowConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
  import { followApi, userApi } from "@/api/api"
  import { useAuthStore } from "@/stores/auth"
  import type { UserProfile } from "@/types/user"
  import type { ListType } from "@/types"
  import { computed, ref, watch } from "vue"
  import { useRoute } from "vue-router"

  import Avatar from "primevue/avatar"

  import {
    Star,
    Bookmark,
    Heart,
    Trophy,
    User as UserIcon,
    Settings,
    UserPlus,
    UserCheck,
    Clock,
    Monitor,
    TvMinimalPlay,
  } from "lucide-vue-next"

  import ProfileInfo from "@/components/profile/ProfileInfo.vue"
  import UserReviews from "@/components/profile/UserReviews.vue"
  import UserWatchlist from "@/components/profile/UserWatchlist.vue"
  import UserWatched from "@/components/profile/UserWatched.vue"
  import UserLikes from "@/components/profile/UserLikes.vue"
  import UserAchievements from "@/components/profile/UserAchievements.vue"
  import EditProfile from "@/components/profile/EditProfile.vue"
  import ConfirmModal from "@/components/profile/components/ConfirmModal.vue"

  const route = useRoute()
  const auth = useAuthStore()

  const user = ref<UserProfile | null>(null)
  const loading = ref(true)
  const showEdit = ref(false)

  type FollowStatus = "none" | "pending" | "accepted"

  const followStatus = ref<FollowStatus>("none")
  const followLoading = ref(false)
  const showUnfollowConfirm = ref(false)

  const viewedUserId = computed(() => {
    const raw = route.params.userId
    const value = Array.isArray(raw) ? raw[0] : raw
    return Number(value)
  })

  const currentUserId = computed(() => auth.user?.id ?? null)

  const notMe = computed(() => {
    return currentUserId.value !== viewedUserId.value
  })

  const isFollowing = computed(() => followStatus.value === "accepted")
  const isPendingFollow = computed(() => followStatus.value === "pending")

  const followButtonLabel = computed(() => {
    if (followLoading.value) return "กำลังโหลด..."
    if (isFollowing.value) return "กำลังติดตาม"
    if (isPendingFollow.value) return "ส่งคำขอแล้ว"
    return "Follow"
  })

  const followButtonIcon = computed(() => {
    if (isFollowing.value) return UserCheck
    if (isPendingFollow.value) return Clock
    return UserPlus
  })

  const followButtonClass = computed(() => {
    if (isFollowing.value) return "btn-following"
    if (isPendingFollow.value) return "btn-pending"
    return "btn-primary"
  })

  type TabKey =
    | "profile"
    | "reviews"
    | "watchlist"
    | "likes"
    | "watched"
    | "achievements"

  const activeTab = ref<TabKey>("reviews")

  const tabs = computed(() => [
    { key: "reviews" as TabKey, label: "Reviews", icon: Star },
    { key: "watchlist" as TabKey, label: "Watchlist", icon: Bookmark },
    { key: "likes" as TabKey, label: "Likes", icon: Heart, count: undefined },
    { key: "watched" as TabKey, label: "Watched", icon: TvMinimalPlay },
    { key: "achievements" as TabKey, label: "Achievements", icon: Trophy },
  ])

  const componentMap: Record<TabKey, unknown> = {
    profile: ProfileInfo,
    reviews: UserReviews,
    watchlist: UserWatchlist,
    likes: UserLikes,
    watched: UserWatched,
    achievements: UserAchievements,
  }

  const activeComponent = computed(() => componentMap[activeTab.value])

  const activeProps = computed(() => {
    const typeMap: Record<string, ListType | undefined> = {
      watchlist: "watchlist",
      likes: "likes",
      watched: "watched",
    }

    return {
      userId: viewedUserId.value,
      listType: typeMap[activeTab.value],
    }
  })

  const activeComponentKey = computed(() => {
    return `${activeTab.value}-${viewedUserId.value}`
  })

  const joinYear = computed(() =>
    user.value?.created_at
      ? new Date(user.value.created_at).getFullYear()
      : "—",
  )

  function normalizeFollowStatus(data: unknown): FollowStatus {
    const value = data as {
      status?: string
      is_following?: boolean
      following?: boolean
      is_pending?: boolean
    }

    if (value.status === "accepted" || value.status === "following") {
      return "accepted"
    }

    if (value.status === "pending") {
      return "pending"
    }

    if (value.is_following || value.following) {
      return "accepted"
    }

    if (value.is_pending) {
      return "pending"
    }

    return "none"
  }

  async function loadProfile() {
    if (!Number.isInteger(viewedUserId.value) || viewedUserId.value <= 0) {
      user.value = null
      loading.value = false
      return
    }

    try {
      loading.value = true
      const res = await userApi.getProfile(viewedUserId.value)
      user.value = res.data.user
    } catch (err) {
      user.value = null
      console.error("fetchUserProfile failed:", err)
    } finally {
      loading.value = false
    }
  }

  async function loadFollowStatus() {
    if (!auth.user || !notMe.value) {
      followStatus.value = "none"
      return
    }

    try {
      const res = await followApi.getFollowStatus(viewedUserId.value)
      followStatus.value = normalizeFollowStatus(res.data)
    } catch (err) {
      followStatus.value = "none"
      console.error("getFollowStatus failed:", err)
    }
  }

  async function loadProfilePage() {
    await loadProfile()
    await loadFollowStatus()
  }

  async function handleEditClose() {
    showEdit.value = false
    await loadProfile()
  }

  async function handleFollowButtonClick() {
    if (followLoading.value) return

    if (isFollowing.value) {
      showUnfollowConfirm.value = true
      return
    }

    if (isPendingFollow.value) {
      return
    }

    await followUser()
  }

  async function followUser() {
    try {
      followLoading.value = true

      const res = await followApi.follow(viewedUserId.value)
      followStatus.value = normalizeFollowStatus(res.data)

      await loadProfile()
    } catch (err) {
      console.error("follow failed:", err)
    } finally {
      followLoading.value = false
    }
  }

  async function confirmUnfollow() {
    try {
      followLoading.value = true

      await followApi.unfollow(viewedUserId.value)

      followStatus.value = "none"
      showUnfollowConfirm.value = false

      await loadProfile()
    } catch (err) {
      console.error("unfollow failed:", err)
    } finally {
      followLoading.value = false
    }
  }

  watch(viewedUserId, loadProfilePage, { immediate: true })
</script>

<style scoped>
  /* ─────────────────────────────────────────
   Design Tokens
  ───────────────────────────────────────── */
  .profile-root {
    --c-bg: #080808;
    --c-surface: #111111;
    --c-card: #161616;
    --c-border: rgba(255, 255, 255, 0.06);
    --c-border-h: rgba(255, 255, 255, 0.12);
    --c-red: #e1251b;
    --c-red-dim: rgba(225, 37, 27, 0.08);
    --c-text: #f0f0f0;
    --c-sub: #8a8a8e;
    --c-muted: #3a3a3c;
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
    margin: 0;
  }
  .empty-sub {
    font-size: 0.8rem;
    color: var(--c-sub);
    margin: 0;
  }

  /* ─────────── Layout ─────────── */
  .profile-layout {
    position: relative;
    z-index: 1;
    max-width: 1100px;
    margin: 0 auto;
    padding: 0 24px 80px;
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
    position: relative;
    padding: 56px 0 40px;
    margin-bottom: 0;
  }

  /* Subtle red-tinted backdrop at top */
  .hero-backdrop {
    position: absolute;
    top: 0;
    left: -24px;
    right: -24px;
    height: 220px;
    background: linear-gradient(
      180deg,
      rgba(225, 37, 27, 0.06) 0%,
      transparent 100%
    );
    pointer-events: none;
  }

  .hero-inner {
    position: relative;
    display: flex;
    align-items: flex-end;
    gap: 32px;
  }

  /* ── Avatar ── */
  .avatar-shell {
    position: relative;
    flex-shrink: 0;
  }

  /* Outer padding ring — creates the "halo" gap */
  .avatar-ring {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    padding: 3px;
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.18),
      rgba(255, 255, 255, 0.04)
    );
  }

  /* Inner circle — clips the image perfectly */
  .avatar-inner {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    overflow: hidden;
    background: var(--c-card);
    border: 2px solid var(--c-bg); /* gap between image and outer ring */
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .avatar-img {
    width: 100% !important;
    height: 100% !important;
  }
  .avatar-img :deep(img) {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 50%;
    display: block;
  }

  .avatar-fallback {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--c-muted);
  }

  .level-chip {
    position: absolute;
    bottom: -2px;
    left: 50%;
    transform: translateX(-50%);
    background: var(--c-surface);
    border: 1px solid var(--c-border-h);
    color: #fff;
    font-size: 0.55rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    padding: 2px 8px;
    border-radius: 20px;
    white-space: nowrap;
    text-transform: uppercase;
  }

  /* ── Identity ── */
  .hero-identity {
    flex: 1;
    min-width: 0;
    padding-bottom: 6px;
  }

  .display-name {
    font-family: var(--font-display);
    font-size: 2rem;
    font-weight: 700;
    color: #fff;
    margin: 0 0 4px;
    line-height: 1.1;
    letter-spacing: -0.03em;
  }

  .username {
    font-size: 0.85rem;
    color: var(--c-sub);
    margin: 0 0 10px;
    letter-spacing: -0.01em;
  }

  .bio {
    font-size: 0.88rem;
    color: rgba(240, 240, 240, 0.72);
    line-height: 1.65;
    margin: 0 0 18px;
    max-width: 520px;
    font-weight: 400;
  }
  .bio--empty {
    color: var(--c-muted);
    font-style: italic;
  }

  /* ── Inline hero stats (Instagram-style) ── */
  .hero-stats {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .hstat {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .hstat-val {
    font-family: var(--font-display);
    font-size: 1.1rem;
    font-weight: 700;
    color: #fff;
    letter-spacing: -0.025em;
    line-height: 1;
  }
  .hstat-lbl {
    font-size: 0.62rem;
    font-weight: 500;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--c-sub);
  }
  .hstat-sep {
    width: 1px;
    height: 30px;
    background: var(--c-border);
    flex-shrink: 0;
  }

  /* ── Buttons ── */
  .hero-actions {
    display: flex;
    gap: 8px;
    align-self: flex-end;
    padding-bottom: 6px;
    flex-shrink: 0;
  }

  .btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-family: var(--font-ui);
    font-size: 0.78rem;
    font-weight: 500;
    padding: 9px 18px;
    border-radius: 8px;
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
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(225, 37, 27, 0.3);
  }
  .btn-ghost {
    background: transparent;
    color: var(--c-text);
    border: 1px solid var(--c-border-h);
  }
  .btn-ghost:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.2);
  }

  /* ─────────── Divider ─────────── */
  .divider {
    height: 1px;
    background: var(--c-border);
    margin-bottom: 28px;
  }

  /* ─────────── Content Grid ─────────── */
  .content-grid {
    display: grid;
    grid-template-columns: 188px 1fr 188px;
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
    font-size: 0.58rem;
    font-weight: 700;
    letter-spacing: 0.09em;
    text-transform: uppercase;
    color: var(--c-muted);
    margin: 0 0 12px 8px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 9px;
    width: 100%;
    padding: 9px 10px;
    border-radius: 7px;
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
    background: rgba(255, 255, 255, 0.06);
    font-weight: 500;
  }
  .nav-count {
    margin-left: auto;
    font-size: 0.58rem;
    font-weight: 700;
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
    background: rgba(0, 0, 0, 0.72);
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
    border-radius: 14px;
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
  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }

  .btn-following {
    background: rgba(255, 255, 255, 0.06);
    color: #fff;
    border: 1px solid rgba(255, 255, 255, 0.16);
  }

  .btn-following:hover:not(:disabled) {
    background: rgba(225, 37, 27, 0.12);
    color: #ff6b6b;
    border-color: rgba(225, 37, 27, 0.32);
  }

  .btn-pending {
    background: rgba(255, 184, 0, 0.1);
    color: #ffb800;
    border: 1px solid rgba(255, 184, 0, 0.22);
  }

  /* ─────────── Responsive ─────────── */
  @media (max-width: 960px) {
    .content-grid {
      grid-template-columns: 1fr;
    }
    .nav-aside {
      border-right: none;
      border-bottom: 1px solid var(--c-border);
      padding-right: 0;
      padding-bottom: 16px;
      margin-bottom: 24px;
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
  }

  @media (max-width: 720px) {
    .profile-hero {
      padding: 40px 0 32px;
    }
    .hero-inner {
      flex-direction: column;
      align-items: flex-start;
      gap: 20px;
    }
    .hero-actions {
      width: 100%;
      align-self: auto;
      padding-bottom: 0;
    }
    .hero-actions .btn {
      flex: 1;
      justify-content: center;
    }
    .display-name {
      font-size: 1.65rem;
    }
    .avatar-ring {
      width: 96px;
      height: 96px;
    }
  }

  @media (max-width: 480px) {
    .profile-layout {
      padding: 0 16px 60px;
    }
    .display-name {
      font-size: 1.4rem;
    }
    .avatar-ring {
      width: 84px;
      height: 84px;
    }
    .hero-stats {
      gap: 14px;
    }
    .hstat-val {
      font-size: 0.95rem;
    }
  }
</style>
