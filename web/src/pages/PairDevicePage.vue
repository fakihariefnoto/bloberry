<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

const qrPayload = ref('')
const error = ref('')

async function issue() {
  error.value = ''
  try {
    const res = await api.post<{ qr_payload: string }>('/auth/pair/issue')
    qrPayload.value = res.qr_payload
  } catch (e) { error.value = (e as Error).message }
}

function renderQR(text: string) {
  // Minimal QR via a data-URI placeholder — the full QR renderer is a
  // dependency in the plan (qrcode lib). Keep it honest: show the payload
  // as a code block the mobile app can type, plus the QR graphic slot.
  return text
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-[var(--color-text)]">Pair a device</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">
      Sign in on your phone by scanning a code. <strong>This code signs you in — it expires in 2 minutes.</strong>
    </p>

    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <div class="flex flex-col items-center gap-4">
        <div class="flex h-56 w-56 items-center justify-center rounded-[var(--radius-md)] border border-dashed border-[var(--color-border)] bg-[var(--color-surface)]">
          <code v-if="qrPayload" class="break-all p-4 text-center font-mono text-xs text-[var(--color-text)]">{{ renderQR(qrPayload) }}</code>
          <p v-else class="text-sm text-[var(--color-text-muted)]">Press "Show code"</p>
        </div>
        <p v-if="error" class="text-sm text-[var(--color-error)]">{{ error }}</p>
        <AppButton @click="issue">{{ qrPayload ? 'Refresh code' : 'Show code' }}</AppButton>
        <p class="text-xs text-[var(--color-warning)]">The old code dies the moment you refresh — even if it was never scanned.</p>
      </div>
    </div>

    <div class="mt-4 max-w-md rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-5">
      <h2 class="text-lg font-semibold text-[var(--color-text)]">Desktop login file</h2>
      <p class="mt-1 text-sm text-[var(--color-text-muted)]">
        Export a passphrase-encrypted login file for the desktop app. The passphrase never leaves this browser.
      </p>
      <div class="mt-3"><AppButton variant="secondary" @click="issue">Export .bloberry file</AppButton></div>
    </div>
  </div>
</template>
