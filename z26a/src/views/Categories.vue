<template>
  <div class="p-3 h-full flex flex-col">
    <!-- Toolbar -->
    <div class="win-toolbar mb-2">
      <button class="win-toolbar-btn" @click="showModal()">
        <img :src="icons.add" /> 新建分类
      </button>
      <div class="win-toolbar-separator"></div>
      <button class="win-toolbar-btn" @click="fetchCategories">
        <img :src="icons.refresh" /> 刷新
      </button>
    </div>

    <!-- Categories ListView -->
    <div class="win-listview flex-1">
      <div class="win-listview-header">
        <div style="width:50px"></div>
        <div style="width:150px">分类名称</div>
        <div style="width:80px">排序</div>
        <div style="width:80px">商品数</div>
        <div style="width:150px">创建时间</div>
        <div style="flex:1">操作</div>
      </div>
      
      <div 
        v-for="cat in categories" 
        :key="cat._id" 
        class="win-listview-item"
        @dblclick="showModal(cat)"
      >
        <div style="width:50px">
          <img 
            :src="cat.icon || icons.category" 
            style="width:32px;height:32px;object-fit:cover" 
            class="border border-gray-400"
          />
        </div>
        <div style="width:150px" class="font-medium">
          {{ cat.name }}
        </div>
        <div style="width:80px">{{ cat.sort || 0 }}</div>
        <div style="width:80px">{{ cat.productCount || 0 }}</div>
        <div style="width:150px" class="text-xs text-gray-600">
          {{ formatTime(cat.createdAt) }}
        </div>
        <div style="flex:1" class="flex gap-1">
          <button class="win-btn win-btn-sm" @click.stop="showModal(cat)">
            <img :src="icons.edit" style="width:12px;height:12px" /> 编辑
          </button>
          <button class="win-btn win-btn-sm" @click.stop="deleteCategory(cat)">
            <img :src="icons.delete" style="width:12px;height:12px" /> 删除
          </button>
        </div>
      </div>

      <div v-if="categories.length === 0" class="text-center py-8 text-gray-500">
        <img :src="icons.category32" style="width:48px;height:48px;margin:0 auto 8px" />
        <div>暂无分类数据</div>
      </div>
    </div>

    <!-- Category Modal -->
    <div v-if="modalVisible" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="win-window" style="width: 380px;">
        <div class="win-titlebar">
          <img :src="icons.category" class="win-titlebar-icon" />
          <span>{{ editingId ? '编辑分类' : '新建分类' }}</span>
          <div class="win-titlebar-buttons">
            <button class="win-titlebar-btn" @click="closeModal">×</button>
          </div>
        </div>
        
        <div class="bg-[#c0c0c0] p-4">
          <form @submit.prevent="saveCategory">
            <table class="w-full">
              <tr>
                <td class="py-2 pr-2 text-right whitespace-nowrap">
                  <label>分类名称(<u>N</u>):</label>
                </td>
                <td class="py-2">
                  <input v-model="form.name" type="text" class="win-input w-full" required />
                </td>
              </tr>
              <tr>
                <td class="py-2 pr-2 text-right whitespace-nowrap">
                  <label>图标URL:</label>
                </td>
                <td class="py-2">
                  <input v-model="form.icon" type="text" class="win-input w-full" placeholder="https://..." />
                </td>
              </tr>
              <tr>
                <td class="py-2 pr-2 text-right whitespace-nowrap">
                  <label>排序:</label>
                </td>
                <td class="py-2">
                  <input v-model.number="form.sort" type="number" class="win-input w-full" />
                </td>
              </tr>
            </table>

            <div class="win-groupbox mt-3 p-2">
              <div class="win-groupbox-title">提示</div>
              <p class="text-xs text-gray-600">排序值越大，分类越靠前显示</p>
            </div>

            <div class="flex justify-end gap-2 mt-4 pt-3 border-t border-gray-400">
              <button type="submit" class="win-btn" style="min-width:80px">
                <img :src="icons.save" style="width:16px;height:16px" /> 保存
              </button>
              <button type="button" class="win-btn" style="min-width:80px" @click="closeModal">
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
import { ref, reactive, onMounted } from 'vue'
import api from '@/api'
import icons from '@/assets/icons'

const categories = ref([])
const modalVisible = ref(false)
const editingId = ref(null)

const form = reactive({
  name: '',
  icon: '',
  sort: 0,
})

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const fetchCategories = async () => {
  try {
    const res = await api.get('/categories')
    categories.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch categories:', e)
  }
}

const showModal = (cat = null) => {
  if (cat) {
    editingId.value = cat._id
    form.name = cat.name
    form.icon = cat.icon || ''
    form.sort = cat.sort || 0
  } else {
    editingId.value = null
    form.name = ''
    form.icon = ''
    form.sort = 0
  }
  modalVisible.value = true
}

const closeModal = () => {
  modalVisible.value = false
  editingId.value = null
}

const saveCategory = async () => {
  try {
    if (editingId.value) {
      await api.put(`/categories/${editingId.value}`, form)
    } else {
      await api.post('/categories', form)
    }
    closeModal()
    fetchCategories()
  } catch (e) {
    alert('保存失败: ' + (e.error || e.message))
  }
}

const deleteCategory = async (cat) => {
  if (!confirm(`确定要删除分类 "${cat.name}" 吗？`)) return
  try {
    await api.delete(`/categories/${cat._id}`)
    fetchCategories()
  } catch (e) {
    alert('删除失败: ' + (e.error || e.message))
  }
}

onMounted(fetchCategories)
</script>
