<template>
  <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs">
    <div class="relative w-full max-w-sm rounded-3xl bg-white p-6 shadow-2xl border border-slate-200/90 animate-in fade-in zoom-in duration-200">
      <button
        @click="$emit('close')"
        class="absolute top-4 right-4 w-8 h-8 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-500 hover:text-slate-700 flex items-center justify-center transition cursor-pointer"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <div class="text-center mb-4">
        <span class="inline-block px-3 py-1 rounded-full text-[10px] font-mono font-bold tracking-wider bg-indigo-50 text-indigo-600 border border-indigo-200/80 mb-2">
          QR ASSET VIEWER
        </span>
        <h3 class="text-sm font-bold text-slate-900 truncate max-w-[240px] mx-auto">
          {{ title || '分享已就绪' }}
        </h3>
      </div>

      <div class="flex justify-center p-4 rounded-2xl bg-white border border-slate-100 shadow-inner mb-4">
        <qrcode-vue :value="currentShareUrl" :size="200" level="H" class="rounded-lg" />
      </div>

      <div class="flex items-center p-1 rounded-xl bg-slate-100 border border-slate-200 mb-3">
        <button
          type="button"
          @click="activeLinkType = 'bridge'"
          class="flex-1 py-1.5 rounded-lg text-xs font-bold transition flex items-center justify-center gap-1 cursor-pointer"
          :class="activeLinkType === 'bridge' ? 'bg-emerald-600 text-white shadow-2xs' : 'text-slate-600 hover:text-slate-900'"
        >
          <span>🛡️ 微信防封</span>
        </button>
        <button
          type="button"
          @click="activeLinkType = 'direct'"
          class="flex-1 py-1.5 rounded-lg text-xs font-bold transition flex items-center justify-center gap-1 cursor-pointer"
          :class="activeLinkType === 'direct' ? 'bg-indigo-600 text-white shadow-2xs' : 'text-slate-600 hover:text-slate-900'"
        >
          <span>⚡ 原站直连</span>
        </button>
      </div>

      <div class="mb-5">
        <div class="flex items-center rounded-xl border border-slate-200 bg-slate-50/80 overflow-hidden focus-within:border-indigo-500 focus-within:ring-1 focus-within:ring-indigo-500">
          <input
            type="text"
            readonly
            :value="currentShareUrl"
            class="flex-1 px-3 py-2 text-xs text-slate-700 bg-transparent outline-none font-mono truncate"
          />
          <button
            type="button"
            @click="copyUrl"
            class="px-3.5 py-2 text-xs font-bold bg-indigo-600 hover:bg-indigo-700 text-white transition active:scale-95 cursor-pointer"
          >
            {{ copied ? '已复制' : '复制' }}
          </button>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <button
          type="button"
          @click="saveQrImage"
          class="w-full py-2.5 px-3 rounded-xl bg-slate-100 hover:bg-slate-200 border border-slate-300 text-xs font-bold text-slate-800 transition active:scale-95 cursor-pointer text-center"
        >
          保存二维码
        </button>
        <a
          :href="currentShareUrl"
          target="_blank"
          class="w-full py-2.5 px-3 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold transition active:scale-95 text-center flex items-center justify-center cursor-pointer shadow-sm"
        >
          打开落地页
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import QrcodeVue from 'qrcode.vue'

const props = defineProps<{
  visible: boolean
  title: string
  directUrl: string
  bridgeUrl: string
}>()

defineEmits(['close'])

const activeLinkType = ref<'bridge' | 'direct'>('bridge')
const copied = ref(false)

const currentShareUrl = computed(() => {
  return activeLinkType.value === 'bridge' && props.bridgeUrl ? props.bridgeUrl : props.directUrl
})

function copyUrl() {
  if (!currentShareUrl.value) return
  navigator.clipboard.writeText(currentShareUrl.value).then(() => {
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  })
}

function saveQrImage() {
  const canvas = document.querySelector('canvas')
  if (!canvas) return
  const link = document.createElement('a')
  link.download = `${props.title || 'qrcode'}.png`
  link.href = canvas.toDataURL('image/png')
  link.click()
}
</script>
