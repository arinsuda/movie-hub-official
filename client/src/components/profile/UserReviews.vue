<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center border-b pb-2 mb-4">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
        รีวิวของฉัน ({{ reviews.length }})
      </h2>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-500">
      กำลังโหลดรายการรีวิว...
    </div>

    <div
      v-else-if="reviews.length === 0"
      class="bg-white dark:bg-gray-800 rounded-xl p-8 text-center border border-gray-100 dark:border-gray-700"
    >
      <p class="text-gray-500 dark:text-gray-400">คุณยังไม่เคยเขียนรีวิวเลย</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-4">
      <div
        v-for="review in reviews"
        :key="review.id"
        class="bg-white dark:bg-gray-800 p-5 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col justify-between"
      >
        <div>
          <div class="flex justify-between items-start mb-2">
            <div>
              <h3
                class="font-semibold text-gray-900 dark:text-white hover:text-blue-600 cursor-pointer"
              >
                {{ review.targetName }}
              </h3>
              <span class="text-xs text-gray-400">{{ review.createdAt }}</span>
            </div>

            <div class="flex items-center space-x-1 text-amber-500">
              <span class="text-sm font-bold">{{ review.rating }}.0</span>
              <svg class="w-4 h-4 fill-current" viewBox="0 0 20 20">
                <path
                  d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z"
                />
              </svg>
            </div>
          </div>

          <p class="text-gray-600 dark:text-gray-300 text-sm line-clamp-3">
            {{ review.content }}
          </p>
        </div>

        <div
          class="flex justify-end space-x-3 pt-4 mt-2 border-t border-gray-50 dark:border-gray-700 text-xs font-medium"
        >
          <button
            @click="handleEdit(review.id)"
            class="text-blue-600 dark:text-blue-400 hover:underline"
          >
            แก้ไข
          </button>
          <button
            @click="handleDelete(review.id)"
            class="text-red-500 hover:underline"
          >
            ลบ
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref } from "vue"

  // รับ userId มาจากหน้าหลัก เพื่อเอาไป query รีวิวของ user คนนี้
  const props = defineProps<{
    userId: number
  }>()

  const loading = ref(false)

  // จำลอง Interface ของ Review
  interface ReviewItem {
    id: number
    targetName: string
    rating: number
    content: string
    createdAt: string
  }

  const reviews = ref<ReviewItem[]>([])

  onMounted(async () => {
    try {
      loading.value = true
      // TODO: เรียกใช้งาน API จริง เช่น const res = await reviewApi.getUserReviews(props.userId)
      // ด้านล่างนี้คือข้อมูล Mock สำหรับทดสอบ UI
      await new Promise(resolve => setTimeout(resolve, 500))
      reviews.value = [
        {
          id: 1,
          targetName: "Sample Project / Product A",
          rating: 5,
          content:
            "ระบบใช้งานได้ดีมากครับ ตอบโจทย์การทำงาน คล่องตัวสุดๆ แนะนำเลยสำหรับคนที่กำลังมองหา Solution ด้านนี้",
          createdAt: "2026-05-15",
        },
        {
          id: 2,
          targetName: "Course: Advanced Go Backend",
          rating: 4,
          content:
            "เนื้อหาแน่นมาก ได้เข้าใจเรื่อง Concurrency กับการทำ Clean Architecture แบบจริงจัง แต่หักคะแนนเสียงไมค์นิดนึงครับ",
          createdAt: "2026-04-20",
        },
      ]
    } catch (err) {
      console.error("Fetch reviews failed:", err)
    } finally {
      loading.value = false
    }
  })

  const handleEdit = (id: number) => {
    console.log("Edit review ID:", id)
    // เพิ่ม Logic เปิด Modal แก้ไขรีวิว
  }

  const handleDelete = (id: number) => {
    if (confirm("คุณแน่ใจใช่ไหมที่จะลบรีวิวนี้?")) {
      reviews.value = reviews.value.filter(item => item.id !== id)
      console.log("Deleted review ID:", id)
    }
  }
</script>
