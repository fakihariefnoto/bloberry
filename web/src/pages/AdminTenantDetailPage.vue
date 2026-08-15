<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

interface Tenant { id: string; name: string; slug: string; status: string; quota_bytes: number; quota_objects: number; used_bytes: number; used_objects: number; default_backend_id?: string }

const route = useRoute()
const router = useRouter()
const tenantId = route.params.tenantId as string
const t = ref<Tenant | null>(null)
const quotaBytes = ref(0)
const saved = ref(false)
const error = ref('')

async function load() {
  t.value = await api.get<Tenant>(`/tenants/${tenantId}`)
  quotaBytes.value = t.value?.quota_bytes || 0
}
onMounted(load)

async function suspend() {
  await api.patch(`/tenants/${tenantId}`, { status: t.value?.status === 'suspended' ? 'active' : 'suspended' })
  load()
}

async function saveQuota() {
  try {
    await api.patch(`/tenants/${tenantId}`, { quota_bytes: quotaBytes.value })
    saved.value = true
    load()
  } catch (e) { error.value = (e as Error).message }
}

const back = () => router.push({ name: 'admin-tenants' })
</script>

<template>
  <div>
    <button class="mb-4 flex items-center gap-1 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="back">
      <ArrowLeft class="h-4 w-4" /> Back to tenants
    </button>

    <div v-if="t" class="max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <h1 class="text-xl font-bold text-[var(--color-text)]">{{ t.name }}</h1>
      <p class="mt-1 text-sm text-[var(--color-text-muted)]">/{{ t.slug }}</p>

      <div class="mt-4 flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Quota (bytes)</label>
          <input v-model.number="quotaBytes" type="number" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        </div>
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <p v-if="saved" class="text-xs text-[var(--color-success)]">Saved.</p>
        <div class="flex gap-2">
          <AppButton size="sm" @click="saveQuota">Save quota</AppButton>
          <AppButton size="sm" variant="destructive" @click="suspend">
            {{ t.status === 'suspended' ? 'Unsuspend' : 'Suspend' }}
          </AppButton>
        </div>
        <p class="text-xs text-[var(--color-text-muted)]">
          {{ t.used_bytes }} bytes used across {{ t.used_objects }} objects. Status: <span :class="t.status === 'suspended' ? 'text-[var(--color-error)]' : 'text-[var(--color-success)]'">{{ t.status }}</span>
        </p>
      </div>
    </div>
  </div>
</template>
