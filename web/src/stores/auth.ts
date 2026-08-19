import { defineStore } from 'pinia'
import { api, setToken, tryRefresh } from '../lib/api'

interface User {
  id: string
  email: string
  display_name: string
  platform_role?: string
  email_verified?: boolean
  settings?: { default_tenant_id?: string; locale?: string; notifications_enabled?: boolean }
}

interface AuthData {
  access_token: string
  refresh_token: string
  expires_in?: number
  user?: User
  totp_required?: boolean
  pending?: string
}

const ACCESS_KEY = 'bloberry.access'
const REFRESH_KEY = 'bloberry.refresh'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    accessToken: localStorage.getItem(ACCESS_KEY) || null,
    refreshToken: localStorage.getItem(REFRESH_KEY) || null,
    smtpConfigured: false,
  }),
  getters: {
    isAuthenticated: (s) => !!s.accessToken,
  },
  actions: {
    applySession(data: AuthData) {
      this.accessToken = data.access_token
      this.refreshToken = data.refresh_token
      this.user = data.user || null
      localStorage.setItem(ACCESS_KEY, data.access_token)
      localStorage.setItem(REFRESH_KEY, data.refresh_token)
      setToken(data.access_token)
    },
    clear() {
      this.user = null
      this.accessToken = null
      this.refreshToken = null
      localStorage.removeItem(ACCESS_KEY)
      localStorage.removeItem(REFRESH_KEY)
      setToken(null)
    },
    async login(email: string, password: string) {
      const data = await api.post<AuthData>('/auth/login', { email, password, platform: 'web' })
      if (data.totp_required) return { totpRequired: true, pending: data.pending }
      this.applySession(data)
      return { totpRequired: false }
    },
    async verifyTotp(pending: string, code: string) {
      const data = await api.post<AuthData>('/auth/login/verify-totp', { pending, code })
      this.applySession(data)
    },
    async requestOtp(email: string) {
      await api.post('/auth/otp/request', { email })
    },
    async verifyOtp(email: string, code: string) {
      const data = await api.post<AuthData>('/auth/otp/verify', { email, code, platform: 'web' })
      this.applySession(data)
    },
    async signup(inviteToken: string, email: string, password: string, displayName: string) {
      const data = await api.post<AuthData>('/auth/signup', {
        invite_token: inviteToken, email, password, display_name: displayName, platform: 'web',
      })
      this.applySession(data)
    },
    async refresh() {
      if (!this.refreshToken) return false
      try {
        const data = await api.post<AuthData>('/auth/refresh', { refresh_token: this.refreshToken })
        this.applySession(data)
        return true
      } catch {
        this.clear()
        return false
      }
    },
    async logout() {
      try {
        if (this.refreshToken) await api.post('/auth/logout', { refresh_token: this.refreshToken })
      } finally {
        this.clear()
      }
    },
    async me() {
      const u = await api.get<User>('/users/me')
      this.user = u
      return u
    },
    async loadSmtpConfig() {
      try {
        const st = await api.get<{ smtp_configured?: boolean }>('/setup/status')
        this.smtpConfigured = st?.smtp_configured === true
      } catch {
        this.smtpConfigured = false
      }
    },
    bootstrap() {
      if (this.accessToken) setToken(this.accessToken)
    },
    // Called once at app boot: if a stored access token is missing or expired,
    // rotate it via the refresh token so a page refresh never lands on a 401.
    // Uses the interceptor's single-flight tryRefresh so this never races with
    // concurrent 401-triggered refreshes (which would consume the rotating
    // refresh token twice).
    async restoreSession() {
      if (!this.accessToken) return
      setToken(this.accessToken)
      if (isJwtExpired(this.accessToken)) {
        const ok = await tryRefresh()
        if (!ok) this.clear()
      }
    },
    syncFromStorage() {
      this.accessToken = localStorage.getItem(ACCESS_KEY)
      this.refreshToken = localStorage.getItem(REFRESH_KEY)
      if (this.accessToken) setToken(this.accessToken)
    },
  },
})

// Decodes the exp claim of a JWT without verifying (verification happens on
// the server; this only decides whether to refresh proactively).
function isJwtExpired(token: string): boolean {
  try {
    const payload = token.split('.')[1]
    if (!payload) return false
    const b64 = payload.replace(/-/g, '+').replace(/_/g, '/')
    const padded = b64 + '='.repeat((4 - (b64.length % 4)) % 4)
    const bin = atob(padded)
    const bytes = Uint8Array.from(bin, (c) => c.charCodeAt(0))
    const decoded = JSON.parse(new TextDecoder().decode(bytes))
    if (typeof decoded.exp !== 'number') return false
    return decoded.exp * 1000 <= Date.now()
  } catch {
    return false
  }
}
