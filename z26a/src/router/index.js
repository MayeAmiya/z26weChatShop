import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    name: 'dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/products',
    name: 'products',
    component: () => import('@/views/Products.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/products/:id',
    name: 'product-detail',
    component: () => import('@/views/ProductDetail.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/orders',
    name: 'orders',
    component: () => import('@/views/Orders.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/users',
    name: 'users',
    component: () => import('@/views/Users.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/categories',
    name: 'categories',
    component: () => import('@/views/Categories.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/home-config',
    name: 'home-config',
    component: () => import('@/views/HomeConfig.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/product-analysis',
    name: 'product-analysis',
    component: () => import('@/views/ProductAnalysis.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/customer-analysis',
    name: 'customer-analysis',
    component: () => import('@/views/CustomerAnalysis.vue'),
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 检查token是否过期
  if (to.meta.requiresAuth) {
    if (!authStore.isLoggedIn) {
      // token不存在或已过期
      authStore.logout()
      next({ name: 'login' })
      return
    }
  }
  
  if (to.meta.guest && authStore.isLoggedIn) {
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

export default router
