<template>
  <!-- 全屏播放页：独立盖层，阻止任何穿透点击 -->
  <div class="np-page" role="dialog" aria-modal="true">
    <!-- 背景：使用封面做模糊/饱和，且不拦截点击 -->
    <div class="bg" :style="bgStyle" aria-hidden="true"></div>

    <!-- 顶部栏：返回 + 标题/歌手 + 设置按钮（高层级，确保可点） -->
    <div class="header">
      <button class="back-btn" @click="router.back()" title="返回 (Esc)">
        <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
          <path d="M15 6l-6 6 6 6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
      <div class="titlebar">
        <div class="song-title" :title="title">{{ title }}</div>
        <div class="song-artist" :title="artist">{{ artist }}</div>
      </div>
      <div class="header-spacer"></div>
      <button class="gear-btn" @click="showLyricCtrl = !showLyricCtrl" title="歌词设置">
        <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
          <path d="M12 8a4 4 0 1 1 0 8 4 4 0 0 1 0-8z" fill="none" stroke="currentColor" stroke-width="2"/>
          <path d="M3 12h3M18 12h3M12 3v3M12 18v3M5.2 5.2l2.1 2.1M16.7 16.7l2.1 2.1M5.2 18.8l2.1-2.1M16.7 7.3l2.1-2.1" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>
    </div>

    <!-- 歌词设置面板 -->
    <LyricControls v-if="showLyricCtrl" :right="16" :top="56" @close="showLyricCtrl = false" />

    <!-- 主体布局：左封面 + 右歌词/控制 -->
    <div class="content">
      <!-- 左侧封面 -->
      <div class="left">
        <div class="cover-wrap">
          <img v-if="coverSrc" :src="coverSrc" alt="cover" @error="onCoverErr" />
          <div v-else class="cover-placeholder" aria-label="no cover">
            <svg viewBox="0 0 24 24" width="48" height="48" aria-hidden="true">
              <path d="M9 3v10.5a3.5 3.5 0 1 0 2 3.15V7h6V3H9z" fill="currentColor"/>
            </svg>
          </div>
        </div>
        <div v-if="coverSrc" class="cover-reflect"><img :src="coverSrc" alt="reflect" /></div>
      </div>

      <!-- 右侧歌词与控制 -->
      <div class="right">
        <div class="lyric-wrap">
          <LyricDisplay
            v-if="hasLyrics"
            :lyrics="store.lyrics"
            :currentTime="store.playerStatus.position"
            :baseFontSize="lyricUi.fontSize || 16"
            :blurOthers="lyricUi.blurOthers === true"
            :isPlaying="store.isPlaying"
          />
          <div v-else class="no-lyrics">暂无歌词</div>
        </div>

        <!-- 进度条与控制按钮 -->
        <div class="controls-wrapper">
          <div class="seek">
            <div class="time">{{ formatTime(store.playerStatus.position) }}</div>
            <input class="seek-bar" type="range" min="0" :max="store.playerStatus.duration || 0" step="1" :value="localPos" @input="onSeekInput" @change="onSeekCommit" />
            <div class="time">{{ formatTime(store.playerStatus.duration) }}</div>
          </div>
          <div class="controls">
            <button class="ctrl" :disabled="!canJump" @click="store.prevSong()" title="上一首"><svg viewBox="0 0 24 24" width="20" height="20"><path d="M6 5v14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M18 6l-9 6 9 6V6z" fill="currentColor"/></svg></button>
            <button v-if="!store.isPlaying" class="ctrl play-pause" @click="store.resumeSong()" title="播放"><svg viewBox="0 0 24 24" width="24" height="24"><path d="M8 5l12 7-12 7V5z" fill="currentColor"/></svg></button>
            <button v-else class="ctrl play-pause" @click="store.pauseSong()" title="暂停"><svg viewBox="0 0 24 24" width="24" height="24"><path d="M7 6h4v12H7zM13 6h4v12h-4z" fill="currentColor"/></svg></button>
            <button class="ctrl" :disabled="!canJump" @click="store.nextSong()" title="下一首"><svg viewBox="0 0 24 24" width="20" height="20"><path d="M18 5v14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M6 6l9 6-9 6V6z" fill="currentColor"/></svg></button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import { useLyricUiStore } from '../stores/lyric'
import { getLyrics } from '../api/music'
import LyricDisplay from '../components/LyricDisplay.vue'
import LyricControls from '../components/LyricControls.vue'

const router = useRouter()
const store = usePlayerStore()
const lyricUi = useLyricUiStore()
lyricUi.load()

const showLyricCtrl = ref(false)

// Data
const title = computed(() => store.currentSong?.title || '未选择歌曲')
const artist = computed(() => store.currentSong?.artist || '未知歌手')
const hadCoverErr = ref(false)
const coverSrc = computed(() => {
  const s = store.currentSong
  if (!s || hadCoverErr.value) return ''
  return s.id ? `http://localhost:8080/api/cover/${s.id}` : (s.cover_url || '')
})
function onCoverErr(){ hadCoverErr.value = true }

const bgStyle = computed(() => ({ backgroundImage: coverSrc.value ? `url(${coverSrc.value})` : 'none' }))
const hasLyrics = computed(() => store.lyrics && Array.isArray(store.lyrics.lines) && store.lyrics.lines.length > 0)

