<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Mail, Lock, ArrowRight, Boxes } from 'lucide-vue-next'
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
  <div class="flex min-h-screen items-center justify-center bg-[#012624] px-4">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <span class="mx-auto flex h-12 w-12 items-center justify-center rounded-[16px] bg-[#003734]">
          <Boxes class="h-6 w-6 text-[#cbfffc]" />
        </span>
        <h1 class="mt-4 text-[28px] font-medium leading-[1.0] tracking-[-0.02em] text-white">Bloberry</h1>
        <p class="mt-2 text-sm text-[#bbc7c6]">Sign in to your storage</p>
      </div>

      <div class="rounded-[16px] bg-[#003734] p-8">
        <form class="flex flex-col gap-4" @submit.prevent="submit">
          <AppInput v-model="email" label="Email" type="email" placeholder="you@company.com" :icon="Mail" />
          <AppInput v-model="password" label="Password" type="password" placeholder="••••••••" :icon="Lock" :error="error" />
          <AppButton type="submit" variant="gradient" :loading="loading">
            Log in
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </form>

        <div class="mt-5 flex flex-col gap-1 text-sm">
          <button class="text-left text-[#cbfffc] hover:underline" @click="loginOtp">Use a code instead</button>
          <button class="text-left text-[#cbfffc] hover:underline" @click="loginForgot">Forgot password?</button>
        </div>
      </div>

      <p class="mt-6 text-center text-xs text-[#707777]">Powered by Bloberry · storage-agnostic object service</p>
    </div>
  </div>
</template>
