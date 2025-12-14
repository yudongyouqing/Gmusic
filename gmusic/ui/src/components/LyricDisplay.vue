<template>
  <div class="lyric-display">
    <!-- 原生滚动容器，通过 CSS pointer-events 控制是否可滚动 -->
    <div ref="scrollRef" class="ld-scroll" :class="{ 'is-locked': isPlaying }">
      <div
        v-for="(line, idx) in lines"
        :key="idx"
        :ref="el => { if (el) lineRefs[idx] = el }"
        class="lyric-line"
        :class="{ current: idx === currentLineIndex, blur: blurOthers && idx !== currentLineIndex }"
      >
        {{ line.text || '\u00A0' }}
      </div>
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

const lines = computed(() => props.lyrics?.lines || [])

const currentLineIndex = computed(() => {
  const ctMs = Math.floor((props.currentTime || 0) * 1000)
  let raw = -1
  const arr = lines.value
  for (let i = 0; i < arr.length; i++) {
    if ((arr[i]?.time || 0) <= ctMs) raw = i
    else break
  }
  if (raw < 0) return 0
  let j = raw
  while (j > 0 && !((arr[j]?.text || '').trim())) j--
  return j
})

const scrollRef = ref(null)
const lineRefs = ref([])

function scrollToCurrent(animate = true) {
  const el = scrollRef.value
  const row = lineRefs.value[currentLineIndex.value]
  if (!el || !row) return

  const targetScrollTop = row.offsetTop - (el.clientHeight / 2 - row.offsetHeight / 2)

  if (Math.abs(targetScrollTop - el.scrollTop) < 1) return

  if (animate) {
    el.scrollTo({ top: targetScrollTop, behavior: 'smooth' })
  } else {
    el.scrollTop = targetScrollTop
  }
}

// 核心逻辑：只要行号变化，就滚动到当前行
watch(currentLineIndex, () => {
  // 仅当播放时才自动滚动
  if (props.isPlaying) {
    nextTick(() => scrollToCurrent(true))
  }
})

// 播放状态变化：暂停时可滚动，播放时立即跳回并锁定
watch(() => props.isPlaying, (playing) => {
  if (playing) {
    nextTick(() => scrollToCurrent(false))
  }
})

// 初始化与尺寸变化时，无动画对齐
function initialAlign() {
  nextTick(() => scrollToCurrent(false))
}

let ro = null
onMounted(() => {
  initialAlign()
  if (window.ResizeObserver) {
    ro = new ResizeObserver(initialAlign)
    if (scrollRef.value) ro.observe(scrollRef.value)
  } else {
    window.addEventListener('resize', initialAlign, { passive: true })
  }
})

onBeforeUnmount(() => {
  if (ro) { try { ro.disconnect() } catch {} ro = null }
  window.removeEventListener('resize', initialAlign)
})
</script>

<style scoped>
.lyric-display { position: relative; height: 100%; overflow: hidden; }
.ld-scroll {
  height: 100%;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  padding: 0 8px; /* 左右留白，上下不留白以精确居中 */
  scroll-behavior: smooth;
}

/* 播放时锁定滚动，用户无法手动操作 */
.ld-scroll.is-locked {
  pointer-events: none;
}

.lyric-line {
  padding: 12px 10px;
  text-align: center;
  color:#666;
  font-size: var(--lyric-base-size,16px);
  line-height:1.6;
  transition: color .25s ease, font-size .25s ease, filter .25s ease, opacity .25s ease;
}
.lyric-line.current { color:#222; font-weight:600; font-size: calc(var(--lyric-base-size,16px) + 2px); }
.lyric-line.blur { filter: blur(1px); opacity: .6; }
</style>
