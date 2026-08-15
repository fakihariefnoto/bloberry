<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import QRCode from 'qrcode'
import { Lock, ShieldCheck, Copy, Check } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import PageHeader from '../components/ui/PageHeader.vue'

// ---- password ----
const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const pwError = ref('')
const pwSuccess = ref('')
const pwLoading = ref(false)

async function changePassword() {
  pwError.value = ''
  pwSuccess.value = ''
  if (newPassword.value !== confirmPassword.value) {
    pwError.value = 'Passwords do not match'
    return
  }
  pwLoading.value = true
  try {
    await api.post('/users/me/password', { current_password: currentPassword.value, new_password: newPassword.value })
    pwSuccess.value = 'Password changed.'
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (e) { pwError.value = (e as Error).message } finally { pwLoading.value = false }
}

// ---- 2FA ----
const totpStep = ref<'idle' | 'provisioned' | 'enabling' | 'done'>('idle')
const totpSecret = ref('')
const otpauthUrl = ref('')
const backupCodes = ref<string[]>([])
const totpCode = ref('')
const totpError = ref('')
const totpLoading = ref(false)
const copied = ref(false)

async function provision() {
  totpError.value = ''
  totpLoading.value = true
  try {
    const res = await api.post<{ secret: string; otpauth_url: string }>('/auth/totp/provision')
    totpSecret.value = res.secret
    otpauthUrl.value = res.otpauth_url
    totpStep.value = 'provisioned'
  } catch (e) { totpError.value = (e as Error).message } finally { totpLoading.value = false }
}

async function enable() {
  totpError.value = ''
  totpLoading.value = true
  try {
    const res = await api.post<{ backup_codes: string[] }>('/auth/totp/enable', { code: totpCode.value })
    backupCodes.value = res.backup_codes || []
    totpStep.value = 'done'
  } catch (e) { totpError.value = (e as Error).message } finally { totpLoading.value = false }
}

// Render the TOTP QR into the SVG element when the otpauth URL is ready.
watch(otpauthUrl, async (url) => {
  if (!url) return
  await nextTick()
  const el = document.getElementById('totp-qr')
  if (el) QRCode.toCanvas(el, url, { width: 192, margin: 1 })
})

async function copyCodes() {
  await navigator.clipboard?.writeText(backupCodes.value.join('\n'))
  copied.value = true
  setTimeout(() => (copied.value = false), 1500)
}
</script>

<template>
  <div>
    <PageHeader title="Account settings" description="Password, security and two-factor authentication." />

    <!-- Password -->
    <div class="max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <div class="flex items-center gap-2">
        <Lock class="h-4 w-4 text-[var(--color-primary)]" />
        <h2 class="text-lg font-semibold text-[var(--color-text)]">Change password</h2>
      </div>
      <div class="mt-5 flex flex-col gap-4">
        <AppInput v-model="currentPassword" label="Current password" type="password" placeholder="Current password" autocomplete="current-password" />
        <AppInput v-model="newPassword" label="New password" type="password" placeholder="New password" autocomplete="new-password" />
        <AppInput v-model="confirmPassword" label="Confirm new password" type="password" placeholder="Confirm new password" :error="pwError" autocomplete="new-password" />
        <p v-if="pwSuccess" class="text-xs text-[var(--color-success)]">{{ pwSuccess }}</p>
        <div class="flex justify-end"><AppButton :loading="pwLoading" @click="changePassword">Change password</AppButton></div>
      </div>
    </div>

    <!-- 2FA -->
    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <div class="flex items-center gap-2">
        <ShieldCheck class="h-4 w-4 text-[var(--color-primary)]" />
        <h2 class="text-lg font-semibold text-[var(--color-text)]">Two-factor authentication</h2>
      </div>

      <!-- idle -->
      <div v-if="totpStep === 'idle'" class="mt-4 flex flex-col gap-3">
        <p class="text-sm leading-relaxed text-[var(--color-text-muted)]">
          Add an authenticator app (Google Authenticator, 1Password, Authy…) to require a one-time code on every login.
        </p>
        <p v-if="totpError" class="text-xs text-[var(--color-error)]">{{ totpError }}</p>
        <div class="flex justify-end"><AppButton :loading="totpLoading" @click="provision">Set up 2FA</AppButton></div>
      </div>

      <!-- provisioned: show secret + QR, enter code -->
      <div v-else-if="totpStep === 'provisioned'" class="mt-4 flex flex-col gap-4">
        <p class="text-sm text-[var(--color-text-muted)]">
          Scan this QR in your authenticator app, or enter the secret manually.
        </p>
        <div class="flex justify-center">
          <canvas id="totp-qr" class="h-48 w-48 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-white" />
        </div>
        <div class="flex items-center justify-center gap-2 rounded-[var(--radius-sm)] bg-[var(--color-surface)] p-2">
          <code class="font-mono text-xs text-[var(--color-text)]">{{ totpSecret }}</code>
        </div>
        <AppInput v-model="totpCode" label="Code from authenticator" placeholder="6-digit code" :error="totpError" />
        <div class="flex justify-end"><AppButton :loading="totpLoading" @click="enable">Verify &amp; enable</AppButton></div>
      </div>

      <!-- done: backup codes -->
      <div v-else-if="totpStep === 'done'" class="mt-4 flex flex-col gap-4">
        <p class="rounded-[var(--radius-sm)] border border-[var(--color-warning)] bg-[var(--color-warning)]/10 px-3 py-2 text-xs text-[var(--color-warning)]">
          Save these backup codes. Each is single-use — store them somewhere safe.
        </p>
        <div class="grid grid-cols-2 gap-2">
          <code v-for="c in backupCodes" :key="c" class="rounded-[var(--radius-sm)] bg-[var(--color-surface)] p-2 text-center font-mono text-xs text-[var(--color-text)]">{{ c }}</code>
        </div>
        <div class="flex justify-end">
          <AppButton variant="secondary" @click="copyCodes">
            <Check v-if="copied" class="mr-2 h-4 w-4 text-[var(--color-success)]" />
            <Copy v-else class="mr-2 h-4 w-4" />
            Copy codes
          </AppButton>
        </div>
        <p class="text-xs text-[var(--color-success)]">2FA is now enabled.</p>
      </div>
    </div>
  </div>
</template>
