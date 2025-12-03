import http from '../service/http'

export const getSongs = () => http.get('/songs')
export const getSongById = (id) => http.get(`/songs/${id}`)
export const searchSongs = (q) => http.get('/songs/search', { params: { q } })

export const play = (file_path) => http.post('/player/play', { file_path })
export const pause = () => http.post('/player/pause')
export const resume = () => http.post('/player/resume')
export const stop = () => http.post('/player/stop')
export const setVolume = (volume) => http.post('/player/volume', { volume })
export const status = () => http.get('/player/status')

export const getLyrics = (songID) => http.get(`/lyrics/${songID}`)
export const scan = (dir_path, workers = 4) => http.post('/scan', { dir_path, workers })

