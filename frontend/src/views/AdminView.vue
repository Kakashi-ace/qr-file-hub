<template>
  <div class="min-h-screen bg-slate-50 bg-grid-slate flex flex-col justify-between">
    <div>
      <Navbar />

      <main class="w-full max-w-5xl mx-auto px-4 pb-16">
        <div v-if="!isAuthenticated" class="max-w-md mx-auto pt-10">
          <div class="bg-white rounded-3xl p-8 border border-slate-200/80 shadow-xl shadow-slate-200/50 text-center">
            <div class="w-12 h-12 rounded-2xl bg-slate-900 text-white flex items-center justify-center mx-auto mb-4">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <h2 class="text-xl font-bold text-slate-900 mb-1">控制台身份认证</h2>
            <p class="text-xs text-slate-400 mb-6 font-mono">ASSET OPERATIONS DESK</p>

            <form @submit.prevent="handleLogin">
              <input
                type="password"
                v-model="loginPassword"
                placeholder="请输入管理密码"
                class="w-full px-4 py-3 rounded-xl border border-slate-200 bg-slate-50 text-sm focus:outline-none focus:border-slate-900 focus:bg-white transition mb-4 text-center font-mono"
              />
              <button
                type="submit"
                :disabled="isLoggingIn"
                class="w-full py-3 px-4 rounded-xl bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold transition active:scale-95 cursor-pointer disabled:opacity-50"
              >
                {{ isLoggingIn ? '验证中...' : '确认进入管理台' }}
              </button>
            </form>
          </div>
        </div>

        <div v-else>
          <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 mb-8">
            <div>
              <span class="inline-block px-3 py-1 rounded-full text-[10px] font-mono font-bold tracking-wider bg-slate-100 text-slate-600 border border-slate-200 mb-2">
                ASSET OPERATIONS DESK
              </span>
              <h1 class="text-2xl sm:text-3xl font-black text-slate-900 tracking-tight">
                云端文件资产控制台
              </h1>
            </div>

            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="triggerCleanExpired"
                class="h-9 px-3.5 rounded-xl bg-white hover:bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-700 shadow-2xs transition active:scale-95 cursor-pointer flex items-center gap-1.5"
              >
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
                <span>清理全部已过期文件</span>
              </button>
              <button
                type="button"
                @click="handleLogout"
                class="h-9 px-3 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-600 text-xs font-semibold transition active:scale-95 cursor-pointer"
              >
                锁定退出
              </button>
            </div>
          </div>

          <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-xs">
              <span class="text-xs font-medium text-slate-400">总托管资产数</span>
              <h3 class="text-2xl sm:text-3xl font-black font-mono text-slate-900 mt-2">
                {{ stats.total_files }}
              </h3>
            </div>
            <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-xs">
              <span class="text-xs font-medium text-slate-400">R2 累计存储占用</span>
              <h3 class="text-2xl sm:text-3xl font-black font-mono text-slate-900 mt-2">
                {{ formatFileSize(stats.total_storage) }}
              </h3>
            </div>
            <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-xs">
              <span class="text-xs font-medium text-slate-400">累计扫码下载量</span>
              <h3 class="text-2xl sm:text-3xl font-black font-mono text-slate-900 mt-2">
                {{ stats.total_downloads }}
              </h3>
            </div>
            <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-xs">
              <span class="text-xs font-medium text-slate-400">已过期失效文件</span>
              <h3 class="text-2xl sm:text-3xl font-black font-mono text-rose-600 mt-2">
                {{ stats.expired_count }}
              </h3>
            </div>
          </div>

          <div class="bg-white rounded-3xl p-6 border border-slate-200/80 shadow-xs mb-8">
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 pb-5 border-b border-slate-100 mb-5">
              <div class="flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full bg-slate-900"></span>
                <h3 class="text-base font-bold text-slate-900">GitHub 防红跳板矩阵与健康检测池</h3>
                <span class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-slate-100 text-slate-600 border border-slate-200">微信 / QQ 拦截探针</span>
              </div>
              <div class="flex items-center gap-2">
                <button
                  type="button"
                  @click="probeAllBridges"
                  :disabled="isProbing"
                  class="h-8 px-3 rounded-xl bg-white hover:bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-700 shadow-2xs transition active:scale-95 cursor-pointer disabled:opacity-50"
                >
                  {{ isProbing ? '体检中...' : '一键体检全部跳板' }}
                </button>
                <button
                  type="button"
                  @click="promptAddBridge"
                  class="h-8 px-3 rounded-xl bg-slate-900 hover:bg-slate-800 text-white text-xs font-semibold shadow-2xs transition active:scale-95 cursor-pointer"
                >
                  新建 1 个跳板
                </button>
              </div>
            </div>

            <div class="overflow-x-auto">
              <table class="w-full text-left text-xs border-collapse">
                <thead>
                  <tr class="border-b border-slate-100 text-slate-400 font-mono text-[11px] bg-slate-50/50">
                    <th class="py-3 px-3">角色状态</th>
                    <th class="py-3 px-3">跳板仓库与域名</th>
                    <th class="py-3 px-3">微信检测状态</th>
                    <th class="py-3 px-3">QQ 检测状态</th>
                    <th class="py-3 px-3 text-right">跳板管理操作</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr v-for="b in bridges" :key="b.id">
                    <td class="py-3 px-3">
                      <span
                        class="px-2 py-0.5 rounded-md font-bold text-[10px]"
                        :class="b.is_primary ? 'bg-indigo-50 text-indigo-700 border border-indigo-200' : 'bg-slate-100 text-slate-500'"
                      >
                        {{ b.is_primary ? '当前主跳板' : '备用跳板' }}
                      </span>
                    </td>
                    <td class="py-3 px-3 font-mono font-medium text-slate-800">
                      {{ b.domain }}
                    </td>
                    <td class="py-3 px-3">
                      <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-600 border border-emerald-200">
                        <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                        正常通过
                      </span>
                    </td>
                    <td class="py-3 px-3">
                      <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-600 border border-emerald-200">
                        <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                        正常通过
                      </span>
                    </td>
                    <td class="py-3 px-3 text-right">
                      <div class="flex items-center justify-end gap-1.5">
                        <button
                          v-if="!b.is_primary"
                          type="button"
                          @click="setPrimaryBridge(b.id)"
                          class="px-2.5 py-1 rounded-lg bg-slate-100 hover:bg-slate-200 text-slate-700 text-[11px] font-semibold transition cursor-pointer"
                        >
                          设为主跳板
                        </button>
                        <button
                          v-if="b.is_custom"
                          type="button"
                          @click="deleteBridge(b.id)"
                          class="px-2.5 py-1 rounded-lg bg-rose-50 hover:bg-rose-100 text-rose-600 text-[11px] font-semibold transition cursor-pointer"
                        >
                          删除
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div class="bg-white rounded-3xl p-6 border border-slate-200/80 shadow-xs">
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 pb-5 border-b border-slate-100 mb-5">
              <div class="flex items-center gap-2.5">
                <span class="w-2.5 h-2.5 rounded-full bg-slate-900"></span>
                <h3 class="text-base font-bold text-slate-900">托管文件清单</h3>
                <span class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-slate-100 text-slate-600 border border-slate-200">实时同步</span>

                <div v-if="selectedFileIds.length > 0" class="flex items-center gap-2 ml-4">
                  <span class="text-xs text-slate-500 font-medium">已选 {{ selectedFileIds.length }} 项</span>
                  <button
                    type="button"
                    @click="triggerBatchDelete"
                    class="px-2.5 py-1 rounded-lg bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition active:scale-95 cursor-pointer"
                  >
                    批量删除
                  </button>
                </div>
              </div>

              <div class="flex items-center gap-2 w-full sm:w-auto">
                <input
                  type="text"
                  v-model="searchKeyword"
                  @input="fetchFiles"
                  placeholder="搜索文件名或ID..."
                  class="px-3 py-1.5 rounded-xl border border-slate-200 bg-slate-50 text-xs focus:outline-none focus:border-slate-900 focus:bg-white transition"
                />
                <button
                  type="button"
                  @click="fetchFiles"
                  class="h-8 px-3 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-semibold transition cursor-pointer shrink-0"
                >
                  刷新
                </button>
              </div>
            </div>

            <div class="overflow-x-auto">
              <table class="w-full text-left text-xs border-collapse">
                <thead>
                  <tr class="border-b border-slate-100 text-slate-400 font-mono text-[11px] bg-slate-50/50">
                    <th class="py-3 px-3 w-8">
                      <input
                        type="checkbox"
                        :checked="isAllSelected"
                        @change="toggleSelectAll"
                        class="rounded border-slate-300 text-slate-900 focus:ring-0"
                      />
                    </th>
                    <th class="py-3 px-3">文件名与专属链接</th>
                    <th class="py-3 px-3">大小</th>
                    <th class="py-3 px-3">有效状态</th>
                    <th class="py-3 px-3">下载控制</th>
                    <th class="py-3 px-3">上传时间</th>
                    <th class="py-3 px-3">访问量</th>
                    <th class="py-3 px-3 text-right">管理操作</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr v-for="f in files" :key="f.file_id">
                    <td class="py-3 px-3">
                      <input
                        type="checkbox"
                        v-model="selectedFileIds"
                        :value="f.file_id"
                        class="rounded border-slate-300 text-slate-900 focus:ring-0"
                      />
                    </td>
                    <td class="py-3 px-3 max-w-xs">
                      <div class="flex items-center gap-2">
                        <span class="px-1.5 py-0.5 rounded text-[10px] font-mono uppercase bg-slate-100 text-slate-600 font-bold">
                          {{ f.mode }}
                        </span>
                        <span class="font-bold text-slate-900 truncate">{{ f.file_name }}</span>
                      </div>
                      <div class="text-[11px] text-slate-400 font-mono truncate mt-0.5">
                        ID: {{ f.file_id }}
                      </div>
                    </td>
                    <td class="py-3 px-3 font-mono text-slate-600">
                      {{ formatFileSize(f.file_size) }}
                    </td>
                    <td class="py-3 px-3">
                      <span
                        class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                        :class="isExpired(f.expire_at) ? 'bg-rose-50 text-rose-600 border border-rose-200' : 'bg-slate-100 text-slate-600 border border-slate-200'"
                      >
                        {{ isExpired(f.expire_at) ? '已失效' : '永久有效' }}
                      </span>
                    </td>
                    <td class="py-3 px-3">
                      <span
                        class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                        :class="f.allow_download ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-500'"
                      >
                        {{ f.allow_download ? '可下载' : '仅预览' }}
                      </span>
                    </td>
                    <td class="py-3 px-3 font-mono text-slate-400 text-[11px]">
                      {{ formatDate(f.created_at) }}
                    </td>
                    <td class="py-3 px-3 font-mono font-bold text-slate-700">
                      {{ f.view_count }} 次
                    </td>
                    <td class="py-3 px-3 text-right">
                      <div class="flex items-center justify-end gap-1.5">
                        <router-link
                          :to="`/s/${f.file_id}`"
                          target="_blank"
                          class="px-2.5 py-1 rounded-lg bg-slate-100 hover:bg-slate-200 text-slate-700 text-[11px] font-semibold transition cursor-pointer"
                        >
                          查看
                        </router-link>
                        <button
                          type="button"
                          @click="openQrModalForFile(f)"
                          class="px-2.5 py-1 rounded-lg bg-indigo-50 hover:bg-indigo-100 text-indigo-700 text-[11px] font-semibold transition cursor-pointer"
                        >
                          二维码
                        </button>
                        <button
                          type="button"
                          @click="triggerDeleteSingle(f.file_id)"
                          class="px-2.5 py-1 rounded-lg bg-rose-50 hover:bg-rose-100 text-rose-600 text-[11px] font-semibold transition cursor-pointer"
                        >
                          删除
                        </button>
                      </div>
                    </td>
                  </tr>

                  <tr v-if="files.length === 0">
                    <td colspan="8" class="py-8 text-center text-slate-400">
                      暂无托管文件资产
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </main>
    </div>

    <ModalQr
      :visible="showQrModal"
      :title="activeQrItem.title"
      :direct-url="activeQrItem.directUrl"
      :bridge-url="activeQrItem.bridgeUrl"
      @close="showQrModal = false"
    />

    <ProgressModal
      :visible="progressModal.visible"
      :title="progressModal.title"
      :message="progressModal.message"
      :progress="progressModal.progress"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Navbar from '@/components/Navbar.vue'
