<template>
  <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs">
    <div class="relative w-full max-w-sm rounded-3xl bg-white p-6 shadow-2xl border border-slate-200/90 text-center animate-in fade-in zoom-in duration-200">
      <div class="w-12 h-12 rounded-2xl bg-indigo-50 text-indigo-600 flex items-center justify-center mx-auto mb-4 border border-indigo-100">
        <svg class="w-6 h-6 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <h3 class="text-base font-bold text-slate-900 mb-1">{{ title }}</h3>
      <p class="text-xs text-slate-500 mb-5">{{ message }}</p>

      <div class="w-full bg-slate-100 rounded-full h-2.5 overflow-hidden mb-2">
        <div
          class="bg-indigo-600 h-2.5 rounded-full transition-all duration-300"
          :style="{ width: `${progress}%` }"
        ></div>
      </div>

      <div class="flex justify-between items-center text-[11px] font-mono font-medium text-slate-400">
        <span>执行中...</span>
        <span>{{ progress }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  visible: boolean
  title: string
  message: string
  progress: number
}>()
</script>
