<template>
  <div class="np-page" role="dialog" aria-modal="true">
    <div class="bg" :style="bgStyle" aria-hidden="true"></div>

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

    <LyricControls v-if="showLyricCtrl" :right="16" :top="56" @close="showLyricCtrl = false" />

    <div class="content">
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

      <div class="right">
        <div class="lyric-wrap">
                    <LyricDisplay
            v-if="hasLyrics"
            :lyrics="store.lyrics"
            :currentTime="store.playerStatus.position"
            :baseFontSize="lyricUi.fontSize || 20"
            :fontWeight="lyricUi.fontWeight || 400"
            :blurOthers="lyricUi.blurOthers === true"
            :isPlaying="store.isPlaying"
            :showTranslation="lyricUi.showTranslation === true"
            :translationScale="lyricUi.translationScale || 80"
            @seek="onLyricSeek"
          />
          <div v-else class="no-lyrics">暂无歌词</div>
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
import { useSettingsStore } from '../stores/settings'
import { getLyrics } from '../api/music'
import LyricDisplay from '../components/LyricDisplay.vue'
import LyricControls from '../components/LyricControls.vue'

const router = useRouter()
const store = usePlayerStore()
const lyricUi = useLyricUiStore()
const settings = useSettingsStore()
lyricUi.load()
settings.load()

const showLyricCtrl = ref(false)

const title = computed(() => store.currentSong?.title || '未选择歌曲')
const artist = computed(() => store.currentSong?.artist || '未知歌手')
const hadCoverErr = ref(false)
const coverSrc = computed(() => {
  const s = store.currentSong
  if (!s || hadCoverErr.value) return ''
  return s.id ? `http://localhost:8080/api/cover/${s.id}` : (s.cover_url || '')
})
function onCoverErr(){ hadCoverErr.value = true }

const bgStyle = computed(() => ({
  backgroundImage: settings.nowPlayingBackgroundUrl ? `url(${settings.nowPlayingBackgroundUrl})` : (coverSrc.value ? `url(${coverSrc.value})` : 'none'),
  '--np-bg-blur': `${lyricUi.backgroundBlur || 22}px`
}))
const hasLyrics = computed(() => store.lyrics && Array.isArray(store.lyrics.lines) && store.lyrics.lines.length > 0)

async function ensureLyrics(){
  try{
    if(!store.currentSong) return
    const need = !store.lyrics || !Array.isArray(store.lyrics?.lines) || store.lyrics.lines.length === 0
    if(need){
      const { data } = await getLyrics(store.currentSong.id)
      if(data) store.lyrics = data
    }
  }catch(_){}
}

let timer = null
function setupPoll(){
  clearInterval(timer)
  if (store.isPlaying) timer = setInterval(() => store.refreshStatus().catch(() => {}), 600)
}
watch(() => store.isPlaying, setupPoll)
watch(() => store.currentSong?.id, async () => { hadCoverErr.value = false; setupPoll(); await ensureLyrics() })

onMounted(async () => {
  await store.refreshStatus()
  setupPoll()
  await ensureLyrics()
  window.addEventListener('keydown', onEsc)
})
onBeforeUnmount(() => { clearInterval(timer); window.removeEventListener('keydown', onEsc) })

function onEsc(e){ if(e.key === 'Escape') router.back() }

function onLyricSeek(time) {
  store.seekTo(time)
}

</script>

<style scoped>
.np-page { position: fixed; inset: 0; z-index: 2000; width:100%; height:100%; overflow:hidden; }
.bg { position:absolute; inset:0; background-color: #000; background-position:center; background-size:cover; filter: blur(var(--np-bg-blur, 22px)) saturate(160%); transform: scale(1.06); opacity:1; z-index:0; pointer-events:none; }

.header { position: absolute; left: 16px; right: 16px; top: 12px; z-index: 20001; display:flex; align-items:center; gap: 12px; }
.back-btn, .gear-btn { width: 40px; height: 40px; border-radius: 10px; border:1px solid rgba(0,0,0,0.06); background: rgba(255,255,255,0.96); display:flex; align-items:center; justify-content:center; cursor:pointer; box-shadow: 0 6px 18px rgba(0,0,0,0.18); }
.back-btn:hover, .gear-btn:hover { background:#fff; }
.titlebar { display:flex; flex-direction:column; gap:4px; background: rgba(255,255,255,0.6); padding: 6px 10px; border-radius: 8px; border: 1px solid rgba(0,0,0,0.06); }
.song-title { font-weight:700; color:#222; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; max-width: 52vw; }
.song-artist { font-size:12px; color:#555; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; max-width: 52vw; }
.header-spacer { flex: 1 1 auto; }

.content { position: relative; z-index:1; display:grid; grid-template-columns: 2fr 3fr; gap: 18px; height:100%; padding: 72px 24px 24px; box-sizing: border-box; }

.left { display:flex; flex-direction:column; align-items:center; justify-content: center; }
.cover-wrap { width: 100%; max-width: 380px; aspect-ratio: 1; border-radius: 12px; background: rgba(0,0,0,0.08); overflow:hidden; display:flex; align-items:center; justify-content:center; box-shadow: 0 20px 50px rgba(0,0,0,0.18); }
.cover-wrap img { width:100%; height:100%; object-fit:cover; display:block; }
.cover-placeholder { font-size:48px; opacity:.8; }
.cover-reflect { width:100%; max-width:380px; height:64px; overflow:hidden; mask-image: linear-gradient(to bottom, rgba(0,0,0,.35), rgba(0,0,0,0)); -webkit-mask-image: linear-gradient(to bottom, rgba(0,0,0,.35), rgba(0,0,0,0)); }
.cover-reflect img { width:100%; height:128px; object-fit:cover; transform: scaleY(-1); opacity:.35; filter: blur(2px); transform-origin: top; }

.right { display:flex; flex-direction:column; min-width:0; justify-content: center; }
.lyric-wrap { flex: 0 1 600px; min-height: 200px; mask-image: linear-gradient(to bottom, transparent, black 20%, black 80%, transparent); -webkit-mask-image: linear-gradient(to bottom, transparent, black 20%, black 80%, transparent); }
.no-lyrics { height:100%; display:flex; align-items:center; justify-content:center; color:#fff; }

@media (max-width: 980px) {
  .content { grid-template-columns: 1fr; padding-top: 68px; }
  .left { display: none; }
}
</style>