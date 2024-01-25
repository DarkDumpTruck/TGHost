import { createRouter, createWebHistory } from 'vue-router'
import GameView from '../views/GameView.vue'
import RoomsView from '../views/RoomsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/rooms'
    },
    {
      path: '/rooms',
      name: 'rooms',
      component: RoomsView
    },
    {
      path: '/game/:roomId/:playerId',
      name: 'game',
      component: GameView,
    },
  ]
})

export default router
