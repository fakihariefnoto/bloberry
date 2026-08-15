<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../lib/api'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const route = useRoute()
const router = useRouter()
const password = ref('')
const confirm = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  if (password.value !== confirm.value) {
    error.value = 'Passwords do not match'
    return
  }
  loading.value = true
  error.value = ''
  try {
    await api.post('/auth/reset-password', { token: route.query.token, new_password: password.value })
    router.push({ name: 'login' })
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
      <h1 class="text-xl font-bold text-[var(--color-text)]">Set a new password</h1>
      <div class="mt-6 flex flex-col gap-4">
        <AppInput v-model="password" label="New password" type="password" />
        <AppInput v-model="confirm" label="Confirm password" type="password" :error="error" />
        <AppButton :loading="loading" @click="submit">Set password</AppButton>
      </div>
    </div>
  </div>
</template>
