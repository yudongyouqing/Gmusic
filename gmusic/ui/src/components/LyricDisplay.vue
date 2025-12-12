<template>
  <div class="lyric-display" ref="containerRef">
    <div class="lyric-list" :style="scrollStyle">
      <div
        v-for="(line, idx) in lines"
        :key="idx"
        :ref="el => { if (el) lineRefs[idx] = el }"
        class="lyric-line"
        :class="{ current: idx === currentLineIndex }"
      >
        {{ line.text || '\u00A0' }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, onBeforeUpdate, nextTick, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  lyrics: { type: Object, default: null },
  currentTime: { type: Number, default: 0 }, // 秒
  anchorRatio: { type: Number, default: 0.30 } // 当前行相对容器高度的锚点位置(0~1)，默认略高于中线
})

const lines = computed(() => props.lyrics?.lines || [])
const currentLineIndex = computed(() => {
  const ctMs = Math.floor((props.currentTime || 0) * 1000)
  // 先找到时间上“应该到达的行”
  let raw = -1
  for (let i = 0; i < lines.value.length; i++) {
    if (lines.value[i].time <= ctMs) raw = i
    else break
  }
  if (raw < 0) return 0
  // 再向前回退到最近一条非空歌词，实现“无缝链接”
  let j = raw
  while (j > 0 && !(lines.value[j]?.text || '').trim()) j--
  return j
})

const containerRef = ref(null)
const lineRefs = ref([])
const scrollY = ref(0)
const noAnim = ref(true) // 尺寸变化/初始对齐时不做动画，避免“莫名滑动”

// 在 DOM 更新前清空 refs
onBeforeUpdate(() => { lineRefs.value = [] })

function centerTo(idx, animate){
  const container = containerRef.value
  const activeLine = lineRefs.value[idx]
  if (!container || !activeLine) return
  const containerHeight = container.clientHeight
  const lineTop = activeLine.offsetTop
  const lineHeight = activeLine.offsetHeight
  const anchor = Math.min(0.9, Math.max(0.1, props.anchorRatio ?? 0.35))
  // 让当前行“中线”对齐到容器的 anchor 比例位置（默认 35% 稍偏上）
  const targetY = containerHeight * anchor
  let offset = targetY - lineTop - lineHeight / 2

  // 边界钳制：不允许第一行被推到视野上方或最后一行被推到下方
  const first = lineRefs.value[0]
  const last = lineRefs.value[lineRefs.value.length - 1]
  const padTop = 8, padBottom = 8
  if (first && last) {
    const minOffset = padTop - first.offsetTop // 使第一行顶到 padTop
    const maxOffset = (containerHeight - padBottom) - (last.offsetTop + last.offsetHeight) // 使最后一行底到容器底部上方 padBottom
    // 注意：minOffset 通常是正值（往下推），maxOffset 通常是负值（往上拉）。
    // 我们希望 offset 处于 [maxOffset, minOffset] 区间内。
    offset = Math.max(maxOffset, Math.min(minOffset, offset))
  }

  noAnim.value = !animate
  scrollY.value = offset
}

const lastIdx = ref(-1)
const lastTime = ref(0)

watch(() => props.currentTime, (t)=>{ lastTime.value = t || 0; nextTick(() => ensureVisible()) })

watch(currentLineIndex, (newIdx) => {
  // 小步前进时才使用动画；大跳转/回退/初始化均无动画，避免“莫名一大段滑动”
  const delta = lastIdx.value >= 0 ? (newIdx - lastIdx.value) : 0
  const animate = (delta === 1) && (props.currentTime >= lastTime.value)
  lastIdx.value = newIdx
  nextTick(() => centerTo(newIdx, animate))
})

// 歌词整体刷新：无动画对齐
watch(() => props.lyrics, () => nextTick(() => centerTo(currentLineIndex.value, false)))

let ro = null
onMounted(() => {
  nextTick(() => centerTo(currentLineIndex.value, false))
  // 监听容器尺寸变化：无动画居中，避免字体/布局变化导致“慢慢下滑”
  if (window.ResizeObserver) {
    ro = new ResizeObserver(() => centerTo(currentLineIndex.value, false))
    if (containerRef.value) ro.observe(containerRef.value)
  } else {
    window.addEventListener('resize', onResize, { passive: true })
  }
})
function onResize(){ centerTo(currentLineIndex.value, false) }

onBeforeUnmount(() => {
  if (ro) { try{ ro.disconnect() }catch(_){} ro = null }
  window.removeEventListener('resize', onResize)
})

const scrollStyle = computed(() => ({
  transform: `translateY(${scrollY.value}px)`,
  transition: noAnim.value ? 'none' : 'transform 0.35s cubic-bezier(.2,.7,.2,1)',
  willChange: 'transform'
}))
</script>

<style scoped>
.lyric-display { height: 100%; overflow: hidden; position: relative; }
.lyric-list { width: 100%; }

.lyric-line {
  padding: 12px 10px;
  text-align: center;
  color: #666;
  font-size: 16px;
  line-height: 1.6;
  transition: color .25s ease, font-size .25s ease, transform .25s ease;
}

.lyric-line.current { color: #222; font-weight: 600; font-size: 18px; transform: translateZ(0) scale(1.02); }
</style>
