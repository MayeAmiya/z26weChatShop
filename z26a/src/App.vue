<template>
  <div class="min-h-screen">
    <!-- Not logged in: show login window -->
    <div v-if="!authStore.isLoggedIn" class="min-h-screen flex items-center justify-center p-4">
      <router-view />
    </div>
    
    <!-- Logged in: Windows Explorer style -->
    <div v-else class="h-screen flex flex-col pb-7">
      <!-- Main Window -->
      <div class="win-window flex-1 m-2 flex flex-col overflow-hidden">
        <!-- Title Bar -->
        <div class="win-titlebar">
          <img :src="icons.computer" class="win-titlebar-icon" />
          <span>Z26B 商城管理系统</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn">_</button>
            <button class="win-titlebar-btn">□</button>
            <button class="win-titlebar-btn" @click="logout">×</button>
          </div>
        </div>

        <!-- Menu Bar -->
        <div class="win-menubar">
          <span class="win-menu-item" @click="logout" style="cursor:pointer">登出(L)</span>
        </div>

        <!-- Toolbar -->
        <div class="win-toolbar">
          <button class="win-toolbar-btn" @click="$router.go(-1)">
            <img :src="icons.folder" /> 返回
          </button>
          <div class="win-toolbar-separator"></div>
          <button class="win-toolbar-btn" @click="refresh">
            <img :src="icons.refresh" /> 刷新
          </button>
          <div class="win-toolbar-separator"></div>
          <button class="win-toolbar-btn" @click="$router.push('/')">
            <img :src="icons.dashboard" /> 仪表盘
          </button>
          <button class="win-toolbar-btn" @click="$router.push('/products')">
            <img :src="icons.product" /> 商品
          </button>
          <button class="win-toolbar-btn" @click="$router.push('/orders')">
            <img :src="icons.order" /> 订单
          </button>
          <button class="win-toolbar-btn" @click="$router.push('/users')">
            <img :src="icons.users" /> 用户
          </button>
          <button class="win-toolbar-btn" @click="$router.push('/categories')">
            <img :src="icons.category" /> 分类
          </button>
        </div>

        <!-- Address Bar -->
        <div class="win-address-bar">
          <label>地址(D):</label>
          <input type="text" class="win-input flex-1" :value="currentPath" readonly />
        </div>

        <!-- Content Area -->
        <div class="flex flex-1 overflow-hidden">
          <!-- Left Panel - Navigation Tree -->
          <div class="win-treeview w-48 flex-shrink-0">
            <div 
              v-for="item in menuItems" 
              :key="item.path"
              class="win-tree-item"
              :class="{ selected: $route.path === item.path || $route.path.startsWith(item.path + '/') }"
              @click="$router.push(item.path)"
            >
              <img :src="item.icon" class="win-tree-icon" />
              <span>{{ item.name }}</span>
            </div>
          </div>

          <!-- Right Panel - Main Content -->
          <div class="flex-1 overflow-auto bg-white">
            <router-view />
          </div>
        </div>

        <!-- Status Bar -->
        <div class="win-statusbar">
          <div class="win-statusbar-panel flex-1">
            {{ statusText }}
          </div>
          <div class="win-statusbar-panel w-32">
            {{ authStore.admin?.username || '管理员' }}
          </div>
          <div class="win-statusbar-panel w-24">
            {{ currentTime }}
          </div>
        </div>
      </div>

      <!-- Taskbar -->
      <div class="win-taskbar">
        <button class="win-start-btn">
          <img :src="icons.windows" />
          <span>开始</span>
        </button>
        <div class="win-taskbar-programs">
          <div class="win-taskbar-item active">
            <img :src="icons.computer" />
            <span>Z26B 商城管理系统</span>
          </div>
        </div>
        <div class="win-tray">
          <img :src="icons.clock" style="width:16px;height:16px" />
          {{ currentTime }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import icons from '@/assets/icons'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const currentTime = ref('')
let timeInterval = null

const menuItems = [
  { path: '/', name: '仪表盘', icon: icons.dashboard },
  { path: '/products', name: '商品管理', icon: icons.product },
  { path: '/orders', name: '订单管理', icon: icons.order },
  { path: '/users', name: '用户管理', icon: icons.users },
  { path: '/categories', name: '分类管理', icon: icons.category },
  { path: '/home-config', name: '首页配置', icon: icons.home },
]

// 页面访问历史
const pageHistory = ref([])

// 监听路由变化，记录访问顺序
router.afterEach((to) => {
  const item = menuItems.find(m => to.path === m.path || to.path.startsWith(m.path + '/'))
  if (item && (pageHistory.value.length === 0 || pageHistory.value[pageHistory.value.length - 1] !== item.name)) {
    // 避免重复添加相同页面
    pageHistory.value.push(item.name)
    // 最多保留5个历史记录
    if (pageHistory.value.length > 5) {
      pageHistory.value.shift()
    }
  }
})

const currentPath = computed(() => {
  if (pageHistory.value.length === 0) {
    const item = menuItems.find(m => route.path === m.path || route.path.startsWith(m.path + '/'))
    return item?.name || '未知'
  }
  return pageHistory.value.join(' > ')
})

const statusText = computed(() => {
  if (route.path === '/') return '欢迎使用 Z26B 商城管理系统'
  if (route.path.startsWith('/products')) return '商品管理 - 管理所有商品信息'
  if (route.path.startsWith('/orders')) return '订单管理 - 查看和处理订单'
  if (route.path.startsWith('/users')) return '用户管理 - 管理用户账户'
  if (route.path.startsWith('/categories')) return '分类管理 - 管理商品分类'
  if (route.path.startsWith('/home-config')) return '首页配置 - 管理轮播图和推荐商品'
  return '就绪'
})

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

const refresh = () => {
  window.location.reload()
}

const logout = () => {
  if (confirm('确定要退出登录吗？')) {
    authStore.logout()
    router.push('/login')
  }
}

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timeInterval) clearInterval(timeInterval)
})
</script>
