<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

interface Application { id: string; name: string; description: string }

const router = useRouter()
const apps = ref<Application[]>([])
const loading = ref(false)
const showCreate = ref(false)
const name = ref('')
const desc = ref('')
const error = ref('')

async function load() {
  loading.value = true
  try { apps.value = await api.get<Application[]>('/applications') } catch { apps.value = [] } finally { loading.value = false }
}
onMounted(load)

async function create() {
  error.value = ''
  try {
    await api.post('/applications', { name: name.value, description: desc.value })
    showCreate.value = false
    name.value = ''
    desc.value = ''
    load()
  } catch (e) { error.value = (e as Error).message }
}

const open = (app: Application) => router.push({ name: 'application-detail', params: { appId: app.id } })
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[var(--color-text)]">Applications</h1>
        <p class="mt-1 text-sm text-[var(--color-text-muted)]">Non-human principals that authenticate with access keys.</p>
      </div>
      <AppButton size="sm" @click="showCreate = !showCreate">Register application</AppButton>
    </div>

    <div v-if="showCreate" class="mb-4 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
      <div class="flex flex-col gap-3">
        <input v-model="name" placeholder="Name" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <input v-model="desc" placeholder="Description" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <div><AppButton size="sm" @click="create">Create</AppButton></div>
      </div>
    </div>

    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Name</th>
            <th class="px-3 py-2">Description</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in apps" :key="a.id" class="cursor-pointer border-t border-[var(--color-border)] hover:bg-[var(--color-surface)]" @click="open(a)">
            <td class="px-3 py-2.5 font-medium text-[var(--color-text)]">{{ a.name }}</td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ a.description }}</td>
          </tr>
          <tr v-if="!loading && !apps.length">
            <td colspan="2" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No applications yet — register one to issue access keys.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
