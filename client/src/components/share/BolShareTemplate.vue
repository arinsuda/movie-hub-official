<template>
  <div class="bol-share-card" :class="[`style-${style}`]">
    <!-- Ambient Background Glows -->
    <div class="bg-glow bg-glow-1" :style="headerGlowStyle"></div>
    <div class="bg-glow bg-glow-2"></div>

    <!-- Header Section -->
    <div class="card-header">
      <div class="header-branding">
        <div class="brand-logo">
          <span class="logo-re">RE</span><span class="logo-mov">MOV</span>
        </div>
        <div class="header-badge">
          <span class="badge-dot"></span>
          BEST {{ mediaType === 'movie' ? 'MOVIES' : 'TV SERIES' }} OF LIFE
        </div>
      </div>
      <div v-if="showUsername && username" class="username-pill">
        <span class="by-prefix">curated by</span>
        <span class="username-text">@{{ username }}</span>
      </div>
    </div>

    <!-- Main Rank Cards List -->
    <div class="ranks-container">
      <div
        v-for="group in topRanks"
        :key="group.rank"
        class="rank-row"
        :style="getRankRowStyle(group)"
      >
        <!-- STYLE 1: POSTER GRID -->
        <template v-if="style === 'grid'">
          <div class="row-left">
            <div class="rank-badge-grid" :class="[`rank-badge-${group.rank}`]">
              <span class="rank-hash">#</span>{{ group.rank }}
            </div>
          </div>
          <div class="row-content-grid">
            <div class="posters-grid">
              <div
                v-for="item in group.items.slice(0, 4)"
                :key="item.id"
                class="poster-item"
              >
                <img
                  v-if="getPosterUrl(item)"
                  :src="getPosterUrl(item)"
                  :alt="item.media.title"
                  class="poster-img"
                />
                <div v-else class="poster-fallback">
                  <span>{{ item.media.title.charAt(0) }}</span>
                </div>
              </div>
              <div v-if="group.items.length > 4" class="more-count-badge">
                +{{ group.items.length - 4 }}
              </div>
            </div>
            <div class="titles-summary">
              {{ formatTitles(group.items) }}
            </div>
          </div>
        </template>

        <!-- STYLE 2: DOMINANT COLOR -->
        <template v-else-if="style === 'color'">
          <div class="color-row-inner">
            <div class="rank-badge-color" :style="getRankColorAccentStyle(group)">
              <span class="rank-num">#{{ group.rank }}</span>
            </div>
            <div class="color-row-body">
              <div class="color-titles">
                {{ formatTitles(group.items) }}
              </div>
              <div class="color-item-count">
                {{ group.items.length }} {{ group.items.length === 1 ? 'item' : 'items' }}
              </div>
            </div>
            <!-- Main item thumbnail preview -->
            <div class="color-thumb-wrapper" v-if="group.items[0] && getPosterUrl(group.items[0])">
              <img
                :src="getPosterUrl(group.items[0])"
                :alt="group.items[0].media.title"
                class="color-thumb-img"
              />
            </div>
          </div>
        </template>

        <!-- STYLE 3: COLOR BARS -->
        <template v-else-if="style === 'bars'">
          <div class="bars-row-inner">
            <div class="rank-badge-bars">
              #{{ group.rank }}
            </div>
            <div class="bars-container">
              <div
                v-for="item in group.items.slice(0, 5)"
                :key="item.id"
                class="bar-item"
                :style="getBarItemStyle(item)"
              >
                <span class="bar-title">{{ item.media.title }}</span>
              </div>
              <div v-if="group.items.length > 5" class="bar-more-tag">
                +{{ group.items.length - 5 }}
              </div>
            </div>
          </div>
        </template>
      </div>

      <!-- Empty state if fewer than 5 ranks -->
      <div v-if="topRanks.length === 0" class="empty-ranks">
        No {{ mediaType === 'movie' ? 'movies' : 'TV series' }} ranked yet
      </div>
    </div>

    <!-- Footer Section -->
    <div class="card-footer">
      <div class="footer-divider"></div>
      <div class="footer-content">
        <div class="footer-logo">
          <span class="logo-re">RE</span><span class="logo-mov">MOV</span>
        </div>
        <span class="footer-url">remov.app</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import type { BMOLItemResponse } from "@/types/movie"
import type { MediaType } from "@/types/common"
import type { ExtractedColor } from "@/utils/extractDominantColor"

export interface RankGroup {
  rank: number
  items: BMOLItemResponse[]
}

const props = defineProps<{
  mediaType: MediaType
  ranks: RankGroup[]
  style: "grid" | "color" | "bars"
  username?: string
  showUsername?: boolean
  posterBlobUrls: Record<string, string>
  dominantColors: Record<string, ExtractedColor>
}>()

