<template>
  <div class="scan-panel">
    <div class="row">
      <input
        v-model="dirPath"
        class="scan-input"
        type="text"
        placeholder="输入要扫描的目录，例如 D:/Music 或 F:/无损音乐"
      />
      <input v-model.number="workers" class="scan-workers" type="number" min="1" max="16" />
      <button class="scan-btn" :disabled="loading || !dirPath" @click="startScan">{{ loading ? '扫描中...' : '开始扫描' }}</button>
    </div>
    <div class="tips">
      <p>提示：Windows 路径请使用正斜杠，例如 D:/Music；或在界面中直接粘贴路径后替换 \ 为 /。</p>
      <p v-if="message" class="msg">{{ message }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { usePlayerStore } from '../stores/player'

const store = usePlayerStore()

const dirPath = ref('D:/Music')
const workers = ref(4)
const loading = ref(false)
const message = ref('')

async function startScan() {
  if (!dirPath.value) return
  loading.value = true
  message.value = ''
  try {
    await store.scanDir(dirPath.value, workers.value || 4)
    message.value = '扫描任务已启动，列表会自动刷新。'
  } catch (e) {
    message.value = '扫描失败，请检查目录是否存在以及权限。'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.scan-panel { background: #fff; border-radius: 12px; padding: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.row { display: flex; gap: 8px; align-items: center; }
.scan-input { flex: 1; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px; }
.scan-workers { width: 80px; padding: 8px; border: 1px solid #ddd; border-radius: 8px; }
.scan-btn { padding: 10px 16px; border: none; border-radius: 8px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #fff; cursor: pointer; }
.scan-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.tips { margin-top: 6px; color: #666; font-size: 12px; }
.msg { color: #333; margin-top: 4px; }
</style>

