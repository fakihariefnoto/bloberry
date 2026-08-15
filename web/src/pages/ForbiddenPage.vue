<script setup lang="ts">
import { ShieldX } from 'lucide-vue-next'
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-[var(--color-background)] px-4 text-center">
    <ShieldX class="h-12 w-12 text-[var(--color-error)]" />
    <h1 class="mt-4 text-lg font-semibold text-[var(--color-text)]">You don't have access</h1>
    <p class="mt-1 max-w-sm text-sm text-[var(--color-text-muted)]">
      This area requires a higher role. If you think this is a mistake, ask a tenant admin.
    </p>
    <RouterLink to="/app" class="mt-4 text-sm text-[var(--color-primary)] hover:underline">Back to Files</RouterLink>
  </div>
</template>