import ModalQr from '@/components/ModalQr.vue'
import ProgressModal from '@/components/ProgressModal.vue'
import { apiClient, FileItem, BridgeItem, AdminStats } from '@/api/client'

const isAuthenticated = ref(false)
const loginPassword = ref('')
const isLoggingIn = ref(false)

const stats = ref<AdminStats>({
  total_files: 0,
  total_storage: 0,
  total_downloads: 0,
  expired_count: 0
})

const bridges = ref<BridgeItem[]>([])
const isProbing = ref(false)

const files = ref<FileItem[]>([])
const selectedFileIds = ref<string[]>([])
const searchKeyword = ref('')

const showQrModal = ref(false)
const activeQrItem = ref({
  title: '',
  directUrl: '',
  bridgeUrl: ''
})

const progressModal = ref({
  visible: false,
  title: '',
  message: '',
  progress: 0
})

const isAllSelected = computed(() => {
  return files.value.length > 0 && selectedFileIds.value.length === files.value.length
})

function formatFileSize(bytes: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN', { hour12: false })
}

function isExpired(expireAt?: string) {
  if (!expireAt) return false
  return new Date(expireAt).getTime() < Date.now()
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedFileIds.value = []
  } else {
    selectedFileIds.value = files.value.map(f => f.file_id)
  }
}

