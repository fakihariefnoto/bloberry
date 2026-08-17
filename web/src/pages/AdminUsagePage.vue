<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'
import PageHeader from '../components/ui/PageHeader.vue'

interface Snap {
  tenant_id: string
  name: string
  slug: string
  bytes_stored: number
  object_count: number
  storage_cost: number
  has_rate_card: boolean
}

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
    <PageHeader title="Install usage" description="Live storage and object counts per project." />

    <div class="mt-4 overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Project</th>
            <th class="px-3 py-2">Storage</th>
            <th class="px-3 py-2">Objects</th>
            <th class="px-3 py-2">Est. cost</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in snaps" :key="s.tenant_id" class="cursor-pointer border-t border-[var(--color-border)] hover:bg-[var(--color-surface)]" @click="router.push({ name: 'admin-tenant-detail', params: { tenantId: s.tenant_id } })">
            <td class="px-3 py-2.5"><span class="font-medium text-[var(--color-text)]">{{ s.name }}</span> <span class="text-xs text-[var(--color-text-muted)]">/{{ s.slug }}</span></td>
            <td class="px-3 py-2.5">{{ formatBytes(s.bytes_stored) }}</td>
            <td class="px-3 py-2.5">{{ s.object_count }}</td>
            <td class="px-3 py-2.5">
              <span v-if="s.has_rate_card">{{ `$${s.storage_cost.toFixed(2)}` }}</span>
              <span v-else class="text-[var(--color-text-muted)]">no rate card</span>
            </td>
          </tr>
          <tr v-if="!loading && !snaps.length">
            <td colspan="4" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No projects yet.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
