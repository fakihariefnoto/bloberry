<script setup lang="ts">
defineProps<{
  variant?: 'primary' | 'secondary' | 'destructive' | 'ghost' | 'gradient'
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
      variant === 'gradient' && 'text-[#012624] font-medium uppercase tracking-[0.08em] text-sm rounded-md',
      (disabled || loading) && 'bg-[var(--color-disabled)] text-[var(--color-on-disabled)] cursor-not-allowed',
    ]"
    :style="variant === 'gradient' ? { background: 'linear-gradient(90deg, rgb(0,130,124) 0%, rgb(203,255,252) 60%, rgb(250,209,255) 100%)' } : undefined"
  >
    <span v-if="loading" class="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
    <slot v-else />
  </button>
</template>
