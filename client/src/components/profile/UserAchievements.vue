<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center border-b pb-2 mb-4">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
        ความสำเร็จของฉัน
      </h2>
      <span
        class="text-sm bg-blue-50 text-blue-600 dark:bg-blue-950/40 dark:text-blue-400 px-3 py-1 rounded-full font-medium"
      >
        ปลดล็อกแล้ว {{ unlockedCount }} / {{ achievements.length }}
      </span>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-500">
      กำลังโหลดข้อมูลความสำเร็จ...
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      <div
        v-for="item in achievements"
        :key="item.id"
        class="bg-white dark:bg-gray-800 p-5 rounded-xl border shadow-sm transition-all flex items-start space-x-4"
        :class="
          item.isUnlocked
            ? 'border-blue-100 dark:border-gray-700 hover:shadow-md'
            : 'border-gray-100 dark:border-gray-800 opacity-60 bg-gray-50/50 dark:bg-gray-800/50'
        "
      >
        <div
          class="p-3 rounded-xl flex-shrink-0 text-2xl"
          :class="
            item.isUnlocked
              ? item.bgColorClass
              : 'bg-gray-200 text-gray-400 dark:bg-gray-700'
          "
        >
          <span>{{ item.icon }}</span>
        </div>

        <div class="space-y-1 min-w-0">
          <div class="flex items-center space-x-2">
            <h3
              class="font-semibold text-sm sm:text-base text-gray-900 dark:text-white truncate"
            >
              {{ item.title }}
            </h3>
            <svg
              v-if="item.isUnlocked"
              class="w-4 h-4 text-green-500 flex-shrink-0"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clip-rule="evenodd"
              />
            </svg>
          </div>

          <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2">
            {{ item.description }}
          </p>

          <p
            v-if="item.isUnlocked && item.unlockedAt"
            class="text-[10px] text-gray-400 pt-1"
          >
            ปลดล็อกเมื่อ {{ item.unlockedAt }}
          </p>
          <p
            v-else-if="!item.isUnlocked"
            class="text-[10px] text-amber-600 dark:text-amber-400 font-medium pt-1"
          >
            🔒 ยังไม่บรรลุเป้าหมาย
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref, computed } from "vue"

  const props = defineProps<{
    userId: number
  }>()

  const loading = ref(false)

  interface Achievement {
    id: number
    title: string
    description: string
    icon: string
    bgColorClass: string
    isUnlocked: boolean
    unlockedAt?: string
  }

  const achievements = ref<Achievement[]>([])

  // คำนวณจำนวน Achievement ที่ปลดล็อกแล้ว
  const unlockedCount = computed(() => {
    return achievements.value.filter(item => item.isUnlocked).length
  })

  onMounted(async () => {
    try {
      loading.value = true
      // TODO: เรียกใช้งาน API เช่น const res = await achievementApi.getAchievements(props.userId)
      await new Promise(resolve => setTimeout(resolve, 500))

      // ข้อมูล Mock ระบบความสำเร็จ
      achievements.value = [
        {
          id: 1,
          title: "สมาชิกหน้าใหม่",
          description: "สมัครสมาชิกและกรอกข้อมูลโปรไฟล์เสร็จสมบูรณ์",
          icon: "🎉",
          bgColorClass: "bg-blue-50 text-blue-600 dark:bg-blue-950/50",
          isUnlocked: true,
          unlockedAt: "2026-01-10",
        },
        {
          id: 2,
          title: "นักวิจารณ์มือโปร",
          description: "เขียนรีวิวคอนเทนต์หรือสินค้าสะสมครบ 5 ครั้ง",
          icon: "✍️",
          bgColorClass: "bg-purple-50 text-purple-600 dark:bg-purple-950/50",
          isUnlocked: true,
          unlockedAt: "2026-05-15",
        },
        {
          id: 3,
          title: "นักส่องยามค่ำคืน",
          description: "เพิ่มรายการเฝ้าดู (Watchlist) ครบ 10 รายการ",
          icon: "👀",
          bgColorClass: "bg-amber-50 text-amber-600 dark:bg-amber-950/50",
          isUnlocked: false,
        },
        {
          id: 4,
          title: "สายซัพพอร์ต",
          description: "กดถูกใจเนื้อหาที่ชื่นชอบครบ 20 ครั้ง",
          icon: "❤️",
          bgColorClass: "bg-red-50 text-red-600 dark:bg-red-950/50",
          isUnlocked: true,
          unlockedAt: "2026-06-01",
        },
        {
          id: 5,
          title: "ผู้ใช้ระดับ Elite",
          description: "ล็อกอินเข้าใช้งานต่อเนื่องกันครบ 30 วัน",
          icon: "🔥",
          bgColorClass: "bg-orange-50 text-orange-600 dark:bg-orange-950/50",
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
