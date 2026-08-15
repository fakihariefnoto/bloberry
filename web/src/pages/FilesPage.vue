<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Folder, FileText, ChevronRight, Upload, RefreshCw, MoreHorizontal, Globe, HardDrive } from 'lucide-vue-next'
import { api } from '../lib/api'
import { useTenantStore } from '../stores/tenant'
import AppButton from '../components/ui/AppButton.vue'

interface FolderRec { id: string; name: string; path: string }
interface ObjectRec { id: string; name: string; size_bytes: number; visibility: string; content_type: string; state: string; backend_name?: string; backend_driver?: string }

const route = useRoute()
const router = useRouter()
const tenants = useTenantStore()

const folderId = computed(() => (route.params.folderId as string) || '')
const folders = ref<FolderRec[]>([])
const objects = ref<ObjectRec[]>([])
const loading = ref(false)
const error = ref('')
const dragOver = ref(false)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get<{ folders: FolderRec[]; objects: ObjectRec[] }>(`/folders/${folderId.value || 'root'}/children`)
    folders.value = res.folders || []
    objects.value = res.objects || []
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
  loadBackends()
})
watch(folderId, load)

function openFolder(id: string) {
  router.push({ name: 'files', params: { folderId: id } })
}

function openFile(id: string) {
  router.push({ name: 'file-detail', params: { fileId: id } })
}

function formatBytes(n: number) {
  if (!n) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (n >= 1024 && i < units.length - 1) { n /= 1024; i++ }
  return `${n.toFixed(n >= 10 || i === 0 ? 0 : 1)} ${units[i]}`
}

// Upload queue — presigned PUT for each selected file.
const uploading = ref(false)
const uploads = ref<{ name: string; progress: number; state: string; error?: string }[]>([])
const fileInput = ref<HTMLInputElement | null>(null)

// Backend switcher — a tenant can write to any install-level or own backend.
interface BackendRec { id: string; name: string; driver: string; tenant_id?: string }
const backends = ref<BackendRec[]>([])
const selectedBackend = ref('')
const currentBackendName = ref('')

async function loadBackends() {
  try {
    backends.value = await api.get<BackendRec[]>('/backends')
    const saved = localStorage.getItem(`bloberry.backend.${tenants.currentId || ''}`) || ''
    if (saved && backends.value.some((b) => b.id === saved)) {
      selectedBackend.value = saved
    } else {
      selectedBackend.value = ''
    }
    updateBackendName()
  } catch { backends.value = [] }
}

function switchBackend() {
  localStorage.setItem(`bloberry.backend.${tenants.currentId || ''}`, selectedBackend.value)
  updateBackendName()
}

function updateBackendName() {
  const b = backends.value.find((x) => x.id === selectedBackend.value)
  currentBackendName.value = b ? `${b.name} (${b.driver})` : ''
}

function openPicker() {
  fileInput.value?.click()
}

async function createFolder() {
  const name = window.prompt('New folder name')
  if (!name || !name.trim()) return
  try {
    await api.post('/folders', { name: name.trim(), parent_id: folderId.value || 'root' })
    load()
  } catch (e) {
    alert((e as Error).message)
  }
}

async function onFileInput(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files) return
  await uploadFiles(Array.from(files))
  ;(e.target as HTMLInputElement).value = ''
}

async function onDrop(e: DragEvent) {
  dragOver.value = false
  if (e.dataTransfer?.files.length) await uploadFiles(Array.from(e.dataTransfer.files))
}