async function handleLogin() {
  if (!loginPassword.value) return
  isLoggingIn.value = true
  try {
    const res = await apiClient.post<{ token: string }>('/admin/login', {
      password: loginPassword.value
    })
    localStorage.setItem('qr_admin_token', res.data.token)
    isAuthenticated.value = true
    await loadAdminData()
  } catch (err: any) {
    alert(err?.response?.data?.error || '密码错误')
  } finally {
    isLoggingIn.value = false
  }
}

function handleLogout() {
  localStorage.removeItem('qr_admin_token')
  isAuthenticated.value = false
}

async function checkAuthStatus() {
  const token = localStorage.getItem('qr_admin_token')
  if (!token) return
  try {
    await apiClient.get('/admin/check')
    isAuthenticated.value = true
    await loadAdminData()
  } catch (err) {
    localStorage.removeItem('qr_admin_token')
    isAuthenticated.value = false
  }
}

async function loadAdminData() {
  await Promise.all([fetchStats(), fetchBridges(), fetchFiles()])
}

async function fetchStats() {
  try {
    const res = await apiClient.get<AdminStats>('/admin/stats')
    stats.value = res.data
  } catch (err) {}
}

async function fetchBridges() {
  try {
    const res = await apiClient.get<{ data: BridgeItem[] }>('/admin/bridges')
    bridges.value = res.data.data
  } catch (err) {}
}

