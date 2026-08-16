<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { KeyRound, Plus, Copy, Check, Building2, Layers } from 'lucide-vue-next'
import { api } from '../lib/api'
import { useAuthStore } from '../stores/auth'
import AppButton from '../components/ui/AppButton.vue'
import AppModal from '../components/ui/AppModal.vue'
import AppInput from '../components/ui/AppInput.vue'
import PageHeader from '../components/ui/PageHeader.vue'

interface KeyRec {
  id: string
  name?: string
  tenant_id: string
  tenant_name?: string
  prefix: string
  last_four: string
  permissions: string[]
  scope_folder_ids: string[]
  application_name?: string
  application_id?: string
  revoked_at?: string
  expires_at?: string
}
interface AppRec { id: string; name: string }
interface TenantRec { id: string; name: string; slug: string }

const auth = useAuthStore()
const isAdmin = auth.user?.platform_role === 'platform_admin'

const keys = ref<KeyRec[]>([])
const apps = ref<AppRec[]>([])
const tenants = ref<TenantRec[]>([])
const loading = ref(false)

const showCreate = ref(false)
const keyType = ref<'tenant' | 'app'>('tenant')
const keyName = ref('')
const tenantId = ref('')
const appId = ref('')
const perms = ref(['read'])
const permOptions = ['read', 'write', 'delete', 'share', 'admin']
const scopeFolders = ref('')
const creating = ref(false)
const error = ref('')
const createdSecret = ref('')
const copied = ref(false)

async function load() {
  loading.value = true
  try { keys.value = await api.get<KeyRec[]>('/keys') } catch { keys.value = [] } finally { loading.value = false }
}
async function loadMeta() {
  try { apps.value = await api.get<AppRec[]>('/applications') } catch { apps.value = [] }
  if (isAdmin) {
    try { tenants.value = await api.get<TenantRec[]>('/admin/tenants') } catch { tenants.value = [] }
  }
}
onMounted(() => { load(); loadMeta() })

function openCreate() {
  error.value = ''
  createdSecret.value = ''
  keyType.value = 'tenant'
  keyName.value = ''
  tenantId.value = isAdmin && tenants.value.length ? tenants.value[0].id : ''
  appId.value = apps.value[0]?.id || ''
  perms.value = ['read']
  scopeFolders.value = ''
  showCreate.value = true
}

async function createKey() {
  error.value = ''
  creating.value = true
  const scope = scopeFolders.value.split(',').map((s) => s.trim()).filter(Boolean)
  try {
    const res = keyType.value === 'tenant'
      ? await api.post<{ secret: string }>('/keys', {
          name: keyName.value,
          tenant_id: isAdmin ? tenantId.value : undefined,
          permissions: perms.value,
          scope_folder_ids: scope,
        })
      : await api.post<{ secret: string }>(`/applications/${appId.value}/keys`, {
          name: keyName.value,
          permissions: perms.value,
          scope_folder_ids: scope,
        })
    createdSecret.value = res.secret
    load()
  } catch (e) { error.value = (e as Error).message } finally { creating.value = false }
}

async function revoke(id: string) {
  if (!window.confirm('Revoke this key? Existing integrations using it will stop working.')) return
  try {
    await api.delete(`/keys/${id}`)
    load()
  } catch (e) { alert((e as Error).message) }
}

function copySecret() {
  navigator.clipboard?.writeText(createdSecret.value)
  copied.value = true
  setTimeout(() => (copied.value = false), 1500)
}
</script>

