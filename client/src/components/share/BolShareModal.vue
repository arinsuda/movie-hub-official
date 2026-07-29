<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="modelValue"
        class="share-modal-overlay"
        role="dialog"
        aria-modal="true"
        :aria-label="$t('share.bolShareTitle')"
        @click.self="close"
        @keydown.esc="close"
        tabindex="-1"
        ref="overlayRef"
      >
        <div class="share-modal">
          <!-- Header -->
          <div class="modal-header">
            <h2 class="modal-title">
              <i class="pi pi-trophy" />
              {{ $t('share.bolShareTitle') }}
            </h2>
            <button class="btn-close" @click="close" :aria-label="$t('common.close') || 'Close'">
              <i class="pi pi-times" />
            </button>
          </div>

          <!-- Controls Section: Style Selector & Username Toggle -->
          <div class="controls-section">
            <div class="style-selector">
              <button
                class="style-btn"
                :class="{ active: selectedStyle === 'grid' }"
                @click="selectedStyle = 'grid'"
              >
                <i class="pi pi-th-large" />
                <span>{{ $t('share.styleGrid') }}</span>
              </button>
              <button
                class="style-btn"
                :class="{ active: selectedStyle === 'color' }"
                @click="selectedStyle = 'color'"
              >
                <i class="pi pi-palette" />
                <span>{{ $t('share.styleColor') }}</span>
              </button>
              <button
                class="style-btn"
                :class="{ active: selectedStyle === 'bars' }"
                @click="selectedStyle = 'bars'"
              >
                <i class="pi pi-bars" />
                <span>{{ $t('share.styleBars') }}</span>
              </button>
            </div>

            <!-- Username Toggle Switch -->
            <div class="toggle-row">
              <span class="toggle-label">{{ $t('share.showUsername') }}</span>
              <label class="toggle-switch">
                <input type="checkbox" v-model="showUsername" />
                <span class="slider"></span>
              </label>
            </div>
          </div>

          <!-- Preview Area -->
          <div class="preview-container">
            <div class="preview-scaler" :style="scalerStyle">
              <div ref="exportNodeRef">
                <BolShareTemplate
                  :media-type="mediaType"
                  :ranks="ranks"
                  :style="selectedStyle"
                  :username="username"
                  :show-username="showUsername"
                  :poster-blob-urls="posterBlobUrls"
                  :dominant-colors="dominantColors"
                />
              </div>
            </div>

            <!-- Loading overlay -->
            <div v-if="isPreparing || isExporting" class="export-overlay">
              <i class="pi pi-spin pi-spinner" />
              <span>{{ isPreparing ? $t('share.preparingImage') : $t('share.generating') }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="modal-actions">
            <button
              class="btn-action btn-download"
              :disabled="isPreparing || isExporting"
              @click="handleDownload"
            >
              <i class="pi pi-download" />
              {{ $t('share.downloadImage') }}
            </button>
            <button
              class="btn-action btn-secondary"
              @click="handleCopyCaption"
            >
              <i class="pi pi-copy" />
              {{ $t('share.copyCaption') }}
            </button>
            <button
              class="btn-action btn-secondary"
              @click="handleCopyLink"
            >
              <i class="pi pi-link" />
              {{ $t('share.copyLink') }}
            </button>
          </div>

          <!-- Format Label -->
          <div class="format-label">
            <i class="pi pi-mobile" />
            {{ $t('share.storyFormat') }} · 9:16 (1080×1920)
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed, onUnmounted, nextTick } from "vue"
import type { RankGroup } from "./BolShareTemplate.vue"
import type { MediaType } from "@/types/common"
import type { ExtractedColor } from "@/utils/extractDominantColor"
import { extractDominantColor } from "@/utils/extractDominantColor"
import { useShareImage } from "@/composables/useShareImage"
import BolShareTemplate from "./BolShareTemplate.vue"

const props = defineProps<{
  modelValue: boolean
  mediaType: MediaType
  ranks: RankGroup[]
  username?: string
}>()

const emit = defineEmits<{
  "update:modelValue": [value: boolean]
}>()

const share = useShareImage()

const overlayRef = ref<HTMLElement | null>(null)
const exportNodeRef = ref<HTMLElement | null>(null)

const selectedStyle = ref<"grid" | "color" | "bars">("grid")
const showUsername = ref(true)

const posterBlobUrls = ref<Record<string, string>>({})
const dominantColors = ref<Record<string, ExtractedColor>>({})
const isPreparing = ref(false)
const isExporting = share.isExporting

// Scale the 1080px export node down to 360px for preview
const PREVIEW_WIDTH = 360
const scalerStyle = computed(() => ({
  transform: `scale(${PREVIEW_WIDTH / 1080})`,
  transformOrigin: "top left",
  width: "1080px",
  height: "1920px",
}))

function close() {
  emit("update:modelValue", false)
}

