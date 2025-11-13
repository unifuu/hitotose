import { createRouter, createWebHistory } from 'vue-router'
import Game from '../components/Game.vue'
import Login from '../components/Login.vue'
import Rank from '../components/Rank.vue'

const routes = [
  { path: '/:pathMatch(.*)*', redirect: '/' },
  {
    path: '/',
    name: 'Home',
    component: Game,
  },
  {
    path: '/game',
    name: 'Game',
    component: Game,
  },
  {
    path: '/fuu',
    name: 'Login',
    component: Login,
  },
  {
    path: '/rank',
    name: 'Rank',
    component: Rank,
  }
]

// Create the router instance
const router = createRouter({
  history: createWebHistory(), // Use HTML5 History API
  routes,
})

export default router