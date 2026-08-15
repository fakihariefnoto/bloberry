<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

interface AuditEvent {
  id: string
  action: string
  principal_type: string
  principal_id: string
  target_type?: string
  target_id?: string
  ip?: string
  created_at: string
}

const events = ref<AuditEvent[]>([])
const loading = ref(false)
const action = ref('')
const limit = 50

async function load() {
  loading.value = true
  try {
    const q = new URLSearchParams()
    if (action.value) q.set('action', action.value)
    q.set('limit', String(limit))
    events.value = await api.get<AuditEvent[]>(`/audit?${q}`)
  } catch { events.value = [] } finally { loading.value = false }
}
onMounted(load)
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[var(--color-text)]">Audit log</h1>
        <p class="mt-1 text-sm text-[var(--color-text-muted)]">What happened in this tenant, in order.</p>
      </div>
      <div class="flex gap-2">
        <input v-model="action" placeholder="filter by action" class="h-9 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-xs" @keyup.enter="load" />
        <AppButton size="sm" @click="load">Apply</AppButton>
      </div>
    </div>

    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">When</th>
            <th class="px-3 py-2">Action</th>
            <th class="px-3 py-2">Principal</th>
            <th class="px-3 py-2">Target</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in events" :key="e.id" class="border-t border-[var(--color-border)]">
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ new Date(e.created_at).toLocaleString() }}</td>
            <td class="px-3 py-2.5 font-mono text-xs">{{ e.action }}</td>
            <td class="px-3 py-2.5 font-mono text-xs">{{ e.principal_id.slice(0, 8) }}…</td>
            <td class="px-3 py-2.5 font-mono text-xs">{{ e.target_type ? `${e.target_type}:${(e.target_id || '').slice(0, 8)}` : '—' }}</td>
          </tr>
          <tr v-if="!loading && !events.length">
            <td colspan="4" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No audit events in range</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
