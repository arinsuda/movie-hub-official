<template>
  <div class="list-page max-w-6xl mx-auto px-4 py-8">
    <!-- Page Header & Tabs -->
    <div
      class="page-header border-b border-gray-200 dark:border-gray-700 pb-4 mb-6"
    >
      <h1
        class="page-title text-3xl font-bold text-gray-900 dark:text-white mb-4"
      >
        My Profile
      </h1>
      <div class="tabs flex space-x-2 overflow-x-auto">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="tab-btn px-4 py-2 text-sm font-medium rounded-lg transition-colors whitespace-nowrap"
          :class="
            activeTab === tab.key
              ? 'bg-primary text-white bg-blue-600'
              : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800'
          "
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Tab Contents -->
    <div class="tab-content">
      <div v-if="loading" class="flex justify-center py-12">
        <span class="text-gray-500">กำลังโหลดข้อมูล...</span>
      </div>

      <template v-else-if="user">
        <!-- เงื่อนไขการแสดงผลแต่ละ Tab โดยดึง Component มาใช้ -->
        <ProfileInfo v-if="activeTab === 'profile'" :user="user" />
        <UserReviews v-else-if="activeTab === 'reviews'" :user-id="userId" />
        <UserWatchlist
          v-else-if="activeTab === 'watchlist'"
          :user-id="userId"
        />
        <UserLikes v-else-if="activeTab === 'likes'" :user-id="userId" />
        <UserAchievements
          v-else-if="activeTab === 'achievements'"
          :user-id="userId"
        />
      </template>

      <div v-else class="text-center py-12 text-gray-500">
        ไม่พบข้อมูลผู้ใช้งาน กรุณาเข้าสู่ระบบใหม่อีกครั้ง
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { authApi } from "@/api/api"
  import { useAuthStore } from "@/stores/auth"
  import type { UserProfile } from "@/types"
  import { onMounted, ref } from "vue"

  import ProfileInfo from "@/components/profile/ProfileInfo.vue"
  import UserReviews from "@/components/profile/UserReviews.vue"
  import UserWatchlist from "@/components/profile/UserWatchlist.vue"
  import UserLikes from "@/components/profile/UserLikes.vue"
  import UserAchievements from "@/components/profile/UserAchievements.vue"

  const auth = useAuthStore()
  const userId = auth.user?.id || 0
  const user = ref<UserProfile | null>(null)
  const loading = ref<boolean>(true)

  type TabKey = "profile" | "reviews" | "watchlist" | "likes" | "achievements"
  const activeTab = ref<TabKey>("profile")
  const tabs: { key: TabKey; label: string }[] = [
    { key: "profile", label: "ข้อมูลส่วนตัว" },
    { key: "reviews", label: "รีวิวของฉัน" },
    { key: "watchlist", label: "รายการเฝ้าดู" },
    { key: "likes", label: "รายการที่ชอบ" },
    { key: "achievements", label: "ความสำเร็จ" },
  ]

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
