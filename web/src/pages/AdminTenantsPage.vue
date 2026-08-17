<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { HardDrive } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'
import PageHeader from '../components/ui/PageHeader.vue'
import AppInput from '../components/ui/AppInput.vue'

interface Tenant { id: string; name: string; slug: string; status: string; used_bytes: number; used_objects: number; default_storage_id?: string }
interface Backend { id: string; name: string; driver: string; health_status: string }

const router = useRouter()
const tenants = ref<Tenant[]>([])
const backends = ref<Backend[]>([])
const loading = ref(false)
const showCreate = ref(false)
const name = ref('')
const slug = ref('')
const backendId = ref('')
const error = ref('')

async function load() {
  loading.value = true
  try { tenants.value = await api.get<Tenant[]>('/admin/tenants') } catch { tenants.value = [] } finally { loading.value = false }
}
async function loadBackends() {
  try { backends.value = await api.get<Backend[]>('/admin/backends') } catch { backends.value = [] }
}
onMounted(() => { load(); loadBackends() })

async function create() {
  error.value = ''
  try {
    await api.post('/tenants', {
      name: name.value,
      slug: slug.value,
      default_storage_id: backendId.value || undefined,
    })
    showCreate.value = false
    name.value = ''
    slug.value = ''
    backendId.value = ''
    load()
  } catch (e) { error.value = (e as Error).message }
}

function backendLabel(b: Backend) {
  return `${b.name} (${b.driver})`
}
</script>

<template>
  <div>
    <PageHeader title="Projects" description="Every project on this install.">
      <AppButton size="sm" @click="showCreate = !showCreate">Create project</AppButton>
    </PageHeader>

    <!-- Create modal -->
    <div v-if="showCreate" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="showCreate = false">
      <div class="w-full max-w-md rounded-[var(--radius-lg)] bg-[var(--color-surface-raised)] p-6 shadow-[var(--shadow-lg)]">
        <h2 class="text-lg font-semibold text-[var(--color-text)]">Create project</h2>
        <div class="mt-5 flex flex-col gap-4">
          <AppInput v-model="name" label="Name" placeholder="Acme Corp" />
          <AppInput v-model="slug" label="Slug" placeholder="acme" hint="Used in URLs and short links" />

          <div class="flex flex-col gap-1">
            <label class="text-xs font-medium text-[var(--color-text-muted)]">Storage engine</label>
            <select v-model="backendId" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
              <option value="">Assign later</option>
              <option v-for="b in backends" :key="b.id" :value="b.id">{{ backendLabel(b) }}</option>
            </select>
            <p v-if="!backends.length" class="text-xs text-[var(--color-warning)]">No storage engines registered yet — register one in Storage engines.</p>
          </div>

          <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>

          <div class="flex justify-end gap-2">
            <AppButton variant="ghost" @click="showCreate = false">Cancel</AppButton>
            <AppButton @click="create">Create project</AppButton>
          </div>
        </div>
      </div>
    </div>

    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Project</th>
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
            <td class="px-3 py-2.5">
              <span v-if="t.default_storage_id" class="inline-flex items-center gap-1.5 text-xs text-[var(--color-text-muted)]">
                <HardDrive class="h-3.5 w-3.5" /> {{ backends.find((b) => b.id === t.default_storage_id)?.name || t.default_storage_id.slice(0, 8) }}
              </span>
              <span v-else class="text-xs text-[var(--color-text-muted)]">—</span>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ t.used_objects }}</td>
          </tr>
          <tr v-if="!loading && !tenants.length">
            <td colspan="4" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No projects yet.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
