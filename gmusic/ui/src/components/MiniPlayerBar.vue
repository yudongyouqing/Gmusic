<template>
  <div
    class="mini-bar"
    :class="{ 'mini-bar--active': store.isPlaying, 'mini-bar--disabled': !store.currentSong }"
    @click="goNowPlaying"
    @touchstart.passive="onTouchStart"
    @touchmove.passive="onTouchMove"
  >
    <div class="mini-bar__content">
      <!-- 封面缩略图（优先 /api/cover/:id，失败回退占位） -->
      <div class="cover">
        <img v-if="coverSrc" :src="coverSrc" alt="cover" @error="onCoverErr" />
        <div v-else class="cover__placeholder">🎵</div>
      </div>

      <!-- Ticker：默认显示当前歌词（或 标题 - 歌手），鼠标悬停时隐藏 -->
      <div class="ticker" :title="tickerText">
        <div class="ticker__inner">{{ tickerText }}</div>
      </div>

      <!-- 控制区：鼠标悬停 mini-bar 时显示，含上下首/播放/音量/播放顺序 -->
      <div class="controls-wrap" @click.stop>
        <button class="btn small" :class="{active: store.playMode==='shuffle'}" title="播放顺序：随机/列表循环" @click="toggleMode">
          <span v-if="store.playMode==='shuffle'">🔀</span>
          <span v-else>🔁</span>
        </button>
        <button class="btn" :disabled="!canJump" title="上一首" @click="store.prevSong">⏮️</button>
        <button class="btn" :title="store.isPlaying ? '暂停' : '播放'" @click="togglePlay">
          <span v-if="store.isPlaying">⏸️</span>
          <span v-else>▶️</span>
        </button>
        <button class="btn" :disabled="!canJump" title="下一首" @click="store.nextSong">⏭️</button>

        <!-- 音量 -->
        <div class="volume">
          <span class="vol-icon">🔊</span>
          <input class="vol-range" type="range" min="0" max="100" step="1" v-model.number="localVolume" @input="onVolumeChange" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'

const store = usePlayerStore()
const router = useRouter()

// 计算封面 URL：优先使用后端封面接口，失败回退占位
const hadCoverErr = ref(false)
const coverSrc = computed(() => {
  const s = store.currentSong
  if (!s || hadCoverErr.value) return ''
  return s.id ? `http://localhost:8080/api/cover/${s.id}` : (s.cover_url || '')
})
function onCoverErr(){ hadCoverErr.value = true }
watch(() => store.currentSong?.id, () => { hadCoverErr.value = false })

// 歌词 Ticker 文本：当前时间对应的歌词行；无歌词则使用 标题 - 歌手
const tickerText = computed(() => {
  const s = store.currentSong
  if (!s) return '未选择歌曲'
  const lines = store.lyrics?.lines || []
  if (lines.length === 0) return `${s.title || '未知标题'} - ${s.artist || '未知歌手'}`
  const ms = Math.floor((store.playerStatus?.position || 0) * 1000)
  // 二分查找当前歌词行
  let left = 0, right = lines.length - 1, idx = 0
  while (left <= right) {
    const mid = (left + right) >> 1
    if (lines[mid].time <= ms) { idx = mid; left = mid + 1 } else { right = mid - 1 }
  }
  const text = lines[idx]?.text || ''
  return text || `${s.title || '未知标题'} - ${s.artist || '未知歌手'}`
})

const progressPercent = computed(() => {
  const d = store.playerStatus?.duration || 0
  const p = store.playerStatus?.position || 0
  return d > 0 ? Math.min(100, Math.max(0, (p / d) * 100)) : 0
})

const listLen = computed(() => (store.songList() || []).length)
const canJump = computed(() => listLen.value > 0)

