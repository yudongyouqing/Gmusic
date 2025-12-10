<template>
  <div class="player">
    <div class="player-cover">
      <img v-if="currentSong?.cover_url" :src="currentSong.cover_url" alt="Cover" />
      <div v-else class="no-cover" aria-label="no cover">
        <svg viewBox="0 0 24 24" width="48" height="48" aria-hidden="true">
          <path d="M9 3v10.5a3.5 3.5 0 1 0 2 3.15V7h6V3H9z" fill="currentColor"/>
        </svg>
      </div>
    </div>

    <div class="player-info">
      <h2>{{ currentSong?.title || '未选择歌曲' }}</h2>
      <p>{{ currentSong?.artist || '未知歌手' }}</p>
      <p class="album">{{ currentSong?.album || '未知专辑' }}</p>
    </div>

    <div class="progress-bar">
      <span class="time">{{ formatTime(playerStatus.position) }}</span>
      <div class="progress">
        <div
          class="progress-fill"
          :style="{ width: `${progressPercent}%` }"
        ></div>
      </div>
      <span class="time">{{ formatTime(playerStatus.duration) }}</span>
    </div>

    <div class="player-controls">
      <button class="control-btn" @click="$emit('stop')" title="停止">
        <svg viewBox="0 0 24 24" width="22" height="22" aria-hidden="true">
          <rect x="6" y="6" width="12" height="12" rx="2" fill="currentColor" />
        </svg>
      </button>

      <button
        v-if="!isPlaying"
        class="control-btn play-btn"
        @click="$emit('resume')"
        title="播放"
      >
        <svg viewBox="0 0 24 24" width="26" height="26" aria-hidden="true">
          <path d="M8 5l12 7-12 7V5z" fill="currentColor" />
        </svg>
      </button>
      <button
        v-else
        class="control-btn pause-btn"
        @click="$emit('pause')"
        title="暂停"
      >
        <svg viewBox="0 0 24 24" width="26" height="26" aria-hidden="true">
          <path d="M7 6h4v12H7zM13 6h4v12h-4z" fill="currentColor" />
        </svg>
      </button>
    </div>

    <div class="volume-control">
      <span aria-hidden="true">
        <svg viewBox="0 0 24 24" width="18" height="18">
          <path d="M11 5l-5 4H4v6h2l5 4V5z" fill="currentColor"/>
          <path d="M15 9a4 4 0 0 1 0 6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/>
          <path d="M17.5 7a7 7 0 0 1 0 10" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/>
        </svg>
      </span>
      <input
        type="range"
        min="0"
        max="100"
        v-model.number="volume"
        @input="onVolumeChange"
        class="volume-slider"
      />
      <span>{{ volume }}%</span>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import './Player.css'

const props = defineProps({
  currentSong: Object,
  isPlaying: Boolean,
  playerStatus: { type: Object, default: () => ({ position: 0, duration: 0 }) }
})
const emit = defineEmits(['pause', 'resume', 'stop', 'set-volume'])

const volume = ref(100)
const progressPercent = computed(() => {
  const d = props.playerStatus?.duration || 0
  const p = props.playerStatus?.position || 0
  return d > 0 ? Math.min(100, Math.max(0, (p / d) * 100)) : 0
})

function onVolumeChange() {
  emit('set-volume', volume.value)
}

function formatTime(seconds) {
  const s = Math.max(0, Math.floor(seconds || 0))
  const mins = Math.floor(s / 60)
  const secs = s % 60
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
}
</script>

