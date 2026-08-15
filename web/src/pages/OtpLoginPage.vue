<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const auth = useAuthStore()
const router = useRouter()

const email = ref('')
const code = ref('')
const sent = ref(false)
const error = ref('')
const loading = ref(false)

async function request() {
  error.value = ''
  try {
    await auth.requestOtp(email.value)
    sent.value = true
  } catch (e) {
    error.value = (e as Error).message
  }
}

async function verify() {
  loading.value = true
  try {
    await auth.verifyOtp(email.value, code.value)
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
    <div class="w-full max-w-sm rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6 shadow-[var(--shadow-md)]">
      <h1 class="text-xl font-bold text-[var(--color-text)]">Email login code</h1>
      <p class="mt-1 text-sm text-[var(--color-text-muted)]">We'll email you a 6-digit code.</p>

      <div v-if="!sent" class="mt-6 flex flex-col gap-4">
        <AppInput v-model="email" label="Email" type="email" :error="error" />
        <AppButton @click="request">Send code</AppButton>
      </div>
      <div v-else class="mt-6 flex flex-col gap-4">
        <AppInput v-model="code" label="6-digit code" :error="error" />
        <AppButton :loading="loading" @click="verify">Verify</AppButton>
      </div>
    </div>
  </div>
</template>
