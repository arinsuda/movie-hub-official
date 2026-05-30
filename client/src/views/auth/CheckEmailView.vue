<template>
  <div class="min-h-screen bg-[#0f0f0f] flex items-center justify-center px-4">
    <div class="w-full max-w-md text-center">
      <div class="text-6xl mb-6">📬</div>
      <h1 class="text-2xl font-bold text-white mb-2">ตรวจสอบอีเมลของคุณ</h1>
      <p class="text-gray-400 mb-2">เราได้ส่งลิงก์ยืนยันไปที่</p>
      <p class="text-white font-semibold mb-8">{{ email || "อีเมลของคุณ" }}</p>
      <p class="text-gray-500 text-sm mb-8">
        ไม่ได้รับอีเมล? ตรวจสอบโฟลเดอร์ Spam<br />หรือขอส่งใหม่ด้านล่าง
      </p>
      <button
        class="inline-block px-8 py-3 bg-[#e50914] text-white font-bold rounded-lg hover:bg-[#f40612] transition-colors disabled:opacity-50"
        :disabled="resendStatus === 'sending' || resendStatus === 'sent'"
        @click="resend"
      >
        <span v-if="resendStatus === 'idle'">ส่งลิงก์ใหม่</span>
        <span v-else-if="resendStatus === 'sending'">กำลังส่ง...</span>
        <span v-else-if="resendStatus === 'sent'">ส่งแล้ว! ตรวจสอบอีเมล</span>
        <span v-else>เกิดข้อผิดพลาด ลองใหม่</span>
      </button>
      <p class="mt-6">
        <RouterLink
          to="/login"
          class="text-sm text-gray-500 hover:text-white transition-colors"
        >
          กลับไปหน้าเข้าสู่ระบบ
        </RouterLink>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from "vue"
  import { RouterLink } from "vue-router"
  import { useAuthStore } from "@/stores/auth"
  import { authApi } from "@/api/api"

  const authStore = useAuthStore()
  const email = ref(authStore.user?.email ?? "")

  type ResendStatus = "idle" | "sending" | "sent" | "failed"
  const resendStatus = ref<ResendStatus>("idle")

  async function resend() {
    if (!email.value) {
      const input = window.prompt("กรุณากรอกอีเมลของคุณ")
      if (!input) return
      email.value = input
    }

    resendStatus.value = "sending"
    try {
      await authApi.resendVerification(email.value)
      resendStatus.value = "sent"
    } catch {
      resendStatus.value = "failed"
    }
  }
</script>
