<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Lock, ArrowRight, Boxes } from 'lucide-vue-next'
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
  <div class="flex min-h-screen items-center justify-center bg-[#012624] px-4">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <span class="mx-auto flex h-12 w-12 items-center justify-center rounded-[16px] bg-[#003734]">
          <Boxes class="h-6 w-6 text-[#cbfffc]" />
        </span>
        <h1 class="mt-4 text-[28px] font-medium leading-[1.0] tracking-[-0.02em] text-white">Set a new password</h1>
      </div>

      <div class="rounded-[16px] bg-[#003734] p-8">
        <div class="flex flex-col gap-4">
          <AppInput v-model="password" label="New password" type="password" :icon="Lock" dark />
          <AppInput v-model="confirm" label="Confirm password" type="password" :icon="Lock" dark :error="error" />
          <AppButton variant="gradient" :loading="loading" @click="submit">
            Set password
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </div>
      </div>
    </div>
  </div>
</template>
