<template>
  <div class="search-bar">
    <input
      type="text"
      class="search-input"
      :value="keyword"
      placeholder="搜索歌曲、歌手、专辑..."
      @input="onInput"
    />
    <button v-if="keyword" class="clear-btn" @click="clear">✕</button>
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

