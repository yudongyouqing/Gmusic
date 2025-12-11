import { createRouter, createWebHistory } from 'vue-router'

import Library from '../views/Library.vue'
import NowPlaying from '../views/NowPlaying.vue'
import Queue from '../views/Queue.vue'

const routes = [
  { path: '/', name: 'library', component: Library },
  { path: '/now-playing', name: 'nowPlaying', component: NowPlaying },
  { path: '/queue', name: 'queue', component: Queue }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

