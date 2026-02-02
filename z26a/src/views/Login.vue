<template>
  <div class="win-window" style="width: 400px;">
    <!-- Title Bar -->
    <div class="win-titlebar">
      <img :src="icons.key" class="win-titlebar-icon" />
      <span>登录到 Z26B 管理系统</span>
      <div class="win-titlebar-buttons">
        <button class="win-titlebar-btn">×</button>
      </div>
    </div>

    <!-- Dialog Content -->
    <div class="p-4 bg-[#c0c0c0]">
      <div class="flex gap-4 mb-4">
        <img :src="icons.key32" style="width:48px;height:48px" />
        <div class="flex-1">
          <p class="mb-2">请输入您的管理员账号信息</p>
          <p class="text-xs text-gray-600">连接到: Z26B 商城管理系统</p>
        </div>
      </div>

      <form @submit.prevent="handleLogin">
        <table class="w-full">
          <tr>
            <td class="py-1 pr-2 text-right whitespace-nowrap">
              <label>用户名(<u>U</u>):</label>
            </td>
            <td class="py-1">
              <input
                v-model="form.email"
                type="email"
                class="win-input w-full"
                placeholder="admin@z26b.com"
                required
              />
            </td>
          </tr>
          <tr>
            <td class="py-1 pr-2 text-right whitespace-nowrap">
              <label>密码(<u>P</u>):</label>
            </td>
            <td class="py-1">
              <input
                v-model="form.password"
                type="password"
                class="win-input w-full"
                placeholder="••••••••"
                required
              />
            </td>
          </tr>
        </table>

        <div v-if="error" class="win-panel p-2 mt-3 text-red-700 text-sm flex items-center gap-2">
          <img :src="icons.warning" style="width:16px;height:16px" />
          {{ error }}
        </div>

        <!-- Buttons -->
        <div class="flex justify-end gap-2 mt-4 pt-3 border-t border-gray-400">
          <button type="submit" class="win-btn" :disabled="loading" style="min-width:80px">
            <img :src="icons.check" v-if="!loading" style="width:16px;height:16px" />
            {{ loading ? '登录中...' : '确定' }}
          </button>
          <button type="button" class="win-btn" style="min-width:80px" @click="resetForm">
            <img :src="icons.cross" style="width:16px;height:16px" />
            取消
          </button>
        </div>
      </form>

      <div class="mt-4 pt-3 border-t border-gray-400">
        <div class="win-groupbox p-2">
          <div class="win-groupbox-title">提示</div>
          <p class="text-xs">默认管理员账号: admin@z26b.com / admin123</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import icons from '@/assets/icons'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref('')

const form = reactive({
  email: '',
  password: '',
})

const resetForm = () => {
  form.email = ''
  form.password = ''
  error.value = ''
}

const handleLogin = async () => {
  error.value = ''
  loading.value = true
  
  try {
    await authStore.login(form)
    router.push('/')
  } catch (e) {
    error.value = e.error || '登录失败，请检查邮箱和密码'
  } finally {
    loading.value = false
  }
}
</script>
