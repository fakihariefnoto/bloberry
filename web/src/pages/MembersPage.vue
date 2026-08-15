<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../lib/api'
import { useTenantStore } from '../stores/tenant'
import AppButton from '../components/ui/AppButton.vue'

interface Membership { id: string; user_id: string; role: string; created_at: string }
interface Invitation { id: string; email: string; role: string }

const tenants = useTenantStore()
const members = ref<Membership[]>([])
const invites = ref<Invitation[]>([])
const loading = ref(false)
const showInvite = ref(false)
const email = ref('')
const role = ref('member')
const error = ref('')

async function load() {
  loading.value = true
  try {
    members.value = await api.get<Membership[]>(`/tenants/${tenants.currentId}/members`)
    invites.value = await api.get<Invitation[]>(`/tenants/${tenants.currentId}/invitations`)
  } catch { members.value = []; invites.value = [] } finally { loading.value = false }
}
onMounted(load)

async function invite() {
  error.value = ''
  try {
    await api.post(`/tenants/${tenants.currentId}/invitations`, { email: email.value, role: role.value })
    showInvite.value = false
    email.value = ''
    load()
  } catch (e) { error.value = (e as Error).message }
}

async function removeMember(id: string) {
  await api.delete(`/tenants/${tenants.currentId}/members/${id}`)
  load()
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
        <p class="mt-1 text-sm text-[var(--color-text-muted)]">Who has access to {{ tenants.current?.name }}.</p>
      </div>
      <AppButton size="sm" @click="showInvite = !showInvite">Invite member</AppButton>
    </div>

    <div v-if="showInvite" class="mb-4 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
      <div class="flex flex-col gap-3">
        <input v-model="email" placeholder="email@company.com" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <select v-model="role" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm">
          <option value="viewer">Viewer</option>
          <option value="member">Member</option>
          <option value="tenant_admin">Admin</option>
        </select>
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <div><AppButton size="sm" @click="invite">Send invitation</AppButton></div>
      </div>
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
            <td class="px-3 py-2.5 font-mono text-xs">{{ m.user_id }}</td>
            <td class="px-3 py-2.5">
              <select :value="m.role" class="rounded-[var(--radius-sm)] border border-[var(--color-border)] bg-[var(--color-surface)] px-2 py-1 text-xs" @change="setRole(m.id, ($event.target as HTMLSelectElement).value)">
                <option value="viewer">viewer</option>
                <option value="member">member</option>
                <option value="tenant_admin">tenant_admin</option>
                <option value="tenant_owner">tenant_owner</option>
              </select>
            </td>
            <td class="px-3 py-2.5 text-right">
              <button class="text-xs text-[var(--color-error)] hover:underline" @click="removeMember(m.id)">Remove</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2 class="mb-2 mt-6 text-lg font-semibold text-[var(--color-text)]">Pending invitations</h2>
    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <tbody>
          <tr v-for="i in invites" :key="i.id" class="border-t border-[var(--color-border)]">
            <td class="px-3 py-2.5">{{ i.email }}</td>
            <td class="px-3 py-2.5">{{ i.role }}</td>
          </tr>
          <tr v-if="!invites.length"><td class="px-4 py-8 text-center text-sm text-[var(--color-text-muted)]">No pending invitations</td></tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
