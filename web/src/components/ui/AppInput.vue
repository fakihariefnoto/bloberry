<script setup lang="ts">
defineProps<{
  label?: string
  error?: string
  hint?: string
  type?: string
  placeholder?: string
  modelValue?: string
}>()
defineEmits<{ 'update:modelValue': [string] }>()
</script>

<template>
  <div class="flex flex-col gap-1">
    <label v-if="label" class="text-xs font-medium text-[var(--color-text-muted)]">{{ label }}</label>
    <input
      :type="type || 'text'"
      :placeholder="placeholder"
      :value="modelValue"
      class="h-12 w-full rounded-[var(--radius-md)] border bg-[var(--color-surface)] px-4 text-sm outline-none transition-[border-color] duration-150"
      :class="error ? 'border-[var(--color-error)]' : 'border-[var(--color-border)] focus:border-[var(--color-primary)] focus:border-2'"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <p v-if="error" class="text-xs text-[var(--color-error)]">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-[var(--color-text-muted)]">{{ hint }}</p>
  </div>
</template>
