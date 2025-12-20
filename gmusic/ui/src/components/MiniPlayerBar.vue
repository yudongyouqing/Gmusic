<template>
  <div class="mini-bar" :class="{ 'mini-bar--active': store.isPlaying, 'mini-bar--disabled': !store.currentSong }"
    @click="goNowPlaying" @touchstart.passive="onTouchStart" @touchmove.passive="onTouchMove">
    <div class="mini-bar__content">
      <!-- 封面缩略图 -->
      <div class="cover">
        <img v-if="coverSrc" :src="coverSrc" alt="cover" @error="onCoverErr" />
        <div v-else class="cover__placeholder" aria-label="no cover">
          <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
            <path d="M9 3v10.5a3.5 3.5 0 1 0 2 3.15V7h6V3H9z" fill="currentColor" />
          </svg>
        </div>
      </div>

      <!-- Meta: 标题与歌手 -->
      <div class="meta" :title="metaTitle">
        <div class="meta-title">{{ store.currentSong?.title || '未选择歌曲' }}</div>
        <div class="meta-artist">{{ store.currentSong?.artist || '未知歌手' }}</div>
      </div>

      <!-- Ticker：默认显示当前歌词 -->
      <div class="ticker" :title="tickerText">
        <transition name="slide-down" mode="out-in">
          <div :key="tickerText" class="ticker__inner">{{ tickerText }}</div>
        </transition>
      </div>

      <!-- 控制区：鼠标悬停 mini-bar 时显示 -->
      <div class="controls-wrap" @click.stop>
        <button class="btn small" :class="{ active: store.playMode === 'shuffle' }" title="播放顺序" @click="toggleMode">
          <svg v-if="store.playMode === 'shuffle'" viewBox="0 0 24 24" width="18" height="18" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 6h4l4 6h4" />
            <path d="M20 6l-2-2m2 2l-2 2" />
            <path d="M4 18h4l4-6h4" />
            <path d="M20 18l-2-2m2 2l-2 2" />
          </svg>
          <svg v-else viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"
            stroke-linecap="round" stroke-linejoin="round">
            <path d="M8 6h8a4 4 0 0 1 4 4v1" />
            <path d="M20 11l-2-2m2 2l-2 2" />
            <path d="M16 18H8a4 4 0 0 1-4-4v-1" />
            <path d="M4 13l2-2m-2 2l2 2" />
          </svg>
        </button>
        <button class="btn" :disabled="!canJump" title="上一首" @click="store.prevSong()"><svg viewBox="0 0 24 24"
            width="18" height="18">
            <path d="M6 5v14" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
            <path d="M18 6l-9 6 9 6V6z" fill="currentColor" />
          </svg></button>
        <button class="btn" :title="store.isPlaying ? '暂停' : '播放'" @click="togglePlay"><svg v-if="store.isPlaying"
            viewBox="0 0 24 24" width="18" height="18">
            <path d="M7 6h4v12H7zM13 6h4v12h-4z" fill="currentColor" />
          </svg><svg v-else viewBox="0 0 24 24" width="18" height="18">
            <path d="M8 5l12 7-12 7V5z" fill="currentColor" />
          </svg></button>
        <button class="btn" :disabled="!canJump" title="下一首" @click="store.nextSong()"><svg viewBox="0 0 24 24"
            width="18" height="18">
            <path d="M18 5v14" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
            <path d="M6 6l9 6-9 6V6z" fill="currentColor" />
          </svg></button>

        <!-- 音量 -->
        <div class="volume">
          <span class="vol-icon"><svg viewBox="0 0 24 24" width="16" height="16">
              <path d="M11 5l-5 4H4v6h2l5 4V5z" fill="currentColor" />
              <path d="M15 9a4 4 0 0 1 0 6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" />
              <path d="M17.5 7a7 7 0 0 1 0 10" stroke="currentColor" stroke-width="2" fill="none"
                stroke-linecap="round" />
            </svg></span>
          <input class="vol-range" type="range" min="0" max="100" step="1" v-model.number="localVolume"
            @input="onVolumeChange" />
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

const hadCoverErr = ref(false)
const coverSrc = computed(() => {
  const s = store.currentSong
  if (!s || hadCoverErr.value) return ''
  return s.id ? `http://localhost:8080/api/cover/${s.id}` : (s.cover_url || '')
})
function onCoverErr() { hadCoverErr.value = true }
watch(() => store.currentSong?.id, () => { hadCoverErr.value = false })

