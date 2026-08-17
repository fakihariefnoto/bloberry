<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Boxes, FolderOpen, Share2, ListChecks, KeyRound, Users, ScrollText,
  Gauge, Settings, Building2, HardDrive, BarChart3, LogOut, User, Smartphone, ChevronDown, Key, ArrowLeftRight,
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
  { name: 'api-keys', label: 'API keys', icon: Key },
  { name: 'shares', label: 'Shares', icon: Share2 },
  { name: 'jobs', label: 'Jobs', icon: ListChecks },
  { name: 'transfers', label: 'Transfers', icon: ArrowLeftRight },
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
        { name: 'admin-backends', label: 'Storage engines', icon: HardDrive },
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
  }
}

const pageTitle = computed(() => {
  const map: Record<string, string> = {
    files: 'Files', shares: 'Shares', jobs: 'Jobs', transfers: 'Transfers', applications: 'Applications',
    'api-keys': 'API keys',
    'application-detail': 'Application', members: 'Members', audit: 'Audit log',
    usage: 'Usage', 'tenant-settings': 'Tenant settings', profile: 'Profile',
    'account-settings': 'Account settings', 'pair-device': 'Pair a device',
    'admin-tenants': 'Tenants', 'admin-tenant-detail': 'Tenant',
    'admin-backends': 'Storage engines', 'admin-backend-detail': 'Backend',
    'admin-usage': 'Install usage',
  }
  return map[String(route.name)] || 'Bloberry'
})
</script>

