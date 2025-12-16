<template>
  <div class="page page-settings">
    <div class="card card--compact">
      <TopBar title="设置" />
    </div>

    <div class="card card--fill" style="margin-top:12px">
      <div class="setting-item">
        <label for="bg-url">播放页背景图片 URL</label>
        <input id="bg-url" type="text" v-model="settings.nowPlayingBackgroundUrl" @change="settings.save()" placeholder="输入图片的 URL，留空则使用封面" />
        <p class="tip">输入你想要的图片 URL 作为播放页的背景。清除并保存后，将恢复为使用当前歌曲的封面作为背景。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useSettingsStore } from '../stores/settings'
import TopBar from '../components/TopBar.vue'

const settings = useSettingsStore()
settings.load()

</script>

<style scoped>
.page-settings { display: flex; flex-direction: column; height: 100%; min-height: 0; }
.setting-item {
  margin-bottom: 24px;
}
label {
  display: block;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}
input[type="text"] {
  width: 100%;
  max-width: 600px;
  padding: 10px 12px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 14px;
}
.tip {
  font-size: 12px;
  color: #666;
  margin-top: 8px;
  max-width: 600px;
}
</style>

