import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import AdminView from '@/views/AdminView.vue'
import LandingView from '@/views/LandingView.vue'
import LoginView from '@/views/LoginView.vue'
import { getAuthToken } from '@/api/client'

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
    path: '/login',
    name: 'Login',
    component: LoginView
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

router.beforeEach((to, _from, next) => {
  const token = getAuthToken()
  if (to.name === 'Landing') {
    return next()
  }
  if (to.name === 'Login') {
    if (token) {
      return next('/')
    }
    return next()
  }
  if (!token) {
    return next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
  }
  next()
})
