<template>
  <RouterView v-slot="{ Component }">
    <Transition name="auth-flow" mode="out-in">
      <component :is="Component" />
    </Transition>
  </RouterView>
</template>

<script setup lang="ts">
  import { RouterView } from "vue-router"
</script>

<style>
  /* ─── Auth Flow Smooth Page Transition ─── */
  .auth-flow-enter-active,
  .auth-flow-leave-active {
    transition:
      opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1),
      transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .auth-flow-leave-to {
    opacity: 0;
    transform: translateY(-12px);
  }

  .auth-flow-enter-from {
    opacity: 0;
    transform: translateY(12px);
  }

  .auth-flow-enter-to,
  .auth-flow-leave-from {
    opacity: 1;
    transform: translateY(0);
  }
</style>
