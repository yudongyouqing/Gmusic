<template>
  <div class="theme-pop" :style="popStyle" @click.stop>
    <div class="row">
      <label>主题</label>
      <div class="chips">
        <button class="chip" :class="{active: local.theme==='current'}" @click="local.theme='current'; apply()">当前风格</button>
        <button class="chip" :class="{active: local.theme==='glass'}" @click="local.theme='glass'; apply()">毛玻璃</button>
      </div>
    </div>

    <div class="row">
      <label>透明度 {{ local.alpha.toFixed(2) }}</label>
      <input type="range" min="0.20" max="0.90" step="0.01" v-model.number="local.alpha" @input="apply" />
    </div>

    <div class="row">
      <label>模糊 {{ local.blur }}px</label>
      <input type="range" min="6" max="40" step="1" v-model.number="local.blur" @input="apply" />
    </div>

    <div class="row">
      <label>饱和度 {{ local.saturate }}%</label>
      <input type="range" min="100" max="220" step="5" v-model.number="local.saturate" @input="apply" />
    </div>

    <div class="row">
      <label>歌词字号 {{ local.lyricFontSize }}px</label>
      <input type="range" min="12" max="18" step="1" v-model.number="local.lyricFontSize" @input="apply" />
    </div>

    <div class="row">
      <label>底栏歌词水平偏移 {{ local.tickerOffsetX }}px</label>
      <input type="range" min="-400" max="400" step="1" v-model.number="local.tickerOffsetX" @input="apply" />
    </div>

    <div class="actions">
      <button class="btn" @click="reset">重置</button>
      <button class="btn primary" @click="save">保存</button>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed } from 'vue'
import { useUiStore } from '../stores/ui'

const emit = defineEmits(['close'])

const props = defineProps({
  pos: { type: Object, default: null },
  z: { type: Number, default: 5000 },
})

const ui = useUiStore()
const local = reactive({ theme: ui.theme, alpha: ui.alpha, blur: ui.blur, saturate: ui.saturate, lyricFontSize: ui.lyricFontSize, tickerOffsetX: ui.tickerOffsetX })

function apply(){
  ui.theme = local.theme
  ui.alpha = local.alpha
  ui.blur = local.blur
  ui.saturate = local.saturate
  ui.lyricFontSize = local.lyricFontSize
  ui.tickerOffsetX = local.tickerOffsetX
  ui.applyTheme()
}

function save(){
  apply()
  ui.saveTheme()
  emit('close') // 通知父组件关闭
}

function reset(){
  if(local.theme==='glass'){ local.alpha=0.36; local.blur=22; local.saturate=180 } 
  else { local.alpha=0.45; local.blur=18; local.saturate=160 }
  local.lyricFontSize = 14
  local.tickerOffsetX = 0 // 重置底栏歌词水平偏移
  apply()
}

const popStyle = computed(()=>{
  if (!props.pos) return { position:'fixed', right:'24px', top:'72px', zIndex: props.z }
  const { top, left, right, height } = props.pos
  const width = 280
  const x = Math.max(12, Math.min(right - width, left))
  const y = top + height + 8
  return { position:'fixed', left: x+'px', top: y+'px', zIndex: props.z }
})
</script>

<style scoped>
.theme-pop{
  width: 280px;
  padding: 14px;
  border-radius: 12px;
  background: var(--mica-surface);
  backdrop-filter: blur(var(--mica-blur)) saturate(var(--mica-saturate));
  -webkit-backdrop-filter: blur(var(--mica-blur)) saturate(var(--mica-saturate));
  border: 1px solid var(--mica-border);
  box-shadow: 0 20px 40px rgba(0,0,0,0.25);
}
.row{ margin-bottom: 10px; }
.row:last-of-type{ margin-bottom: 0; }
label{ display:block; font-size:12px; color:#555; margin-bottom:6px; }
.chips{ display:flex; gap:8px; }
.chip{ padding: 6px 10px; border-radius: 999px; border: 1px solid rgba(0,0,0,.08); background: rgba(255,255,255,.7); cursor:pointer; }
.chip.active{ background: linear-gradient(135deg,#667eea22,#764ba222); border-color: rgba(0,0,0,.15); }
.actions{ display:flex; justify-content:flex-end; gap:8px; margin-top:12px; }
.btn{ padding:8px 12px; border-radius:8px; border:1px solid rgba(0,0,0,.08); background: rgba(255,255,255,.8); cursor:pointer; }
.btn.primary{ background: linear-gradient(135deg,#667eea,#764ba2); color:#fff; border-color: transparent; }
input[type="range"]{ width:100%; accent-color:#667eea; }
</style>
