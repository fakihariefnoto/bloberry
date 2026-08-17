<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Folder, FileText, ChevronRight, Upload, RefreshCw, Globe, HardDrive, Trash2 } from 'lucide-vue-next'
import { api } from '../lib/api'
import { useTenantStore } from '../stores/tenant'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppModal from '../components/ui/AppModal.vue'

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

async function load() {
  loading.value = true
  error.value = ''
  try {
    const q = selectedBackend.value ? `?storage_id=${encodeURIComponent(selectedBackend.value)}` : ''
    const res = await api.get<{ folders: FolderRec[]; objects: ObjectRec[] }>(`/folders/${folderId.value || 'root'}/children${q}`)
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
    let chosen = ''
    if (saved && backends.value.some((b) => b.id === saved)) {
      chosen = saved
    } else if (tenants.current?.default_storage_id && backends.value.some((b) => b.id === tenants.current!.default_storage_id)) {
      chosen = tenants.current.default_storage_id
    } else if (backends.value.length) {
      chosen = backends.value[0].id // default to the tenant's first engine
    }
    selectedBackend.value = chosen
    updateBackendName()
    await load()
  } catch { backends.value = [] }
}

function switchBackend() {
  localStorage.setItem(`bloberry.backend.${tenants.currentId || ''}`, selectedBackend.value)
  updateBackendName()
  load() // refresh the listing so files from the newly-selected engine show
}

function updateBackendName() {
  const b = backends.value.find((x) => x.id === selectedBackend.value)
  currentBackendName.value = b ? `${b.name} (${b.driver})` : ''
}

// Folder creation modal state
const showCreateFolder = ref(false)
const newFolderName = ref('')
const folderError = ref('')
const folderCreating = ref(false)

function openCreateFolder() {
  newFolderName.value = ''
  folderError.value = ''
  showCreateFolder.value = true
}

function closeCreateFolder() {
  if (folderCreating.value) return
  showCreateFolder.value = false
  newFolderName.value = ''
  folderError.value = ''
}

async function createFolder() {
  const name = newFolderName.value.trim()
  if (!name) {
    folderError.value = 'Folder name is required.'
    return
  }
  folderError.value = ''
  folderCreating.value = true
  try {
    await api.post('/folders', { name, parent_id: folderId.value || 'root' })
    showCreateFolder.value = false
    newFolderName.value = ''
    folderError.value = ''
    load()
  } catch (e) {
    folderError.value = (e as Error).message
  } finally {
    folderCreating.value = false
  }
}

async function onFileInput(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files) return
  await uploadFiles(Array.from(files))
  ;(e.target as HTMLInputElement).value = ''
}

// name-conflict resolution
const showConflict = ref(false)
const conflictFile = ref<File | null>(null)

async function uploadFiles(files: File[], overwrite = false) {
  uploading.value = true
  for (const f of files) {
    const item: { name: string; progress: number; state: string; error?: string } = { name: f.name, progress: 0, state: 'pending' }
    uploads.value.push(item)
    try {
      const res = await api.post<{ file_id: string; upload_url: string; headers?: Record<string, string> }>('/objects/presign-put', {
        folder_id: folderId.value || 'root',
        name: f.name,
        storage_id: selectedBackend.value || undefined,
        overwrite,
        size: f.size,
        content_type: f.type,
      })
      item.state = 'uploading'
      ;(item as { _url?: string })._url = res.upload_url
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
      const err = e as Error & { code?: string }
      if (err.code === 'name_conflict') {
        item.state = 'failed'
        item.error = 'A file with this name already exists.'
        conflictFile.value = f
        showConflict.value = true
        break // pause the queue until the user decides
      }
      item.state = 'failed'
      const presign = (item as { _url?: string })._url
      item.error = presign ? `${err.message} — presigned URL: ${presign}` : err.message
    }
  }
  uploading.value = false
  load()
}

// Upload modal state — user picks drag-and-drop OR a file picker.
const showUpload = ref(false)
const uploadDragOver = ref(false)

function openUpload() {
  showUpload.value = true
  uploadDragOver.value = false
}

function onUploadDrop(e: DragEvent) {
  uploadDragOver.value = false
  if (e.dataTransfer?.files.length) {
    showUpload.value = false
    uploadFiles(Array.from(e.dataTransfer.files))
  }
}

async function pickAndUpload() {
  showUpload.value = false
  fileInput.value?.click()
}

