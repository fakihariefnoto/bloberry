import { defineStore } from 'pinia'
import { api } from '../lib/api'

export interface Tenant {
  id: string
  name: string
  slug: string
  quota_bytes: number
  quota_objects: number
  used_bytes: number
  used_objects: number
  status: string
  default_storage_id?: string
  storage_engines?: string[]
  upload_policy?: { mode: string; extensions?: string[] }
}

export interface Membership {
  id: string
  user_id: string
  tenant_id: string
  role: string
}

export const useTenantStore = defineStore('tenant', {
  state: () => ({
    list: [] as Tenant[],
    currentId: localStorage.getItem('bloberry.tenant') || '',
    memberships: [] as Membership[],
  }),
  getters: {
    current(state): Tenant | undefined {
      return state.list.find((t) => t.id === state.currentId)
    },
    currentRole(state): string {
      return state.memberships.find((m) => m.tenant_id === state.currentId)?.role || ''
    },
  },
  actions: {
    async load() {
      this.list = await api.get<Tenant[]>('/tenants')
      if (!this.currentId && this.list.length) {
        this.currentId = this.list[0].id
        localStorage.setItem('bloberry.tenant', this.currentId)
      }
      // The stored project may no longer exist (deleted) or the list may not
      // contain it yet — always fall back to a valid id instead of showing
      // "no access" after a refresh.
      if (this.currentId && !this.list.some((t) => t.id === this.currentId) && this.list.length) {
        this.currentId = this.list[0].id
        localStorage.setItem('bloberry.tenant', this.currentId)
      }
    },
    switchTo(id: string) {
      this.currentId = id
      localStorage.setItem('bloberry.tenant', id)
    },
  },
})
