<template>
  <div class="p-3 h-full flex flex-col overflow-auto">
    <!-- Toolbar -->
    <div class="win-toolbar mb-2">
      <button class="win-toolbar-btn" @click="$router.push('/products')">
        <img :src="icons.folder" /> 返回列表
      </button>
      <div class="win-toolbar-separator"></div>
      <button class="win-toolbar-btn" @click="saveProduct" :disabled="saving">
        <img :src="icons.save" /> {{ saving ? '保存中...' : '保存' }}
      </button>
      <button class="win-toolbar-btn" @click="fetchProduct" v-if="!isNew">
        <img :src="icons.refresh" /> 刷新
      </button>
    </div>

    <div class="flex-1 overflow-auto">
      <!-- Basic Info -->
      <div class="win-groupbox mb-3">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.product" style="width:16px;height:16px" />
          基本信息
        </div>
        <div class="p-3">
          <table class="w-full">
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap w-24">
                <label>商品名称:</label>
              </td>
              <td class="py-2" colspan="3">
                <input v-model="form.title" type="text" class="win-input w-full" required />
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap">
                <label>商品分类:</label>
              </td>
              <td class="py-2" colspan="3">
                <select v-model="form.categoryId" class="win-select" style="width:200px">
                  <option value="">-- 无分类 --</option>
                  <option v-for="cat in categories" :key="cat._id" :value="cat._id">{{ cat.name }}</option>
                </select>
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap align-top">
                <label>商品标签:</label>
              </td>
              <td class="py-2" colspan="3">
                <div class="win-tag-list">
                  <!-- 已添加的标签列表 -->
                  <div v-for="(tag, index) in productTags" :key="tag._id || index" class="win-tag-item">
                    <span class="win-tag-color" :style="{ backgroundColor: tag.color || '#888' }"></span>
                    <input 
                      v-if="tag.editing" 
                      v-model="tag.name" 
                      type="text" 
                      class="win-input flex-1" 
                      @blur="finishEditTag(index)"
                      @keyup.enter="finishEditTag(index)"
                      ref="tagInputRefs"
                      style="font-size:11px;padding:1px 4px;"
                    />
                    <span v-else class="flex-1 cursor-pointer" @dblclick="startEditTag(index)">{{ tag.name }}</span>
                    <input 
                      type="color" 
                      v-model="tag.color" 
                      class="win-color-picker"
                      title="选择颜色"
                    />
                    <button class="win-tag-btn" @click="removeTag(index)" title="删除">×</button>
                  </div>
                  <!-- 添加按钮 -->
                  <div class="win-tag-add" @click="addNewTag">
                    <span>＋ 添加标签</span>
                  </div>
                </div>
                <div class="text-xs text-gray-500 mt-1">双击标签名称可编辑，标签会在保存商品时一起保存</div>
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap align-top">
                <label>商品图片:</label>
              </td>
              <td class="py-2" colspan="3">
                <ImageUploader v-model="form.images" :max-count="9" multiple />
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap align-top">
                <label>商品描述:</label>
              </td>
              <td class="py-2" colspan="3">
                <RichTextEditor 
                  v-model="form.description" 
                  placeholder="请输入商品详细描述..."
                  style="min-height: 300px;"
                />
              </td>
            </tr>
            <tr>
              <td class="py-2 pr-3 text-right whitespace-nowrap">
                <label>状态:</label>
              </td>
              <td class="py-2" colspan="3">
                <label class="flex items-center gap-2">
                  <input v-model="form.status" type="checkbox" :true-value="1" :false-value="0" />
                  <span>{{ form.status === 1 ? '已上架' : '已下架' }}</span>
                </label>
              </td>
            </tr>
          </table>
        </div>
      </div>

      <!-- SKU Management -->
      <div class="win-groupbox" v-if="!isNew">
        <div class="win-groupbox-title flex items-center gap-1">
          <img :src="icons.category" style="width:16px;height:16px" />
          规格管理 (SKU)
        </div>
        <div class="p-3">
          <div class="win-toolbar mb-2">
            <button class="win-toolbar-btn" @click="addSku">
              <img :src="icons.add" /> 添加规格
            </button>
          </div>

          <div class="win-listview" style="max-height:200px">
            <div class="win-listview-header">
              <div style="width:200px">规格名称</div>
              <div style="width:100px">价格</div>
              <div style="width:100px">库存</div>
              <div style="flex:1">操作</div>
            </div>
            <div v-for="sku in skus" :key="sku._id" class="win-listview-item">
              <div style="width:200px">{{ sku.description || '默认规格' }}</div>
              <div style="width:100px" class="text-red-600">¥{{ sku.price }}</div>
              <div style="width:100px">{{ sku.count }}</div>
              <div style="flex:1" class="flex gap-1">
                <button class="win-btn win-btn-sm" @click="editSku(sku)">
                  <img :src="icons.edit" style="width:12px;height:12px" /> 编辑
                </button>
                <button class="win-btn win-btn-sm" @click="deleteSku(sku)">
                  <img :src="icons.delete" style="width:12px;height:12px" /> 删除
                </button>
              </div>
            </div>
            <div v-if="skus.length === 0" class="text-center py-4 text-gray-500">
              暂无规格，请添加
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- SKU Modal -->
    <div v-if="showSkuModal" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="win-window" style="width: 380px;">
        <div class="win-titlebar">
          <img :src="icons.category" class="win-titlebar-icon" />
          <span>{{ editingSkuId ? '编辑规格' : '添加规格' }}</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn" @click="closeSkuModal">×</button>
          </div>
        </div>
        
        <div class="bg-[#c0c0c0] p-4">
          <form @submit.prevent="saveSku">
            <table class="w-full">
              <tr>
                <td class="py-2 pr-2 text-right whitespace-nowrap">
                  <label>规格名称:</label>
                </td>
                <td class="py-2">
                  <input v-model="skuForm.specInfo" type="text" class="win-input w-full" placeholder="如: 红色 / 大号" />
                </td>
              </tr>
              <tr>
                <td class="py-2 pr-2 text-right whitespace-nowrap">
                  <label>价格:</label>
                </td>
                <td class="py-2">
                  <div class="flex items-center gap-1">
                    <span>¥</span>
                    <input v-model.number="skuForm.price" type="number" step="0.01" class="win-input flex-1" required />
                  </div>
                </td>
              </tr>
              <tr>
                <td class="py-2 pr-2 text-right whitespace-nowrap">
                  <label>库存:</label>
                </td>
                <td class="py-2">
                  <input v-model.number="skuForm.stockQuantity" type="number" class="win-input w-full" required />
                </td>
              </tr>
            </table>

            <div class="flex justify-end gap-2 mt-4 pt-3 border-t border-gray-400">
              <button type="submit" class="win-btn" style="min-width:80px">
                <img :src="icons.save" style="width:16px;height:16px" /> 保存
              </button>
              <button type="button" class="win-btn" style="min-width:80px" @click="closeSkuModal">
                <img :src="icons.cross" style="width:16px;height:16px" /> 取消
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import ImageUploader from '@/components/ImageUploader.vue'
import RichTextEditor from '@/components/RichTextEditor.vue'
import icons from '@/assets/icons'

