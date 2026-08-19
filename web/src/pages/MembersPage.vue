<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Mail, UserPlus, Trash2 } from 'lucide-vue-next'
import { api } from '../lib/api'
import { useTenantStore } from '../stores/tenant'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppModal from '../components/ui/AppModal.vue'

interface Membership { id: string; user_id: string; email?: string; display_name?: string; role: string; created_at: string }
interface Invitation { id: string; email: string; role: string }

const tenants = useTenantStore()
const members = ref<Membership[]>([])
const invites = ref<Invitation[]>([])
const loading = ref(false)
const error = ref('')

const showAdd = ref(false)
const email = ref('')
const role = ref('member')
const method = ref<'invite' | 'password' | 'activation'>('invite')
const password = ref('')
const adding = ref(false)
const generated = ref('')
const smtpConfigured = ref(true)

const showRemove = ref(false)
const removeTarget = ref<Membership | null>(null)
const removing = ref(false)

const roleLabels: Record<string, string> = {
  viewer: 'Viewer', member: 'Member', tenant_admin: 'Admin', tenant_owner: 'Owner',
}

async function load() {
  loading.value = true
  try {
    members.value = await api.get<Membership[]>(`/tenants/${tenants.currentId}/members`)
    invites.value = await api.get<Invitation[]>(`/tenants/${tenants.currentId}/invitations`)
  } catch { members.value = []; invites.value = [] } finally { loading.value = false }
}
onMounted(async () => {
  await load()
  try {
    const st = await api.get<{ smtp_configured?: boolean }>('/setup/status')
    smtpConfigured.value = st?.smtp_configured === true
  } catch { smtpConfigured.value = false }
})

function openAdd() {
  email.value = ''
  password.value = ''
  generated.value = ''
  error.value = ''
  method.value = smtpConfigured.value ? 'invite' : 'password'
  showAdd.value = true
}

async function addMember() {
  error.value = ''
  generated.value = ''
  adding.value = true
  try {
    const body: Record<string, unknown> = { email: email.value, role: role.value, method: method.value }
    if (method.value === 'password') {
      if (password.value.length < 8) {
        error.value = 'Password must be at least 8 characters.'
        adding.value = false
        return
      }
      body.password = password.value
    }
    await api.post(`/tenants/${tenants.currentId}/members`, body)
    if (method.value === 'password') {
      generated.value = password.value
    }
    showAdd.value = false
    email.value = ''
    password.value = ''
    await load()
  } catch (e) { error.value = (e as Error).message } finally { adding.value = false }
}

function askRemove(m: Membership) {
  removeTarget.value = m
  showRemove.value = true
}

async function doRemove() {
  const m = removeTarget.value
  if (!m) return
  removing.value = true
  try {
    await api.delete(`/tenants/${tenants.currentId}/members/${m.id}`)
    showRemove.value = false
    removeTarget.value = null
    load()
  } catch (e) { error.value = (e as Error).message; showRemove.value = false } finally { removing.value = false }
}

async function setRole(id: string, r: string) {
  await api.patch(`/tenants/${tenants.currentId}/members/${id}`, { role: r })
  load()
}

