<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Mail, Lock, ArrowRight, Boxes, ArrowLeft, ShieldCheck, KeyRound, Server } from 'lucide-vue-next'
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
  <div class="flex min-h-screen bg-[var(--color-background)]">
    <!-- Brand panel -->
    <div class="relative hidden w-[44%] flex-col justify-between overflow-hidden border-r border-[var(--color-border)] bg-[var(--color-primary)] p-10 lg:flex">
      <div class="pointer-events-none absolute inset-0 opacity-30" style="background: radial-gradient(60% 50% at 20% 10%, rgba(255,255,255,0.35), transparent 60%), radial-gradient(50% 40% at 90% 90%, rgba(0,0,0,0.25), transparent 60%)" />
      <div class="relative flex items-center gap-2.5">
        <span class="flex h-9 w-9 items-center justify-center rounded-[var(--radius-md)] bg-white/15 backdrop-blur">
          <Boxes class="h-4 w-4 text-white" />
        </span>
        <span class="text-sm font-semibold tracking-tight text-white">Bloberry</span>
      </div>

      <div class="relative">
        <h1 class="text-balance text-3xl font-bold leading-tight tracking-tight text-white">
          Your files. Your storage. Your rules.
        </h1>
        <p class="mt-3 max-w-sm text-sm leading-relaxed text-white/80">
          One dashboard, six providers, and a permission model that matches how your team actually works.
        </p>
        <ul class="mt-8 space-y-4">
          <li class="flex items-center gap-3 text-sm text-white/90">
            <ShieldCheck class="h-4 w-4 shrink-0" /> Envelope-encrypted credentials
          </li>
          <li class="flex items-center gap-3 text-sm text-white/90">
            <KeyRound class="h-4 w-4 shrink-0" /> Folder-level permissions
          </li>
          <li class="flex items-center gap-3 text-sm text-white/90">
            <Server class="h-4 w-4 shrink-0" /> Self-hosted, one binary
          </li>
        </ul>
      </div>

      <p class="relative text-xs text-white/70">Storage-agnostic object service</p>
    </div>

    <!-- Form panel -->
    <div class="flex flex-1 flex-col items-center justify-center px-6 py-10">
      <div class="w-full max-w-sm">
        <button class="mb-8 inline-flex items-center gap-1.5 text-sm text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-primary)]" @click="home">
          <ArrowLeft class="h-4 w-4" /> Back to home
        </button>

        <div class="mb-8 lg:hidden">
          <span class="inline-flex h-10 w-10 items-center justify-center rounded-[var(--radius-md)] bg-[var(--color-primary)] shadow-[var(--shadow-sm)]">
            <Boxes class="h-5 w-5 text-[var(--color-on-primary)]" />
          </span>
        </div>

        <h1 class="text-2xl font-bold tracking-tight text-[var(--color-text)]">Welcome back</h1>
        <p class="mt-1.5 text-sm text-[var(--color-text-muted)]">Sign in to your storage.</p>

        <form class="mt-8 flex flex-col gap-4" @submit.prevent="submit">
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

        <div class="mt-6 border-t border-[var(--color-border)] pt-5 text-center">
          <button class="text-sm font-medium text-[var(--color-primary)] hover:underline" @click="loginOtp">Use a one-time code instead</button>
        </div>
      </div>
    </div>
  </div>
</template>
