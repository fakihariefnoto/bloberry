<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { HardDrive, Plus, X, CheckCircle2, Loader2, Eye, EyeOff } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'
import PageHeader from '../components/ui/PageHeader.vue'
import AppInput from '../components/ui/AppInput.vue'

interface Backend { id: string; name: string; driver: string; health_status: string; tenant_id?: string }

const router = useRouter()
const backends = ref<Backend[]>([])
const loading = ref(false)
const showCreate = ref(false)
const submitting = ref(false)
const error = ref('')
const createdId = ref('')
const showSecrets = ref(false)

// ---- form state ----
const name = ref('')
const driver = ref('s3')
const configFields = ref<Record<string, string>>({})
const credentialFields = ref<Record<string, string>>({})

interface FieldDef { key: string; label: string; placeholder?: string; secret?: boolean; required?: boolean; hint?: string }

// Per-driver field presets — no raw JSON in the UI.
const DRIVERS: Record<string, { label: string; config: FieldDef[]; creds: FieldDef[] }> = {
  s3: {
    label: 'AWS S3',
    config: [
      { key: 'endpoint', label: 'Endpoint', placeholder: 'https://s3.us-east-1.amazonaws.com', hint: 'Override for MinIO, B2, Spaces, Wasabi' },
      { key: 'region', label: 'Region', placeholder: 'us-east-1', required: true },
      { key: 'bucket', label: 'Bucket', placeholder: 'my-bucket', required: true },
      { key: 'prefix', label: 'Prefix', placeholder: 'bloberry/' },
    ],
    creds: [
      { key: 'access_key_id', label: 'Access key ID', required: true },
      { key: 'secret_access_key', label: 'Secret access key', secret: true, required: true },
    ],
  },
  r2: {
    label: 'Cloudflare R2',
    config: [
      { key: 'endpoint', label: 'Account endpoint', placeholder: 'https://<account>.r2.cloudflarestorage.com', required: true, hint: 'Account-scoped endpoint — do NOT include the bucket name (no mybucket. prefix). Region is always auto.' },
      { key: 'bucket', label: 'Bucket', placeholder: 'my-bucket', required: true },
      { key: 'prefix', label: 'Prefix', placeholder: 'bloberry/' },
    ],
    creds: [
      { key: 'access_key_id', label: 'Access key ID', required: true },
      { key: 'secret_access_key', label: 'Secret access key', secret: true, required: true },
    ],
  },
  oss: {
    label: 'Alibaba OSS',
    config: [
      { key: 'endpoint', label: 'Endpoint', placeholder: 'https://oss-cn-hangzhou.aliyuncs.com', required: true },
      { key: 'bucket', label: 'Bucket', placeholder: 'my-bucket', required: true },
      { key: 'prefix', label: 'Prefix', placeholder: 'bloberry/' },
    ],
    creds: [
      { key: 'access_key_id', label: 'Access key ID', required: true },
      { key: 'secret_access_key', label: 'Secret access key', secret: true, required: true },
    ],
  },
  gcs: {
    label: 'Google Cloud Storage',
    config: [
      { key: 'bucket', label: 'Bucket', placeholder: 'my-bucket', required: true },
      { key: 'prefix', label: 'Prefix', placeholder: 'bloberry/' },
    ],
    creds: [
      { key: 'service_account_json', label: 'Service account JSON', secret: true, required: true, hint: 'Paste the full service-account key JSON' },
    ],
  },
  azblob: {
    label: 'Azure Blob',
    config: [
      { key: 'endpoint', label: 'Endpoint', placeholder: 'https://<account>.blob.core.windows.net' },
      { key: 'account_name', label: 'Account name', placeholder: 'mystorage', required: true },
      { key: 'container', label: 'Container', placeholder: 'bloberry', required: true },
      { key: 'prefix', label: 'Prefix', placeholder: 'bloberry/' },
    ],
    creds: [
      { key: 'shared_key', label: 'Account key', secret: true, required: true },
    ],
  },
  disk: {
    label: 'Local disk',
    config: [
      { key: 'root', label: 'Storage path', placeholder: '/var/lib/bloberry/objects', required: true, hint: 'Absolute path on the server volume' },
    ],
    creds: [],
  },
}

