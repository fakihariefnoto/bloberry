<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const loginOtp = () => router.push({ name: 'otp-login' })
const loginForgot = () => router.push({ name: 'forgot-password' })

async function submit() {
  error.value = ''
  loading.value = true
  try {
    const res = await auth.login(email.value, password.value)
    if (res.totpRequired) {
      router.push({ name: 'login', query: { ...route.query, totp: res.pending } })
      return
    }
    router.push((route.query.next as string) || { name: 'files' })
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-[var(--color-background)] px-4">
    <div class="w-full max-w-sm rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6 shadow-[var(--shadow-md)]">
      <h1 class="text-2xl font-bold text-[var(--color-text)]">Bloberry</h1>
      <p class="mt-1 text-sm text-[var(--color-text-muted)]">Sign in to your storage</p>

      <form class="mt-6 flex flex-col gap-4" @submit.prevent="submit">
        <AppInput v-model="email" label="Email" type="email" placeholder="you@company.com" />
        <AppInput v-model="password" label="Password" type="password" placeholder="••••••••" :error="error" />
        <AppButton type="submit" :loading="loading">Log in</AppButton>
      </form>

      <div class="mt-4 flex flex-col gap-1 text-sm">
        <button class="text-left text-[var(--color-primary)] hover:underline" @click="loginOtp">Use a code instead</button>
        <button class="text-left text-[var(--color-primary)] hover:underline" @click="loginForgot">Forgot password?</button>
      </div>
    </div>
  </div>
</template>
