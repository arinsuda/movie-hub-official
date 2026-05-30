<template>
  <div class="min-h-screen bg-[#0f0f0f] flex items-center justify-center px-4">
    <div class="w-full max-w-md text-center">
      <!-- Loading -->
      <template v-if="status === 'loading'">
        <div class="animate-spin text-5xl mb-6">⏳</div>
        <p class="text-gray-400">กำลังยืนยันอีเมล...</p>
      </template>

      <!-- Success -->
      <template v-else-if="status === 'success'">
        <div class="text-6xl mb-6">🎬</div>
        <h1 class="text-2xl font-bold text-white mb-2">ยืนยันอีเมลสำเร็จ!</h1>
        <p class="text-gray-400 mb-8">ตอนนี้คุณสามารถเข้าสู่ระบบได้แล้ว</p>
        <RouterLink
          to="/login"
          class="inline-block px-8 py-3 bg-[#e50914] text-white font-bold rounded-lg hover:bg-[#f40612] transition-colors"
        >
          เข้าสู่ระบบ
        </RouterLink>
      </template>

      <!-- Error -->
      <template v-else-if="status === 'error'">
        <div class="text-6xl mb-6">❌</div>
        <h1 class="text-2xl font-bold text-white mb-2">
          ลิงก์ไม่ถูกต้องหรือหมดอายุ
        </h1>
        <p class="text-gray-400 mb-8">ลองขอลิงก์ยืนยันอีเมลใหม่อีกครั้ง</p>
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
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from "vue"
  import { useRoute, RouterLink } from "vue-router"
  import { authApi } from "@/api/api"

  const route = useRoute()

  type Status = "loading" | "success" | "error"
  type ResendStatus = "idle" | "sending" | "sent" | "failed"

  const status = ref<Status>("loading")
  const resendStatus = ref<ResendStatus>("idle")

  // เก็บ email ไว้สำหรับ resend (ดึงได้จาก query param ถ้า backend ส่งมา)
  const email = ref("")

  onMounted(async () => {
    const token = route.query.token as string

    if (!token) {
      status.value = "error"
      return
    }

    try {
      await authApi.verifyEmail(token)
      status.value = "success"
    } catch {
      status.value = "error"
    }
  })

  async function resend() {
    if (!email.value) {
      // ถ้าไม่มี email ให้ redirect ไปหน้า login แล้วบอก user กรอกเอง
      // หรือจะเปิด prompt ก็ได้ — ตอนนี้ใช้ prompt ก่อน
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