const currentDriver = computed(() => DRIVERS[driver.value] || DRIVERS.s3)

function resetForm() {
  name.value = ''
  driver.value = 's3'
  configFields.value = {}
  credentialFields.value = {}
  error.value = ''
  createdId.value = ''
  showSecrets.value = false
}

const canSubmit = computed(() => {
  if (!name.value.trim()) return false
  const c = currentDriver.value
  for (const f of [...c.config, ...c.creds]) {
    if (f.required && !((configFields.value[f.key] ?? credentialFields.value[f.key]) || '').trim()) return false
  }
  return true
})

async function create() {
  error.value = ''
  submitting.value = true
  try {
    const res = await api.post<{ id: string }>('/admin/backends', {
      name: name.value.trim(),
      driver: driver.value,
      config: buildObj(configFields.value),
      credentials: buildObj(credentialFields.value),
    })
    createdId.value = res.id
    showCreate.value = false
    resetForm()
    await load()
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    submitting.value = false
  }
}

function buildObj(fields: Record<string, string>): Record<string, string> {
  const out: Record<string, string> = {}
  for (const [k, v] of Object.entries(fields)) {
    if (v !== undefined && v !== '') out[k] = v
  }
  return out
}

async function load() {
  loading.value = true
  try { backends.value = await api.get<Backend[]>('/admin/backends') } catch { backends.value = [] } finally { loading.value = false }
}
onMounted(load)

async function checkHealth(id: string) {
  await api.post(`/admin/backends/${id}/health`)
  await load()
}

const open = (b: Backend) => router.push({ name: 'admin-backend-detail', params: { backendId: b.id } })

const scopeLabel = (b: Backend) => (b.tenant_id ? 'BYO tenant' : 'Install-level')
</script>

