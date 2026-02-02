<template>
  <div class="p-3 h-full flex flex-col overflow-auto">
    <!-- Toolbar -->
    <div class="win-toolbar mb-2">
      <button class="win-toolbar-btn" @click="saveAll" :disabled="saving">
        <img :src="icons.save" /> {{ saving ? '保存中...' : '保存所有配置' }}
      </button>
      <button class="win-toolbar-btn" @click="fetchAll">
        <img :src="icons.refresh" /> 刷新
      </button>
    </div>

    <div class="flex-1 overflow-auto space-y-3">
      <!-- 轮播图管理 -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.image" style="width:16px;height:16px" />
          轮播图管理
        </div>
        <div class="p-3">
          <div class="win-toolbar mb-2">
            <button class="win-toolbar-btn" @click="addBanner">
              <img :src="icons.add" /> 添加轮播图
            </button>
          </div>

          <div class="win-listview" style="max-height: 300px">
            <div class="win-listview-header">
              <div style="width: 80px">预览</div>
              <div style="width: 200px">标题</div>
              <div style="flex: 1">跳转链接</div>
              <div style="width: 60px">排序</div>
              <div style="width: 60px">状态</div>
              <div style="width: 120px">操作</div>
            </div>
            <draggable
              v-model="banners"
              item-key="_id"
              handle=".drag-handle"
              @end="onBannerDragEnd"
            >
              <template #item="{ element: banner, index }">
                <div class="win-listview-item">
                  <div style="width: 80px">
                    <img
                      v-if="banner.imageUrl"
                      :src="banner.imageUrl"
                      style="width: 60px; height: 30px; object-fit: cover; border: 1px solid #808080"
                    />
                    <span v-else class="text-gray-400 text-xs">无图片</span>
                  </div>
                  <div style="width: 200px" class="truncate">{{ banner.title || '未命名' }}</div>
                  <div style="flex: 1" class="truncate text-xs text-gray-500">
                    {{ banner.linkUrl || '无链接' }}
                  </div>
                  <div style="width: 60px" class="drag-handle cursor-move">
                    {{ banner.sort || index + 1 }}
                  </div>
                  <div style="width: 60px">
                    <span
                      :class="banner.enabled ? 'text-green-600' : 'text-gray-400'"
                    >
                      {{ banner.enabled ? '启用' : '禁用' }}
                    </span>
                  </div>
                  <div style="width: 120px" class="flex gap-1">
                    <button class="win-btn win-btn-sm" @click="editBanner(banner, index)">
                      编辑
                    </button>
                    <button class="win-btn win-btn-sm" @click="deleteBanner(index)">
                      删除
                    </button>
                  </div>
                </div>
              </template>
            </draggable>
            <div v-if="banners.length === 0" class="text-center py-4 text-gray-500">
              暂无轮播图，点击上方按钮添加
            </div>
          </div>
        </div>
      </div>

      <!-- 推荐商品管理 -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.product" style="width:16px;height:16px" />
          推荐商品管理
        </div>
        <div class="p-3">
          <div class="win-toolbar mb-2">
            <button class="win-toolbar-btn" @click="openProductSelector">
              <img :src="icons.add" /> 添加推荐商品
            </button>
          </div>

          <div class="win-listview" style="max-height: 300px">
            <div class="win-listview-header">
              <div style="width: 60px">图片</div>
              <div style="flex: 1">商品名称</div>
              <div style="width: 80px">价格</div>
              <div style="width: 80px">标签</div>
              <div style="width: 60px">排序</div>
              <div style="width: 100px">操作</div>
            </div>
            <draggable
              v-model="recommendedProducts"
              item-key="productId"
              handle=".drag-handle"
              @end="onProductDragEnd"
            >
              <template #item="{ element: item, index }">
                <div class="win-listview-item">
                  <div style="width: 60px">
                    <img
                      v-if="item.product?.cover_image"
                      :src="item.product.cover_image"
                      style="width: 40px; height: 40px; object-fit: cover; border: 1px solid #808080"
                    />
                  </div>
                  <div style="flex: 1" class="truncate">
                    {{ item.product?.name || item.productId }}
                  </div>
                  <div style="width: 80px" class="text-red-600">
                    ¥{{ item.product?.min_price || '-' }}
                  </div>
                  <div style="width: 80px">
                    <select
                      v-model="item.tag"
                      class="win-select"
                      style="width: 70px; font-size: 11px"
                    >
                      <option value="">无</option>
                      <option value="hot">热卖</option>
                      <option value="new">新品</option>
                      <option value="sale">特价</option>
                    </select>
                  </div>
                  <div style="width: 60px" class="drag-handle cursor-move">
                    {{ item.sort || index + 1 }}
                  </div>
                  <div style="width: 100px" class="flex gap-1">
                    <button class="win-btn win-btn-sm" @click="removeRecommended(index)">
                      移除
                    </button>
                  </div>
                </div>
              </template>
            </draggable>
            <div v-if="recommendedProducts.length === 0" class="text-center py-4 text-gray-500">
              暂无推荐商品，点击上方按钮添加
            </div>
          </div>
        </div>
      </div>

      <!-- 首页富文本内容 -->
      <div class="win-groupbox">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.file" style="width:16px;height:16px" />
          首页内容编辑
        </div>
        <div class="p-3">
          <div class="mb-2 flex items-center gap-2">
            <label>内容标题:</label>
            <input
              v-model="homeContent.title"
              type="text"
              class="win-input flex-1"
              placeholder="可选，如：活动公告"
            />
            <label class="flex items-center gap-1">
              <input v-model="homeContent.enabled" type="checkbox" />
              <span>启用</span>
            </label>
          </div>
          
          <!-- 富文本工具栏扩展：插入商品按钮 -->
          <div class="mb-2">
            <button class="win-btn" @click="openProductSelectorForRichText">
              <img :src="icons.product" style="width:16px;height:16px" /> 插入商品
            </button>
            <span class="text-xs text-gray-500 ml-2">在富文本中插入商品卡片</span>
          </div>

          <RichTextEditor
            ref="homeContentEditorRef"
            v-model="homeContent.content"
            placeholder="编辑首页展示内容，支持图文混排..."
          />

          <div class="mt-2 flex justify-end">
            <button class="win-btn" @click="saveHomeContent">
              <img :src="icons.save" style="width:16px;height:16px" /> 保存内容
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 轮播图编辑弹窗 -->
    <div v-if="showBannerModal" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="win-window" style="width: 500px">
        <div class="win-titlebar">
          <img :src="icons.image" class="win-titlebar-icon" />
          <span>{{ editingBannerIndex >= 0 ? '编辑轮播图' : '添加轮播图' }}</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn" @click="closeBannerModal">×</button>
          </div>
        </div>

        <div class="bg-[#c0c0c0] p-4">
          <table class="w-full">
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap w-20">
                <label>标题:</label>
              </td>
              <td class="py-2">
                <input v-model="bannerForm.title" type="text" class="win-input w-full" />
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap align-top">
                <label>图片:</label>
              </td>
              <td class="py-2">
                <ImageUploader v-model="bannerForm.images" :max-count="1" />
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap">
                <label>跳转链接:</label>
              </td>
              <td class="py-2">
                <input
                  v-model="bannerForm.linkUrl"
                  type="text"
                  class="win-input w-full"
                  placeholder="/pages/goods/details/index?spuId=xxx"
                />
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap">
                <label>排序:</label>
              </td>
              <td class="py-2">
                <input
                  v-model.number="bannerForm.sort"
                  type="number"
                  class="win-input"
                  style="width: 80px"
                />
                <span class="text-xs text-gray-500 ml-2">数字越小越靠前</span>
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap">
                <label>状态:</label>
              </td>
              <td class="py-2">
                <label class="flex items-center gap-2">
                  <input v-model="bannerForm.enabled" type="checkbox" />
                  <span>启用</span>
                </label>
              </td>
            </tr>
          </table>

          <div class="flex justify-end gap-2 mt-4 pt-3 border-t border-gray-400">
            <button class="win-btn" style="min-width: 80px" @click="saveBanner">
              <img :src="icons.save" style="width: 16px; height: 16px" /> 确定
            </button>
            <button class="win-btn" style="min-width: 80px" @click="closeBannerModal">
              取消
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 商品选择弹窗 -->
    <div v-if="showProductSelector" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="win-window" style="width: 600px; max-height: 80vh">
        <div class="win-titlebar">
          <img :src="icons.product" class="win-titlebar-icon" />
          <span>选择商品</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn" @click="closeProductSelector">×</button>
          </div>
        </div>

        <div class="bg-[#c0c0c0] p-4">
          <!-- 搜索 -->
          <div class="mb-3">
            <input
              v-model="productSearch"
              type="text"
              class="win-input w-full"
              placeholder="输入商品名称搜索..."
              @input="searchProducts"
            />
            <div class="text-xs text-gray-500 mt-1">输入关键词后自动搜索</div>
          </div>

          <!-- 商品列表 -->
          <div class="win-listview" style="height: 300px; overflow-y: auto">
            <div class="win-listview-header">
              <div style="width: 40px"></div>
              <div style="width: 50px">图片</div>
              <div style="flex: 1">商品名称</div>
              <div style="width: 80px">价格</div>
            </div>
            <div
              v-for="product in availableProducts"
              :key="product._id"
              class="win-listview-item cursor-pointer"
              :class="{ 'bg-[#000080] text-white': selectedProducts.includes(product._id) }"
              @click="toggleProductSelection(product._id)"
            >
              <div style="width: 40px">
                <input
                  type="checkbox"
                  :checked="selectedProducts.includes(product._id)"
                  @click.stop
                  @change="toggleProductSelection(product._id)"
                />
              </div>
              <div style="width: 50px">
                <img
                  v-if="product.cover_image"
                  :src="product.cover_image"
                  style="width: 32px; height: 32px; object-fit: cover"
                />
              </div>
              <div style="flex: 1" class="truncate">{{ product.name }}</div>
              <div style="width: 80px">¥{{ product.min_price || '-' }}</div>
            </div>
            <div v-if="availableProducts.length === 0" class="text-center py-4 text-gray-500">
              {{ productSearch ? '未找到匹配的商品' : '请输入商品名称搜索' }}
            </div>
          </div>

          <div class="flex justify-between items-center mt-4 pt-3 border-t border-gray-400">
            <span class="text-xs text-gray-600">已选择 {{ selectedProducts.length }} 个商品</span>
            <div class="flex gap-2">
              <button class="win-btn" style="min-width: 80px" @click="confirmProductSelection">
                确定
              </button>
              <button class="win-btn" style="min-width: 80px" @click="closeProductSelector">
                取消
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import draggable from 'vuedraggable'
import api from '@/api'
import icons from '@/assets/icons'
import ImageUploader from '@/components/ImageUploader.vue'
import RichTextEditor from '@/components/RichTextEditor.vue'

