<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="modelValue" class="spotlight-backdrop" @click.self="closeSpotlight">
        <div class="spotlight-modal">
          
          <!-- SELECTING ZONE (Top) -->
          <div v-if="selectedSpotlightItems.length > 0" class="spotlight-selecting-zone">
            <div
              v-for="item in selectedSpotlightItems"
              :key="item.id"
              class="selecting-chip"
            >
              <img
                v-if="item.poster_path"
                :src="`${TMDB_IMG}${item.poster_path}`"
                :alt="getMediaTitle(item)"
                class="selecting-chip-poster"
              />
              <span class="selecting-chip-title">{{ getMediaTitle(item) }}</span>
              <button
                class="btn-remove-selecting"
                :disabled="saving"
                @click="unpinSpotlightItem(item.id)"
              >✕</button>
            </div>
          </div>

          <!-- SPOTLIGHT BAR (Middle) -->
          <div class="spotlight-bar">
            <div class="spotlight-search-wrapper">
              <Search :size="18" class="spotlight-search-icon" />
              <input
                ref="spotlightInput"
                v-model="spotlightQuery"
                type="text"
                :placeholder="$t('library.bmol.searchPlaceholder')"
                class="spotlight-search-input"
                :disabled="saving"
              />
            </div>

            <!-- Select Multiple Mode Checkbox/Toggle -->
            <div class="spotlight-toggle-wrapper">
              <label class="switch-container">
                <input
                  type="checkbox"
                  v-model="selectMultipleMode"
                  :disabled="saving"
                />
                <span class="switch-slider" />
              </label>
              <span class="switch-label">{{ $t("library.bmol.selectMultiple") }}</span>
            </div>
          </div>

          <!-- SEARCH RESULTS (Bottom) -->
          <div class="spotlight-results">
            <div v-if="spotlightSearchLoading" class="spotlight-loading text-center py-4">
              <div class="loading-ring small-ring mx-auto" />
            </div>
            <div v-else-if="spotlightResults.length > 0" class="spotlight-results-list">
              <div
                v-for="res in spotlightResults"
                :key="res.id"
                class="spotlight-result-row"
                :class="{ 'spotlight-result-row--selected': isSpotlightItemSelected(res.id) }"
                @click="!saving && handleResultClick(res)"
              >
                <img
                  v-if="res.poster_path"
                  :src="`${TMDB_IMG}${res.poster_path}`"
                  :alt="getMediaTitle(res)"
                  class="result-row-poster"
                />
                <div class="result-row-poster-fallback" v-else>
                  <Film :size="16" />
                </div>
                <div class="result-row-info">
                  <p class="result-row-title">{{ getMediaTitle(res) }}</p>
                  <p class="result-row-date">{{ getMediaDate(res) || '—' }}</p>
                </div>
                <div class="result-row-indicator">
                  <span v-if="isSpotlightItemSelected(res.id)" class="indicator-check">✓</span>
                  <span v-else class="indicator-plus">+</span>
                </div>
              </div>
            </div>
            <div v-else-if="spotlightQuery.trim() !== ''" class="spotlight-empty-results">
              <p>{{ $t("library.bmol.errorEmptySearch") }}</p>
            </div>
          </div>

          <!-- MODAL ACTION FOOTER -->
          <div class="spotlight-footer">
            <button
              class="btn-spotlight-cancel"
              :disabled="saving"
              @click="closeSpotlight"
            >Close</button>
            <button
              v-if="selectMultipleMode"
              class="btn-spotlight-save"
              :disabled="selectedSpotlightItems.length === 0 || saving"
              @click="saveSpotlightItems"
            >
              {{ saving ? "Saving..." : $t("library.bmol.addSelected", { count: selectedSpotlightItems.length }) }}
            </button>
          </div>

        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, onBeforeUnmount, nextTick } from "vue"
import { Search, Film } from "lucide-vue-next"
import { movieApi } from "@/api/api"
import type { Movie, TVSeries } from "@/types"

const props = defineProps<{
  modelValue: boolean
  rank: number | null
  mediaType: "movie" | "tv"
  saving: boolean
}>()

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void
  (e: "save", selectedMedia: Array<Movie | TVSeries>): void
}>()

const TMDB_IMG = "https://image.tmdb.org/t/p/w342"

