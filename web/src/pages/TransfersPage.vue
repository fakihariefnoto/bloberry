<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { ArrowRight, RefreshCw, HardDrive } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'
import PageHeader from '../components/ui/PageHeader.vue'

interface Backend { id: string; name: string; driver: string }
interface Job { id: string; kind: string; state: string; progress_done: number; progress_total: number; failure_message?: string; created_at: string; payload?: Record<string, any> }

const backends = ref<Backend[]>([])
const source = ref('')
const target = ref('')
const error = ref('')
const starting = ref(false)
const jobs = ref<Job[]>([])
let timer: number | undefined

async function loadBackends() {
  try {
    backends.value = await api.get<Backend[]>('/backends')
    if (!source.value && backends.value.length) source.value = backends.value[0].id
    if (!target.value && backends.value.length > 1) target.value = backends.value[1].id
  } catch { backends.value = [] }
}
async function loadJobs() {
  try {
    const all = await api.get<Job[]>('/jobs')
    jobs.value = (all || []).filter((j) => j.kind === 'transfer')
  } catch { jobs.value = [] }
}
onMounted(() => { loadBackends(); loadJobs(); timer = window.setInterval(loadJobs, 3000) })
onUnmounted(() => window.clearInterval(timer))

const engineName = (id: string) => backends.value.find((b) => b.id === id)?.name || id.slice(0, 8)

async function startTransfer() {
  error.value = ''
  starting.value = true
  try {
    await api.post('/transfers', { source_storage_id: source.value, target_storage_id: target.value })
    loadJobs()
  } catch (e) { error.value = (e as Error).message } finally { starting.value = false }
}

const pct = (j: Job) => (j.progress_total > 0 ? Math.round((j.progress_done / j.progress_total) * 100) : 0)
</script>

<template>
  <div>
    <PageHeader title="Transfers" description="Copy objects from one storage engine to another. Runs as a background job." />

    <!-- Transfer form -->
    <div class="mb-6 max-w-2xl rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <div class="grid items-end gap-4 sm:grid-cols-[1fr_auto_1fr]">
        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">From</label>
          <select v-model="source" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
            <option v-for="b in backends" :key="b.id" :value="b.id">{{ b.name }} ({{ b.driver }})</option>
          </select>
        </div>
        <div class="flex h-12 items-center justify-center pb-1 text-[var(--color-primary)]">
          <ArrowRight class="h-5 w-5" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">To</label>
          <select v-model="target" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
            <option v-for="b in backends" :key="b.id" :value="b.id">{{ b.name }} ({{ b.driver }})</option>
          </select>
        </div>
      </div>
      <p v-if="error" class="mt-3 rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
      <div class="mt-4 flex justify-end">
        <AppButton :loading="starting" :disabled="!source || !target || source === target" @click="startTransfer">
          <RefreshCw class="mr-2 h-4 w-4" /> Start transfer
        </AppButton>
      </div>
      <p class="mt-3 text-xs text-[var(--color-text-muted)]">
        Copies every active object from the source engine to the target, then re-points each file's metadata. Existing files keep working throughout.
      </p>
    </div>

    <!-- Transfers list -->
    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Source → Target</th>
            <th class="px-3 py-2">State</th>
            <th class="px-3 py-2">Progress</th>
            <th class="px-3 py-2">Started</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="j in jobs" :key="j.id" class="border-t border-[var(--color-border)]">
            <td class="px-3 py-2.5">
              <span class="inline-flex items-center gap-1.5 text-xs">
                <HardDrive class="h-3.5 w-3.5 text-[var(--color-text-muted)]" />
                {{ engineName((j.payload as any)?.source_storage_id || '') }}
                <ArrowRight class="h-3 w-3 text-[var(--color-text-muted)]" />
                {{ engineName((j.payload as any)?.target_storage_id || '') }}
              </span>
            </td>
            <td class="px-3 py-2.5">
              <span :class="j.state === 'succeeded' ? 'text-[var(--color-success)]' : j.state === 'failed' ? 'text-[var(--color-error)]' : j.state === 'running' ? 'text-[var(--color-primary)]' : 'text-[var(--color-warning)]'">{{ j.state }}</span>
              <p v-if="j.failure_message" class="text-xs text-[var(--color-error)]">{{ j.failure_message }}</p>
            </td>
            <td class="px-3 py-2.5">
              <div class="flex items-center gap-2">
                <div class="h-1.5 w-32 overflow-hidden rounded-full bg-[var(--color-border)]">
                  <div v-if="j.progress_total > 0" class="h-full bg-[var(--color-primary)]" :style="{ width: `${pct(j)}%` }" />
                </div>
                <span v-if="j.progress_total > 0" class="text-xs text-[var(--color-text-muted)]">{{ j.progress_done }}/{{ j.progress_total }}</span>
                <span v-else class="text-xs text-[var(--color-text-muted)]">—</span>
              </div>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ new Date(j.created_at).toLocaleString() }}</td>
          </tr>
          <tr v-if="!jobs.length">
            <td colspan="4" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No transfers yet.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
