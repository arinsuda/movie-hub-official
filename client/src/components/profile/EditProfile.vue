<template>
  <div class="edit-root">
    <header class="edit-header">
      <div class="edit-title-group">
        <span class="edit-eyebrow">Profile</span>
        <h2 class="edit-title">Edit Profile</h2>
      </div>
      <button class="close-btn" aria-label="Close" @click="handleClose">
        <X :size="14" />
      </button>
    </header>

    <form class="edit-form" @submit.prevent="handleSave">
      <!-- Avatar Upload -->
      <div class="avatar-upload">
        <div class="avatar-preview">
          <img
            v-if="form.avatar_url"
            :src="form.avatar_url"
            alt="Avatar"
            class="preview-img"
          />
          <div v-else class="preview-fallback"><UserIcon :size="24" /></div>
        </div>
        <div class="avatar-info">
          <span class="field-label">Profile Photo</span>
          <label class="upload-btn">
            <Upload :size="12" />
            <span>Upload new</span>
            <input
              type="file"
              accept="image/*"
              class="sr-only"
              @change="handleFileChange"
            />
          </label>
        </div>
      </div>

      <div class="field-divider" />

      <!-- Display Name -->
      <div class="field-group">
        <label class="field-label" for="ep-display-name">Display Name</label>
        <input
          id="ep-display-name"
          v-model="form.display_name"
          type="text"
          class="field-input"
          placeholder="Your display name"
          maxlength="60"
        />
      </div>

      <!-- Username -->
      <div class="field-group">
        <label class="field-label" for="ep-username">Username</label>
        <div class="input-prefix-wrap">
          <span class="input-prefix">@</span>
          <input
            id="ep-username"
            v-model="form.username"
            type="text"
            class="field-input field-input--prefixed"
            placeholder="username"
            maxlength="30"
          />
        </div>
      </div>

      <!-- ─── Email ─────────────────────────────────────────── -->
      <div class="field-group">
        <label class="field-label" for="ep-email">Email</label>
        <div class="email-row">
          <input
            id="ep-email"
            v-model="form.email"
            type="email"
            class="field-input email-input"
            :class="{ 'field-input--unlocked': emailStep === 'input' }"
            :disabled="emailStep !== 'input'"
            placeholder="your@email.com"
          />

          <!-- ยังไม่ verified → Verify Email -->
          <button
            v-if="!isEmailVerified && emailStep === 'idle'"
            type="button"
            class="email-action-btn email-action-btn--verify"
            @click="handleResendVerification"
          >
            <Mail :size="12" /> Verify Email
          </button>

          <!-- idle + verified → Change Email -->
          <button
            v-else-if="emailStep === 'idle'"
            type="button"
            class="email-action-btn"
            @click="handleRequestOTP"
          >
            <Pencil :size="12" /> Change Email
          </button>

          <!-- requesting OTP -->
          <button
            v-else-if="emailStep === 'requesting'"
            type="button"
            class="email-action-btn"
            disabled
          >
            <span class="btn-spinner" /> Sending…
          </button>

          <!-- verify OTP → แสดงใน otp-box ด้านล่าง, ปุ่ม Cancel -->
          <button
            v-else-if="emailStep === 'otp'"
            type="button"
            class="email-action-btn email-action-btn--cancel"
            @click="cancelEmailChange"
          >
            <X :size="12" /> Cancel
          </button>

          <!-- input email ใหม่ -->
          <template v-else-if="emailStep === 'input'">
            <button
              type="button"
              class="email-action-btn email-action-btn--confirm"
              :disabled="
                !form.email || form.email === initialForm.email || emailSaving
              "
              @click="handleUpdateEmail"
            >
              <span v-if="emailSaving" class="btn-spinner" />
              <template v-else><Mail :size="12" /> Save Email</template>
            </button>
            <button
              type="button"
              class="email-action-btn email-action-btn--cancel"
              @click="cancelEmailChange"
            >
              <X :size="12" /> Cancel
            </button>
          </template>

          <!-- waiting for link click → spinner -->
          <button
            v-else-if="emailStep === 'waiting'"
            type="button"
            class="email-action-btn email-action-btn--waiting"
            disabled
          >
            <span class="btn-spinner" /> Waiting…
          </button>

          <!-- done ✓ -->
          <button
            v-else-if="emailStep === 'done'"
            type="button"
            class="email-action-btn email-action-btn--done"
            disabled
          >
            <CheckCircle :size="13" /> Verified!
          </button>
        </div>

        <!-- OTP box -->
        <transition name="otp-slide">
          <div v-if="emailStep === 'otp'" class="otp-box">
            <p class="otp-hint">
              Enter the 6-digit code sent to
              <strong>{{ props.user?.email }}</strong>
            </p>
            <div class="otp-row">
              <input
                v-model="otpValue"
                type="text"
                inputmode="numeric"
                maxlength="6"
                class="field-input otp-input"
                placeholder="000000"
                @input="otpValue = otpValue.replace(/\D/g, '')"
              />
              <button
                type="button"
                class="email-action-btn email-action-btn--confirm"
                :disabled="otpValue.length !== 6 || otpConfirming"
                @click="handleVerifyOTP"
              >
                <span v-if="otpConfirming" class="btn-spinner" />
                <span v-else>Confirm</span>
              </button>
            </div>
            <p v-if="otpError" class="otp-error">{{ otpError }}</p>
          </div>
        </transition>

        <!-- waiting hint -->
        <transition name="otp-slide">
          <div v-if="emailStep === 'waiting'" class="otp-box otp-box--waiting">
            <div class="waiting-row">
              <span class="waiting-spinner" />
              <p class="otp-hint" style="margin: 0">
                Verification link sent to <strong>{{ form.email }}</strong
                >.<br />
                Click the link in your email to confirm.
              </p>
            </div>
            <button
              type="button"
              class="resend-btn"
              @click="handleResendEmailVerification"
            >
              Resend link
            </button>
          </div>
        </transition>

        <!-- email hint while inputting -->
        <p v-if="emailStep === 'input'" class="email-hint">
          <Info :size="11" /> Type your new email then save. A verification link
          will be sent.
        </p>
      </div>
      <!-- ─────────────────────────────────────────────────────── -->

      <!-- Bio -->
      <div class="field-group">
        <label class="field-label" for="ep-bio">Bio</label>
        <textarea
          id="ep-bio"
          v-model="form.bio"
          class="field-textarea"
          placeholder="Say something about yourself…"
          rows="3"
          maxlength="200"
        />
        <span class="char-count">{{ form.bio?.length ?? 0 }} / 200</span>
      </div>

      <!-- Date of Birth -->
      <div class="field-group">
        <label class="field-label" for="ep-dob">Date of Birth</label>
        <input
          id="ep-dob"
          v-model="form.date_of_birth"
          type="date"
          class="field-input field-input--date"
        />
      </div>

      <!-- Gender -->
      <div class="field-group">
        <label class="field-label">Gender</label>
        <div class="radio-group">
          <label
            v-for="opt in genderOptions"
            :key="opt.value"
            class="radio-option"
            :class="{ 'radio-option--active': form.gender === opt.value }"
          >
            <input
              type="radio"
              :value="opt.value"
              v-model="form.gender"
              class="sr-only"
            />
            {{ opt.label }}
          </label>
        </div>
      </div>

      <!-- Private Account -->
      <div class="field-group">
        <div class="toggle-row">
          <div class="toggle-info">
            <span class="toggle-label">Private Account</span>
            <span class="toggle-desc"
              >Only followers can see your activity</span
            >
          </div>
          <button
            type="button"
            class="toggle-btn"
            :class="{ 'toggle-btn--on': form.is_private }"
            @click="form.is_private = !form.is_private"
            :aria-pressed="form.is_private"
          >
            <span class="toggle-knob" />
          </button>
        </div>
      </div>
    </form>

    <div class="edit-footer">
      <button
        type="button"
        class="footer-btn footer-btn--cancel"
        @click="handleClose"
      >
        Cancel
      </button>
      <button
        type="submit"
        class="footer-btn footer-btn--save"
        :disabled="saving || !isDirty"
        @click="handleSave"
      >
        <span v-if="saving" class="save-spinner" />
        <span v-else>Save Changes</span>
      </button>
    </div>

    <ConfirmModal
      v-model="showEmailConfirmModal"
      list-type="email_change"
      :item-name="props.user?.email ?? ''"
      @confirm="confirmRequestOTP"
      @cancel="showEmailConfirmModal = false"
    />
  </div>
