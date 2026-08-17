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
const adding = ref(false)

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
onMounted(load)

async function addByEmail() {
  error.value = ''
  adding.value = true
  try {
    await api.post(`/tenants/${tenants.currentId}/members`, { email: email.value, role: role.value })
    showAdd.value = false
    email.value = ''
    load()
  } catch (e) { error.value = (e as Error).message } finally { adding.value = false }
}

async function invite() {
  error.value = ''
  adding.value = true
  try {
    await api.post(`/tenants/${tenants.currentId}/invitations`, { email: email.value, role: role.value })
    showAdd.value = false
    email.value = ''
    load()
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
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[var(--color-text)]">Members</h1>
        <p class="mt-1 text-sm text-[var(--color-text-muted)]">Who has access to {{ tenants.current?.name }} and what they can do.</p>
      </div>
      <AppButton size="sm" @click="showAdd = true"><UserPlus class="mr-1.5 h-4 w-4" /> Add member</AppButton>
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

    <!-- Add member modal: add an existing user by email -->
    <AppModal :open="showAdd" title="Add a member" description="Add someone who already has an account by email, or invite them by email." @close="showAdd = false">
      <div class="flex flex-col gap-3">
        <AppInput v-model="email" label="Email" placeholder="jane@acme.com" autofocus />
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
        <AppButton variant="ghost" @click="showAdd = false">Cancel</AppButton>
        <AppButton variant="secondary" :disabled="adding" @click="invite"><Mail class="mr-1.5 h-4 w-4" /> Invite</AppButton>
        <AppButton :loading="adding" @click="addByEmail"><UserPlus class="mr-1.5 h-4 w-4" /> Add</AppButton>
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
