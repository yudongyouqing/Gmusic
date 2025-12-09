<template>
  <div
    class="mini-bar"
    :class="{ 'mini-bar--active': store.isPlaying, 'mini-bar--disabled': !store.currentSong }"
    @click="goNowPlaying"
    @touchstart.passive="onTouchStart"
    @touchmove.passive="onTouchMove"
  >
    <div class="mini-bar__content">
      <!-- 封面缩略图 -->
      <div class="cover">
        <img v-if="store.currentSong?.cover_url" :src="store.currentSong.cover_url" alt="cover" />
        <div v-else class="cover__placeholder">🎵</div>
      </div>

      <!-- 标题/歌手 + 进度条 -->
      <div class="meta">
        <div class="title" :title="(store.currentSong?.title || '未选择歌曲')">
          {{ store.currentSong?.title || '未选择歌曲' }}
        </div>
        <div class="artist" :title="(store.currentSong?.artist || '点击进入播放界面')">
          {{ store.currentSong?.artist || '点击进入播放界面' }}
        </div>
        <div class="progress">
          <div class="progress__fill" :style="{ width: progressPercent + '%' }"></div>
        </div>
      </div>

      <!-- 控制/预留操作区 -->
      <div class="actions" @click.stop>
        <button class="btn" :disabled="!canJump" title="上一首" @click="store.prevSong">⏮️</button>
        <button class="btn" :title="store.isPlaying ? '暂停' : '播放'" @click="togglePlay">
          <span v-if="store.isPlaying">⏸️</span>
          <span v-else>▶️</span>
        </button>
        <button class="btn" :disabled="!canJump" title="下一首" @click="store.nextSong">⏭️</button>
        <!-- 预留更多操作位：喜欢/更多 -->
        <button class="btn" :disabled="!store.currentSong" title="喜欢（预留）">💖</button>
        <button class="btn" title="更多（预留）">⋯</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'

const store = usePlayerStore()
const router = useRouter()

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

function goNowPlaying() {
  router.push('/now-playing')
}

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
.meta { flex: 1 1 auto; min-width: 0; display:flex; flex-direction:column; gap: 4px; }
.title { font-weight: 600; color: #222; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.artist { font-size: 12px; color: #666; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.progress { position: relative; height: 4px; background: rgba(0,0,0,0.08); border-radius: 2px; overflow:hidden; }
.progress__fill { position:absolute; left:0; top:0; bottom:0; width:0; background: linear-gradient(90deg, #667eea 0%, #764ba2 100%); }
.actions { display:flex; gap: 8px; align-items:center; }
.btn { width: 36px; height: 36px; border-radius: 8px; border: 1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.6); cursor: pointer; display:flex; align-items:center; justify-content:center; }
.btn[disabled] { opacity: .6; cursor: not-allowed; }
.btn:hover { background: rgba(255,255,255,0.8); }
@media (max-width: 768px) { :root{ --bottom-bar-height: 76px; } .btn { width: 40px; height: 40px; } }
</style>
