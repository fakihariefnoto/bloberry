<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, KeyRound, Copy } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

const route = useRoute()
const router = useRouter()
const appId = route.params.appId as string

interface AccessKey { id: string; last_four: string; prefix: string; permissions: string[]; scope_folder_ids: string[]; revoked_at?: string; expires_at?: string }
interface Application { id: string; name: string; description: string }

const app = ref<Application | null>(null)
const keys = ref<AccessKey[]>([])
const error = ref('')
const showCreate = ref(false)
const secret = ref('')
const perms = ref(['read'])

async function load() {
  app.value = await api.get<Application>(`/applications/${appId}`)
  keys.value = await api.get<AccessKey[]>(`/applications/${appId}/keys`)
}
onMounted(load)

async function createKey() {
  error.value = ''
  try {
    const res = await api.post<{ secret: string }>(`/applications/${appId}/keys`, { permissions: perms.value })
    secret.value = res.secret
    showCreate.value = true
    load()
  } catch (e) { error.value = (e as Error).message }
}

async function revoke(id: string) {
  await api.delete(`/applications/${appId}/keys/${id}`)
  load()
}

const back = () => router.push({ name: 'applications' })
function copyText(t: string) {
  window.navigator.clipboard?.writeText(t)
}
</script>

<template>
  <div>
    <button class="mb-4 flex items-center gap-1 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="back">
      <ArrowLeft class="h-4 w-4" /> Back to applications
    </button>

    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-[var(--color-text)]">{{ app?.name }}</h1>
        <p v-if="app?.description" class="mt-1 text-sm text-[var(--color-text-muted)]">{{ app.description }}</p>
      </div>
      <AppButton size="sm" @click="createKey"><KeyRound class="h-4 w-4" /> Create key</AppButton>
    </div>

    <p v-if="error" class="mb-3 text-sm text-[var(--color-error)]">{{ error }}</p>

    <!-- Key-created modal: secret shown exactly once (PRD D5) -->
    <div v-if="showCreate" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="showCreate = false">
      <div class="w-full max-w-md rounded-[var(--radius-lg)] bg-[var(--color-surface-raised)] p-6 shadow-[var(--shadow-lg)]">
        <h2 class="text-lg font-semibold text-[var(--color-text)]">Access key created</h2>
        <p class="mt-1 text-sm text-[var(--color-warning)]">You won't see this again. Copy it now.</p>
        <div class="mt-4 flex items-center gap-2 rounded-[var(--radius-sm)] bg-[var(--color-surface)] p-3">
          <code class="flex-1 break-all font-mono text-xs text-[var(--color-text)]">{{ secret }}</code>
          <button class="text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="copyText(secret)">
            <Copy class="h-4 w-4" />
          </button>
        </div>
        <div class="mt-4"><AppButton @click="showCreate = false">I've saved it</AppButton></div>
      </div>
    </div>

    <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Key</th>
            <th class="px-3 py-2">Permissions</th>
            <th class="px-3 py-2">State</th>
            <th class="px-3 py-2"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="k in keys" :key="k.id" class="border-t border-[var(--color-border)]">
            <td class="px-3 py-2.5 font-mono text-xs">{{ k.prefix }}••••{{ k.last_four }}</td>
            <td class="px-3 py-2.5">{{ k.permissions.join(', ') }}</td>
            <td class="px-3 py-2.5">
              <span v-if="k.revoked_at" class="text-xs text-[var(--color-error)]">Revoked</span>
              <span v-else class="text-xs text-[var(--color-success)]">Active</span>
            </td>
            <td class="px-3 py-2.5">
              <button v-if="!k.revoked_at" class="text-xs text-[var(--color-error)] hover:underline" @click="revoke(k.id)">Revoke</button>
            </td>
          </tr>
          <tr v-if="!keys.length">
            <td colspan="4" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No access keys yet</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
