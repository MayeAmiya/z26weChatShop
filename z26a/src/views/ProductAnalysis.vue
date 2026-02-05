<template>
  <div class="p-4">
    <!-- 页面标题 -->
    <div class="win-groupbox mb-4">
      <div class="win-groupbox-title flex items-center gap-2">
        <img :src="icons.chart32" style="width:20px;height:20px" />
        商品分析
      </div>
      <div class="p-3">
        <p class="text-sm text-gray-600">分析商品数据，了解商品表现和销售趋势。</p>
      </div>
    </div>

    <!-- 概览统计 -->
    <div class="win-groupbox mb-4">
      <div class="win-groupbox-title">概览统计</div>
      <div class="grid grid-cols-6 gap-4 p-3">
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">总商品数</div>
          <div class="text-xl font-bold">{{ overview.totalProducts || 0 }}</div>
        </div>
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">总浏览量</div>
          <div class="text-xl font-bold">{{ overview.totalViews || 0 }}</div>
        </div>
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">总销量</div>
          <div class="text-xl font-bold text-green-700">{{ overview.totalSales || 0 }}</div>
        </div>
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">总营收</div>
          <div class="text-xl font-bold text-green-700">¥{{ (overview.totalRevenue || 0).toFixed(2) }}</div>
        </div>
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">平均转化率</div>
          <div class="text-xl font-bold">{{ ((overview.avgConversionRate || 0) * 100).toFixed(1) }}%</div>
        </div>
        <div class="win-panel p-3 text-center">
          <div class="text-xs text-gray-600">平均评分</div>
          <div class="text-xl font-bold text-orange-600">{{ (overview.avgScore || 0).toFixed(1) }}</div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <!-- 热销商品排行 -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.chart" style="width:16px;height:16px" />
          热销商品 Top 10
        </div>
        <div class="p-2">
          <div class="win-toolbar mb-2">
            <button class="win-toolbar-btn" @click="fetchData">
              <img :src="icons.refresh" /> 刷新
            </button>
          </div>
          <div class="win-listview" style="height: 300px">
            <div class="win-listview-header">
              <div style="width:40px">#</div>
              <div style="flex:1">商品名称</div>
              <div style="width:80px">销量</div>
              <div style="width:100px">营收</div>
            </div>
            <div 
              v-for="(item, index) in topSalesProducts" 
              :key="item._id" 
              class="win-listview-item"
            >
              <div style="width:40px">{{ index + 1 }}</div>
              <div style="flex:1" class="truncate">{{ item.spu?.name || '-' }}</div>
              <div style="width:80px">{{ item.totalSales || 0 }}</div>
              <div style="width:100px" class="text-green-700">¥{{ (item.totalRevenue || 0).toFixed(2) }}</div>
            </div>
            <div v-if="topSalesProducts.length === 0" class="text-center py-4 text-gray-500">
              <img :src="icons.info" style="width:16px;height:16px;vertical-align:middle" />
              暂无数据
            </div>
          </div>
        </div>
      </div>

      <!-- 浏览量排行 -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.search" style="width:16px;height:16px" />
          浏览量 Top 10
        </div>
        <div class="p-2">
          <div class="win-toolbar mb-2">
            <button class="win-toolbar-btn" @click="fetchData">
              <img :src="icons.refresh" /> 刷新
            </button>
          </div>
          <div class="win-listview" style="height: 300px">
            <div class="win-listview-header">
              <div style="width:40px">#</div>
              <div style="flex:1">商品名称</div>
              <div style="width:80px">浏览量</div>
              <div style="width:80px">转化率</div>
            </div>
            <div 
              v-for="(item, index) in topViewProducts" 
              :key="item._id" 
              class="win-listview-item"
            >
              <div style="width:40px">{{ index + 1 }}</div>
              <div style="flex:1" class="truncate">{{ item.spu?.name || '-' }}</div>
              <div style="width:80px">{{ item.totalViews || 0 }}</div>
              <div style="width:80px">{{ ((item.conversionRate || 0) * 100).toFixed(1) }}%</div>
            </div>
            <div v-if="topViewProducts.length === 0" class="text-center py-4 text-gray-500">
              <img :src="icons.info" style="width:16px;height:16px;vertical-align:middle" />
              暂无数据
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 商品统计列表 -->
    <div class="win-groupbox mt-4">
      <div class="win-groupbox-title flex items-center gap-1">
        <img :src="icons.folder" style="width:16px;height:16px" />
        商品统计明细
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
            <div style="flex:1">商品名称</div>
            <div style="width:80px">浏览量</div>
            <div style="width:80px">加购数</div>
            <div style="width:80px">销量</div>
            <div style="width:100px">营收</div>
            <div style="width:80px">评论数</div>
            <div style="width:60px">评分</div>
            <div style="width:80px">转化率</div>
          </div>
          <div 
            v-for="item in productStats" 
            :key="item._id" 
            class="win-listview-item"
          >
            <div style="flex:1" class="truncate">{{ item.spu?.name || '-' }}</div>
            <div style="width:80px">{{ item.totalViews || 0 }}</div>
            <div style="width:80px">{{ item.totalCarts || 0 }}</div>
            <div style="width:80px">{{ item.totalSales || 0 }}</div>
            <div style="width:100px" class="text-green-700">¥{{ (item.totalRevenue || 0).toFixed(2) }}</div>
            <div style="width:80px">{{ item.totalComments || 0 }}</div>
            <div style="width:60px">{{ (item.avgScore || 0).toFixed(1) }}</div>
            <div style="width:80px">{{ ((item.conversionRate || 0) * 100).toFixed(1) }}%</div>
          </div>
          <div v-if="productStats.length === 0" class="text-center py-4 text-gray-500">
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
const topSalesProducts = ref([])
const topViewProducts = ref([])
const productStats = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const [overviewRes, topSalesRes, topViewsRes, statsRes] = await Promise.all([
      api.get('/crm/product/overview'),
      api.get('/crm/product/top', { params: { orderBy: 'total_revenue DESC', limit: 10 } }),
      api.get('/crm/product/top', { params: { orderBy: 'total_views DESC', limit: 10 } }),
      api.get('/crm/product/list', { params: { page: page.value, pageSize: pageSize.value } }),
    ])
    overview.value = overviewRes.data || {}
    topSalesProducts.value = topSalesRes.data || []
    topViewProducts.value = topViewsRes.data || []
    productStats.value = statsRes.data?.list || []
    total.value = statsRes.data?.total || 0
  } catch (e) {
    console.error('Failed to fetch product analysis data:', e)
  } finally {
    loading.value = false
  }
}

const refreshAll = async () => {
  loading.value = true
  try {
    await api.post('/crm/product/refresh-all')
    await fetchData()
  } catch (e) {
    console.error('Failed to refresh:', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>
