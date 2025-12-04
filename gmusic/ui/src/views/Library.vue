<template>
  <div class="page page-library">
    <!-- 扫描面板：在这里输入目录并启动扫描 -->
    <ScanPanel />

    <div style="height:16px"></div>
    <SearchBar @search="onSearch" />
    <div style="height:16px"></div>
    <SongList :songs="store.songList()" :currentSong="store.currentSong" @select="onSelect" />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { usePlayerStore } from '../stores/player'
import SearchBar from '../components/SearchBar.vue'
import SongList from '../components/SongList.vue'
import ScanPanel from '../components/ScanPanel.vue'

const store = usePlayerStore()

onMounted(() => {
  store.fetchSongs()
})

function onSearch(keyword) {
  store.doSearch(keyword)
}

function onSelect(song) {
  store.playSong(song)
}
</script>
