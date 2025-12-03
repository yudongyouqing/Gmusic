<template>
  <div class="app">
    <div class="container">
      <header class="header">
        <h1>GMusic 本地音乐播放器</h1>
        <nav class="nav">
          <RouterLink to="/" class="nav-link">音乐库</RouterLink>
          <span class="sep">|</span>
          <RouterLink to="/now-playing" class="nav-link">正在播放</RouterLink>
        </nav>
      </header>

      <div class="main-content">
        <RouterView />
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
