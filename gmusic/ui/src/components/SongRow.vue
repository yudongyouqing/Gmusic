<template>
  <div class="song-row" :class="{ active }" @click="$emit('select', song)">
    <div class="cover">
      <img v-if="hasCover" :src="coverSrc" alt="cover" @error="onErr" />
      <div v-else class="cover__placeholder" aria-label="no cover">
        <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
          <path d="M9 3v10.5a3.5 3.5 0 1 0 2 3.15V7h6V3H9z" fill="currentColor"/>
        </svg>
      </div>
      <div v-if="active" class="playing-indicator"></div>
    </div>

    <div class="meta">
      <div class="title" :title="song.title">{{ song.title || '未知标题' }}</div>
      <div class="artist" :title="song.artist">{{ song.artist || '未知歌手' }}</div>
    </div>

    <div class="album" :title="song.album">{{ song.album || '' }}</div>
    <div class="duration">{{ formatDuration(song.duration) }}</div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  song: { type: Object, required: true },
  active: { type: Boolean, default: false },
})

const hadError = ref(false)
const hasCover = computed(() => !!props.song?.cover_url && !hadError.value)
const coverSrc = computed(() => {
  // 优先使用后端的 cover 代理接口，避免本地绝对路径跨域问题
  return props.song?.id ? `http://localhost:8080/api/cover/${props.song.id}` : ''
})

function onErr() { hadError.value = true }

function formatDuration(seconds) {
  const s = Math.max(0, Math.floor(seconds || 0))
  const mins = Math.floor(s / 60)
  const secs = s % 60
  return s ? `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}` : '--:--'
}
</script>

<style scoped>
.song-row {
  display: grid;
  grid-template-columns: 56px 1.5fr 1fr 60px;
  align-items: center;
  gap: 12px;
  height: 64px;
  padding: 8px 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: background .15s ease, transform .06s ease;
}
.song-row:hover { background: rgba(255,255,255,0.5); }
.song-row.active { background: linear-gradient(180deg, rgba(102,126,234,0.12), rgba(118,75,162,0.12)); }

.cover { position: relative; width: 48px; height: 48px; border-radius: 8px; overflow: hidden; background: rgba(0,0,0,0.08); display:flex; align-items:center; justify-content:center; }
.cover img { width: 100%; height: 100%; object-fit: cover; display:block; }
.cover__placeholder { font-size: 18px; opacity: .8; }
.playing-indicator { position: absolute; right: 6px; bottom: 6px; width: 10px; height: 10px; border-radius: 50%; background: #5b7bfe; box-shadow: 0 0 0 6px rgba(91,123,254,0.25); }

.meta { min-width: 0; display:flex; flex-direction: column; gap:4px; }
.title { color:#222; font-weight:600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.artist { color:#666; font-size: 12px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.album { color:#666; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.duration { color:#666; text-align:right; font-variant-numeric: tabular-nums; }
</style>

