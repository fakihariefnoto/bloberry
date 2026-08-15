<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Activity } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

interface Backend {
  id: string; name: string; driver: string; config: Record<string, unknown>
  health_status: string; health_error?: string; health_checked_at?: string
  rate_card?: Record<string, number>
}

const route = useRoute()
const router = useRouter()
const backendId = route.params.backendId as string
const b = ref<Backend | null>(null)
const error = ref('')

async function load() {
  try { b.value = await api.get<Backend>(`/admin/backends/${backendId}`) } catch (e) { error.value = (e as Error).message }
}
onMounted(load)

async function checkHealth() {
  try {
    b.value = await api.post<Backend>(`/admin/backends/${backendId}/health`)
  } catch (e) { error.value = (e as Error).message }
}

async function remove() {
  try {
    await api.delete(`/admin/backends/${backendId}`)
    router.push({ name: 'admin-backends' })
  } catch (e) { error.value = (e as Error).message }
}

const back = () => router.push({ name: 'admin-backends' })
</script>

<template>
  <div>
    <button class="mb-4 flex items-center gap-1 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="back">
      <ArrowLeft class="h-4 w-4" /> Back to backends
    </button>

    <p v-if="error" class="mb-3 text-sm text-[var(--color-error)]">{{ error }}</p>

    <div v-if="b" class="max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-bold text-[var(--color-text)]">{{ b.name }}</h1>
          <p class="mt-1 font-mono text-xs text-[var(--color-text-muted)]">{{ b.driver }}</p>
        </div>
        <span :class="b.health_status === 'healthy' ? 'text-[var(--color-success)]' : b.health_status === 'unreachable' ? 'text-[var(--color-error)]' : 'text-[var(--color-warning)]'">{{ b.health_status }}</span>
      </div>

      <div class="mt-4 flex flex-col gap-3">
        <pre class="overflow-x-auto rounded-[var(--radius-sm)] bg-[var(--color-surface)] p-3 font-mono text-xs text-[var(--color-text)]">{{ JSON.stringify(b.config, null, 2) }}</pre>
        <p v-if="b.health_error" class="text-xs text-[var(--color-error)]">Provider error: {{ b.health_error }}</p>
        <div class="flex gap-2">
          <AppButton size="sm" @click="checkHealth"><Activity class="h-4 w-4" /> Check health</AppButton>
          <AppButton size="sm" variant="destructive" @click="remove">Delete</AppButton>
        </div>
      </div>
    </div>
  </div>
</template>