</template>

<script setup lang="ts">
  import { reactive, ref, computed, onUnmounted } from "vue"
  import {
    X,
    Upload,
    User as UserIcon,
    Pencil,
    Mail,
    Info,
    CheckCircle,
  } from "lucide-vue-next"
  import type { UserProfile } from "@/types/user"
  import { authApi, userApi } from "@/api/api"
  import { useAuthStore } from "@/stores/auth"
  import ConfirmModal from "@/components/profile/components/ConfirmModal.vue"

  const props = defineProps<{ user: UserProfile | null }>()
  const emit = defineEmits<{ close: [] }>()

  const authStore = useAuthStore()

  // ─── email step machine ───────────────────────────────────────────────────────
  // idle → requesting → otp → input → waiting → done → idle
  type EmailStep = "idle" | "requesting" | "otp" | "input" | "waiting" | "done"
  const emailStep = ref<EmailStep>("idle")
  // ─────────────────────────────────────────────────────────────────────────────

  const saving = ref(false)
  const avatarFile = ref<File | null>(null)
  const emailSaving = ref(false)
  const showEmailConfirmModal = ref(false)
  const otpValue = ref("")
  const otpError = ref("")
  const otpConfirming = ref(false)

  // polling
  let pollTimer: ReturnType<typeof setInterval> | null = null
  const POLL_INTERVAL = 3000 // ms

  const genderOptions = [
    { value: "male", label: "Male" },
    { value: "female", label: "Female" },
    { value: "other", label: "Other" },
  ]

  const initialForm = {
    display_name: props.user?.display_name ?? "",
    username: props.user?.username ?? "",
    bio: props.user?.bio ?? "",
    avatar_url: props.user?.avatar_url ?? "",
    date_of_birth: props.user?.date_of_birth
      ? props.user.date_of_birth.slice(0, 10)
      : "",
    gender: props.user?.gender ?? "",
    is_private: props.user?.is_private ?? false,
    email: props.user?.email ?? "",
  }
  const form = reactive({ ...initialForm })

  const isEmailVerified = computed(() => !!props.user?.verified_email_at)

  const isDirty = computed(
    () =>
      avatarFile.value !== null ||
      form.display_name !== initialForm.display_name ||
      form.username !== initialForm.username ||
      form.bio !== initialForm.bio ||
      form.date_of_birth !== initialForm.date_of_birth ||
      form.gender !== initialForm.gender ||
      form.is_private !== initialForm.is_private,
  )

  // ─── helpers ─────────────────────────────────────────────────────────────────
  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  function resetEmailFlow() {
    stopPolling()
    emailStep.value = "idle"
    otpValue.value = ""
    otpError.value = ""
    form.email = initialForm.email
  }

  /** เริ่ม poll ตรวจว่า backend เปลี่ยน email + verified แล้วหรือยัง */
  function startPolling(newEmail: string) {
    stopPolling()
    pollTimer = setInterval(async () => {
      try {
        if (!props.user?.id) return
        const { data } = await userApi.getProfile(props.user.id)
        const fresh = data.user
        if (fresh.email === newEmail && fresh.verified_email_at) {
          stopPolling()
          await authStore.fetchMe()
          emailStep.value = "done"
          setTimeout(() => {
            emailStep.value = "idle"
          }, 2500)
        }
      } catch {
        /* ignore */
      }
    }, POLL_INTERVAL)
  }
  // ─────────────────────────────────────────────────────────────────────────────

  function handleFileChange(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    avatarFile.value = file
    const reader = new FileReader()
    reader.onload = ev => {
      form.avatar_url = ev.target?.result as string
    }
    reader.readAsDataURL(file)
  }

  // step 1 – request OTP
  function handleRequestOTP() {
    if (!props.user?.id) return
    showEmailConfirmModal.value = true
  }

  async function confirmRequestOTP() {
    showEmailConfirmModal.value = false
    if (!props.user?.id) return
    emailStep.value = "requesting"
    otpError.value = ""
    try {
      await userApi.requestEmailChange(props.user.id)
      emailStep.value = "otp"
    } catch (err: any) {
      otpError.value = err?.response?.data?.error ?? "Failed to send OTP."
      emailStep.value = "idle"
    }
  }

  // step 2 – verify OTP → unlock input
  async function handleVerifyOTP() {
    if (!props.user?.id || otpValue.value.length !== 6) return
    otpConfirming.value = true
    otpError.value = ""
    try {
      await userApi.verifyEmailChange(props.user.id, otpValue.value)
      emailStep.value = "input"
      otpValue.value = ""
      form.email = "" // ให้ user พิมพ์ email ใหม่
    } catch (err: any) {
      otpError.value = err?.response?.data?.error ?? "Invalid OTP."
    } finally {
      otpConfirming.value = false
    }
  }

  // step 3 – save new email → backend ส่ง verification link → poll
  async function handleUpdateEmail() {
    if (!props.user?.id || !form.email) return
    emailSaving.value = true
    otpError.value = ""
    try {
      await userApi.updateEmail(props.user.id, form.email)
      // ไม่ปิด modal, ไม่ reset — เข้า waiting
      emailStep.value = "waiting"
      startPolling(form.email)
    } catch (err: any) {
      otpError.value = err?.response?.data?.error ?? "Failed to update email."
    } finally {
      emailSaving.value = false
    }
  }

  // resend verification link ขณะ waiting
  async function handleResendEmailVerification() {
    if (!form.email) return
    try {
      await authApi.resendVerification(form.email)
    } catch {
      /* show toast ถ้ามี */
    }
  }

  function cancelEmailChange() {
    resetEmailFlow()
  }

  // ถ้ายังอยู่กลางคัน flow → ยกเลิกทั้งหมด (email กลับเป็นเดิม)
  function handleClose() {
    if (emailStep.value !== "idle" && emailStep.value !== "done") {
      resetEmailFlow()
    }
    emit("close")
  }

  // verify email (กรณียังไม่ verified ตั้งแต่แรก)
  async function handleResendVerification() {
    if (!props.user?.email) return
    try {
      await authApi.resendVerification(props.user.email)
      alert(`Verification email sent to ${props.user.email}.`)
    } catch {
      alert("Failed to resend verification email.")
    }
  }

  async function handleSave() {
    saving.value = true
    try {
      const payload = new FormData()
      if (form.display_name !== initialForm.display_name)
        payload.append("display_name", form.display_name)
      if (form.bio !== initialForm.bio) payload.append("bio", form.bio)
      if (form.date_of_birth !== initialForm.date_of_birth)
        payload.append("date_of_birth", form.date_of_birth)
      if (form.gender !== initialForm.gender)
        payload.append("gender", form.gender)
      if (form.is_private !== initialForm.is_private)
        payload.append("is_private", String(form.is_private))
      if (avatarFile.value) payload.append("avatar", avatarFile.value)
      await userApi.updateProfile(props.user!.id, payload)
      await authStore.fetchMe()
      // ไม่ emit close — ให้ user กดปิดเอง
    } catch (err) {
      console.error("Save profile failed:", err)
    } finally {
      saving.value = false
    }
  }

  onUnmounted(() => stopPolling())
