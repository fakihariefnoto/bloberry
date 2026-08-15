<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)

async function change() {
  error.value = ''
  success.value = ''
  if (newPassword.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }
  loading.value = true
  try {
    await api.post('/users/me/password', { current_password: currentPassword.value, new_password: newPassword.value })
    success.value = 'Password changed.'
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (e) { error.value = (e as Error).message } finally { loading.value = false }
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-[var(--color-text)]">Account settings</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">Password, security and 2FA.</p>

    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <h2 class="text-lg font-semibold text-[var(--color-text)]">Change password</h2>
      <div class="mt-4 flex flex-col gap-4">
        <input v-model="currentPassword" type="password" placeholder="Current password" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <input v-model="newPassword" type="password" placeholder="New password" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <input v-model="confirmPassword" type="password" placeholder="Confirm new password" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <p v-if="success" class="text-xs text-[var(--color-success)]">{{ success }}</p>
        <div><AppButton :loading="loading" @click="change">Change password</AppButton></div>
      </div>
    </div>

    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <h2 class="text-lg font-semibold text-[var(--color-text)]">Two-factor authentication</h2>
      <p class="mt-1 text-sm text-[var(--color-text-muted)]">
        TOTP codes from an authenticator app. Backup codes shown once when you enable it.
      </p>
    </div>
  </div>
</template>
