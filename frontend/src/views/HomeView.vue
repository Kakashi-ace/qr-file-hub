<template>
  <div class="min-h-screen bg-slate-50 bg-grid-slate flex flex-col justify-between">
    <div>
      <Navbar />

      <main class="w-full max-w-2xl mx-auto px-4 pb-12">
        <div class="text-center mb-8">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-[11px] font-mono font-bold tracking-widest bg-indigo-50 text-indigo-600 border border-indigo-200/60 mb-3 shadow-2xs">
            HIGH-SPEED ASSET DISTRIBUTION
          </div>
          <h1 class="text-3xl sm:text-4xl font-black text-slate-900 tracking-tight">
            文件分发与 Markdown 在线制码
          </h1>
        </div>

        <div class="bg-white rounded-3xl p-6 sm:p-8 border border-slate-200/80 shadow-xl shadow-slate-200/50">
          <div class="flex items-center justify-center p-1 rounded-xl bg-slate-100/90 border border-slate-200/80 mb-6 max-w-md mx-auto">
            <button
              type="button"
              @click="mode = 'file'"
              class="flex-1 py-2 px-3 rounded-lg text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer"
              :class="mode === 'file' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-600 hover:text-slate-900'"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
              <span>文件上传</span>
            </button>
            <button
              type="button"
              @click="mode = 'text'"
              class="flex-1 py-2 px-3 rounded-lg text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer"
              :class="mode === 'text' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-600 hover:text-slate-900'"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span>Markdown 编辑</span>
            </button>
            <button
              type="button"
              @click="mode = 'url'"
              class="flex-1 py-2 px-3 rounded-lg text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer"
              :class="mode === 'url' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-600 hover:text-slate-900'"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
              </svg>
              <span>网址防红</span>
            </button>
          </div>

          <div v-if="mode === 'file'" class="mb-6">
            <div
              @dragover.prevent="isDragging = true"
              @dragleave.prevent="isDragging = false"
              @drop.prevent="onDropFile"
              @click="triggerFileInput"
              class="rounded-2xl p-8 sm:p-12 text-center cursor-pointer border-2 border-dashed transition-all relative overflow-hidden"
              :class="isDragging ? 'border-indigo-500 bg-indigo-50/50' : 'border-slate-300 bg-slate-50/70 hover:border-indigo-400 hover:bg-indigo-50/20'"
            >
              <input type="file" ref="fileInputRef" class="hidden" @change="onFileSelected" />
              <div v-if="!selectedFile">
                <div class="w-16 h-16 rounded-2xl bg-indigo-50 text-indigo-600 mx-auto flex items-center justify-center mb-4 border border-indigo-200/60 shadow-sm">
                  <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                  </svg>
                </div>
                <p class="text-base font-semibold text-slate-800">
                  拖拽文件到此处，或 <span class="text-indigo-600 underline decoration-indigo-400/60 font-bold">点击选择文件</span>
                </p>
              </div>

              <div v-else class="flex items-center justify-between p-4 rounded-xl bg-white border border-slate-200 shadow-sm text-left">
                <div class="flex items-center gap-3 overflow-hidden">
                  <div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold text-xs shrink-0">
                    {{ selectedFile.name.split('.').pop()?.toUpperCase() || 'FILE' }}
                  </div>
                  <div class="truncate">
                    <p class="text-sm font-bold text-slate-900 truncate">{{ selectedFile.name }}</p>
                    <p class="text-xs text-slate-400 font-mono">{{ formatFileSize(selectedFile.size) }}</p>
                  </div>
                </div>
                <button
                  type="button"
                  @click.stop="selectedFile = null"
                  class="px-3 py-1.5 rounded-lg text-xs font-semibold text-slate-500 hover:text-rose-600 hover:bg-rose-50 transition cursor-pointer"
                >
                  更换
                </button>
              </div>
            </div>
          </div>

          <div v-else-if="mode === 'text'" class="mb-6">
            <div class="mb-4">
              <label class="block text-xs font-bold text-slate-700 mb-1.5">文档标题（选填）</label>
              <input
                type="text"
                v-model="markdownTitle"
                placeholder="输入文档标题（选填）"
                class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:bg-white transition"
              />
            </div>

            <div class="border border-slate-200 rounded-2xl overflow-hidden bg-slate-50/50">
              <div class="flex items-center justify-between p-2 border-b border-slate-200 bg-white">
                <div class="flex items-center gap-1">
                  <button
                    type="button"
                    @click="editorTab = 'source'"
                    class="py-1 px-3 rounded-lg text-xs font-bold transition cursor-pointer"
                    :class="editorTab === 'source' ? 'bg-indigo-600 text-white' : 'text-slate-600 hover:text-slate-900'"
                  >
                    编辑源码
                  </button>
                  <button
                    type="button"
                    @click="editorTab = 'preview'"
                    class="py-1 px-3 rounded-lg text-xs font-bold transition cursor-pointer"
                    :class="editorTab === 'preview' ? 'bg-indigo-600 text-white' : 'text-slate-600 hover:text-slate-900'"
                  >
                    实时渲染预览
                  </button>
                </div>
                <div class="flex items-center gap-2">
                  <button
                    type="button"
                    @click="fillMarkdownExample"
                    class="text-xs text-indigo-600 hover:text-indigo-700 font-bold cursor-pointer"
                  >
                    填入示例模版
                  </button>
                  <span class="text-slate-300">|</span>
                  <button
                    type="button"
                    @click="markdownContent = ''"
                    class="text-xs text-slate-400 hover:text-rose-600 cursor-pointer"
                  >
                    清空
                  </button>
                </div>
              </div>

              <div v-if="editorTab === 'source'" class="p-3">
                <textarea
                  v-model="markdownContent"
                  rows="10"
                  placeholder="在此输入 Markdown 标题与正文..."
                  class="w-full p-2 bg-transparent text-sm text-slate-800 placeholder:text-slate-400 font-mono resize-y focus:outline-none"
                ></textarea>
                <div class="text-right text-[11px] text-slate-400 font-mono pt-2">
                  {{ markdownContent.length }} 字
                </div>
              </div>

              <div v-else class="p-5 max-h-96 overflow-y-auto prose prose-slate prose-sm bg-white" v-html="renderedMarkdown"></div>
            </div>
          </div>

          <div v-else-if="mode === 'url'" class="mb-6">
            <label class="block text-xs font-bold text-slate-700 mb-1.5">跳转目标网址</label>
            <input
              type="url"
              v-model="targetUrl"
              placeholder="https://example.com"
              class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:bg-white transition mb-3"
            />
            <label class="block text-xs font-bold text-slate-700 mb-1.5">网址备注标题（选填）</label>
            <input
              type="text"
              v-model="urlTitle"
              placeholder="输入网址备注标题（选填）"
              class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:bg-white transition"
            />
          </div>

          <div class="p-4 rounded-2xl bg-slate-50 border border-slate-200/80 mb-6">
            <div class="flex items-center gap-2 mb-3">
              <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
              <span class="text-xs font-bold text-slate-800">专属高信誉防红网桥</span>
            </div>
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-2">
              <select
                v-model="selectedBridge"
                class="flex-1 px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs font-medium text-slate-800 focus:outline-none focus:border-indigo-500"
              >
                <option v-for="b in bridges" :key="b.domain" :value="b.domain">
                  {{ b.is_primary ? '⭐️ [推荐] ' : '🛡️ [跳板] ' }}{{ b.domain }}
                </option>
              </select>
              <button
                type="button"
                @click="promptCustomBridge"
                class="px-3.5 py-2 rounded-xl bg-white hover:bg-slate-50 border border-slate-200 text-xs font-bold text-slate-700 shadow-2xs transition active:scale-95 cursor-pointer shrink-0"
              >
                + 自定义
              </button>
            </div>
          </div>

          <div class="mb-6">
            <label class="block text-xs font-bold text-slate-700 mb-2">二维码与链接有效期</label>
            <div class="grid grid-cols-3 sm:grid-cols-6 gap-2">
              <button
                v-for="opt in expireOptions"
                :key="opt.value"
                type="button"
                @click="expireType = opt.value"
                class="py-2 px-2 rounded-xl text-xs font-bold transition-all border text-center cursor-pointer"
                :class="expireType === opt.value ? 'bg-indigo-50 border-indigo-500 text-indigo-700 shadow-2xs' : 'bg-white border-slate-200 text-slate-600 hover:border-slate-300'"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>

          <div v-if="mode !== 'url'" class="flex items-center justify-between p-3.5 rounded-2xl bg-white border border-slate-200 mb-4">
            <span class="text-xs font-bold text-slate-800">
              {{ mode === 'file' ? '允许访客下载原文件' : '允许访客下载 Markdown 原文件' }}
            </span>
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="allowDownload = !allowDownload"
                class="w-12 h-6 rounded-full transition-colors relative cursor-pointer"
                :class="allowDownload ? 'bg-indigo-600' : 'bg-slate-300'"
              >
                <span
                  class="absolute top-1 left-1 w-4 h-4 rounded-full bg-white transition-transform"
                  :class="allowDownload ? 'translate-x-6' : 'translate-x-0'"
                ></span>
              </button>
              <span
                class="text-[11px] font-bold px-2 py-0.5 rounded-md"
                :class="allowDownload ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-500'"
              >
                {{ allowDownload ? '已开启' : '已关闭' }}
              </span>
            </div>
          </div>

          <div class="flex items-center justify-between p-3.5 rounded-2xl bg-white border border-slate-200 mb-6">
            <div class="flex items-center gap-2">
              <span class="text-xs font-bold text-slate-800">启用 GitHub 微信防封网桥</span>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-600 border border-emerald-200">推荐</span>
            </div>
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="enableBridge = !enableBridge"
                class="w-12 h-6 rounded-full transition-colors relative cursor-pointer"
                :class="enableBridge ? 'bg-indigo-600' : 'bg-slate-300'"
              >
                <span
                  class="absolute top-1 left-1 w-4 h-4 rounded-full bg-white transition-transform"
                  :class="enableBridge ? 'translate-x-6' : 'translate-x-0'"
                ></span>
              </button>
              <span
                class="text-[11px] font-bold px-2 py-0.5 rounded-md"
                :class="enableBridge ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-500'"
              >
                {{ enableBridge ? '已开启' : '已关闭' }}
              </span>
            </div>
          </div>

          <button
            type="button"
            @click="handleSubmit"
            :disabled="isSubmitting"
            class="w-full py-3.5 px-4 rounded-2xl bg-indigo-600 hover:bg-indigo-700 active:scale-[0.99] text-white text-sm font-bold shadow-lg shadow-indigo-600/25 transition-all flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
          >
            <svg v-if="isSubmitting" class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ submitButtonText }}</span>
          </button>
        </div>
      </main>
    </div>

    <ModalQr
      :visible="showQrModal"
      :title="qrResult.title"
      :direct-url="qrResult.directUrl"
      :bridge-url="qrResult.bridgeUrl"
      @close="showQrModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import axios from 'axios'
