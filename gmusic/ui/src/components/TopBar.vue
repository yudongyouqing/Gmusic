<template>
  <div class="topbar">
    <div class="left">
      <div class="title">{{ title }}</div>
      <div class="count" v-if="count !== undefined">{{ count }}</div>
    </div>
    <div class="center">
      <button class="chip" :title="playModeText" @click="store.togglePlayMode()">
        <span class="chip-inner" v-if="store.playMode === 'loop'">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M8 6h8a4 4 0 0 1 4 4v1" />
            <path d="M20 11l-2-2m2 2l-2 2" />
            <path d="M16 18H8a4 4 0 0 1-4-4v-1" />
            <path d="M4 13l2-2m-2 2l2 2" />
          </svg>
          <span>列表循环</span>
        </span>
        <span class="chip-inner" v-else-if="store.playMode === 'shuffle'">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M4 6h4l4 6h4" />
            <path d="M20 6l-2-2m2 2l-2 2" />
            <path d="M4 18h4l4-6h4" />
            <path d="M20 18l-2-2m2 2l-2 2" />
          </svg>
          <span>随机播放</span>
        </span>
        <span class="chip-inner" v-else-if="store.playMode === 'single'">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M8 6h8a4 4 0 0 1 4 4v1" />
            <path d="M20 11l-2-2m2 2l-2 2" />
            <path d="M16 18H8a4 4 0 0 1-4-4v-1" />
            <path d="M4 13l2-2m-2 2l2 2" />
            <path d="M12 8v8" />
          </svg>
          <span>单曲循环</span>
        </span>
      </button>
    </div>
    <div class="right" ref="rightRef">
      <button class="icon-btn" title="补全时长" @click="onRefreshDurations">
        <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
          <path d="M20 6v6h-6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M20 12a8 8 0 1 1-2.34-5.66L20 8" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
      <button class="icon-btn" ref="sortBtnRef" title="排序" @click.stop="toggleSort">
        <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
          <path d="M8 6h10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <path d="M8 12h6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <path d="M8 18h3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <path d="M5 7v10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <path d="M3 9l2-2 2 2" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M3 15l2 2 2-2" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
      <div v-if="showSort" class="sort-menu">
        <div class="group">
          <div class="title">标题</div>
          <button @click="chooseSort('title','asc')">升序</button>
          <button @click="chooseSort('title','desc')">降序</button>
        </div>
        <div class="group">
          <div class="title">歌手</div>
          <button @click="chooseSort('artist','asc')">升序</button>
          <button @click="chooseSort('artist','desc')">降序</button>
        </div>
        <div class="group">
          <div class="title">专辑</div>
          <button @click="chooseSort('album','asc')">升序</button>
          <button @click="chooseSort('album','desc')">降序</button>
        </div>
        <div class="group">
          <button class="wide" @click="chooseSort('custom',store.sortDir)">自定义排序（拖拽）</button>
        </div>
      </div>
      <button class="icon-btn" title="搜索">
        <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
          <circle cx="11" cy="11" r="7" stroke="currentColor" stroke-width="2" fill="none"/>
          <path d="M20 20l-3.5-3.5" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>
      <button class="icon-btn" ref="themeBtnRef" title="主题/毛玻璃" @click.stop="toggleTheme">
        <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
          <path d="M13.5 19H12a8 8 0 1 1 0-16h.5a3 3 0 1 1 0 6H11a2 2 0 1 0 0 4h2.5a2.5 2.5 0 1 1 0 5z" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <circle cx="8" cy="11" r="1.2" fill="currentColor"/>
          <circle cx="11" cy="8" r="1.2" fill="currentColor"/>
          <circle cx="14.5" cy="10" r="1.2" fill="currentColor"/>
          <circle cx="15.5" cy="13" r="1.2" fill="currentColor"/>
        </svg>
      </button>
    </div>

    <!-- Teleport 到 body，避免被父级 overflow/backdrop-filter 裁剪 -->
    <teleport to="body">
      <ThemeSwitcher v-if="showTheme" :pos="themeBtnRect" :z="5000" @close="showTheme = false" />
    </teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { usePlayerStore } from '../stores/player'
import { refreshDurations } from '../api/music'
import ThemeSwitcher from './ThemeSwitcher.vue'

