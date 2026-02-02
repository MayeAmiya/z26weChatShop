<template>
  <div class="p-3 h-full flex flex-col">
    <!-- Toolbar -->
    <div class="win-toolbar mb-2">
      <button class="win-toolbar-btn" @click="$router.push('/products/new')">
        <img :src="icons.add" /> 新建商品
      </button>
      <div class="win-toolbar-separator"></div>
      <button class="win-toolbar-btn" @click="fetchProducts">
        <img :src="icons.refresh" /> 刷新
      </button>
      <button class="win-toolbar-btn" @click="showTagManager = true">
        <img :src="icons.category" /> 管理标签
      </button>
      <div class="flex-1"></div>
      <div class="flex items-center gap-2">
        <img :src="icons.search" style="width:16px;height:16px" />
        <input
          v-model="search"
          type="text"
          placeholder="搜索商品..."
          class="win-input"
          style="width: 150px"
          @input="debounceFetch"
        />
        <select v-model="categoryFilter" class="win-select" @change="fetchProducts">
          <option value="">全部分类</option>
          <option v-for="cat in categories" :key="cat._id" :value="cat._id">{{ cat.name }}</option>
        </select>
        <select v-model="statusFilter" class="win-select" @change="fetchProducts">
          <option value="">全部状态</option>
          <option value="ENABLED">上架中</option>
          <option value="DISABLED">已下架</option>
        </select>
      </div>
    </div>

    <!-- Products ListView -->
    <div class="win-listview flex-1">
      <div class="win-listview-header">
        <div style="width:50px"></div>
        <div style="width:180px">商品名称</div>
        <div style="width:100px">分类</div>
        <div style="width:150px">标签</div>
        <div style="width:80px">状态</div>
        <div style="flex:1">操作</div>
      </div>
      
      <div 
        v-for="product in products" 
        :key="product._id" 
        class="win-listview-item"
        @dblclick="$router.push(`/products/${product._id}`)"
      >
        <div style="width:50px">
          <img 
            :src="product.cover_image || icons.image" 
            style="width:32px;height:32px;object-fit:cover" 
            class="border border-gray-400"
          />
        </div>
        <div style="width:180px">
          <div class="font-medium truncate" :title="product.name">{{ product.name }}</div>
          <div class="text-xs text-gray-500">ID: {{ product._id?.slice(-8) }}</div>
        </div>
        <div style="width:100px">
          <span class="win-badge win-badge-blue" v-if="product.category">
            {{ product.category.name }}
          </span>
          <span v-else class="text-gray-400">-</span>
        </div>
        <div style="width:150px" class="flex flex-wrap gap-1">
          <span 
            v-for="tag in product.tags?.slice(0, 3)" 
            :key="tag._id" 
            class="win-tag"
            :style="{ backgroundColor: tag.color || '#e0e0e0' }"
          >
            {{ tag.name }}
          </span>
          <span v-if="product.tags?.length > 3" class="text-xs text-gray-500">+{{ product.tags.length - 3 }}</span>
        </div>
        <div style="width:80px">
          <span :class="product.status === 'ENABLED' ? 'win-badge win-badge-green' : 'win-badge win-badge-gray'">
            {{ product.status === 'ENABLED' ? '上架' : '下架' }}
          </span>
        </div>
        <div style="flex:1" class="flex gap-1">
          <button class="win-btn win-btn-sm" @click.stop="$router.push(`/products/${product._id}`)">
            <img :src="icons.edit" style="width:12px;height:12px" /> 编辑
          </button>
          <button 
            class="win-btn win-btn-sm" 
            @click.stop="toggleStatus(product)"
          >
            <img :src="product.status === 'ENABLED' ? icons.cross : icons.check" style="width:12px;height:12px" />
            {{ product.status === 'ENABLED' ? '下架' : '上架' }}
          </button>
          <button class="win-btn win-btn-sm" @click.stop="deleteProduct(product)">
            <img :src="icons.delete" style="width:12px;height:12px" /> 删除
          </button>
        </div>
      </div>

      <div v-if="products.length === 0" class="text-center py-8 text-gray-500">
        <img :src="icons.folder32" style="width:48px;height:48px;margin:0 auto 8px" />
        <div>暂无商品数据</div>
      </div>
    </div>

    <!-- Pagination -->
    <div class="win-toolbar mt-2 justify-between">
      <div class="text-sm flex items-center gap-1">
        <img :src="icons.info" style="width:16px;height:16px" />
        共 {{ total }} 个商品
      </div>
      <div class="flex gap-1">
        <button class="win-btn win-btn-sm" :disabled="page === 1" @click="page--; fetchProducts()">
          &lt; 上一页
        </button>
        <span class="win-panel px-3 py-1 text-sm">
          第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) || 1 }} 页
        </span>
        <button 
          class="win-btn win-btn-sm" 
          :disabled="page >= Math.ceil(total / pageSize)" 
          @click="page++; fetchProducts()"
        >
          下一页 &gt;
        </button>
      </div>
    </div>

    <!-- Tag Manager Dialog -->
    <div v-if="showTagManager" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="win-window" style="width: 450px; max-height: 80vh;">
        <div class="win-titlebar">
          <img :src="icons.category" class="win-titlebar-icon" />
          <span>标签管理</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn" @click="showTagManager = false">×</button>
          </div>
        </div>
        <div class="bg-[#c0c0c0] p-3">
          <div class="win-toolbar mb-2">
            <input v-model="newTagName" class="win-input" placeholder="新标签名称" style="width:120px" />
            <input v-model="newTagColor" type="color" class="win-input" style="width:40px;padding:1px" />
            <button class="win-btn win-btn-sm" @click="createTag">
              <img :src="icons.add" style="width:12px;height:12px" /> 添加
            </button>
          </div>
          <div class="win-listview" style="height:200px">
            <div v-for="tag in tags" :key="tag._id" class="win-listview-item">
              <span class="win-tag" :style="{ backgroundColor: tag.color || '#e0e0e0' }">{{ tag.name }}</span>
              <span class="flex-1"></span>
              <button class="win-btn win-btn-sm" @click="deleteTag(tag)">
                <img :src="icons.delete" style="width:12px;height:12px" />
              </button>
            </div>
            <div v-if="tags.length === 0" class="text-center py-4 text-gray-500">
              暂无标签
            </div>
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