import Navbar from '@/components/Navbar.vue'
import ModalQr from '@/components/ModalQr.vue'
import { apiClient, BridgeItem, PresignUploadResponse, ShareResultResponse } from '@/api/client'

const mode = ref<'file' | 'text' | 'url'>('file')
const isDragging = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)

const markdownTitle = ref('')
const markdownContent = ref('')
const editorTab = ref<'source' | 'preview'>('source')

const targetUrl = ref('')
const urlTitle = ref('')

const expireType = ref('never')
const expireOptions = [
  { label: '永久有效', value: 'never' },
  { label: '1 小时', value: '1h' },
  { label: '24 小时 (1天)', value: '24h' },
  { label: '7 天', value: '7d' },
  { label: '30 天', value: '30d' },
  { label: '自定义时间', value: 'custom' },
]

const allowDownload = ref(false)
const enableBridge = ref(true)

const bridges = ref<BridgeItem[]>([
  { id: 1, domain: 'https://kakashi-ace.github.io', is_active: true, is_primary: true, is_custom: false, wechat_status: 'normal', qq_status: 'normal', latency_ms: 68 },
  { id: 2, domain: 'https://kakashi-ace.github.io/bridge-v4zahz/', is_active: true, is_primary: false, is_custom: false, wechat_status: 'normal', qq_status: 'normal', latency_ms: 75 }
])
const selectedBridge = ref('https://kakashi-ace.github.io')

