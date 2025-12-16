import { defineStore } from 'pinia'

const KEY = 'gmusic:settings'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    nowPlayingBackgroundUrl: '', // 播放页自定义背景图片 URL
  }),
  actions: {
    load(){
      try{
        const raw = localStorage.getItem(KEY)
        if(!raw) return
        const data = JSON.parse(raw)
        if(data && typeof data==='object'){
          if(typeof data.nowPlayingBackgroundUrl === 'string') {
            this.nowPlayingBackgroundUrl = data.nowPlayingBackgroundUrl
          }
        }
      }catch{}
    },
    save(){
      try{
        const dataToSave = {
          nowPlayingBackgroundUrl: this.nowPlayingBackgroundUrl
        }
        localStorage.setItem(KEY, JSON.stringify(dataToSave))
      }catch{}
    }
  }
})