<template>
  <div>
    <PageHeader title="Storage engines" description="Credential sets per driver. Assign a project to one via its settings.">
      <AppButton size="sm" @click="showCreate = !showCreate">
        <Plus class="h-4 w-4" /> Register storage engine
      </AppButton>
    </PageHeader>

    <!-- Create modal -->
    <div v-if="showCreate" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="showCreate = false">
      <div class="flex max-h-[90vh] w-full max-w-lg flex-col overflow-hidden rounded-[var(--radius-lg)] bg-[var(--color-surface-raised)] shadow-[var(--shadow-lg)]">
        <div class="flex items-center justify-between border-b border-[var(--color-border)] px-6 py-4">
          <h2 class="text-lg font-semibold text-[var(--color-text)]">Register storage engine</h2>
          <button class="rounded p-1 text-[var(--color-text-muted)] hover:text-[var(--color-text)]" @click="showCreate = false">
            <X class="h-5 w-5" />
          </button>
        </div>

        <div class="flex-1 space-y-5 overflow-y-auto px-6 py-5">
          <AppInput v-model="name" label="Name" placeholder="e.g. s3-eu-prod" hint="How this storage engine appears in the project dropdown" />

          <div class="flex flex-col gap-1">
            <label class="text-xs font-medium text-[var(--color-text-muted)]">Provider</label>
            <div class="grid grid-cols-3 gap-2">
              <button
                v-for="(d, key) in DRIVERS"
                :key="key"
                type="button"
                class="flex h-11 items-center justify-center rounded-[var(--radius-md)] border text-xs font-medium transition-colors"
                :class="driver === key
                  ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)] text-[var(--color-primary)]'
                  : 'border-[var(--color-border)] bg-[var(--color-surface)] text-[var(--color-text-muted)] hover:border-[var(--color-primary)]'"
                @click="driver = key"
              >
                {{ d.label }}
              </button>
            </div>
          </div>

          <div v-if="currentDriver.config.length" class="space-y-4">
            <p class="text-xs font-semibold uppercase tracking-wide text-[var(--color-text-muted)]">Configuration</p>
            <AppInput
              v-for="f in currentDriver.config"
              :key="f.key"
              v-model="configFields[f.key]"
              :label="f.label"
              :placeholder="f.placeholder"
              :hint="f.hint"
            />
          </div>

          <div v-if="currentDriver.creds.length" class="space-y-4">
            <div class="flex items-center justify-between">
              <p class="text-xs font-semibold uppercase tracking-wide text-[var(--color-text-muted)]">Credentials</p>
              <button class="text-xs text-[var(--color-primary)] hover:underline" @click="showSecrets = !showSecrets">
                {{ showSecrets ? 'Hide secrets' : 'Reveal secrets' }}
              </button>
            </div>
            <div v-for="f in currentDriver.creds" :key="f.key" class="flex flex-col gap-1">
              <label class="text-xs font-medium text-[var(--color-text-muted)]">{{ f.label }}</label>
              <div class="relative">
                <input
                  v-model="credentialFields[f.key]"
                  :type="f.secret && !showSecrets ? 'password' : 'text'"
                  :placeholder="f.placeholder"
                  class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-4 font-mono text-xs outline-none focus:border-[var(--color-primary)] focus:border-2"
                />
                <button
                  v-if="f.secret"
                  type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--color-text-muted)] hover:text-[var(--color-text)]"
                  @click="showSecrets = !showSecrets"
                >
                  <EyeOff v-if="showSecrets" class="h-4 w-4" />
                  <Eye v-else class="h-4 w-4" />
                </button>
              </div>
              <p v-if="f.hint" class="text-xs text-[var(--color-text-muted)]">{{ f.hint }}</p>
            </div>
          </div>

          <p class="text-xs text-[var(--color-text-muted)]">
            Credentials are envelope-encrypted before they touch the database and never shown again.
          </p>

          <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
        </div>

        <div class="flex items-center justify-end gap-2 border-t border-[var(--color-border)] px-6 py-4">
          <AppButton variant="ghost" @click="showCreate = false">Cancel</AppButton>
          <AppButton :loading="submitting" :disabled="!canSubmit" @click="create">
            <CheckCircle2 class="mr-2 h-4 w-4" /> Register storage engine
          </AppButton>
        </div>
      </div>
    </div>

    <!-- Success toast -->
    <div v-if="createdId" class="fixed bottom-6 right-6 z-50 flex items-center gap-2 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] px-4 py-3 shadow-[var(--shadow-md)]">
      <CheckCircle2 class="h-5 w-5 text-[var(--color-success)]" />
      <span class="text-sm text-[var(--color-text)]">Backend registered</span>
      <button class="ml-2 text-xs font-medium text-[var(--color-primary)] hover:underline" @click="createdId = ''">Dismiss</button>
    </div>

    <!-- List -->
    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Name</th>
            <th class="px-3 py-2">Driver</th>
            <th class="px-3 py-2">Health</th>
            <th class="px-3 py-2">Scope</th>
            <th class="px-3 py-2 text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in backends" :key="b.id" class="cursor-pointer border-t border-[var(--color-border)] hover:bg-[var(--color-surface)]" @click="open(b)">
            <td class="px-3 py-2.5"><span class="flex items-center gap-2 font-medium text-[var(--color-text)]"><HardDrive class="h-4 w-4 text-[var(--color-text-muted)]" /> {{ b.name }}</span></td>
            <td class="px-3 py-2.5 font-mono text-xs text-[var(--color-text-muted)]">{{ b.driver }}</td>
            <td class="px-3 py-2.5">
              <span
                class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium"
                :class="b.health_status === 'healthy'
                  ? 'bg-[var(--color-success)]/10 text-[var(--color-success)]'
                  : b.health_status === 'unreachable'
                    ? 'bg-[var(--color-error)]/10 text-[var(--color-error)]'
                    : 'bg-[var(--color-warning)]/10 text-[var(--color-warning)]'"
              >
                <Loader2 v-if="b.health_status === 'checking'" class="h-3 w-3 animate-spin" />
                {{ b.health_status }}
              </span>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ scopeLabel(b) }}</td>
            <td class="px-3 py-2.5 text-right">
              <button class="text-xs font-medium text-[var(--color-primary)] hover:underline" @click.stop="checkHealth(b.id)">
                Check health
              </button>
            </td>
          </tr>
          <tr v-if="!loading && !backends.length">
            <td colspan="5" class="px-4 py-12 text-center">
              <HardDrive class="mx-auto h-8 w-8 text-[var(--color-text-muted)]" />
              <p class="mt-3 text-sm font-medium text-[var(--color-text)]">No storage engines yet</p>
              <p class="mt-1 text-xs text-[var(--color-text-muted)]">Register one to connect a provider.</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
