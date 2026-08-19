<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Mail, Lock, ArrowRight, Boxes, ArrowLeft } from 'lucide-vue-next'
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

onMounted(() => auth.loadSmtpConfig())

const loginOtp = () => router.push({ name: 'otp-login' })
const loginForgot = () => router.push({ name: 'forgot-password' })
const home = () => router.push({ name: 'landing' })

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
  <div class="relative flex min-h-screen items-center justify-center overflow-hidden bg-[var(--color-background)] px-4 py-10">
    <!-- Ambient background -->
    <div class="pointer-events-none fixed inset-0">
      <div class="absolute inset-0 opacity-70" style="background: radial-gradient(55% 45% at 50% 0%, rgba(139,125,235,0.12), transparent 70%)" />
      <div class="absolute inset-0 opacity-50" style="background: radial-gradient(45% 35% at 85% 80%, rgba(30,20,107,0.25), transparent 70%)" />
    </div>

    <div class="relative w-full max-w-sm">
      <!-- Back link -->
      <button class="mb-8 inline-flex items-center gap-1.5 text-sm text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-primary)]" @click="home">
        <ArrowLeft class="h-4 w-4" /> Back to home
      </button>

      <!-- Brand mark -->
      <div class="mb-6 flex flex-col items-center text-center">
        <span class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-[var(--color-primary)] to-[var(--color-accent-deep)] shadow-[0_0_24px_rgba(139,125,235,0.35)]">
          <Boxes class="h-6 w-6 text-[var(--color-on-primary)]" />
        </span>
        <h1 class="mt-4 text-2xl font-bold tracking-tight text-[var(--color-text)]">Welcome back</h1>
        <p class="mt-1.5 text-sm text-[var(--color-text-muted)]">Sign in to your storage.</p>
      </div>

      <!-- Form -->
      <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-8 shadow-[var(--shadow-md)]">
        <form class="flex flex-col gap-4" @submit.prevent="submit">
          <AppInput v-model="email" label="Email" type="email" placeholder="you@company.com" :icon="Mail" autocomplete="username" />
          <AppInput v-model="password" label="Password" type="password" placeholder="Your password" :icon="Lock" :error="error" autocomplete="current-password" />
          <div class="flex items-center justify-end">
            <button type="button" class="text-xs font-medium text-[var(--color-primary)] hover:underline" @click="loginForgot">Forgot password?</button>
          </div>
          <AppButton type="submit" :loading="loading">
            Sign in
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </form>

        <div v-if="auth.smtpConfigured" class="mt-6 border-t border-[var(--color-border)] pt-5 text-center">
          <button class="text-sm font-medium text-[var(--color-primary)] hover:underline" @click="loginOtp">Use a one-time code instead</button>
        </div>
      </div>

      <p class="mt-6 text-center text-xs text-[var(--color-text-muted)]">Storage-agnostic object service · self-hosted</p>
    </div>
  </div>
</template>
