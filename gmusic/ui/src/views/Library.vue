<template>
  <div class="page page-library">
    <!-- 顶部工具栏 -->
    <div class="card card--compact">
      <TopBar :title="'歌曲'" :count="store.songList().length" />
    </div>

    <!-- 搜索和扫描一行布局 -->
    <div class="card card--compact" style="margin-top:12px">
      <div class="search-scan-row">
        <!-- 左边：搜索框 -->
        <div class="search-wrapper">
          <SearchBar @search="onSearch" />
        </div>
        <!-- 右边：扫描面板 -->
        <div class="scan-wrapper">
      <ScanPanel />
    </div>
      </div>
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

<style scoped>
.search-scan-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  width: 100%;
  box-sizing: border-box;
}

.search-wrapper {
  flex: 0 1 400px;
  min-width: 200px;
}

.scan-wrapper {
  flex: 1 1 auto;
  min-width: 0;
}

@media (max-width: 1024px) {
  .search-scan-row {
    flex-direction: column;
  }

  .search-wrapper {
    flex: 1 1 auto;
    width: 100%;
  }

  .scan-wrapper {
    flex: 1 1 auto;
    width: 100%;
  }
}
</style>
