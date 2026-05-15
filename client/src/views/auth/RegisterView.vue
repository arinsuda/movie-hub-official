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
        <h1 class="auth-title">REGISTER</h1>
        <form @submit.prevent="handleRegister" class="auth-form" novalidate>
          <div class="form-group">
            <label class="form-label">Username</label>
            <div class="input-wrapper">
              <span class="input-icon"><IconUser /></span>
              <input
                v-model="form.username"
                type="text"
                class="form-input"
                autocomplete="username"
                :disabled="isLoading"
              />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Email address</label>
            <div class="input-wrapper">
              <span class="input-icon"><IconUser /></span>
              <input
                v-model="form.email"
                type="email"
                class="form-input"
                autocomplete="email"
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
                autocomplete="new-password"
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
          <div class="form-group">
            <label class="form-label">Confirm password</label>
            <div class="input-wrapper">
              <span class="input-icon"><IconLock /></span>
              <input
                v-model="form.confirm_password"
                :type="showConfirm ? 'text' : 'password'"
                class="form-input"
                autocomplete="new-password"
                :disabled="isLoading"
              />
              <button
                type="button"
                class="input-toggle"
                @click="showConfirm = !showConfirm"
                tabindex="-1"
              >
                <IconEyeOff v-if="showConfirm" /><IconEye v-else />
              </button>
            </div>
          </div>
          <p v-if="errorMsg" class="form-error">{{ errorMsg }}</p>
          <button type="submit" class="btn-primary" :disabled="isLoading">
            <span v-if="isLoading" class="spinner" /><span v-else
              >REGISTER</span
            >
          </button>
          <p class="form-footer">
            <RouterLink to="/login" class="link-underline"
              >I already have an account.</RouterLink
            >
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import {
  User as IconUser,
  Lock as IconLock,
  Eye as IconEye,
  EyeOff as IconEyeOff,
} from "lucide-vue-next";

const router = useRouter();
const authStore = useAuthStore();

const form = ref({
  username: "",
  email: "",
  password: "",
  confirm_password: "",
});
const showPassword = ref(false);
const showConfirm = ref(false);
const isLoading = ref(false);
const errorMsg = ref("");

async function handleRegister() {
  errorMsg.value = "";
  const { username, email, password, confirm_password } = form.value;
  if (!username || !email || !password || !confirm_password) {
    errorMsg.value = "Please fill in all fields";
    return;
  }
  if (password !== confirm_password) {
    errorMsg.value = "Passwords do not match";
    return;
  }
  if (password.length < 8) {
    errorMsg.value = "Password must be at least 8 characters";
    return;
  }
  isLoading.value = true;
  try {
    await authStore.register(form.value);
    router.push({ name: "verify-email" });
  } catch (e: any) {
    errorMsg.value =
      e?.response?.data?.error ?? "Registration failed. Please try again.";
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
  background:
    url("https://image.tmdb.org/t/p/w300/qJ2tW6WMUDux911r6m7haRef0WH.jpg") 0 0 /
      20% auto repeat-y,
    url("https://image.tmdb.org/t/p/w300/udDclJoHjfjb8Ekgsd4FDteOkCU.jpg") 20%
      0 / 20% auto repeat-y,
    url("https://image.tmdb.org/t/p/w300/rAiYTfKGqDCRIIqo664sY9XZIvQ.jpg") 40%
      0 / 20% auto repeat-y,
    url("https://image.tmdb.org/t/p/w300/b0PlSFdDwbyK0cf5RxwDpaOJQvQ.jpg") 60%
      0 / 20% auto repeat-y,
    url("https://image.tmdb.org/t/p/w300/74xTEgt7R36Fpooo50r9T25onhq.jpg") 80%
      0 / 20% auto repeat-y;
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
  padding: 2.5rem;
}
.auth-title {
  font-family: "Impact", "Arial Black", sans-serif;
  font-size: 2rem;
  font-weight: 900;
  letter-spacing: 3px;
  text-align: center;
  color: #ffffff;
  margin-bottom: 1.5rem;
}
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
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
