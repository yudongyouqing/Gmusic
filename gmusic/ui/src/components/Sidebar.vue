<template>
  <aside class="sidebar">
    <div class="brand">
      <div class="logo">🎵</div>
      <div class="name">GMusic</div>
    </div>

    <nav class="nav">
      <RouterLink to="/" class="nav-item" active-class="active">
        <span class="icon">📄</span>
        <span class="text">歌曲</span>
      </RouterLink>
      <RouterLink to="/now-playing" class="nav-item" active-class="active">
        <span class="icon">▶️</span>
        <span class="text">正在播放</span>
      </RouterLink>
      <a class="nav-item disabled" title="占位">
        <span class="icon">🎸</span>
        <span class="text">曲风</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon">💿</span>
        <span class="text">专辑</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon">🎤</span>
        <span class="text">艺术家</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon">📁</span>
        <span class="text">文件夹</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon">🎚️</span>
        <span class="text">音质库</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon">⚙️</span>
        <span class="text">设置</span>
      </a>
      <div class="divider"></div>
      <a class="nav-item disabled" title="占位">
        <span class="icon">⭐</span>
        <span class="text">新建歌单</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon">🔥</span>
        <span class="text">追梦歌单</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon">💔</span>
        <span class="text">苦情歌单</span>
      </a>
    </nav>
  </aside>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.sidebar {
  width: 220px;
  flex: 0 0 220px;
  height: 100%;
  background: var(--mica-surface, rgba(255,255,255,.45));
  backdrop-filter: blur(var(--mica-blur, 16px)) saturate(var(--mica-saturate,160%));
  -webkit-backdrop-filter: blur(var(--mica-blur, 16px)) saturate(var(--mica-saturate,160%));
  border-right: 1px solid var(--mica-border, rgba(255,255,255,.35));
  padding: 12px;
  box-sizing: border-box;
}
.brand { display:flex; align-items:center; gap: 10px; padding: 8px 10px; margin-bottom: 8px; }
.logo { width: 28px; height: 28px; border-radius: 6px; background: rgba(0,0,0,0.08); display:flex; align-items:center; justify-content:center; }
.name { font-weight: 700; color:#222; }
.nav { display:flex; flex-direction:column; gap: 4px; padding: 6px 0; }
.nav-item { display:flex; align-items:center; gap: 10px; padding: 10px 12px; border-radius: 8px; color:#333; text-decoration:none; }
.nav-item:hover { background: rgba(255,255,255,.55); }
.nav-item.active { background: rgba(255,255,255,.7); font-weight:600; }
.nav-item.disabled { opacity:.65; cursor: default; }
.icon { width: 18px; text-align:center; }
.divider { height: 1px; background: rgba(0,0,0,0.08); margin: 8px 4px; }
</style>
