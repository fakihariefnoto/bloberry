<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'

interface Snap { tenant_id: string; bytes_stored: number; object_count: number; egress_bytes: number; estimated_cost: number; period: string }

const router = useRouter()
const snaps = ref<Snap[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try { snaps.value = await api.get<Snap[]>('/admin/usage') } catch { snaps.value = [] } finally { loading.value = false }
}
onMounted(load)

function formatBytes(n: number) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (n >= 1024 && i < units.length - 1) { n /= 1024; i++ }
  return `${n.toFixed(n >= 10 || i === 0 ? 0 : 1)} ${units[i]}`
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-[var(--color-text)]">Install usage</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">Latest snapshot per tenant.</p>

    <div class="mt-4 overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Tenant</th>
            <th class="px-3 py-2">Storage</th>
            <th class="px-3 py-2">Objects</th>
            <th class="px-3 py-2">Egress</th>
            <th class="px-3 py-2">Est. cost</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in snaps" :key="s.tenant_id + s.period" class="cursor-pointer border-t border-[var(--color-border)] hover:bg-[var(--color-surface)]" @click="router.push({ name: 'admin-tenant-detail', params: { tenantId: s.tenant_id } })">
            <td class="px-3 py-2.5 font-mono text-xs">{{ s.tenant_id.slice(0, 8) }}…</td>
            <td class="px-3 py-2.5">{{ formatBytes(s.bytes_stored) }}</td>
            <td class="px-3 py-2.5">{{ s.object_count }}</td>
            <td class="px-3 py-2.5">{{ formatBytes(s.egress_bytes) }}</td>
            <td class="px-3 py-2.5">{{ s.estimated_cost ? `$${s.estimated_cost.toFixed(2)}` : 'unknown' }}</td>
          </tr>
          <tr v-if="!loading && !snaps.length">
            <td colspan="5" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No usage snapshots yet.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
