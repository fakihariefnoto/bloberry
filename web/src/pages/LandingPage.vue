<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowUpRight, Boxes, KeyRound, ShieldCheck, Server, Globe, HardDrive, UploadCloud } from 'lucide-vue-next'
import { api } from '../lib/api'
import StorageConstellation from '../components/StorageConstellation.vue'

const router = useRouter()
const needsSetup = ref(false)
const checking = ref(true)

onMounted(async () => {
  try {
    const st = await api.get<{ needs_setup: boolean }>('/setup/status')
    needsSetup.value = st.needs_setup
  } catch {
    needsSetup.value = false
  } finally {
    checking.value = false
  }
})

function primaryCta() {
  router.push(needsSetup.value ? { name: 'setup' } : { name: 'login' })
}
</script>

<template>
  <div class="min-h-screen bg-[var(--color-background)] text-[var(--color-text)]">
    <!-- Nav -->
    <header class="relative z-10 mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
      <div class="flex items-center gap-2">
        <span class="flex h-9 w-9 items-center justify-center rounded-[var(--radius-md)] bg-[var(--color-primary-subtle)]">
          <Boxes class="h-4 w-4 text-[var(--color-primary)]" />
        </span>
        <span class="text-sm font-semibold tracking-wide text-[var(--color-text)]">Bloberry</span>
      </div>
      <nav class="hidden items-center gap-7 text-sm text-[var(--color-text-muted)] md:flex">
        <a href="#features" class="transition-colors hover:text-[var(--color-primary)]">Features</a>
        <a href="#how" class="transition-colors hover:text-[var(--color-primary)]">How it works</a>
        <a href="#security" class="transition-colors hover:text-[var(--color-primary)]">Security</a>
      </nav>
      <div class="flex items-center gap-3">
        <button
          class="h-10 rounded-[var(--radius-md)] bg-[var(--color-primary)] px-5 text-sm font-semibold text-[var(--color-on-primary)] transition-opacity hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] focus:ring-offset-2"
          @click="primaryCta"
        >
          {{ checking ? 'Loading' : needsSetup ? 'Get started' : 'Sign in' }}
        </button>
      </div>
    </header>

    <!-- Hero with storage constellation -->
    <section class="relative mx-auto max-w-6xl px-6 pb-24 pt-16 md:pt-24">
      <!-- Wave background — clipped to the hero, never overlapping lower cards -->
      <div class="pointer-events-none absolute inset-0 overflow-hidden">
        <StorageConstellation />
      </div>

      <div class="relative z-10 mx-auto max-w-3xl text-center">
        <p class="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--color-primary)]">Storage-agnostic object service</p>
        <h1 class="mt-6 text-[clamp(2.4rem,5.5vw,4rem)] font-bold leading-tight tracking-tight text-[var(--color-text)]">
          One API, one dashboard,
          <span class="text-[var(--color-primary)]">any storage</span>
        </h1>
        <p class="mx-auto mt-6 max-w-xl text-base leading-relaxed text-[var(--color-text-muted)]">
          Bloberry sits between your applications and S3, R2, OSS, GCS, Azure Blob or plain disk —
          with real folders, per-folder permissions, application keys and temporary links over any provider.
        </p>
        <div class="mt-8 flex items-center justify-center gap-4">
          <button
            class="h-12 rounded-[var(--radius-md)] bg-[var(--color-primary)] px-7 text-sm font-semibold text-[var(--color-on-primary)] shadow-[var(--shadow-sm)] transition-opacity hover:opacity-90"
            @click="primaryCta"
          >
            {{ needsSetup ? 'Set up Bloberry' : 'Sign in to dashboard' }}
          </button>
          <a href="#features" class="px-2 py-3 text-sm font-semibold text-[var(--color-primary)] hover:underline">Explore</a>
        </div>
      </div>
    </section>

    <!-- Stats -->
    <section class="mx-auto max-w-6xl px-6">
      <div class="grid grid-cols-2 gap-8 rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface)] p-8 md:grid-cols-4">
        <div class="text-center">
          <p class="text-4xl font-bold text-[var(--color-primary)]">6</p>
          <p class="mt-2 text-xs font-medium uppercase tracking-wide text-[var(--color-text-muted)]">Storage drivers</p>
        </div>
        <div class="text-center">
          <p class="text-4xl font-bold text-[var(--color-primary)]">5GB</p>
          <p class="mt-2 text-xs font-medium uppercase tracking-wide text-[var(--color-text-muted)]">Max object size</p>
        </div>
        <div class="text-center">
          <p class="text-4xl font-bold text-[var(--color-primary)]">100%</p>
          <p class="mt-2 text-xs font-medium uppercase tracking-wide text-[var(--color-text-muted)]">Self-hosted</p>
        </div>
        <div class="text-center">
          <p class="text-4xl font-bold text-[var(--color-primary)]">1</p>
          <p class="mt-2 text-xs font-medium uppercase tracking-wide text-[var(--color-text-muted)]">Binary to deploy</p>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section id="features" class="mx-auto max-w-6xl px-6 py-24">
      <p class="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--color-primary)]">Features</p>
      <div class="mt-8 grid gap-6 md:grid-cols-3">
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-7">
          <div class="flex items-start justify-between">
            <HardDrive class="h-6 w-6 text-[var(--color-primary)]" />
            <span class="flex h-8 w-8 items-center justify-center rounded-[var(--radius-sm)] bg-[var(--color-primary-subtle)]">
              <ArrowUpRight class="h-4 w-4 text-[var(--color-primary)]" />
            </span>
          </div>
          <h3 class="mt-6 text-lg font-semibold text-[var(--color-text)]">Any provider</h3>
          <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">S3, R2, OSS, GCS, Azure Blob or local disk behind one interface.</p>
        </div>
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-7">
          <div class="flex items-start justify-between">
            <KeyRound class="h-6 w-6 text-[var(--color-primary)]" />
            <span class="flex h-8 w-8 items-center justify-center rounded-[var(--radius-sm)] bg-[var(--color-primary-subtle)]">
              <ArrowUpRight class="h-4 w-4 text-[var(--color-primary)]" />
            </span>
          </div>
          <h3 class="mt-6 text-lg font-semibold text-[var(--color-text)]">Real permissions</h3>
          <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">Roles, folder-level grants and scoped application keys.</p>
        </div>
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-7">
          <div class="flex items-start justify-between">
            <UploadCloud class="h-6 w-6 text-[var(--color-primary)]" />
            <span class="flex h-8 w-8 items-center justify-center rounded-[var(--radius-sm)] bg-[var(--color-primary-subtle)]">
              <ArrowUpRight class="h-4 w-4 text-[var(--color-primary)]" />
            </span>
          </div>
          <h3 class="mt-6 text-lg font-semibold text-[var(--color-text)]">Bytes bypass you</h3>
          <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">Presigned uploads go straight to the provider. Your server never touches the payload.</p>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section id="how" class="mx-auto max-w-6xl px-6 pb-24">
      <p class="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--color-primary)]">How it works</p>
      <div class="mt-8 grid gap-8 md:grid-cols-3">
        <div class="flex gap-4">
          <span class="text-3xl font-bold text-[var(--color-primary)]">01</span>
          <div>
            <h4 class="text-base font-semibold text-[var(--color-text)]">Connect a backend</h4>
            <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">Register any provider — credentials are envelope-encrypted at rest.</p>
          </div>
        </div>
        <div class="flex gap-4">
          <span class="text-3xl font-bold text-[var(--color-primary)]">02</span>
          <div>
            <h4 class="text-base font-semibold text-[var(--color-text)]">Give your team folders</h4>
            <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">Real folders with per-folder grants, roles and tenants.</p>
          </div>
        </div>
        <div class="flex gap-4">
          <span class="text-3xl font-bold text-[var(--color-primary)]">03</span>
          <div>
            <h4 class="text-base font-semibold text-[var(--color-text)]">Share and integrate</h4>
            <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">Signed links, short URLs, SDKs and a CLI that script the same API.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Security -->
    <section id="security" class="mx-auto max-w-6xl px-6 pb-24">
      <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface)] px-10 py-16 text-center">
        <ShieldCheck class="mx-auto h-8 w-8 text-[var(--color-primary)]" />
        <p class="mt-4 text-xs font-semibold uppercase tracking-[0.12em] text-[var(--color-primary)]">Security</p>
        <h2 class="mx-auto mt-3 max-w-md text-3xl font-bold tracking-tight text-[var(--color-text)]">
          Credentials encrypted. Access audited.
        </h2>
        <div class="mx-auto mt-10 grid max-w-2xl gap-8 md:grid-cols-3">
          <div class="text-center">
            <Globe class="mx-auto h-5 w-5 text-[var(--color-primary)]" />
            <p class="mt-3 text-sm text-[var(--color-text-muted)]">Envelope encryption at rest</p>
          </div>
          <div class="text-center">
            <KeyRound class="mx-auto h-5 w-5 text-[var(--color-primary)]" />
            <p class="mt-3 text-sm text-[var(--color-text-muted)]">Two-factor authentication</p>
          </div>
          <div class="text-center">
            <Server class="mx-auto h-5 w-5 text-[var(--color-primary)]" />
            <p class="mt-3 text-sm text-[var(--color-text-muted)]">Full audit log per tenant</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="mx-auto max-w-6xl px-6 pb-12">
      <div class="flex flex-col items-center justify-between gap-3 border-t border-[var(--color-border)] pt-6 md:flex-row">
        <p class="text-xs font-semibold uppercase tracking-wide text-[var(--color-text-muted)]">Bloberry</p>
        <p class="text-xs text-[var(--color-text-muted)]">Storage-agnostic object service — self-hosted</p>
      </div>
    </footer>
  </div>
</template>
