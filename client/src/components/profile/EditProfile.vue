<template>
  <div class="edit-root">
    <header class="edit-header">
      <h2 class="edit-title">Edit Profile</h2>
      <button class="close-btn" aria-label="Close" @click="$emit('close')">
        <X :size="16" />
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

      <div class="edit-footer">
        <button
          type="button"
          class="footer-btn footer-btn--cancel"
          @click="$emit('close')"
        >
          Cancel
        </button>
        <button
          type="submit"
          class="footer-btn footer-btn--save"
          :disabled="saving || !isDirty"
        >
          <span v-if="saving" class="save-spinner" />
          <span v-else>Save Changes</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
  import { reactive, ref, computed } from "vue"
  import { X, Upload, User as UserIcon } from "lucide-vue-next"
  import type { UserProfile } from "@/types/user"
  import { userApi } from "@/api/api"
  import { useAuthStore } from "@/stores/auth"

  const authStore = useAuthStore()
  const props = defineProps<{ user: UserProfile | null }>()
  const emit = defineEmits<{ close: [] }>()

  const saving = ref(false)
  const avatarFile = ref<File | null>(null)

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
  }

  const form = reactive({ ...initialForm })

  const isDirty = computed(() => {
    return (
      avatarFile.value !== null ||
      form.display_name !== initialForm.display_name ||
      form.username !== initialForm.username ||
      form.bio !== initialForm.bio ||
      form.date_of_birth !== initialForm.date_of_birth ||
      form.gender !== initialForm.gender ||
      form.is_private !== initialForm.is_private
    )
  })

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

  async function handleSave() {
    saving.value = true
    try {
      const payload = new FormData()
      payload.append("display_name", form.display_name)
      payload.append("bio", form.bio)
      payload.append("date_of_birth", form.date_of_birth)
      payload.append("gender", form.gender)
      payload.append("is_private", String(form.is_private))
      if (avatarFile.value) {
        payload.append("avatar", avatarFile.value)
      }
      await userApi.updateProfile(props.user!.id, payload)

      // ✅ sync store ให้ตรงกับ DB
      await authStore.fetchMe()

      emit("close")
    } catch (err) {
      console.error("Save profile failed:", err)
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
  .edit-root {
    --c-surface: #111111;
    --c-card: #161616;
    --c-border: rgba(255, 255, 255, 0.07);
    --c-red: #e1251b;
    --c-text: #f0f0f0;
    --c-sub: #8a8a8e;
    --c-muted: #3a3a3c;
    --font:
      -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue",
      system-ui, sans-serif;

    padding: 24px;
    font-family: var(--font);
    color: var(--c-text);
    background: var(--c-surface);
    border-radius: 14px;
  }

  /* Header */
  .edit-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24px;
  }
  .edit-title {
    font-size: 1rem;
    font-weight: 600;
    margin: 0;
    letter-spacing: -0.01em;
  }
  .close-btn {
    width: 30px;
    height: 30px;
    border-radius: 8px;
    border: 1px solid var(--c-border);
    background: var(--c-card);
    color: var(--c-sub);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
  }
  .close-btn:hover {
    color: var(--c-text);
    border-color: rgba(255, 255, 255, 0.14);
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
    border-color: rgba(255, 255, 255, 0.14);
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
  }
  .field-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .field-label {
    font-size: 0.7rem;
    font-weight: 600;
    letter-spacing: 0.07em;
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

  /* Date input */
  .field-input--date {
    color-scheme: dark;
  }

  .char-count {
    font-size: 0.65rem;
    color: var(--c-muted);
    text-align: right;
    margin-top: -4px;
  }

  /* Gender radio */
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
    border-color: rgba(255, 255, 255, 0.14);
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
    padding: 12px 14px;
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
</style>