</script>

<style scoped>
  .edit-root {
    --c-surface: #111111;
    --c-card: #161616;
    --c-border: rgba(255, 255, 255, 0.07);
    --c-border-h: rgba(255, 255, 255, 0.14);
    --c-red: #e1251b;
    --c-green: #30d158;
    --c-yellow: #ffd60a;
    --c-text: #f0f0f0;
    --c-sub: #8a8a8e;
    --c-muted: #3a3a3c;
    --font:
      "Inter", -apple-system, BlinkMacSystemFont, "SF Pro Text",
      "Helvetica Neue", system-ui, sans-serif;
    --ease: cubic-bezier(0.16, 1, 0.3, 1);
    padding: 24px;
    display: flex;
    flex-direction: column;
    max-height: 90vh;
    overflow: hidden;
    font-family: var(--font);
    color: var(--c-text);
    background: var(--c-surface);
    border-radius: 14px;
  }

  /* Header */
  .edit-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 24px;
  }
  .edit-title-group {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .edit-eyebrow {
    font-size: 0.58rem;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--c-muted);
  }
  .edit-title {
    font-size: 1.05rem;
    font-weight: 600;
    margin: 0;
    letter-spacing: -0.01em;
    color: #fff;
  }
  .close-btn {
    width: 28px;
    height: 28px;
    border-radius: 7px;
    border: 1px solid var(--c-border);
    background: var(--c-card);
    color: var(--c-sub);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
    flex-shrink: 0;
  }
  .close-btn:hover {
    color: var(--c-text);
    border-color: var(--c-border-h);
  }

  /* Avatar */
  .avatar-upload {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
  }
  .avatar-preview {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background: var(--c-card);
    border: 1px solid var(--c-border);
    overflow: hidden;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--c-muted);
  }
  .preview-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .avatar-info {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .upload-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--c-text);
    background: var(--c-card);
    border: 1px solid var(--c-border);
    padding: 6px 12px;
    border-radius: 7px;
    cursor: pointer;
    transition: all 0.2s;
  }
  .upload-btn:hover {
    border-color: var(--c-border-h);
  }
  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
  }
  .field-divider {
    height: 1px;
    background: var(--c-border);
    margin-bottom: 20px;
  }

  /* Form */
  .edit-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
    overflow-y: auto;
    padding-bottom: 8px;
  }
  .field-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .field-label {
    font-size: 0.68rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--c-muted);
  }
  .field-input,
  .field-textarea {
    background: var(--c-card);
    border: 1px solid var(--c-border);
    border-radius: 8px;
    padding: 10px 14px;
    font-family: var(--font);
    font-size: 0.875rem;
    color: var(--c-text);
    outline: none;
    transition: border-color 0.2s;
    width: 100%;
    box-sizing: border-box;
  }
  .field-input:focus,
  .field-textarea:focus {
    border-color: rgba(255, 255, 255, 0.2);
  }
  .field-input:disabled {
    opacity: 0.45;
    cursor: not-allowed;
  }
  .field-input--unlocked {
    border-color: rgba(255, 255, 255, 0.2);
    background: #1a1a1a;
  }
  .field-textarea {
    resize: none;
    line-height: 1.55;
  }
  .input-prefix-wrap {
    position: relative;
  }
  .input-prefix {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--c-sub);
    font-size: 0.875rem;
    pointer-events: none;
  }
  .field-input--prefixed {
    padding-left: 26px;
  }
  .field-input--date {
    color-scheme: dark;
  }
  .char-count {
    font-size: 0.65rem;
    color: var(--c-muted);
    text-align: right;
    margin-top: -4px;
  }

  /* Email row */
  .email-row {
    display: flex;
    gap: 8px;
    align-items: stretch;
  }
  .email-input {
    flex: 1;
    min-width: 0;
  }

  .email-action-btn {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    font-family: var(--font);
    font-size: 0.72rem;
    font-weight: 600;
    white-space: nowrap;
    padding: 0 12px;
    border-radius: 8px;
    border: 1px solid var(--c-border);
    background: var(--c-card);
    color: var(--c-text);
    cursor: pointer;
    transition: all 0.18s;
    flex-shrink: 0;
  }
  .email-action-btn:hover:not(:disabled) {
    border-color: var(--c-border-h);
  }
  .email-action-btn:disabled {
    opacity: 0.45;
    cursor: not-allowed;
  }

  .email-action-btn--verify {
    border-color: rgba(255, 214, 10, 0.3);
    color: var(--c-yellow);
  }
  .email-action-btn--verify:hover:not(:disabled) {
    background: rgba(255, 214, 10, 0.08);
    border-color: rgba(255, 214, 10, 0.5);
  }

  .email-action-btn--cancel {
    color: var(--c-sub);
  }
  .email-action-btn--confirm {
    background: var(--c-red);
    border-color: var(--c-red);
    color: #fff;
    padding: 0 16px;
  }
  .email-action-btn--confirm:hover:not(:disabled) {
    background: #ff3b30;
  }

  /* OTP box */
  .otp-box {
    background: var(--c-card);
    border: 1px solid var(--c-border);
    border-radius: 10px;
    padding: 14px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  .otp-hint {
    font-size: 0.78rem;
    color: var(--c-sub);
    margin: 0;
    line-height: 1.5;
  }
  .otp-hint strong {
    color: var(--c-text);
    font-weight: 500;
  }
  .otp-row {
    display: flex;
    gap: 8px;
  }
  .otp-input {
    flex: 1;
    letter-spacing: 0.25em;
    font-size: 1rem;
    text-align: center;
  }
  .otp-error {
    font-size: 0.72rem;
    color: #ff453a;
    margin: 0;
  }

  /* Email hint */
  .email-hint {
    display: flex;
    align-items: center;
    gap: 5px;
    font-size: 0.7rem;
    color: var(--c-sub);
    margin: 0;
  }

  /* OTP slide animation */
  .otp-slide-enter-active,
  .otp-slide-leave-active {
    transition: all 0.22s var(--ease);
  }
  .otp-slide-enter-from,
  .otp-slide-leave-to {
    opacity: 0;
    transform: translateY(-6px);
  }

  /* Spinner */
  .btn-spinner {
    width: 11px;
    height: 11px;
    border: 1.5px solid rgba(255, 255, 255, 0.3);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
    display: inline-block;
  }

  /* Gender */
  .radio-group {
    display: flex;
    gap: 8px;
  }
  .radio-option {
    flex: 1;
    text-align: center;
    padding: 8px 0;
    font-size: 0.8rem;
    font-weight: 500;
    border-radius: 8px;
    border: 1px solid var(--c-border);
    background: var(--c-card);
    color: var(--c-sub);
    cursor: pointer;
    transition: all 0.18s;
    user-select: none;
  }
  .radio-option:hover {
    color: var(--c-text);
    border-color: var(--c-border-h);
  }
  .radio-option--active {
    border-color: var(--c-red);
    color: var(--c-text);
    background: rgba(225, 37, 27, 0.1);
  }

  /* Toggle */
  .toggle-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--c-card);
    border: 1px solid var(--c-border);
    border-radius: 10px;
    padding: 13px 14px;
  }
  .toggle-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .toggle-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--c-text);
  }
  .toggle-desc {
    font-size: 0.72rem;
    color: var(--c-sub);
  }
  .toggle-btn {
    width: 42px;
    height: 24px;
    border-radius: 99px;
    border: none;
    background: var(--c-muted);
    cursor: pointer;
    position: relative;
    transition: background 0.2s;
    flex-shrink: 0;
  }
  .toggle-btn--on {
    background: var(--c-red);
  }
  .toggle-knob {
    position: absolute;
    top: 3px;
    left: 3px;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: #fff;
    transition: transform 0.2s;
    display: block;
  }
  .toggle-btn--on .toggle-knob {
    transform: translateX(18px);
  }

  /* Footer */
  .edit-footer {
    display: flex;
    gap: 10px;
    justify-content: flex-end;
    margin-top: 8px;
    padding-top: 16px;
    border-top: 1px solid var(--c-border);
    flex-shrink: 0;
  }
  .footer-btn {
    font-family: var(--font);
    font-size: 0.82rem;
    font-weight: 500;
    padding: 9px 18px;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .footer-btn--cancel {
    background: var(--c-card);
    color: var(--c-sub);
    border: 1px solid var(--c-border);
  }
  .footer-btn--cancel:hover {
    color: var(--c-text);
    border-color: var(--c-border-h);
  }
  .footer-btn--save {
    background: var(--c-red);
    color: #fff;
    min-width: 110px;
    justify-content: center;
  }
  .footer-btn--save:hover:not(:disabled) {
    background: #ff3b30;
    box-shadow: 0 4px 14px rgba(225, 37, 27, 0.35);
  }
  .footer-btn--save:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .save-spinner {
    width: 14px;
    height: 14px;
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

  /* waiting state */
  .email-action-btn--waiting {
    color: var(--c-sub);
    gap: 6px;
  }

  .email-action-btn--done {
    border-color: rgba(48, 209, 88, 0.4);
    color: var(--c-green);
    background: rgba(48, 209, 88, 0.08);
  }

  .otp-box--waiting {
    border-color: rgba(255, 255, 255, 0.1);
  }

  .waiting-row {
    display: flex;
    align-items: flex-start;
    gap: 12px;
  }

  .waiting-spinner {
    width: 18px;
    height: 18px;
    border: 2px solid rgba(255, 255, 255, 0.15);
    border-top-color: var(--c-sub);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    flex-shrink: 0;
    margin-top: 2px;
  }

  .resend-btn {
    font-family: var(--font);
    font-size: 0.7rem;
    color: var(--c-sub);
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    text-decoration: underline;
    text-underline-offset: 2px;
    align-self: flex-start;
  }
  .resend-btn:hover {
    color: var(--c-text);
  }
</style>