async function fetchFiles() {
  try {
    const res = await apiClient.get<{ items: FileItem[] }>('/admin/files', {
      params: { keyword: searchKeyword.value }
    })
    files.value = res.data.items || []
  } catch (err) {}
}

async function probeAllBridges() {
  isProbing.value = true
  try {
    const res = await apiClient.post<{ data: BridgeItem[] }>('/admin/bridges/probe')
    bridges.value = res.data.data
  } catch (err) {
  } finally {
    isProbing.value = false
  }
}

function promptAddBridge() {
  const custom = window.prompt('请输入新的跳板域名：')
  if (custom && custom.trim()) {
    apiClient.post('/admin/bridges', { domain: custom.trim() }).then(() => {
      fetchBridges()
    })
  }
}

async function setPrimaryBridge(id: number) {
  try {
    await apiClient.put(`/admin/bridges/${id}/primary`)
    await fetchBridges()
  } catch (err) {}
}

async function deleteBridge(id: number) {
  if (!confirm('确定删除该跳板吗？')) return
  try {
    await apiClient.delete(`/admin/bridges/${id}`)
    await fetchBridges()
  } catch (err) {}
}

function openQrModalForFile(f: FileItem) {
  const direct = `${window.location.origin}/s/${f.file_id}`
  const bridge = f.bridge_domain ? `${f.bridge_domain.replace(/\/$/, '')}/?f=${f.file_id}` : direct
  activeQrItem.value = {
    title: f.file_name,
    directUrl: direct,
    bridgeUrl: bridge
  }
  showQrModal.value = true
}

async function runWithProgress(title: string, message: string, task: () => Promise<void>) {
  progressModal.value = {
    visible: true,
    title,
    message,
    progress: 15
  }
  const timer = setInterval(() => {
    if (progressModal.value.progress < 85) {
      progressModal.value.progress += 15
    }
  }, 100)

  try {
    await task()
    progressModal.value.progress = 100
    await new Promise(r => setTimeout(r, 300))
  } finally {
    clearInterval(timer)
    progressModal.value.visible = false
  }
}

async function triggerDeleteSingle(fileId: string) {
  if (!confirm('确定彻底删除该资产文件吗？')) return
  await runWithProgress('正在物理删除', '正在从云端存储与数据库清除记录...', async () => {
    await apiClient.delete(`/admin/files/${fileId}`)
    await loadAdminData()
  })
}

async function triggerBatchDelete() {
  if (selectedFileIds.value.length === 0) return
  if (!confirm(`确定彻底批量删除选中的 ${selectedFileIds.value.length} 个资产文件吗？`)) return
  await runWithProgress('正在批量删除', '正在并发销毁存储对象与元数据...', async () => {
    await apiClient.post('/admin/files/batch-delete', {
      file_ids: selectedFileIds.value
    })
    selectedFileIds.value = []
    await loadAdminData()
  })
}

async function triggerCleanExpired() {
  if (!confirm('确定清理释放所有已失效的过期文件吗？')) return
  await runWithProgress('正在深度清理', '正在扫描并销毁全部失效对象...', async () => {
    await apiClient.post('/admin/files/clean-expired')
    await loadAdminData()
  })
}

onMounted(() => {
  checkAuthStatus()
})
</script>
