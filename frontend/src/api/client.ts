import axios from 'axios'

export function getAuthToken(): string {
  return localStorage.getItem('qr_auth_token') || localStorage.getItem('qr_admin_token') || ''
}

export function setAuthToken(token: string, user?: any) {
  localStorage.setItem('qr_auth_token', token)
  localStorage.setItem('qr_admin_token', token)
  if (user) {
    localStorage.setItem('qr_auth_user', JSON.stringify(user))
  }
}

export function clearAuthToken() {
  localStorage.removeItem('qr_auth_token')
  localStorage.removeItem('qr_admin_token')
  localStorage.removeItem('qr_auth_user')
}

export function getCurrentUser(): { id?: number; username?: string; role?: string } | null {
  const userStr = localStorage.getItem('qr_auth_user')
  if (!userStr) return null
  try {
    return JSON.parse(userStr)
  } catch {
    return null
  }
}

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000
})

apiClient.interceptors.request.use((config) => {
  const token = getAuthToken()
  if (token) {
    if (config.headers && typeof config.headers.set === 'function') {
      config.headers.set('Authorization', `Bearer ${token}`)
    } else {
      config.headers = config.headers || {}
      config.headers['Authorization'] = `Bearer ${token}`
    }
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error?.response?.status === 401) {
      const path = window.location.pathname
      if (path !== '/login' && !path.startsWith('/s/')) {
        clearAuthToken()
        window.location.href = `/login?redirect=${encodeURIComponent(path + window.location.search)}`
      }
    }
    return Promise.reject(error)
  }
)

export interface PresignUploadResponse {
  file_id: string
  r2_key: string
  upload_url: string
}

export interface ShareResultResponse {
  file_id: string
  file_name?: string
  target_url?: string
  direct_url: string
  bridge_url: string
  expire_at?: string
}

export interface FileItem {
  id: number
  file_id: string
  file_name: string
  file_size: number
  mime_type: string
  mode: string
  target_url?: string
  text_content?: string
  allow_download: boolean
  bridge_domain: string
  expire_type: string
  expire_at?: string
  view_count: number
  download_count: number
  created_by?: string
  created_at: string
}

export interface BridgeItem {
  id: number
  domain: string
  is_active: boolean
  is_primary: boolean
  is_custom: boolean
  wechat_status: string
  qq_status: string
  latency_ms: number
  last_checked_at?: string
}

export interface AdminStats {
  total_files: number
  total_storage: number
  total_downloads: number
  expired_count: number
}
