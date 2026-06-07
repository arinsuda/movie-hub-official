<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center border-b pb-2 mb-4">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
        รายการที่ชอบ ({{ likedItems.length }})
      </h2>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-500">
      กำลังโหลดรายการที่ชอบ...
    </div>

    <div
      v-else-if="likedItems.length === 0"
      class="bg-white dark:bg-gray-800 rounded-xl p-8 text-center border border-gray-100 dark:border-gray-700"
    >
      <p class="text-gray-500 dark:text-gray-400">คุณยังไม่มีรายการที่ถูกใจ</p>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="item in likedItems"
        :key="item.id"
        class="bg-white dark:bg-gray-800 p-4 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 flex items-center justify-between transition-all hover:border-gray-200 dark:hover:border-gray-600"
      >
        <div class="flex items-center space-x-4 min-w-0">
          <div
            class="p-3 rounded-lg bg-red-50 dark:bg-red-950/30 text-red-500 flex-shrink-0"
          >
            <svg class="w-6 h-6 fill-current" viewBox="0 0 20 20">
              <path
                d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z"
              />
            </svg>
          </div>

          <div class="min-w-0">
            <h3
              class="font-medium text-sm sm:text-base text-gray-900 dark:text-white truncate hover:text-blue-600 cursor-pointer"
            >
              {{ item.title }}
            </h3>
            <p class="text-xs text-gray-400 mt-0.5">
              โดย {{ item.author }} • ถูกใจเมื่อ {{ item.likedAt }}
            </p>
          </div>
        </div>

        <button
          @click="handleUnlike(item.id)"
          class="flex-shrink-0 ml-4 px-3 py-1.5 border border-red-200 text-red-500 hover:bg-red-50 dark:border-red-900/50 dark:hover:bg-red-950/20 rounded-lg text-xs font-medium transition-colors"
          title="ยกเลิกการชอบ"
        >
          เลิกชอบ
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref } from "vue"

  const props = defineProps<{
    userId: number
  }>()

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
      // TODO: เรียกใช้งาน API สำหรับดึงรายการ Like ของ User เช่น const res = await likesApi.getLikes(props.userId)
      await new Promise(resolve => setTimeout(resolve, 500))

      // ข้อมูล Mock สำหรับแสดงผลหน้ารายการที่ชอบ
      likedItems.value = [
        {
          id: 201,
          title: "How to Optimize SQL Queries in MySQL for Production",
          author: "Tech Blogger",
          likedAt: "2026-06-03",
        },
        {
          id: 202,
          title: "State Management in Flutter: Comprehensive Guide to Riverpod",
          author: "Dev community",
          likedAt: "2026-05-29",
        },
        {
          id: 203,
          title: "Designing Scalable ERP System Database Models",
          author: "Database Architect",
          likedAt: "2026-05-12",
        },
      ]
    } catch (err) {
      console.error("Fetch liked items failed:", err)
    } finally {
      loading.value = false
    }
  })

  const handleUnlike = (id: number) => {
    // นำรายการออกจาก UI ทันที
    likedItems.value = likedItems.value.filter(item => item.id !== id)
    console.log("Unliked item ID:", id)
    // TODO: เรียกใช้งาน API เพื่อบันทึกลงฐานข้อมูล เช่น likesApi.unlike(id)
  }
</script>