// 状态
const loading = ref(false)
const saving = ref(false)

// 轮播图
const banners = ref([])
const showBannerModal = ref(false)
const editingBannerIndex = ref(-1)
const bannerForm = ref({
  title: '',
  images: [],
  linkUrl: '',
  sort: 1,
  enabled: true,
})

// 推荐商品
const recommendedProducts = ref([])
const showProductSelector = ref(false)
const productSearch = ref('')
const allProducts = ref([])
const selectedProducts = ref([])
const productSelectorMode = ref('recommend') // 'recommend' or 'richtext'

// 首页富文本内容
const homeContentEditorRef = ref(null)
const homeContent = ref({
  key: 'main',
  title: '',
  content: '',
  enabled: true,
  priority: 0,
})

// 搜索结果（后端已排除已推荐的商品）
const availableProducts = computed(() => allProducts.value)

// 获取所有配置
const fetchAll = async () => {
  loading.value = true
  try {
    // 获取轮播图配置
    const bannerRes = await api.get('/home-config/banners')
    // 转换后端数据格式
    banners.value = (bannerRes.data || []).map(b => ({
      _id: b.id,
      title: b.title,
      imageUrl: b.image,
      linkUrl: b.link,
      sort: b.priority,
      enabled: true,
    }))

    // 获取推荐商品配置
    const recommendRes = await api.get('/home-config/recommended')
    recommendedProducts.value = (recommendRes.data || []).map(p => ({
      productId: p.productId,
      product: {
        _id: p.productId,
        name: p.productName,
        cover_image: p.productImage,
        min_price: p.productPrice,
      },
      tag: (p.tags && p.tags[0]) || '',
      sort: p.priority,
    }))

    // 获取所有商品（用于选择）
    const productsRes = await api.get('/products', { params: { pageSize: 1000 } })
    allProducts.value = (productsRes.data?.list || productsRes.data || []).map(p => ({
      _id: p._id,
      name: p.name,
      cover_image: p.cover_image,
      min_price: p.min_price || p.price,
    }))

    // 获取首页富文本内容
    try {
      const contentRes = await api.get('/home-config/contents/main')
      if (contentRes._id) {
        homeContent.value = {
          key: contentRes.key || 'main',
          title: contentRes.title || '',
          content: contentRes.content || '',
          enabled: contentRes.enabled !== false,
          priority: contentRes.priority || 0,
        }
      }
    } catch (e) {
      // 忽略，使用默认空值
    }
  } catch (error) {
    console.error('获取配置失败:', error)
  } finally {
    loading.value = false
  }
}

