<script setup lang="ts">
import { FileQuestion } from 'lucide-vue-next'
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-[var(--color-background)] px-4 text-center">
    <FileQuestion class="h-12 w-12 text-[var(--color-text-muted)]" />
    <h1 class="mt-4 text-lg font-semibold text-[var(--color-text)]">Page not found</h1>
    <p class="mt-1 max-w-sm text-sm text-[var(--color-text-muted)]">
      That address doesn't exist on this install.
    </p>
    <RouterLink to="/app" class="mt-4 text-sm text-[var(--color-primary)] hover:underline">Back to Files</RouterLink>
  </div>
</template>
