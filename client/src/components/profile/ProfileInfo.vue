<template>
  <div
    class="bg-gray-950/40 backdrop-blur-md rounded-2xl p-6 border border-gray-800/60 shadow-2xl relative overflow-hidden"
  >
    <div
      class="absolute -top-24 -left-24 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl pointer-events-none"
    ></div>
    <div
      class="absolute -bottom-24 -right-24 w-48 h-48 bg-purple-500/10 rounded-full blur-3xl pointer-events-none"
    ></div>

    <div
      class="flex flex-col md:flex-row items-center md:items-start gap-8 relative z-10"
    >
      <div class="flex flex-col items-center space-y-4 flex-shrink-0">
        <div
          class="relative w-32 h-32 rounded-full overflow-hidden border-2 border-gray-700/50 bg-gray-900 flex items-center justify-center text-gray-400 group shadow-inner"
        >
          <img
            v-if="user.avatar_url"
            :src="user.avatar_url"
            alt="User Avatar"
            class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-110"
          />
          <span
            v-else
            class="text-4xl font-bold tracking-wider bg-gradient-to-tr from-blue-400 to-purple-400 bg-clip-text text-transparent"
          >
            {{ user.username?.charAt(0).toUpperCase() }}
          </span>
          <div
            class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center cursor-pointer"
          >
            <svg
              class="w-6 h-6 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
              />
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
          </div>
        </div>
        <button
          class="text-xs font-medium text-gray-400 hover:text-white transition-colors tracking-wide bg-gray-900/80 px-3 py-1.5 rounded-full border border-gray-800"
        >
          เปลี่ยนรูปโปรไฟล์
        </button>
      </div>

      <div class="flex-1 w-full space-y-6">
        <h2
          class="text-lg font-bold uppercase tracking-wider text-white border-b border-gray-800/80 pb-3"
        >
          ACCOUNT DETAILS
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label
              class="block text-xs font-semibold uppercase tracking-wider text-gray-400 mb-2"
              >Username</label
            >
            <div class="relative">
              <input
                type="text"
                v-model="formData.username"
                disabled
                class="w-full px-4 py-2.5 bg-gray-900/60 border border-gray-800/80 rounded-xl text-gray-500 cursor-not-allowed font-medium select-none text-sm"
              />
              <span class="absolute right-3 top-3 text-gray-600"> 🔒 </span>
            </div>
          </div>

          <div>
            <label
              class="block text-xs font-semibold uppercase tracking-wider text-gray-400 mb-2"
              >Display Name</label
            >
            <input
              type="text"
              v-model="formData.displayName"
              class="w-full px-4 py-2.5 bg-gray-900/40 border border-gray-800 rounded-xl text-white font-medium focus:border-blue-500/80 focus:bg-gray-900/90 focus:ring-4 focus:ring-blue-500/10 outline-none transition-all text-sm placeholder-gray-600"
              placeholder="Enter your public display name..."
            />
          </div>
        </div>

        <div class="pt-4 flex justify-end">
          <button
            @click="handleUpdateProfile"
            :disabled="isSaving"
            class="px-6 py-2.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 active:scale-95 text-white font-semibold rounded-xl text-xs uppercase tracking-widest shadow-lg shadow-blue-500/20 transition-all disabled:opacity-40 disabled:pointer-events-none"
          >
            {{ isSaving ? "Saving changes..." : "Save Changes" }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from "vue"
  import type { UserProfile } from "@/types"

  const props = defineProps<{
    user: UserProfile
  }>()

  const isSaving = ref(false)

  const formData = ref({
    username: props.user.username || "",
    displayName: props.user.display_name || "",
  })

  const handleUpdateProfile = async () => {
    try {
      isSaving.value = true
      await new Promise(resolve => setTimeout(resolve, 800))
      alert("บันทึกข้อมูลส่วนตัวเรียบร้อยแล้ว!")
    } catch (err) {
      console.error("Update profile failed:", err)
      alert("เกิดข้อผิดพลาดในการบันทึกข้อมูล")
    } finally {
      isSaving.value = false
    }
  }
</script>