<template>
  <div>
    <PageHeader title="API keys" description="Credentials for the SDKs and integrations, scoped to a tenant.">
      <AppButton size="sm" @click="openCreate"><Plus class="h-4 w-4" /> Create key</AppButton>
    </PageHeader>

    <!-- Create modal -->
    <AppModal :open="showCreate" title="Create API key" description="The secret is shown once — copy it now." @close="showCreate = false">
      <div v-if="!createdSecret" class="flex flex-col gap-4">
        <AppInput v-model="keyName" label="Name" placeholder="e.g. CI uploads, backup worker" hint="A name to identify this key in the list" />

        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Key type</label>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              class="h-11 rounded-[var(--radius-md)] border px-2 text-xs font-medium"
              :class="keyType === 'tenant'
                ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)] text-[var(--color-primary)]'
                : 'border-[var(--color-border)] text-[var(--color-text-muted)]'"
              @click="keyType = 'tenant'"
            >
              Tenant key
            </button>
            <button
              type="button"
              class="h-11 rounded-[var(--radius-md)] border px-2 text-xs font-medium"
              :class="keyType === 'app'
                ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)] text-[var(--color-primary)]'
                : 'border-[var(--color-border)] text-[var(--color-text-muted)]'"
              @click="keyType = 'app'"
            >
              Application key
            </button>
          </div>
          <p class="text-xs text-[var(--color-text-muted)]">
            <span v-if="keyType === 'tenant'">
              <Building2 class="mr-1 inline h-3 w-3" /> Tenant key: acts directly on the tenant, scoped to folders and permissions. Best for SDKs and CI.
            </span>
            <span v-else>
              <Layers class="mr-1 inline h-3 w-3" /> Application key: belongs to an application (a non-human principal) you manage in Applications.
            </span>
          </p>
        </div>

        <div v-if="isAdmin" class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Tenant</label>
          <select v-model="tenantId" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
            <option v-for="t in tenants" :key="t.id" :value="t.id">{{ t.name }} ({{ t.slug }})</option>
          </select>
        </div>

        <div v-if="keyType === 'app'" class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Application</label>
          <select v-model="appId" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
            <option v-for="a in apps" :key="a.id" :value="a.id">{{ a.name }}</option>
          </select>
          <p v-if="!apps.length" class="text-xs text-[var(--color-warning)]">
            No applications yet — register one in Applications, or use a Tenant key instead.
          </p>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Folder scope</label>
          <input v-model="scopeFolders" placeholder="Folder IDs, comma-separated (blank = whole tenant)" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-xs outline-none focus:border-[var(--color-primary)] focus:border-2" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Permissions</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="p in permOptions"
              :key="p"
              type="button"
              class="h-9 rounded-[var(--radius-md)] border px-3 text-xs font-medium"
              :class="perms.includes(p)
                ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)] text-[var(--color-primary)]'
                : 'border-[var(--color-border)] text-[var(--color-text-muted)]'"
              @click="perms.includes(p) ? (perms = perms.filter((x) => x !== p)) : perms.push(p)"
            >
              {{ p }}
            </button>
          </div>
        </div>

        <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
        <div class="flex justify-end gap-2">
          <AppButton variant="ghost" @click="showCreate = false">Cancel</AppButton>
          <AppButton :loading="creating" @click="createKey">Create key</AppButton>
        </div>
      </div>

      <div v-else class="flex flex-col gap-4">
        <p class="rounded-[var(--radius-sm)] border border-[var(--color-warning)] bg-[var(--color-warning)]/10 px-3 py-2 text-xs text-[var(--color-warning)]">
          You won't see this secret again. Use it as the bearer token in your SDK client.
        </p>
        <div class="flex items-center gap-2 rounded-[var(--radius-sm)] bg-[var(--color-surface)] p-3">
          <code class="flex-1 break-all font-mono text-xs text-[var(--color-text)]">{{ createdSecret }}</code>
          <button class="text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="copySecret">
            <Check v-if="copied" class="h-4 w-4 text-[var(--color-success)]" />
            <Copy v-else class="h-4 w-4" />
          </button>
        </div>
        <div class="flex justify-end">
          <AppButton @click="showCreate = false">Done</AppButton>
        </div>
      </div>
    </AppModal>

    <!-- List -->
    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Key</th>
            <th class="px-3 py-2">Name</th>
            <th class="px-3 py-2">Tenant</th>
            <th class="px-3 py-2">Type</th>
            <th class="px-3 py-2">Permissions</th>
            <th class="px-3 py-2">State</th>
            <th class="px-3 py-2 text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="k in keys" :key="k.id" class="border-t border-[var(--color-border)]">
            <td class="px-3 py-2.5 font-mono text-xs text-[var(--color-text)]">{{ k.prefix }}••••{{ k.last_four }}</td>
            <td class="px-3 py-2.5 text-[var(--color-text)]">{{ k.name || '—' }}</td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ k.tenant_name || k.tenant_id.slice(0, 8) }}</td>
            <td class="px-3 py-2.5">
              <span v-if="k.application_name" class="inline-flex items-center gap-1.5 text-[var(--color-text-muted)]">
                <Layers class="h-3.5 w-3.5" /> {{ k.application_name }}
              </span>
              <span v-else class="inline-flex items-center gap-1.5 text-[var(--color-text-muted)]">
                <Building2 class="h-3.5 w-3.5" /> Tenant
              </span>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ (k.permissions || []).join(', ') }}</td>
            <td class="px-3 py-2.5">
              <span v-if="k.revoked_at" class="text-xs text-[var(--color-error)]">Revoked</span>
              <span v-else class="text-xs text-[var(--color-success)]">Active</span>
            </td>
            <td class="px-3 py-2.5 text-right">
              <button v-if="!k.revoked_at" class="text-xs font-medium text-[var(--color-error)] hover:underline" @click="revoke(k.id)">Revoke</button>
            </td>
          </tr>
          <tr v-if="!loading && !keys.length">
            <td colspan="7" class="px-4 py-12 text-center">
              <KeyRound class="mx-auto h-8 w-8 text-[var(--color-text-muted)]" />
              <p class="mt-3 text-sm font-medium text-[var(--color-text)]">No API keys yet</p>
              <p class="mt-1 text-xs text-[var(--color-text-muted)]">Create one to authenticate your SDK clients.</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