// 轮询播放状态（播放时每600ms刷新，暂停时停止）
let timer = null
const setupPoll = () => {
  clearInterval(timer)
  if (store.isPlaying) {
    timer = setInterval(() => store.refreshStatus().catch(() => {}), 600)
  }
}
watch(() => store.isPlaying, setupPoll)
watch(() => store.currentSong?.id, setupPoll)

onMounted(setupPoll)
onBeforeUnmount(() => clearInterval(timer))

function togglePlay() {
  if (!store.currentSong) return
  if (store.isPlaying) store.pauseSong()
  else store.resumeSong()
}
function toggleMode(){ store.setPlayMode(store.playMode === 'shuffle' ? 'loop' : 'shuffle') }

// 音量
const localVolume = ref(80)
function onVolumeChange(){ store.setVolumePercent(localVolume.value) }

function goNowPlaying() { router.push('/now-playing') }

// 简单上滑手势：上滑超过 40px 进入播放页
const startY = ref(0)
function onTouchStart(e) { startY.value = e.touches?.[0]?.clientY || 0 }
function onTouchMove(e) {
  const y = e.touches?.[0]?.clientY || 0
  if (startY.value && startY.value - y > 40) { startY.value = 0; router.push('/now-playing') }
}
</script>

<style scoped>
:root{ --bottom-bar-height: 72px; }
.mini-bar { position: fixed; left: 0; right: 0; bottom: 0; height: var(--bottom-bar-height); z-index: 1000; background: var(--mica-surface, rgba(255,255,255,0.45)); backdrop-filter: blur(var(--mica-blur, 16px)) saturate(var(--mica-saturate,160%)); -webkit-backdrop-filter: blur(var(--mica-blur, 16px)) saturate(var(--mica-saturate,160%)); border-top: 1px solid var(--mica-border, rgba(255,255,255,0.35)); box-shadow: 0 -8px 24px rgba(0,0,0,0.18); }
.mini-bar--disabled { opacity: .9; }
.mini-bar__content { height: 100%; display: flex; align-items: center; gap: 12px; padding: 10px 14px calc(10px + env(safe-area-inset-bottom, 0px)); box-sizing: border-box; }
.cover { width: 44px; height: 44px; border-radius: 8px; overflow: hidden; flex: 0 0 auto; background: rgba(0,0,0,0.1); display:flex; align-items:center; justify-content:center; }
.cover img { width: 100%; height: 100%; object-fit: cover; }
.cover__placeholder { font-size: 20px; opacity: .8; }

/* 歌词 Ticker 默认显示，悬停 mini-bar 时隐藏 */
.ticker { flex: 1 1 auto; min-width: 0; overflow: hidden; white-space: nowrap; }
.ticker__inner { display: inline-block; padding-left: 8px; color:#222; animation: marquee 12s linear infinite; }
@keyframes marquee { from { transform: translateX(0); } to { transform: translateX(-50%); } }

/* 控制区默认隐藏，悬停时显示 */
.controls-wrap { display: none; align-items: center; gap: 8px; }
.mini-bar:hover .controls-wrap { display: flex; }
.mini-bar:hover .ticker { display: none; }

.btn { width: 36px; height: 36px; border-radius: 8px; border: 1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.6); cursor: pointer; display:flex; align-items:center; justify-content:center; }
.btn.small { width: 32px; height: 32px; }
.btn.active { background: linear-gradient(135deg,#667eea22,#764ba222); }
.btn[disabled] { opacity: .6; cursor: not-allowed; }
.btn:hover { background: rgba(255,255,255,0.8); }

.volume { display:flex; align-items:center; gap:6px; min-width: 120px; }
.vol-icon { font-size: 14px; }
.vol-range { width: 100px; accent-color:#667eea; }

@media (max-width: 900px) {
  .volume { display:none; } /* 窄屏隐藏音量滑块，保留按钮区域 */
}
@media (max-width: 768px) { :root{ --bottom-bar-height: 76px; } .btn { width: 40px; height: 40px; } }
</style>
