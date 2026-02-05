<template>
  <div class="p-4">
    <!-- 页面标题 -->
    <div class="win-groupbox mb-4">
      <div class="win-groupbox-title flex items-center gap-2">
        <img :src="icons.users32" style="width:20px;height:20px" />
        客户分析
      </div>
      <div class="p-3">
        <p class="text-sm text-gray-600">分析客户行为数据，了解客户价值和消费习惯。</p>
      </div>
    </div>

    <!-- 概览统计 -->
    <div class="win-groupbox mb-4">
      <div class="win-groupbox-title">概览统计</div>
      <div class="grid grid-cols-4 gap-4 p-3">
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">总客户数</div>
          <div class="text-xl font-bold">{{ overview.totalCustomers || 0 }}</div>
        </div>
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">活跃客户（30天）</div>
          <div class="text-xl font-bold text-blue-600">{{ overview.activeCustomers || 0 }}</div>
        </div>
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">总消费金额</div>
          <div class="text-xl font-bold text-green-700">¥{{ (overview.totalSpent || 0).toFixed(2) }}</div>
        </div>
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">平均客单价</div>
          <div class="text-xl font-bold">¥{{ (overview.avgOrderValue || 0).toFixed(2) }}</div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <!-- 高价值客户排行 -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.money" style="width:16px;height:16px" />
          高价值客户 Top 10
        </div>
        <div class="p-2">
          <div class="win-toolbar mb-2">
            <button class="win-toolbar-btn" @click="fetchData">
              <img :src="icons.refresh" /> 刷新
            </button>
          </div>
          <div class="win-listview" style="height: 280px">
            <div class="win-listview-header">
              <div style="width:40px">#</div>
              <div style="flex:1">客户昵称</div>
              <div style="width:80px">订单数</div>
              <div style="width:100px">消费金额</div>
            </div>
            <div 
              v-for="(item, index) in topSpentCustomers" 
              :key="item._id" 
              class="win-listview-item"
            >
              <div style="width:40px">{{ index + 1 }}</div>
              <div style="flex:1" class="truncate">{{ item.user?.nickName || '匿名用户' }}</div>
              <div style="width:80px">{{ item.totalOrders || 0 }}</div>
              <div style="width:100px" class="text-green-700">¥{{ (item.totalSpent || 0).toFixed(2) }}</div>
            </div>
            <div v-if="topSpentCustomers.length === 0" class="text-center py-4 text-gray-500">
              <img :src="icons.info" style="width:16px;height:16px;vertical-align:middle" />
              暂无数据
            </div>
          </div>
        </div>
      </div>

      <!-- 客户等级分布 -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.chart" style="width:16px;height:16px" />
          客户等级分布
        </div>
        <div class="p-4">
          <table class="w-full text-sm">
            <tr class="border-b border-gray-300">
              <td class="py-2 font-bold">等级</td>
              <td class="py-2 font-bold text-right">客户数</td>
              <td class="py-2 font-bold text-right">占比</td>
            </tr>
            <tr v-for="level in levelDistribution" :key="level.level" class="border-b border-gray-300">
              <td class="py-2 flex items-center gap-2">
                <span :class="levelBadgeClass(level.level)">{{ levelText(level.level) }}</span>
              </td>
              <td class="py-2 text-right">{{ level.count }}</td>
              <td class="py-2 text-right">{{ level.percentage }}%</td>
            </tr>
            <tr v-if="levelDistribution.length === 0">
              <td colspan="3" class="py-4 text-center text-gray-500">暂无统计数据</td>
            </tr>
          </table>
        </div>
      </div>
    </div>

    <!-- 客户统计列表 -->
    <div class="win-groupbox mt-4">
      <div class="win-groupbox-title flex items-center gap-1">
        <img :src="icons.folder" style="width:16px;height:16px" />
        客户统计明细
      </div>
      <div class="p-2">
        <div class="win-toolbar mb-2">
          <button class="win-toolbar-btn" @click="fetchData" :disabled="loading">
            <img :src="icons.refresh" /> 刷新
          </button>
          <button class="win-toolbar-btn" @click="refreshAll" :disabled="loading">
            <img :src="icons.refresh" /> 重新统计
          </button>
          <button class="win-toolbar-btn">
            <img :src="icons.save" /> 导出
          </button>
        </div>
        <div class="win-listview" style="height: 250px">
          <div class="win-listview-header">
            <div style="flex:1">客户昵称</div>
            <div style="width:80px">等级</div>
            <div style="width:80px">订单数</div>
            <div style="width:100px">总消费</div>
            <div style="width:80px">平均客单</div>
            <div style="width:80px">退款次数</div>
            <div style="width:80px">浏览次数</div>
            <div style="width:100px">最后活跃</div>
          </div>
          <div 
            v-for="item in customerStats" 
            :key="item._id" 
            class="win-listview-item"
          >
            <div style="flex:1" class="truncate">{{ item.user?.nickName || '匿名用户' }}</div>
            <div style="width:80px">
              <span :class="levelBadgeClass(item.customerLevel)">{{ levelText(item.customerLevel) }}</span>
            </div>
            <div style="width:80px">{{ item.totalOrders || 0 }}</div>
            <div style="width:100px" class="text-green-700">¥{{ (item.totalSpent || 0).toFixed(2) }}</div>
            <div style="width:80px">¥{{ (item.avgOrderValue || 0).toFixed(2) }}</div>
            <div style="width:80px">{{ item.totalRefunds || 0 }}</div>
            <div style="width:80px">{{ item.totalViews || 0 }}</div>
            <div style="width:100px">{{ formatTime(item.lastActiveAt) }}</div>
          </div>
          <div v-if="customerStats.length === 0" class="text-center py-4 text-gray-500">
            <img :src="icons.info" style="width:16px;height:16px;vertical-align:middle" />
            暂无统计数据
          </div>
        </div>
        <!-- 分页 -->
        <div class="flex justify-end mt-2 gap-2">
          <button class="win-btn" :disabled="page <= 1" @click="page--; fetchData()">上一页</button>
          <span class="px-2 py-1">第 {{ page }} 页</span>
          <button class="win-btn" @click="page++; fetchData()">下一页</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import icons from '@/assets/icons'

