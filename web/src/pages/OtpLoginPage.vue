<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Mail, KeyRound, ArrowRight, Boxes } from 'lucide-vue-next'
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
  <div class="flex min-h-screen items-center justify-center bg-[#012624] px-4">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <span class="mx-auto flex h-12 w-12 items-center justify-center rounded-[16px] bg-[#003734]">
          <Boxes class="h-6 w-6 text-[#cbfffc]" />
        </span>
        <h1 class="mt-4 text-[28px] font-medium leading-[1.0] tracking-[-0.02em] text-white">Email login code</h1>
        <p class="mt-2 text-sm text-[#bbc7c6]">We'll email you a 6-digit code.</p>
      </div>

      <div class="rounded-[16px] bg-[#003734] p-8">
        <div v-if="!sent" class="flex flex-col gap-4">
          <AppInput v-model="email" label="Email" type="email" :icon="Mail" dark :error="error" />
          <AppButton variant="gradient" @click="request">
            Send code
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </div>
        <div v-else class="flex flex-col gap-4">
          <AppInput v-model="code" label="6-digit code" :icon="KeyRound" dark :error="error" />
          <AppButton variant="gradient" :loading="loading" @click="verify">
            Verify
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </div>
      </div>
    </div>
  </div>
</template>