function copyGenerated() {
  navigator.clipboard?.writeText(generated.value)
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[var(--color-text)]">Members</h1>
        <p class="mt-1 text-sm text-[var(--color-text-muted)]">Who has access to {{ tenants.current?.name }} and what they can do.</p>
      </div>
      <AppButton size="sm" @click="openAdd"><UserPlus class="mr-1.5 h-4 w-4" /> Add member</AppButton>
    </div>

    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
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
                @change="setRole(m.id, ($event.target as HTMLSelectElement).value)"
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
          <tr v-if="!loading && !members.length">
            <td colspan="3" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No members yet.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2 class="mb-2 mt-6 text-lg font-semibold text-[var(--color-text)]">Pending invitations</h2>
    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <tbody>
          <tr v-for="i in invites" :key="i.id" class="border-t border-[var(--color-border)]">
            <td class="px-3 py-2.5">
              <span class="flex items-center gap-2"><Mail class="h-3.5 w-3.5 text-[var(--color-text-muted)]" /> {{ i.email }}</span>
            </td>
            <td class="px-3 py-2.5 text-xs text-[var(--color-text-muted)]">{{ roleLabels[i.role] || i.role }}</td>
          </tr>
          <tr v-if="!invites.length"><td class="px-4 py-8 text-center text-sm text-[var(--color-text-muted)]">No pending invitations</td></tr>
        </tbody>
      </table>
    </div>

    <!-- Add member modal -->
    <AppModal :open="showAdd" title="Add a member" description="Choose how this member gets access." @close="showAdd = false">
      <div class="flex flex-col gap-3">
        <AppInput v-model="email" label="Email" placeholder="jane@acme.com" autofocus />

        <div v-if="smtpConfigured" class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Delivery method</label>
          <div class="grid grid-cols-1 gap-1.5">
            <label class="flex cursor-pointer items-start gap-2.5 rounded-[var(--radius-md)] border p-3" :class="method === 'invite' ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)]' : 'border-[var(--color-border)]'" @click="method = 'invite'">
              <input type="radio" :checked="method === 'invite'" class="mt-0.5 h-4 w-4 accent-[var(--color-primary)]" @change="method = 'invite'" />
              <span class="text-sm text-[var(--color-text)]">Send invitation email <span class="block text-xs text-[var(--color-text-muted)]">Member receives a sign-up link by email.</span></span>
            </label>
            <label class="flex cursor-pointer items-start gap-2.5 rounded-[var(--radius-md)] border p-3" :class="method === 'password' ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)]' : 'border-[var(--color-border)]'" @click="method = 'password'">
              <input type="radio" :checked="method === 'password'" class="mt-0.5 h-4 w-4 accent-[var(--color-primary)]" @change="method = 'password'" />
              <span class="text-sm text-[var(--color-text)]">Set a password for them <span class="block text-xs text-[var(--color-text-muted)]">Share it securely yourself. No email needed.</span></span>
            </label>
          </div>
        </div>

        <div v-if="!smtpConfigured" class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Delivery method</label>
          <div class="grid grid-cols-1 gap-1.5">
            <label class="flex cursor-pointer items-start gap-2.5 rounded-[var(--radius-md)] border p-3" :class="method === 'password' ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)]' : 'border-[var(--color-border)]'" @click="method = 'password'">
              <input type="radio" :checked="method === 'password'" class="mt-0.5 h-4 w-4 accent-[var(--color-primary)]" @change="method = 'password'" />
              <span class="text-sm text-[var(--color-text)]">Set a password for them <span class="block text-xs text-[var(--color-text-muted)]">Share it securely yourself. No email server needed.</span></span>
            </label>
            <label class="flex cursor-pointer items-start gap-2.5 rounded-[var(--radius-md)] border p-3" :class="method === 'activation' ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)]' : 'border-[var(--color-border)]'" @click="method = 'activation'">
              <input type="radio" :checked="method === 'activation'" class="mt-0.5 h-4 w-4 accent-[var(--color-primary)]" @change="method = 'activation'" />
              <span class="text-sm text-[var(--color-text)]">Activation — member sets their own password <span class="block text-xs text-[var(--color-text-muted)]">They enter their email at your /activate page and choose a password. One-time only.</span></span>
            </label>
          </div>
        </div>

        <div v-if="method === 'password'" class="flex flex-col gap-1">
          <AppInput v-model="password" label="Password" type="password" placeholder="At least 8 characters" autocomplete="new-password" />
          <p class="text-xs text-[var(--color-text-muted)]">You must share this password securely — it will not be emailed.</p>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Role</label>
          <select v-model="role" class="h-12 w-full rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
            <option value="viewer">Viewer — read-only access</option>
            <option value="member">Member — upload and manage files</option>
            <option value="tenant_admin">Admin — manage members and settings</option>
            <option value="tenant_owner">Owner — full control</option>
          </select>
        </div>

        <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
      </div>
      <template #footer>
        <AppButton variant="ghost" :disabled="adding" @click="showAdd = false">Cancel</AppButton>
        <AppButton :loading="adding" @click="addMember"><UserPlus class="mr-1.5 h-4 w-4" /> Add member</AppButton>
      </template>
    </AppModal>

    <!-- Generated password confirm -->
    <AppModal :open="!!generated" title="Member created" description="Share this password securely with the new member. It won't be shown again." @close="generated = ''">
      <div class="flex items-center gap-2 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] p-3">
        <code class="flex-1 select-all font-mono text-sm text-[var(--color-text)]">{{ generated }}</code>
        <AppButton size="sm" variant="secondary" @click="copyGenerated">Copy</AppButton>
      </div>
      <template #footer>
        <AppButton @click="generated = ''">Done</AppButton>
      </template>
    </AppModal>

    <!-- Remove member confirm -->
    <AppModal :open="showRemove" title="Remove member" :description="`Remove ${removeTarget?.display_name || removeTarget?.email || 'this member'} from ${tenants.current?.name}? They'll lose access immediately.`" @close="showRemove = false">
      <template #footer>
        <AppButton variant="ghost" :disabled="removing" @click="showRemove = false">Cancel</AppButton>
        <AppButton variant="destructive" :loading="removing" @click="doRemove"><Trash2 class="mr-1.5 h-4 w-4" /> Remove</AppButton>
      </template>
    </AppModal>
  </div>
</template>
