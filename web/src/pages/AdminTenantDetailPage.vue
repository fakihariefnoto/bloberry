<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, HardDrive, PauseCircle, CheckCircle2, UserPlus, Mail, Trash2 } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppModal from '../components/ui/AppModal.vue'

interface Tenant {
  id: string; name: string; slug: string; status: string
  quota_bytes: number; quota_objects: number
  used_bytes: number; used_objects: number
  default_storage_id?: string
}
interface Backend { id: string; name: string; driver: string; health_status: string }

const route = useRoute()
const router = useRouter()
const tenantId = route.params.tenantId as string

const t = ref<Tenant | null>(null)
const backends = ref<Backend[]>([])
const storageId = ref('')
const quotaBytes = ref('0')
const quotaObjects = ref('0')
const saved = ref('')
const error = ref('')
const saving = ref(false)

async function load() {
  t.value = await api.get<Tenant>(`/tenants/${tenantId}`)
  if (t.value) {
    quotaBytes.value = String(t.value.quota_bytes)
    quotaObjects.value = String(t.value.quota_objects)
    storageId.value = t.value.default_storage_id || ''
  }
}
async function loadBackends() {
  try { backends.value = await api.get<Backend[]>('/admin/backends') } catch { backends.value = [] }
}

async function saveAll() {
  error.value = ''
  saved.value = ''
  saving.value = true
  try {
    await api.patch(`/tenants/${tenantId}`, {
      quota_bytes: Number(quotaBytes.value) || 0,
      quota_objects: Number(quotaObjects.value) || 0,
      default_storage_id: storageId.value || undefined,
    })
    saved.value = 'Saved.'
    await load()
  } catch (e) { error.value = (e as Error).message } finally { saving.value = false }
}

async function toggleStatus() {
  await api.patch(`/tenants/${tenantId}`, { status: t.value?.status === 'suspended' ? 'active' : 'suspended' })
  await load()
}

function formatBytes(n: number) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (n >= 1024 && i < units.length - 1) { n /= 1024; i++ }
  return `${n.toFixed(n >= 10 || i === 0 ? 0 : 1)} ${units[i]}`
}

const usedPct = () => (t.value?.quota_bytes ? Math.round(((t.value.used_bytes || 0) / t.value.quota_bytes) * 100) : 0)
const back = () => router.push({ name: 'admin-tenants' })

// --- Member management (platform admin / tenant admin+owner) ---
interface Membership { id: string; user_id: string; email?: string; display_name?: string; role: string; created_at: string }
interface Invitation { id: string; email: string; role: string }

const members = ref<Membership[]>([])
const invites = ref<Invitation[]>([])
const showAdd = ref(false)
const addEmail = ref('')
const addRole = ref('member')
const memberError = ref('')
const showRemove = ref(false)
const removeTarget = ref<Membership | null>(null)
const addingMember = ref(false)

async function loadMembers() {
  try {
    members.value = await api.get<Membership[]>(`/tenants/${tenantId}/members`)
    invites.value = await api.get<Invitation[]>(`/tenants/${tenantId}/invitations`)
  } catch { members.value = []; invites.value = [] }
}

async function addMemberByEmail() {
  memberError.value = ''
  addingMember.value = true
  try {
    await api.post(`/tenants/${tenantId}/members`, { email: addEmail.value, role: addRole.value })
    showAdd.value = false
    addEmail.value = ''
    loadMembers()
  } catch (e) { memberError.value = (e as Error).message } finally { addingMember.value = false }
}

async function inviteMember() {
  memberError.value = ''
  addingMember.value = true
  try {
    await api.post(`/tenants/${tenantId}/invitations`, { email: addEmail.value, role: addRole.value })
    showAdd.value = false
    addEmail.value = ''
    loadMembers()
  } catch (e) { memberError.value = (e as Error).message } finally { addingMember.value = false }
}

function askRemove(m: Membership) {
  removeTarget.value = m
  showRemove.value = true
}

async function doRemoveMember() {
  const m = removeTarget.value
  if (!m) return
  try {
    await api.delete(`/tenants/${tenantId}/members/${m.id}`)
    showRemove.value = false
    removeTarget.value = null
    loadMembers()
  } catch (e) { memberError.value = (e as Error).message; showRemove.value = false }
}

async function setMemberRole(id: string, role: string) {
  await api.patch(`/tenants/${tenantId}/members/${id}`, { role })
  loadMembers()
}

onMounted(() => { load(); loadBackends(); loadMembers() })
</script>

