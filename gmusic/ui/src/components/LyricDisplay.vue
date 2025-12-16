<template>
  <div class="lyric-display">
    <div ref="scrollRef" class="ld-scroll">
      <div
        v-for="(line, idx) in lines"
        :key="idx"
        :ref="el => { if (el) lineRefs[idx] = el }"
        class="lyric-line"
        :style="lineStyle"
        :class="{ current: idx === currentLineIndex, blur: blurOthers && idx !== currentLineIndex }"
        @click="onLineClick(line)"
      >
        <span>{{ line.text || '\u00A0' }}</span>
        <span v-if="showTranslation" class="translation">(翻译占位)</span>
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
  fontWeight: { type: Number, default: 400 },
  blurOthers: { type: Boolean, default: false },
  isPlaying: { type: Boolean, default: false },
  showTranslation: { type: Boolean, default: false },
  translationScale: { type: Number, default: 80 },
})

const emit = defineEmits(['seek'])

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

const lineStyle = computed(() => ({
  '--lyric-base-size': `${props.baseFontSize}px`,
  '--lyric-font-weight': props.fontWeight,
  '--lyric-translation-scale': `${props.translationScale / 100}`
}))

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

watch(currentLineIndex, () => {
  if (props.isPlaying) {
    nextTick(() => scrollToCurrent(true))
  }
})

watch(() => props.isPlaying, (playing) => {
  if (playing) {
    nextTick(() => scrollToCurrent(false))
  }
})

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

function onLineClick(line) {
  if (line && typeof line.time === 'number') {
    emit('seek', line.time / 1000)
  }
}
</script>

<style scoped>
.lyric-display { position: relative; height: 100%; }
.ld-scroll {
  height: 100%;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  padding: 0 8px;
  scroll-behavior: smooth;
  /* Hide scrollbar */
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
}
.ld-scroll::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}

.lyric-line {
  padding: 10px;
  text-align: center;
  color: #fff;
  font-size: var(--lyric-base-size, 16px);
  font-weight: var(--lyric-font-weight, 400);
  line-height: 1.6;
  transition: color .25s ease, font-size .25s ease, font-weight .25s ease, filter .25s ease, opacity .25s ease;
  cursor: pointer;
}

.lyric-line:hover {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.lyric-line.current {
  color: #fff;
  font-weight: bold;
  font-size: calc(var(--lyric-base-size, 16px) + 2px);
  text-shadow: 0 0 5px rgba(255, 255, 255, 0.5);
}

.lyric-line.blur {
  filter: blur(1px);
  opacity: .5;
}

.translation {
  display: block;
  margin-top: 4px;
  font-size: calc(var(--lyric-base-size, 16px) * var(--lyric-translation-scale, 0.8));
  opacity: 0.8;
}
</style>
