<template>
  <div class="image-uploader">
    <!-- 已上传的图片列表 -->
    <div class="flex flex-wrap gap-2 mb-2">
      <div
        v-for="(url, index) in modelValue"
        :key="index"
        class="relative group"
        style="border: 2px solid #808080; border-right-color: #ffffff; border-bottom-color: #ffffff;"
      >
        <img :src="url" style="width:64px;height:64px;object-fit:cover;display:block" />
        <div class="absolute inset-0 bg-black bg-opacity-50 opacity-0 group-hover:opacity-100 flex items-center justify-center">
          <button
            type="button"
            @click="removeImage(index)"
            class="win-btn win-btn-sm"
            style="padding:2px 4px"
          >
            <img :src="icons.delete" style="width:14px;height:14px" />
          </button>
        </div>
      </div>

      <!-- 上传按钮 -->
      <label
        v-if="!maxCount || modelValue.length < maxCount"
        class="flex flex-col items-center justify-center cursor-pointer"
        style="width:64px;height:64px;border:2px dashed #808080;background:#ffffff"
        :class="{ 'opacity-50 cursor-not-allowed': uploading }"
      >
        <img v-if="!uploading" :src="icons.upload" style="width:20px;height:20px" />
        <div v-else class="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-600"></div>
        <span class="text-xs text-gray-600 mt-1">{{ uploading ? '上传中' : '上传' }}</span>
        <input
          type="file"
          :accept="accept"
          :multiple="multiple"
          class="hidden"
          :disabled="uploading"
          @change="handleFileChange"
        />
      </label>
    </div>

    <!-- 提示信息 -->
    <div v-if="tip" class="text-xs text-gray-500">{{ tip }}</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '@/api'
import icons from '@/assets/icons'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => [],
  },
  maxCount: {
    type: Number,
    default: 0, // 0 表示不限制
  },
  multiple: {
    type: Boolean,
    default: false,
  },
  accept: {
    type: String,
    default: 'image/jpeg,image/png,image/gif,image/webp',
  },
  tip: {
    type: String,
    default: '支持 jpg/png/gif/webp 格式，单个文件不超过 5MB',
  },
})

const emit = defineEmits(['update:modelValue'])

const uploading = ref(false)

const handleFileChange = async (e) => {
  const files = e.target.files
  if (!files || files.length === 0) return

  uploading.value = true

  try {
    const formData = new FormData()

    if (props.multiple && files.length > 1) {
      // 批量上传
      for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i])
      }
      const res = await api.post('/upload/images', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      
      if (res.uploaded && res.uploaded.length > 0) {
        const newUrls = res.uploaded.map(item => item.url)
        emit('update:modelValue', [...props.modelValue, ...newUrls])
      }
      
      if (res.errors && res.errors.length > 0) {
        alert('部分图片上传失败：\n' + res.errors.join('\n'))
      }
    } else {
      // 单个上传
      formData.append('file', files[0])
      const res = await api.post('/upload/image', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      
      if (res.url) {
        emit('update:modelValue', [...props.modelValue, res.url])
      }
    }
  } catch (err) {
    alert(err.error || '上传失败')
  } finally {
    uploading.value = false
    // 清空 input，允许重复选择同一文件
    e.target.value = ''
  }
}

const removeImage = (index) => {
  const newValue = [...props.modelValue]
  newValue.splice(index, 1)
  emit('update:modelValue', newValue)
}
</script>
