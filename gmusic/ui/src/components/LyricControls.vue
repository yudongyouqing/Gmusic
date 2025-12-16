<template>
  <div class="lyric-ctrl" :style="posStyle" @click.stop>
    <div class="row head">
      <div class="title">播放页面歌词</div>
      <button class="close" @click="$emit('close')">×</button>
    </div>

    <div class="row">
      <label>背景模糊 {{ local.backgroundBlur }}px</label>
      <input type="range" min="0" max="50" step="1" v-model.number="local.backgroundBlur" @input="apply" />
    </div>

    <div class="row">
      <label>字体大小 {{ local.fontSize }}px</label>
      <input type="range" min="12" max="32" step="1" v-model.number="local.fontSize" @input="apply" />
    </div>

    <div class="row">
      <label>字体粗细 {{ local.fontWeight }}</label>
      <input type="range" min="100" max="900" step="100" v-model.number="local.fontWeight" @input="apply" />
    </div>

    <div class="row toggle">
      <label>模糊非当前歌词</label>
      <label class="switch">
        <input type="checkbox" v-model="local.blurOthers" @change="apply" />
        <span class="slider"></span>
      </label>
    </div>

    <div class="row toggle">
      <label>歌词翻译 (占位)</label>
      <label class="switch">
        <input type="checkbox" v-model="local.showTranslation" @change="apply" />
        <span class="slider"></span>
      </label>
    </div>

    <div class="row">
      <label>歌词翻译字体大小缩放 {{ local.translationScale }}%</label>
      <input type="range" min="50" max="100" step="5" v-model.number="local.translationScale" @input="apply" />
    </div>

    <div class="actions">
      <button class="btn" @click="reset">重置</button>
      <button class="btn primary" @click="save">保存</button>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed } from 'vue'
import { useLyricUiStore } from '../stores/lyric'

const props = defineProps({
  right: { type: Number, default: 16 },
  top: { type: Number, default: 12 }
})

const emit = defineEmits(['close'])
const store = useLyricUiStore()
store.load()

const local = reactive({
  fontSize: store.fontSize,
  fontWeight: store.fontWeight,
  blurOthers: store.blurOthers,
  showTranslation: store.showTranslation,
  translationScale: store.translationScale,
  backgroundBlur: store.backgroundBlur
})

function apply(){
  store.fontSize = local.fontSize
  store.fontWeight = local.fontWeight
  store.blurOthers = local.blurOthers
  store.showTranslation = local.showTranslation
  store.translationScale = local.translationScale
  store.backgroundBlur = local.backgroundBlur
}

function save(){ 
  apply()
  store.save()
  emit('close')
}

function reset(){
  local.fontSize = 20
  local.fontWeight = 400
  local.blurOthers = false
  local.showTranslation = false
  local.translationScale = 80
  local.backgroundBlur = 22
  apply()
}

const posStyle = computed(()=>({ position:'fixed', right: props.right+'px', top: props.top+'px', zIndex: 20001 }))
</script>

<style scoped>
.lyric-ctrl{ width: 320px; padding: 14px; border-radius: 12px; background: var(--mica-surface, rgba(255,255,255,.7)); backdrop-filter: blur(var(--mica-blur,18px)) saturate(var(--mica-saturate,160%)); border:1px solid var(--mica-border, rgba(255,255,255,.35)); box-shadow: 0 10px 24px rgba(0,0,0,.18); }
.row{ margin-bottom: 12px; }
.row:last-of-type{ margin-bottom: 0; }
.head{ display:flex; align-items:center; justify-content:space-between; margin-bottom: 10px; }
.title{ font-weight:700; color:#222; }
.close{ width:28px; height:28px; border-radius:8px; border:1px solid rgba(0,0,0,0.06); background: rgba(255,255,255,.9); cursor:pointer; }
label{ display:block; font-size:14px; color:#333; margin-bottom:8px; }
input[type="range"]{ width:100%; accent-color:#667eea; }
.toggle{ display:flex; align-items:center; justify-content:space-between; }
.switch{ position:relative; width:46px; height:24px; }
.switch input{ display:none; }
.slider{ position:absolute; inset:0; background:#d0d0d0; border-radius:999px; transition:.2s; }
.slider::after{ content:''; position:absolute; left:3px; top:3px; width:18px; height:18px; background:#fff; border-radius:50%; transition:.2s; box-shadow: 0 2px 6px rgba(0,0,0,.2); }
.switch input:checked + .slider{ background:#5b7bfe; }
.switch input:checked + .slider::after{ transform: translateX(22px); }
.actions{ display:flex; justify-content:flex-end; gap:8px; margin-top:16px; }
.btn{ padding:8px 14px; border-radius:8px; border:1px solid rgba(0,0,0,.08); background: rgba(255,255,255,.9); cursor:pointer; }
.btn.primary{ background: linear-gradient(135deg,#667eea,#764ba2); color:#fff; border-color: transparent; }
</style>
