<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../lib/api'
import { useTenantStore } from '../stores/tenant'
import AppButton from '../components/ui/AppButton.vue'

const tenants = useTenantStore()
const name = ref('')
const quotaBytes = ref(0)
const error = ref('')
const saved = ref(false)

onMounted(() => {
  if (tenants.current) {
    name.value = tenants.current.name
    quotaBytes.value = tenants.current.quota_bytes
  }
})

async function save() {
  error.value = ''
  saved.value = false
  try {
    await api.patch(`/tenants/${tenants.currentId}`, { name: name.value, quota_bytes: quotaBytes.value })
    saved.value = true
    await tenants.load()
  } catch (e) { error.value = (e as Error).message }
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-[var(--color-text)]">Tenant settings</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">Configuration for {{ tenants.current?.name }}.</p>

    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Tenant name</label>
          <input v-model="name" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Quota (bytes, 0 = unlimited)</label>
          <input v-model.number="quotaBytes" type="number" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        </div>
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <p v-if="saved" class="text-xs text-[var(--color-success)]">Saved.</p>
        <div><AppButton @click="save">Save changes</AppButton></div>
      </div>
    </div>
  </div>
</template>
