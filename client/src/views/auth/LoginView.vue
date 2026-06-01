<template>
  <div class="auth-page">
    <div class="auth-bg"><div class="auth-bg__overlay" /></div>
    <div class="auth-container" ref="containerRef">
      <div class="auth-brand" ref="brandRef">
        <div class="brand-logo">
          <span class="brand-movie">RE</span><span class="brand-hub">MOV</span>
        </div>
        <p class="brand-tagline">
          Welcome to Remov – Your Ultimate Destination<br />
          for Movie Reviews and Recommendations!
        </p>
      </div>
      <div class="auth-divider" ref="dividerRef" />
      <div class="auth-form-panel" ref="formPanelRef">
        <h1 class="auth-title">LOGIN</h1>
        <form @submit.prevent="handleLogin" class="auth-form" novalidate>
          <div class="form-group form-item">
            <label class="form-label">Username or email</label>
            <div class="input-wrapper">
              <span class="input-icon"><IconUser :size="16" /></span>
              <input
                v-model="form.identifier"
                type="text"
                class="form-input"
                autocomplete="username"
                :disabled="isLoading"
              />
            </div>
          </div>
          <div class="form-group form-item">
            <label class="form-label">Password</label>
            <div class="input-wrapper">
              <span class="input-icon"><IconLock :size="16" /></span>
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                class="form-input"
                autocomplete="current-password"
                :disabled="isLoading"
              />
              <button
                type="button"
                class="input-toggle"
                @click="showPassword = !showPassword"
                tabindex="-1"
              >
                <IconEyeOff v-if="showPassword" :size="16" /><IconEye
                  v-else
                  :size="16"
                />
              </button>
            </div>
          </div>
          <Transition name="alert">
            <div v-if="errorMsg" class="form-alert" role="alert">
              <span class="form-alert__icon">
                <IconAlertCircle :size="15" />
              </span>
              <span class="form-alert__text">{{ errorMsg }}</span>
              <button
                class="form-alert__close"
                type="button"
                @click="errorMsg = ''"
              >
                <IconX :size="13" />
              </button>
            </div>
          </Transition>
          <div class="form-forgot form-item">
            <RouterLink to="/forgot-password" class="link-muted"
              >forget password</RouterLink
            >
          </div>
          <button
            type="submit"
            class="btn-primary form-item"
            :disabled="isLoading"
          >
            <span v-if="isLoading" class="spinner" /><span v-else>LOGIN</span>
          </button>
          <p class="form-footer form-item">
            <RouterLink to="/register" class="link-underline"
              >I don't have a account.</RouterLink
            >
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from "vue"
  import { useRouter, useRoute } from "vue-router"
  import { useAuthStore } from "@/stores/auth"
  import { gsap } from "gsap"
  import {
    User as IconUser,
    Lock as IconLock,
    Eye as IconEye,
    EyeOff as IconEyeOff,
    AlertCircle as IconAlertCircle, // เพิ่ม
    X as IconX,
  } from "lucide-vue-next"

  const router = useRouter()
  const route = useRoute()
  const authStore = useAuthStore()

  const form = ref({ identifier: "", password: "" })
  const showPassword = ref(false)
  const isLoading = ref(false)
  const errorMsg = ref("")

  // --- Refs ---
  const containerRef = ref<HTMLElement | null>(null)
  const brandRef = ref<HTMLElement | null>(null)
  const dividerRef = ref<HTMLElement | null>(null)
  const formPanelRef = ref<HTMLElement | null>(null)

  onMounted(() => {
    const tl = gsap.timeline({ defaults: { ease: "power3.out" } })

    // ซ่อนก่อน (set ก่อน animate)
    gsap.set(containerRef.value, { opacity: 0, y: 32 })
    gsap.set(brandRef.value, { opacity: 0, x: -28 })
    gsap.set(dividerRef.value, { scaleY: 0, transformOrigin: "top center" })
    gsap.set(formPanelRef.value, { opacity: 0, x: 28 })
    gsap.set(".form-item", { opacity: 0, y: 14 })

    tl.to(containerRef.value, { opacity: 1, y: 0, duration: 0.55 })
      .to(brandRef.value, { opacity: 1, x: 0, duration: 0.5 }, "-=0.3")
      .to(dividerRef.value, { scaleY: 1, duration: 0.45 }, "-=0.25")
      .to(formPanelRef.value, { opacity: 1, x: 0, duration: 0.5 }, "-=0.35")
      .to(
        ".form-item",
        {
          opacity: 1,
          y: 0,
          duration: 0.4,
          stagger: 0.08,
        },
        "-=0.25",
      )
  })

  async function handleLogin() {
    errorMsg.value = ""
    if (!form.value.identifier || !form.value.password) {
      errorMsg.value = "Please fill in all fields"
      // shake animation เมื่อ error
      gsap.fromTo(
        formPanelRef.value,
        { x: -8 },
        { x: 0, duration: 0.4, ease: "elastic.out(1, 0.3)" },
      )
      return
    }
    isLoading.value = true
    try {
      await authStore.login(form.value)
      const redirect = route.query.redirect as string | undefined
      // fade out ก่อน navigate
      await gsap.to(containerRef.value, {
        opacity: 0,
        y: -20,
        duration: 0.35,
        ease: "power2.in",
      })
      router.push(redirect ?? { name: "home" })
    } catch (e: any) {
      const status = e?.response?.status

      if (status === 401) {
        errorMsg.value = "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง"
      } else if (status === 429) {
        errorMsg.value = "พยายามเข้าสู่ระบบบ่อยเกินไป กรุณารอสักครู่"
      } else if (status >= 500) {
        errorMsg.value = "เกิดข้อผิดพลาดจากเซิร์ฟเวอร์ กรุณาลองใหม่ภายหลัง"
      } else {
        errorMsg.value = "เกิดข้อผิดพลาด กรุณาลองใหม่"
      }

      gsap.fromTo(
        formPanelRef.value,
        { x: -8 },
        { x: 0, duration: 0.45, ease: "elastic.out(1, 0.3)" },
      )
    } finally {
      isLoading.value = false
    }
  }
