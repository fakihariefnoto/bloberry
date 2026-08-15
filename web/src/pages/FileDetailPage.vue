<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Download, Globe, Lock } from 'lucide-vue-next'
import { api } from '../lib/api'
import AppButton from '../components/ui/AppButton.vue'

interface ObjectRec {
  id: string
  name: string
  size_bytes: number
  content_type: string
  visibility: string
  state: string
  folder_id: string
  created_at: string
  updated_at: string
}

const route = useRoute()
const router = useRouter()
const fileId = route.params.fileId as string

const obj = ref<ObjectRec | null>(null)
const error = ref('')
const shareUrl = ref('')

async function load() {
  try {
    obj.value = await api.get<ObjectRec>(`/objects/${fileId}`)
  } catch (e) {
    error.value = (e as Error).message
  }
}
onMounted(load)

function formatBytes(n: number) {
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  while (n >= 1024 && i < units.length - 1) { n /= 1024; i++ }
  return `${n.toFixed(n >= 10 || i === 0 ? 0 : 1)} ${units[i]}`
}

async function download() {
  const a = document.createElement('a')
  a.href = `/objects/${fileId}/download`
  a.download = obj.value?.name || 'download'
  document.body.appendChild(a)
  a.click()
  a.remove()
}

async function createShare() {
  try {
    const res = await api.post<{ url: string }>('/shares', { object_id: fileId, ttl: 3600 })
    shareUrl.value = res.url
  } catch (e) {
    error.value = (e as Error).message
  }
}

async function togglePublic() {
  try {
    obj.value = await api.patch<ObjectRec>(`/objects/${fileId}/visibility`, { visibility: obj.value?.visibility === 'public' ? 'private' : 'public' })
  } catch (e) {
    error.value = (e as Error).message
  }
}

const back = () => router.push({ name: 'files' })
function copyText(t: string) {
  window.navigator.clipboard?.writeText(t)
}
</script>

<template>
  <div>
    <button class="mb-4 flex items-center gap-1 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="back">
      <ArrowLeft class="h-4 w-4" /> Back to files
    </button>

    <p v-if="error" class="mb-3 text-sm text-[var(--color-error)]">{{ error }}</p>

    <div v-if="obj" class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-6">
      <div class="flex items-start justify-between">
        <div>
          <h1 class="text-xl font-bold text-[var(--color-text)]">{{ obj.name }}</h1>
          <p class="mt-1 text-xs text-[var(--color-text-muted)]">
            {{ formatBytes(obj.size_bytes) }} · {{ obj.content_type || 'application/octet-stream' }}
          </p>
        </div>
        <button
          class="flex items-center gap-1 rounded-full px-3 py-1 text-xs font-medium"
          :class="obj.visibility === 'public' ? 'bg-[var(--color-warning)]/15 text-[var(--color-warning)]' : 'bg-[var(--color-surface)] text-[var(--color-text-muted)]'"
          @click="togglePublic"
        >
          <Globe v-if="obj.visibility === 'public'" class="h-3.5 w-3.5" />
          <Lock v-else class="h-3.5 w-3.5" />
          {{ obj.visibility === 'public' ? 'Public' : 'Private' }}
        </button>
      </div>

      <div class="mt-6 flex gap-2">
        <AppButton @click="download"><Download class="h-4 w-4" /> Download</AppButton>
        <AppButton variant="secondary" @click="createShare">Create share link</AppButton>
      </div>

      <div v-if="shareUrl" class="mt-4 flex items-center gap-2">
        <code class="flex-1 truncate rounded-[var(--radius-sm)] bg-[var(--color-surface)] px-3 py-2 text-xs font-mono text-[var(--color-text)]">{{ shareUrl }}</code>
        <AppButton size="sm" variant="secondary" @click="copyText(shareUrl)">Copy</AppButton>
      </div>

      <dl class="mt-6 grid grid-cols-2 gap-4 border-t border-[var(--color-border)] pt-4 text-sm">
        <div><dt class="text-xs text-[var(--color-text-muted)]">File ID</dt><dd class="font-mono text-xs">{{ obj.id }}</dd></div>
        <div><dt class="text-xs text-[var(--color-text-muted)]">State</dt><dd>{{ obj.state }}</dd></div>
        <div><dt class="text-xs text-[var(--color-text-muted)]">Uploaded</dt><dd>{{ new Date(obj.created_at).toLocaleString() }}</dd></div>
        <div><dt class="text-xs text-[var(--color-text-muted)]">Modified</dt><dd>{{ new Date(obj.updated_at).toLocaleString() }}</dd></div>
      </dl>
    </div>
  </div>
</template>
