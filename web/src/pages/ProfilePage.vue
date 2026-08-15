<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Mail, Globe, CheckCircle2, Building2, ShieldCheck } from 'lucide-vue-next'
import { api } from '../lib/api'
import { useAuthStore } from '../stores/auth'
import { useTenantStore } from '../stores/tenant'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import PageHeader from '../components/ui/PageHeader.vue'

const auth = useAuthStore()
const tenants = useTenantStore()

const displayName = ref(auth.user?.display_name || '')
const locale = ref(auth.user?.settings?.locale || 'en')
const notifications = ref(auth.user?.settings?.notifications_enabled ?? true)
const saved = ref('')
const error = ref('')
const saving = ref(false)

onMounted(async () => {
  await auth.me()
  await tenants.load()
  displayName.value = auth.user?.display_name || ''
  locale.value = auth.user?.settings?.locale || 'en'
  notifications.value = auth.user?.settings?.notifications_enabled ?? true
})

async function save() {
  error.value = ''
  saved.value = ''
  saving.value = true
  try {
    await api.patch('/users/me', { display_name: displayName.value, locale: locale.value, notifications_enabled: notifications.value })
    saved.value = 'Saved.'
    await auth.me()
  } catch (e) { error.value = (e as Error).message } finally { saving.value = false }
}

const initial = (auth.user?.display_name || auth.user?.email || '?').charAt(0).toUpperCase()
</script>

<template>
  <div>
    <PageHeader title="Profile" :description="auth.user?.email" />

    <!-- Identity card -->
    <div class="max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <div class="flex items-center gap-4">
        <span class="flex h-16 w-16 items-center justify-center rounded-full bg-[var(--color-primary-subtle)] text-2xl font-bold text-[var(--color-primary)]">
          {{ initial }}
        </span>
        <div>
          <p class="text-lg font-bold text-[var(--color-text)]">{{ auth.user?.display_name || '—' }}</p>
          <p class="flex items-center gap-1.5 text-sm text-[var(--color-text-muted)]">
            <Mail class="h-3.5 w-3.5" /> {{ auth.user?.email }}
            <CheckCircle2 v-if="auth.user?.email_verified" class="h-3.5 w-3.5 text-[var(--color-success)]" />
          </p>
          <p v-if="auth.user?.platform_role" class="mt-1 inline-flex items-center gap-1 rounded-full bg-[var(--color-primary-subtle)] px-2 py-0.5 text-xs font-medium text-[var(--color-primary)]">
            <ShieldCheck class="h-3 w-3" /> Platform admin
          </p>
        </div>
      </div>

      <div class="mt-5 border-t border-[var(--color-border)] pt-4">
        <p class="text-xs font-medium text-[var(--color-text-muted)]">Memberships</p>
        <ul class="mt-2 space-y-2">
          <li v-for="t in tenants.list" :key="t.id" class="flex items-center justify-between text-sm">
            <span class="flex items-center gap-2 text-[var(--color-text)]">
              <Building2 class="h-4 w-4 text-[var(--color-primary)]" /> {{ t.name }}
            </span>
            <span class="rounded-full bg-[var(--color-surface)] px-2 py-0.5 text-xs text-[var(--color-text-muted)]">
              {{ tenants.memberships.find((m) => m.tenant_id === t.id)?.role || '—' }}
            </span>
          </li>
          <li v-if="!tenants.list.length" class="text-sm text-[var(--color-text-muted)]">No memberships.</li>
        </ul>
      </div>
    </div>

    <!-- Settings -->
    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <h2 class="text-lg font-semibold text-[var(--color-text)]">Settings</h2>
      <div class="mt-5 flex flex-col gap-4">
        <AppInput v-model="displayName" label="Display name" placeholder="Your name" />
        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Locale</label>
          <select v-model="locale" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm outline-none focus:border-[var(--color-primary)] focus:border-2">
            <option value="en">English</option>
            <option value="id">Bahasa Indonesia</option>
          </select>
        </div>
        <label class="flex items-center justify-between rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-3">
          <span class="flex items-center gap-2 text-sm text-[var(--color-text)]">
            <Globe class="h-4 w-4 text-[var(--color-primary)]" /> Notifications
          </span>
          <input v-model="notifications" type="checkbox" class="h-4 w-4 accent-[var(--color-primary)]" />
        </label>
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <p v-if="saved" class="text-xs text-[var(--color-success)]">{{ saved }}</p>
        <div class="flex justify-end"><AppButton :loading="saving" @click="save">Save changes</AppButton></div>
      </div>
    </div>
  </div>
</template>
