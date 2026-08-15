<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Mail, Lock, User, Building2, AtSign, ArrowRight, ArrowLeft, CheckCircle2, Boxes, HardDrive } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppInput from '../components/ui/AppInput.vue'
import AppButton from '../components/ui/AppButton.vue'

const router = useRouter()
const step = ref(1)
const loading = ref(false)
const error = ref('')
const done = ref(false)
const stillNeeded = ref(true)

const email = ref('')
const password = ref('')
const displayName = ref('')
const tenantName = ref('')
const tenantSlug = ref('')

const steps = ['Admin account', 'Your tenant', 'Done']

onMounted(async () => {
  try {
    const st = await api.get<{ needs_setup: boolean }>('/setup/status')
    stillNeeded.value = st.needs_setup
  } catch {
    stillNeeded.value = true
  }
})

function next() {
  error.value = ''
  if (step.value === 1) {
    if (!email.value || !password.value) {
      error.value = 'Email and password are required.'
      return
    }
    if (password.value.length < 8) {
      error.value = 'Password must be at least 8 characters.'
      return
    }
  }
  if (step.value === 2) {
    if (!tenantName.value) {
      error.value = 'Tenant name is required.'
      return
    }
  }
  step.value += 1
}

function back() {
  error.value = ''
  step.value -= 1
}

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
    step.value = 3
  } catch (e) {
    error.value = (e as Error).message
    step.value = 2
  } finally {
    loading.value = false
  }
}

const goLogin = () => router.push({ name: 'login' })
</script>

<template>
  <div class="flex min-h-screen flex-col bg-[var(--color-background)]">
    <!-- Header -->
    <header class="border-b border-[var(--color-border)]">
      <div class="mx-auto flex h-16 max-w-2xl items-center justify-between px-6">
        <div class="flex items-center gap-2.5">
          <span class="flex h-8 w-8 items-center justify-center rounded-[var(--radius-md)] bg-[var(--color-primary)]">
            <Boxes class="h-4 w-4 text-[var(--color-on-primary)]" />
          </span>
          <span class="text-sm font-semibold tracking-tight text-[var(--color-text)]">Bloberry</span>
        </div>
        <p class="text-xs text-[var(--color-text-muted)]">One-time setup</p>
      </div>
    </header>

    <main class="flex flex-1 items-start justify-center px-6 py-10">
      <div class="w-full max-w-lg">
        <!-- Stepper -->
        <ol class="mb-10 flex items-center gap-2">
          <li v-for="(label, i) in steps" :key="label" class="flex flex-1 flex-col gap-1.5">
            <span class="flex items-center gap-2">
              <span
                class="flex h-7 w-7 items-center justify-center rounded-full text-xs font-semibold"
                :class="i + 1 < step ? 'bg-[var(--color-success)] text-white' : i + 1 === step ? 'bg-[var(--color-primary)] text-[var(--color-on-primary)]' : 'bg-[var(--color-surface)] text-[var(--color-text-muted)]'"
              >
                <CheckCircle2 v-if="i + 1 < step" class="h-4 w-4" />
                <span v-else>{{ i + 1 }}</span>
              </span>
              <span class="hidden text-xs font-medium text-[var(--color-text-muted)] sm:inline">{{ label }}</span>
            </span>
            <span v-if="i < steps.length - 1" class="h-px bg-[var(--color-border)]" :class="i + 1 < step ? '!bg-[var(--color-success)]' : ''" />
          </li>
        </ol>

        <!-- Step 1: admin account -->
        <div v-if="step === 1" class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-8 shadow-[var(--shadow-sm)]">
          <h1 class="text-xl font-bold tracking-tight text-[var(--color-text)]">Create your admin account</h1>
          <p class="mt-1.5 text-sm text-[var(--color-text-muted)]">This account becomes the platform administrator.</p>
          <div class="mt-6 flex flex-col gap-4">
            <AppInput v-model="email" label="Email" type="email" placeholder="admin@company.com" :icon="Mail" autocomplete="username" />
            <AppInput v-model="password" label="Password" type="password" placeholder="At least 8 characters" :icon="Lock" :error="step === 1 ? error : ''" autocomplete="new-password" />
            <AppInput v-model="displayName" label="Display name" placeholder="Your name" :icon="User" />
          </div>
          <div class="mt-8 flex justify-end">
            <AppButton @click="next">
              Continue
              <ArrowRight class="ml-2 h-4 w-4" />
            </AppButton>
          </div>
        </div>

        <!-- Step 2: tenant -->
        <div v-else-if="step === 2" class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-8 shadow-[var(--shadow-sm)]">
          <h1 class="text-xl font-bold tracking-tight text-[var(--color-text)]">Name your first tenant</h1>
          <p class="mt-1.5 text-sm text-[var(--color-text-muted)]">A tenant is your first workspace. You can add more later.</p>
          <div class="mt-6 flex flex-col gap-4">
            <AppInput v-model="tenantName" label="Tenant name" placeholder="Acme Corp" :icon="Building2" :error="step === 2 ? error : ''" />
            <AppInput v-model="tenantSlug" label="Tenant slug" placeholder="acme" :icon="AtSign" hint="Used in URLs and short links. Leave blank to derive from the name." />
          </div>
          <div class="mt-8 flex items-center justify-between">
            <AppButton variant="ghost" @click="back">
              <ArrowLeft class="mr-2 h-4 w-4" /> Back
            </AppButton>
            <AppButton :loading="loading" @click="submit">
              Create admin &amp; tenant
              <ArrowRight class="ml-2 h-4 w-4" />
            </AppButton>
          </div>
        </div>

        <!-- Step 3: done -->
        <div v-else class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-10 text-center shadow-[var(--shadow-sm)]">
          <span class="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-[var(--color-success)]/15">
            <CheckCircle2 class="h-7 w-7 text-[var(--color-success)]" />
          </span>
          <h1 class="mt-5 text-2xl font-bold tracking-tight text-[var(--color-text)]">You're ready</h1>
          <p class="mx-auto mt-2 max-w-sm text-sm leading-relaxed text-[var(--color-text-muted)]">
            Platform admin, tenant and a local disk backend are created. Sign in to upload your first file.
          </p>
          <div class="mt-6 flex items-center justify-center gap-2 rounded-[var(--radius-md)] bg-[var(--color-surface)] px-4 py-3 text-xs text-[var(--color-text-muted)]">
            <HardDrive class="h-4 w-4 text-[var(--color-primary)]" />
            Local disk backend · install-level
          </div>
          <AppButton class="mt-8" @click="goLogin">
            Continue to sign in
            <ArrowRight class="ml-2 h-4 w-4" />
          </AppButton>
        </div>
      </div>
    </main>
  </div>
</template>
