<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  FolderOpen, Share2, ListChecks, KeyRound, Users, ScrollText,
  Gauge, Settings, Building2, HardDrive, BarChart3, LogOut, User, Smartphone,
} from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useTenantStore } from '../stores/tenant'

const auth = useAuthStore()
const tenants = useTenantStore()
const router = useRouter()
const route = useRoute()

onMounted(async () => {
  if (!auth.user) {
    try { await auth.me() } catch { /* ignore */ }
  }
  await tenants.load()
})

const isPlatformAdmin = computed(() => !!auth.user?.platform_role)
const role = computed(() => tenants.currentRole)

const mainNav = [
  { name: 'files', label: 'Files', icon: FolderOpen },
  { name: 'shares', label: 'Shares', icon: Share2 },
  { name: 'jobs', label: 'Jobs', icon: ListChecks },
]

const adminNav = computed(() =>
  ['tenant_admin', 'tenant_owner'].includes(role.value)
    ? [
        { name: 'applications', label: 'Applications', icon: KeyRound },
        { name: 'members', label: 'Members', icon: Users },
        { name: 'audit', label: 'Audit log', icon: ScrollText },
        { name: 'usage', label: 'Usage', icon: Gauge },
        { name: 'tenant-settings', label: 'Settings', icon: Settings },
      ]
    : [],
)

const platformNav = computed(() =>
  isPlatformAdmin.value
    ? [
        { name: 'admin-tenants', label: 'Tenants', icon: Building2 },
        { name: 'admin-backends', label: 'Storage backends', icon: HardDrive },
        { name: 'admin-usage', label: 'Install usage', icon: BarChart3 },
      ]
    : [],
)

const userNav = [
  { name: 'profile', label: 'Profile', icon: User },
  { name: 'account-settings', label: 'Account settings', icon: Settings },
  { name: 'pair-device', label: 'Pair a device', icon: Smartphone },
]

function isActive(name: string) {
  return route.name === name
}

async function logout() {
  await auth.logout()
  router.push({ name: 'login' })
}

function switchTenant(id: string) {
  if (id !== tenants.currentId) {
    tenants.switchTo(id)
    router.push({ name: 'files' })
  }
}
</script>

<template>
  <div class="flex h-screen bg-[var(--color-background)]">
    <!-- Sidebar -->
    <aside class="flex w-[260px] flex-col border-r border-[var(--color-border)] bg-[var(--color-surface)]">
      <div class="px-4 py-3">
        <select
          v-model="tenants.currentId"
          class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] px-3 text-sm"
          @change="switchTenant(tenants.currentId)"
        >
          <option v-for="t in tenants.list" :key="t.id" :value="t.id">{{ t.name }}</option>
        </select>
      </div>

      <nav class="flex-1 overflow-y-auto px-2 py-2">
        <p class="px-2 pb-1 pt-2 text-[11px] font-semibold uppercase tracking-wide text-[var(--color-text-muted)]">Workspace</p>
        <RouterLink
          v-for="item in mainNav"
          :key="item.name"
          :to="{ name: item.name }"
          class="flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2.5 text-sm"
          :class="isActive(item.name) ? 'bg-[var(--color-primary-subtle)] font-semibold text-[var(--color-primary)]' : 'text-[var(--color-text-muted)] hover:bg-[var(--color-surface-raised)]'"
        >
          <component :is="item.icon" class="h-5 w-5" />
          {{ item.label }}
        </RouterLink>

        <template v-if="adminNav.length">
          <p class="px-2 pb-1 pt-3 text-[11px] font-semibold uppercase tracking-wide text-[var(--color-text-muted)]">Administration</p>
          <RouterLink
            v-for="item in adminNav"
            :key="item.name"
            :to="{ name: item.name }"
            class="flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2.5 text-sm"
            :class="isActive(item.name) ? 'bg-[var(--color-primary-subtle)] font-semibold text-[var(--color-primary)]' : 'text-[var(--color-text-muted)] hover:bg-[var(--color-surface-raised)]'"
          >
            <component :is="item.icon" class="h-5 w-5" />
            {{ item.label }}
          </RouterLink>
        </template>

        <template v-if="platformNav.length">
          <p class="px-2 pb-1 pt-3 text-[11px] font-semibold uppercase tracking-wide text-[var(--color-text-muted)]">Platform</p>
          <RouterLink
            v-for="item in platformNav"
            :key="item.name"
            :to="{ name: item.name }"
            class="flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2.5 text-sm"
            :class="isActive(item.name) ? 'bg-[var(--color-primary-subtle)] font-semibold text-[var(--color-primary)]' : 'text-[var(--color-text-muted)] hover:bg-[var(--color-surface-raised)]'"
          >
            <component :is="item.icon" class="h-5 w-5" />
            {{ item.label }}
          </RouterLink>
        </template>
      </nav>

      <div class="border-t border-[var(--color-border)] px-2 py-2">
        <p class="truncate px-3 pb-1 text-sm font-semibold text-[var(--color-text)]">{{ auth.user?.display_name || auth.user?.email }}</p>
        <RouterLink
          v-for="item in userNav"
          :key="item.name"
          :to="{ name: item.name }"
          class="flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2 text-sm text-[var(--color-text-muted)] hover:bg-[var(--color-surface-raised)]"
        >
          <component :is="item.icon" class="h-4 w-4" />
          {{ item.label }}
        </RouterLink>
        <button
          class="flex w-full items-center gap-3 rounded-[var(--radius-md)] px-3 py-2 text-sm text-[var(--color-error)] hover:bg-[var(--color-surface-raised)]"
          @click="logout"
        >
          <LogOut class="h-4 w-4" />
          Log out
        </button>
      </div>
    </aside>

    <!-- Main -->
    <main class="flex-1 overflow-y-auto">
      <div class="mx-auto max-w-6xl px-6 py-6">
        <RouterView />
      </div>
    </main>
  </div>
</template>
