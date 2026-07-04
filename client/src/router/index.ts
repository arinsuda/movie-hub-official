import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

function hasGenres(genres: string | null | undefined): boolean {
  if (!genres) return false;
  if (genres === "null") return false;
  try {
    const parsed = JSON.parse(genres);
    return Array.isArray(parsed) && parsed.length > 0;
  } catch {
    return false;
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior: () => ({ top: 0 }),
  routes: [
    // ── Auth ──────────────────────────────────────────────────
    {
      path: "/login",
      name: "login",
      component: () => import("@/views/auth/LoginView.vue"),
      meta: { guestOnly: true },
    },
    {
      path: "/register",
      name: "register",
      component: () => import("@/views/auth/RegisterView.vue"),
      meta: { guestOnly: true },
    },
    {
      path: "/check-email",
      name: "check-email",
      component: () => import("@/views/auth/CheckEmailView.vue"),
    },
    {
      path: "/auth/verify-email",
      name: "verify-email",
      component: () => import("@/views/auth/VerifyEmailView.vue"),
    },
    {
      path: "/forgot-password",
      name: "forgot-password",
      component: () => import("@/views/auth/ForgotPasswordView.vue"),
      meta: { guestOnly: true },
    },
    {
      path: "/reset-password",
      name: "reset-password",
      component: () => import("@/views/auth/ResetPasswordView.vue"),
      meta: { guestOnly: true },
    },

    // ── Onboarding (protected, no MainLayout/Navbar) ───────────
    {
      path: "/onboarding",
      name: "onboarding",
      component: () => import("@/views/onboarding/FavoriteGenreView.vue"),
      meta: { requiresAuth: true },
    },

    // ── Main (protected, with MainLayout + Navbar) ─────────────
    {
      path: "/",
      component: () => import("@/layouts/MainLayout.vue"),
      meta: { requiresAuth: true },
      children: [
        {
          path: "",
          name: "home",
          component: () => import("@/views/HomeView.vue"),
        },
        {
          path: "search",
          name: "search-results",
          component: () => import("@/views/SearchResultView.vue"),
        },
        {
          path: "movies",
          name: "movies",
          component: () => import("@/views/movie/MoviesView.vue"),
        },
        {
          path: "movies/:id",
          name: "movie-detail",
          component: () => import("@/views/movie/MovieDetailView.vue"),
        },
        {
          path: "tv",
          name: "tv",
          component: () => import("@/views/tv/TVSeriesView.vue"),
        },
        {
          path: "tv/:id",
          name: "tv-detail",
          component: () => import("@/views/tv/TVDetailView.vue"),
        },
        {
          path: "upcoming",
          name: "upcoming",
          component: () => import("@/views/movie/Upcoming.vue"),
        },
        {
          path: "about-us",
          name: "about-us",
          component: () => import("@/views/aboutus/AboutUsView.vue"),
        },
        {
          path: "donate",
          name: "donate",
          component: () => import("@/views/aboutus/DonateView.vue"),
        },
        {
          path: "users/:userId",
          name: "user-profile",
          component: () => import("@/views/user/UserProfileView.vue"),
        },
        {
          path: "users/:userId/library",
          name: "user-library",
          component: () => import("@/views/library/LibraryView.vue"),
        },
        {
          path: "users/:userId/achievements",
          name: "user-achievements",
          component: () => import("@/views/achievement/AchievementsView.vue"),
        },
        {
          path: "feed",
          name: "feed",
          component: () => import("@/views/FeedView.vue"),
        },
      ],
    },

    // ── 404 ───────────────────────────────────────────────────
    {
      path: "/:pathMatch(.*)*",
      name: "not-found",
      component: () => import("@/views/NotFoundView.vue"),
    },
  ],
});

router.beforeEach(async (to) => {
  const authStore = useAuthStore();

  if (!authStore.user && to.meta.requiresAuth) {
    try {
      await authStore.fetchMe();
    } catch (error) {
      console.error(error);
    }
  }

  // ยังไม่ login → ไป login
  if (to.meta.requiresAuth && !authStore.isLoggedIn)
    return { name: "login", query: { redirect: to.fullPath } };

  // login แล้วเข้า guest-only page → ไป home
  if (to.meta.guestOnly && authStore.isLoggedIn) return { name: "home" };

  // ✅ ต้อง requiresAuth ด้วย ไม่งั้น guest route ก็โดน redirect
  if (
    to.meta.requiresAuth &&
    authStore.isLoggedIn &&
    !hasGenres(authStore.user?.favorite_genres) &&
    to.name !== "onboarding"
  ) {
    return { name: "onboarding" };
  }

  // ✅ ใช้ hasGenres ให้สอดคล้องกัน
  if (
    authStore.isLoggedIn &&
    hasGenres(authStore.user?.favorite_genres) &&
    to.name === "onboarding"
  ) {
    return { name: "home" };
  }
});

export default router;
