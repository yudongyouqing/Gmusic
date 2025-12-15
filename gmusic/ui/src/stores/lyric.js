import { defineStore } from 'pinia'

const KEY = 'gmusic:lyric-ui'

export const useLyricUiStore = defineStore('lyricUi', {
  state: () => ({
    fontSize: 20,       // px，基础字号 (已调大)
    fontWeight: 400,    // 100-900，字体粗细
    blurOthers: false,  // 是否模糊非当前行
    showTranslation: false, // 是否显示翻译
    translationScale: 80, // %，翻译字体缩放
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
          translationScale: this.translationScale
        }
        localStorage.setItem(KEY, JSON.stringify(dataToSave))
      }catch{}
    }
  }
})