// 保存所有配置
const saveAll = async () => {
  saving.value = true
  try {
    // 分别保存轮播图
    for (const banner of banners.value) {
      const data = {
        image: banner.imageUrl,
        title: banner.title,
        link: banner.linkUrl,
        priority: banner.sort,
      }
      if (banner._id) {
        await api.put(`/home-config/banners/${banner._id}`, data)
      } else {
        const res = await api.post('/home-config/banners', data)
        banner._id = res.data._id
      }
    }

    // 更新轮播图排序
    const bannerIds = banners.value.filter(b => b._id).map(b => b._id)
    if (bannerIds.length > 0) {
      await api.post('/home-config/banners/reorder', { ids: bannerIds })
    }

    // 保存推荐商品
    // 先获取当前已有的推荐商品，对比添加/删除
    const currentRes = await api.get('/home-config/recommended')
    const currentIds = (currentRes.data || []).map(p => p.productId)
    const newIds = recommendedProducts.value.map(p => p.productId)
    
    // 添加新的推荐
    for (const product of recommendedProducts.value) {
      if (!currentIds.includes(product.productId)) {
        await api.post('/home-config/recommended', {
          productId: product.productId,
          tags: product.tag ? [product.tag] : [],
        })
      } else {
        // 更新标签
        await api.put(`/home-config/recommended/${product.productId}`, {
          tags: product.tag ? [product.tag] : [],
        })
      }
    }

    // 移除取消的推荐
    for (const id of currentIds) {
      if (!newIds.includes(id)) {
        await api.delete(`/home-config/recommended/${id}`)
      }
    }

    // 更新排序
    if (newIds.length > 0) {
      await api.post('/home-config/recommended/reorder', { ids: newIds })
    }

    alert('保存成功！')
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败: ' + (error.error || error.message))
  } finally {
    saving.value = false
  }
}

