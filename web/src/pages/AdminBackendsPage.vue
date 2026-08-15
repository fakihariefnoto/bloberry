<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { HardDrive, Plus } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

interface Backend { id: string; name: string; driver: string; health_status: string; tenant_id?: string }

const router = useRouter()
const backends = ref<Backend[]>([])
const loading = ref(false)
const showCreate = ref(false)
const name = ref('')
const driver = ref('s3')
const config = ref('{}')
const credentials = ref('{}')
const error = ref('')

async function load() {
  loading.value = true
  try { backends.value = await api.get<Backend[]>('/admin/backends') } catch { backends.value = [] } finally { loading.value = false }
}
onMounted(load)

async function create() {
  error.value = ''
  try {
    await api.post('/admin/backends', {
      name: name.value, driver: driver.value,
      config: JSON.parse(config.value || '{}'),
      credentials: JSON.parse(credentials.value || '{}'),
    })
    showCreate.value = false
    load()
  } catch (e) { error.value = (e as Error).message }
}

const open = (b: Backend) => router.push({ name: 'admin-backend-detail', params: { backendId: b.id } })
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[var(--color-text)]">Storage backends</h1>
        <p class="mt-1 text-sm text-[var(--color-text-muted)]">Credential sets per driver. Assign a tenant to one via its settings.</p>
      </div>
      <AppButton size="sm" @click="showCreate = !showCreate"><Plus class="h-4 w-4" /> Register backend</AppButton>
    </div>

    <div v-if="showCreate" class="mb-4 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
      <div class="flex flex-col gap-3">
        <input v-model="name" placeholder="Name (e.g. s3-eu-prod)" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <select v-model="driver" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm">
          <option value="s3">S3</option>
          <option value="r2">Cloudflare R2</option>
          <option value="oss">Alibaba OSS</option>
          <option value="gcs">Google Cloud Storage</option>
          <option value="azblob">Azure Blob</option>
          <option value="disk">Local disk</option>
        </select>
        <input v-model="config" placeholder='Config JSON ({"bucket":"…"})' class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 font-mono text-xs" />
        <input v-model="credentials" placeholder='Credentials JSON ({"access_key_id":"…"})' class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 font-mono text-xs" />
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <div><AppButton size="sm" @click="create">Register</AppButton></div>
      </div>
    </div>

    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Name</th>
            <th class="px-3 py-2">Driver</th>
            <th class="px-3 py-2">Health</th>
            <th class="px-3 py-2">Scope</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in backends" :key="b.id" class="cursor-pointer border-t border-[var(--color-border)] hover:bg-[var(--color-surface)]" @click="open(b)">
            <td class="px-3 py-2.5"><span class="flex items-center gap-2 font-medium text-[var(--color-text)]"><HardDrive class="h-4 w-4 text-[var(--color-text-muted)]" /> {{ b.name }}</span></td>
            <td class="px-3 py-2.5 font-mono text-xs">{{ b.driver }}</td>
            <td class="px-3 py-2.5">
              <span :class="b.health_status === 'healthy' ? 'text-[var(--color-success)]' : b.health_status === 'unreachable' ? 'text-[var(--color-error)]' : 'text-[var(--color-warning)]'">{{ b.health_status }}</span>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ b.tenant_id ? 'BYO tenant' : 'Install-level' }}</td>
          </tr>
          <tr v-if="!loading && !backends.length">
            <td colspan="4" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No storage backends registered.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