// Take top 5 ranks only
const topRanks = computed(() => props.ranks.slice(0, 5))

const DEFAULT_COLOR: ExtractedColor = {
  r: 229,
  g: 9,
  b: 20,
  hex: "#e50914",
  rgbStr: "rgb(229, 9, 20)",
}

// Poster URL getter
function getPosterUrl(item: BMOLItemResponse | undefined): string | undefined {
  if (!item) return undefined
  const key = `${item.media_type}:${item.media.id}`
  return props.posterBlobUrls[key] || item.media.poster_url || undefined
}

// Dominant color getter for a single item
function getItemColor(item: BMOLItemResponse | undefined): ExtractedColor {
  if (!item) return DEFAULT_COLOR
  const key = `${item.media_type}:${item.media.id}`
  return props.dominantColors[key] || DEFAULT_COLOR
}

// Format titles list as "Title 1 · Title 2"
function formatTitles(items: BMOLItemResponse[]): string {
  if (!items || items.length === 0) return ""
  return items.map((i) => i.media.title).join("  ·  ")
}

// Top header glow style based on #1 rank color if available
const headerGlowStyle = computed(() => {
  const firstRank = topRanks.value[0]
  const firstItem = firstRank?.items[0]
  if (firstItem) {
    const col = getItemColor(firstItem)
    return {
      background: `radial-gradient(circle, rgba(${col.r}, ${col.g}, ${col.b}, 0.35) 0%, rgba(0, 0, 0, 0) 70%)`,
    }
  }
  return {
    background: "radial-gradient(circle, rgba(229, 9, 20, 0.3) 0%, rgba(0, 0, 0, 0) 70%)",
  }
})

// Dynamic row background for Style 2 (Dominant Color)
function getRankRowStyle(group: RankGroup) {
  const firstItem = group.items[0]
  if (props.style !== "color" || !firstItem) return {}
  const col = getItemColor(firstItem)
  return {
    background: `linear-gradient(135deg, rgba(${col.r}, ${col.g}, ${col.b}, 0.28) 0%, rgba(${Math.max(0, col.r - 40)}, ${Math.max(0, col.g - 40)}, ${Math.max(0, col.b - 40)}, 0.08) 100%)`,
    borderLeft: `4px solid rgba(${col.r}, ${col.g}, ${col.b}, 0.9)`,
    boxShadow: `0 8px 32px rgba(${col.r}, ${col.g}, ${col.b}, 0.15)`,
  }
}

// Accent badge style for Style 2
function getRankColorAccentStyle(group: RankGroup) {
  const firstItem = group.items[0]
  if (!firstItem) return {}
  const col = getItemColor(firstItem)
  return {
    background: `linear-gradient(135deg, rgba(${col.r}, ${col.g}, ${col.b}, 0.9) 0%, rgba(${col.r}, ${col.g}, ${col.b}, 0.6) 100%)`,
    boxShadow: `0 4px 14px rgba(${col.r}, ${col.g}, ${col.b}, 0.4)`,
  }
}

// Bar item style for Style 3
function getBarItemStyle(item: BMOLItemResponse) {
  const col = getItemColor(item)
  return {
    background: `linear-gradient(90deg, rgba(${col.r}, ${col.g}, ${col.b}, 0.85) 0%, rgba(${col.r}, ${col.g}, ${col.b}, 0.4) 100%)`,
    borderColor: `rgba(${col.r}, ${col.g}, ${col.b}, 0.6)`,
  }
}
</script>

<style scoped>
.bol-share-card {
  width: 1080px;
  height: 1920px;
  background-color: #0b0b0f;
  background-image:
    radial-gradient(at 0% 0%, rgba(30, 20, 50, 0.8) 0px, transparent 50%),
    radial-gradient(at 100% 100%, rgba(20, 15, 35, 0.8) 0px, transparent 50%);
  color: #ffffff;
  font-family: 'Inter', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 80px 72px 70px 72px;
  box-sizing: border-box;
  overflow: hidden;
}

/* Ambient Glows */
.bg-glow {
  position: absolute;
  pointer-events: none;
  border-radius: 50%;
}
.bg-glow-1 {
  top: -150px;
  right: -150px;
  width: 800px;
  height: 800px;
}
.bg-glow-2 {
  bottom: -200px;
  left: -200px;
  width: 900px;
  height: 900px;
  background: radial-gradient(circle, rgba(229, 9, 20, 0.15) 0%, rgba(0, 0, 0, 0) 70%);
}