// 轮播图操作
const addBanner = () => {
  editingBannerIndex.value = -1
  bannerForm.value = {
    title: '',
    images: [],
    linkUrl: '',
    sort: banners.value.length + 1,
    enabled: true,
  }
  showBannerModal.value = true
}

const editBanner = (banner, index) => {
  editingBannerIndex.value = index
  bannerForm.value = {
    title: banner.title || '',
    images: banner.imageUrl ? [banner.imageUrl] : [],
    linkUrl: banner.linkUrl || '',
    sort: banner.sort || index + 1,
    enabled: banner.enabled !== false,
  }
  showBannerModal.value = true
}

const saveBanner = () => {
  const banner = {
    title: bannerForm.value.title,
    imageUrl: bannerForm.value.images[0] || '',
    linkUrl: bannerForm.value.linkUrl,
    sort: bannerForm.value.sort,
    enabled: bannerForm.value.enabled,
  }

  if (editingBannerIndex.value >= 0) {
    banners.value[editingBannerIndex.value] = { ...banners.value[editingBannerIndex.value], ...banner }
  } else {
    banners.value.push(banner)
  }

  closeBannerModal()
}

const deleteBanner = (index) => {
  if (confirm('确定删除这个轮播图吗？')) {
    banners.value.splice(index, 1)
  }
}