const products = ref([])
const categories = ref([])
const tags = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const search = ref('')
const statusFilter = ref('')
const categoryFilter = ref('')

const showTagManager = ref(false)
const newTagName = ref('')
const newTagColor = ref('#3b82f6')

let debounceTimer = null
const debounceFetch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    fetchProducts()
  }, 500)
}

const fetchProducts = async () => {
  try {
    const params = {
      page: page.value,
      pageSize,
    }
    if (search.value) params.keyword = search.value
    if (statusFilter.value) params.status = statusFilter.value
    if (categoryFilter.value) params.categoryId = categoryFilter.value

    const res = await api.get('/products', { params })
    products.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error('Failed to fetch products:', e)
  }
}

const fetchCategories = async () => {
  try {
    const res = await api.get('/categories')
    categories.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch categories:', e)
  }
}

const fetchTags = async () => {
  try {
    const res = await api.get('/tags')
    tags.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch tags:', e)
  }
}

const createTag = async () => {
  if (!newTagName.value.trim()) return
  try {
    await api.post('/tags', { name: newTagName.value.trim(), color: newTagColor.value })
    newTagName.value = ''
    fetchTags()
  } catch (e) {
    alert('创建失败: ' + (e.error || e.message))
  }
}

const deleteTag = async (tag) => {
  if (!confirm(`确定要删除标签 "${tag.name}" 吗？`)) return
  try {
    await api.delete(`/tags/${tag._id}`)
    fetchTags()
  } catch (e) {
    alert('删除失败: ' + (e.error || e.message))
  }
}

const toggleStatus = async (product) => {
  try {
    await api.put(`/products/${product._id}/toggle-status`)
    fetchProducts()
  } catch (e) {
    alert('操作失败: ' + (e.error || e.message))
  }
}

const deleteProduct = async (product) => {
  if (!confirm(`确定要删除商品 "${product.name}" 吗？`)) return
  try {
    await api.delete(`/products/${product._id}`)
    fetchProducts()
  } catch (e) {
    alert('删除失败: ' + (e.error || e.message))
  }
}

onMounted(() => {
  fetchCategories()
  fetchTags()
  fetchProducts()
})
</script>
