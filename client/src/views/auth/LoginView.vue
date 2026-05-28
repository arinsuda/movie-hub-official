<template>
  <div class="auth-page">
    <div class="auth-bg"><div class="auth-bg__overlay" /></div>
    <div class="auth-container">
      <div class="auth-brand">
        <div class="brand-logo">
          <span class="brand-movie">MOVIE</span
          ><span class="brand-hub">HUB</span>
        </div>
        <p class="brand-tagline">
          Welcome to MovieHub – Your Ultimate Destination<br />
          for Movie Reviews and Recommendations!
        </p>
      </div>
      <div class="auth-divider" />
      <div class="auth-form-panel">
        <h1 class="auth-title">LOGIN</h1>
        <form @submit.prevent="handleLogin" class="auth-form" novalidate>
          <div class="form-group">
            <label class="form-label">Username</label>
            <div class="input-wrapper">
              <span class="input-icon"><IconUser /></span>
              <input
                v-model="form.identifier"
                type="text"
                class="form-input"
                autocomplete="username"
                :disabled="isLoading"
              />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Password</label>
            <div class="input-wrapper">
              <span class="input-icon"><IconLock /></span>
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
                <IconEyeOff v-if="showPassword" /><IconEye v-else />
              </button>
            </div>
          </div>
          <p v-if="errorMsg" class="form-error">{{ errorMsg }}</p>
          <div class="form-forgot">
            <RouterLink to="/forgot-password" class="link-muted"
              >forget password</RouterLink
            >
          </div>
          <button type="submit" class="btn-primary" :disabled="isLoading">
            <span v-if="isLoading" class="spinner" /><span v-else>LOGIN</span>
          </button>
          <p class="form-footer">
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
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import {
  User as IconUser,
  Lock as IconLock,
  Eye as IconEye,
  EyeOff as IconEyeOff,
} from "lucide-vue-next";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const form = ref({ identifier: "", password: "" });
const showPassword = ref(false);
const isLoading = ref(false);
const errorMsg = ref("");

async function handleLogin() {
  errorMsg.value = "";
  if (!form.value.identifier || !form.value.password) {
    errorMsg.value = "Please fill in all fields";
    return;
  }
  isLoading.value = true;
  try {
    await authStore.login(form.value);
    const redirect = route.query.redirect as string | undefined;
    router.push(redirect ?? { name: "home" });
  } catch (e: any) {
    errorMsg.value =
      e?.response?.data?.error ?? "Login failed. Please try again.";
  } finally {
    isLoading.value = false;
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
  font-family: "Impact", "Arial Black", sans-serif;
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
  font-family: "Impact", "Arial Black", sans-serif;
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
  pointer-events: none;
  width: 18px;
  height: 18px;
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
  padding: 0;
  width: 18px;
  height: 18px;
  transition: color 0.2s;
}
.input-toggle:hover {
  color: #222;
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
  font-family: "Impact", "Arial Black", sans-serif;
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
</style>
