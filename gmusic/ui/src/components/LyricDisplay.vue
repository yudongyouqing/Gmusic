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
        {{ line.text || '&nbsp;' }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, onBeforeUpdate } from 'vue'

const props = defineProps({
  lyrics: { type: Object, default: null },
  currentTime: { type: Number, default: 0 } // 秒
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

// 在 DOM 更新前清空 refs
onBeforeUpdate(() => {
  lineRefs.value = []
})

watch(currentLineIndex, (newIdx) => {
  const container = containerRef.value
  const activeLine = lineRefs.value[newIdx]
  if (!container || !activeLine) return

  const containerHeight = container.clientHeight
  const lineTop = activeLine.offsetTop
  const lineHeight = activeLine.offsetHeight

  // 计算目标偏移量，使当前行垂直居中
  const offset = containerHeight / 2 - lineTop - lineHeight / 2
  scrollY.value = offset
})

const scrollStyle = computed(() => ({
  transform: `translateY(${scrollY.value}px)`
}))
</script>

<style scoped>
.lyric-display {
  height: 100%;
  overflow: hidden;
  position: relative;
}

.lyric-list {
  width: 100%;
  transition: transform 0.5s ease-out; /* 平滑滚动动画 */
}

.lyric-line {
  padding: 12px 10px;
  text-align: center;
  color: #666;
  font-size: 16px;
  line-height: 1.6;
  transition: all 0.3s ease;
}

.lyric-line.current {
  color: #222;
  font-weight: 600;
  font-size: 18px;
  transform: scale(1.05);
}
</style>
