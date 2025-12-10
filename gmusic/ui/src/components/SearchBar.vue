<template>
  <div class="search-bar">
    <input
      type="text"
      class="search-input"
      :value="keyword"
      placeholder="搜索歌曲、歌手、专辑..."
      @input="onInput"
    />
    <button v-if="keyword" class="clear-btn" @click="clear" aria-label="清空">
      <svg viewBox="0 0 24 24" width="16" height="16" aria-hidden="true">
        <path d="M6 6l12 12M18 6L6 18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import './SearchBar.css'

const emit = defineEmits(['search'])
const keyword = ref('')

function onInput(e) {
  keyword.value = e.target.value
  emit('search', keyword.value)
}

function clear() {
  keyword.value = ''
  emit('search', '')
}
</script>

