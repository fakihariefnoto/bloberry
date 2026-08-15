<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Mail, Lock, User, Building2, AtSign, ArrowRight, CheckCircle2, Boxes } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const done = ref(false)
const stillNeeded = ref(true)

const email = ref('')
const password = ref('')
const displayName = ref('')
const tenantName = ref('')
const tenantSlug = ref('')

onMounted(async () => {
  try {
    const st = await api.get<{ needs_setup: boolean }>('/setup/status')
    stillNeeded.value = st.needs_setup
  } catch {
    stillNeeded.value = true
  }
})

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await api.post('/setup', {
      email: email.value,
      password: password.value,
      display_name: displayName.value,
      tenant_name: tenantName.value,
      tenant_slug: tenantSlug.value || tenantName.value.toLowerCase().replace(/\s+/g, '-'),
    })
    done.value = true
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}

const goLogin = () => router.push({ name: 'login' })
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-[#012624] px-4">
    <div class="w-full max-w-md">
      <div class="mb-8 text-center">
        <span class="mx-auto flex h-12 w-12 items-center justify-center rounded-[16px] bg-[#003734]">
          <Boxes class="h-6 w-6 text-[#cbfffc]" />
        </span>
        <h1 class="mt-4 text-[28px] font-medium leading-[1.0] tracking-[-0.02em] text-white">Set up Bloberry</h1>
        <p class="mt-2 text-sm text-[#bbc7c6]">One-time first-run — create your admin account and first tenant.</p>
      </div>

      <div v-if="done" class="rounded-[16px] bg-[#003734] p-8 text-center">
        <CheckCircle2 class="mx-auto h-10 w-10 text-[#cbfffc]" />
        <h2 class="mt-4 text-xl font-medium text-white">You're ready</h2>
        <p class="mt-2 text-sm text-[#bbc7c6]">Platform admin, tenant and local disk backend created.</p>
        <button
          class="mt-6 w-full rounded-md py-4 text-xs font-medium uppercase tracking-[0.08em] text-[#012624]"
          style="background: linear-gradient(90deg, rgb(0,130,124) 0%, rgb(203,255,252) 60%, rgb(250,209,255) 100%)"
          @click="goLogin"
        >
          Continue to sign in
        </button>
      </div>

      <div v-else class="flex flex-col gap-4 rounded-[16px] bg-[#003734] p-8">
        <p v-if="!stillNeeded" class="mb-2 rounded-md bg-[rgba(3,81,75,0.5)] p-3 text-xs text-[#edfffe]">
          This install is already configured.
          <button class="text-[#cbfffc] underline" @click="goLogin">Go to sign in</button>
        </p>

        <AppInput v-model="email" label="Email" type="email" placeholder="admin@company.com" :icon="Mail" />
        <AppInput v-model="password" label="Password" type="password" placeholder="Create a strong password" :icon="Lock" />
        <AppInput v-model="displayName" label="Display name" placeholder="Your name" :icon="User" />
        <AppInput v-model="tenantName" label="Tenant name" placeholder="Acme Corp" :icon="Building2" />
        <AppInput v-model="tenantSlug" label="Tenant slug" placeholder="acme (blank = derived)" :icon="AtSign" :hint="'Used in URLs and short links'" />

        <p v-if="error" class="text-xs text-[#fde9ff]">{{ error }}</p>

        <AppButton :loading="loading" :disabled="!stillNeeded" @click="submit">
          Create admin &amp; tenant
          <ArrowRight class="ml-2 h-4 w-4" />
        </AppButton>
      </div>
    </div>
  </div>
</template>
