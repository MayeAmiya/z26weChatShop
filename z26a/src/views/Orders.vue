<template>
  <div class="p-3 h-full flex flex-col">
    <!-- Toolbar -->
    <div class="win-toolbar mb-2">
      <button class="win-toolbar-btn" @click="fetchOrders">
        <img :src="icons.refresh" /> 刷新
      </button>
      <div class="win-toolbar-separator"></div>
      <div class="flex items-center gap-2">
        <img :src="icons.search" style="width:16px;height:16px" />
        <input
          v-model="search"
          type="text"
          placeholder="搜索订单号..."
          class="win-input"
          style="width: 150px"
          @input="debounceFetch"
        />
        <select v-model="statusFilter" class="win-select" @change="fetchOrders">
          <option value="">全部状态</option>
          <option value="TO_PAY">待支付</option>
          <option value="TO_SEND">待发货</option>
          <option value="TO_RECEIVE">待收货</option>
          <option value="FINISHED">已完成</option>
          <option value="CANCELED">已取消</option>
          <option value="RETURN_FINISH">已退款</option>
        </select>
      </div>
    </div>

    <!-- Orders ListView -->
    <div class="win-listview flex-1">
      <div class="win-listview-header">
        <div style="width:100px">订单号</div>
        <div style="width:100px">用户</div>
        <div style="width:80px">金额</div>
        <div style="width:80px">状态</div>
        <div style="width:150px">下单时间</div>
        <div style="flex:1">操作</div>
      </div>
      
      <div 
        v-for="order in orders" 
        :key="order._id" 
        class="win-listview-item"
        @dblclick="viewOrder(order)"
      >
        <div style="width:100px">
          <div class="font-medium">{{ order._id?.slice(-8) }}</div>
          <div class="text-xs text-gray-500">{{ order.skuDetails?.length || 0 }} 件</div>
        </div>
        <div style="width:100px">
          <div class="flex items-center gap-1">
            <img :src="icons.user" style="width:14px;height:14px" />
            {{ order.user?.nickName || '用户' }}
          </div>
          <div class="text-xs text-gray-500">{{ order.user?.phone }}</div>
        </div>
        <div style="width:80px" class="text-red-600 font-bold">
          ¥{{ order.finalPrice }}
        </div>
        <div style="width:80px">
          <span :class="statusBadgeClass(order.status)">{{ statusText(order.status) }}</span>
        </div>
        <div style="width:150px" class="text-xs text-gray-600">
          {{ formatTime(order.createdAt) }}
        </div>
        <div style="flex:1" class="flex gap-1">
          <button class="win-btn win-btn-sm" @click.stop="viewOrder(order)">
            <img :src="icons.info" style="width:12px;height:12px" /> 详情
          </button>
          <button 
            v-if="order.status === 'TO_SEND'"
            class="win-btn win-btn-sm"
            @click.stop="shipOrder(order)"
          >
            <img :src="icons.check" style="width:12px;height:12px" /> 发货
          </button>
          <button 
            v-if="['TO_SEND', 'TO_RECEIVE'].includes(order.status)"
            class="win-btn win-btn-sm"
            @click.stop="refundOrder(order)"
          >
            <img :src="icons.money" style="width:12px;height:12px" /> 退款
          </button>
        </div>
      </div>

      <div v-if="orders.length === 0" class="text-center py-8 text-gray-500">
        <img :src="icons.order32" style="width:48px;height:48px;margin:0 auto 8px" />
        <div>暂无订单数据</div>
      </div>
    </div>

    <!-- Pagination -->
    <div class="win-toolbar mt-2 justify-between">
      <div class="text-sm flex items-center gap-1">
        <img :src="icons.info" style="width:16px;height:16px" />
        共 {{ total }} 个订单
      </div>
      <div class="flex gap-1">
        <button class="win-btn win-btn-sm" :disabled="page === 1" @click="page--; fetchOrders()">
          &lt; 上一页
        </button>
        <span class="win-panel px-3 py-1 text-sm">
          第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) || 1 }} 页
        </span>
        <button 
          class="win-btn win-btn-sm" 
          :disabled="page >= Math.ceil(total / pageSize)" 
          @click="page++; fetchOrders()"
        >
          下一页 &gt;
        </button>
      </div>
    </div>

    <!-- Order detail dialog -->
    <div v-if="selectedOrder" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="win-window" style="width: 600px; max-height: 80vh;">
        <div class="win-titlebar">
          <img :src="icons.order" class="win-titlebar-icon" />
          <span>订单详情 - {{ selectedOrder._id?.slice(-8) }}</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn" @click="selectedOrder = null">×</button>
          </div>
        </div>
        
        <div class="bg-[#c0c0c0] p-3 overflow-y-auto" style="max-height: calc(80vh - 30px);">
          <!-- Order Info -->
          <div class="win-groupbox mb-3">
            <div class="win-groupbox-title flex items-center gap-1">
              <img :src="icons.info" style="width:16px;height:16px" />
              订单信息
            </div>
            <table class="text-sm w-full">
              <tr>
                <td class="py-1 pr-4 text-gray-600">订单号:</td>
                <td class="py-1">{{ selectedOrder._id }}</td>
                <td class="py-1 pr-4 text-gray-600">状态:</td>
                <td class="py-1"><span :class="statusBadgeClass(selectedOrder.status)">{{ statusText(selectedOrder.status) }}</span></td>
              </tr>
              <tr>
                <td class="py-1 pr-4 text-gray-600">下单时间:</td>
                <td class="py-1">{{ formatTime(selectedOrder.createdAt) }}</td>
                <td class="py-1 pr-4 text-gray-600">支付时间:</td>
                <td class="py-1">{{ formatTime(selectedOrder.payedAt) || '-' }}</td>
              </tr>
            </table>
          </div>

          <!-- Address Info -->
          <div class="win-groupbox mb-3">
            <div class="win-groupbox-title flex items-center gap-1">
              <img :src="icons.user" style="width:16px;height:16px" />
              收货信息
            </div>
            <div class="text-sm p-2">
              <div class="flex items-center gap-2 mb-1">
                <strong>{{ selectedOrder.address?.name }}</strong>
                <span>{{ selectedOrder.address?.phone }}</span>
              </div>
              <div class="text-gray-600">
                {{ selectedOrder.address?.province }}{{ selectedOrder.address?.city }}{{ selectedOrder.address?.district }}{{ selectedOrder.address?.detailAddress }}
              </div>
            </div>
          </div>

          <!-- Products -->
          <div class="win-groupbox mb-3">
            <div class="win-groupbox-title flex items-center gap-1">
              <img :src="icons.product" style="width:16px;height:16px" />
              商品列表
            </div>
            <div class="win-listview" style="max-height:150px">
              <div v-for="item in selectedOrder.skuDetails" :key="item.skuId" class="win-listview-item">
                <img :src="item.goodsThumb" style="width:40px;height:40px;object-fit:cover" class="border border-gray-400 mr-2" />
                <div class="flex-1">
                  <div class="text-sm">{{ item.goodsName }}</div>
                  <div class="text-xs text-gray-500">{{ item.specInfo }}</div>
                </div>
                <div class="text-right">
                  <div class="text-red-600">¥{{ item.goodsPrice }}</div>
                  <div class="text-xs">x{{ item.quantity }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Price Summary -->
          <div class="win-panel p-2">
            <table class="text-sm w-full">
              <tr>
                <td class="text-gray-600">商品总额:</td>
                <td class="text-right">¥{{ selectedOrder.totalPrice }}</td>
              </tr>
              <tr>
                <td class="text-gray-600">运费:</td>
                <td class="text-right">¥{{ selectedOrder.freightFee || 0 }}</td>
              </tr>
              <tr class="font-bold">
                <td>实付金额:</td>
                <td class="text-right text-red-600">¥{{ selectedOrder.finalPrice }}</td>
              </tr>
            </table>
          </div>

          <!-- Buttons -->
          <div class="flex justify-end gap-2 mt-3">
            <button class="win-btn" @click="selectedOrder = null">
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

const orders = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const search = ref('')
const statusFilter = ref('')
const selectedOrder = ref(null)

let debounceTimer = null
const debounceFetch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    fetchOrders()
  }, 500)
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

