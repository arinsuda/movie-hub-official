<template>
  <div class="watch-status-container">
    <button 
      class="watch-status-btn" 
      :class="{ 'watched': summary?.has_watched }"
      @click="showPanel = true"
    >
      <Eye v-if="summary?.has_watched" class="icon" />
      <EyeOff v-else class="icon" />
      <span v-if="summary?.has_watched">Watched {{ summary.watch_count }} times</span>
      <span v-else>Mark as Watched</span>
    </button>
    <button class="log-btn" @click="showModal = true" title="Log a watch">
      <Plus class="icon" />
    </button>

    <WatchLogModal 
      v-if="showModal" 
      :media-id="mediaId" 
      :media-type="mediaType" 
      @close="showModal = false"
      @logged="onLogged"
    />
    
    <WatchHistoryPanel
      v-if="showPanel"
      :media-id="mediaId"
      :media-type="mediaType"
      @close="showPanel = false"
      @changed="fetchSummary"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Eye, EyeOff, Plus } from 'lucide-vue-next'
import { watchLogApi } from '@/api/endpoints/watchLog'
import type { WatchSummaryResponse } from '@/types'
import WatchLogModal from './WatchLogModal.vue'
import WatchHistoryPanel from './WatchHistoryPanel.vue'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  mediaId: number
  mediaType: "movie" | "tv"
}>()

const auth = useAuthStore()
const summary = ref<WatchSummaryResponse | null>(null)
const showModal = ref(false)
const showPanel = ref(false)

async function fetchSummary() {
  if (!auth.user) return
  try {
    const res = await watchLogApi.getMyWatchLogs(props.mediaType, props.mediaId)
    summary.value = res.data.summary
  } catch (err) {
    console.error('Failed to fetch watch summary:', err)
  }
}

function onLogged() {
  showModal.value = false
  fetchSummary()
}

watch(() => props.mediaId, fetchSummary)
onMounted(fetchSummary)
</script>

<style scoped>
.watch-status-container {
  display: inline-flex;
  gap: 0.5rem;
  align-items: center;
}
.watch-status-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.875rem;
  font-weight: 500;
}
.watch-status-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}
.watch-status-btn.watched {
  background: rgba(79, 70, 229, 0.1);
  border-color: rgba(79, 70, 229, 0.3);
  color: #a5b4fc;
}
.log-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
}
.log-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}
.icon {
  width: 1.25rem;
  height: 1.25rem;
}
</style>