const metaTitle = computed(() => {
  const s = store.currentSong
  if (!s) return '未选择歌曲'
  return `${s.title || '未知标题'} - ${s.artist || '未知歌手'}`
})

const tickerText = computed(() => {
  const s = store.currentSong
  if (!s) return ''
  const lines = store.lyrics?.lines || []
  if (lines.length === 0) return `${s.title || '未知标题'} - ${s.artist || '未知歌手'}`
  const ms = Math.floor((store.playerStatus?.position || 0) * 1000)
  let idx = -1
  for (let i = 0; i < lines.length; i++) {
    if (lines[i].time <= ms) idx = i
    else break
  }
  let j = idx
  while (j >= 0 && !(lines[j]?.text || '').trim()) j--
  if (j >= 0) return lines[j].text
  return `${s.title || '未知标题'} - ${s.artist || '未知歌手'}`
})

const canJump = computed(() => (store.songList()?.length || 0) > 0)

let timer = null
const setupPoll = () => {
  clearInterval(timer)
  if (store.isPlaying) timer = setInterval(() => store.refreshStatus().catch(() => { }), 600)
}
watch(() => store.isPlaying, setupPoll)
watch(() => store.currentSong?.id, setupPoll)

onMounted(setupPoll)
onBeforeUnmount(() => clearInterval(timer))

function togglePlay(e) { e.stopPropagation(); if (!store.currentSong) return; if (store.isPlaying) store.pauseSong(); else store.resumeSong() }
function toggleMode(e) { e.stopPropagation(); store.togglePlayMode() }

const localVolume = ref(80)
function onVolumeChange() { store.setVolumePercent(localVolume.value) }

// 直接跳转，不再使用过渡动画
function goNowPlaying() { router.push('/now-playing') }

// 触摸滑动进入播放页
const startY = ref(0)
function onTouchStart(e) { startY.value = e.touches?.[0]?.clientY || 0 }
function onTouchMove(e) {
  const y = e.touches?.[0]?.clientY || 0
  if (startY.value && startY.value - y > 40) { startY.value = 0; router.push('/now-playing') }
}
</script>

<style scoped>
.mini-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  height: 72px;
  z-index: 1000;
  background: var(--mica-surface, rgba(255, 255, 255, 0.45));
  backdrop-filter: blur(var(--mica-blur, 16px)) saturate(var(--mica-saturate, 160%));
  border-top: 1px solid var(--mica-border, rgba(255, 255, 255, 0.35));
  box-shadow: 0 -8px 24px rgba(0, 0, 0, 0.18);
}

.mini-bar__content {
  height: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px calc(10px + env(safe-area-inset-bottom, 0px));
  box-sizing: border-box;
}

.cover {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  overflow: hidden;
  flex: 0 0 auto;
  background: rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover__placeholder {
  opacity: .8;
}

.meta {
  flex: 0 1 320px;
  min-width: 160px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.meta-title {
  font-size: 14px;
  font-weight: 600;
  color: #222;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta-artist {
  font-size: 12px;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
}

.ticker {
  flex: 1 1 auto;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  position: relative;
  height: 1.5em;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ticker__inner {
  width: 100%;
  text-align: center;
  color: #222;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: var(--ticker-font-size, 14px);
  transform: translateX(var(--mini-ticker-offset, 0px));
  transition: transform .2s ease;
}

.controls-wrap {
  display: none;
  align-items: center;
  gap: 8px;
}

.mini-bar:hover .controls-wrap {
  display: flex;
  flex: 1 1 auto;
  justify-content: center;
}

.mini-bar:hover .ticker {
  display: none;
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.4s ease-out;
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn.small {
  width: 32px;
  height: 32px;
}

.btn.active {
  background: linear-gradient(135deg, #667eea22, #764ba222);
}

.btn[disabled] {
  opacity: .6;
  cursor: not-allowed;
}

.btn:hover {
  background: rgba(255, 255, 255, 0.8);
}

.volume {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 120px;
}

.vol-range {
  width: 100px;
  accent-color: #667eea;
}

@media (max-width: 900px) {
  .volume {
    display: none;
  }
}

@media (max-width: 768px) {
  .btn {
    width: 40px;
    height: 40px;
  }
}
</style>