const spotlightQuery = ref("")
const spotlightResults = ref<Array<Movie | TVSeries>>([])
const spotlightSearchLoading = ref(false)
const selectedSpotlightItems = ref<Array<Movie | TVSeries>>([])
const selectMultipleMode = ref(false)
const spotlightInput = ref<HTMLInputElement | null>(null)

let spotlightSearchTimeout: ReturnType<typeof setTimeout> | null = null
let latestSearchRequest = 0
let previousBodyOverflow = ""

function invalidateSearch(): void {
  latestSearchRequest++
  if (spotlightSearchTimeout) {
    clearTimeout(spotlightSearchTimeout)
    spotlightSearchTimeout = null
  }
}

function closeSpotlight() {
  emit("update:modelValue", false)
  invalidateSearch()
  spotlightQuery.value = ""
  spotlightResults.value = []
  selectedSpotlightItems.value = []
  selectMultipleMode.value = false
}

function isSpotlightItemSelected(id: number) {
  return selectedSpotlightItems.value.some(item => item.id === id)
}

function unpinSpotlightItem(id: number) {
  selectedSpotlightItems.value = selectedSpotlightItems.value.filter(item => item.id !== id)
}

function handleResultClick(media: Movie | TVSeries) {
  if (selectMultipleMode.value) {
    if (isSpotlightItemSelected(media.id)) {
      unpinSpotlightItem(media.id)
    } else {
      selectedSpotlightItems.value.push(media)
    }
  } else {
    selectedSpotlightItems.value = [media]
    saveSpotlightItems()
  }
}

function saveSpotlightItems() {
  if (selectedSpotlightItems.value.length === 0) return
  emit("save", [...selectedSpotlightItems.value])
}

function getMediaTitle(media: Movie | TVSeries): string {
  if ("title" in media) {
    return media.title
  }
  return media.name
}

function getMediaDate(media: Movie | TVSeries): string {
  if ("release_date" in media) {
    return media.release_date
  }
  return media.first_air_date
}

// Watch query for search
watch(spotlightQuery, (newVal) => {
  invalidateSearch()
  if (!newVal.trim()) {
    spotlightResults.value = []
    return
  }
  spotlightSearchLoading.value = true
  const currentRequest = ++latestSearchRequest

  spotlightSearchTimeout = setTimeout(async () => {
    try {
      let res
      if (props.mediaType === "movie") {
        res = await movieApi.search(newVal)
      } else {
        res = await movieApi.searchSeries(newVal)
      }

      if (currentRequest === latestSearchRequest) {
        spotlightResults.value = res.data.results.slice(0, 5)
      }
    } catch (err) {
      console.error("Spotlight search failed:", err)
    } finally {
      if (currentRequest === latestSearchRequest) {
        spotlightSearchLoading.value = false
      }
    }
  }, 350)
})

// Watch modelValue for body lock & input focus
watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    previousBodyOverflow = document.body.style.overflow
    document.body.style.overflow = "hidden"
    nextTick(() => {
      spotlightInput.value?.focus()
    })
  } else {
    const activeBackdrops = document.querySelectorAll(".modal-backdrop, .spotlight-backdrop")
    if (activeBackdrops.length <= 1) {
      document.body.style.overflow = previousBodyOverflow || ""
    }
    invalidateSearch()
  }
})

// Watch mediaType to clear pending searches
watch(() => props.mediaType, () => {
  invalidateSearch()
  spotlightQuery.value = ""
  spotlightResults.value = []
})

onBeforeUnmount(() => {
  const activeBackdrops = document.querySelectorAll(".modal-backdrop, .spotlight-backdrop")
  if (activeBackdrops.length <= 1) {
    document.body.style.overflow = previousBodyOverflow || ""
  }
  invalidateSearch()
})
</script>

<style scoped>
/* Modal Fade Transition */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.22s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

/* SPOTLIGHT QUICK ADD STYLING */
.spotlight-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  z-index: 99999;
  padding-top: 10vh;
}

.spotlight-modal {
  width: 90%;
  max-width: 600px;
  background: #121214;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  box-shadow: 0 32px 64px rgba(0, 0, 0, 0.8);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 75vh;
}

/* Selecting Zone (Top) */
.spotlight-selecting-zone {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 14px 18px;
  background: rgba(255, 255, 255, 0.01);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  max-height: 120px;
  overflow-y: auto;
}

