import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router/index.js'
import App from './App.vue'
import './App.css'

// 先创建 app 与 pinia
const app = createApp(App)
const pinia = createPinia()
app.use(pinia)

// 读取并应用主题（本地缓存）
import { useUiStore } from './stores/ui'
const ui = useUiStore(pinia)
ui.loadTheme()
ui.applyTheme()

// 再挂载路由与应用
app.use(router)
app.mount('#root')
