<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, HardDrive, PauseCircle, CheckCircle2 } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'

interface Tenant {
  id: string; name: string; slug: string; status: string
  quota_bytes: number; quota_objects: number
  used_bytes: number; used_objects: number
  default_storage_id?: string
}
interface Backend { id: string; name: string; driver: string; health_status: string }

const route = useRoute()
const router = useRouter()
const tenantId = route.params.tenantId as string

const t = ref<Tenant | null>(null)
const backends = ref<Backend[]>([])
const storageId = ref('')
const quotaBytes = ref('0')
const quotaObjects = ref('0')
const saved = ref('')
const error = ref('')
const saving = ref(false)

async function load() {
  t.value = await api.get<Tenant>(`/tenants/${tenantId}`)
  if (t.value) {
    quotaBytes.value = String(t.value.quota_bytes)
    quotaObjects.value = String(t.value.quota_objects)
    storageId.value = t.value.default_storage_id || ''
  }
}
async function loadBackends() {
  try { backends.value = await api.get<Backend[]>('/admin/backends') } catch { backends.value = [] }
}
onMounted(() => { load(); loadBackends() })

async function saveAll() {
  error.value = ''
  saved.value = ''
  saving.value = true
  try {
    await api.patch(`/tenants/${tenantId}`, {
      quota_bytes: Number(quotaBytes.value) || 0,
      quota_objects: Number(quotaObjects.value) || 0,
      default_storage_id: storageId.value || undefined,
    })
    saved.value = 'Saved.'
    await load()
  } catch (e) { error.value = (e as Error).message } finally { saving.value = false }
}

async function toggleStatus() {
  await api.patch(`/tenants/${tenantId}`, { status: t.value?.status === 'suspended' ? 'active' : 'suspended' })
  await load()
}

function formatBytes(n: number) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (n >= 1024 && i < units.length - 1) { n /= 1024; i++ }
  return `${n.toFixed(n >= 10 || i === 0 ? 0 : 1)} ${units[i]}`
}

const usedPct = () => (t.value?.quota_bytes ? Math.round(((t.value.used_bytes || 0) / t.value.quota_bytes) * 100) : 0)
const back = () => router.push({ name: 'admin-tenants' })
</script>

<template>
  <div>
    <button class="mb-4 flex items-center gap-1 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="back">
      <ArrowLeft class="h-4 w-4" /> Back to projects
    </button>

    <div v-if="t">
      <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-[var(--color-text)]">{{ t.name }}</h1>
          <p class="mt-1 text-sm text-[var(--color-text-muted)]">/{{ t.slug }}</p>
        </div>
        <AppButton size="sm" :variant="t.status === 'suspended' ? 'secondary' : 'destructive'" @click="toggleStatus">
          <PauseCircle v-if="t.status !== 'suspended'" class="h-4 w-4" />
          <CheckCircle2 v-else class="h-4 w-4" />
          {{ t.status === 'suspended' ? 'Unsuspend' : 'Suspend' }}
        </AppButton>
      </div>

      <!-- Stat cards -->
      <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
          <p class="text-xs font-medium text-[var(--color-text-muted)]">Storage used</p>
          <p class="mt-2 text-2xl font-bold tracking-tight text-[var(--color-text)]">{{ formatBytes(t.used_bytes) }}</p>
          <p v-if="t.quota_bytes" class="mt-2 text-xs text-[var(--color-text-muted)]">
            {{ usedPct() }}% of {{ formatBytes(t.quota_bytes) }}
          </p>
          <div v-if="t.quota_bytes" class="mt-3 h-1.5 overflow-hidden rounded-full bg-[var(--color-border)]">
            <div
              class="h-full"
              :class="usedPct() > 80 ? 'bg-[var(--color-warning)]' : 'bg-[var(--color-primary)]'"
              :style="{ width: `${Math.min(100, usedPct())}%` }"
            />
          </div>
        </div>
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
          <p class="text-xs font-medium text-[var(--color-text-muted)]">Objects</p>
          <p class="mt-2 text-2xl font-bold tracking-tight text-[var(--color-text)]">{{ t.used_objects.toLocaleString() }}</p>
          <p v-if="t.quota_objects" class="mt-2 text-xs text-[var(--color-text-muted)]">quota {{ t.quota_objects.toLocaleString() }}</p>
        </div>
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
          <p class="text-xs font-medium text-[var(--color-text-muted)]">Status</p>
          <p class="mt-2 inline-flex items-center gap-2 text-2xl font-bold tracking-tight" :class="t.status === 'suspended' ? 'text-[var(--color-error)]' : 'text-[var(--color-success)]'">
            <span class="h-2.5 w-2.5 rounded-full" :class="t.status === 'suspended' ? 'bg-[var(--color-error)]' : 'bg-[var(--color-success)]'" />
            {{ t.status }}
          </p>
        </div>
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
          <p class="text-xs font-medium text-[var(--color-text-muted)]">Storage engine</p>
          <p class="mt-2 flex items-center gap-2 text-2xl font-bold tracking-tight text-[var(--color-text)]">
            <HardDrive class="h-5 w-5 text-[var(--color-primary)]" />
            {{ backends.find((b) => b.id === t?.default_storage_id)?.driver || 'none' }}
          </p>
          <p class="mt-2 truncate text-xs text-[var(--color-text-muted)]">{{ backends.find((b) => b.id === t?.default_storage_id)?.name || 'Not assigned' }}</p>
        </div>
      </div>

      <!-- Settings -->
      <div class="mt-6 max-w-lg rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
        <h2 class="text-lg font-semibold text-[var(--color-text)]">Configuration</h2>
        <div class="mt-5 flex flex-col gap-4">
          <div class="flex flex-col gap-1">
            <label class="text-xs font-medium text-[var(--color-text-muted)]">Storage engine</label>
            <select v-model="storageId" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
              <option value="">No engine assigned</option>
              <option v-for="b in backends" :key="b.id" :value="b.id">{{ b.name }} ({{ b.driver }})</option>
            </select>
            <p class="flex items-center gap-1.5 text-xs text-[var(--color-text-muted)]">
              <HardDrive class="h-3.5 w-3.5" />
              Applies to new objects — existing files keep their engine.
            </p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <AppInput v-model="quotaBytes" label="Quota (bytes, 0 = unlimited)" type="number" placeholder="0" />
            <AppInput v-model="quotaObjects" label="Object quota (0 = unlimited)" type="number" placeholder="0" />
          </div>

          <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
          <p v-if="saved" class="text-xs text-[var(--color-success)]">{{ saved }}</p>
          <div class="flex justify-end"><AppButton :loading="saving" @click="saveAll">Save changes</AppButton></div>
        </div>
      </div>
    </div>
  </div>
</template>
