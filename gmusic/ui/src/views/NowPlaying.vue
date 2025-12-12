<template>
  <!-- 全屏播放页：独立盖层，阻止任何穿透点击 -->
  <div class="np-overlay" role="dialog" aria-modal="true" @click.stop>
    <!-- 背景：使用封面做模糊/饱和，且不拦截点击 -->
    <div class="bg" :style="bgStyle" aria-hidden="true"></div>

    <!-- 顶部栏：返回 + 标题/歌手（返回键始终可见） -->
    <div class="header">
      <button class="back-btn" @click="onBack" title="返回">
        <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
          <path d="M15 6l-6 6 6 6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
      <div class="titlebar">
        <div class="song-title" :title="title">{{ title }}</div>
        <div class="song-artist" :title="artist">{{ artist }}</div>
      </div>
    </div>

    <!-- 主体布局：左封面 + 右歌词/控制 -->
    <div class="content">
      <!-- 左侧封面 -->
      <div class="left">
        <div class="cover-wrap" ref="coverMainRef">
          <img v-if="coverSrc" :src="coverSrc" alt="cover" @error="onCoverErr" />
          <div v-else class="cover-placeholder" aria-label="no cover">
            <svg viewBox="0 0 24 24" width="48" height="48" aria-hidden="true">
              <path d="M9 3v10.5a3.5 3.5 0 1 0 2 3.15V7h6V3H9z" fill="currentColor"/>
            </svg>
          </div>
        </div>
        <!-- 倒影 -->
        <div v-if="coverSrc" class="cover-reflect"><img :src="coverSrc" alt="reflect" /></div>
      </div>

      <!-- 右侧歌词与控制 -->
      <div class="right">
        <div class="lyric-wrap">
          <LyricDisplay
            v-if="hasLyrics"
            :lyrics="store.lyrics"
            :currentTime="store.playerStatus.position"
            :anchorRatio="0.35"
          />
          <div v-else class="no-lyrics">暂无歌词</div>
        </div>

        <!-- 进度条控制 -->
        <div class="seek">
          <div class="time">{{ formatTime(store.playerStatus.position) }}</div>
          <input class="seek-bar" type="range" min="0" :max="store.playerStatus.duration || 0" step="1"
            :value="localPos" @input="onSeekInput" @change="onSeekCommit" />
          <div class="time">{{ formatTime(store.playerStatus.duration) }}</div>
        </div>

        <!-- 播放控制按钮 -->
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

    <!-- 缩回到底栏的过渡层（在同页内统一管理） -->
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
import { getLyrics } from '../api/music'

const router = useRouter()
const store = usePlayerStore()

// 基本信息
const title = computed(() => store.currentSong?.title || '未选择歌曲')
const artist = computed(() => store.currentSong?.artist || '未知歌手')

// 封面
const hadCoverErr = ref(false)
const coverSrc = computed(() => {
  const s = store.currentSong
  if (!s || hadCoverErr.value) return ''
  return s.id ? `http://localhost:8080/api/cover/${s.id}` : (s.cover_url || '')
})
function onCoverErr(){ hadCoverErr.value = true }

// 背景：封面模糊/饱和
const bgStyle = computed(() => ({
  backgroundImage: coverSrc.value ? `url(${coverSrc.value})` : 'none'
}))

// 歌词可见性
const hasLyrics = computed(() => {
  const l = store.lyrics
  return !!(l && Array.isArray(l.lines) && l.lines.length)
})

// 本地进度
const localPos = ref(0)
const dragging = ref(false)
watch(() => store.playerStatus.position, v => { if (!dragging.value) localPos.value = Math.floor(v || 0) })
function onSeekInput(e){ dragging.value = true; localPos.value = Number(e.target.value || 0) }
async function onSeekCommit(e){ const val = Number(e.target.value || 0); await store.seekTo(val); dragging.value = false }
const canJump = computed(() => (store.songList()?.length || 0) > 0)

// 兜底拉取歌词
async function ensureLyrics(){
  try{
    if(!store.currentSong) return
    const need = !store.lyrics || !Array.isArray(store.lyrics?.lines) || store.lyrics.lines.length === 0
    if(need){
      const { data } = await getLyrics(store.currentSong.id)
      if(data) store.lyrics = data
    }
  }catch(_){ /* ignore */ }
}

// 轮询状态
let timer = null
function setupPoll(){
  clearInterval(timer)
  if (store.isPlaying) timer = setInterval(() => store.refreshStatus().catch(() => {}), 600)
}
watch(() => store.isPlaying, setupPoll)
watch(() => store.currentSong?.id, async () => { hadCoverErr.value = false; setupPoll(); await ensureLyrics() })

onMounted(async () => {
  localPos.value = Math.floor(store.playerStatus.position || 0)
  setupPoll()
  await ensureLyrics()
  window.addEventListener('keydown', onEsc)
})
onBeforeUnmount(() => { clearInterval(timer); window.removeEventListener('keydown', onEsc) })