// Bulk selection
const selected = ref<Set<string>>(new Set())
const showBulkDelete = ref(false)

function toggleSelect(id: string) {
  const s = new Set(selected.value)
  if (s.has(id)) s.delete(id)
  else s.add(id)
  selected.value = s
}

function toggleSelectAll(e: Event) {
  const checked = (e.target as HTMLInputElement).checked
  if (checked) {
    selected.value = new Set(objects.value.map((o) => o.id))
  } else {
    selected.value = new Set()
  }
}

async function confirmBulkDelete() {
  if (!selected.value.size) return
  showBulkDelete.value = true
}

async function doBulkDelete() {
  const ids = Array.from(selected.value)
  showBulkDelete.value = false
  for (const id of ids) {
    try { await api.delete(`/objects/${id}`) } catch { /* continue */ }
  }
  selected.value = new Set()
  load()
}

function skipConflict() {
  showConflict.value = false
  conflictFile.value = null
}

async function overwriteConflict() {
  const f = conflictFile.value
  showConflict.value = false
  conflictFile.value = null
  if (f) await uploadFiles([f], true)
}

// Folder deletion — subtree delete via DELETE /folders/{id}.
const showDeleteFolder = ref(false)
const folderToDelete = ref<FolderRec | null>(null)
const deletingFolder = ref(false)

function askDeleteFolder(f: FolderRec) {
  folderToDelete.value = f
  showDeleteFolder.value = true
}

