<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowRight, Boxes, KeyRound, ShieldCheck, Server, Globe,
  HardDrive, Database, Cloud, Lock, Sparkles, Zap, Terminal, Layers,
} from 'lucide-vue-next'
import { api } from '../lib/api'
import StorageConstellation from '../components/StorageConstellation.vue'

const router = useRouter()
const needsSetup = ref(false)
const checking = ref(true)

onMounted(async () => {
  // Router guard redirects authenticated users to /app/files before this
  // mounts, so the landing never renders for them.
  document.documentElement.dataset.theme = 'dark'
  try {
    const st = await api.get<{ needs_setup: boolean }>('/setup/status')
    needsSetup.value = st.needs_setup
    // App already set up and user not signed in → straight to login.
    if (!st.needs_setup) {
      router.replace({ name: 'login' })
      return
    }
  } catch {
    needsSetup.value = false
  } finally {
    checking.value = false
  }
})

onUnmounted(() => {
  delete document.documentElement.dataset.theme
})

function primaryCta() {
  router.push(needsSetup.value ? { name: 'setup' } : { name: 'login' })
}
</script>

<template>
  <div class="min-h-screen bg-[#07060f] text-white antialiased">
    <!-- Ambient background -->
    <div class="pointer-events-none fixed inset-0">
      <div class="absolute inset-0 opacity-60" style="background: radial-gradient(60% 45% at 50% 0%, rgba(139,125,235,0.18), transparent 70%)" />
      <div class="absolute inset-0 opacity-40" style="background: radial-gradient(45% 35% at 85% 60%, rgba(30,20,107,0.5), transparent 70%)" />
      <div class="absolute inset-0 opacity-30" style="background: radial-gradient(40% 30% at 10% 80%, rgba(139,125,235,0.12), transparent 70%)" />
    </div>

    <!-- Nav -->
    <header class="sticky top-0 z-30 border-b border-white/5 bg-[#07060f]/70 backdrop-blur-xl">
      <div class="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
        <a href="#" class="flex items-center gap-2.5">
          <img src="/bloberry-icon.png" alt="Bloberry" class="h-8 w-8 rounded-lg object-cover" />
          <span class="text-[15px] font-semibold tracking-tight">Bloberry</span>
        </a>
        <nav class="hidden items-center gap-8 text-sm text-white/60 md:flex">
          <a href="#platform" class="transition-colors hover:text-white">Platform</a>
          <a href="#storage" class="transition-colors hover:text-white">Storage</a>
          <a href="#security" class="transition-colors hover:text-white">Security</a>
          <a href="#selfhost" class="transition-colors hover:text-white">Self-host</a>
        </nav>
        <div class="flex items-center gap-3">
          <button
            class="inline-flex h-9 items-center gap-1.5 rounded-lg bg-gradient-to-r from-[#8b7deb] to-[#6f5ce0] px-4 text-sm font-semibold text-white shadow-[0_0_20px_rgba(139,125,235,0.35)] transition-opacity hover:opacity-90"
            @click="primaryCta"
          >
            {{ checking ? 'Loading…' : needsSetup ? 'Get started' : 'Sign in' }}
          </button>
        </div>
      </div>
    </header>

    <!-- Hero -->
    <section class="relative">
      <div class="pointer-events-none absolute inset-0">
        <StorageConstellation />
      </div>
      <div class="pointer-events-none absolute inset-x-0 bottom-0 h-48 bg-gradient-to-b from-transparent to-[#07060f]" />

      <div class="relative mx-auto max-w-6xl px-6 pb-20 pt-20 md:pt-24">
        <div class="mx-auto max-w-3xl text-center">
          <p class="mx-auto inline-flex items-center gap-2 rounded-full border border-[#8b7deb]/30 bg-[#8b7deb]/10 px-3.5 py-1.5 text-xs font-medium text-[#c4b5fd]">
            <Sparkles class="h-3 w-3" />
            Self-hosted · your credentials, your bytes
          </p>
          <h1 class="mt-6 text-balance text-[clamp(2.4rem,6vw,4.2rem)] font-bold leading-[1.05] tracking-tight">
            One API over
            <span class="bg-gradient-to-r from-[#a78bfa] via-[#8b7deb] to-[#4c3fd4] bg-clip-text text-transparent">any object storage.</span>
          </h1>
          <p class="mx-auto mt-6 max-w-xl text-pretty text-base leading-relaxed text-white/60">
            Real folders, per-folder permissions and application keys on top of S3, R2, OSS, GCS,
            Azure Blob or plain disk — without moving your bytes.
          </p>
          <div class="mt-9 flex flex-wrap items-center justify-center gap-3">
            <button
              class="group inline-flex h-12 items-center gap-2 rounded-lg bg-gradient-to-r from-[#8b7deb] to-[#4c3fd4] px-7 text-sm font-semibold text-white shadow-[0_0_32px_rgba(139,125,235,0.4)] transition-all hover:shadow-[0_0_48px_rgba(139,125,235,0.6)]"
              @click="primaryCta"
            >
              {{ needsSetup ? 'Set up your instance' : 'Enter the dashboard' }}
              <ArrowRight class="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
            </button>
            <a
              href="#platform"
              class="inline-flex h-12 items-center gap-2 rounded-lg border border-white/10 bg-white/5 px-6 text-sm font-semibold text-white/80 backdrop-blur transition-colors hover:border-[#8b7deb]/50 hover:text-white"
            >
              Explore the platform
            </a>
          </div>

          <div class="mt-10 flex items-center justify-center gap-2 text-xs text-white/40">
            <Terminal class="h-3.5 w-3.5" />
            <span class="font-mono">bloberry://assets — one path for every provider</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Providers -->
    <section class="relative border-y border-white/5 bg-white/[0.02]">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-center gap-x-12 gap-y-4 px-6 py-9">
        <span class="inline-flex items-center gap-2 text-sm font-medium text-white/50"><Cloud class="h-4 w-4 text-[#8b7deb]" /> AWS S3</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-white/50"><Globe class="h-4 w-4 text-[#8b7deb]" /> Cloudflare R2</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-white/50"><Database class="h-4 w-4 text-[#8b7deb]" /> Google GCS</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-white/50"><Boxes class="h-4 w-4 text-[#8b7deb]" /> Alibaba OSS</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-white/50"><HardDrive class="h-4 w-4 text-[#8b7deb]" /> Azure Blob</span>
        <span class="inline-flex items-center gap-2 text-sm font-medium text-white/50"><Server class="h-4 w-4 text-[#8b7deb]" /> Local disk</span>
      </div>
    </section>

    <!-- Platform bento -->
    <section id="platform" class="relative mx-auto max-w-6xl px-6 py-24">
      <div class="grid gap-4 md:grid-cols-6">
        <div class="group relative overflow-hidden rounded-2xl border border-white/10 bg-white/[0.03] p-7 backdrop-blur transition-colors hover:border-[#8b7deb]/40 md:col-span-4">
          <div class="pointer-events-none absolute -right-20 -top-20 h-56 w-56 rounded-full bg-[#8b7deb]/10 blur-3xl" />
          <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#8b7deb]/30 to-[#2a1a8f]/30">
            <KeyRound class="h-5 w-5 text-[#c4b5fd]" />
          </span>
          <h3 class="mt-5 text-xl font-bold tracking-tight">Permissions that match your folders</h3>
          <p class="mt-2 max-w-md text-sm leading-relaxed text-white/60">
            Five human roles, folder-level grants, and scoped application keys — one permission model
            across web, mobile, CLI and SDKs.
          </p>
          <div class="mt-5 flex flex-wrap gap-2">
            <span class="rounded-full border border-white/10 bg-white/5 px-2.5 py-1 text-xs font-medium text-white/70">Roles &amp; grants</span>
            <span class="rounded-full border border-white/10 bg-white/5 px-2.5 py-1 text-xs font-medium text-white/70">Scoped access keys</span>
            <span class="rounded-full border border-white/10 bg-white/5 px-2.5 py-1 text-xs font-medium text-white/70">Two-factor auth</span>
          </div>
        </div>

        <div class="flex flex-col justify-between rounded-2xl border border-white/10 bg-white/[0.03] p-7 backdrop-blur transition-colors hover:border-[#8b7deb]/40 md:col-span-2">
          <div>
            <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#8b7deb]/30 to-[#2a1a8f]/30">
              <Zap class="h-5 w-5 text-[#c4b5fd]" />
            </span>
            <h3 class="mt-5 text-lg font-bold tracking-tight">Bytes bypass you</h3>
            <p class="mt-2 text-sm leading-relaxed text-white/60">
              Presigned uploads go straight to the provider. Your server never touches the payload.
            </p>
          </div>
          <p class="mt-5 font-mono text-xs text-white/40">PUT → provider, not server</p>
        </div>

        <div class="rounded-2xl border border-white/10 bg-white/[0.03] p-7 backdrop-blur transition-colors hover:border-[#8b7deb]/40 md:col-span-2">
          <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#8b7deb]/30 to-[#2a1a8f]/30">
            <ShieldCheck class="h-5 w-5 text-[#c4b5fd]" />
          </span>
          <h3 class="mt-5 text-lg font-bold tracking-tight">Encrypted at rest</h3>
          <p class="mt-2 text-sm leading-relaxed text-white/60">
            Provider keys are envelope-encrypted; the key never lives in the database.
          </p>
        </div>

        <div class="flex flex-col justify-between rounded-2xl border border-white/10 bg-white/[0.03] p-7 backdrop-blur transition-colors hover:border-[#8b7deb]/40 md:col-span-4">
          <div class="flex flex-wrap items-start justify-between gap-6">
            <div>
              <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#8b7deb]/30 to-[#2a1a8f]/30">
                <Server class="h-5 w-5 text-[#c4b5fd]" />
              </span>
              <h3 class="mt-5 text-xl font-bold tracking-tight">One binary to deploy</h3>
              <p class="mt-2 max-w-md text-sm leading-relaxed text-white/60">
                The API, dashboard and short links ship as a single Go binary. From a fresh VPS to a
                running install in about fifteen minutes.
              </p>
            </div>
            <div class="grid grid-cols-3 gap-8">
              <div class="text-center">
                <p class="bg-gradient-to-r from-[#a78bfa] to-[#8b7deb] bg-clip-text text-3xl font-bold tracking-tight text-transparent">1</p>
                <p class="mt-1 text-xs text-white/40">binary</p>
              </div>
              <div class="text-center">
                <p class="bg-gradient-to-r from-[#a78bfa] to-[#8b7deb] bg-clip-text text-3xl font-bold tracking-tight text-transparent">5 GiB</p>
                <p class="mt-1 text-xs text-white/40">max object</p>
              </div>
              <div class="text-center">
                <p class="bg-gradient-to-r from-[#a78bfa] to-[#8b7deb] bg-clip-text text-3xl font-bold tracking-tight text-transparent">6</p>
                <p class="mt-1 text-xs text-white/40">drivers</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Storage drivers -->
    <section id="storage" class="border-y border-white/5 bg-white/[0.02]">
      <div class="mx-auto max-w-6xl px-6 py-24">
        <div class="grid gap-10 lg:grid-cols-5 lg:items-center">
          <div class="lg:col-span-2">
            <p class="inline-flex items-center gap-1.5 font-mono text-xs uppercase tracking-widest text-[#8b7deb]">
              <Layers class="h-3.5 w-3.5" /> Storage layer
            </p>
            <h2 class="mt-4 text-2xl font-bold tracking-tight md:text-3xl">Your storage, your rules.</h2>
            <p class="mt-3 text-base leading-relaxed text-white/60">
              Every driver implements one interface — presigning, multipart, streaming, health — so
              switching providers is a config change, not a rewrite.
            </p>
          </div>
          <div class="grid gap-3 sm:grid-cols-2 lg:col-span-3">
            <div class="rounded-xl border border-white/10 bg-[#0b0a16] p-4">
              <div class="flex items-center gap-2.5">
                <Cloud class="h-4 w-4 text-[#8b7deb]" />
                <p class="text-sm font-semibold">S3 &amp; R2</p>
              </div>
              <p class="mt-1.5 text-xs leading-relaxed text-white/50">AWS S3, Cloudflare R2, MinIO, Backblaze, Spaces, Wasabi</p>
            </div>
            <div class="rounded-xl border border-white/10 bg-[#0b0a16] p-4">
              <div class="flex items-center gap-2.5">
                <Database class="h-4 w-4 text-[#8b7deb]" />
                <p class="text-sm font-semibold">Google Cloud Storage</p>
              </div>
              <p class="mt-1.5 text-xs leading-relaxed text-white/50">Service-account signing, IAM-aware</p>
            </div>
            <div class="rounded-xl border border-white/10 bg-[#0b0a16] p-4">
              <div class="flex items-center gap-2.5">
                <Boxes class="h-4 w-4 text-[#8b7deb]" />
                <p class="text-sm font-semibold">Alibaba OSS &amp; Azure</p>
              </div>
              <p class="mt-1.5 text-xs leading-relaxed text-white/50">First-party SDKs, native signature schemes</p>
            </div>
            <div class="rounded-xl border border-white/10 bg-[#0b0a16] p-4">
              <div class="flex items-center gap-2.5">
                <Server class="h-4 w-4 text-[#8b7deb]" />
                <p class="text-sm font-semibold">Local disk</p>
              </div>
              <p class="mt-1.5 text-xs leading-relaxed text-white/50">Plain VPS volume, HMAC-signed presigned URLs</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section id="selfhost" class="mx-auto max-w-6xl px-6 py-24">
      <div class="grid gap-8 md:grid-cols-3">
        <div class="relative rounded-2xl border border-white/10 bg-white/[0.02] p-7">
          <span class="bg-gradient-to-r from-[#a78bfa] to-[#8b7deb] bg-clip-text font-mono text-sm font-bold text-transparent">01</span>
          <h3 class="mt-4 text-lg font-semibold tracking-tight">Run one binary</h3>
          <p class="mt-2 text-sm leading-relaxed text-white/60">A single Go server embeds the dashboard, API and worker. systemd keeps it alive.</p>
        </div>
        <div class="relative rounded-2xl border border-white/10 bg-white/[0.02] p-7">
          <span class="bg-gradient-to-r from-[#a78bfa] to-[#8b7deb] bg-clip-text font-mono text-sm font-bold text-transparent">02</span>
          <h3 class="mt-4 text-lg font-semibold tracking-tight">Connect a provider</h3>
          <p class="mt-2 text-sm leading-relaxed text-white/60">Register any storage engine once. Credentials are encrypted before they touch the database.</p>
        </div>
        <div class="relative rounded-2xl border border-white/10 bg-white/[0.02] p-7">
          <span class="bg-gradient-to-r from-[#a78bfa] to-[#8b7deb] bg-clip-text font-mono text-sm font-bold text-transparent">03</span>
          <h3 class="mt-4 text-lg font-semibold tracking-tight">Invite your team</h3>
          <p class="mt-2 text-sm leading-relaxed text-white/60">Real folders with per-folder grants, plus keys for every integration.</p>
        </div>
      </div>
    </section>

    <!-- Security -->
    <section id="security" class="border-y border-white/5 bg-white/[0.02]">
      <div class="mx-auto max-w-6xl px-6 py-24">
        <div class="grid gap-10 md:grid-cols-2 md:items-center">
          <div>
            <p class="inline-flex items-center gap-1.5 font-mono text-xs uppercase tracking-widest text-[#8b7deb]">
              <ShieldCheck class="h-3.5 w-3.5" /> Security
            </p>
            <h2 class="mt-4 text-2xl font-bold tracking-tight md:text-3xl">Built for the trust self-hosting demands.</h2>
            <p class="mt-3 text-base leading-relaxed text-white/60">
              Every decision in the data path is auditable — and every secret is encrypted before storage.
            </p>
          </div>
          <ul class="divide-y divide-white/5 rounded-2xl border border-white/10 bg-[#0b0a16]">
            <li class="flex items-start gap-3 px-6 py-5">
              <Lock class="mt-0.5 h-5 w-5 shrink-0 text-[#8b7deb]" />
              <div><p class="text-sm font-semibold">Envelope encryption at rest</p><p class="mt-1 text-xs leading-relaxed text-white/50">A database dump cannot decrypt provider credentials.</p></div>
            </li>
            <li class="flex items-start gap-3 px-6 py-5">
              <KeyRound class="mt-0.5 h-5 w-5 shrink-0 text-[#8b7deb]" />
              <div><p class="text-sm font-semibold">Two-factor authentication</p><p class="mt-1 text-xs leading-relaxed text-white/50">TOTP with single-use backup codes, on every human login.</p></div>
            </li>
            <li class="flex items-start gap-3 px-6 py-5">
              <ShieldCheck class="mt-0.5 h-5 w-5 shrink-0 text-[#8b7deb]" />
              <div><p class="text-sm font-semibold">Full audit log</p><p class="mt-1 text-xs leading-relaxed text-white/50">Every access-key and grant action is queryable per project.</p></div>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="relative mx-auto max-w-6xl px-6 py-24">
      <div class="relative overflow-hidden rounded-2xl border border-[#8b7deb]/20 bg-gradient-to-br from-[#1e146b]/40 via-[#0b0a16] to-[#0b0a16] px-8 py-16 text-center">
        <div class="pointer-events-none absolute -left-24 -top-24 h-72 w-72 rounded-full bg-[#8b7deb]/15 blur-3xl" />
        <div class="pointer-events-none absolute -bottom-24 -right-24 h-72 w-72 rounded-full bg-[#2a1a8f]/40 blur-3xl" />
        <div class="relative">
          <h2 class="mx-auto max-w-md text-2xl font-bold tracking-tight md:text-3xl">
            {{ needsSetup ? 'Set up your instance.' : 'Ready when you are.' }}
          </h2>
          <p class="mx-auto mt-3 max-w-md text-sm leading-relaxed text-white/60">
            {{ needsSetup ? 'Create your admin account and first project in about a minute.' : 'Sign in to manage your folders, keys and shares.' }}
          </p>
          <button
            class="group mt-8 inline-flex h-12 items-center gap-2 rounded-lg bg-gradient-to-r from-[#8b7deb] to-[#4c3fd4] px-7 text-sm font-semibold text-white shadow-[0_0_32px_rgba(139,125,235,0.4)] transition-all hover:shadow-[0_0_48px_rgba(139,125,235,0.6)]"
            @click="primaryCta"
          >
            {{ needsSetup ? 'Get started' : 'Sign in' }}
            <ArrowRight class="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
          </button>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-white/5">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 py-8 md:flex-row">
        <div class="flex items-center gap-2.5">
          <img src="/bloberry-icon.png" alt="Bloberry" class="h-7 w-7 rounded-md object-cover" />
          <span class="text-sm font-semibold">Bloberry</span>
        </div>
        <p class="text-xs text-white/40">Storage-agnostic object service — self-hosted, single binary.</p>
      </div>
    </footer>
  </div>
</template>
