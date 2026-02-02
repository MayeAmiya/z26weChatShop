<template>
  <div class="file-uploader">
    <!-- 已上传的文件列表 -->
    <div class="win-listview mb-2" v-if="modelValue.length > 0" style="max-height: 150px">
      <div class="win-listview-header">
        <div style="width: 32px"></div>
        <div style="flex: 1">文件名</div>
        <div style="width: 80px">大小</div>
        <div style="width: 60px">操作</div>
      </div>
      <div
        v-for="(file, index) in modelValue"
        :key="index"
        class="win-listview-item"
      >
        <div style="width: 32px">
          <img :src="getFileIcon(file.type)" style="width: 16px; height: 16px" />
        </div>
        <div style="flex: 1" class="truncate">
          <a :href="file.url" target="_blank" class="text-blue-600 hover:underline">
            {{ file.name }}
          </a>
        </div>
        <div style="width: 80px" class="text-gray-500 text-xs">
          {{ formatFileSize(file.size) }}
        </div>
        <div style="width: 60px">
          <button
            type="button"
            class="win-btn win-btn-sm"
            @click="removeFile(index)"
          >
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- 上传按钮 -->
    <label
      class="win-btn inline-flex items-center gap-1 cursor-pointer"
      :class="{ 'opacity-50 cursor-not-allowed': uploading }"
    >
      <img :src="icons.upload" style="width: 16px; height: 16px" />
      <span>{{ uploading ? '上传中...' : '上传文件' }}</span>
      <input
        type="file"
        :accept="accept"
        :multiple="multiple"
        class="hidden"
        :disabled="uploading"
        @change="handleFileChange"
      />
    </label>

    <div class="text-xs text-gray-500 mt-1">{{ tip }}</div>
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
  multiple: {
    type: Boolean,
    default: true,
  },
  accept: {
    type: String,
    default: '.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.zip,.rar',
  },
  tip: {
    type: String,
    default: '支持 PDF、Office 文档、压缩包等，单个文件不超过 50MB',
  },
  maxSize: {
    type: Number,
    default: 50 * 1024 * 1024, // 50MB
  },
})

const emit = defineEmits(['update:modelValue'])

const uploading = ref(false)

const getFileIcon = (type) => {
  if (type?.includes('pdf')) return icons.file
  if (type?.includes('word') || type?.includes('document')) return icons.file
  if (type?.includes('sheet') || type?.includes('excel')) return icons.file
  if (type?.includes('presentation') || type?.includes('powerpoint')) return icons.file
  return icons.file
}

const formatFileSize = (bytes) => {
  if (!bytes) return '未知'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const handleFileChange = async (e) => {
  const files = e.target.files
  if (!files || files.length === 0) return

  // 检查文件大小
  for (const file of files) {
    if (file.size > props.maxSize) {
      alert(`文件 "${file.name}" 超过大小限制 (${formatFileSize(props.maxSize)})`)
      return
    }
  }

  uploading.value = true

  try {
    const uploadedFiles = []

    for (const file of files) {
      const formData = new FormData()
      formData.append('file', file)

      const res = await api.post('/upload/image', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })

      if (res.url) {
        uploadedFiles.push({
          name: file.name,
          url: res.url,
          type: file.type,
          size: file.size,
        })
      }
    }

    emit('update:modelValue', [...props.modelValue, ...uploadedFiles])
  } catch (error) {
    console.error('文件上传失败:', error)
    alert('文件上传失败: ' + (error.error || error.message))
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

const removeFile = (index) => {
  const newFiles = [...props.modelValue]
  newFiles.splice(index, 1)
  emit('update:modelValue', newFiles)
}
</script>
