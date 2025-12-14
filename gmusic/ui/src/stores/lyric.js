import { defineStore } from 'pinia'

const KEY = 'gmusic:lyric-ui'

export const useLyricUiStore = defineStore('lyricUi', {
  state: () => ({
    anchor: 0.35,       // 0~1，当前行锚点，相对容器高度
    fontSize: 16,       // px，基础字号
    blurOthers: false,  // 是否模糊非当前行
    showGuide: true     // 是否显示锚点引导线
  }),
  actions: {
    load(){
      try{
        const raw = localStorage.getItem(KEY)
        if(!raw) return
        const data = JSON.parse(raw)
        if(data && typeof data==='object'){
          if(typeof data.anchor==='number') this.anchor = data.anchor
          if(typeof data.fontSize==='number') this.fontSize = data.fontSize
          if(typeof data.blurOthers==='boolean') this.blurOthers = data.blurOthers
          if(typeof data.showGuide==='boolean') this.showGuide = data.showGuide
        }
      }catch{}
    },
    save(){
      try{ localStorage.setItem(KEY, JSON.stringify({ anchor:this.anchor, fontSize:this.fontSize, blurOthers:this.blurOthers, showGuide:this.showGuide })) }catch{}
    }
  }
})
