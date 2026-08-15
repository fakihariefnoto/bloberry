<script setup lang="ts">
import { ref } from 'vue'
import { Mail, ArrowRight, Boxes, CheckCircle2 } from 'lucide-vue-next'
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
  <div class="flex min-h-screen items-center justify-center bg-[#012624] px-4">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <span class="mx-auto flex h-12 w-12 items-center justify-center rounded-[16px] bg-[#003734]">
          <Boxes class="h-6 w-6 text-[#cbfffc]" />
        </span>
        <h1 class="mt-4 text-[28px] font-medium leading-[1.0] tracking-[-0.02em] text-white">Reset your password</h1>
        <p v-if="!sent" class="mt-2 text-sm text-[#bbc7c6]">Enter your email and we'll send a reset link.</p>
        <p v-else class="mt-2 text-sm text-[#cbfffc]">If that email exists, a reset link is on its way.</p>
      </div>

      <div v-if="!sent" class="rounded-[16px] bg-[#003734] p-8">
        <div class="flex flex-col gap-4">
          <AppInput v-model="email" label="Email" type="email" :icon="Mail" dark :error="error" />
          <AppButton variant="gradient" :loading="loading" @click="submit">
            Send reset link
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </div>
      </div>

      <div v-else class="rounded-[16px] bg-[#003734] p-8 text-center">
        <CheckCircle2 class="mx-auto h-8 w-8 text-[#cbfffc]" />
        <p class="mt-3 text-sm text-[#edfffe]">Check your inbox. The link is valid for 30 minutes.</p>
      </div>
    </div>
  </div>
</template>