const isSubmitting = ref(false)
const showQrModal = ref(false)
const qrResult = ref({
  title: '',
  directUrl: '',
  bridgeUrl: ''
})

const renderedMarkdown = computed(() => {
  if (!markdownContent.value) return '<p class="text-slate-400">暂无预览内容</p>'
  const raw = marked.parse(markdownContent.value) as string
  return DOMPurify.sanitize(raw)
})

const submitButtonText = computed(() => {
  if (isSubmitting.value) return '正在生成中...'
  if (mode.value === 'file') return '立即上传并生成二维码'
  if (mode.value === 'text') return '立即生成 Markdown 二维码'
  return '立即生成防红链接与二维码'
})

function formatFileSize(bytes: number) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function triggerFileInput() {
  fileInputRef.value?.click()
}

function onFileSelected(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    selectedFile.value = target.files[0]
  }
}

function onDropFile(e: DragEvent) {
  isDragging.value = false
  if (e.dataTransfer?.files && e.dataTransfer.files[0]) {
    selectedFile.value = e.dataTransfer.files[0]
  }
}

function fillMarkdownExample() {
  markdownTitle.value = '技术方案说明文档'
  markdownContent.value = `# 项目架构设计\n\n- 前端：Cloudflare Pages (Vue 3 + Tailwind CSS)\n- 后端：Render (Go + Gin)\n- 数据库：TiDB Cloud\n- 存储：Cloudflare R2\n\n> 具备极高安全等级与防拦截性能。`
}

