import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import AdminView from '@/views/AdminView.vue'
import LandingView from '@/views/LandingView.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomeView
  },
  {
    path: '/admin',
    name: 'Admin',
    component: AdminView
  },
  {
    path: '/s/:fileId',
    name: 'Landing',
    component: LandingView
  }
]

export const router = createRouter({
  history: createWebHistory(),
  routes
})
