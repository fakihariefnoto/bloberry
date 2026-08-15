<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

interface Tenant { id: string; name: string; slug: string; status: string; used_bytes: number; used_objects: number }

const router = useRouter()
const tenants = ref<Tenant[]>([])
const loading = ref(false)
const showCreate = ref(false)
const name = ref('')
const slug = ref('')
const error = ref('')

async function load() {
  loading.value = true
  try { tenants.value = await api.get<Tenant[]>('/admin/tenants') } catch { tenants.value = [] } finally { loading.value = false }
}
onMounted(load)

async function create() {
  error.value = ''
  try {
    await api.post('/tenants', { name: name.value, slug: slug.value })
    showCreate.value = false
    name.value = ''
    slug.value = ''
    load()
  } catch (e) { error.value = (e as Error).message }
}

function formatBytes(n: number) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (n >= 1024 && i < units.length - 1) { n /= 1024; i++ }
  return `${n.toFixed(n >= 10 || i === 0 ? 0 : 1)} ${units[i]}`
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[var(--color-text)]">Tenants</h1>
        <p class="mt-1 text-sm text-[var(--color-text-muted)]">Every tenant on this install.</p>
      </div>
      <AppButton size="sm" @click="showCreate = !showCreate">Create tenant</AppButton>
    </div>

    <div v-if="showCreate" class="mb-4 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
      <div class="flex flex-col gap-3">
        <input v-model="name" placeholder="Name" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <input v-model="slug" placeholder="slug" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <div><AppButton size="sm" @click="create">Create</AppButton></div>
      </div>
    </div>

    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Tenant</th>
            <th class="px-3 py-2">Status</th>
            <th class="px-3 py-2">Storage</th>
            <th class="px-3 py-2">Objects</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in tenants" :key="t.id" class="cursor-pointer border-t border-[var(--color-border)] hover:bg-[var(--color-surface)]" @click="router.push({ name: 'admin-tenant-detail', params: { tenantId: t.id } })">
            <td class="px-3 py-2.5"><span class="font-medium text-[var(--color-text)]">{{ t.name }}</span> <span class="text-xs text-[var(--color-text-muted)]">/{{ t.slug }}</span></td>
            <td class="px-3 py-2.5">
              <span :class="t.status === 'suspended' ? 'text-[var(--color-error)]' : 'text-[var(--color-success)]'">{{ t.status }}</span>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ formatBytes(t.used_bytes) }}</td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ t.used_objects }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
