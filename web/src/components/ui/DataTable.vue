<script setup lang="ts">
// Generic data table per design/style-guide.md → Data table.
// Rows at size.row-height, sticky header, hover surface, selected primary-subtle.
defineProps<{
  columns: { key: string; label: string }[]
  rows: any[]
  loading?: boolean
  emptyTitle?: string
  emptyBody?: string
  onRowClick?: (row: any) => void
}>()
</script>

<template>
  <div class="overflow-hidden rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface)]">
    <div v-if="loading" class="flex flex-col gap-2 p-4">
      <div v-for="n in 5" :key="n" class="h-11 animate-pulse rounded bg-[var(--color-border)]" />
    </div>
    <table v-else class="w-full border-collapse text-sm">
      <thead>
        <tr>
          <th
            v-for="c in columns"
            :key="c.key"
            class="sticky top-0 border-b border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-2 text-left text-xs font-medium text-[var(--color-text-muted)]"
          >
            {{ c.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="row in rows"
          :key="row.id"
          class="cursor-pointer border-b border-[var(--color-border)] transition-colors duration-150 hover:bg-[var(--color-surface)]"
          @click="onRowClick?.(row)"
        >
          <td v-for="c in columns" :key="c.key" class="px-3 py-2.5 text-[var(--color-text)]">
            <slot :name="c.key" :row="row">
              {{ row[c.key] }}
            </slot>
          </td>
        </tr>
        <tr v-if="rows.length === 0">
          <td :colspan="columns.length" class="px-4 py-10 text-center">
            <p class="text-sm font-semibold text-[var(--color-text)]">{{ emptyTitle || 'Nothing here' }}</p>
            <p class="mt-1 text-xs text-[var(--color-text-muted)]">{{ emptyBody }}</p>
            <slot name="empty" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
