<template>
  <div class="rich-text-editor">
    <div style="border: 2px solid #808080; border-right-color: #ffffff; border-bottom-color: #ffffff;">
      <Toolbar
        style="border-bottom: 1px solid #808080"
        :editor="editorRef"
        :defaultConfig="toolbarConfig"
        mode="default"
      />
      <Editor
        style="height: 400px; overflow-y: auto; background: #fff;"
        v-model="valueHtml"
        :defaultConfig="editorConfig"
        mode="default"
        @onCreated="handleCreated"
        @onChange="handleChange"
      />
    </div>
    <div class="text-xs text-gray-500 mt-1">
      支持插入图片、视频，可设置文字大小、颜色、对齐方式等
    </div>
  </div>
</template>

<script setup>
import '@wangeditor/editor/dist/css/style.css'
import { ref, shallowRef, onBeforeUnmount, watch } from 'vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import api from '@/api'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])

const editorRef = shallowRef()
const valueHtml = ref(props.modelValue)

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  if (newVal !== valueHtml.value) {
    valueHtml.value = newVal
  }
})

// 工具栏配置
const toolbarConfig = {
  excludeKeys: [
    'group-video', // 可以开启，但需要配置上传
  ],
  toolbarKeys: [
    'headerSelect',
    'blockquote',
    '|',
    'bold',
    'underline',
    'italic',
    'through',
    'color',
    'bgColor',
    '|',
    'fontSize',
    'fontFamily',
    'lineHeight',
    '|',
    'bulletedList',
    'numberedList',
    'todo',
    '|',
    'justifyLeft',
    'justifyCenter',
    'justifyRight',
    'justifyJustify',
    '|',
    'insertLink',
    {
      key: 'group-image',
      title: '图片',
      iconSvg: '<svg viewBox="0 0 1024 1024"><path d="M959.877 128l0.123 0.123v767.775l-0.123 0.122H64.102l-0.122-0.122V128.123l0.122-0.123h895.775zM960 64H64C28.795 64 0 92.795 0 128v768c0 35.205 28.795 64 64 64h896c35.205 0 64-28.795 64-64V128c0-35.205-28.795-64-64-64zM832 288c0 53.020-42.98 96-96 96s-96-42.98-96-96 42.98-96 96-96 96 42.98 96 96zM896 832H128V704l224-384 256 320h64l224-192z"></path></svg>',
      menuKeys: ['insertImage', 'uploadImage'],
    },
    {
      key: 'group-video',
      title: '视频',
      menuKeys: ['insertVideo', 'uploadVideo'],
    },
    'insertTable',
    'codeBlock',
    '|',
    'undo',
    'redo',
    '|',
    'fullScreen',
  ],
}

// 编辑器配置
const editorConfig = {
  placeholder: '请输入商品详情描述...',
  MENU_CONF: {
    // 上传图片配置
    uploadImage: {
      // 自定义上传
      async customUpload(file, insertFn) {
        try {
          const formData = new FormData()
          formData.append('file', file)
          
          const res = await api.post('/upload/image', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
          })
          
          if (res.url) {
            insertFn(res.url, file.name, res.url)
          }
        } catch (error) {
          console.error('图片上传失败:', error)
          alert('图片上传失败: ' + (error.error || error.message))
        }
      },
      maxFileSize: 10 * 1024 * 1024, // 10MB
      allowedFileTypes: ['image/*'],
    },
    // 上传视频配置
    uploadVideo: {
      async customUpload(file, insertFn) {
        try {
          const formData = new FormData()
          formData.append('file', file)
          
          const res = await api.post('/upload/image', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
          })
          
          if (res.url) {
            insertFn(res.url)
          }
        } catch (error) {
          console.error('视频上传失败:', error)
          alert('视频上传失败: ' + (error.error || error.message))
        }
      },
      maxFileSize: 100 * 1024 * 1024, // 100MB
      allowedFileTypes: ['video/*'],
    },
  },
}

const handleCreated = (editor) => {
  editorRef.value = editor
}

const handleChange = (editor) => {
  const html = editor.getHtml()
  emit('update:modelValue', html)
}

// 暴露编辑器实例供外部使用
const getEditor = () => editorRef.value

defineExpose({
  getEditor,
  editorRef,
})

// 组件销毁时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor) {
    editor.destroy()
  }
})
</script>

<style scoped>
.rich-text-editor :deep(.w-e-toolbar) {
  background-color: #c0c0c0 !important;
  border-bottom: 1px solid #808080 !important;
}

.rich-text-editor :deep(.w-e-bar-item button) {
  background-color: #c0c0c0;
}

.rich-text-editor :deep(.w-e-bar-item button:hover) {
  background-color: #d4d4d4;
}
</style>
