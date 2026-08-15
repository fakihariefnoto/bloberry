<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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
    <div class="w-full max-w-sm rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6 shadow-[var(--shadow-md)]">
      <h1 class="text-xl font-bold text-[var(--color-text)]">You're invited</h1>
      <p class="mt-1 text-sm text-[var(--color-text-muted)]">Create your account to join this workspace.</p>
      <div class="mt-6 flex flex-col gap-4">
        <AppInput v-model="email" label="Email" type="email" />
        <AppInput v-model="displayName" label="Display name" />
        <AppInput v-model="password" label="Password" type="password" :error="error" />
        <AppButton :loading="loading" @click="accept">Create account</AppButton>
      </div>
      <p class="mt-4 text-sm text-[var(--color-text-muted)]">
        Already have an account?
        <button class="text-[var(--color-primary)] hover:underline" @click="login">Log in</button>
      </p>
    </div>
  </div>
</template>
