<template>
  <teleport to="body">
    <div v-if="visible" class="np-overlay" @click.stop>
      <div class="np-bg" :class="{ enter: started && mode==='expand' }" :style="bgStyle"></div>
      <img class="np-cover" :src="coverSrc" :style="coverStyle" @transitionend="onCoverEnd" />
      <div class="np-lyric-panel" :class="{ enter: started }"></div>
    </div>
  </teleport>
</template>

<script setup>
import { onMounted, ref, computed, nextTick } from 'vue'

const bgStyle = computed(() => ({
  backgroundImage: props.coverSrc ? `url(${props.coverSrc})` : 'none'
}))

const props = defineProps({
  startRect: { type: Object, required: true },
  targetRect: { type: Object, default: null },
  mode: { type: String, default: 'expand' }, // 'expand' | 'collapse'
  coverSrc: { type: String, default: '' }
})
const emit = defineEmits(['done'])

const visible = ref(true)
const started = ref(false)
let finished = false

const coverStyle = ref({})

function computeTargetRect(){
  if (props.targetRect) return props.targetRect
  const vw = window.innerWidth
  const vh = window.innerHeight
  const padX = Math.max(24, vw * 0.06)
  const size = Math.min(vw * 0.38, 520)
  const left = padX
  const top = Math.max(24, (vh - size) / 2)
  return { left, top, width: size, height: size }
}

onMounted(async () => {
  const s = props.startRect || { left: 0, top: 0, width: 0, height: 0 }
  // 初始：放在起点
  coverStyle.value = {
    position: 'fixed',
    left: s.left + 'px', top: s.top + 'px', width: s.width + 'px', height: s.height + 'px',
    borderRadius: '12px', objectFit: 'cover', zIndex: 1,
    transition: 'all 460ms cubic-bezier(.2,.7,.2,1)'
  }
  await nextTick()
  // 目标：根据模式决定
  const t = computeTargetRect()
  coverStyle.value = {
    ...coverStyle.value,
    left: t.left + 'px', top: t.top + 'px', width: t.width + 'px', height: t.height + 'px'
  }
  // 背景与右侧面板淡入/淡出
  started.value = true
})

function onCoverEnd(){
  if (finished) return
  finished = true
  emit('done')
  visible.value = false
}
</script>

<style scoped>
.np-overlay{ position: fixed; inset: 0; z-index: 5000; pointer-events: none; }
.np-bg{ position: absolute; inset: 0; opacity: 0; transition: opacity 380ms ease; filter: blur(18px) saturate(160%); background-position: center; background-size: cover; }
.np-bg.enter{ opacity: .95; }

.np-cover{ box-shadow: 0 20px 50px rgba(0,0,0,.28); border-radius: 12px; }

/* 右侧歌词面板的占位淡入，仅做氛围动画，真正歌词在路由页 */
.np-lyric-panel{ position: fixed; right: 6%; top: 15%; width: 42%; height: 70%; border-radius: 16px; background: rgba(255,255,255,0.08); border: 1px solid rgba(255,255,255,0.15); backdrop-filter: blur(12px) saturate(160%); opacity: 0; transform: translateY(10px); transition: all 380ms ease 80ms; }
.np-lyric-panel.enter{ opacity: 1; transform: translateY(0); }
</style>

