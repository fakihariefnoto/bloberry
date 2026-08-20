<script setup lang="ts">
import { ref } from 'vue'
import { Mail, ArrowRight, CheckCircle2 } from 'lucide-vue-next'
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
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <img src="/bloberry-icon.png" alt="Bloberry" class="mx-auto h-12 w-12 rounded-[var(--radius-lg)] object-cover" />
        <h1 class="mt-4 text-3xl font-bold tracking-tight text-[var(--color-text)]">Reset your password</h1>
        <p v-if="!sent" class="mt-2 text-sm text-[var(--color-text-muted)]">Enter your email and we'll send a reset link.</p>
        <p v-else class="mt-2 text-sm text-[var(--color-success)]">If that email exists, a reset link is on its way.</p>
      </div>

      <div v-if="!sent" class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-8 shadow-[var(--shadow-md)]">
        <div class="flex flex-col gap-4">
          <AppInput v-model="email" label="Email" type="email" :icon="Mail" :error="error" />
          <AppButton :loading="loading" @click="submit">
            Send reset link
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </div>
      </div>

      <div v-else class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-8 text-center">
        <CheckCircle2 class="mx-auto h-8 w-8 text-[var(--color-success)]" />
        <p class="mt-3 text-sm text-[var(--color-text-muted)]">Check your inbox. The link is valid for 30 minutes.</p>
      </div>
    </div>
  </div>
</template>
