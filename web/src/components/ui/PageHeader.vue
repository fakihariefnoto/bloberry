<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div class="max-w-xl">
      <h1 class="text-2xl font-bold tracking-tight text-[var(--color-text)]">{{ title }}</h1>
      <p v-if="description" class="mt-1.5 text-sm leading-relaxed text-[var(--color-text-muted)]">{{ description }}</p>
    </div>
    <div class="flex items-center gap-2">
      <slot />
    </div>
  </div>
</template>