const overview = ref({})
const topSpentCustomers = ref([])
const levelDistribution = ref([])
const customerStats = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

const levelText = (level) => {
  const map = {
    normal: '普通',
    silver: '白银',
    gold: '黄金',
    platinum: '铂金',
    diamond: '钻石',
    vip: 'VIP',
  }
  return map[level] || level
}

const levelBadgeClass = (level) => {
  const map = {
    normal: 'win-badge win-badge-gray',
    silver: 'win-badge win-badge-blue',
    gold: 'win-badge win-badge-yellow',
    platinum: 'win-badge win-badge-purple',
    diamond: 'win-badge win-badge-red',
    vip: 'win-badge win-badge-green',
  }
  return map[level] || 'win-badge'
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN')
}

const fetchData = async () => {
  loading.value = true
  try {
    const [overviewRes, topRes, levelRes, statsRes] = await Promise.all([
      api.get('/crm/customer/overview'),
      api.get('/crm/customer/top', { params: { orderBy: 'total_spent DESC', limit: 10 } }),
      api.get('/crm/customer/level-distribution'),
      api.get('/crm/customer/list', { params: { page: page.value, pageSize: pageSize.value } }),
    ])
    overview.value = overviewRes.data || {}
    topSpentCustomers.value = topRes.data || []
    levelDistribution.value = levelRes.data || []
    customerStats.value = statsRes.data?.list || []
    total.value = statsRes.data?.total || 0
  } catch (e) {
    console.error('Failed to fetch customer analysis data:', e)
  } finally {
    loading.value = false
  }
}

const refreshAll = async () => {
  loading.value = true
  try {
    await api.post('/crm/customer/refresh-all')
    await fetchData()
  } catch (e) {
    console.error('Failed to refresh:', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>
