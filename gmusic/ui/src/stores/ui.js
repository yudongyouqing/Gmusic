import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', {
  state: () => ({
    theme: 'current', // 'current' | 'glass'
    alpha: 0.45,      // 透明度 0-1
    blur: 18,         // 模糊 px
    saturate: 160,    // 饱和 %
    lyricFontSize: 14, // 歌词字号 px
    lyricAnchor: 0.35, // 歌词锚点（0~1，越小越靠上）
    tickerOffsetX: 0,  // 底栏歌词水平偏移(px，负为左移，正为右移)
    overlayRect: null,
    overlayCover: ''
  }),
  actions: {
    applyTheme() {
      const presets = {
        current: { alpha: this.alpha, blur: this.blur, saturate: this.saturate },
        glass:   { alpha: Math.min(0.36, this.alpha), blur: Math.max(22, this.blur), saturate: Math.max(180, this.saturate) },
      }
      const p = presets[this.theme] || presets.current
      const surface = `rgba(255,255,255,${p.alpha})`
      const surfaceStrong = `rgba(255,255,255,${Math.min(0.9, p.alpha + 0.15)})`
      const blur = `${p.blur}px`
      const sat = `${p.saturate}%`
      const lfs = `${this.lyricFontSize}px`
      const toX = `${this.tickerOffsetX}px`
      const anchor = `${this.lyricAnchor}`

      const root = document.documentElement
      root.style.setProperty('--mica-surface', surface)
      root.style.setProperty('--mica-surface-strong', surfaceStrong)
      root.style.setProperty('--mica-blur', blur)
      root.style.setProperty('--mica-saturate', sat)
      root.style.setProperty('--lyric-font-size', lfs)
      root.style.setProperty('--mini-ticker-offset', toX)
      root.style.setProperty('--lyric-anchor', anchor)
      root.style.setProperty('--mica-border', `rgba(255,255,255,${Math.max(0.25, p.alpha - 0.1)})`)
    },
    saveTheme() {
      const data = {
        theme: this.theme,
        alpha: this.alpha,
        blur: this.blur,
        saturate: this.saturate,
        lyricFontSize: this.lyricFontSize,
        lyricAnchor: this.lyricAnchor,
        tickerOffsetX: this.tickerOffsetX
      }
      localStorage.setItem('gmusic:theme', JSON.stringify(data))
    },
    loadTheme() {
      try {
        const raw = localStorage.getItem('gmusic:theme')
        if (raw) {
          const data = JSON.parse(raw)
          if (data && typeof data === 'object') {
            this.theme = data.theme ?? this.theme
            this.alpha = typeof data.alpha === 'number' ? data.alpha : this.alpha
            this.blur = typeof data.blur === 'number' ? data.blur : this.blur
            this.saturate = typeof data.saturate === 'number' ? data.saturate : this.saturate
            this.lyricFontSize = typeof data.lyricFontSize === 'number' ? data.lyricFontSize : this.lyricFontSize
            this.lyricAnchor = typeof data.lyricAnchor === 'number' ? data.lyricAnchor : this.lyricAnchor
            this.tickerOffsetX = typeof data.tickerOffsetX === 'number' ? data.tickerOffsetX : this.tickerOffsetX
          }
        }
      } catch {}
    }
  }
})
