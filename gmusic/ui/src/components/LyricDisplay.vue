<template>
  <div class="lyric-display">
    <div class="lyric-title" v-if="lyrics?.title || lyrics?.artist">
      <h4 v-if="lyrics?.title">{{ lyrics.title }}</h4>
      <p v-if="lyrics?.artist">{{ lyrics.artist }}</p>
    </div>

    <div class="lyric-lines" v-if="lines && lines.length">
      <div
        v-for="(line, idx) in windowLines"
        :key="startIdx + idx"
        class="lyric-line"
        :class="{ current: startIdx + idx === currentLineIndex }"
      >
        <span class="lyric-time">{{ line.time_str }}</span>
        <span class="lyric-text">{{ line.text }}</span>
      </div>
    </div>

    <div class="no-lyrics" v-else>暂无歌词</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import './LyricDisplay.css'

const props = defineProps({
  lyrics: { type: Object, default: null },
  currentTime: { type: Number, default: 0 } // 秒
})

const lines = computed(() => props.lyrics?.lines || [])
const currentLineIndex = computed(() => {
  const ctMs = Math.floor((props.currentTime || 0) * 1000)
  let idx = 0
  for (let i = 0; i < lines.value.length; i++) {
    if (lines.value[i].time <= ctMs) idx = i
    else break
  }
  return idx
})

const windowSize = 5
const startIdx = computed(() => Math.max(0, currentLineIndex.value - windowSize))
const endIdx = computed(() => Math.min(lines.value.length, currentLineIndex.value + windowSize + 1))
const windowLines = computed(() => lines.value.slice(startIdx.value, endIdx.value))
</script>

