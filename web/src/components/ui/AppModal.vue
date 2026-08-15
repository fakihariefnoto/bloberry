<script setup lang="ts">
import { X } from 'lucide-vue-next'

defineProps<{
  open: boolean
  title: string
  description?: string
  width?: string
}>()

const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="emit('close')">
      <div
        class="flex max-h-[90vh] w-full flex-col overflow-hidden rounded-[var(--radius-lg)] bg-[var(--color-surface-raised)] shadow-[var(--shadow-lg)]"
        :style="{ maxWidth: width || '28rem' }"
      >
        <div class="flex items-start justify-between gap-4 border-b border-[var(--color-border)] px-6 py-4">
          <div>
            <h2 class="text-lg font-semibold text-[var(--color-text)]">{{ title }}</h2>
            <p v-if="description" class="mt-0.5 text-sm text-[var(--color-text-muted)]">{{ description }}</p>
          </div>
          <button class="rounded p-1 text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-text)]" @click="emit('close')">
            <X class="h-5 w-5" />
          </button>
        </div>
        <div class="flex-1 overflow-y-auto px-6 py-5">
          <slot />
        </div>
        <div v-if="$slots.footer" class="flex items-center justify-end gap-2 border-t border-[var(--color-border)] px-6 py-4">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </Teleport>
</template>
