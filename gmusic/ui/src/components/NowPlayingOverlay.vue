<template>
  <teleport to="body">
    <!-- 仅封面动画：无背景、无占位、不可拦截事件 -->
    <div v-if="visible" class="np-overlay" aria-hidden="true">
      <img class="np-cover" :src="coverSrc" :style="coverStyle" />
    </div>
  </teleport>
</template>

<script setup>
import { onMounted, ref, nextTick } from 'vue'

const props = defineProps({
  startRect: { type: Object, required: true },
  targetRect: { type: Object, default: null },
  mode: { type: String, default: 'expand' }, // 'expand' | 'collapse'
  coverSrc: { type: String, default: '' }
})
const emit = defineEmits(['done'])

const visible = ref(true)
let finished = false
const coverStyle = ref({})

function computeTargetRect(){
  if (props.targetRect) return props.targetRect
  const vw = window.innerWidth
  const vh = window.innerHeight
  const padX = Math.max(24, vw * 0.06)
  const size = Math.min(vw * 0.35, 520)
  const left = padX
  const top = Math.max(24, (vh - size) / 2)
  return { left, top, width: size, height: size }
}

onMounted(async () => {
  const s = props.startRect || { left: 0, top: 0, width: 0, height: 0 }
  // 初始：起点样式（更短过渡）
  coverStyle.value = {
    position: 'fixed',
    left: s.left + 'px', top: s.top + 'px', width: s.width + 'px', height: s.height + 'px',
    borderRadius: '12px', objectFit: 'cover', zIndex: 1,
    transition: 'all 240ms cubic-bezier(.2,.7,.2,1)'
  }
  // 兜底：强制结束（避免个别浏览器不触发过渡事件）
  setTimeout(() => { if (!finished) onCoverEnd() }, 280)
  await nextTick()
  // 目标：左侧封面最终位置
  const t = computeTargetRect()
  coverStyle.value = {
    ...coverStyle.value,
    left: t.left + 'px', top: t.top + 'px', width: t.width + 'px', height: t.height + 'px'
  }
})

function onCoverEnd(){
  if (finished) return
  finished = true
  emit('done')
  visible.value = false
}
</script>

<style scoped>
.np-overlay{ position: fixed; inset: 0; z-index: 10000; pointer-events: none; }
.np-cover{ box-shadow: 0 16px 40px rgba(0,0,0,.22); border-radius: 12px; }
</style>

