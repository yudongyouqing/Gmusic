<template>
  <div class="song-list">
    <!-- 表头（可选） -->
    <div class="list-head">
      <div class="col cover">封面</div>
      <div class="col title">标题 / 歌手</div>
      <div class="col album">专辑</div>
      <div class="col duration">时长</div>
    </div>

    <div class="songs">
      <SongRow
        v-for="s in songs"
        :key="s.id"
        :song="s"
        :active="currentSong && currentSong.id === s.id"
        @select="() => onSelect(s)"
      />
      <div v-if="!songs || songs.length === 0" class="empty">暂无歌曲</div>
    </div>
  </div>
</template>

<script setup>
import SongRow from './SongRow.vue'

const props = defineProps({
  songs: { type: Array, default: () => [] },
  currentSong: Object
})

const emit = defineEmits(['select'])
function onSelect(song) { emit('select', song) }
</script>

<style scoped>
.song-list {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0; /* 关键：允许在父容器中拉伸 */
}

.list-head {
  display: grid;
  grid-template-columns: 56px 1.5fr 1fr 60px;
  align-items: center;
  gap: 12px;
  height: 36px;
  padding: 0 12px;
  color: #666;
  font-size: 12px;
}
.list-head .col.cover { text-align: left; }
.list-head .col.duration { text-align: right; }

.songs {
  flex: 1 1 0;          /* 关键：占满剩余高度 */
  min-height: 0;        /* 关键：允许内部滚动 */
  overflow-y: auto;
  padding-bottom: 6px;  /* 与底部条留出极小安全距离 */
}

.empty { text-align: center; color: #999; padding: 40px 0; }
</style>