// Seek bar
const localPos = ref(0)
const dragging = ref(false)
watch(() => store.playerStatus.position, v => { if (!dragging.value) localPos.value = Math.floor(v || 0) })
function onSeekInput(e){ dragging.value = true; localPos.value = Number(e.target.value || 0) }
async function onSeekCommit(e){ const val = Number(e.target.value || 0); await store.seekTo(val); dragging.value = false }
const canJump = computed(() => (store.songList()?.length || 0) > 0)

// Lyrics
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

// Lifecycle & Timers
let timer = null
function setupPoll(){
  clearInterval(timer)
  if (store.isPlaying) timer = setInterval(() => store.refreshStatus().catch(() => {}), 600)
}
watch(() => store.isPlaying, setupPoll)
watch(() => store.currentSong?.id, async () => { hadCoverErr.value = false; setupPoll(); await ensureLyrics() })

onMounted(async () => {
  await store.refreshStatus() // Enter with latest status
  localPos.value = Math.floor(store.playerStatus.position || 0)
  setupPoll()
  await ensureLyrics()
  window.addEventListener('keydown', onEsc)
})
onBeforeUnmount(() => { clearInterval(timer); window.removeEventListener('keydown', onEsc) })

function onEsc(e){ if(e.key === 'Escape') router.back() }

function formatTime(seconds){
  const s = Math.max(0, Math.floor(seconds || 0))
  const mins = Math.floor(s / 60); const secs = s % 60
  return `${String(mins).padStart(2,'0')}:${String(secs).padStart(2,'0')}`
}
</script>

<style scoped>
.np-page { position: fixed; inset: 0; z-index: 2000; width:100%; height:100%; overflow:hidden; }
.bg { position:absolute; inset:0; background-position:center; background-size:cover; filter: blur(22px) saturate(160%); transform: scale(1.06); opacity:.95; z-index:0; pointer-events:none; }

.header { position: absolute; left: 16px; right: 16px; top: 12px; z-index: 20000; display:flex; align-items:center; gap: 12px; }
.back-btn, .gear-btn { width: 40px; height: 40px; border-radius: 10px; border:1px solid rgba(0,0,0,0.06); background: rgba(255,255,255,0.96); display:flex; align-items:center; justify-content:center; cursor:pointer; box-shadow: 0 6px 18px rgba(0,0,0,0.18); }
.back-btn:hover, .gear-btn:hover { background:#fff; }
.titlebar { display:flex; flex-direction:column; gap:4px; background: rgba(255,255,255,0.6); padding: 6px 10px; border-radius: 8px; border: 1px solid rgba(0,0,0,0.06); }
.song-title { font-weight:700; color:#222; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; max-width: 52vw; }
.song-artist { font-size:12px; color:#555; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; max-width: 52vw; }
.header-spacer { flex: 1 1 auto; }

.content { position: relative; z-index:1; display:grid; grid-template-columns: 420px 1fr; gap: 18px; height:100%; padding: 72px 24px 24px; box-sizing: border-box; }

.left { display:flex; flex-direction:column; align-items:center; justify-content: center; }
.cover-wrap { width: 100%; max-width: 380px; aspect-ratio: 1; border-radius: 12px; background: rgba(0,0,0,0.08); overflow:hidden; display:flex; align-items:center; justify-content:center; box-shadow: 0 20px 50px rgba(0,0,0,0.18); }
.cover-wrap img { width:100%; height:100%; object-fit:cover; display:block; }
.cover-placeholder { font-size:48px; opacity:.8; }
.cover-reflect { width:100%; max-width:380px; height:64px; overflow:hidden; mask-image: linear-gradient(to bottom, rgba(0,0,0,.35), rgba(0,0,0,0)); -webkit-mask-image: linear-gradient(to bottom, rgba(0,0,0,.35), rgba(0,0,0,0)); }
.cover-reflect img { width:100%; height:128px; object-fit:cover; transform: scaleY(-1); opacity:.35; filter: blur(2px); transform-origin: top; }

.right { display:flex; flex-direction:column; min-width:0; }
/* 修复：移除 overflow:hidden，允许子元素滚动 */
.lyric-wrap { flex:1; min-height:0; background: rgba(255,255,255,0.55); border:1px solid rgba(0,0,0,0.06); border-radius: 12px; backdrop-filter: blur(10px) saturate(160%); -webkit-backdrop-filter: blur(10px) saturate(160%); }
.no-lyrics { height:100%; display:flex; align-items:center; justify-content:center; color:#666; }

.controls-wrapper { margin-top: 16px; }
.seek { display:grid; grid-template-columns: 56px 1fr 56px; gap:10px; align-items:center; }
.time { color:#555; font-variant-numeric: tabular-nums; text-align:center; }
.seek-bar { width:100%; accent-color:#667eea; }

.controls { display:flex; gap:16px; justify-content:center; padding-top: 12px; }
.ctrl { width:44px; height:44px; border-radius:50%; border:1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.8); cursor:pointer; display:flex; align-items:center; justify-content:center; }
.ctrl:hover { background: rgba(255,255,255,0.95); }
.ctrl.play-pause { width: 56px; height: 56px; }

@media (max-width: 980px) {
  .content { grid-template-columns: 1fr; padding-top: 68px; }
  .left { display: none; }
}
</style>
