<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../lib/api'
import { useTenantStore } from '../stores/tenant'

interface UsageSnapshot {
  bytes_stored: number
  object_count: number
  egress_bytes: number
  period: string
}

const tenants = useTenantStore()
const snap = ref<UsageSnapshot | null>(null)
const cost = ref<{ has_rate_card: boolean; total: number } | null>(null)

async function load() {
  try { snap.value = await api.get<UsageSnapshot>('/usage/me') } catch { snap.value = null }
  try { cost.value = await api.get<{ has_rate_card: boolean; total: number }>('/usage/estimated-cost') } catch { cost.value = null }
}
onMounted(load)

function formatBytes(n: number) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (n >= 1024 && i < units.length - 1) { n /= 1024; i++ }
  return `${n.toFixed(n >= 10 || i === 0 ? 0 : 1)} ${units[i]}`
}

const percent = (bytes: number, quota: number) => (quota > 0 ? Math.round((bytes / quota) * 100) : 0)
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-[var(--color-text)]">Usage</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">Storage and bandwidth for {{ tenants.current?.name }}.</p>

    <div class="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-3">
      <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
        <p class="text-xs font-medium text-[var(--color-text-muted)]">Storage</p>
        <p class="mt-1 text-3xl font-bold text-[var(--color-text)]">{{ formatBytes(snap?.bytes_stored || 0) }}</p>
        <div v-if="tenants.current?.quota_bytes" class="mt-3 h-1.5 overflow-hidden rounded-full bg-[var(--color-border)]">
          <div
            class="h-full"
            :class="percent(snap?.bytes_stored || 0, tenants.current.quota_bytes) > 80 ? 'bg-[var(--color-warning)]' : 'bg-[var(--color-primary)]'"
            :style="{ width: `${Math.min(100, percent(snap?.bytes_stored || 0, tenants.current.quota_bytes))}%` }"
          />
        </div>
        <p v-if="tenants.current?.quota_bytes" class="mt-1 text-xs text-[var(--color-text-muted)]">
          {{ percent(snap?.bytes_stored || 0, tenants.current.quota_bytes) }}% of quota
        </p>
      </div>

      <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
        <p class="text-xs font-medium text-[var(--color-text-muted)]">Objects</p>
        <p class="mt-1 text-3xl font-bold text-[var(--color-text)]">{{ snap?.object_count || 0 }}</p>
      </div>

      <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
        <p class="text-xs font-medium text-[var(--color-text-muted)]">Egress (estimated)</p>
        <p class="mt-1 text-3xl font-bold text-[var(--color-text)]">{{ formatBytes(snap?.egress_bytes || 0) }}</p>
        <p class="mt-1 text-xs text-[var(--color-text-muted)]">±10% — redirect downloads are estimated</p>
      </div>
    </div>

    <div class="mt-4 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <p class="text-xs font-medium text-[var(--color-text-muted)]">Estimated monthly cost</p>
      <p class="mt-1 text-3xl font-bold text-[var(--color-text)]">
        {{ cost?.has_rate_card ? `$${(cost.total || 0).toFixed(2)}` : 'Unknown' }}
      </p>
      <p class="mt-1 text-xs text-[var(--color-text-muted)]">
        {{ cost?.has_rate_card ? 'From your storage backend rate card.' : 'No rate card configured — cost is never shown as zero.' }}
      </p>
    </div>
  </div>
</template>
