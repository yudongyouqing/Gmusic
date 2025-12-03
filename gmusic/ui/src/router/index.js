import { createRouter, createWebHistory } from 'vue-router'

import Library from '../views/Library.vue'
import NowPlaying from '../views/NowPlaying.vue'

const routes = [
  { path: '/', name: 'library', component: Library },
  { path: '/now-playing', name: 'nowPlaying', component: NowPlaying }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

