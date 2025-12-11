import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
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

  // 播放模式：loop（列表循环）/ shuffle（随机）/ single（单曲循环）
  const playMode = ref('loop')

  // 排序 & 自定义顺序
  const sortMode = ref('title') // 'title' | 'artist' | 'album' | 'custom'
  const sortDir = ref('asc')    // 'asc' | 'desc'
  const customOrder = ref([])   // 保存自定义顺序的 id 列表

  // 随机播放队列（独立于排序）
  const queue = ref([])         // id 列表，表示实际播放顺序（仅 shuffle 使用）
  const queueIndex = ref(-1)    // 指向 queue 中当前歌曲位置

  // 记录已探测过的歌曲，避免重复请求
  const probed = ref(new Set())

  // 基础列表：搜索结果优先，否则全量
  const baseList = computed(() => searchResults.value || songs.value)

  // 应用排序后的列表（不含随机）
  const orderedList = computed(() => {
    const list = (baseList.value || []).slice()
    if (!list.length) return []
    if (sortMode.value === 'custom' && customOrder.value?.length) {
      const idxMap = new Map(customOrder.value.map((id, i) => [id, i]))
      list.sort((a, b) => {
        const ia = idxMap.has(a.id) ? idxMap.get(a.id) : Number.MAX_SAFE_INTEGER
        const ib = idxMap.has(b.id) ? idxMap.get(b.id) : Number.MAX_SAFE_INTEGER
        return ia - ib
      })
      return list
    }
    const key = sortMode.value
    const dir = sortDir.value === 'desc' ? -1 : 1
    list.sort((a, b) => {
      const av = (a?.[key] || '').toString()
      const bv = (b?.[key] || '').toString()
      return av.localeCompare(bv, 'zh-CN') * dir
    })
    return list
  })

  // 旧接口兼容：当前用于渲染与顺序计算的列表
  const songList = () => orderedList.value

  function getCurrentIndexIn(list) {
    if (!currentSong.value) return -1
    return (list || []).findIndex(s => s.id === currentSong.value.id)
  }

  function persistOrder() {
    try {
      localStorage.setItem('gmusic:order', JSON.stringify({ sortMode: sortMode.value, sortDir: sortDir.value, customOrder: customOrder.value }))
    } catch {}
  }
  function loadOrder() {
    try {
      const raw = localStorage.getItem('gmusic:order')
      if (!raw) return
      const data = JSON.parse(raw)
      if (data && typeof data === 'object') {
        if (data.sortMode) sortMode.value = data.sortMode
        if (data.sortDir) sortDir.value = data.sortDir
        if (Array.isArray(data.customOrder)) customOrder.value = data.customOrder
      }
    } catch {}
  }

  function ensureCustomOrderCoversAll() {
    const ids = new Set(customOrder.value)
    const list = songs.value || []
    const missing = list.filter(s => !ids.has(s.id)).map(s => s.id)
    if (missing.length) customOrder.value = customOrder.value.concat(missing)
  }

  function setSort(mode, dir = sortDir.value) {
    const modes = ['title', 'artist', 'album', 'custom']
    if (!modes.includes(mode)) return
    sortMode.value = mode
    sortDir.value = dir
    if (mode === 'custom' && customOrder.value.length === 0) ensureCustomOrderCoversAll()
    persistOrder()
  }
  function setCustomOrder(ids) {
    customOrder.value = Array.isArray(ids) ? ids.slice() : []
    persistOrder()
  }

  // 队列：在 shuffle 模式下生成一次，按队列切歌；其他模式直接用 orderedList
  function buildQueueFrom(list, startId) {
    const ids = (list || []).map(s => s.id)
    // Fisher-Yates 洗牌
    for (let i = ids.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1))
      ;[ids[i], ids[j]] = [ids[j], ids[i]]
    }
    // 将当前歌曲放到队首并设置 index=0
    if (startId) {
      const p = ids.indexOf(startId)
      if (p > 0) { ids.splice(p, 1); ids.unshift(startId) }
    }
    queue.value = ids
    queueIndex.value = 0
  }

  function getQueueListForView() {
    if (playMode.value === 'shuffle' && queue.value.length) {
      const idSet = new Set(queue.value)
      const byId = new Map((songs.value || []).map(s => [s.id, s]))
      const arr = queue.value.map(id => byId.get(id)).filter(Boolean)
      // 如果有搜索结果在影响 baseList，视图仍展示实际队列（与需求一致）
      return arr
    }
    return orderedList.value
  }

  // actions
  async function fetchSongs() {
    const { data } = await getSongs()
    songs.value = data.songs || []
    if (sortMode.value === 'custom') ensureCustomOrderCoversAll()
    setTimeout(() => updateMissingDurations().catch(() => {}), 0)
  }

  async function updateMissingDurations() {
    const list = baseList.value || []
    const targets = list.filter(s => (!s.duration || s.duration <= 0) && !probed.value.has(s.id))
    if (!targets.length) return

    const limit = 3
    let idx = 0
    async function worker() {
      while (idx < targets.length) {
        const cur = targets[idx++]
        probed.value.add(cur.id)
        try {
          const { data: info } = await audioInfoById(cur.id)
          if (info?.duration > 0) {
            const apply = (arr) => {
              if (!arr) return
              const i = arr.findIndex(x => x.id === cur.id)
              if (i >= 0) arr[i] = { ...arr[i], duration: info.duration }
            }
            apply(songs.value)
            apply(searchResults.value)
            if (currentSong.value && currentSong.value.id === cur.id) {
              currentSong.value = { ...currentSong.value, duration: info.duration }
              playerStatus.value = { ...playerStatus.value, duration: info.duration }
            }
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
    setTimeout(() => updateMissingDurations().catch(() => {}), 0)
  }

  async function playSong(song, opts = {}) {
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
      if (!song.duration || song.duration <= 0) {
        try {
          const { data: info } = await audioInfoById(song.id)
          if (info?.duration > 0) {
            currentSong.value = { ...currentSong.value, duration: info.duration }
            playerStatus.value = { ...playerStatus.value, duration: info.duration }
            const list = orderedList.value || []
            const idx = list.findIndex(s => s.id === song.id)
            if (idx >= 0) list[idx] = { ...list[idx], duration: info.duration }
            await updateSong(song.id, { duration: info.duration })
          }
        } catch (_) { /* 忽略探测失败 */ }
      }

      // 在 shuffle 模式下，保证队列存在并以当前歌曲为队首；可选择保留现有队列
      if (playMode.value === 'shuffle' && !opts.keepQueue) {
        buildQueueFrom(orderedList.value, song.id)
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
    const list = orderedList.value
    if (idx < 0 || idx >= list.length) return
    await playSong(list[idx])
  }

  function getSeqNextPrev(delta) {
    const list = orderedList.value
    if (!list.length) return -1
    const cur = getCurrentIndexIn(list)
    if (cur < 0) return -1
    const next = (cur + delta + list.length) % list.length
    return next
  }

  async function nextSong() {
    if (playMode.value === 'single') {
      const idx = getCurrentIndexIn(orderedList.value)
      if (idx >= 0) await playByIndex(idx)
      return
    }
    if (playMode.value === 'shuffle') {
      if (!queue.value.length) buildQueueFrom(orderedList.value, currentSong.value?.id)
      queueIndex.value = Math.min(queueIndex.value + 1, queue.value.length - 1)
      const nextId = queue.value[queueIndex.value]
      const song = (songs.value || []).find(s => s.id === nextId)
      if (song) await playSong(song)
      return
    }
    const nextIdx = getSeqNextPrev(1)
    if (nextIdx >= 0) await playByIndex(nextIdx)
  }

  async function prevSong() {
    if (playMode.value === 'single') {
      const idx = getCurrentIndexIn(orderedList.value)
      if (idx >= 0) await playByIndex(idx)
      return
    }
    if (playMode.value === 'shuffle') {
      if (!queue.value.length) buildQueueFrom(orderedList.value, currentSong.value?.id)
      queueIndex.value = Math.max(queueIndex.value - 1, 0)
      const prevId = queue.value[queueIndex.value]
      const song = (songs.value || []).find(s => s.id === prevId)
      if (song) await playSong(song)
      return
    }
    const prevIdx = getSeqNextPrev(-1)
    if (prevIdx >= 0) await playByIndex(prevIdx)
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

  async function refreshStatus() {
    try {
      const { data } = await status()
      playerStatus.value = data
      if (import.meta.env && import.meta.env.DEV) console.log('[player/status]', data)
      // 播放结束自动下一首
      if (isPlaying.value && data.duration > 0 && data.position >= data.duration - 0.5) {
        await nextSong()
      }
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

  async function scanDir(dirPath, workers = 4) { await scan(dirPath, workers); setTimeout(fetchSongs, 2000); setTimeout(fetchSongs, 5000) }

  function setPlayMode(mode) {
    const modes = ['loop', 'shuffle', 'single']
    if (!modes.includes(mode)) return
    const prev = playMode.value
    playMode.value = mode
    if (mode === 'shuffle' && prev !== 'shuffle') {
      buildQueueFrom(orderedList.value, currentSong.value?.id)
    }
  }

  function togglePlayMode() {
    const modes = ['loop', 'shuffle', 'single']
    const idx = modes.indexOf(playMode.value)
    setPlayMode(modes[(idx + 1) % modes.length])
  }

  // 初始化加载排序配置
  loadOrder()

  return {
    songs, searchResults, currentSong, isPlaying, lyrics, playerStatus, playPending, playMode,
    sortMode, sortDir, customOrder, queue, queueIndex,
    songList, orderedList, getQueueListForView,
    fetchSongs, updateMissingDurations, doSearch, playSong, playByIndex, nextSong, prevSong, pauseSong, resumeSong, stopSong, setVolumePercent, refreshStatus, seekTo, scanDir, setPlayMode, togglePlayMode,
    setSort, setCustomOrder,
  }
})
