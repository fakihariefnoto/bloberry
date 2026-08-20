<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Mail, Lock, User, ArrowRight } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const email = ref('')
const password = ref('')
const displayName = ref('')
const error = ref('')
const loading = ref(false)

async function accept() {
  loading.value = true
  error.value = ''
  try {
    await auth.signup(route.params.token as string, email.value, password.value, displayName.value)
    router.push({ name: 'files' })
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}

const login = () => router.push({ name: 'login' })
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-[var(--color-background)] px-4">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <img src="/bloberry-icon.png" alt="Bloberry" class="mx-auto h-12 w-12 rounded-[var(--radius-lg)] object-cover" />
        <h1 class="mt-4 text-3xl font-bold tracking-tight text-[var(--color-text)]">You're invited</h1>
        <p class="mt-2 text-sm text-[var(--color-text-muted)]">Create your account to join this workspace.</p>
      </div>

      <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-8 shadow-[var(--shadow-md)]">
        <div class="flex flex-col gap-4">
          <AppInput v-model="email" label="Email" type="email" :icon="Mail" />
          <AppInput v-model="displayName" label="Display name" :icon="User" />
          <AppInput v-model="password" label="Password" type="password" :icon="Lock" :error="error" />
          <AppButton :loading="loading" @click="accept">
            Create account
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </div>
        <p class="mt-5 text-sm text-[var(--color-text-muted)]">
          Already have an account?
          <button class="text-[var(--color-primary)] hover:underline" @click="login">Log in</button>
        </p>
      </div>
    </div>
  </div>
</template>
