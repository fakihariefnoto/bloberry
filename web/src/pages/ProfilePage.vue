<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../lib/api'
import { useAuthStore } from '../stores/auth'
import AppButton from '../components/ui/AppButton.vue'

const auth = useAuthStore()
const displayName = ref(auth.user?.display_name || '')
const locale = ref(auth.user?.settings?.locale || 'en')
const saved = ref(false)
const error = ref('')

async function save() {
  error.value = ''
  saved.value = false
  try {
    await api.patch('/users/me', { display_name: displayName.value, locale: locale.value })
    saved.value = true
    await auth.me()
  } catch (e) { error.value = (e as Error).message }
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-[var(--color-text)]">Profile</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">{{ auth.user?.email }}</p>

    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Display name</label>
          <input v-model="displayName" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-xs font-medium text-[var(--color-text-muted)]">Locale</label>
          <select v-model="locale" class="h-12 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 text-sm">
            <option value="en">English</option>
            <option value="id">Bahasa Indonesia</option>
          </select>
        </div>
        <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
        <p v-if="saved" class="text-xs text-[var(--color-success)]">Saved.</p>
        <div><AppButton @click="save">Save</AppButton></div>
      </div>
    </div>
  </div>
</template>