<template>
  <div>
    <button class="mb-4 flex items-center gap-1 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="back">
      <ArrowLeft class="h-4 w-4" /> Back to projects
    </button>

    <div v-if="t">
      <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-[var(--color-text)]">{{ t.name }}</h1>
          <p class="mt-1 text-sm text-[var(--color-text-muted)]">/{{ t.slug }}</p>
        </div>
        <AppButton size="sm" :variant="t.status === 'suspended' ? 'secondary' : 'destructive'" @click="toggleStatus">
          <PauseCircle v-if="t.status !== 'suspended'" class="h-4 w-4" />
          <CheckCircle2 v-else class="h-4 w-4" />
          {{ t.status === 'suspended' ? 'Unsuspend' : 'Suspend' }}
        </AppButton>
      </div>

      <!-- Stat cards -->
      <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
          <p class="text-xs font-medium text-[var(--color-text-muted)]">Storage used</p>
          <p class="mt-2 text-2xl font-bold tracking-tight text-[var(--color-text)]">{{ formatBytes(t.used_bytes) }}</p>
          <p v-if="t.quota_bytes" class="mt-2 text-xs text-[var(--color-text-muted)]">
            {{ usedPct() }}% of {{ formatBytes(t.quota_bytes) }}
          </p>
          <div v-if="t.quota_bytes" class="mt-3 h-1.5 overflow-hidden rounded-full bg-[var(--color-border)]">
            <div
              class="h-full"
              :class="usedPct() > 80 ? 'bg-[var(--color-warning)]' : 'bg-[var(--color-primary)]'"
              :style="{ width: `${Math.min(100, usedPct())}%` }"
            />
          </div>
        </div>
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
          <p class="text-xs font-medium text-[var(--color-text-muted)]">Objects</p>
          <p class="mt-2 text-2xl font-bold tracking-tight text-[var(--color-text)]">{{ t.used_objects.toLocaleString() }}</p>
          <p v-if="t.quota_objects" class="mt-2 text-xs text-[var(--color-text-muted)]">quota {{ t.quota_objects.toLocaleString() }}</p>
        </div>
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
          <p class="text-xs font-medium text-[var(--color-text-muted)]">Status</p>
          <p class="mt-2 inline-flex items-center gap-2 text-2xl font-bold tracking-tight" :class="t.status === 'suspended' ? 'text-[var(--color-error)]' : 'text-[var(--color-success)]'">
            <span class="h-2.5 w-2.5 rounded-full" :class="t.status === 'suspended' ? 'bg-[var(--color-error)]' : 'bg-[var(--color-success)]'" />
            {{ t.status }}
          </p>
        </div>
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
          <p class="text-xs font-medium text-[var(--color-text-muted)]">Storage engine</p>
          <p class="mt-2 flex items-center gap-2 text-2xl font-bold tracking-tight text-[var(--color-text)]">
            <HardDrive class="h-5 w-5 text-[var(--color-primary)]" />
            {{ backends.find((b) => b.id === t?.default_storage_id)?.driver || 'none' }}
          </p>
          <p class="mt-2 truncate text-xs text-[var(--color-text-muted)]">{{ backends.find((b) => b.id === t?.default_storage_id)?.name || 'Not assigned' }}</p>
        </div>
      </div>

      <!-- Settings -->
      <div class="mt-6 max-w-lg rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
        <h2 class="text-lg font-semibold text-[var(--color-text)]">Configuration</h2>
        <div class="mt-5 flex flex-col gap-4">
          <div class="flex flex-col gap-1">
            <label class="text-xs font-medium text-[var(--color-text-muted)]">Storage engine</label>
            <select v-model="storageId" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
              <option value="">No engine assigned</option>
              <option v-for="b in backends" :key="b.id" :value="b.id">{{ b.name }} ({{ b.driver }})</option>
            </select>
            <p class="flex items-center gap-1.5 text-xs text-[var(--color-text-muted)]">
              <HardDrive class="h-3.5 w-3.5" />
              Applies to new objects — existing files keep their engine.
            </p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <AppInput v-model="quotaBytes" label="Quota (bytes, 0 = unlimited)" type="number" placeholder="0" />
            <AppInput v-model="quotaObjects" label="Object quota (0 = unlimited)" type="number" placeholder="0" />
          </div>

          <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
          <p v-if="saved" class="text-xs text-[var(--color-success)]">{{ saved }}</p>
          <div class="flex justify-end"><AppButton :loading="saving" @click="saveAll">Save changes</AppButton></div>
        </div>
      </div>

      <!-- Members -->
      <div class="mt-6 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-lg font-semibold text-[var(--color-text)]">Members</h2>
            <p class="mt-1 text-xs text-[var(--color-text-muted)]">Who has access to {{ t.name }}.</p>
          </div>
          <AppButton size="sm" @click="showAdd = true"><UserPlus class="mr-1.5 h-4 w-4" /> Add member</AppButton>
        </div>

        <div class="mt-4 overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
          <table class="w-full text-sm">
            <thead>
              <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
                <th class="px-3 py-2">User</th>
                <th class="px-3 py-2">Role</th>
                <th class="px-3 py-2"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in members" :key="m.id" class="border-t border-[var(--color-border)]">
                <td class="px-3 py-2.5">
                  <div class="flex items-center gap-2.5">
                    <span class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-[var(--color-primary-subtle)] text-xs font-semibold text-[var(--color-primary)]">
                      {{ (m.display_name || m.email || m.user_id).charAt(0).toUpperCase() }}
                    </span>
                    <div class="min-w-0">
                      <p class="truncate text-sm font-medium text-[var(--color-text)]">{{ m.display_name || '—' }}</p>
                      <p class="truncate text-xs text-[var(--color-text-muted)]">{{ m.email || m.user_id }}</p>
                    </div>
                  </div>
                </td>
                <td class="px-3 py-2.5">
                  <select
                    :value="m.role"
                    class="rounded-[var(--radius-sm)] border border-[var(--color-border)] bg-[var(--color-surface)] px-2 py-1 text-xs text-[var(--color-text-muted)] outline-none focus:border-[var(--color-primary)]"
                    @change="setMemberRole(m.id, ($event.target as HTMLSelectElement).value)"
                  >
                    <option value="viewer">Viewer</option>
                    <option value="member">Member</option>
                    <option value="tenant_admin">Admin</option>
                    <option value="tenant_owner">Owner</option>
                  </select>
                </td>
                <td class="px-3 py-2.5 text-right">
                  <button class="text-xs text-[var(--color-error)] hover:underline" @click="askRemove(m)">Remove</button>
                </td>
              </tr>
              <tr v-if="!members.length">
                <td colspan="3" class="px-4 py-8 text-center text-sm text-[var(--color-text-muted)]">No members yet.</td>
              </tr>
            </tbody>
          </table>
        </div>

        <h3 class="mb-2 mt-5 text-sm font-semibold text-[var(--color-text)]">Pending invitations</h3>
        <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
          <table class="w-full text-sm">
            <tbody>
              <tr v-for="i in invites" :key="i.id" class="border-t border-[var(--color-border)]">
                <td class="px-3 py-2.5"><span class="flex items-center gap-2"><Mail class="h-3.5 w-3.5 text-[var(--color-text-muted)]" /> {{ i.email }}</span></td>
                <td class="px-3 py-2.5 text-xs text-[var(--color-text-muted)]">{{ i.role }}</td>
              </tr>
              <tr v-if="!invites.length"><td class="px-4 py-6 text-center text-sm text-[var(--color-text-muted)]">No pending invitations</td></tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Add member modal -->
      <AppModal :open="showAdd" title="Add a member" description="Add someone who already has an account by email, or invite them." @close="showAdd = false">
        <div class="flex flex-col gap-3">
          <AppInput v-model="addEmail" label="Email" placeholder="jane@acme.com" autofocus />
          <div class="flex flex-col gap-1">
            <label class="text-xs font-medium text-[var(--color-text-muted)]">Role</label>
            <select v-model="addRole" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
              <option value="viewer">Viewer — read-only access</option>
              <option value="member">Member — upload and manage files</option>
              <option value="tenant_admin">Admin — manage members and settings</option>
              <option value="tenant_owner">Owner — full control</option>
            </select>
          </div>
          <p v-if="memberError" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ memberError }}</p>
        </div>
        <template #footer>
          <AppButton variant="ghost" @click="showAdd = false">Cancel</AppButton>
          <AppButton variant="secondary" :disabled="addingMember" @click="inviteMember"><Mail class="mr-1.5 h-4 w-4" /> Invite</AppButton>
          <AppButton :loading="addingMember" @click="addMemberByEmail"><UserPlus class="mr-1.5 h-4 w-4" /> Add</AppButton>
        </template>
      </AppModal>

      <!-- Remove member confirm -->
      <AppModal :open="showRemove" title="Remove member" :description="`Remove ${removeTarget?.display_name || removeTarget?.email || 'this member'}? They'll lose access immediately.`" @close="showRemove = false">
        <template #footer>
          <AppButton variant="ghost" @click="showRemove = false">Cancel</AppButton>
          <AppButton variant="destructive" @click="doRemoveMember"><Trash2 class="mr-1.5 h-4 w-4" /> Remove</AppButton>
        </template>
      </AppModal>
    </div>
  </div>
</template>