/* Header */
.card-header {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 40px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-branding {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.brand-logo {
  font-size: 52px;
  font-weight: 900;
  letter-spacing: 2px;
  line-height: 1;
}

.logo-re {
  color: #e50914;
}

.logo-mov {
  color: #ffffff;
}

.header-badge {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 4px;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
}

.badge-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: #e50914;
  box-shadow: 0 0 12px #e50914;
}

.username-pill {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(16px);
  padding: 14px 28px;
  border-radius: 40px;
}

.by-prefix {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.5);
  text-transform: uppercase;
  letter-spacing: 2px;
}

.username-text {
  font-size: 26px;
  font-weight: 700;
  color: #ffffff;
}

/* Ranks Container */
.ranks-container {
  position: relative;
  z-index: 2;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  margin: 40px 0;
  gap: 24px;
}

.rank-row {
  border-radius: 24px;
  transition: all 0.2s ease;
}

/* STYLE 1: POSTER GRID */
.style-grid .rank-row {
  background: rgba(255, 255, 255, 0.035);
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(12px);
  padding: 24px 32px;
  display: flex;
  align-items: center;
  gap: 32px;
}

.rank-badge-grid {
  width: 76px;
  height: 76px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 38px;
  font-weight: 900;
  color: #ffffff;
}

.rank-badge-1 {
  background: linear-gradient(135deg, #ffd700 0%, #ff8c00 100%);
  color: #000000;
  box-shadow: 0 8px 24px rgba(255, 215, 0, 0.35);
  border: none;
}

.rank-badge-2 {
  background: linear-gradient(135deg, #e0e0e0 0%, #9e9e9e 100%);
  color: #000000;
  box-shadow: 0 6px 20px rgba(220, 220, 220, 0.25);
  border: none;
}

.rank-badge-3 {
  background: linear-gradient(135deg, #cd7f32 0%, #8b4513 100%);
  color: #ffffff;
  box-shadow: 0 6px 20px rgba(205, 127, 50, 0.25);
  border: none;
}

.rank-hash {
  font-size: 24px;
  opacity: 0.8;
  margin-right: 2px;
}

.row-content-grid {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.posters-grid {
  display: flex;
  align-items: center;
  gap: 16px;
}

.poster-item {
  width: 80px;
  height: 120px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.15);
  flex-shrink: 0;
}

.poster-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.poster-fallback {
  width: 100%;
  height: 100%;
  background: #1e1e24;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.4);
}

.more-count-badge {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 10px 18px;
  font-size: 20px;
  font-weight: 700;
  color: #ffffff;
}

.titles-summary {
  font-size: 26px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.95);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* STYLE 2: DOMINANT COLOR */
.style-color .rank-row {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(16px);
  padding: 24px 32px;
}

.color-row-inner {
  display: flex;
  align-items: center;
  gap: 28px;
}

.rank-badge-color {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 900;
  color: #ffffff;
  flex-shrink: 0;
}

.color-row-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.color-titles {
  font-size: 32px;
  font-weight: 800;
  color: #ffffff;
  line-height: 1.25;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.6);
}

.color-item-count {
  font-size: 18px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
  text-transform: uppercase;
  letter-spacing: 2px;
}

.color-thumb-wrapper {
  width: 90px;
  height: 135px;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
  border: 2px solid rgba(255, 255, 255, 0.2);
  flex-shrink: 0;
}

.color-thumb-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* STYLE 3: COLOR BARS */
.style-bars .rank-row {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 20px 28px;
}

.bars-row-inner {
  display: flex;
  align-items: center;
  gap: 28px;
}

.rank-badge-bars {
  font-size: 36px;
  font-weight: 900;
  color: #e50914;
  min-width: 60px;
  flex-shrink: 0;
}

.bars-container {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 14px;
  overflow: hidden;
}

.bar-item {
  flex: 1;
  height: 64px;
  border-radius: 16px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  box-sizing: border-box;
  border: 1px solid rgba(255, 255, 255, 0.2);
  min-width: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.bar-title {
  font-size: 22px;
  font-weight: 700;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
}

.bar-more-tag {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  padding: 12px 18px;
  font-size: 20px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.8);
  flex-shrink: 0;
}

.empty-ranks {
  text-align: center;
  padding: 80px 0;
  font-size: 28px;
  color: rgba(255, 255, 255, 0.4);
}

/* Footer */
.card-footer {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.footer-divider {
  width: 100%;
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.15) 50%, transparent 100%);
}

.footer-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.footer-logo {
  font-size: 36px;
  font-weight: 900;
  letter-spacing: 2px;
}

.footer-url {
  font-size: 22px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 2px;
}
</style>
