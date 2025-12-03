<template>
  <div class="player">
    <div class="player-cover">
      <img v-if="currentSong?.cover_url" :src="currentSong.cover_url" alt="Cover" />
      <div v-else class="no-cover">🎵</div>
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
      <button class="control-btn" @click="$emit('stop')" title="停止">⏹️</button>

      <button
        v-if="!isPlaying"
        class="control-btn play-btn"
        @click="$emit('resume')"
        title="播放"
      >
        ▶️
      </button>
      <button
        v-else
        class="control-btn pause-btn"
        @click="$emit('pause')"
        title="暂停"
      >
        ⏸️
      </button>
    </div>

    <div class="volume-control">
      <span>🔊</span>
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

