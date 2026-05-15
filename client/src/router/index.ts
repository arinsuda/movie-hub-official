import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

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
      path: "/verify-email",
      name: "verify-email",
      component: () => import("@/views/auth/VerifyEmailView.vue"),
    },

    // ── Main (protected) ──────────────────────────────────────
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
  if (to.meta.requiresAuth && !authStore.isLoggedIn)
    return { name: "login", query: { redirect: to.fullPath } };
  if (to.meta.guestOnly && authStore.isLoggedIn) return { name: "home" };
});

export default router;
