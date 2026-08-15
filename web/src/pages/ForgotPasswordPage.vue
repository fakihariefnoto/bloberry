<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../lib/api'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const email = ref('')
const sent = ref(false)
const error = ref('')
const loading = ref(false)

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await api.post('/auth/forgot-password', { email: email.value })
    sent.value = true
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
      <h1 class="text-xl font-bold text-[var(--color-text)]">Reset your password</h1>
      <p v-if="!sent" class="mt-1 text-sm text-[var(--color-text-muted)]">Enter your email and we'll send a reset link.</p>
      <p v-else class="mt-2 text-sm text-[var(--color-success)]">If that email exists, a reset link is on its way.</p>

      <div v-if="!sent" class="mt-6 flex flex-col gap-4">
        <AppInput v-model="email" label="Email" type="email" :error="error" />
        <AppButton :loading="loading" @click="submit">Send reset link</AppButton>
      </div>
    </div>
  </div>
</template>
