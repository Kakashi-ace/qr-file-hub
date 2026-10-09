<template>
  <header class="w-full max-w-4xl mx-auto pt-6 px-4 mb-8">
    <div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white/80 backdrop-blur-md rounded-2xl p-4 border border-slate-200/80 shadow-xs">
      <router-link to="/" class="flex items-center gap-3 cursor-pointer group">
        <div class="w-10 h-10 rounded-xl bg-indigo-600 text-white flex items-center justify-center font-bold text-lg shadow-sm group-hover:scale-105 transition-transform">
          Q
        </div>
        <div class="flex items-center gap-2">
          <span class="text-base font-bold text-slate-900 tracking-tight">QRFile 共享中心</span>
          <span class="px-2 py-0.5 rounded-full text-[10px] font-mono font-semibold bg-slate-100 text-slate-600 border border-slate-200">ENTERPRISE</span>
        </div>
      </router-link>

      <nav class="flex items-center gap-1.5 bg-slate-100/80 p-1 rounded-xl border border-slate-200/60">
        <router-link
          to="/"
          class="flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-all"
          :class="$route.path === '/' ? 'bg-white text-slate-900 shadow-2xs font-bold' : 'text-slate-600 hover:text-slate-900'"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span>制作二维码</span>
        </router-link>
        <router-link
          to="/admin"
          class="flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-all"
          :class="$route.path === '/admin' ? 'bg-white text-slate-900 shadow-2xs font-bold' : 'text-slate-600 hover:text-slate-900'"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
          </svg>
          <span>后台文件管理</span>
        </router-link>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
const $route = useRoute()
</script>
