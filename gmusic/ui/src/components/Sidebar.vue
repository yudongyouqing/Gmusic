<template>
  <aside class="sidebar">
    <div class="brand">
      <div class="logo" aria-label="logo">
        <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
          <path d="M9 3v10.5a3.5 3.5 0 1 0 2 3.15V7h6V3H9z" fill="currentColor"/>
        </svg>
      </div>
      <div class="name">GMusic</div>
    </div>

    <nav class="nav">
      <RouterLink to="/" class="nav-item" active-class="active">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V9l-5-6z" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            <path d="M14 3v6h6" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
          </svg>
        </span>
        <span class="text">歌曲</span>
      </RouterLink>
      <RouterLink to="/now-playing" class="nav-item" active-class="active">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M8 5l12 7-12 7V5z" fill="currentColor" />
          </svg>
        </span>
        <span class="text">正在播放</span>
      </RouterLink>
      <RouterLink to="/queue" class="nav-item" active-class="active">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M4 7h14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <path d="M4 12h14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <path d="M4 17h10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </span>
        <span class="text">播放列表</span>
      </RouterLink>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M9 3v10.5a3.5 3.5 0 1 0 2 3.15V7h6V3H9z" fill="currentColor"/>
          </svg>
        </span>
        <span class="text">曲风</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <circle cx="12" cy="12" r="8" stroke="currentColor" stroke-width="2" fill="none"/>
            <circle cx="12" cy="12" r="2" fill="currentColor"/>
          </svg>
        </span>
        <span class="text">专辑</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <rect x="9" y="3" width="6" height="11" rx="3" fill="currentColor"/>
            <path d="M12 14v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <path d="M8 20h8" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </span>
        <span class="text">艺术家</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M3 7h6l2 2h10v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7z" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
          </svg>
        </span>
        <span class="text">文件夹</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M4 4v16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <path d="M10 4v16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <path d="M16 4v16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <circle cx="4" cy="10" r="2" fill="currentColor"/>
            <circle cx="10" cy="7" r="2" fill="currentColor"/>
            <circle cx="16" cy="14" r="2" fill="currentColor"/>
          </svg>
        </span>
        <span class="text">音质库</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M12 8a4 4 0 1 1 0 8 4 4 0 0 1 0-8z" fill="none" stroke="currentColor" stroke-width="2"/>
            <path d="M3 12h3M18 12h3M12 3v3M12 18v3M5.2 5.2l2.1 2.1M16.7 16.7l2.1 2.1M5.2 18.8l2.1-2.1M16.7 7.3l2.1-2.1" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </span>
        <span class="text">设置</span>
      </a>
      <div class="divider"></div>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M12 3l3 6 6 .9-4.5 4.3L18 21l-6-3.2L6 21l1.5-6.8L3 9.9 9 9l3-6z" fill="currentColor"/>
          </svg>
        </span>
        <span class="text">新建歌单</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M12 3c1.5 3-1 4.5-1 6.5A3.5 3.5 0 0 0 14.5 13c0 2.5-2 4.5-4.5 4.5S5.5 15.5 5.5 13c0-2.4 1.5-4.1 3.2-5.8C10.2 5.7 11.3 4.5 12 3z" fill="currentColor"/>
          </svg>
        </span>
        <span class="text">追梦歌单</span>
      </a>
      <a class="nav-item disabled" title="占位">
        <span class="icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path d="M12 21l-6-6a4 4 0 1 1 6-5l1-2 3 2-2 4 4 4-6 3z" fill="currentColor"/>
          </svg>
        </span>
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
