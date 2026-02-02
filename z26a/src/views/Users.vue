<template>
  <div class="p-3 h-full flex flex-col">
    <!-- Toolbar -->
    <div class="win-toolbar mb-2">
      <button class="win-toolbar-btn" @click="fetchUsers">
        <img :src="icons.refresh" /> 刷新
      </button>
      <div class="flex-1"></div>
      <div class="flex items-center gap-2">
        <img :src="icons.search" style="width:16px;height:16px" />
        <input
          v-model="search"
          type="text"
          placeholder="搜索用户..."
          class="win-input"
          style="width: 180px"
          @input="debounceFetch"
        />
      </div>
    </div>

    <!-- Users ListView -->
    <div class="win-listview flex-1">
      <div class="win-listview-header">
        <div style="width:50px"></div>
        <div style="width:150px">用户名</div>
        <div style="width:120px">手机号</div>
        <div style="width:100px">钱包余额</div>
        <div style="width:80px">订单数</div>
        <div style="width:150px">注册时间</div>
        <div style="flex:1">操作</div>
      </div>
      
      <div 
        v-for="user in users" 
        :key="user._id" 
        class="win-listview-item"
      >
        <div style="width:50px">
          <img 
            :src="user.avatarUrl || icons.user32" 
            style="width:32px;height:32px;object-fit:cover" 
            class="border border-gray-400"
          />
        </div>
        <div style="width:150px">
          <div class="font-medium truncate">{{ user.nickName || '未设置昵称' }}</div>
          <div class="text-xs text-gray-500">ID: {{ user._id?.slice(-8) }}</div>
        </div>
        <div style="width:120px">{{ user.phone || '-' }}</div>
        <div style="width:100px" class="text-red-600 font-bold">
          ¥{{ (user.walletBalance || 0).toFixed(2) }}
        </div>
        <div style="width:80px">{{ user.orderCount || 0 }}</div>
        <div style="width:150px" class="text-xs text-gray-600">
          {{ formatTime(user.createdAt) }}
        </div>
        <div style="flex:1" class="flex gap-1">
          <button class="win-btn win-btn-sm" @click.stop="showRechargeModal(user)">
            <img :src="icons.money" style="width:12px;height:12px" /> 充值
          </button>
          <button class="win-btn win-btn-sm" @click.stop="viewUserOrders(user)">
            <img :src="icons.order" style="width:12px;height:12px" /> 订单
          </button>
        </div>
      </div>

      <div v-if="users.length === 0" class="text-center py-8 text-gray-500">
        <img :src="icons.users32" style="width:48px;height:48px;margin:0 auto 8px" />
        <div>暂无用户数据</div>
      </div>
    </div>

    <!-- Pagination -->
    <div class="win-toolbar mt-2 justify-between">
      <div class="text-sm flex items-center gap-1">
        <img :src="icons.info" style="width:16px;height:16px" />
        共 {{ total }} 个用户
      </div>
      <div class="flex gap-1">
        <button class="win-btn win-btn-sm" :disabled="page === 1" @click="page--; fetchUsers()">
          &lt; 上一页
        </button>
        <span class="win-panel px-3 py-1 text-sm">
          第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) || 1 }} 页
        </span>
        <button 
          class="win-btn win-btn-sm" 
          :disabled="page >= Math.ceil(total / pageSize)" 
          @click="page++; fetchUsers()"
        >
          下一页 &gt;
        </button>
      </div>
    </div>

    <!-- Recharge Modal -->
    <div v-if="rechargeUser" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="win-window" style="width: 350px;">
        <div class="win-titlebar">
          <img :src="icons.money" class="win-titlebar-icon" />
          <span>充值余额</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn" @click="rechargeUser = null">×</button>
          </div>
        </div>
        
        <div class="bg-[#c0c0c0] p-4">
          <div class="win-panel p-2 mb-3">
            <div class="flex items-center gap-2 mb-1">
              <img :src="icons.user" style="width:16px;height:16px" />
              <strong>{{ rechargeUser.nickName }}</strong>
            </div>
            <div class="text-sm text-gray-600">当前余额: ¥{{ (rechargeUser.walletBalance || 0).toFixed(2) }}</div>
          </div>

          <form @submit.prevent="doRecharge">
            <table class="w-full">
              <tr>
                <td class="py-1 pr-2 text-right whitespace-nowrap">
                  <label>充值金额:</label>
                </td>
                <td class="py-1">
                  <input 
                    v-model.number="rechargeAmount" 
                    type="number" 
                    step="0.01" 
                    class="win-input w-full" 
                    placeholder="0.00"
                    required 
                  />
                </td>
              </tr>
            </table>

            <div class="flex justify-end gap-2 mt-4 pt-3 border-t border-gray-400">
              <button type="submit" class="win-btn" style="min-width:80px">
                <img :src="icons.check" style="width:16px;height:16px" /> 确定
              </button>
              <button type="button" class="win-btn" style="min-width:80px" @click="rechargeUser = null">
                <img :src="icons.cross" style="width:16px;height:16px" /> 取消
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- User Orders Modal -->
    <div v-if="userOrders" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="win-window" style="width: 550px; max-height: 80vh;">
        <div class="win-titlebar">
          <img :src="icons.order" class="win-titlebar-icon" />
          <span>{{ userOrders.user?.nickName || '用户' }} 的订单</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn" @click="userOrders = null">×</button>
          </div>
        </div>
        
        <div class="bg-[#c0c0c0] p-3 overflow-y-auto" style="max-height: calc(80vh - 30px);">
          <div class="win-listview" style="max-height:300px">
            <div class="win-listview-header">
              <div style="width:100px">订单号</div>
              <div style="width:80px">金额</div>
              <div style="width:80px">状态</div>
              <div style="flex:1">时间</div>
            </div>
            <div v-for="order in userOrders.orders" :key="order._id" class="win-listview-item">
              <div style="width:100px">{{ order._id?.slice(-8) }}</div>
              <div style="width:80px" class="text-red-600">¥{{ order.finalPrice }}</div>
              <div style="width:80px">
                <span :class="statusBadgeClass(order.status)">{{ statusText(order.status) }}</span>
              </div>
              <div style="flex:1" class="text-xs">{{ formatTime(order.createdAt) }}</div>
            </div>
            <div v-if="!userOrders.orders?.length" class="text-center py-4 text-gray-500">
              暂无订单
            </div>
          </div>

          <div class="flex justify-end gap-2 mt-3">
            <button class="win-btn" @click="userOrders = null">
              <img :src="icons.cross" style="width:16px;height:16px" /> 关闭
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import icons from '@/assets/icons'