const formatTime = (time) => {
  if (!time) return ''
  // 如果是秒级时间戳，转为毫秒
  const timestamp = time < 10000000000 ? time * 1000 : time
  return new Date(timestamp).toLocaleString('zh-CN')
}

const fetchOrders = async () => {
  try {
    const params = {
      page: page.value,
      pageSize,
    }
    if (search.value) params.orderNo = search.value
    if (statusFilter.value) params.status = statusFilter.value

    const res = await api.get('/orders', { params })
    const orderList = res.data?.list || []
    
    // 转换订单列表数据格式
    orders.value = orderList.map(order => ({
      ...order,
      skuDetails: order.items?.map(item => ({
        skuId: item._id || item.skuId,
        goodsName: item.sku?.spu?.name || item.sku?.description || '商品',
        goodsThumb: item.sku?.image || item.sku?.spu?.cover_image || '',
        goodsPrice: item.price || 0,
        quantity: item.quantity,
        specInfo: item.sku?.description || '',
      })) || [],
    }))
    total.value = res.data?.total || 0
  } catch (e) {
    console.error('Failed to fetch orders:', e)
  }
}

const viewOrder = async (order) => {
  try {
    const res = await api.get(`/orders/${order._id}`)
    const orderData = res.data?.order || res.data
    
    // 转换 items 为前端期望的 skuDetails 格式
    if (orderData.items) {
      orderData.skuDetails = orderData.items.map(item => ({
        skuId: item._id || item.skuId,
        goodsName: item.sku?.spu?.name || item.sku?.description || '商品',
        goodsThumb: item.sku?.image || item.sku?.spu?.cover_image || '',
        specInfo: item.sku?.description || '',
        goodsPrice: item.price || item.sku?.price || 0,
        quantity: item.quantity,
      }))
    }
    
    // 转换 delivery_info 为 address 格式
    if (orderData.delivery_info) {
      const di = typeof orderData.delivery_info === 'string' 
        ? JSON.parse(orderData.delivery_info)
        : orderData.delivery_info
      orderData.address = {
        name: di.name || '',
        phone: di.phone || '',
        province: di.provinceName || '',
        city: di.cityName || '',
        district: di.districtName || '',
        detailAddress: di.detailAddress || '',
      }
    }
    
    // 合并用户信息
    if (res.data?.user) {
      orderData.user = res.data.user
    }
    
    selectedOrder.value = orderData
  } catch (e) {
    console.error('Failed to fetch order detail:', e)
  }
}

const shipOrder = async (order) => {
  if (!confirm('确定要发货吗？')) return
  try {
    await api.put(`/orders/${order._id}/ship`)
    alert('发货成功')
    fetchOrders()
  } catch (e) {
    alert('操作失败: ' + (e.error || e.message))
  }
}

const refundOrder = async (order) => {
  if (!confirm('确定要退款吗？')) return
  try {
    await api.put(`/orders/${order._id}/refund`)
    alert('退款成功')
    fetchOrders()
  } catch (e) {
    alert('操作失败: ' + (e.error || e.message))
  }
}

onMounted(fetchOrders)
</script>
