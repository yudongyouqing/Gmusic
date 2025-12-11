<template>
  <div class="page page-queue">
    <div class="card card--compact">
      <TopBar :title="'播放列表'" :count="list.length" />
    </div>

    <div class="card card--fill" style="margin-top:12px">
      <SongList :songs="list" :currentSong="store.currentSong" :disableDrag="true" @select="onSelect" />
      <div v-if="store.playMode==='shuffle'" class="tip">提示：当前为随机播放，显示的是实际播放顺序；上一首/下一首按此队列进行。</div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { usePlayerStore } from '../stores/player'
import SongList from '../components/SongList.vue'
import TopBar from '../components/TopBar.vue'

const store = usePlayerStore()
const list = computed(() => store.getQueueListForView())

async function onSelect(song){
  // 在已有队列内跳转，不重建随机队列
  await store.playSong(song, { keepQueue: true })
}
</script>

<style scoped>
.page-queue{ display:flex; flex-direction:column; height:100%; min-height:0; }
.tip{ margin-top:8px; color:#666; font-size:12px; }
</style>