const props = defineProps({
  title: { type: String, default: '歌曲' },
  count: { type: Number, default: undefined }
})

const store = usePlayerStore()

const playModeText = computed(() => {
  switch (store.playMode) {
    case 'loop': return '列表循环'
    case 'shuffle': return '随机播放'
    case 'single': return '单曲循环'
    default: return ''
  }
})

async function onRefreshDurations() {
  try {
    const { data } = await refreshDurations()
    await store.fetchSongs()
    alert(`补全完成：共 ${data.total} 条，更新 ${data.updated} 条，跳过 ${data.skipped} 条`)
  } catch (e) {
    const msg = e?.response?.data?.error || e?.message || '失败'
    alert(`补全失败：${msg}`)
  }
}

// 主题开关 + 锚点定位
const showTheme = ref(false)
const rightRef = ref(null)
const themeBtnRef = ref(null)
const themeBtnRect = ref(null)

// 排序菜单
const showSort = ref(false)
const sortBtnRef = ref(null)
function toggleSort(){ showSort.value = !showSort.value }
function chooseSort(mode, dir){ store.setSort(mode, dir); showSort.value = false }

function updateThemeRect(){
  const btn = themeBtnRef.value
  if(!btn) return
  const rect = btn.getBoundingClientRect()
  themeBtnRect.value = { top: rect.top, left: rect.left, right: rect.right, height: rect.height }
}

function toggleTheme(){
  if(!showTheme.value){ updateThemeRect() }
  showTheme.value = !showTheme.value
}

function onClickOutside(e){
  // 单击页面其他位置关闭面板
  if(showTheme.value){
    const btn = themeBtnRef.value
    if(btn && !btn.contains(e.target)) showTheme.value = false
  }
  if(showSort.value){
    const r = rightRef.value
    if(r && !r.contains(e.target)) showSort.value = false
  }
}

function onWindow(){ if(showTheme.value){ updateThemeRect() } }

onMounted(()=>{
  document.addEventListener('click', onClickOutside)
  window.addEventListener('resize', onWindow)
  window.addEventListener('scroll', onWindow, true)
})

onBeforeUnmount(()=>{
  document.removeEventListener('click', onClickOutside)
  window.removeEventListener('resize', onWindow)
  window.removeEventListener('scroll', onWindow, true)
})
</script>

<style scoped>
.topbar { display:flex; align-items:center; justify-content:space-between; gap: 12px; height: 56px; padding: 0 12px; position: relative; }
.left { display:flex; align-items:baseline; gap:10px; }
.title { font-size: 22px; font-weight: 700; color:#222; }
.count { font-size: 14px; color:#666; }
.center { display:flex; gap:8px; }
.chip { padding: 8px 12px; border-radius: 999px; border: 1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.7); cursor:pointer; }
.chip.active { background: linear-gradient(135deg,#667eea22,#764ba222); border-color: rgba(0,0,0,0.15); }
.right { display:flex; gap:8px; position: relative; }
.icon-btn { width: 36px; height: 36px; border-radius: 8px; border:1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.7); cursor: pointer; }
.icon-btn:hover, .chip:hover { background: rgba(255,255,255,0.9); }
.chip .chip-inner { display:flex; align-items:center; gap:6px; }
.chip svg { display:block; }

/* 排序菜单 */
.sort-menu { position: absolute; top: 44px; right: 0; background: var(--mica-surface); backdrop-filter: blur(var(--mica-blur)) saturate(var(--mica-saturate)); -webkit-backdrop-filter: blur(var(--mica-blur)) saturate(var(--mica-saturate)); border:1px solid var(--mica-border); border-radius: 10px; padding: 8px; z-index: 4000; box-shadow: 0 10px 24px rgba(0,0,0,0.15); min-width: 220px; }
.sort-menu .group { display:flex; align-items:center; gap:8px; padding: 4px 0; }
.sort-menu .group .title { width: 42px; color:#666; font-size:12px; }
.sort-menu .group button { padding:6px 8px; border-radius:8px; border:1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.8); cursor:pointer; font-size:12px; }
.sort-menu .group button.wide { flex:1; text-align:center; }
.sort-menu .group button:hover { background: rgba(255,255,255,0.95); }
</style>
