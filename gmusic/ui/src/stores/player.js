import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSongs, searchSongs, play, pause, resume, stop, setVolume, status, getLyrics } from '../api/music'

export const usePlayerStore = defineStore('player', () => {
  // state
  const songs = ref([])
  const searchResults = ref(null)
  const currentSong = ref(null)
  const isPlaying = ref(false)
  const lyrics = ref(null)
  const playerStatus = ref({ position: 0, duration: 0 })

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
    currentSong.value = song
    await play(song.file_path)
    isPlaying.value = true
    try {
      const { data } = await getLyrics(song.id)
      lyrics.value = data
    } catch {
      lyrics.value = null
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

  return {
    // state
    songs, searchResults, currentSong, isPlaying, lyrics, playerStatus,
    // getters
    songList,
    // actions
    fetchSongs, doSearch, playSong, pauseSong, resumeSong, stopSong, setVolumePercent, refreshStatus,
  }
})

