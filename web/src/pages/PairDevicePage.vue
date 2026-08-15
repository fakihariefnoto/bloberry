<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import QRCode from 'qrcode'
import { Smartphone, Download, RefreshCw, AlertTriangle, Copy, Check } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppModal from '../components/ui/AppModal.vue'
import PageHeader from '../components/ui/PageHeader.vue'

const qrPayload = ref('')
const qrError = ref('')
const qrLoading = ref(false)
const copied = ref(false)

async function issuePair() {
  qrError.value = ''
  qrLoading.value = true
  try {
    const res = await api.post<{ qr_payload: string }>('/auth/pair/issue')
    qrPayload.value = res.qr_payload
  } catch (e) { qrError.value = (e as Error).message } finally { qrLoading.value = false }
}

watch(qrPayload, async (payload) => {
  if (!payload) return
  await nextTick()
  const el = document.getElementById('pair-qr')
  if (el) QRCode.toCanvas(el, payload, { width: 224, margin: 1 })
})

// ---- desktop config export ----
const showExport = ref(false)
const exportPass = ref('')
const exportPass2 = ref('')
const exportError = ref('')
const exportLoading = ref(false)
const exported = ref(false)

function openExport() {
  exportPass.value = ''
  exportPass2.value = ''
  exportError.value = ''
  exported.value = false
  showExport.value = true
}

async function doExport() {
  exportError.value = ''
  if (exportPass.value.length < 8) {
    exportError.value = 'Passphrase must be at least 8 characters.'
    return
  }
  if (exportPass.value !== exportPass2.value) {
    exportError.value = 'Passphrases do not match.'
    return
  }
  exportLoading.value = true
  try {
    const res = await api.post<{ payload: string }>('/auth/config/issue')
    const { encrypted, salt, iv } = await encryptPayload(res.payload, exportPass.value)
    const file = JSON.stringify({ v: 1, salt, iv, data: encrypted }, null, 2)
    const blob = new Blob([file], { type: 'application/json' })
    const a = document.createElement('a')
    a.href = URL.createObjectURL(blob)
    a.download = 'bloberry.bloberry'
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(a.href)
    exported.value = true
  } catch (e) { exportError.value = (e as Error).message } finally { exportLoading.value = false }
}

async function encryptPayload(plaintext: string, passphrase: string) {
  const enc = new TextEncoder()
  const salt = crypto.getRandomValues(new Uint8Array(16))
  const iv = crypto.getRandomValues(new Uint8Array(12))
  const keyMaterial = await crypto.subtle.importKey('raw', enc.encode(passphrase), 'PBKDF2', false, ['deriveKey'])
  const key = await crypto.subtle.deriveKey(
    { name: 'PBKDF2', salt, iterations: 600000, hash: 'SHA-256' },
    keyMaterial,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt'],
  )
  const ciphertext = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, key, enc.encode(plaintext))
  return {
    encrypted: btoa(String.fromCharCode(...new Uint8Array(ciphertext))),
    salt: btoa(String.fromCharCode(...salt)),
    iv: btoa(String.fromCharCode(...iv)),
  }
}

function copyPayload() {
  navigator.clipboard?.writeText(qrPayload.value)
  copied.value = true
  setTimeout(() => (copied.value = false), 1500)
}
</script>

<template>
  <div>
    <PageHeader title="Pair a device" description="Sign in on your phone or the desktop app without re-entering your password." />

    <!-- QR pairing -->
    <div class="max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <div class="flex items-center gap-2">
        <Smartphone class="h-4 w-4 text-[var(--color-primary)]" />
        <h2 class="text-lg font-semibold text-[var(--color-text)]">Mobile app</h2>
      </div>
      <p class="mt-2 text-sm text-[var(--color-text-muted)]">
        Scan this code in the Bloberry mobile app. <strong class="text-[var(--color-warning)]">It signs you in — it expires in 2 minutes.</strong>
      </p>
      <div class="mt-4 flex flex-col items-center gap-4">
        <div class="flex h-60 w-60 items-center justify-center rounded-[var(--radius-md)] border border-dashed border-[var(--color-border)] bg-[var(--color-surface)]">
          <canvas v-if="qrPayload" id="pair-qr" class="h-56 w-56" />
          <p v-else class="px-6 text-center text-sm text-[var(--color-text-muted)]">Press "Show code" to generate a QR</p>
        </div>
        <p v-if="qrError" class="text-sm text-[var(--color-error)]">{{ qrError }}</p>
        <div class="flex items-center gap-2">
          <AppButton :loading="qrLoading" @click="issuePair">
            <RefreshCw v-if="qrPayload" class="mr-2 h-4 w-4" />
            {{ qrPayload ? 'Refresh code' : 'Show code' }}
          </AppButton>
          <button v-if="qrPayload" class="text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="copyPayload">
            <Check v-if="copied" class="h-4 w-4 text-[var(--color-success)]" />
            <Copy v-else class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Desktop config export -->
    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <div class="flex items-center gap-2">
        <Download class="h-4 w-4 text-[var(--color-primary)]" />
        <h2 class="text-lg font-semibold text-[var(--color-text)]">Desktop login file</h2>
      </div>
      <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">
        Export a passphrase-encrypted login file for the desktop app. The passphrase is used only in this browser to encrypt — it never leaves your machine.
      </p>
      <div class="mt-4 flex justify-end"><AppButton variant="secondary" @click="openExport">Export .bloberry file</AppButton></div>
    </div>

    <!-- Export modal -->
    <AppModal :open="showExport" title="Export desktop login file" description="Choose a passphrase to encrypt the file. You'll enter it again in the desktop app." @close="showExport = false">
      <div class="flex flex-col gap-4">
        <AppInput v-model="exportPass" label="Passphrase" type="password" placeholder="At least 8 characters" autocomplete="new-password" />
        <AppInput v-model="exportPass2" label="Confirm passphrase" type="password" placeholder="Repeat passphrase" :error="exportError" autocomplete="new-password" />
        <div class="flex items-start gap-2 rounded-[var(--radius-sm)] bg-[var(--color-surface)] p-3">
          <AlertTriangle class="mt-0.5 h-4 w-4 shrink-0 text-[var(--color-warning)]" />
          <p class="text-xs leading-relaxed text-[var(--color-text-muted)]">
            The file is valid for 24 hours and signs you in to this account. Losing the passphrase means the file is unusable.
          </p>
        </div>
        <p v-if="exported" class="text-xs text-[var(--color-success)]">File downloaded. Import it in the desktop app.</p>
        <div class="flex justify-end gap-2">
          <AppButton variant="ghost" @click="showExport = false">Cancel</AppButton>
          <AppButton :loading="exportLoading" @click="doExport"><Download class="mr-2 h-4 w-4" /> Export &amp; download</AppButton>
        </div>
      </div>
    </AppModal>
  </div>
</template>