async function uploadFiles(files: File[]) {
  uploading.value = true
  for (const f of files) {
    const item: { name: string; progress: number; state: string; error?: string } = { name: f.name, progress: 0, state: 'pending' }
    uploads.value.push(item)
    try {
      const res = await api.post<{ file_id: string; upload_url: string; headers?: Record<string, string> }>('/objects/presign-put', {
        folder_id: folderId.value || 'root',
        name: f.name,
        backend_id: selectedBackend.value || undefined,
        size: f.size,
        content_type: f.type,
      })
      item.state = 'uploading'
      const xhr = new XMLHttpRequest()
      xhr.open('PUT', res.upload_url)
      if (res.headers) Object.entries(res.headers).forEach(([k, v]) => xhr.setRequestHeader(k, v))
      xhr.upload.onprogress = (ev) => {
        if (ev.lengthComputable) item.progress = Math.round((ev.loaded / ev.total) * 100)
      }
      await new Promise<void>((resolve, reject) => {
        xhr.onload = () => (xhr.status >= 200 && xhr.status < 300 ? resolve() : reject(new Error(`Upload failed (${xhr.status})`)))
        xhr.onerror = () => reject(new Error('Network error'))
        xhr.send(f)
      })
      await api.post(`/objects/${res.file_id}/complete`, { etag: 'x' })
      item.state = 'done'
      item.progress = 100
    } catch (e) {
      item.state = 'failed'
      item.error = (e as Error).message
    }
  }
  uploading.value = false
  load()
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div class="flex items-center gap-1 text-sm">
        <button class="text-[var(--color-text-muted)] hover:text-[var(--color-primary)]" @click="router.push({ name: 'files' })">Root</button>
        <ChevronRight v-if="folderId" class="h-4 w-4 text-[var(--color-text-muted)]" />
        <span v-if="folderId" class="font-semibold text-[var(--color-text)]">…</span>
      </div>
      <div class="flex items-center gap-2">
        <select
          v-if="backends.length"
          v-model="selectedBackend"
          class="h-9 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)] px-2.5 text-xs font-medium text-[var(--color-text-muted)] outline-none focus:border-[var(--color-primary)]"
          @change="switchBackend"
        >
          <option value="">Default backend</option>
          <option v-for="b in backends" :key="b.id" :value="b.id">{{ b.name }} ({{ b.driver }})</option>
        </select>
        <AppButton variant="secondary" size="sm" @click="load"><RefreshCw class="h-4 w-4" /> Refresh</AppButton>
        <AppButton variant="secondary" size="sm" @click="createFolder"><Folder class="h-4 w-4" /> New folder</AppButton>
        <input ref="fileInput" type="file" multiple class="hidden" @change="onFileInput" />
        <AppButton size="sm" :disabled="uploading" @click="openPicker"><Upload class="h-4 w-4" /> Upload</AppButton>
      </div>
    </div>

    <div
      class="flex min-h-[200px] flex-col rounded-[var(--radius-md)] border-2 border-dashed p-8 text-center transition-colors duration-150"
      :class="dragOver ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)]' : 'border-[var(--color-border)] bg-[var(--color-background)]'"
      @dragover.prevent="dragOver = true"
      @dragleave.prevent="dragOver = false"
      @drop.prevent="onDrop"
    >
      <Upload class="mx-auto h-6 w-6 text-[var(--color-text-muted)]" />
      <p class="mt-2 text-sm text-[var(--color-text)]">Drag files here or press Upload</p>
      <p class="mt-1 text-xs text-[var(--color-text-muted)]">
        Max 5 GB per file<span v-if="currentBackendName"> · uploading to {{ currentBackendName }}</span>
      </p>
    </div>

    <p v-if="error" class="mt-3 text-sm text-[var(--color-error)]">{{ error }}</p>

    <div class="mt-4 overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="px-3 py-2">Name</th>
            <th class="px-3 py-2">Size</th>
            <th class="px-3 py-2">Storage</th>
            <th class="px-3 py-2">Status</th>
            <th class="px-3 py-2"></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="f in folders"
            :key="f.id"
            class="cursor-pointer border-t border-[var(--color-border)] transition-colors hover:bg-[var(--color-surface)]"
            @click="openFolder(f.id)"
          >
            <td class="px-3 py-2.5">
              <span class="flex items-center gap-2"><Folder class="h-4 w-4 text-[var(--color-primary)]" /> {{ f.name }}</span>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">—</td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">—</td>
            <td class="px-3 py-2.5"></td>
            <td class="px-3 py-2.5"><MoreHorizontal class="h-4 w-4 text-[var(--color-text-muted)]" /></td>
          </tr>
          <tr
            v-for="o in objects"
            :key="o.id"
            class="cursor-pointer border-t border-[var(--color-border)] transition-colors hover:bg-[var(--color-surface)]"
            @click="openFile(o.id)"
          >
            <td class="px-3 py-2.5">
              <span class="flex items-center gap-2">
                <FileText class="h-4 w-4 text-[var(--color-text-muted)]" /> {{ o.name }}
              </span>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">{{ formatBytes(o.size_bytes) }}</td>
            <td class="px-3 py-2.5">
              <span class="inline-flex items-center gap-1.5 text-xs text-[var(--color-text-muted)]">
                <HardDrive class="h-3.5 w-3.5" />
                <span class="font-mono">{{ o.backend_driver || '—' }}</span>
                <span v-if="o.backend_name" class="text-[var(--color-text-muted)]/70">· {{ o.backend_name }}</span>
              </span>
            </td>
            <td class="px-3 py-2.5">
              <span
                v-if="o.visibility === 'public'"
                class="inline-flex items-center gap-1 rounded-full bg-[var(--color-warning)]/15 px-2 py-0.5 text-xs font-medium text-[var(--color-warning)]"
              >
                <Globe class="h-3 w-3" /> Public
              </span>
            </td>
            <td class="px-3 py-2.5"></td>
          </tr>
          <tr v-if="!loading && !folders.length && !objects.length">
            <td colspan="5" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No files here yet</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Upload queue -->
    <div v-if="uploads.length" class="fixed bottom-4 right-4 w-80 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4 shadow-[var(--shadow-md)]">
      <p class="text-sm font-semibold text-[var(--color-text)]">
        {{ uploads.filter((u) => u.state === 'pending' || u.state === 'uploading').length }} uploading · {{ uploads.filter((u) => u.state === 'done').length }} done
      </p>
      <div class="mt-3 flex flex-col gap-2">
        <div v-for="u in uploads" :key="u.name" class="flex items-center justify-between gap-2 text-xs">
          <span class="truncate text-[var(--color-text)]">{{ u.name }}</span>
          <span
            :class="u.state === 'failed' ? 'text-[var(--color-error)]' : u.state === 'done' ? 'text-[var(--color-success)]' : 'text-[var(--color-text-muted)]'"
          >
            {{ u.state === 'done' ? 'Done' : u.state === 'failed' ? (u.error || 'Failed') : `${u.progress}%` }}
          </span>
        </div>
        <div v-if="!uploading && uploads.some((u) => u.state === 'done')" class="mt-1 h-1 overflow-hidden rounded-full bg-[var(--color-border)]">
          <div class="h-full bg-[var(--color-primary)]" :style="{ width: '100%' }" />
        </div>
      </div>
    </div>
  </div>
</template>
