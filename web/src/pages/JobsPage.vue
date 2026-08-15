<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { api } from '../lib/api'

interface Job {
  id: string
  kind: string
  state: string
  progress_done: number
  progress_total: number
  failure_message?: string
  created_at: string
}

const jobs = ref<Job[]>([])
const loading = ref(false)
let timer: number | undefined

async function load() {
  loading.value = true
  try {
    jobs.value = await api.get<Job[]>('/jobs')
  } catch {
    jobs.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
  timer = window.setInterval(load, 5000)
})
onUnmounted(() => window.clearInterval(timer))

function stateClass(s: string) {
  if (s === 'succeeded') return 'text-[var(--color-success)]'
  if (s === 'failed') return 'text-[var(--color-error)]'
  if (s === 'running') return 'text-[var(--color-primary)]'
  return 'text-[var(--color-warning)]'
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-[var(--color-text)]">Jobs</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">Extractions, bundles and large deletes, live.</p>

    <div class="mt-4 overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Type</th>
            <th class="px-3 py-2">State</th>
            <th class="px-3 py-2">Progress</th>
            <th class="px-3 py-2">Created</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="j in jobs" :key="j.id" class="border-t border-[var(--color-border)]">
            <td class="px-3 py-2.5 capitalize">{{ j.kind.replace('_', ' ') }}</td>
            <td class="px-3 py-2.5">
              <span :class="stateClass(j.state)">{{ j.state }}</span>
              <p v-if="j.failure_message" class="text-xs text-[var(--color-error)]">{{ j.failure_message }}</p>
            </td>
            <td class="px-3 py-2.5">
              <div class="flex items-center gap-2">
                <div class="h-1.5 w-32 overflow-hidden rounded-full bg-[var(--color-border)]">
                  <div
                    v-if="j.progress_total > 0"
                    class="h-full bg-[var(--color-primary)]"
                    :style="{ width: `${Math.round((j.progress_done / j.progress_total) * 100)}%` }"
                  />
                </div>
                <span v-if="j.progress_total > 0" class="text-xs text-[var(--color-text-muted)]">{{ j.progress_done }}/{{ j.progress_total }}</span>
              </div>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ new Date(j.created_at).toLocaleString() }}</td>
          </tr>
          <tr v-if="!loading && !jobs.length">
            <td colspan="4" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No jobs yet</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
