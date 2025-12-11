<template>
  <div class="page page-nowplaying">
    <!-- 背景：使用封面高斯模糊/saturate -->
    <div class="bg" :style="bgStyle"></div>
    <button class="back-btn" @click="onBack" title="返回">
      <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
        <path d="M15 6l-6 6 6 6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>
    <div class="np-layout">
      <!-- 左侧大封面与信息 -->
      <div class="panel cover-panel">
        <div class="cover-wrap" ref="coverMainRef">
          <img v-if="coverSrc" :src="coverSrc" alt="cover" @error="onCoverErr" />
          <div v-else class="cover-placeholder" aria-label="no cover">
            <svg viewBox="0 0 24 24" width="48" height="48" aria-hidden="true">
              <path d="M9 3v10.5a3.5 3.5 0 1 0 2 3.15V7h6V3H9z" fill="currentColor"/>
            </svg>
          </div>
        </div>
        <!-- 倒影 -->
        <div v-if="coverSrc" class="cover-reflect">
          <img :src="coverSrc" alt="reflect" />
        </div>
        <div class="info">
          <div class="title">{{ store.currentSong?.title || '未选择歌曲' }}</div>
          <div class="artist">{{ store.currentSong?.artist || '未知歌手' }}</div>
        </div>
      </div>

      <!-- 右侧歌词与控制 -->
      <div class="panel lyric-panel">
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
          <button class="ctrl" :disabled="!canJump" @click="store.prevSong" title="上一首">
            <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
              <path d="M6 5v14" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
              <path d="M18 6l-9 6 9 6V6z" fill="currentColor" />
            </svg>
          </button>
          <button v-if="!store.isPlaying" class="ctrl" @click="store.resumeSong" title="播放">
            <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
              <path d="M8 5l12 7-12 7V5z" fill="currentColor" />
            </svg>
          </button>
          <button v-else class="ctrl" @click="store.pauseSong" title="暂停">
            <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
              <path d="M7 6h4v12H7zM13 6h4v12h-4z" fill="currentColor" />
            </svg>
          </button>
          <button class="ctrl" :disabled="!canJump" @click="store.nextSong" title="下一首">
            <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
              <path d="M18 5v14" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
              <path d="M6 6l9 6-9 6V6z" fill="currentColor" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 缩回到底栏的过渡层 -->
    <NowPlayingOverlay
      v-if="showOverlay && coverSrc"
      :startRect="overlayStartRect"
      :targetRect="overlayTargetRect"
      :coverSrc="coverSrc"
      mode="collapse"
      @done="onOverlayDone"
    />
  </div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, ref, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import LyricDisplay from '../components/LyricDisplay.vue'
import NowPlayingOverlay from '../components/NowPlayingOverlay.vue'

const router = useRouter()
const store = usePlayerStore()

const hadCoverErr = ref(false)
const coverSrc = computed(() => {
  const s = store.currentSong
  if (!s || hadCoverErr.value) return ''
  return s.id ? `http://localhost:8080/api/cover/${s.id}` : (s.cover_url || '')
})
function onCoverErr() { hadCoverErr.value = true }

// 背景：封面高斯模糊/饱和，轻缩放
const bgStyle = computed(() => ({
  backgroundImage: coverSrc.value ? `url(${coverSrc.value})` : 'none'
}))

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

// 缩回到底栏动画
const coverMainRef = ref(null)
const showOverlay = ref(false)
const overlayStartRect = ref(null)
const overlayTargetRect = ref(null)

async function onBack(){
  try{
    await nextTick()
    const srcEl = coverMainRef.value
    const dstEl = document.querySelector('[data-mini-cover]')
    if(srcEl && dstEl && coverSrc.value){
      const r1 = srcEl.getBoundingClientRect()
      const r2 = dstEl.getBoundingClientRect()
      overlayStartRect.value = { left: r1.left, top: r1.top, width: r1.width, height: r1.height }
      overlayTargetRect.value = { left: r2.left, top: r2.top, width: r2.width, height: r2.height }
      showOverlay.value = true
      return
    }
  }catch(_){}
  router.back()
}
function onOverlayDone(){
  showOverlay.value = false
  router.back()
}
</script>

<style scoped>
.page-nowplaying { position: relative; width:100%; height:100%; min-height:0; }
/* 背景为封面：模糊/饱和/轻缩放 */
.bg{ position: absolute; inset:0; background-position:center; background-size:cover; filter: blur(26px) saturate(170%); transform: scale(1.06); opacity:.95; z-index:0; pointer-events:none; }
/* 返回按钮 */
.back-btn{ position:absolute; left:12px; top:12px; z-index:2; width:36px; height:36px; border-radius:8px; border:1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.7); display:flex; align-items:center; justify-content:center; cursor:pointer; }
.back-btn:hover{ background: rgba(255,255,255,0.9); }

.np-layout { position: relative; z-index:1; display:grid; grid-template-columns: 380px 1fr; gap: 16px; height:100%; }
.panel { background: var(--mica-surface); backdrop-filter: blur(var(--mica-blur)) saturate(var(--mica-saturate)); -webkit-backdrop-filter: blur(var(--mica-blur)) saturate(var(--mica-saturate)); border:1px solid var(--mica-border); border-radius: 12px; padding: 16px; box-sizing: border-box; }
.cover-panel { display:flex; flex-direction:column; align-items:center; gap: 16px; }
.cover-wrap { width: 100%; aspect-ratio:1; max-width: 520px; border-radius: 12px; background: rgba(0,0,0,0.08); overflow:hidden; display:flex; align-items:center; justify-content:center; }
.cover-wrap img { width:100%; height:100%; object-fit:cover; }
.cover-placeholder { font-size: 48px; opacity:.8; }
/* 封面倒影 */
.cover-reflect{ width:100%; max-width:520px; height:64px; overflow:hidden; border-radius:0 0 12px 12px; mask-image: linear-gradient(to bottom, rgba(0,0,0,.35), rgba(0,0,0,0)); -webkit-mask-image: linear-gradient(to bottom, rgba(0,0,0,.35), rgba(0,0,0,0)); }
.cover-reflect img{ width:100%; height:128px; object-fit:cover; transform: scaleY(-1); opacity:.35; filter: blur(2px); transform-origin: top; }

.info { text-align:center; }
.title { font-size: 20px; font-weight:700; color:#222; }
.artist { font-size: 13px; color:#666; margin-top: 4px; }

.lyric-panel { display:flex; flex-direction:column; height:100%; min-height:0; }
.lyric-wrap { flex:1; min-height:0; overflow:hidden; /* 关键：隐藏滚动条，由内部 transform 控制 */ }
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
