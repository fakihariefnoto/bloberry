<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowRight, ArrowUpRight, Boxes, KeyRound, ShieldCheck, Server, Globe,
  HardDrive, UploadCloud, Database, Cloud, Lock,
} from 'lucide-vue-next'
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
    <header class="sticky top-0 z-20 border-b border-[var(--color-border)]/60 bg-[var(--color-background)]/80 backdrop-blur">
      <div class="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
        <a href="#" class="flex items-center gap-2.5">
          <span class="flex h-8 w-8 items-center justify-center rounded-[var(--radius-md)] bg-[var(--color-primary)] shadow-[var(--shadow-sm)]">
            <Boxes class="h-4 w-4 text-[var(--color-on-primary)]" />
          </span>
          <span class="text-[15px] font-semibold tracking-tight text-[var(--color-text)]">Bloberry</span>
        </a>
        <nav class="hidden items-center gap-8 text-sm text-[var(--color-text-muted)] md:flex">
          <a href="#features" class="transition-colors hover:text-[var(--color-text)]">Platform</a>
          <a href="#drivers" class="transition-colors hover:text-[var(--color-text)]">Storage</a>
          <a href="#security" class="transition-colors hover:text-[var(--color-text)]">Security</a>
          <a href="#pricing" class="transition-colors hover:text-[var(--color-text)]">Self-host</a>
        </nav>
        <button
          class="inline-flex h-9 items-center gap-1.5 rounded-[var(--radius-md)] bg-[var(--color-primary)] px-4 text-sm font-semibold text-[var(--color-on-primary)] transition-opacity hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] focus:ring-offset-2"
          @click="primaryCta"
        >
          {{ checking ? 'Loading…' : needsSetup ? 'Get started' : 'Sign in' }}
        </button>
      </div>
    </header>

    <!-- Hero -->
    <section class="relative overflow-hidden">
      <div class="pointer-events-none absolute inset-0">
        <StorageConstellation />
      </div>
      <div class="pointer-events-none absolute inset-x-0 bottom-0 h-40 bg-gradient-to-b from-transparent to-[var(--color-background)]" />

      <div class="relative mx-auto max-w-6xl px-6 pb-16 pt-16 md:pt-20">
        <div class="mx-auto max-w-3xl text-center">
          <p class="mx-auto inline-flex items-center gap-2 rounded-full border border-[var(--color-border)] bg-[var(--color-surface-raised)] px-3 py-1 text-xs font-medium text-[var(--color-text-muted)]">
            <Lock class="h-3 w-3 text-[var(--color-primary)]" />
            Self-hosted · your credentials, your bytes
          </p>
          <h1 class="mt-5 text-balance text-[clamp(2.2rem,5vw,3.6rem)] font-bold leading-[1.08] tracking-tight text-[var(--color-text)]">
            One API over any object storage.
          </h1>
          <p class="mx-auto mt-5 max-w-xl text-pretty text-base leading-relaxed text-[var(--color-text-muted)]">
            Bloberry gives you real folders, per-folder permissions and application keys on top of
            S3, R2, OSS, GCS, Azure Blob or plain disk — without moving your bytes.
          </p>
          <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
            <button
              class="inline-flex h-12 items-center gap-2 rounded-[var(--radius-md)] bg-[var(--color-primary)] px-6 text-sm font-semibold text-[var(--color-on-primary)] shadow-[var(--shadow-md)] transition-opacity hover:opacity-90"
              @click="primaryCta"
            >
              {{ needsSetup ? 'Set up your instance' : 'Sign in to dashboard' }}
              <ArrowRight class="h-4 w-4" />
            </button>
            <a href="#features" class="inline-flex h-12 items-center gap-2 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] px-5 text-sm font-semibold text-[var(--color-text)] transition-colors hover:border-[var(--color-primary)] hover:text-[var(--color-primary)]">
              Explore the platform
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- Providers strip (driver identity) -->
    <section class="border-y border-[var(--color-border)] bg-[var(--color-surface)]">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-center gap-x-10 gap-y-4 px-6 py-8">
        <p class="w-full text-center text-xs font-medium uppercase tracking-wider text-[var(--color-text-muted)] md:w-auto md:text-left">
          One interface, six providers
        </p>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-[var(--color-text-muted)]"><Cloud class="h-4 w-4" /> AWS S3</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-[var(--color-text-muted)]"><Globe class="h-4 w-4" /> Cloudflare R2</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-[var(--color-text-muted)]"><Database class="h-4 w-4" /> Google GCS</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-[var(--color-text-muted)]"><Boxes class="h-4 w-4" /> Alibaba OSS</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-[var(--color-text-muted)]"><HardDrive class="h-4 w-4" /> Azure Blob</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-[var(--color-text-muted)]"><Server class="h-4 w-4" /> Local disk</span>
      </div>
    </section>

    <!-- Features: asymmetric bento -->
    <section id="features" class="mx-auto max-w-6xl px-6 py-20 md:py-24">
      <div class="grid gap-4 md:grid-cols-6">
        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-7 md:col-span-4">
          <span class="inline-flex h-10 w-10 items-center justify-center rounded-[var(--radius-sm)] bg-[var(--color-primary-subtle)]">
            <KeyRound class="h-5 w-5 text-[var(--color-primary)]" />
          </span>
          <h3 class="mt-5 text-xl font-bold tracking-tight text-[var(--color-text)]">Permissions that match your folders</h3>
          <p class="mt-2 max-w-md text-sm leading-relaxed text-[var(--color-text-muted)]">
            Five human roles, folder-level grants, and scoped application keys — all resolved in one
            permission model across web, mobile, CLI and SDKs.
          </p>
          <div class="mt-5 flex flex-wrap gap-2">
            <span class="rounded-full bg-[var(--color-surface)] px-2.5 py-1 text-xs font-medium text-[var(--color-text-muted)]">Roles &amp; grants</span>
            <span class="rounded-full bg-[var(--color-surface)] px-2.5 py-1 text-xs font-medium text-[var(--color-text-muted)]">Scoped access keys</span>
            <span class="rounded-full bg-[var(--color-surface)] px-2.5 py-1 text-xs font-medium text-[var(--color-text-muted)]">Two-factor auth</span>
          </div>
        </div>

        <div class="flex flex-col justify-between rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface)] p-7 md:col-span-2">
          <div>
            <span class="inline-flex h-10 w-10 items-center justify-center rounded-[var(--radius-sm)] bg-[var(--color-primary-subtle)]">
              <UploadCloud class="h-5 w-5 text-[var(--color-primary)]" />
            </span>
            <h3 class="mt-5 text-lg font-bold tracking-tight text-[var(--color-text)]">Bytes bypass you</h3>
            <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">
              Presigned uploads go straight to the provider. Your server never touches the payload.
            </p>
          </div>
          <a href="#pricing" class="mt-5 inline-flex items-center gap-1.5 text-sm font-semibold text-[var(--color-primary)] hover:underline">
            How it works <ArrowUpRight class="h-4 w-4" />
          </a>
        </div>

        <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-7 md:col-span-2">
          <span class="inline-flex h-10 w-10 items-center justify-center rounded-[var(--radius-sm)] bg-[var(--color-primary-subtle)]">
            <ShieldCheck class="h-5 w-5 text-[var(--color-primary)]" />
          </span>
          <h3 class="mt-5 text-lg font-bold tracking-tight text-[var(--color-text)]">Credentials, encrypted</h3>
          <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">
            Provider keys are envelope-encrypted at rest with a key that never lives in the database.
          </p>
        </div>

        <div class="flex flex-col justify-between rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-surface)] p-7 md:col-span-4">
          <div>
            <span class="inline-flex h-10 w-10 items-center justify-center rounded-[var(--radius-sm)] bg-[var(--color-primary-subtle)]">
              <Server class="h-5 w-5 text-[var(--color-primary)]" />
            </span>
            <h3 class="mt-5 text-xl font-bold tracking-tight text-[var(--color-text)]">One binary to deploy</h3>
            <p class="mt-2 max-w-md text-sm leading-relaxed text-[var(--color-text-muted)]">
              The API, the dashboard and short links ship as a single Go binary. From a fresh VPS to a
              running install in about fifteen minutes.
            </p>
          </div>
          <div class="mt-5 grid grid-cols-3 gap-3 border-t border-[var(--color-border)] pt-5">
            <div>
              <p class="text-2xl font-bold tracking-tight text-[var(--color-text)]">1</p>
              <p class="mt-1 text-xs text-[var(--color-text-muted)]">binary</p>
            </div>
            <div>
              <p class="text-2xl font-bold tracking-tight text-[var(--color-text)]">5 GiB</p>
              <p class="mt-1 text-xs text-[var(--color-text-muted)]">max object</p>
            </div>
            <div>
              <p class="text-2xl font-bold tracking-tight text-[var(--color-text)]">6</p>
              <p class="mt-1 text-xs text-[var(--color-text-muted)]">drivers</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Drivers -->
    <section id="drivers" class="border-y border-[var(--color-border)] bg-[var(--color-surface)]">
      <div class="mx-auto max-w-6xl px-6 py-16 md:py-20">
        <div class="max-w-2xl">
          <h2 class="text-2xl font-bold tracking-tight text-[var(--color-text)] md:text-3xl">Your storage, your rules.</h2>
          <p class="mt-3 text-base leading-relaxed text-[var(--color-text-muted)]">
            Every driver implements one interface — presigning, multipart, streaming and health checks —
            so switching providers is a configuration change, not a rewrite.
          </p>
        </div>
        <div class="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div class="flex items-start gap-3 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
            <Cloud class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
            <div><p class="text-sm font-semibold text-[var(--color-text)]">S3 &amp; R2</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">AWS S3, Cloudflare R2, MinIO, Backblaze, Spaces, Wasabi</p></div>
          </div>
          <div class="flex items-start gap-3 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
            <Database class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
            <div><p class="text-sm font-semibold text-[var(--color-text)]">Google Cloud Storage</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">Service-account signing, IAM-aware</p></div>
          </div>
          <div class="flex items-start gap-3 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
            <Boxes class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
            <div><p class="text-sm font-semibold text-[var(--color-text)]">Alibaba OSS &amp; Azure Blob</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">First-party SDKs, native signature schemes</p></div>
          </div>
          <div class="flex items-start gap-3 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
            <Server class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
            <div><p class="text-sm font-semibold text-[var(--color-text)]">Local disk</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">Plain VPS volume, HMAC-signed presigned URLs</p></div>
          </div>
          <div class="flex items-start gap-3 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
            <KeyRound class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
            <div><p class="text-sm font-semibold text-[var(--color-text)]">Application keys</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">Scoped to folders and permissions, revocable</p></div>
          </div>
          <div class="flex items-start gap-3 rounded-[var(--radius-md)] border border-[var(--color-border)] bg-[var(--color-surface-raised)] p-4">
            <Globe class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
            <div><p class="text-sm font-semibold text-[var(--color-text)]">Short links &amp; shares</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">Signed links, short URLs, revocable anytime</p></div>
          </div>
        </div>
      </div>
    </section>

    <!-- How it works (3 steps) -->
    <section id="pricing" class="mx-auto max-w-6xl px-6 py-20 md:py-24">
      <div class="max-w-2xl">
        <h2 class="text-2xl font-bold tracking-tight text-[var(--color-text)] md:text-3xl">From VPS to first file in minutes.</h2>
      </div>
      <div class="mt-10 grid gap-8 md:grid-cols-3">
        <div class="relative">
          <p class="text-sm font-bold text-[var(--color-primary)]">01</p>
          <h3 class="mt-3 text-lg font-semibold tracking-tight text-[var(--color-text)]">Run one binary</h3>
          <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">A single Go server embeds the dashboard, API and worker. systemd keeps it alive.</p>
        </div>
        <div class="relative">
          <p class="text-sm font-bold text-[var(--color-primary)]">02</p>
          <h3 class="mt-3 text-lg font-semibold tracking-tight text-[var(--color-text)]">Connect a provider</h3>
          <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">Register any storage backend once. Credentials are encrypted before they touch the database.</p>
        </div>
        <div class="relative">
          <p class="text-sm font-bold text-[var(--color-primary)]">03</p>
          <h3 class="mt-3 text-lg font-semibold tracking-tight text-[var(--color-text)]">Invite your team</h3>
          <p class="mt-2 text-sm leading-relaxed text-[var(--color-text-muted)]">Real folders with per-folder grants, plus keys for every integration.</p>
        </div>
      </div>
    </section>

    <!-- Security -->
    <section id="security" class="border-y border-[var(--color-border)] bg-[var(--color-surface)]">
      <div class="mx-auto max-w-6xl px-6 py-16 md:py-20">
        <div class="grid gap-10 md:grid-cols-2 md:items-center">
          <div>
            <h2 class="text-2xl font-bold tracking-tight text-[var(--color-text)] md:text-3xl">Built for the trust self-hosting demands.</h2>
            <p class="mt-3 text-base leading-relaxed text-[var(--color-text-muted)]">
              Every decision in the data path is auditable — and every secret is encrypted before storage.
            </p>
            <a href="#features" class="mt-6 inline-flex items-center gap-1.5 text-sm font-semibold text-[var(--color-primary)] hover:underline">
              Read how it works <ArrowUpRight class="h-4 w-4" />
            </a>
          </div>
          <ul class="divide-y divide-[var(--color-border)] border-y border-[var(--color-border)]">
            <li class="flex items-start gap-3 py-4">
              <Lock class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
              <div><p class="text-sm font-semibold text-[var(--color-text)]">Envelope encryption at rest</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">A database dump cannot decrypt provider credentials.</p></div>
            </li>
            <li class="flex items-start gap-3 py-4">
              <KeyRound class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
              <div><p class="text-sm font-semibold text-[var(--color-text)]">Two-factor authentication</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">TOTP with single-use backup codes, on every human login.</p></div>
            </li>
            <li class="flex items-start gap-3 py-4">
              <ShieldCheck class="mt-0.5 h-5 w-5 shrink-0 text-[var(--color-primary)]" />
              <div><p class="text-sm font-semibold text-[var(--color-text)]">Full audit log</p><p class="mt-1 text-xs leading-relaxed text-[var(--color-text-muted)]">Every access-key and grant action is queryable per tenant.</p></div>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="mx-auto max-w-6xl px-6 py-20 md:py-24">
      <div class="rounded-[var(--radius-lg)] border border-[var(--color-border)] bg-[var(--color-primary-subtle)] px-8 py-14 text-center">
        <h2 class="mx-auto max-w-md text-2xl font-bold tracking-tight text-[var(--color-text)] md:text-3xl">
          {{ needsSetup ? 'Set up your instance.' : 'Ready when you are.' }}
        </h2>
        <p class="mx-auto mt-3 max-w-md text-sm leading-relaxed text-[var(--color-text-muted)]">
          {{ needsSetup ? 'Create your admin account and first tenant in about a minute.' : 'Sign in to manage your folders, keys and shares.' }}
        </p>
        <button
          class="mt-7 inline-flex h-12 items-center gap-2 rounded-[var(--radius-md)] bg-[var(--color-primary)] px-6 text-sm font-semibold text-[var(--color-on-primary)] shadow-[var(--shadow-md)] transition-opacity hover:opacity-90"
          @click="primaryCta"
        >
          {{ needsSetup ? 'Get started' : 'Sign in' }}
          <ArrowRight class="h-4 w-4" />
        </button>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-[var(--color-border)]">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 py-8 md:flex-row">
        <div class="flex items-center gap-2.5">
          <span class="flex h-7 w-7 items-center justify-center rounded-[var(--radius-sm)] bg-[var(--color-primary)]">
            <Boxes class="h-3.5 w-3.5 text-[var(--color-on-primary)]" />
          </span>
          <span class="text-sm font-semibold text-[var(--color-text)]">Bloberry</span>
        </div>
        <p class="text-xs text-[var(--color-text-muted)]">Storage-agnostic object service — self-hosted, single binary.</p>
      </div>
    </footer>
  </div>
</template>
