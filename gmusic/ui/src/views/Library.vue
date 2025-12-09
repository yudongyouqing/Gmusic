<template>
  <div class="page page-library">
    <!-- 顶部工具栏 -->
    <div class="card card--compact">
      <TopBar :title="'歌曲'" :count="store.songList().length" />
    </div>

    <!-- 扫描面板卡片 -->
    <div class="card card--compact">
      <ScanPanel />
    </div>

    <!-- 搜索框卡片 -->
    <div class="card card--compact" style="margin-top:12px">
      <SearchBar @search="onSearch" />
    </div>

    <!-- 列表卡片（填充剩余高度，内部可滚动） -->
    <div class="card card--fill" style="margin-top:12px">
      <SongList :songs="store.songList()" :currentSong="store.currentSong" @select="onSelect" />
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { usePlayerStore } from '../stores/player'
import SearchBar from '../components/SearchBar.vue'
import SongList from '../components/SongList.vue'
import ScanPanel from '../components/ScanPanel.vue'
import TopBar from '../components/TopBar.vue'

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