</script>

<style scoped>
  .auth-page {
    position: relative;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    background: #0a0a0a;
  }
  .auth-bg {
    position: absolute;
    inset: 0;
    background: url("@/assets/bg.png") center/cover no-repeat;
    filter: brightness(0.35) saturate(0.7);
  }
  .auth-bg__overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(
      to right,
      rgba(10, 10, 10, 0.5) 0%,
      rgba(10, 10, 10, 0.1) 50%,
      rgba(10, 10, 10, 0.5) 100%
    );
  }
  .auth-container {
    position: relative;
    z-index: 1;
    display: flex;
    align-items: center;
    width: 100%;
    max-width: 860px;
    min-height: 460px;
    background: rgba(12, 12, 12, 0.6);
    backdrop-filter: blur(24px);
    -webkit-backdrop-filter: blur(24px);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 16px;
    overflow: hidden;
    margin: 2rem;
  }
  .auth-brand {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem 2.5rem;
    text-align: center;
  }
  .brand-logo {
    font-family: "Noto Sans Thai", "Arial Black", sans-serif;
    font-size: 3.5rem;
    font-weight: 900;
    letter-spacing: -1px;
    line-height: 1;
    margin-bottom: 1.5rem;
  }
  .brand-movie {
    color: #ffffff;
  }
  .brand-hub {
    color: #e50914;
  }
  .brand-tagline {
    font-size: 0.875rem;
    color: #a3a3a3;
    line-height: 1.7;
    max-width: 280px;
  }
  .auth-divider {
    width: 1px;
    align-self: stretch;
    background: rgba(255, 255, 255, 0.08);
    flex-shrink: 0;
  }
  .auth-form-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 3rem 2.5rem;
  }
  .auth-title {
    font-family: "Noto Sans Thai", "Arial Black", sans-serif;
    font-size: 2rem;
    font-weight: 900;
    letter-spacing: 3px;
    text-align: center;
    color: #ffffff;
    margin-bottom: 2rem;
  }
  .auth-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }
  .form-label {
    font-size: 0.82rem;
    font-weight: 500;
    color: #d1d1d1;
  }
  .input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }
  .input-icon {
    position: absolute;
    left: 12px;
    color: #777;
    display: flex;
    align-items: center;
    justify-content: center;
    pointer-events: none;
    width: 16px;
    height: 16px;
  }
  .form-input {
    width: 100%;
    padding: 0.65rem 2.8rem;
    background: rgba(228, 228, 228, 0.92);
    border: none;
    border-radius: 9999px;
    font-size: 0.9rem;
    color: #111;
    outline: none;
    transition:
      background 0.2s,
      box-shadow 0.2s;
  }
  .form-input:focus {
    background: #f0f0f0;
    box-shadow: 0 0 0 2px rgba(229, 9, 20, 0.45);
  }
  .form-input:disabled {
    opacity: 0.6;
  }
  .input-toggle {
    position: absolute;
    right: 12px;
    background: none;
    border: none;
    cursor: pointer;
    color: #666;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    width: 16px;
    height: 16px;
    transition: color 0.2s;
  }
  .input-toggle:hover {
    color: #e50914;
  }
  .form-error {
    font-size: 0.8rem;
    color: #e50914;
    text-align: center;
    margin: 0;
  }
  .form-forgot {
    text-align: right;
    margin-top: -0.25rem;
  }
  .link-muted {
    font-size: 0.8rem;
    color: #a3a3a3;
    text-decoration: none;
    transition: color 0.2s;
  }
  .link-muted:hover {
    color: #fff;
  }
  .btn-primary {
    margin-top: 0.25rem;
    width: 100%;
    padding: 0.75rem;
    background: #e50914;
    color: #fff;
    border: none;
    border-radius: 9999px;
    font-family: "Noto Sans Thai", "Arial Black", sans-serif;
    font-size: 1.1rem;
    font-weight: 900;
    letter-spacing: 2px;
    cursor: pointer;
    transition:
      background 0.2s,
      transform 0.1s;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 48px;
  }
  .btn-primary:hover:not(:disabled) {
    background: #f40612;
    transform: translateY(-1px);
  }
  .btn-primary:active:not(:disabled) {
    transform: translateY(0);
  }
  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .spinner {
    width: 20px;
    height: 20px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }
  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
  .form-footer {
    text-align: center;
    margin: 0;
  }
  .link-underline {
    font-size: 0.85rem;
    color: #d1d1d1;
    text-underline-offset: 3px;
    transition: color 0.2s;
  }
  .link-underline:hover {
    color: #fff;
  }
  @media (max-width: 640px) {
    .auth-container {
      flex-direction: column;
      margin: 1rem;
    }
    .auth-divider {
      width: auto;
      height: 1px;
      align-self: stretch;
    }
    .auth-brand {
      padding: 2rem 2rem 1.5rem;
    }
    .auth-form-panel {
      padding: 1.5rem 2rem 2rem;
    }
    .brand-logo {
      font-size: 2.5rem;
    }
  }
  /* ลบ .form-error เดิมออก แล้วใส่อันนี้แทน */
  .form-alert {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0.6rem 0.9rem;
    background: rgba(229, 9, 20, 0.12);
    border: 1px solid rgba(229, 9, 20, 0.35);
    border-radius: 10px;
    color: #ff6b6b;
    font-size: 0.8rem;
    line-height: 1.4;
  }

  .form-alert__icon {
    display: flex;
    flex-shrink: 0;
    color: #e50914;
  }

  .form-alert__text {
    flex: 1;
  }

  .form-alert__close {
    display: flex;
    flex-shrink: 0;
    background: none;
    border: none;
    cursor: pointer;
    color: #ff6b6b;
    padding: 0;
    opacity: 0.7;
    transition: opacity 0.2s;
  }

  .form-alert__close:hover {
    opacity: 1;
  }

  /* Transition animation */
  .alert-enter-active {
    animation: alertIn 0.35s cubic-bezier(0.21, 1.02, 0.73, 1) forwards;
  }

  .alert-leave-active {
    animation: alertOut 0.25s ease-in forwards;
  }

  @keyframes alertIn {
    from {
      opacity: 0;
      transform: translateY(-8px) scaleY(0.92);
      max-height: 0;
    }
    to {
      opacity: 1;
      transform: translateY(0) scaleY(1);
      max-height: 80px;
    }
  }

  @keyframes alertOut {
    from {
      opacity: 1;
      transform: translateY(0) scaleY(1);
      max-height: 80px;
    }
    to {
      opacity: 0;
      transform: translateY(-6px) scaleY(0.94);
      max-height: 0;
    }
  }
</style>
