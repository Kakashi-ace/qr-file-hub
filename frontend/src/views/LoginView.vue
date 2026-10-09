<template>
  <div class="min-h-screen bg-slate-50 bg-grid-slate flex flex-col justify-center items-center px-4 py-12">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="w-14 h-14 rounded-2xl bg-indigo-600 text-white flex items-center justify-center font-black text-2xl shadow-lg shadow-indigo-600/30 mx-auto mb-4">
          Q
        </div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">QRFile 共享中心</h1>
        <p class="text-xs text-slate-500 font-mono mt-1">ENTERPRISE ACCESS MANAGEMENT</p>
      </div>

      <div class="bg-white rounded-3xl p-8 border border-slate-200/80 shadow-xl shadow-slate-200/60">
        <div class="mb-6 pb-4 border-b border-slate-100 flex items-center justify-between">
          <div>
            <h2 class="text-base font-bold text-slate-900">系统身份登录</h2>
            <p class="text-xs text-slate-400 mt-0.5">请使用授权账号登录工作台</p>
          </div>
          <span class="px-2.5 py-1 rounded-full text-[10px] font-mono font-bold bg-indigo-50 text-indigo-600 border border-indigo-100">
            SECURE
          </span>
        </div>

        <div v-if="errorMessage" class="mb-5 p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-600 text-xs font-semibold flex items-center gap-2">
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{{ errorMessage }}</span>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">登录账号</label>
            <div class="relative">
              <input
                type="text"
                v-model="username"
                placeholder="请输入账号（默认 admin）"
                required
                autocomplete="username"
                class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:border-indigo-600 focus:bg-white transition"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">登录密码</label>
            <div class="relative">
              <input
                type="password"
                v-model="password"
                placeholder="请输入密码"
                required
                autocomplete="current-password"
                class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:border-indigo-600 focus:bg-white transition"
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full mt-2 py-3 px-4 rounded-xl bg-indigo-600 hover:bg-indigo-700 active:scale-[0.99] text-white text-xs font-bold shadow-md shadow-indigo-600/20 transition-all flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
          >
            <svg v-if="isLoading" class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ isLoading ? '正在鉴权...' : '确认登录系统' }}</span>
          </button>
        </form>

        <div class="mt-6 pt-4 border-t border-slate-100 text-center">
          <p class="text-[11px] text-slate-400">
            受统一身份安全策略保护，所有操作均记录安全审计日志
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { apiClient, setAuthToken, getAuthToken } from '@/api/client'

const router = useRouter()
const route = useRoute()

const username = ref('admin')
const password = ref('')
const isLoading = ref(false)
const errorMessage = ref('')

async function handleLogin() {
  if (!password.value) {
    errorMessage.value = '请输入登录密码'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    const res = await apiClient.post('/auth/login', {
      username: username.value.trim(),
      password: password.value
    })

    setAuthToken(res.data.token, res.data.user)

    const redirectPath = (route.query.redirect as string) || '/'
    router.replace(redirectPath)
  } catch (err: any) {
    errorMessage.value = err?.response?.data?.error || '登录失败，请检查账号与密码'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  if (getAuthToken()) {
    const redirectPath = (route.query.redirect as string) || '/'
    router.replace(redirectPath)
  }
})
</script>
