import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const admin = ref(JSON.parse(localStorage.getItem('admin_info') || 'null'))

  // 检查token是否过期
  const isTokenExpired = () => {
    if (!token.value) return true
    try {
      // JWT token格式: header.payload.signature
      const payload = JSON.parse(atob(token.value.split('.')[1]))
      // exp是Unix时间戳（秒），比较当前时间
      return payload.exp * 1000 < Date.now()
    } catch (e) {
      return true
    }
  }

  const isLoggedIn = computed(() => !!token.value && !isTokenExpired())

  const login = async (credentials) => {
    const res = await api.post('/login', credentials)
    token.value = res.token
    admin.value = res.admin
    localStorage.setItem('admin_token', res.token)
    localStorage.setItem('admin_info', JSON.stringify(res.admin))
    return res
  }

  const logout = () => {
    token.value = ''
    admin.value = null
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_info')
  }

  // 验证token有效性（调用后端）
  const validateToken = async () => {
    if (!token.value || isTokenExpired()) {
      logout()
      return false
    }
    try {
      await api.get('/profile')
      return true
    } catch (e) {
      logout()
      return false
    }
  }

  const fetchProfile = async () => {
    try {
      const res = await api.get('/profile')
      admin.value = res
      localStorage.setItem('admin_info', JSON.stringify(res))
    } catch (e) {
      logout()
    }
  }

  return {
    token,
    admin,
    isLoggedIn,
    isTokenExpired,
    login,
    logout,
    validateToken,
    fetchProfile,
  }
})