async function doDeleteFolder() {
  const f = folderToDelete.value
  if (!f) return
  deletingFolder.value = true
  try {
    await api.delete(`/folders/${f.id}`)
    showDeleteFolder.value = false
    folderToDelete.value = null
    await load()
  } catch (e) {
    error.value = (e as Error).message
    showDeleteFolder.value = false
  } finally {
    deletingFolder.value = false
  }
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
          <option value="">Default storage engine</option>
          <option v-for="b in backends" :key="b.id" :value="b.id">{{ b.name }} ({{ b.driver }})</option>
        </select>
        <AppButton variant="secondary" size="sm" @click="load"><RefreshCw class="h-4 w-4" /> Refresh</AppButton>
        <AppButton variant="secondary" size="sm" @click="openCreateFolder"><Folder class="h-4 w-4" /> New folder</AppButton>
        <input ref="fileInput" type="file" multiple class="hidden" @change="onFileInput" />
        <AppButton size="sm" :disabled="uploading" @click="openUpload"><Upload class="h-4 w-4" /> Upload</AppButton>
      </div>
    </div>

    <p class="mb-3 text-xs text-[var(--color-text-muted)]">
      {{ currentBackendName ? `Uploading to ${currentBackendName}` : 'Select a storage engine to filter' }}
      <span v-if="objects.length"> · {{ objects.length }} files shown</span>
    </p>

    <p v-if="error" class="mt-3 text-sm text-[var(--color-error)]">{{ error }}</p>

    <!-- Bulk actions -->
    <div v-if="selected.size" class="mb-2 flex items-center gap-3 rounded-[var(--radius-md)] border border-[var(--color-primary)] bg-[var(--color-primary-subtle)] px-3 py-2">
      <span class="text-sm font-medium text-[var(--color-text)]">{{ selected.size }} selected</span>
      <button class="text-xs font-medium text-[var(--color-error)] hover:underline" @click="confirmBulkDelete">Delete selected</button>
      <button class="text-xs font-medium text-[var(--color-text-muted)] hover:underline" @click="selected = new Set()">Clear</button>
    </div>

    <div class="mt-4 overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)]">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-[var(--color-surface)] text-left text-xs font-medium text-[var(--color-text-muted)]">
            <th class="w-10 px-3 py-2">
              <input type="checkbox" class="h-4 w-4 accent-[var(--color-primary)]" :checked="selected.size === objects.length && objects.length > 0" @change="toggleSelectAll" />
            </th>
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
            <td class="px-3 py-2.5"></td>
            <td class="px-3 py-2.5">
              <span class="flex items-center gap-2"><Folder class="h-4 w-4 text-[var(--color-primary)]" /> {{ f.name }}</span>
            </td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">—</td>
            <td class="px-3 py-2.5 text-[var(--color-text-muted)]">—</td>
            <td class="px-3 py-2.5"></td>
            <td class="px-3 py-2.5" @click.stop>
              <button class="text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-error)]" :title="`Delete folder ${f.name}`" @click="askDeleteFolder(f)">
                <Trash2 class="h-4 w-4" />
              </button>
            </td>
          </tr>
          <tr
            v-for="o in objects"
            :key="o.id"
            class="cursor-pointer border-t border-[var(--color-border)] transition-colors hover:bg-[var(--color-surface)]"
            @click="openFile(o.id)"
          >
            <td class="px-3 py-2.5" @click.stop>
              <input type="checkbox" class="h-4 w-4 accent-[var(--color-primary)]" :checked="selected.has(o.id)" @change="toggleSelect(o.id)" />
            </td>
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
            <td colspan="6" class="px-4 py-10 text-center text-sm text-[var(--color-text-muted)]">No files here yet</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Bulk delete confirm -->
    <AppModal :open="showBulkDelete" title="Delete selected files" :description="`Permanently delete ${selected.size} file${selected.size > 1 ? 's' : ''}? This cannot be undone.`" @close="showBulkDelete = false">
      <template #footer>
        <AppButton variant="ghost" @click="showBulkDelete = false">Cancel</AppButton>
        <AppButton variant="destructive" @click="doBulkDelete"><Trash2 class="mr-2 h-4 w-4" /> Delete files</AppButton>
      </template>
    </AppModal>

    <!-- Folder delete confirm -->
    <AppModal :open="showDeleteFolder" title="Delete folder" :description="`Delete ${folderToDelete?.name}? Files inside it are also deleted. This cannot be undone.`" @close="showDeleteFolder = false">
      <template #footer>
        <AppButton variant="ghost" :disabled="deletingFolder" @click="showDeleteFolder = false">Cancel</AppButton>
        <AppButton variant="destructive" :loading="deletingFolder" @click="doDeleteFolder"><Trash2 class="mr-2 h-4 w-4" /> Delete folder</AppButton>
      </template>
    </AppModal>

    <!-- Name conflict modal -->
    <AppModal :open="showConflict" title="A file with this name already exists" description="Do you want to replace the existing file with this upload?">
      <p class="text-sm text-[var(--color-text-muted)]">
        <span class="font-mono text-[var(--color-text)]">{{ conflictFile?.name }}</span> already exists in this folder.
        Replacing it removes the current file. Skipping keeps both uploads in the queue untouched.
      </p>
      <template #footer>
        <AppButton variant="ghost" @click="skipConflict">Skip</AppButton>
        <AppButton variant="destructive" @click="overwriteConflict">Replace file</AppButton>
      </template>
    </AppModal>

    <!-- Upload modal: drag-and-drop OR file picker -->
    <AppModal :open="showUpload" title="Upload files" :description="currentBackendName ? `Uploading to ${currentBackendName}` : 'Choose a storage engine to upload to'" @close="showUpload = false">
      <div
        class="flex min-h-[220px] flex-col items-center justify-center rounded-[var(--radius-md)] border-2 border-dashed p-8 text-center transition-colors duration-150"
        :class="uploadDragOver ? 'border-[var(--color-primary)] bg-[var(--color-primary-subtle)]' : 'border-[var(--color-border)] bg-[var(--color-background)]'"
        @dragover.prevent="uploadDragOver = true"
        @dragleave.prevent="uploadDragOver = false"
        @drop.prevent="onUploadDrop"
      >
        <Upload class="h-8 w-8 text-[var(--color-text-muted)]" />
        <p class="mt-3 text-sm text-[var(--color-text)]">Drag files here</p>
        <p class="mt-1 text-xs text-[var(--color-text-muted)]">Max 5 GB per file</p>
      </div>
      <template #footer>
        <AppButton variant="ghost" @click="showUpload = false">Cancel</AppButton>
        <AppButton @click="pickAndUpload"><Upload class="mr-2 h-4 w-4" /> Choose files</AppButton>
      </template>
    </AppModal>

    <!-- New folder modal -->
    <AppModal :open="showCreateFolder" title="New folder" description="Create a folder in the current directory." @close="closeCreateFolder">
      <AppInput
        v-model="newFolderName"
        label="Folder name"
        placeholder="e.g. Project Files"
        :error="folderError"
        autofocus
        @keyup.enter="createFolder"
      />
      <template #footer>
        <AppButton variant="ghost" :disabled="folderCreating" @click="closeCreateFolder">Cancel</AppButton>
        <AppButton :loading="folderCreating" @click="createFolder">
          Create folder
        </AppButton>
      </template>
    </AppModal>

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
