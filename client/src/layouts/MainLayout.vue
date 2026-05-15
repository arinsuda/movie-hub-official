<template>
  <div class="app-shell">
    <!-- ── Navbar ──────────────────────────────────────────── -->
    <header class="navbar" :class="{ 'navbar--scrolled': scrolled }">
      <RouterLink to="/" class="nav-logo">
        <span class="logo-movie">MOVIE</span><span class="logo-hub">HUB</span>
      </RouterLink>

      <nav class="nav-links">
        <RouterLink
          to="/"
          class="nav-link"
          :class="{ active: route.name === 'home' }"
          >HOME</RouterLink
        >
        <RouterLink
          to="/movies"
          class="nav-link"
          :class="{ active: route.path.startsWith('/movies') }"
          >MOVIES</RouterLink
        >
        <RouterLink
          to="/tv"
          class="nav-link"
          :class="{ active: route.path.startsWith('/tv') }"
          >TV SERIES</RouterLink
        >
      </nav>

      <div class="nav-right">
        <!-- Search -->
        <div class="search-box" :class="{ 'search-box--open': searchOpen }">
          <button class="search-icon-btn" @click="toggleSearch">
            <Search :size="18" />
          </button>
          <input
            v-if="searchOpen"
            ref="searchInput"
            v-model="searchQuery"
            class="search-input"
            placeholder="Search movies, series..."
            @keydown.enter="doSearch"
            @keydown.escape="closeSearch"
          />
        </div>

        <!-- Notification bell -->
        <button class="icon-btn">
          <Bell :size="18" />
        </button>

        <!-- User dropdown -->
        <div class="user-menu" v-if="authStore.user" ref="userMenuRef">
          <button class="user-trigger" @click="userMenuOpen = !userMenuOpen">
            <div class="user-avatar">
              <img
                v-if="authStore.user.avatar_url"
                :src="authStore.user.avatar_url"
                :alt="authStore.user.username"
              />
              <UserIcon v-else :size="16" />
            </div>
            <span class="user-name">{{ authStore.user.username }}</span>
            <ChevronDown
              :size="14"
              :class="{ 'rotate-180': userMenuOpen }"
              class="transition-transform duration-200"
            />
          </button>

          <Transition name="dropdown">
            <div class="dropdown-menu" v-if="userMenuOpen">
              <RouterLink
                :to="`/users/${authStore.user.id}`"
                class="dropdown-item"
                @click="userMenuOpen = false"
              >
                <UserIcon :size="14" />Profile
              </RouterLink>
              <RouterLink
                :to="`/users/${authStore.user.id}/library`"
                class="dropdown-item"
                @click="userMenuOpen = false"
              >
                <BookMarked :size="14" />My Library
              </RouterLink>
              <RouterLink
                to="/feed"
                class="dropdown-item"
                @click="userMenuOpen = false"
              >
                <Rss :size="14" />Feed
              </RouterLink>
              <div class="dropdown-divider" />
              <button
                class="dropdown-item dropdown-item--danger"
                @click="handleLogout"
              >
                <LogOut :size="14" />Log out
              </button>
            </div>
          </Transition>
        </div>
      </div>
    </header>

    <!-- ── Page content ────────────────────────────────────── -->
    <main class="page-content">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import {
  Search,
  Bell,
  User as UserIcon,
  ChevronDown,
  BookMarked,
  Rss,
  LogOut,
} from "lucide-vue-next";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const scrolled = ref(false);
const searchOpen = ref(false);
const searchQuery = ref("");
const searchInput = ref<HTMLInputElement | null>(null);
const userMenuOpen = ref(false);
const userMenuRef = ref<HTMLElement | null>(null);

function onScroll() {
  scrolled.value = window.scrollY > 20;
}

function toggleSearch() {
  searchOpen.value = !searchOpen.value;
  if (searchOpen.value) {
    nextTick(() => searchInput.value?.focus());
  }
}

function closeSearch() {
  searchOpen.value = false;
  searchQuery.value = "";
}