function onEsc(e){ if(e.key === 'Escape') onBack() }

function formatTime(seconds){
  const s = Math.max(0, Math.floor(seconds || 0))
  const mins = Math.floor(s / 60); const secs = s % 60
  return `${String(mins).padStart(2,'0')}:${String(secs).padStart(2,'0')}`
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
      const r1 = srcEl.getBoundingClientRect(); const r2 = dstEl.getBoundingClientRect()
      overlayStartRect.value = { left:r1.left, top:r1.top, width:r1.width, height:r1.height }
      overlayTargetRect.value = { left:r2.left, top:r2.top, width:r2.width, height:r2.height }
      showOverlay.value = true
      return
    }
  }catch(_){ }
  router.back()
}
function onOverlayDone(){ showOverlay.value = false; router.back() }
</script>

<style scoped>
/* 顶层全屏盖层：阻止穿透点击 */
.np-overlay{ position: fixed; inset: 0; z-index: 10000; width:100%; height:100%; min-height:0; overflow:hidden; }

/* 背景不拦截点击（pointer-events:none） */
.bg{ position:absolute; inset:0; background-position:center; background-size:cover; filter: blur(22px) saturate(160%); transform: scale(1.06); opacity:.95; z-index:0; pointer-events:none; }

/* 顶部栏：返回 + 标题 */
.header{ position: absolute; left: 16px; right: 16px; top: 12px; z-index: 2; display:flex; align-items:center; gap: 12px; }
.back-btn{ width: 40px; height: 40px; border-radius: 10px; border:1px solid rgba(0,0,0,0.06); background: rgba(255,255,255,0.96); display:flex; align-items:center; justify-content:center; cursor:pointer; box-shadow: 0 6px 18px rgba(0,0,0,0.18); }
.back-btn:hover{ background:#fff; }
.titlebar{ display:flex; flex-direction:column; gap:4px; background: rgba(255,255,255,0.6); padding: 6px 10px; border-radius: 8px; border: 1px solid rgba(0,0,0,0.06); }
.song-title{ font-weight:700; color:#222; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; max-width: 52vw; }
.song-artist{ font-size:12px; color:#555; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; max-width: 52vw; }

/* 主内容区域 */
.content{ position: relative; z-index:1; display:grid; grid-template-columns: 420px 1fr; gap: 18px; height:100%; padding: 72px 24px 16px; box-sizing: border-box; }

/* 左侧封面 */
.left{ display:flex; flex-direction:column; align-items:center; }
.cover-wrap{ width: 100%; max-width: 520px; aspect-ratio: 1; border-radius: 12px; background: rgba(0,0,0,0.08); overflow:hidden; display:flex; align-items:center; justify-content:center; box-shadow: 0 20px 50px rgba(0,0,0,0.18); }
.cover-wrap img{ width:100%; height:100%; object-fit:cover; display:block; }
.cover-placeholder{ font-size:48px; opacity:.8; }
.cover-reflect{ width:100%; max-width:520px; height:64px; overflow:hidden; border-radius: 0 0 12px 12px; mask-image: linear-gradient(to bottom, rgba(0,0,0,.35), rgba(0,0,0,0)); -webkit-mask-image: linear-gradient(to bottom, rgba(0,0,0,.35), rgba(0,0,0,0)); }
.cover-reflect img{ width:100%; height:128px; object-fit:cover; transform: scaleY(-1); opacity:.35; filter: blur(2px); transform-origin: top; }

/* 右侧歌词与控制 */
.right{ display:flex; flex-direction:column; min-width:0; }
.lyric-wrap{ flex:1; min-height:0; overflow:hidden; background: rgba(255,255,255,0.55); border:1px solid rgba(0,0,0,0.06); border-radius: 12px; backdrop-filter: blur(10px) saturate(160%); -webkit-backdrop-filter: blur(10px) saturate(160%); }
.no-lyrics{ height:100%; display:flex; align-items:center; justify-content:center; color:#666; }

.seek{ display:grid; grid-template-columns: 56px 1fr 56px; gap:10px; align-items:center; margin-top: 10px; }
.time{ color:#555; font-variant-numeric: tabular-nums; text-align:center; }
.seek-bar{ width:100%; accent-color:#667eea; }

.controls{ display:flex; gap:10px; justify-content:center; padding-top: 10px; }
.ctrl{ width:44px; height:44px; border-radius:10px; border:1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.8); cursor:pointer; }
.ctrl:hover{ background: rgba(255,255,255,0.95); }

@media (max-width: 980px){
  .content{ grid-template-columns: 1fr; padding-top: 68px; }
  .left .cover-wrap, .left .cover-reflect{ max-width: 82vw; }
}
</style>