/** Pre-load poster blobs & extract dominant colors */
async function loadMediaData() {
  isPreparing.value = true
  posterBlobUrls.value = {}
  dominantColors.value = {}
  share.resetFallbackTracking()

  try {
    const topRanks = props.ranks.slice(0, 5)
    const loadPromises: Promise<void>[] = []

    for (const group of topRanks) {
      for (const item of group.items) {
        const key = `${item.media_type}:${item.media.id}`
        const posterUrl = item.media.poster_url

        if (!posterUrl) continue

        const promise = (async () => {
          const blobUrl = await share.fetchImageAsBlob(posterUrl)
          if (blobUrl) {
            posterBlobUrls.value[key] = blobUrl
            const color = await extractDominantColor(blobUrl)
            dominantColors.value[key] = color
          } else {
            share.markFallbackUsed()
          }
        })()

        loadPromises.push(promise)
      }
    }

    await Promise.all(loadPromises)
  } catch (err) {
    console.error("Error loading media for BOL share:", err)
  } finally {
    isPreparing.value = false
    share.showFallbackWarningIfNeeded()
  }
}

async function handleDownload() {
  if (!exportNodeRef.value) return
  const blob = await share.generatePng(exportNodeRef.value)
  if (blob) {
    const label = props.mediaType === "movie" ? "Movies" : "TV"
    const filename = share.sanitizeFilename(
      `BOL_${label}_Top5`,
      "REMOV_BOL"
    )
    share.downloadBlob(blob, filename)
  }
}

function handleCopyCaption() {
  const mediaLabel = props.mediaType === "movie" ? "Movies" : "TV Series"
  let text = `🏆 My Top Best ${mediaLabel} of Life on REMOV ✨\n\n`

  const topRanks = props.ranks.slice(0, 5)
  topRanks.forEach((group) => {
    const titles = group.items.map((i) => i.media.title).join(", ")
    text += `#${group.rank}: ${titles}\n`
  })

  if (showUsername.value && props.username) {
    text += `\nCurated by @${props.username}\n`
  }

  text += `\nShared from remov.app ✨`
  share.copyText(text, "share.captionCopied")
}

function handleCopyLink() {
  const link = `${window.location.origin}/library`
  share.copyText(link, "share.linkCopied")
}

// Watch modal state
watch(
  () => props.modelValue,
  async (visible) => {
    if (visible) {
      document.body.style.overflow = "hidden"
      await nextTick()
      loadMediaData()
      overlayRef.value?.focus()
    } else {
      document.body.style.overflow = ""
      posterBlobUrls.value = {}
      dominantColors.value = {}
      share.cleanup()
    }
  },
  { immediate: true }
)

onUnmounted(() => {
  document.body.style.overflow = ""
  share.cleanup()
})
</script>

<style scoped>
.share-modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.share-modal {
  background: #18181c;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 20px;
  max-width: 460px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  box-shadow: 0 28px 90px rgba(0, 0, 0, 0.7);
}

/* Header */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.modal-title {
  color: #ffffff;
  font-size: 1.1rem;
  font-weight: 700;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.modal-title i {
  color: #ffd700;
}

.btn-close {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: #ffffff;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.15s;
}

.btn-close:hover {
  background: rgba(255, 255, 255, 0.14);
}

/* Controls */
.controls-section {
  padding: 16px 24px 8px 24px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.style-selector {
  display: flex;
  gap: 8px;
  background: rgba(0, 0, 0, 0.3);
  padding: 4px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.style-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 8px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.825rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.style-btn:hover {
  color: #ffffff;
}

.style-btn.active {
  background: #e50914;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(229, 9, 20, 0.4);
}

.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 4px;
}

.toggle-label {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.8);
  font-weight: 500;
}

/* Custom Toggle Switch */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  inset: 0;
  background-color: rgba(255, 255, 255, 0.15);
  transition: 0.3s;
  border-radius: 24px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #e50914;
}

input:checked + .slider:before {
  transform: translateX(20px);
}

/* Preview */
.preview-container {
  position: relative;
  margin: 12px 24px;
  border-radius: 14px;
  overflow: hidden;
  background: #08080a;
  width: 360px;
  height: 640px;
  align-self: center;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.preview-scaler {
  pointer-events: none;
}

.export-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #ffffff;
  font-size: 0.9rem;
  z-index: 5;
}

.export-overlay .pi-spinner {
  font-size: 2.2rem;
  color: #e50914;
}

/* Actions */
.modal-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px 24px 16px;
}

.btn-action {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
  border: none;
}

.btn-download {
  background: #e50914;
  color: #ffffff;
}

.btn-download:hover:not(:disabled) {
  background: #f40612;
}

.btn-download:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.05);
  color: #ffffff;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
}

/* Format label */
.format-label {
  text-align: center;
  padding: 8px 24px 20px;
  color: rgba(255, 255, 255, 0.4);
  font-size: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

/* Transition */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