const route = useRoute()
const router = useRouter()

const isNew = computed(() => route.params.id === 'new')
const productId = computed(() => isNew.value ? null : route.params.id)

const categories = ref([])
const productTags = ref([]) // 当前商品的标签列表
const skus = ref([])
const saving = ref(false)

const form = reactive({
  title: '',
  description: '',
  images: [],
  categoryId: '',
  status: 1,
})

// SKU
const showSkuModal = ref(false)
const editingSkuId = ref(null)
const skuForm = reactive({
  specInfo: '',
  price: 0,
  stockQuantity: 0,
})

const fetchProduct = async () => {
  if (isNew.value) return
  try {
    const res = await api.get(`/products/${productId.value}`)
    const { product, skus: skuList } = res.data
    
    // 处理图片数组 - 可能是字符串或数组
    let images = []
    if (product.swiper_images) {
      if (typeof product.swiper_images === 'string') {
        try {
          images = JSON.parse(product.swiper_images)
        } catch {
          images = []
        }
      } else if (Array.isArray(product.swiper_images)) {
        images = product.swiper_images
      }
    }
    
    // 映射后端字段到前端表单字段
    Object.assign(form, {
      title: product.name || '',
      description: product.detail || '',
      images: images,
      categoryId: product.categoryId || '',
      status: product.status === 'ENABLED' ? 1 : 0,
    })
    // 加载标签
    productTags.value = (product.tags || []).map(t => ({
      _id: t._id,
      name: t.name,
      color: t.color || '#888888',
      editing: false
    }))
    skus.value = skuList || []
  } catch (e) {
    console.error('Failed to fetch product:', e)
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

// 标签管理
const addNewTag = () => {
  productTags.value.push({
    _id: null, // 新标签没有ID
    name: '新标签',
    color: '#' + Math.floor(Math.random()*16777215).toString(16).padStart(6, '0'),
    editing: true
  })
  // 等待DOM更新后聚焦
  setTimeout(() => {
    const inputs = document.querySelectorAll('.win-tag-item input[type="text"]')
    if (inputs.length > 0) {
      inputs[inputs.length - 1].focus()
      inputs[inputs.length - 1].select()
    }
  }, 50)
}

const startEditTag = (index) => {
  productTags.value[index].editing = true
  setTimeout(() => {
    const inputs = document.querySelectorAll('.win-tag-item input[type="text"]')
    if (inputs[index]) {
      inputs[index].focus()
      inputs[index].select()
    }
  }, 50)
}

const finishEditTag = (index) => {
  productTags.value[index].editing = false
}

const removeTag = (index) => {
  productTags.value.splice(index, 1)
}

const saveProduct = async () => {
  saving.value = true
  try {
    // 准备标签数据
    const tags = productTags.value.map(t => ({
      _id: t._id,
      name: t.name,
      color: t.color
    }))
    
    // 转换字段名以匹配后端
    const payload = {
      name: form.title,
      detail: form.description,
      images: form.images,
      coverImage: form.images?.[0] || '',
      categoryId: form.categoryId || '',
      status: form.status === 1 ? 'ENABLED' : 'DISABLED',
      tags
    }
    
    if (isNew.value) {
      const res = await api.post('/products', payload)
      router.push(`/products/${res.data._id}`)
    } else {
      await api.put(`/products/${productId.value}`, payload)
      alert('保存成功')
      fetchProduct() // 刷新获取最新标签ID
    }
  } catch (e) {
    alert('保存失败: ' + (e.error || e.message))
  } finally {
    saving.value = false
  }
}

const addSku = () => {
  editingSkuId.value = null
  skuForm.specInfo = ''
  skuForm.price = 0
  skuForm.stockQuantity = 0
  showSkuModal.value = true
}

const editSku = (sku) => {
  editingSkuId.value = sku._id
  skuForm.specInfo = sku.description || ''
  skuForm.price = sku.price || 0
  skuForm.stockQuantity = sku.count || 0
  showSkuModal.value = true
}

const closeSkuModal = () => {
  showSkuModal.value = false
  editingSkuId.value = null
  skuForm.specInfo = ''
  skuForm.price = 0
  skuForm.stockQuantity = 0
}

const saveSku = async () => {
  try {
    // 转换字段名以匹配后端 (description, price, count)
    const data = {
      spuId: productId.value,
      description: skuForm.specInfo,
      price: skuForm.price,
      count: skuForm.stockQuantity,
    }
    if (editingSkuId.value) {
      await api.put(`/skus/${editingSkuId.value}`, data)
    } else {
      await api.post('/skus', data)
    }
    closeSkuModal()
    fetchProduct()
  } catch (e) {
    alert('保存失败: ' + (e.error || e.message))
  }
}

const deleteSku = async (sku) => {
  if (!confirm('确定要删除此规格吗？')) return
  try {
    await api.delete(`/skus/${sku._id}`)
    fetchProduct()
  } catch (e) {
    alert('删除失败: ' + (e.error || e.message))
  }
}

onMounted(() => {
  fetchCategories()
  fetchProduct()
})
</script>
