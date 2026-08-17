<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '../lib/api'
import { useTenantStore } from '../stores/tenant'
import { HardDrive } from 'lucide-vue-next'
import AppButton from '../components/ui/AppButton.vue'
import PageHeader from '../components/ui/PageHeader.vue'
import AppInput from '../components/ui/AppInput.vue'

interface Backend { id: string; name: string; driver: string; tenant_id?: string }

const tenants = useTenantStore()
const name = ref('')
const quotaBytes = ref('0')
const allEngines = ref<Backend[]>([])
const assignedIds = ref<string[]>([])
const defaultId = ref('')
const error = ref('')
const saved = ref('')
const loadingBackends = ref(false)

// Engines the tenant can actually use: install-level + tenant-owned + assigned.
const availableEngines = computed(() => {
  const own = allEngines.value.filter((b) => b.tenant_id === tenants.currentId || !b.tenant_id)
  const assigned = allEngines.value.filter((b) => assignedIds.value.includes(b.id))
  const seen = new Set<string>()
  return [...own, ...assigned].filter((b) => (seen.has(b.id) ? false : (seen.add(b.id), true)))
})

async function loadEngines() {
  loadingBackends.value = true
  try {
    allEngines.value = await api.get<Backend[]>('/admin/backends')
    const current = tenants.current
    if (current) {
      assignedIds.value = current.storage_engines || []
      defaultId.value = current.default_storage_id || ''
    }
  } catch { allEngines.value = [] } finally { loadingBackends.value = false }
}

onMounted(async () => {
  await tenants.load()
  if (tenants.current) {
    name.value = tenants.current.name
    quotaBytes.value = String(tenants.current.quota_bytes)
  }
  await loadEngines()
})

function toggleEngine(id: string, checked: boolean) {
  if (checked) {
    if (!assignedIds.value.includes(id)) assignedIds.value.push(id)
  } else {
    assignedIds.value = assignedIds.value.filter((x) => x !== id)
    if (defaultId.value === id) defaultId.value = ''
  }
}

async function save() {
  error.value = ''
  saved.value = ''
  try {
    const payload: Record<string, unknown> = {
      name: name.value,
      quota_bytes: Number(quotaBytes.value) || 0,
      storage_engines: assignedIds.value,
    }
    if (defaultId.value) payload.default_storage_id = defaultId.value
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

        <!-- Storage engines: pick which engines the tenant may use + the default -->
        <div class="flex flex-col gap-1.5">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Storage engines</label>
          <div v-if="loadingBackends" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] animate-pulse" />
          <div v-else class="divide-y divide-[var(--color-border)] overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)]">
            <label v-for="b in availableEngines" :key="b.id" class="flex cursor-pointer items-center gap-3 px-3 py-2.5 hover:bg-[var(--color-surface-raised)]">
              <input
                type="radio"
                name="default-engine"
                class="h-4 w-4 accent-[var(--color-primary)]"
                :checked="defaultId === b.id"
                :disabled="!assignedIds.includes(b.id)"
                @change="defaultId = b.id"
              />
              <input
                type="checkbox"
                class="h-4 w-4 accent-[var(--color-primary)]"
                :checked="assignedIds.includes(b.id)"
                @change="toggleEngine(b.id, ($event.target as HTMLInputElement).checked)"
              />
              <HardDrive class="h-4 w-4 shrink-0 text-[var(--color-primary)]" />
              <span class="flex-1 text-sm text-[var(--color-text)]">{{ b.name }}</span>
              <span class="rounded-full bg-[var(--color-primary-subtle)] px-2 py-0.5 text-xs text-[var(--color-primary)]">{{ b.driver }}</span>
            </label>
          </div>
          <p class="text-xs leading-relaxed text-[var(--color-text-muted)]">
            Check an engine to let this tenant use it. The radio picks the default for new uploads (shown in Files).
            Install-level and this tenant's own engines are always available.
          </p>
          <p v-if="!availableEngines.length" class="text-xs text-[var(--color-warning)]">No storage engines registered yet — ask a platform admin to add one.</p>
        </div>

        <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
        <p v-if="saved" class="text-xs text-[var(--color-success)]">{{ saved }}</p>
        <div class="flex justify-end"><AppButton @click="save">Save changes</AppButton></div>
      </div>
    </div>
  </div>
</template>
