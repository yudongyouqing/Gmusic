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
      <div
        v-for="(s, i) in songs"
        :key="s.id"
        class="drag-wrap"
        :draggable="isCustom && !isSearching && !disableDrag"
        @dragstart="onDragStart(i)"
        @dragover.prevent="onDragOver(i)"
        @drop.prevent="onDrop(i)"
      >
        <SongRow
          :song="s"
          :active="currentSong && currentSong.id === s.id"
          @select="() => onSelect(s)"
        />
      </div>
      <div v-if="!songs || songs.length === 0" class="empty">暂无歌曲</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { usePlayerStore } from '../stores/player'
import SongRow from './SongRow.vue'

const props = defineProps({
  songs: { type: Array, default: () => [] },
  currentSong: Object,
  disableDrag: { type: Boolean, default: false }
})

const emit = defineEmits(['select'])
function onSelect(song) { emit('select', song) }

const store = usePlayerStore()
const isCustom = computed(() => store.sortMode === 'custom')
const isSearching = computed(() => !!store.searchResults)

const dragFrom = ref(-1)
function onDragStart(i){ dragFrom.value = i }
function onDragOver(i){ /* 可在此添加视觉高亮 */ }
function onDrop(i){
  if (dragFrom.value < 0 || dragFrom.value === i) return
  const ids = props.songs.map(s => s.id)
  const [moved] = ids.splice(dragFrom.value, 1)
  ids.splice(i, 0, moved)
  store.setSort('custom', store.sortDir)
  store.setCustomOrder(ids)
  dragFrom.value = -1
}
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

.drag-wrap[draggable="true"]{ cursor: move; }
.drag-wrap[draggable="true"]:active{ opacity: .85; }

.empty { text-align: center; color: #999; padding: 40px 0; }
</style>
