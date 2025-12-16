import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', {
  state: () => ({
    theme: 'current', // 'current' | 'glass'
    alpha: 0.45,      // 透明度 0-1
    blur: 18,         // 全局模糊
    saturate: 160,    // 饱和 %
    tickerFontSize: 14, // 底部栏歌词字号 px
    tickerOffsetX: 0,  // 底栏歌词水平偏移(px)
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
      const lfs = `${this.tickerFontSize}px`
      const toX = `${this.tickerOffsetX}px`

      const root = document.documentElement
      root.style.setProperty('--mica-surface', surface)
      root.style.setProperty('--mica-surface-strong', surfaceStrong)
      root.style.setProperty('--mica-blur', blur)
      root.style.setProperty('--mica-saturate', sat)
      root.style.setProperty('--ticker-font-size', lfs)
      root.style.setProperty('--mini-ticker-offset', toX)
      root.style.setProperty('--mica-border', `rgba(255,255,255,${Math.max(0.25, p.alpha - 0.1)})`)
    },
    saveTheme() {
      const data = {
        theme: this.theme,
        alpha: this.alpha,
        blur: this.blur,
        saturate: this.saturate,
        tickerFontSize: this.tickerFontSize,
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
            this.tickerFontSize = typeof data.tickerFontSize === 'number' ? data.tickerFontSize : this.tickerFontSize
            this.tickerOffsetX = typeof data.tickerOffsetX === 'number' ? data.tickerOffsetX : this.tickerOffsetX
          }
        }
      } catch {}
    }
  }
})
