<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Mail, Lock, User, Building2, AtSign, ArrowRight, CheckCircle2, Boxes } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const done = ref(false)
const stillNeeded = ref(true)

const email = ref('')
const password = ref('')
const displayName = ref('')
const tenantName = ref('')
const tenantSlug = ref('')

onMounted(async () => {
  try {
    const st = await api.get<{ needs_setup: boolean }>('/setup/status')
    stillNeeded.value = st.needs_setup
  } catch {
    stillNeeded.value = true
  }
})

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await api.post('/setup', {
      email: email.value,
      password: password.value,
      display_name: displayName.value,
      tenant_name: tenantName.value,
      tenant_slug: tenantSlug.value || tenantName.value.toLowerCase().replace(/\s+/g, '-'),
    })
    done.value = true
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}

const goLogin = () => router.push({ name: 'login' })
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-[var(--color-background)] px-4">
    <div class="w-full max-w-md">
      <div class="mb-8 text-center">
        <span class="mx-auto flex h-12 w-12 items-center justify-center rounded-[var(--radius-lg)] bg-[var(--color-primary-subtle)]">
          <Boxes class="h-6 w-6 text-[var(--color-primary)]" />
        </span>
        <h1 class="mt-4 text-3xl font-bold tracking-tight text-[var(--color-text)]">Set up Bloberry</h1>
        <p class="mt-2 text-sm text-[var(--color-text-muted)]">One-time first-run — create your admin account and first tenant.</p>
      </div>

      <div v-if="done" class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-8 text-center shadow-[var(--shadow-md)]">
        <CheckCircle2 class="mx-auto h-10 w-10 text-[var(--color-success)]" />
        <h2 class="mt-4 text-xl font-bold text-[var(--color-text)]">You're ready</h2>
        <p class="mt-2 text-sm text-[var(--color-text-muted)]">Platform admin, tenant and local disk backend created.</p>
        <AppButton class="mt-6 w-full" @click="goLogin">Continue to sign in</AppButton>
      </div>

      <div v-else class="flex flex-col gap-4 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-8 shadow-[var(--shadow-md)]">
        <p v-if="!stillNeeded" class="mb-2 rounded-[var(--radius-sm)] bg-[var(--color-primary-subtle)] p-3 text-xs text-[var(--color-text)]">
          This install is already configured.
          <button class="text-[var(--color-primary)] underline" @click="goLogin">Go to sign in</button>
        </p>

        <AppInput v-model="email" label="Email" type="email" placeholder="admin@company.com" :icon="Mail" />
        <AppInput v-model="password" label="Password" type="password" placeholder="Create a strong password" :icon="Lock" />
        <AppInput v-model="displayName" label="Display name" placeholder="Your name" :icon="User" />
        <AppInput v-model="tenantName" label="Tenant name" placeholder="Acme Corp" :icon="Building2" />
        <AppInput v-model="tenantSlug" label="Tenant slug" placeholder="acme (blank = derived)" :icon="AtSign" :hint="'Used in URLs and short links'" />

        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>

        <AppButton :loading="loading" :disabled="!stillNeeded" @click="submit">
          Create admin &amp; tenant
          <ArrowRight class="ml-2 h-4 w-4" />
        </AppButton>
      </div>
    </div>
  </div>
</template>
