<template>
  <div class="topbar">
    <div class="left">
      <div class="title">{{ title }}</div>
      <div class="count" v-if="count !== undefined">{{ count }}</div>
    </div>
    <div class="center">
      <button class="chip" :class="{ active: store.playMode==='loop' }" @click="store.setPlayMode('loop')">列表循环</button>
      <button class="chip" :class="{ active: store.playMode==='shuffle' }" @click="store.setPlayMode('shuffle')">随机播放</button>
    </div>
    <div class="right">
      <button class="icon-btn" title="补全时长" @click="onRefreshDurations">⟳</button>
      <button class="icon-btn" title="排序（占位）">⇅</button>
      <button class="icon-btn" title="搜索（占位）">🔍</button>
    </div>
  </div>
</template>

<script setup>
import { usePlayerStore } from '../stores/player'
import { refreshDurations } from '../api/music'

const props = defineProps({
  title: { type: String, default: '歌曲' },
  count: { type: Number, default: undefined }
})

const store = usePlayerStore()

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
</script>

<style scoped>
.topbar { display:flex; align-items:center; justify-content:space-between; gap: 12px; height: 56px; padding: 0 12px; }
.left { display:flex; align-items:baseline; gap:10px; }
.title { font-size: 22px; font-weight: 700; color:#222; }
.count { font-size: 14px; color:#666; }
.center { display:flex; gap:8px; }
.chip { padding: 8px 12px; border-radius: 999px; border: 1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.7); cursor:pointer; }
.chip.active { background: linear-gradient(135deg,#667eea22,#764ba222); border-color: rgba(0,0,0,0.15); }
.right { display:flex; gap:8px; }
.icon-btn { width: 36px; height: 36px; border-radius: 8px; border:1px solid rgba(0,0,0,0.08); background: rgba(255,255,255,0.7); cursor: pointer; }
.icon-btn:hover, .chip:hover { background: rgba(255,255,255,0.9); }
</style>
