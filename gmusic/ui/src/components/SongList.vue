<template>
  <div class="song-list">
    <h3>播放列表 ({{ songs?.length || 0 }})</h3>
    <div class="songs">
      <div
        v-for="song in songs"
        :key="song.id"
        class="song-item"
        :class="{ active: currentSong && currentSong.id === song.id }"
        @click="$emit('select', song)"
      >
        <div class="song-info">
          <div class="song-title">{{ song.title }}</div>
          <div class="song-artist">{{ song.artist }}</div>
        </div>
        <div class="song-duration">{{ formatDuration(song.duration) }}</div>
      </div>
      <div v-if="!songs || songs.length === 0" class="empty">暂无歌曲</div>
    </div>
  </div>
</template>

<script setup>
import './SongList.css'

const props = defineProps({
  songs: { type: Array, default: () => [] },
  currentSong: Object
})

function formatDuration(seconds) {
  const s = Math.max(0, Math.floor(seconds || 0))
  const mins = Math.floor(s / 60)
  const secs = s % 60
  return `${mins}:${String(secs).padStart(2, '0')}`
}
</script>