const users = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const search = ref('')
const rechargeUser = ref(null)
const rechargeAmount = ref(0)
const userOrders = ref(null)

let debounceTimer = null
const debounceFetch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    fetchUsers()
  }, 500)
}

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const statusText = (status) => {
  const map = {
    TO_PAY: '待支付',
    TO_SEND: '待发货',
    TO_RECEIVE: '待收货',
    FINISHED: '已完成',
    CANCELED: '已取消',
    RETURN_FINISH: '已退款',
  }
  return map[status] || status
}

const statusBadgeClass = (status) => {
  const map = {
    TO_PAY: 'win-badge win-badge-yellow',
    TO_SEND: 'win-badge win-badge-blue',
    TO_RECEIVE: 'win-badge win-badge-purple',
    FINISHED: 'win-badge win-badge-green',
    CANCELED: 'win-badge win-badge-gray',
    RETURN_FINISH: 'win-badge win-badge-red',
  }
  return map[status] || 'win-badge'
}

const fetchUsers = async () => {
  try {
    const params = {
      page: page.value,
      pageSize,
    }
    if (search.value) params.keyword = search.value

    const res = await api.get('/users', { params })
    users.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error('Failed to fetch users:', e)
  }
}

const showRechargeModal = (user) => {
  rechargeUser.value = user
  rechargeAmount.value = 0
}

const doRecharge = async () => {
  try {
    await api.post(`/users/${rechargeUser.value._id}/recharge`, {
      amount: rechargeAmount.value
    })
    alert('充值成功')
    rechargeUser.value = null
    fetchUsers()
  } catch (e) {
    alert('充值失败: ' + (e.error || e.message))
  }
}

const viewUserOrders = async (user) => {
  try {
    const res = await api.get(`/users/${user._id}/orders`)
    userOrders.value = {
      user,
      orders: res.data || []
    }
  } catch (e) {
    console.error('Failed to fetch user orders:', e)
  }
}

onMounted(fetchUsers)
</script>
