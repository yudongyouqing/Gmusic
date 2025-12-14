<template>
  <div class="lyric-display">
    <div ref="scrollRef" class="ld-scroll" :class="{ 'is-locked': isPlaying }">
      <!-- 虚拟滚动：仅渲染可视区域附近的歌词行 -->
      <div class="scroll-spacer" :style="{ height: `${topPadding}px` }"></div>
      <div
        v-for="line in visibleLines"
        :key="line.idx"
        :ref="el => { if (el) lineRefs[line.idx] = el }"
        class="lyric-line"
        :class="{ current: line.idx === currentLineIndex, blur: blurOthers && line.idx !== currentLineIndex }"
      >
        {{ line.text || '\u00A0' }}
      </div>
      <div class="scroll-spacer" :style="{ height: `${bottomPadding}px` }"></div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  lyrics: { type: Object, default: null },
  currentTime: { type: Number, default: 0 },
  baseFontSize: { type: Number, default: 16 },
  blurOthers: { type: Boolean, default: false },
  isPlaying: { type: Boolean, default: false }
})

const allLines = computed(() => props.lyrics?.lines || [])

const currentLineIndex = computed(() => {
  const ctMs = Math.floor((props.currentTime || 0) * 1000)
  let raw = -1
  for (let i = 0; i < allLines.value.length; i++) {
    if ((allLines.value[i]?.time || 0) <= ctMs) raw = i
    else break
  }
  if (raw < 0) return 0
  let j = raw
  while (j > 0 && !((allLines.value[j]?.text || '').trim())) j--
  return j
})

const scrollRef = ref(null)
const lineRefs = ref([])

// --- 虚拟滚动核心 ---
const visibleRange = ref({ start: 0, end: 20 })
const topPadding = ref(0)
const bottomPadding = ref(0)
const averageLineHeight = ref(38) // 预估行高

const visibleLines = computed(() => {
  return allLines.value.slice(visibleRange.value.start, visibleRange.value.end).map((line, i) => ({
    ...line,
    idx: visibleRange.value.start + i
  }))
})

function updateVisibleRange() {
  const el = scrollRef.value
  if (!el) return

  const buffer = 10 // 上下多渲染10行
  const scrollTop = el.scrollTop
  const clientHeight = el.clientHeight

  const startIndex = Math.max(0, Math.floor(scrollTop / averageLineHeight.value) - buffer)
  const endIndex = Math.min(allLines.value.length, Math.ceil((scrollTop + clientHeight) / averageLineHeight.value) + buffer)

  visibleRange.value = { start: startIndex, end: endIndex }
  topPadding.value = startIndex * averageLineHeight.value
  bottomPadding.value = (allLines.value.length - endIndex) * averageLineHeight.value
}

function onScroll() {
  if (!props.isPlaying) {
    updateVisibleRange()
  }
}
// --- 虚拟滚动结束 ---

function scrollToCurrent(animate = true) {
  const el = scrollRef.value
  if (!el) return

  const targetIndex = currentLineIndex.value
  const targetScrollTop = targetIndex * averageLineHeight.value - (el.clientHeight / 2 - averageLineHeight.value / 2)

  if (Math.abs(targetScrollTop - el.scrollTop) < 1) return

  if (animate) {
    el.scrollTo({ top: targetScrollTop, behavior: 'smooth' })
  } else {
    el.scrollTop = targetScrollTop
  }
  // 滚动后立即更新可视范围
  nextTick(updateVisibleRange)
}

watch(currentLineIndex, () => {
  if (props.isPlaying) {
    scrollToCurrent(true)
  }
})

watch(() => props.isPlaying, (playing) => {
  if (playing) {
    scrollToCurrent(false)
  }
})

function initialAlign() {
  nextTick(() => {
    updateVisibleRange()
    scrollToCurrent(false)
  })
}

let ro = null
onMounted(() => {
  scrollRef.value?.addEventListener('scroll', onScroll)
  initialAlign()
  if (window.ResizeObserver) {
    ro = new ResizeObserver(initialAlign)
    if (scrollRef.value) ro.observe(scrollRef.value)
  } else {
    window.addEventListener('resize', initialAlign, { passive: true })
  }
})

onBeforeUnmount(() => {
  scrollRef.value?.removeEventListener('scroll', onScroll)
  if (ro) { try { ro.disconnect() } catch {} ro = null }
  window.removeEventListener('resize', initialAlign)
})

</script>

<style scoped>
.lyric-display { position: relative; height: 100%; }
.ld-scroll {
  height: 100%;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  scroll-behavior: smooth;
}

/* 播放时锁定滚动，用户无法手动操作 */
.ld-scroll.is-locked {
  pointer-events: none;
}

.lyric-line {
  padding: 10px;
  text-align: center;
  color:#666;
  font-size: var(--lyric-base-size, 16px);
  line-height:1.6;
  transition: color .25s ease, font-size .25s ease, filter .25s ease, opacity .25s ease;
}
.lyric-line.current { color:#222; font-weight:600; font-size: calc(var(--lyric-base-size, 16px) + 2px); }
.lyric-line.blur { filter: blur(1px); opacity: .6; }
</style>
