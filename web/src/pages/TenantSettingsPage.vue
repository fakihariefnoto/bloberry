<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../lib/api'
import { useTenantStore } from '../stores/tenant'
import { HardDrive } from 'lucide-vue-next'
import AppButton from '../components/ui/AppButton.vue'
import PageHeader from '../components/ui/PageHeader.vue'
import AppInput from '../components/ui/AppInput.vue'

interface Backend { id: string; name: string; driver: string; health_status: string }

const tenants = useTenantStore()
const name = ref('')
const quotaBytes = ref('0')
const backends = ref<Backend[]>([])
const backendId = ref('')
const error = ref('')
const saved = ref('')
const loadingBackends = ref(false)

async function loadBackends() {
  loadingBackends.value = true
  try {
    backends.value = await api.get<Backend[]>('/admin/backends')
    if (tenants.current?.default_backend_id) {
      backendId.value = tenants.current.default_backend_id
    }
  } catch { backends.value = [] } finally { loadingBackends.value = false }
}

onMounted(async () => {
  await tenants.load()
  if (tenants.current) {
    name.value = tenants.current.name
    quotaBytes.value = String(tenants.current.quota_bytes)
  }
  await loadBackends()
})

async function save() {
  error.value = ''
  saved.value = ''
  try {
    const payload: Record<string, unknown> = { name: name.value, quota_bytes: Number(quotaBytes.value) || 0 }
    if (backendId.value) payload.default_backend_id = backendId.value
    await api.patch(`/tenants/${tenants.currentId}`, payload)
    saved.value = 'Saved.'
    await tenants.load()
  } catch (e) { error.value = (e as Error).message }
}
</script>

<template>
  <div>
    <PageHeader title="Tenant settings" :description="`Configuration for ${tenants.current?.name}.`" />

    <div class="mt-4 max-w-lg rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <div class="flex flex-col gap-4">
        <AppInput v-model="name" label="Tenant name" placeholder="Acme Corp" />
        <AppInput v-model="quotaBytes" label="Quota (bytes, 0 = unlimited)" type="number" placeholder="0" />

        <!-- Storage backend assignment -->
        <div class="flex flex-col gap-1.5">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Storage backend</label>
          <div v-if="loadingBackends" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] animate-pulse" />
          <select
            v-else
            v-model="backendId"
            class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2"
          >
            <option value="">No backend assigned</option>
            <option v-for="b in backends" :key="b.id" :value="b.id">{{ b.name }} ({{ b.driver }})</option>
          </select>
          <p class="flex items-center gap-1.5 text-xs text-[var(--color-text-muted)]">
            <HardDrive class="h-3.5 w-3.5" />
            Applies to new objects. Existing files keep their current backend — switching never strands them.
          </p>
          <p v-if="!backends.length" class="text-xs text-[var(--color-warning)]">No storage backends registered yet — ask a platform admin to add one.</p>
        </div>

        <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
        <p v-if="saved" class="text-xs text-[var(--color-success)]">{{ saved }}</p>
        <div class="flex justify-end"><AppButton @click="save">Save changes</AppButton></div>
      </div>
    </div>
  </div>
</template>
