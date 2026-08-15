<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  state: 'active' | 'pending' | 'failed' | 'public' | 'private' | 'neutral'
  label?: string
}>()

const cls = computed(() => {
  const map: Record<string, string> = {
    active: 'bg-active/15 text-active',
    pending: 'bg-pending/15 text-pending',
    failed: 'bg-failed/15 text-failed',
    public: 'bg-public/15 text-public',
    private: 'bg-[var(--color-surface)] text-private',
    neutral: 'bg-[var(--color-surface)] text-[var(--color-text-muted)]',
  }
  return map[props.state]
})
</script>

<template>
  <span
    class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
    :class="cls"
  >
    <slot>{{ label }}</slot>
  </span>
</template>
