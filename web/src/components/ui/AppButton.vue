<script setup lang="ts">
defineProps<{
  variant?: 'primary' | 'secondary' | 'destructive' | 'ghost'
  type?: 'button' | 'submit'
  disabled?: boolean
  loading?: boolean
  size?: 'md' | 'sm'
}>()
</script>

<template>
  <button
    :type="type || 'button'"
    :disabled="disabled || loading"
    :class="[
      'inline-flex items-center justify-center gap-2 font-semibold rounded-[var(--radius-md)]',
      'focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] focus:ring-offset-2',
      size === 'sm' ? 'h-9 px-3 text-sm' : 'h-12 px-4 text-sm',
      variant === 'primary' && 'bg-[var(--color-primary)] text-[var(--color-on-primary)]',
      variant === 'secondary' && 'border border-[var(--color-primary)] text-[var(--color-primary)] bg-transparent',
      variant === 'destructive' && 'bg-[var(--color-error)] text-[var(--color-on-primary)]',
      variant === 'ghost' && 'text-[var(--color-text-muted)] hover:bg-[var(--color-surface)]',
      (disabled || loading) && 'bg-[var(--color-disabled)] text-[var(--color-on-disabled)] cursor-not-allowed',
    ]"
  >
    <span v-if="loading" class="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
    <slot v-else />
  </button>
</template>
