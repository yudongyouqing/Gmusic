<template>
  <div class="app">
    <div class="container">
      <!-- 顶部大标题已移除，采用 Sidebar + TopBar 布局 -->
      <div class="main-content">
        <div class="main-shell">
          <Sidebar />
          <div class="content">
            <RouterView />
          </div>
        </div>
      </div>
    </div>

    <!-- 迷你播放器条，固定在页面底部 -->
    <MiniPlayerBar />
  </div>
</template>

<script setup>
import { RouterView } from 'vue-router'
import MiniPlayerBar from './components/MiniPlayerBar.vue'
import Sidebar from './components/Sidebar.vue'
</script>
