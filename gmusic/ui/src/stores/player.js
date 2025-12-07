import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSongs, searchSongs, play, pause, resume, stop, setVolume, status, getLyrics, scan } from '../api/music'

export const usePlayerStore = defineStore('player', () => {
  // state
  const songs = ref([])
  const searchResults = ref(null)
  const currentSong = ref(null)
  const isPlaying = ref(false)
  const lyrics = ref(null)
  const playerStatus = ref({ position: 0, duration: 0 })
  const playPending = ref(false)

  // getters
  const songList = () => searchResults.value || songs.value

  // actions
  async function fetchSongs() {
    const { data } = await getSongs()
    songs.value = data.songs || []
  }

  async function doSearch(keyword) {
    if (!keyword) {
      searchResults.value = null
      return
    }
    const { data } = await searchSongs(keyword)
    searchResults.value = data.songs || []
  }

  async function playSong(song) {
    if (playPending.value) return
    playPending.value = true
    try {
      currentSong.value = song
      await play(song.file_path)
      isPlaying.value = true
      try {
        const { data } = await getLyrics(song.id)
        lyrics.value = data
      } catch {
        lyrics.value = null
      }
    } catch (e) {
      // 前端可视化错误提示
      const msg = e?.response?.data?.error || e?.message || '播放失败'
      // eslint-disable-next-line no-alert
      alert(`播放失败：${msg}`)
      isPlaying.value = false
    } finally {
      playPending.value = false
    }
  }

  async function pauseSong() {
    await pause()
    isPlaying.value = false
  }

  async function resumeSong() {
    await resume()
    isPlaying.value = true
  }

  async function stopSong() {
    await stop()
    isPlaying.value = false
    currentSong.value = null
    lyrics.value = null
    playerStatus.value = { position: 0, duration: 0 }
  }

  async function setVolumePercent(vol) {
    await setVolume(vol / 100)
  }

  async function refreshStatus() {
    const { data } = await status()
    playerStatus.value = data
  }

  // 扫描目录导入歌曲
  async function scanDir(dirPath, workers = 4) {
    await scan(dirPath, workers)
    // 简单轮询刷新列表（2s、5s）
    setTimeout(fetchSongs, 2000)
    setTimeout(fetchSongs, 5000)
  }

  return {
    // state
    songs, searchResults, currentSong, isPlaying, lyrics, playerStatus, playPending,
    // getters
    songList,
    // actions
    fetchSongs, doSearch, playSong, pauseSong, resumeSong, stopSong, setVolumePercent, refreshStatus, scanDir,
  }
})
