<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { KeyRound } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import { useAuthStore } from '../stores/auth'

interface AuthData { access_token: string; refresh_token: string }

const router = useRouter()
const auth = useAuthStore()
const email = ref('')
const password = ref('')
const displayName = ref('')
const error = ref('')
const loading = ref(false)

async function activate() {
  error.value = ''
  if (password.value.length < 8) {
    error.value = 'Password must be at least 8 characters.'
    return
  }
  loading.value = true
  try {
    const res = await api.post<AuthData>('/auth/activate', {
      email: email.value,
      password: password.value,
      display_name: displayName.value || undefined,
      platform: 'web',
    })
    auth.applySession(res)
    await auth.me()
    router.push({ name: 'files' })
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-[var(--color-background)] px-4">
    <div class="w-full max-w-md">
      <div class="mb-6 flex flex-col items-center gap-3">
        <span class="flex h-12 w-12 items-center justify-center rounded-[var(--radius-lg)] bg-[var(--color-primary)]">
          <KeyRound class="h-6 w-6 text-[var(--color-on-primary)]" />
        </span>
        <h1 class="text-2xl font-bold tracking-tight text-[var(--color-text)]">Activate your account</h1>
        <p class="text-center text-sm text-[var(--color-text-muted)]">
          A project admin added you. Set your password to activate — this only works once for this email.
        </p>
      </div>

      <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
        <div class="flex flex-col gap-4">
          <AppInput v-model="email" label="Email" placeholder="jane@acme.com" type="email" autofocus />
          <AppInput v-model="displayName" label="Display name (optional)" placeholder="Jane Doe" />
          <AppInput v-model="password" label="Password" type="password" placeholder="At least 8 characters" autocomplete="new-password" />

          <p v-if="error" class="rounded-[var(--radius-sm)] border border-[var(--color-error)] bg-[var(--color-error)]/10 px-3 py-2 text-xs text-[var(--color-error)]">{{ error }}</p>
          <AppButton :loading="loading" @click="activate">Activate account</AppButton>
          <p class="text-center text-xs text-[var(--color-text-muted)]">
            Already activated? <RouterLink class="text-[var(--color-primary)] hover:underline" to="/login">Sign in</RouterLink>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
