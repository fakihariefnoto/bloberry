<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Mail, Lock, User, ArrowRight, Boxes } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const email = ref('')
const password = ref('')
const displayName = ref('')
const error = ref('')
const loading = ref(false)

async function accept() {
  loading.value = true
  error.value = ''
  try {
    await auth.signup(route.params.token as string, email.value, password.value, displayName.value)
    router.push({ name: 'files' })
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}

const login = () => router.push({ name: 'login' })
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-[#012624] px-4">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <span class="mx-auto flex h-12 w-12 items-center justify-center rounded-[16px] bg-[#003734]">
          <Boxes class="h-6 w-6 text-[#cbfffc]" />
        </span>
        <h1 class="mt-4 text-[28px] font-medium leading-[1.0] tracking-[-0.02em] text-white">You're invited</h1>
        <p class="mt-2 text-sm text-[#bbc7c6]">Create your account to join this workspace.</p>
      </div>

      <div class="rounded-[16px] bg-[#003734] p-8">
        <div class="flex flex-col gap-4">
          <AppInput v-model="email" label="Email" type="email" :icon="Mail" dark />
          <AppInput v-model="displayName" label="Display name" :icon="User" dark />
          <AppInput v-model="password" label="Password" type="password" :icon="Lock" dark :error="error" />
          <AppButton variant="gradient" :loading="loading" @click="accept">
            Create account
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </div>
        <p class="mt-5 text-sm text-[#707777]">
          Already have an account?
          <button class="text-[#cbfffc] hover:underline" @click="login">Log in</button>
        </p>
      </div>
    </div>
  </div>
</template>
