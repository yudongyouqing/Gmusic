<template>
  <div class="page page-nowplaying">
    <div class="np-layout">
      <!-- 左侧大封面 -->
      <div class="panel cover-card">
        <div class="cover-wrap">
          <img v-if="coverSrc" :src="coverSrc" alt="cover" @error="onCoverErr" />
          <div v-else class="cover-placeholder">🎵</div>
        </div>
        <div class="title">{{ store.currentSong?.title || '未选择歌曲' }}</div>
        <div class="artist">{{ store.currentSong?.artist || '未知歌手' }}</div>
      </div>

      <!-- 右侧歌词与控制 -->
      <div class="panel lyric-card">
        <div class="lyric-wrap">
          <LyricDisplay v-if="store.lyrics" :lyrics="store.lyrics" :currentTime="store.playerStatus.position" />
          <div v-else class="no-lyrics">暂无歌词</div>
        </div>

        <!-- 进度条控制 -->
        <div class="seek">
          <div class="time">{{ formatTime(store.playerStatus.position) }}</div>
          <input
            class="seek-bar"
            type="range"
            min="0"
            :max="store.playerStatus.duration || 0"
            step="1"
            :value="localPos"
            @input="onSeekInput"
            @change="onSeekCommit"
          />
          <div class="time">{{ formatTime(store.playerStatus.duration) }}</div>
        </div>

        <!-- 控制按钮 -->
        <div class="controls">
          <button class="ctrl" :disabled="!canJump" @click="store.prevSong" title="上一首">⏮️</button>
          <button v-if="!store.isPlaying" class="ctrl" @click="store.resumeSong" title="播放">▶️</button>
          <button v-else class="ctrl" @click="store.pauseSong" title="暂停">⏸️</button>
          <button class="ctrl" :disabled="!canJump" @click="store.nextSong" title="下一首">⏭️</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { usePlayerStore } from '../stores/player'
import LyricDisplay from '../components/LyricDisplay.vue'

const store = usePlayerStore()

const hadCoverErr = ref(false)
const coverSrc = computed(() => {
  const s = store.currentSong
  if (!s || hadCoverErr.value) return ''
  return s.id ? `http://localhost:8080/api/cover/${s.id}` : (s.cover_url || '')
})
function onCoverErr() { hadCoverErr.value = true }

// 本地进度，便于拖动时不抖动
const localPos = ref(0)
watch(() => store.playerStatus.position, v => { if (!dragging.value) localPos.value = Math.floor(v || 0) })

const dragging = ref(false)
function onSeekInput(e) {
  dragging.value = true
  localPos.value = Number(e.target.value || 0)
}
async function onSeekCommit(e) {
  const val = Number(e.target.value || 0)
  await store.seekTo(val)
  dragging.value = false
}

const canJump = computed(() => (store.songList()?.length || 0) > 0)

// 播放时轮询状态（和底部条一致，双保险）
let timer = null
function setupPoll() {
  clearInterval(timer)
  if (store.isPlaying) {
    timer = setInterval(() => store.refreshStatus().catch(() => {}), 600)
  }
}
watch(() => store.isPlaying, setupPoll)
watch(() => store.currentSong?.id, () => { hadCoverErr.value = false; setupPoll() })

onMounted(() => {
  localPos.value = Math.floor(store.playerStatus.position || 0)
  setupPoll()
})

onBeforeUnmount(() => clearInterval(timer))

function formatTime(seconds) {
  const s = Math.max(0, Math.floor(seconds || 0))
  const mins = Math.floor(s / 60)
  const secs = s % 60
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
}
</script>

<style scoped>
.page-nowplaying { width:100%; height:100%; min-height:0; }
.np-layout { display:grid; grid-template-columns: 380px 1fr; gap: 16px; height:100%; }
.panel { background: var(--mica-surface); backdrop-filter: blur(var(--mica-blur)) saturate(var(--mica-saturate)); -webkit-backdrop-filter: blur(var(--mica-blur)) saturate(var(--mica-saturate)); border:1px solid var(--mica-border); border-radius: 12px; padding: 16px; box-sizing: border-box; }
.cover-card { display:flex; flex-direction:column; align-items:center; gap: 12px; }
.cover-wrap { width: 100%; aspect-ratio:1; max-width: 520px; border-radius: 12px; background: rgba(0,0,0,0.08); overflow:hidden; display:flex; align-items:center; justify-content:center; }
.cover-wrap img { width:100%; height:100%; object-fit:cover; }
.cover-placeholder { font-size: 48px; opacity:.8; }
.title { font-size: 20px; font-weight:700; color:#222; text-align:center; margin-top: 6px; }
.artist { font-size: 13px; color:#666; text-align:center; }

.lyric-card { display:flex; flex-direction:column; height:100%; min-height:0; }
.lyric-wrap { flex:1; min-height:0; overflow:auto; }
.no-lyrics { height:100%; display:flex; align-items:center; justify-content:center; color:#999; }

.seek { display:grid; grid-template-columns: 56px 1fr 56px; gap:10px; align-items:center; margin-top: 8px; }
.time { color:#666; font-variant-numeric: tabular-nums; text-align:center; }
.seek-bar { width:100%; accent-color:#667eea; }
.controls { display:flex; gap:10px; justify-content:center; padding-top: 10px; }
.ctrl { width:44px; height:44px; border-radius:10px; border:1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.7); cursor:pointer; }
.ctrl:hover { background: rgba(255,255,255,0.9); }

@media (max-width: 900px) {
  .np-layout { grid-template-columns: 1fr; }
  .cover-wrap { max-width: 100%; }
}
</style>
