<template>
  <div class="profile-info">
    <div class="section-head">
      <span class="section-label">Personal Information</span>
      <div class="section-rule" />
    </div>

    <div class="info-grid">
      <div class="info-row">
        <span class="info-icon"><User :size="14" aria-hidden="true" /></span>
        <span class="info-key">Display Name</span>
        <span class="info-val">{{ user.display_name || "—" }}</span>
      </div>

      <div class="info-row">
        <span class="info-icon"><AtSign :size="14" aria-hidden="true" /></span>
        <span class="info-key">Username</span>
        <span class="info-val">@{{ user.username }}</span>
      </div>

      <div class="info-row">
        <span class="info-icon"><Shield :size="14" aria-hidden="true" /></span>
        <span class="info-key">Account Role</span>
        <span class="info-val status-premium">Standard Member</span>
      </div>

      <div class="info-row">
        <span class="info-icon"
          ><CalendarDays :size="14" aria-hidden="true"
        /></span>
        <span class="info-key">Joined Date</span>
        <span class="info-val">{{ joinDate }}</span>
      </div>
    </div>

    <template v-if="user.bio">
      <div class="section-head mt">
        <span class="section-label">Biography</span>
        <div class="section-rule" />
      </div>
      <p class="bio-block">{{ user.bio }}</p>
    </template>
  </div>
</template>

<script setup lang="ts">
  import { computed } from "vue"
  import { User, AtSign, Shield, CalendarDays } from "lucide-vue-next"
  import type { UserProfile } from "@/types"

  const props = defineProps<{ user: UserProfile }>()

  const joinDate = computed(() => {
    if (!props.user.created_at) return "—"
    return new Date(props.user.created_at).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    })
  })
</script>

<style scoped>
  .profile-info {
    /* ใช้ CSS Tokens และดีไซน์แบบเดียวกับหน้าหลัก เพื่อความกลมกลืน */
    --surface: #121212;
    --border: rgba(255, 255, 255, 0.05);
    --text: #f5f5f7;
    --sub: #8e8e93;
    --muted: #48484a;
    --ui:
      -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue",
      system-ui, sans-serif;

    color: var(--text);
    font-family: var(--ui);
  }

  .section-head {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.25rem;
  }

  .section-head.mt {
    margin-top: 2.5rem;
  }

  .section-label {
    font-size: 0.7rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--muted);
    white-space: nowrap;
  }

  .section-rule {
    flex: 1;
    height: 1px;
    background: var(--border);
  }

  .info-grid {
    display: flex;
    flex-direction: column;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 0.5rem 1.25rem;
  }

  .info-row {
    display: grid;
    grid-template-columns: 20px 140px 1fr;
    align-items: center;
    gap: 0.75rem;
    padding: 0.85rem 0;
  }

  .info-row:not(:last-child) {
    border-bottom: 1px solid var(--border);
  }

  .info-icon {
    color: var(--sub);
    display: flex;
    align-items: center;
  }

  .info-key {
    font-size: 0.8rem;
    color: var(--sub);
    font-weight: 500;
  }

  .info-val {
    font-size: 0.85rem;
    color: #ffffff;
    font-weight: 400;
  }

  .status-premium {
    color: #ffffff;
    font-weight: 500;
  }

  .bio-block {
    font-size: 0.9rem;
    color: var(--sub);
    line-height: 1.6;
    max-width: 600px;
    margin: 0;
    padding: 0.25rem 0.5rem;
  }

  /* Responsive handling */
  @media (max-width: 600px) {
    .info-row {
      grid-template-columns: 20px 1fr;
      row-gap: 0.25rem;
      padding: 1rem 0;
    }
    .info-val {
      grid-column: 2;
    }
  }
</style>
