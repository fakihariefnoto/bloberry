<script setup lang="ts">
import { Link2Off } from 'lucide-vue-next'
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-[var(--color-background)] px-4 text-center">
    <Link2Off class="h-12 w-12 text-[var(--color-text-muted)]" />
    <h1 class="mt-4 text-lg font-semibold text-[var(--color-text)]">This link has expired</h1>
    <p class="mt-1 max-w-sm text-sm text-[var(--color-text-muted)]">
      This share link is no longer valid. Ask the person who shared it for a fresh one.
    </p>
  </div>
</template>