.selecting-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(225, 37, 27, 0.08);
  border: 1px solid rgba(225, 37, 27, 0.2);
  border-radius: 20px;
  padding: 4px 10px;
}

.selecting-chip-poster {
  width: 16px;
  height: 24px;
  border-radius: 2px;
  object-fit: cover;
}

.selecting-chip-title {
  font-size: 0.78rem;
  font-weight: 600;
  color: #ffffff;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.btn-remove-selecting {
  background: transparent;
  border: none;
  color: #8a8a8e;
  font-size: 0.72rem;
  cursor: pointer;
  padding: 0 2px;
  transition: color 0.2s;
}

.btn-remove-selecting:hover:not(:disabled) {
  color: #ef4444;
}

.btn-remove-selecting:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Spotlight Search Bar */
.spotlight-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.spotlight-search-wrapper {
  display: flex;
  align-items: center;
  flex: 1;
  background: #1c1c1e;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  padding: 0 12px;
  height: 42px;
}

.spotlight-search-icon {
  color: #71717a;
}

.spotlight-search-input {
  background: transparent;
  border: none;
  color: #ffffff;
  font-size: 0.9rem;
  width: 100%;
  outline: none;
  padding-left: 10px;
}

.spotlight-search-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spotlight-toggle-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.switch-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: #8a8a8e;
}

/* Switch styling */
.switch-container {
  position: relative;
  display: inline-block;
  width: 34px;
  height: 20px;
}

.switch-container input {
  opacity: 0;
  width: 0;
  height: 0;
}

.switch-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #27272a;
  transition: .25s;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.switch-slider:before {
  position: absolute;
  content: "";
  height: 12px;
  width: 12px;
  left: 3px;
  bottom: 3px;
  background-color: #8a8a8e;
  transition: .25s;
  border-radius: 50%;
}

.switch-container input:checked + .switch-slider {
  background-color: rgba(225, 37, 27, 0.15);
  border-color: rgba(225, 37, 27, 0.3);
}

.switch-container input:checked + .switch-slider:before {
  transform: translateX(14px);
  background-color: #e1251b;
}

.switch-container input:disabled + .switch-slider {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Search Results */
.spotlight-results {
  flex: 1;
  overflow-y: auto;
  min-height: 150px;
  max-height: 320px;
}

.spotlight-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px 0;
}

.spotlight-results-list {
  display: flex;
  flex-direction: column;
}

.spotlight-result-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 18px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-bottom: 1px solid rgba(255, 255, 255, 0.02);
}

.spotlight-result-row:hover {
  background: rgba(255, 255, 255, 0.03);
}

.spotlight-result-row--selected {
  background: rgba(225, 37, 27, 0.02) !important;
}

.result-row-poster {
  width: 32px;
  height: 48px;
  border-radius: 4px;
  object-fit: cover;
}

.result-row-poster-fallback {
  width: 32px;
  height: 48px;
  border-radius: 4px;
  background: #1c1c1e;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #71717a;
}

.result-row-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.result-row-title {
  font-size: 0.85rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
}

.result-row-date {
  font-size: 0.72rem;
  color: #71717a;
  margin: 4px 0 0 0;
}

.result-row-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.2s;
}

.spotlight-result-row:hover .result-row-indicator {
  border-color: rgba(225, 37, 27, 0.5);
  background: rgba(225, 37, 27, 0.1);
}

.indicator-check {
  color: #ef4444;
  font-weight: 700;
  font-size: 0.8rem;
}

.indicator-plus {
  color: #8a8a8e;
  font-size: 0.9rem;
}

.spotlight-result-row:hover .indicator-plus {
  color: #e1251b;
}

.spotlight-empty-results {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: #71717a;
  font-size: 0.85rem;
}

/* Spotlight Footer */
.spotlight-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 14px 18px;
  background: rgba(255, 255, 255, 0.01);
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.btn-spotlight-cancel {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #a1a1aa;
  padding: 8px 16px;
  font-size: 0.8rem;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-spotlight-cancel:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.04);
  color: #ffffff;
}

.btn-spotlight-cancel:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-spotlight-save {
  background: #e1251b;
  color: #ffffff;
  border: none;
  padding: 8px 16px;
  font-size: 0.8rem;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-spotlight-save:hover:not(:disabled) {
  background: #b81d15;
}

.btn-spotlight-save:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