function promptCustomBridge() {
  const custom = window.prompt('请输入自定义跳板域名（例如 https://your-domain.com）：')
  if (custom && custom.trim()) {
    let domain = custom.trim()
    if (!domain.startsWith('http://') && !domain.startsWith('https://')) {
      domain = 'https://' + domain
    }
    bridges.value.push({
      id: Date.now(),
      domain: domain,
      is_active: true,
      is_primary: false,
      is_custom: true,
      wechat_status: 'normal',
      qq_status: 'normal',
      latency_ms: 0
    })
    selectedBridge.value = domain
  }
}

async function loadBridges() {
  try {
    const res = await apiClient.get<{ data: BridgeItem[] }>('/bridges')
    if (res.data && res.data.data.length > 0) {
      bridges.value = res.data.data
      const primary = res.data.data.find(b => b.is_primary)
      if (primary) {
        selectedBridge.value = primary.domain
      }
    }
  } catch (err) {
  }
}

async function handleSubmit() {
  if (mode.value === 'file') {
    if (!selectedFile.value) {
      alert('请先选择需要上传的文件')
      return
    }
    await uploadFileFlow()
  } else if (mode.value === 'text') {
    if (!markdownContent.value.trim()) {
      alert('请输入 Markdown 正文内容')
      return
    }
    await createMarkdownFlow()
  } else if (mode.value === 'url') {
    if (!targetUrl.value.trim()) {
      alert('请输入跳转目标网址')
      return
    }
    await createUrlFlow()
  }
}

async function uploadFileFlow() {
  const file = selectedFile.value!
  isSubmitting.value = true
  try {
    const presignRes = await apiClient.post<PresignUploadResponse>('/upload/presign', {
      file_name: file.name,
      file_size: file.size,
      mime_type: file.type || 'application/octet-stream'
    })

    const { file_id, r2_key, upload_url } = presignRes.data

    await axios.put(upload_url, file, {
      headers: {
        'Content-Type': file.type || 'application/octet-stream'
      }
    })

    const completeRes = await apiClient.post<ShareResultResponse>('/upload/complete', {
      file_id,
      file_name: file.name,
      file_size: file.size,
      mime_type: file.type || 'application/octet-stream',
      r2_key,
      allow_download: allowDownload.value,
      bridge_domain: enableBridge.value ? selectedBridge.value : '',
      expire_type: expireType.value
    })

    qrResult.value = {
      title: file.name,
      directUrl: completeRes.data.direct_url,
      bridgeUrl: completeRes.data.bridge_url
    }
    showQrModal.value = true
  } catch (err: any) {
    alert(err?.response?.data?.error || '上传或制码失败')
  } finally {
    isSubmitting.value = false
  }
}

async function createMarkdownFlow() {
  isSubmitting.value = true
  try {
    const res = await apiClient.post<ShareResultResponse>('/markdown/create', {
      title: markdownTitle.value,
      content: markdownContent.value,
      allow_download: allowDownload.value,
      bridge_domain: enableBridge.value ? selectedBridge.value : '',
      expire_type: expireType.value
    })

    qrResult.value = {
      title: res.data.file_name || 'Markdown 文档',
      directUrl: res.data.direct_url,
      bridgeUrl: res.data.bridge_url
    }
    showQrModal.value = true
  } catch (err: any) {
    alert(err?.response?.data?.error || '生成 Markdown 失败')
  } finally {
    isSubmitting.value = false
  }
}

async function createUrlFlow() {
  isSubmitting.value = true
  try {
    const res = await apiClient.post<ShareResultResponse>('/url/create', {
      target_url: targetUrl.value,
      title: urlTitle.value,
      bridge_domain: enableBridge.value ? selectedBridge.value : '',
      expire_type: expireType.value
    })

    let bUrl = res.data.bridge_url
    if (bUrl && targetUrl.value && !bUrl.includes('&c=') && !bUrl.includes('?c=')) {
      try {
        const clean = targetUrl.value.trim()
        const b64 = btoa(clean).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
        bUrl = bUrl + (bUrl.includes('?') ? '&c=' : '?c=') + b64
      } catch (e) {}
    }

    qrResult.value = {
      title: urlTitle.value || targetUrl.value,
      directUrl: res.data.direct_url,
      bridgeUrl: bUrl
    }
    showQrModal.value = true
  } catch (err: any) {
    alert(err?.response?.data?.error || '生成防红链接失败')
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  loadBridges()
})
</script>