const closeBannerModal = () => {
  showBannerModal.value = false
}

const onBannerDragEnd = () => {
  banners.value.forEach((b, i) => {
    b.sort = i + 1
  })
}

// 推荐商品操作
const closeProductSelector = () => {
  showProductSelector.value = false
}

const toggleProductSelection = (productId) => {
  const index = selectedProducts.value.indexOf(productId)
  if (index >= 0) {
    selectedProducts.value.splice(index, 1)
  } else {
    selectedProducts.value.push(productId)
  }
}

const removeRecommended = (index) => {
  recommendedProducts.value.splice(index, 1)
}

const onProductDragEnd = () => {
  recommendedProducts.value.forEach((p, i) => {
    p.sort = i + 1
  })
}

// 搜索防抖计时器
let searchTimer = null

const searchProducts = async () => {
  // 防抖处理
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  
  searchTimer = setTimeout(async () => {
    try {
      const params = {}
      if (productSearch.value) {
        params.keyword = productSearch.value
      }
      const products = await api.get('/home-config/products/search', { params })
      // api拦截器已经返回 response.data，所以 products 直接就是数组
      allProducts.value = (products || []).map(p => ({
        _id: p._id,
        name: p.name,
        cover_image: p.coverImage || p.cover_image,
        min_price: p.price || p.min_price,
      }))
    } catch (error) {
      console.error('搜索商品失败:', error)
    }
  }, 300)
}

// ============================================
// 首页富文本内容相关
// ============================================

// 打开商品选择器（用于插入富文本）
const openProductSelectorForRichText = () => {
  productSelectorMode.value = 'richtext'
  selectedProducts.value = []
  productSearch.value = ''
  allProducts.value = []
  showProductSelector.value = true
}

// 修改打开商品选择器函数，支持两种模式
const openProductSelector = () => {
  productSelectorMode.value = 'recommend'
  selectedProducts.value = []
  productSearch.value = ''
  allProducts.value = []
  showProductSelector.value = true
}

// 修改确认选择函数，支持两种模式
const confirmProductSelection = () => {
  if (productSelectorMode.value === 'richtext') {
    // 插入商品到富文本
    insertProductsToRichText()
  } else {
    // 添加推荐商品
    confirmRecommendedProducts()
  }
  closeProductSelector()
}

// 原有的确认推荐商品函数
const confirmRecommendedProducts = () => {
  selectedProducts.value.forEach(productId => {
    const product = allProducts.value.find(p => p._id === productId)
    if (product && !recommendedProducts.value.find(r => r.productId === productId)) {
      recommendedProducts.value.push({
        productId: product._id,
        product: product,
        tag: '',
        sort: recommendedProducts.value.length + 1,
      })
    }
  })
}

// 插入商品到富文本
const insertProductsToRichText = () => {
  if (!selectedProducts.value.length) {
    alert('请先选择要插入的商品')
    return
  }
  
  // 生成商品占位符标记
  const productsMarkers = selectedProducts.value.map(productId => {
    const product = allProducts.value.find(p => p._id === productId)
    if (!product) return ''
    // 格式: [商品:商品ID:商品名称] - 名称仅供管理端显示参考
    return `[商品:${product._id}:${product.name}]`
  }).join(' ')
  
  // 使用编辑器API插入
  if (productsMarkers && homeContentEditorRef.value) {
    const editor = homeContentEditorRef.value.getEditor?.()
    if (editor) {
      editor.focus()
      editor.insertText(productsMarkers)
    } else {
      homeContent.value.content += productsMarkers
    }
  } else {
    homeContent.value.content += productsMarkers
  }
  
  // 关闭选择对话框
  showProductSelector.value = false
  selectedProducts.value = []
}

// 保存首页富文本内容
const saveHomeContent = async () => {
  try {
    await api.post('/home-config/contents', {
      key: homeContent.value.key || 'main',
      title: homeContent.value.title,
      content: homeContent.value.content,
      enabled: homeContent.value.enabled,
      priority: homeContent.value.priority,
    })
    alert('首页内容保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败: ' + (error.error || error.message))
  }
}

onMounted(() => {
  fetchAll()
})
</script>