<template>
  <div class="flex h-screen bg-[var(--color-background)]">
    <!-- Sidebar -->
    <aside class="flex w-[260px] flex-col border-r border-[var(--color-border)] bg-[var(--color-surface)]">
      <div class="flex h-16 items-center gap-2.5 border-b border-[var(--color-border)] px-5">
        <span class="flex h-8 w-8 items-center justify-center rounded-[var(--radius-md)] bg-[var(--color-primary)] shadow-[var(--shadow-sm)]">
          <Boxes class="h-4 w-4 text-[var(--color-on-primary)]" />
        </span>
        <div class="flex flex-col">
          <span class="text-sm font-semibold tracking-tight text-[var(--color-text)]">Bloberry</span>
        </div>
      </div>

      <div class="px-3 py-3">
        <select
          v-model="tenants.currentId"
          class="h-11 w-full appearance-none rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] px-3 pr-8 text-sm font-medium text-[var(--color-text)] focus:border-[var(--color-primary)] focus:outline-none"
          @change="switchTenant(tenants.currentId)"
        >
          <option v-for="t in tenants.list" :key="t.id" :value="t.id">{{ t.name }}</option>
        </select>
      </div>

      <nav class="flex-1 overflow-y-auto px-3 py-2">
        <p class="px-2 pb-1 pt-2 text-[11px] font-semibold uppercase tracking-wider text-[var(--color-text-muted)]">Workspace</p>
        <RouterLink
          v-for="item in mainNav"
          :key="item.name"
          :to="{ name: item.name }"
          class="relative flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2.5 text-sm"
          :class="isActive(item.name) ? 'bg-[var(--color-primary-subtle)] font-semibold text-[var(--color-primary)]' : 'text-[var(--color-text-muted)] hover:bg-[var(--color-surface-raised)] hover:text-[var(--color-text)]'"
        >
          <component :is="item.icon" class="h-5 w-5" />
          {{ item.label }}
          <span v-if="isActive(item.name)" class="absolute left-0 top-1/2 h-6 w-1 -translate-y-1/2 rounded-r-full bg-[var(--color-primary)]" />
        </RouterLink>

        <template v-if="adminNav.length">
          <p class="px-2 pb-1 pt-3 text-[11px] font-semibold uppercase tracking-wider text-[var(--color-text-muted)]">Administration</p>
          <RouterLink
            v-for="item in adminNav"
            :key="item.name"
            :to="{ name: item.name }"
            class="relative flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2.5 text-sm"
            :class="isActive(item.name) ? 'bg-[var(--color-primary-subtle)] font-semibold text-[var(--color-primary)]' : 'text-[var(--color-text-muted)] hover:bg-[var(--color-surface-raised)] hover:text-[var(--color-text)]'"
          >
            <component :is="item.icon" class="h-5 w-5" />
            {{ item.label }}
            <span v-if="isActive(item.name)" class="absolute left-0 top-1/2 h-6 w-1 -translate-y-1/2 rounded-r-full bg-[var(--color-primary)]" />
          </RouterLink>
        </template>

        <template v-if="platformNav.length">
          <p class="px-2 pb-1 pt-3 text-[11px] font-semibold uppercase tracking-wider text-[var(--color-text-muted)]">Platform</p>
          <RouterLink
            v-for="item in platformNav"
            :key="item.name"
            :to="{ name: item.name }"
            class="relative flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2.5 text-sm"
            :class="isActive(item.name) ? 'bg-[var(--color-primary-subtle)] font-semibold text-[var(--color-primary)]' : 'text-[var(--color-text-muted)] hover:bg-[var(--color-surface-raised)] hover:text-[var(--color-text)]'"
          >
            <component :is="item.icon" class="h-5 w-5" />
            {{ item.label }}
            <span v-if="isActive(item.name)" class="absolute left-0 top-1/2 h-6 w-1 -translate-y-1/2 rounded-r-full bg-[var(--color-primary)]" />
          </RouterLink>
        </template>
      </nav>

      <div class="border-t border-[var(--color-border)] p-3">
        <div class="flex items-center gap-2.5 rounded-[var(--radius-md)] px-2 py-2">
          <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-[var(--color-primary-subtle)] text-sm font-semibold text-[var(--color-primary)]">
            {{ (auth.user?.display_name || auth.user?.email || '?').charAt(0).toUpperCase() }}
          </span>
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-medium text-[var(--color-text)]">{{ auth.user?.display_name || 'Account' }}</p>
            <p class="truncate text-xs text-[var(--color-text-muted)]">{{ auth.user?.email }}</p>
          </div>
          <ChevronDown class="h-4 w-4 text-[var(--color-text-muted)]" />
        </div>
        <div class="mt-1 flex flex-col gap-0.5 border-t border-[var(--color-border)] pt-2">
          <RouterLink
            v-for="item in userNav"
            :key="item.name"
            :to="{ name: item.name }"
            class="flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2 text-sm text-[var(--color-text-muted)] hover:bg-[var(--color-surface-raised)] hover:text-[var(--color-text)]"
          >
            <component :is="item.icon" class="h-4 w-4" />
            {{ item.label }}
          </RouterLink>
          <button
            class="flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2 text-sm text-[var(--color-error)] hover:bg-[var(--color-surface-raised)]"
            @click="logout"
          >
            <LogOut class="h-4 w-4" />
            Log out
          </button>
        </div>
      </div>
    </aside>

    <!-- Main -->
    <div class="flex min-w-0 flex-1 flex-col">
      <header class="flex h-16 shrink-0 items-center justify-between border-b border-[var(--color-border)] bg-[var(--color-surface-raised)]/60 px-6 backdrop-blur">
        <h2 class="text-[15px] font-semibold tracking-tight text-[var(--color-text)]">{{ pageTitle }}</h2>
        <div class="flex items-center gap-2 text-xs text-[var(--color-text-muted)]">
          <span class="inline-flex items-center gap-1.5 rounded-full border border-[var(--color-border)] bg-[var(--color-surface)] px-2.5 py-1">
            <span class="h-1.5 w-1.5 rounded-full bg-[var(--color-success)]" />
            {{ tenants.current?.name || 'No tenant' }}
          </span>
        </div>
      </header>
      <main class="flex-1 overflow-y-auto">
        <div class="mx-auto max-w-6xl px-6 py-6">
          <RouterView :key="tenants.currentId" />
        </div>
      </main>
    </div>
  </div>
</template>
