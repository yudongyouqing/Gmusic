import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSongs, searchSongs, play, pause, resume, stop, setVolume, status, getLyrics, scan, seek, audioInfoById, updateSong } from '../api/music'

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

  // 记录已探测过的歌曲，避免重复请求
  const probed = ref(new Set())

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
    // 后台补全缺失时长（非阻塞）
    setTimeout(() => updateMissingDurations().catch(() => {}), 0)
  }

  async function updateMissingDurations() {
    const list = songList() || []
    // 只处理 duration<=0 且未探测过的
    const targets = list.filter(s => (!s.duration || s.duration <= 0) && !probed.value.has(s.id))
    if (!targets.length) return

    const limit = 3 // 并发限制
    let idx = 0
    async function worker() {
      while (idx < targets.length) {
        const cur = targets[idx++]
        probed.value.add(cur.id)
        try {
          const { data: info } = await audioInfoById(cur.id)
          if (info?.duration > 0) {
            // 更新 songs/searchResults 列表里的该条
            const apply = (arr) => {
              if (!arr) return
              const i = arr.findIndex(x => x.id === cur.id)
              if (i >= 0) arr[i] = { ...arr[i], duration: info.duration }
            }
            apply(songs.value)
            apply(searchResults.value)
            // 如当前播放正是该条，也更新当前状态
            if (currentSong.value && currentSong.value.id === cur.id) {
              currentSong.value = { ...currentSong.value, duration: info.duration }
              playerStatus.value = { ...playerStatus.value, duration: info.duration }
            }
            // 持久化到 DB
            await updateSong(cur.id, { duration: info.duration })
          }
        } catch (_) { /* 忽略单条错误 */ }
      }
    }
    const jobs = Array.from({ length: Math.min(limit, targets.length) }, () => worker())
    await Promise.all(jobs)
  }

  async function doSearch(keyword) {
    if (!keyword) {
      searchResults.value = null
      return
    }
    const { data } = await searchSongs(keyword)
    searchResults.value = data.songs || []
    // 搜索结果页也尝试后台补全
    setTimeout(() => updateMissingDurations().catch(() => {}), 0)
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
      // 播放后若时长为 0，则探测并缓存（本地 + 数据库）
      if (!song.duration || song.duration <= 0) {
        try {
          const { data: info } = await audioInfoById(song.id)
          if (info?.duration > 0) {
            // 更新当前播放与列表缓存
            currentSong.value = { ...currentSong.value, duration: info.duration }
            playerStatus.value = { ...playerStatus.value, duration: info.duration }
            const list = songList() || []
            const idx = list.findIndex(s => s.id === song.id)
            if (idx >= 0) list[idx] = { ...list[idx], duration: info.duration }
            // 持久化到数据库
            await updateSong(song.id, { duration: info.duration })
          }
        } catch (_) { /* 忽略探测失败 */ }
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

  // 加可视化日志，便于排查“没看到状态”问题
  async function refreshStatus() {
    try {
      const { data } = await status()
      playerStatus.value = data
      if (import.meta.env && import.meta.env.DEV) console.log('[player/status]', data)
    } catch (e) {
      console.error('[player/status] error', e)
    }
  }

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
    fetchSongs, updateMissingDurations, doSearch, playSong, playByIndex, nextSong, prevSong, pauseSong, resumeSong, stopSong, setVolumePercent, refreshStatus, seekTo, scanDir, setPlayMode,
  }
})
