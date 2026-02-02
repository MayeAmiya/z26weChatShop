<template>
  <div class="p-4">
    <!-- Stats Panel -->
    <div class="win-groupbox mb-4">
      <div class="win-groupbox-title">系统统计</div>
      <div class="grid grid-cols-4 gap-4 p-3">
        <!-- Today Orders -->
        <div class="win-panel p-3 flex items-center gap-3">
          <img :src="icons.order32" style="width:32px;height:32px" />
          <div>
            <div class="text-xs text-gray-600">今日订单</div>
            <div class="text-xl font-bold">{{ stats.todayOrders }}</div>
          </div>
        </div>
        
        <!-- Today Sales -->
        <div class="win-panel p-3 flex items-center gap-3">
          <img :src="icons.money32" style="width:32px;height:32px" />
          <div>
            <div class="text-xs text-gray-600">今日销售额</div>
            <div class="text-xl font-bold text-green-700">¥{{ (stats.todaySales || 0).toFixed(2) }}</div>
          </div>
        </div>
        
        <!-- To Ship Orders -->
        <div class="win-panel p-3 flex items-center gap-3">
          <img :src="icons.cart32" style="width:32px;height:32px" />
          <div>
            <div class="text-xs text-gray-600">待发货订单</div>
            <div class="text-xl font-bold text-orange-600">{{ stats.toShipOrders }}</div>
          </div>
        </div>
        
        <!-- Total Users -->
        <div class="win-panel p-3 flex items-center gap-3">
          <img :src="icons.users32" style="width:32px;height:32px" />
          <div>
            <div class="text-xs text-gray-600">用户总数</div>
            <div class="text-xl font-bold">{{ stats.totalUsers }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <!-- Recent Orders -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.order" style="width:16px;height:16px" />
          最近订单
        </div>
        <div class="p-2">
          <!-- Toolbar -->
          <div class="win-toolbar mb-2">
            <button class="win-toolbar-btn" @click="$router.push('/orders')">
              <img :src="icons.folder" /> 查看全部
            </button>
            <button class="win-toolbar-btn" @click="fetchData">
              <img :src="icons.refresh" /> 刷新
            </button>
          </div>

          <!-- List View -->
          <div class="win-listview" style="height: 200px">
            <div class="win-listview-header">
              <div style="width:100px">订单号</div>
              <div style="width:80px">用户</div>
              <div style="width:80px">金额</div>
              <div style="flex:1">状态</div>
            </div>
            <div 
              v-for="order in recentOrders" 
              :key="order._id" 
              class="win-listview-item"
              @click="$router.push('/orders/' + order._id)"
            >
              <div style="width:100px">{{ order._id?.slice(-8) }}</div>
              <div style="width:80px">{{ order.user?.nickName || '用户' }}</div>
              <div style="width:80px" class="text-red-600">¥{{ order.finalPrice }}</div>
              <div style="flex:1">
                <span :class="statusBadgeClass(order.status)">{{ statusText(order.status) }}</span>
              </div>
            </div>
            <div v-if="recentOrders.length === 0" class="text-center py-4 text-gray-500">
              <img :src="icons.info" style="width:16px;height:16px;vertical-align:middle" />
              暂无订单数据
            </div>
          </div>
        </div>
      </div>

      <!-- Order Status Stats -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.chart" style="width:16px;height:16px" />
          订单状态统计
        </div>
        <div class="p-3">
          <table class="w-full text-sm">
            <tr v-for="stat in orderStats" :key="stat.status" class="border-b border-gray-300">
              <td class="py-2 flex items-center gap-2">
                <img :src="getStatusIcon(stat.status)" style="width:16px;height:16px" />
                {{ statusText(stat.status) }}
              </td>
              <td class="py-2 text-right font-bold">
                {{ stat.count }}
              </td>
            </tr>
            <tr v-if="orderStats.length === 0">
              <td colspan="2" class="py-4 text-center text-gray-500">暂无统计数据</td>
            </tr>
          </table>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="win-groupbox mt-4">
      <div class="win-groupbox-title flex items-center gap-1">
        <img :src="icons.folder" style="width:16px;height:16px" />
        快捷操作
      </div>
      <div class="p-4 flex gap-6">
        <div class="win-desktop-icon" @click="$router.push('/products/new')">
          <img :src="icons.add32" />
          <span>新建商品</span>
        </div>
        <div class="win-desktop-icon" @click="$router.push('/orders')">
          <img :src="icons.order32" />
          <span>处理订单</span>
        </div>
        <div class="win-desktop-icon" @click="$router.push('/users')">
          <img :src="icons.users32" />
          <span>用户管理</span>
        </div>
        <div class="win-desktop-icon" @click="$router.push('/categories')">
          <img :src="icons.category32" />
          <span>分类管理</span>
        </div>
        <div class="win-desktop-icon" @click="$router.push('/home-config')">
          <img :src="icons.dashboard32" />
          <span>首页配置</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import icons from '@/assets/icons'

const stats = ref({
  todayOrders: 0,
  todaySales: 0,
  toShipOrders: 0,
  totalUsers: 0,
})
const recentOrders = ref([])
const orderStats = ref([])

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

const getStatusIcon = (status) => {
  const map = {
    TO_PAY: icons.clock,
    TO_SEND: icons.cart,
    TO_RECEIVE: icons.order,
    FINISHED: icons.check,
    CANCELED: icons.cross,
    RETURN_FINISH: icons.money,
  }
  return map[status] || icons.info
}

const fetchData = async () => {
  try {
    const [statsRes, ordersRes, orderStatsRes] = await Promise.all([
      api.get('/stats'),
      api.get('/orders', { params: { pageSize: 5 } }),
      api.get('/stats/order-status'),
    ])
    stats.value = statsRes.data || {}
    recentOrders.value = ordersRes.data?.list || []
    orderStats.value = orderStatsRes.data || []
  } catch (e) {
    console.error('Failed to fetch dashboard data:', e)
  }
}

onMounted(fetchData)
</script>
