<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../lib/api'

interface ShareLink { id: string; kind: string; url: string; expires_at?: string; hit_count: number }

const links = ref<ShareLink[]>([])
const loading = ref(false)
const error = ref('')

async function load() {
  loading.value = true
  try {
    links.value = await api.get<ShareLink[]>('/shares')
  } catch {
    links.value = []
  } finally {
    loading.value = false
  }
}
onMounted(load)
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-[var(--color-text)]">Shares</h1>
    <p class="mt-1 text-sm text-[var(--color-text-muted)]">Links you've created for your files.</p>

    <p v-if="error" class="mt-3 text-sm text-[var(--color-error)]">{{ error }}</p>

    <div class="mt-4 overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Link</th>
            <th class="px-3 py-2">Kind</th>
            <th class="px-3 py-2">Hits</th>
            <th class="px-3 py-2">Expires</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="l in links" :key="l.id" class="border-t border-[var(--color-border)]">
            <td class="px-3 py-2.5 font-mono text-xs">{{ l.url }}</td>
            <td class="px-3 py-2.5">{{ l.kind }}</td>
            <td class="px-3 py-2.5">{{ l.hit_count }}</td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ l.expires_at ? new Date(l.expires_at).toLocaleString() : 'never' }}</td>
          </tr>
          <tr v-if="!loading && !links.length">
            <td colspan="4" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No share links yet</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
