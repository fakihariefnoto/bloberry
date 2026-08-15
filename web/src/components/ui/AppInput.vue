<script setup lang="ts">
import { computed, ref } from 'vue'
import { Eye, EyeOff } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  label?: string
  error?: string
  hint?: string
  type?: string
  placeholder?: string
  modelValue?: string
  icon?: any // lucide component
  dark?: boolean
  autocomplete?: string
}>(), { type: 'text', dark: false })

const emit = defineEmits<{ 'update:modelValue': [string] }>()

const show = ref(false)
const isPassword = computed(() => props.type === 'password')
const inputType = computed(() => (isPassword.value && show.value ? 'text' : props.type))
</script>

<template>
  <div class="flex flex-col gap-1">
    <label v-if="label" class="text-xs font-medium" :class="dark ? 'text-[#bbc7c6]' : 'text-[var(--color-text-muted)]'">{{ label }}</label>
    <div
      class="relative flex h-12 w-full items-center rounded-[var(--radius-md)] border transition-[border-color] duration-150"
      :class="dark
        ? (error ? 'border-[#e46258]' : 'border-[#707777]/40 bg-[#011d1c] focus-within:border-[#cbfffc] focus-within:border-2')
        : (error ? 'border-[var(--color-error)] bg-[var(--color-surface)]' : 'border-[var(--color-border)] bg-[var(--color-surface)] focus-within:border-[var(--color-primary)] focus-within:border-2')"
    >
      <span v-if="icon" class="pointer-events-none absolute left-3 flex items-center" :class="dark ? 'text-[#707777]' : 'text-[var(--color-text-muted)]'">
        <component :is="icon" class="h-5 w-5" />
      </span>
      <input
        :type="inputType"
        :placeholder="placeholder"
        :value="modelValue"
        class="h-full w-full bg-transparent text-sm outline-none placeholder:text-[#707777]"
        :class="dark ? 'text-white' : 'text-[var(--color-text)]'"
        :style="icon ? { paddingLeft: '2.5rem' } : { paddingLeft: '1rem' }"
        :autocomplete="autocomplete || (isPassword ? 'current-password' : undefined)"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      />
      <button
        v-if="isPassword"
        type="button"
        class="absolute right-3 flex h-8 w-8 items-center justify-center rounded"
        :class="dark ? 'text-[#707777] hover:text-white' : 'text-[var(--color-text-muted)] hover:text-[var(--color-text)]'"
        :aria-label="show ? 'Hide password' : 'Show password'"
        @click="show = !show"
      >
        <EyeOff v-if="show" class="h-4 w-4" />
        <Eye v-else class="h-4 w-4" />
      </button>
      <span v-else-if="icon" class="w-4" />
    </div>
    <p v-if="error" class="text-xs" :class="dark ? 'text-[#e46258]' : 'text-[var(--color-error)]'">{{ error }}</p>
    <p v-else-if="hint" class="text-xs" :class="dark ? 'text-[#707777]' : 'text-[var(--color-text-muted)]'">{{ hint }}</p>
  </div>
</template>
