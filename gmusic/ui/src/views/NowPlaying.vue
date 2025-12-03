<template>
  <div class="page page-nowplaying">
    <Player
      :currentSong="store.currentSong"
      :isPlaying="store.isPlaying"
      :playerStatus="store.playerStatus"
      @pause="store.pauseSong"
      @resume="store.resumeSong"
      @stop="store.stopSong"
      @set-volume="store.setVolumePercent"
    />

    <div style="height:16px"></div>

    <LyricDisplay v-if="store.lyrics" :lyrics="store.lyrics" :currentTime="store.playerStatus.position" />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import Player from '../components/Player.vue'
import LyricDisplay from '../components/LyricDisplay.vue'
import { usePlayerStore } from '../stores/player'

const store = usePlayerStore()

onMounted(() => {
  // 定时刷新播放状态
  setInterval(() => {
    if (store.isPlaying) store.refreshStatus()
  }, 600)
})
</script>

