<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center border-b pb-2 mb-4">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
        รายการเฝ้าดูของฉัน ({{ watchlist.length }})
      </h2>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-500">
      กำลังโหลดรายการเฝ้าดู...
    </div>

    <div
      v-else-if="watchlist.length === 0"
      class="bg-white dark:bg-gray-800 rounded-xl p-8 text-center border border-gray-100 dark:border-gray-700"
    >
      <p class="text-gray-500 dark:text-gray-400">
        ยังไม่มีรายการในวอทช์ลิสต์ของคุณ
      </p>
    </div>

    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
      <div
        v-for="item in watchlist"
        :key="item.id"
        class="group bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col justify-between transition-all hover:shadow-md"
      >
        <div
          class="relative aspect-[3/4] bg-gray-100 dark:bg-gray-700 overflow-hidden"
        >
          <img
            v-if="item.coverUrl"
            :src="item.coverUrl"
            :alt="item.title"
            class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
          <div
            v-else
            class="w-full h-full flex items-center justify-center text-gray-400 text-xs text-center p-2"
          >
            No Image
          </div>

          <span
            class="absolute top-2 left-2 bg-black/60 backdrop-blur-sm text-white text-[10px] px-2 py-0.5 rounded font-medium"
          >
            {{ item.category }}
          </span>
        </div>

        <div class="p-3 flex-1 flex flex-col justify-between">
          <h3
            class="font-medium text-sm text-gray-900 dark:text-white line-clamp-2 mb-2 group-hover:text-blue-600 cursor-pointer"
          >
            {{ item.title }}
          </h3>

          <div
            class="flex items-center justify-between pt-2 border-t border-gray-50 dark:border-gray-700"
          >
            <span class="text-xs text-gray-400"
              >เพิ่มเมื่อ {{ item.addedAt }}</span
            >

            <button
              @click="handleRemove(item.id)"
              class="text-gray-400 hover:text-red-500 transition-colors p-1 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700"
              title="ลบออกจากรายการ"
            >
              <svg
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-4v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
            </button>
          </div>
        </div>
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

  interface WatchlistItem {
    id: number
    title: string
    category: string
    coverUrl?: string
    addedAt: string
  }

  const watchlist = ref<WatchlistItem[]>([])

  onMounted(async () => {
    try {
      loading.value = true
      // TODO: เรียกใช้ API จริง เช่น const res = await watchlistApi.getWatchlist(props.userId)
      await new Promise(resolve => setTimeout(resolve, 500))

      // ข้อมูลจำลองสำหรับจัด UI Layout
      watchlist.value = [
        {
          id: 101,
          title: "Introduction to Nuxt 3 Framework",
          category: "Course",
          coverUrl:
            "https://images.unsplash.com/photo-1516116211223-5c359a36298a?w=400&auto=format&fit=crop&q=60",
          addedAt: "2026-06-01",
        },
        {
          id: 102,
          title: "Tailwind CSS Component Architecture",
          category: "Article",
          coverUrl:
            "https://images.unsplash.com/photo-1507238691740-187a5b1d37b8?w=400&auto=format&fit=crop&q=60",
          addedAt: "2026-05-28",
        },
        {
          id: 103,
          title: "Building Microservices with Go (Fiber)",
          category: "Video",
          coverUrl:
            "https://images.unsplash.com/photo-1618401471353-b98aedd07871?w=400&auto=format&fit=crop&q=60",
          addedAt: "2026-05-20",
        },
      ]
    } catch (err) {
      console.error("Fetch watchlist failed:", err)
    } finally {
      loading.value = false
    }
  })

  const handleRemove = (id: number) => {
    // กรองตัวที่ถูกลบออกจาก State ทันที
    watchlist.value = watchlist.value.filter(item => item.id !== id)
    console.log("Removed from watchlist, ID:", id)
    // TODO: เรียก API ลบออก เช่น watchlistApi.remove(id)
  }
</script>
