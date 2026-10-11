<template>
  <div v-if="file && file.mode === 'url' && file.target_url" class="fixed inset-0 w-screen h-screen z-50 bg-white overflow-hidden">
    <iframe
      :src="file.target_url"
      class="w-full h-full border-none"
      allow="fullscreen; clipboard-read; clipboard-write; autoplay; camera; microphone"
    ></iframe>
  </div>

  <div v-else class="min-h-screen bg-slate-50 flex flex-col justify-center items-center px-4 py-8 pb-[12vh]">
    <div v-if="loading" class="text-center animate-in fade-in duration-300">
      <div class="w-12 h-12 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mx-auto mb-4"></div>
      <p class="text-xs font-mono font-bold text-slate-400 tracking-wider">LOADING ASSET...</p>
    </div>

    <div v-else-if="error" class="w-full max-w-sm bg-white rounded-3xl p-8 border border-slate-200 text-center shadow-xl">
      <div class="w-12 h-12 rounded-2xl bg-rose-50 text-rose-600 flex items-center justify-center mx-auto mb-4">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <h2 class="text-base font-bold text-slate-900 mb-2">访问受限</h2>
      <p class="text-xs text-slate-500 mb-6">{{ error }}</p>
      <router-link
        to="/"
        class="inline-block py-2.5 px-6 rounded-xl bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold transition"
      >
        返回首页
      </router-link>
    </div>

    <div v-else-if="file" class="w-full max-w-2xl bg-white rounded-3xl p-6 sm:p-8 border border-slate-200 shadow-xl shadow-slate-200/50 -translate-y-2 transition-transform">
      <div class="flex items-center justify-between pb-5 border-b border-slate-100 mb-6">
        <div class="flex items-center gap-3 overflow-hidden">
          <div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold text-xs shrink-0">
            {{ file.mode.toUpperCase() }}
          </div>
          <div class="truncate">
            <h1 class="text-base font-bold text-slate-900 truncate">{{ file.file_name }}</h1>
            <p class="text-xs text-slate-400 font-mono">{{ formatFileSize(file.file_size) }}</p>
          </div>
        </div>

        <a
          v-if="file.allow_download && file.download_url"
          :href="file.download_url"
          class="h-9 px-4 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold shadow-sm transition active:scale-95 flex items-center gap-1.5 shrink-0"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          <span>下载文件</span>
        </a>
      </div>

      <div v-if="isImage(file.mime_type)" class="text-center py-4">
        <img
          :src="file.public_url || file.download_url"
          :alt="file.file_name"
          class="max-h-[60vh] max-w-full rounded-2xl mx-auto shadow-sm object-contain"
        />
      </div>

      <div v-else-if="file.mode === 'text'" class="prose prose-slate max-w-none py-2 text-sm max-h-[65vh] overflow-y-auto" v-html="renderedMarkdown"></div>

      <div v-else class="text-center py-12 bg-slate-50 rounded-2xl border border-slate-100">
        <div class="w-16 h-16 rounded-2xl bg-indigo-100 text-indigo-600 flex items-center justify-center mx-auto mb-3">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
        </div>
        <h3 class="text-sm font-bold text-slate-800 mb-1">{{ file.file_name }}</h3>
        <p class="text-xs text-slate-400 font-mono mb-4">{{ formatFileSize(file.file_size) }}</p>
        <p v-if="!file.allow_download" class="text-xs text-slate-400">该文件仅供在线查验，未开启原件下载权限</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { apiClient } from '@/api/client'

const route = useRoute()
const loading = ref(true)
const error = ref('')
const file = ref<any>(null)

function formatFileSize(bytes: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function isImage(mime?: string) {
  return mime && mime.startsWith('image/')
}

const renderedMarkdown = computed(() => {
  if (!file.value?.text_content) return ''
  const raw = marked.parse(file.value.text_content) as string
  return DOMPurify.sanitize(raw)
})

async function fetchFile() {
  const fileId = route.params.fileId as string
  if (!fileId) {
    error.value = '未指定有效的文件ID'
    loading.value = false
    return
  }

  try {
    const res = await apiClient.get(`/share/${fileId}`)
    file.value = res.data

    if (file.value.mode === 'url' && file.value.target_url) {
      const target = file.value.target_url.trim()
      if (/^https?:\/\//i.test(target)) {
        document.title = file.value.file_name || '安全访问中...'
        const ua = navigator.userAgent.toLowerCase()
        const isTencent = ua.includes('micromessenger') || ua.includes('qq/')
        if (!isTencent) {
          try {
            window.location.replace(target)
          } catch (e) {
            window.location.href = target
          }
        }
        return
      }
      error.value = '非法的目标网址协议'
      return
    }
  } catch (err: any) {
    error.value = err?.response?.data?.error || '文件不存在或已失效'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchFile()
})
</script>
