import { defineStore } from 'pinia'

const KEY = 'gmusic:lyric-ui'

export const useLyricUiStore = defineStore('lyricUi', {
  state: () => ({
    fontSize: 20,
    fontWeight: 400,
    blurOthers: false,
    showTranslation: false,
    translationScale: 80,
    backgroundBlur: 22, // 新增：播放页背景模糊
  }),
  actions: {
    load(){
      try{
        const raw = localStorage.getItem(KEY)
        if(!raw) return
        const data = JSON.parse(raw)
        if(data && typeof data==='object'){
          if(typeof data.fontSize==='number') this.fontSize = data.fontSize
          if(typeof data.fontWeight==='number') this.fontWeight = data.fontWeight
          if(typeof data.blurOthers==='boolean') this.blurOthers = data.blurOthers
          if(typeof data.showTranslation==='boolean') this.showTranslation = data.showTranslation
          if(typeof data.translationScale==='number') this.translationScale = data.translationScale
          if(typeof data.backgroundBlur==='number') this.backgroundBlur = data.backgroundBlur
        }
      }catch{}
    },
    save(){
      try{
        const dataToSave = {
          fontSize: this.fontSize,
          fontWeight: this.fontWeight,
          blurOthers: this.blurOthers,
          showTranslation: this.showTranslation,
          translationScale: this.translationScale,
          backgroundBlur: this.backgroundBlur
        }
        localStorage.setItem(KEY, JSON.stringify(dataToSave))
      }catch{}
    }
  }
})
