<script setup lang="ts">
import { PauseCircle } from 'lucide-vue-next'
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-[var(--color-background)] px-4 text-center">
    <PauseCircle class="h-12 w-12 text-[var(--color-warning)]" />
    <h1 class="mt-4 text-lg font-semibold text-[var(--color-text)]">This project is suspended</h1>
    <p class="mt-1 max-w-sm text-sm text-[var(--color-text-muted)]">
      Read access still works. Contact the platform administrator to resume writes.
    </p>
  </div>
</template>
