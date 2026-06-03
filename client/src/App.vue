<template>
  <RouterView v-slot="{ Component }">
    <Transition name="auth-flow" mode="out-in">
      <component :is="Component" />
    </Transition>
  </RouterView>

  <ToastContainer />
</template>

<script setup lang="ts">
  import { RouterView } from "vue-router"
  import ToastContainer from "@/components/common/ToastContainer.vue"
</script>

<style>
  /* ─── Auth Flow Smooth Page Transition ─── */

  /* ระยะเวลาและ Timing Function ในการเล่นแอนิเมชันสลับหน้า */
  .auth-flow-enter-active,
  .auth-flow-leave-active {
    transition:
      opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1),
      transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  /* ตอนที่หน้าเก่ากำลังจะหายไป (Fade Out + เลื่อนขึ้นเล็กน้อย) */
  .auth-flow-leave-to {
    opacity: 0;
    transform: translateY(-12px);
  }

  /* ตอนที่หน้าใหม่กำลังจะปรากฏกาย (Fade In เข้ามาจากด้านล่างเล็กน้อย) */
  .auth-flow-enter-from {
    opacity: 0;
    transform: translateY(12px);
  }

  /* ล้างสไตล์เริ่มต้นเมื่อ Transition ทำงานเสร็จสิ้น */
  .auth-flow-enter-to,
  .auth-flow-leave-from {
    opacity: 1;
    transform: translateY(0);
  }
</style>
