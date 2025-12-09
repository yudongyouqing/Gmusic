import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSongs, searchSongs, play, pause, resume, stop, setVolume, status, getLyrics, scan, seek } from '../api/music'

export const usePlayerStore = defineStore('player', () => {
  // state
  const songs = ref([])
  const searchResults = ref(null)
  const currentSong = ref(null)
  const isPlaying = ref(false)
  const lyrics = ref(null)
  const playerStatus = ref({ position: 0, duration: 0 })
  const playPending = ref(false)

  // 播放模式：loop（列表循环）/ shuffle（随机）
  const playMode = ref('loop')

  // getters
  const songList = () => searchResults.value || songs.value

  function getCurrentList() { return songList() || [] }
  function getCurrentIndex() {
    const list = getCurrentList()
    if (!currentSong.value) return -1
    return list.findIndex(s => s.id === currentSong.value.id)
  }

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
      const msg = e?.response?.data?.error || e?.message || '播放失败'
      alert(`播放失败：${msg}`)
      isPlaying.value = false
    } finally {
      playPending.value = false
    }
  }

  async function playByIndex(idx) {
    const list = getCurrentList()
    if (idx < 0 || idx >= list.length) return
    await playSong(list[idx])
  }

  async function nextSong() {
    const list = getCurrentList()
    if (!list.length) return
    const cur = getCurrentIndex()
    let nextIdx = 0
    if (playMode.value === 'shuffle') {
      if (list.length === 1) nextIdx = 0
      else {
        let r = Math.floor(Math.random() * list.length)
        if (r === cur) r = (r + 1) % list.length
        nextIdx = r
      }
    } else {
      nextIdx = (cur + 1 + list.length) % list.length
    }
    await playByIndex(nextIdx)
  }

  async function prevSong() {
    const list = getCurrentList()
    if (!list.length) return
    const cur = getCurrentIndex()
    let prevIdx = 0
    if (playMode.value === 'shuffle') {
      if (list.length === 1) prevIdx = 0
      else {
        let r = Math.floor(Math.random() * list.length)
        if (r === cur) r = (r + list.length - 1) % list.length
        prevIdx = r
      }
    } else {
      prevIdx = (cur - 1 + list.length) % list.length
    }
    await playByIndex(prevIdx)
  }

  async function pauseSong() { await pause(); isPlaying.value = false }
  async function resumeSong() { await resume(); isPlaying.value = true }

  async function stopSong() {
    await stop()
    isPlaying.value = false
    currentSong.value = null
    lyrics.value = null
    playerStatus.value = { position: 0, duration: 0 }
  }

  async function setVolumePercent(vol) { await setVolume(vol / 100) }

  async function refreshStatus() { const { data } = await status(); playerStatus.value = data }

  async function seekTo(sec) {
    if (sec < 0) sec = 0
    const d = playerStatus.value?.duration || 0
    if (d > 0 && sec > d) sec = d
    await seek(sec)
    await refreshStatus()
  }

  // 扫描目录导入歌曲
  async function scanDir(dirPath, workers = 4) { await scan(dirPath, workers); setTimeout(fetchSongs, 2000); setTimeout(fetchSongs, 5000) }

  function setPlayMode(mode) { if (mode !== 'loop' && mode !== 'shuffle') return; playMode.value = mode }

  return {
    // state
    songs, searchResults, currentSong, isPlaying, lyrics, playerStatus, playPending, playMode,
    // getters
    songList,
    // actions
    fetchSongs, doSearch, playSong, playByIndex, nextSong, prevSong, pauseSong, resumeSong, stopSong, setVolumePercent, refreshStatus, seekTo, scanDir, setPlayMode,
  }
})