function doSearch() {
  if (!searchQuery.value.trim()) return;
  router.push({ name: "movies", query: { q: searchQuery.value.trim() } });
  closeSearch();
}

function onClickOutside(e: MouseEvent) {
  if (userMenuRef.value && !userMenuRef.value.contains(e.target as Node)) {
    userMenuOpen.value = false;
  }
}

async function handleLogout() {
  userMenuOpen.value = false;
  await authStore.logout();
  router.push({ name: "login" });
}

onMounted(() => {
  window.addEventListener("scroll", onScroll, { passive: true });
  document.addEventListener("click", onClickOutside);
});
onUnmounted(() => {
  window.removeEventListener("scroll", onScroll);
  document.removeEventListener("click", onClickOutside);
});
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  background: #141414;
}

/* ── Navbar ─────────────────────────────────────────── */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 2rem;
  padding: 0 2rem;
  height: 64px;
  transition:
    background 0.3s,
    backdrop-filter 0.3s,
    border-color 0.3s;
  border-bottom: 1px solid transparent;
}
.navbar--scrolled {
  background: rgba(20, 20, 20, 0.92);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-color: rgba(255, 255, 255, 0.06);
}

.nav-logo {
  font-family: "Impact", "Arial Black", sans-serif;
  font-size: 1.5rem;
  font-weight: 900;
  text-decoration: none;
  letter-spacing: -0.5px;
  flex-shrink: 0;
}
.logo-movie {
  color: #ffffff;
}
.logo-hub {
  color: #e50914;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}
.nav-link {
  padding: 0.4rem 0.75rem;
  font-size: 0.82rem;
  font-weight: 600;
  letter-spacing: 0.8px;
  color: #a3a3a3;
  text-decoration: none;
  border-radius: 6px;
  transition:
    color 0.2s,
    background 0.2s;
}
.nav-link:hover {
  color: #fff;
}
.nav-link.active {
  color: #fff;
}

.nav-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Search */
.search-box {
  display: flex;
  align-items: center;
  background: transparent;
  border-radius: 9999px;
  overflow: hidden;
  transition: background 0.2s;
}
.search-box--open {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
}
.search-icon-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #a3a3a3;
  display: flex;
  padding: 0.5rem 0.6rem;
  transition: color 0.2s;
}
.search-icon-btn:hover {
  color: #fff;
}
.search-input {
  background: none;
  border: none;
  outline: none;
  color: #fff;
  font-size: 0.875rem;
  width: 200px;
  padding: 0.5rem 0.75rem 0.5rem 0;
}
.search-input::placeholder {
  color: #666;
}

/* Icon button */
.icon-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #a3a3a3;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  transition:
    color 0.2s,
    background 0.2s;
}
.icon-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.08);
}

/* User menu */
.user-menu {
  position: relative;
}
.user-trigger {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 9999px;
  padding: 0.3rem 0.75rem 0.3rem 0.3rem;
  cursor: pointer;
  color: #fff;
  transition: background 0.2s;
}
.user-trigger:hover {
  background: rgba(255, 255, 255, 0.12);
}
.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #2a2a2a;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
}
.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.user-name {
  font-size: 0.85rem;
  font-weight: 500;
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 160px;
  background: #1f1f1f;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 0.4rem;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.6);
}
.dropdown-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.55rem 0.75rem;
  font-size: 0.875rem;
  color: #d1d1d1;
  text-decoration: none;
  border-radius: 6px;
  border: none;
  background: none;
  width: 100%;
  cursor: pointer;
  transition:
    background 0.15s,
    color 0.15s;
  text-align: left;
}
.dropdown-item:hover {
  background: rgba(255, 255, 255, 0.07);
  color: #fff;
}
.dropdown-item--danger:hover {
  background: rgba(229, 9, 20, 0.12);
  color: #e50914;
}
.dropdown-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.07);
  margin: 0.3rem 0;
}

/* Dropdown transition */
.dropdown-enter-active,
.dropdown-leave-active {
  transition:
    opacity 0.15s,
    transform 0.15s;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* ── Page content ──────────────────────────────────── */
.page-content {
  padding-top: 64px;
}
</style>
