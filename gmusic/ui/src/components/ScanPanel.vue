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
.scan-panel { 
  width: 100%;
  box-sizing: border-box;
  background: transparent; 
  border-radius: 0; 
  padding: 0; 
  box-shadow: none; 
}
.row { 
  display: flex; 
  gap: 8px; 
  align-items: center; 
  flex-wrap: wrap; /* 小屏换行，避免溢出 */
}
.scan-input { 
  flex: 1 1 auto; 
  min-width: 0;            /* 允许在 flex 下收缩，避免水平溢出 */
  padding: 10px 12px; 
  border: 1px solid rgba(0,0,0,0.08); 
  border-radius: 8px; 
  background: rgba(255,255,255,0.7);
  box-sizing: border-box;  /* 包含 padding 与边框 */
  font-size: 14px;
  color: #222;
}

.scan-input::placeholder {
  color: #999;
}

.scan-input:focus {
  background: rgba(255,255,255,0.9);
  border-color: rgba(0,0,0,0.15);
  outline: none;
}

.scan-workers { 
  width: 60px; 
  padding: 8px; 
  border: 1px solid rgba(0,0,0,0.08); 
  border-radius: 8px; 
  background: rgba(255,255,255,0.7);
  box-sizing: border-box;
  font-size: 14px;
  color: #222;
}

.scan-workers:focus {
  background: rgba(255,255,255,0.9);
  border-color: rgba(0,0,0,0.15);
  outline: none;
}

.scan-btn { 
  padding: 8px 16px; 
  border: none; 
  border-radius: 8px; 
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
  color: #fff; 
  cursor: pointer; 
  font-size: 14px;
  white-space: nowrap;
  transition: all 0.2s;
}

.scan-btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.scan-btn:disabled { 
  opacity: 0.6; 
  cursor: not-allowed; 
}

.tips { 
  margin-top: 6px; 
  color: #666; 
  font-size: 12px; 
  display: none;
}

.msg { 
  color: #333; 
  margin-top: 4px; 
}

@media (max-width: 1024px) {
  .scan-workers { 
    width: 50px; 
  }
}

@media (max-width: 768px) {
  .scan-workers, .scan-btn { height: 40px; }
}
</style>
