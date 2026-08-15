<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowUpRight, Boxes, KeyRound, ShieldCheck, Server, Globe, HardDrive, UploadCloud } from 'lucide-vue-next'
import { api } from '../lib/api'

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
  <div class="min-h-screen bg-[#012624] text-[#edfffe]">
    <!-- Nav -->
    <header class="mx-auto flex max-w-[1440px] items-center justify-between px-8 py-6">
      <div class="flex items-center gap-2">
        <span class="flex h-8 w-8 items-center justify-center rounded-md bg-[#003734]">
          <Boxes class="h-4 w-4 text-[#cbfffc]" />
        </span>
        <span class="text-sm font-medium tracking-wide text-white">BLOBERRY</span>
      </div>
      <nav class="flex items-center gap-6 text-xs uppercase tracking-[0.12em] text-[#bbc7c6]">
        <a href="#features" class="transition-colors hover:text-white">Features</a>
        <a href="#how" class="transition-colors hover:text-white">How it works</a>
        <a href="#security" class="transition-colors hover:text-white">Security</a>
      </nav>
      <div class="flex items-center gap-3">
        <button
          class="rounded-md px-6 py-3 text-xs font-medium uppercase tracking-[0.08em] text-[#012624]"
          style="background: linear-gradient(90deg, rgb(0,130,124) 0%, rgb(203,255,252) 60%, rgb(250,209,255) 100%)"
          @click="primaryCta"
        >
          {{ checking ? 'Loading' : needsSetup ? 'Get started' : 'Sign in' }}
        </button>
      </div>
    </header>

    <!-- Hero -->
    <section class="relative mx-auto max-w-[1440px] px-8 pb-24 pt-20 text-center">
      <div class="mx-auto max-w-[900px]">
        <p class="text-xs font-medium uppercase tracking-[0.15em] text-[#bbc7c6]">Storage-agnostic object service</p>
        <h1 class="mt-6 text-[clamp(2.8rem,6vw,5.5rem)] font-medium leading-[1.0] tracking-[-0.04em] text-white">
          One API, one dashboard,
          <span class="text-[#fde9ff]">any storage</span>
        </h1>
        <p class="mx-auto mt-8 max-w-[620px] text-base leading-[1.6] text-[#bbc7c6]">
          Bloberry sits between your applications and S3, R2, OSS, GCS, Azure Blob or plain disk —
          giving you real folders, per-folder permissions, application keys and temporary links over any provider.
        </p>
        <div class="mt-10 flex items-center justify-center gap-4">
          <button
            class="rounded-md px-8 py-4 text-xs font-medium uppercase tracking-[0.08em] text-[#012624]"
            style="background: linear-gradient(90deg, rgb(0,130,124) 0%, rgb(203,255,252) 60%, rgb(250,209,255) 100%)"
            @click="primaryCta"
          >
            {{ needsSetup ? 'Set up Bloberry' : 'Sign in to dashboard' }}
          </button>
          <a href="#features" class="px-2 py-4 text-xs uppercase tracking-[0.12em] text-[#bbc7c6] hover:text-white">Explore</a>
        </div>
      </div>

      <!-- Biometric orb -->
      <div class="relative mx-auto mt-20 h-64 w-64">
        <div class="absolute inset-0 animate-pulse rounded-full opacity-40" style="background: radial-gradient(circle, rgba(0,130,124,0.8) 0%, transparent 65%)"></div>
        <div class="absolute inset-8 rounded-full border border-[#00827c]/40"></div>
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="grid grid-cols-3 gap-3 opacity-80">
            <span class="h-2 w-2 rounded-full bg-[#cbfffc]"></span>
            <span class="h-2 w-2 rounded-full bg-[#fde9ff]"></span>
            <span class="h-2 w-2 rounded-full bg-[#cbfffc]"></span>
            <span class="h-2 w-2 rounded-full bg-[#fde9ff]"></span>
            <span class="h-3 w-3 rounded-full bg-white/80"></span>
            <span class="h-2 w-2 rounded-full bg-[#cbfffc]"></span>
            <span class="h-2 w-2 rounded-full bg-[#cbfffc]"></span>
            <span class="h-2 w-2 rounded-full bg-[#fde9ff]"></span>
            <span class="h-2 w-2 rounded-full bg-[#cbfffc]"></span>
          </div>
        </div>
      </div>
    </section>

    <!-- Stats -->
    <section class="mx-auto max-w-[1440px] px-8">
      <div class="grid grid-cols-2 gap-8 md:grid-cols-4">
        <div class="text-center">
          <p class="text-[64px] font-medium leading-[1.0] tracking-[-0.03em] text-[#fde9ff]">6</p>
          <p class="mt-3 text-xs uppercase tracking-[0.055em] text-[#edfffe]">Storage drivers</p>
        </div>
        <div class="text-center">
          <p class="text-[64px] font-medium leading-[1.0] tracking-[-0.03em] text-[#fde9ff]">5GB</p>
          <p class="mt-3 text-xs uppercase tracking-[0.055em] text-[#edfffe]">Max object size</p>
        </div>
        <div class="text-center">
          <p class="text-[64px] font-medium leading-[1.0] tracking-[-0.03em] text-[#fde9ff]">100%</p>
          <p class="mt-3 text-xs uppercase tracking-[0.055em] text-[#edfffe]">Self-hosted</p>
        </div>
        <div class="text-center">
          <p class="text-[64px] font-medium leading-[1.0] tracking-[-0.03em] text-[#fde9ff]">1</p>
          <p class="mt-3 text-xs uppercase tracking-[0.055em] text-[#edfffe]">Binary to deploy</p>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section id="features" class="mx-auto max-w-[1440px] px-8 py-24">
      <p class="text-xs font-medium uppercase tracking-[0.15em] text-[#bbc7c6]">Features</p>
      <div class="mt-8 grid gap-6 md:grid-cols-3">
        <div class="rounded-[16px] bg-[#003734] p-9">
          <div class="flex items-start justify-between">
            <HardDrive class="h-6 w-6 text-[#cbfffc]" />
            <button class="flex h-8 w-8 items-center justify-center rounded-md bg-[rgba(3,81,75,0.5)]">
              <ArrowUpRight class="h-4 w-4 text-white" />
            </button>
          </div>
          <h3 class="mt-8 text-[24px] font-medium leading-[1.0] text-white">Any provider</h3>
          <p class="mt-4 text-[16px] leading-[1.5] text-[#bbc7c6]">S3, R2, OSS, GCS, Azure Blob or local disk behind one interface.</p>
        </div>
        <div class="rounded-[16px] bg-[#003734] p-9">
          <div class="flex items-start justify-between">
            <KeyRound class="h-6 w-6 text-[#cbfffc]" />
            <button class="flex h-8 w-8 items-center justify-center rounded-md bg-[rgba(3,81,75,0.5)]">
              <ArrowUpRight class="h-4 w-4 text-white" />
            </button>
          </div>
          <h3 class="mt-8 text-[24px] font-medium leading-[1.0] text-white">Real permissions</h3>
          <p class="mt-4 text-[16px] leading-[1.5] text-[#bbc7c6]">Roles, folder-level grants and scoped application keys.</p>
        </div>
        <div class="rounded-[16px] bg-[#003734] p-9">
          <div class="flex items-start justify-between">
            <UploadCloud class="h-6 w-6 text-[#cbfffc]" />
            <button class="flex h-8 w-8 items-center justify-center rounded-md bg-[rgba(3,81,75,0.5)]">
              <ArrowUpRight class="h-4 w-4 text-white" />
            </button>
          </div>
          <h3 class="mt-8 text-[24px] font-medium leading-[1.0] text-white">Bytes bypass you</h3>
          <p class="mt-4 text-[16px] leading-[1.5] text-[#bbc7c6]">Presigned uploads go straight to the provider. Your server never touches the payload.</p>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section id="how" class="mx-auto max-w-[1440px] px-8 pb-24">
      <p class="text-xs font-medium uppercase tracking-[0.15em] text-[#bbc7c6]">How it works</p>
      <div class="mt-8 grid gap-6 md:grid-cols-3">
        <div class="flex gap-5">
          <span class="text-[48px] font-medium leading-[1.0] text-[#fde9ff]">01</span>
          <div>
            <h4 class="text-[20px] font-medium text-white">Connect a backend</h4>
            <p class="mt-2 text-sm leading-[1.5] text-[#bbc7c6]">Register any provider — credentials are envelope-encrypted at rest.</p>
          </div>
        </div>
        <div class="flex gap-5">
          <span class="text-[48px] font-medium leading-[1.0] text-[#fde9ff]">02</span>
          <div>
            <h4 class="text-[20px] font-medium text-white">Give your team folders</h4>
            <p class="mt-2 text-sm leading-[1.5] text-[#bbc7c6]">Real folders with per-folder grants, roles and tenants.</p>
          </div>
        </div>
        <div class="flex gap-5">
          <span class="text-[48px] font-medium leading-[1.0] text-[#fde9ff]">03</span>
          <div>
            <h4 class="text-[20px] font-medium text-white">Share and integrate</h4>
            <p class="mt-2 text-sm leading-[1.5] text-[#bbc7c6]">Signed links, short URLs, SDKs and a CLI that script the same API.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Security -->
    <section id="security" class="mx-auto max-w-[1440px] px-8 pb-28">
      <div class="rounded-[16px] bg-[#011d1c] px-10 py-20 text-center">
        <ShieldCheck class="mx-auto h-8 w-8 text-[#cbfffc]" />
        <p class="mt-6 text-xs font-medium uppercase tracking-[0.15em] text-[#bbc7c6]">Security</p>
        <h2 class="mx-auto mt-4 max-w-[520px] text-[clamp(2rem,4vw,3.2rem)] font-medium leading-[1.0] tracking-[-0.04em] text-white">
          Credentials encrypted. Access audited.
        </h2>
        <div class="mx-auto mt-10 grid max-w-[720px] gap-6 text-left md:grid-cols-3">
          <div class="text-center">
            <Globe class="mx-auto h-5 w-5 text-[#fde9ff]" />
            <p class="mt-3 text-sm text-[#edfffe]">Envelope encryption at rest</p>
          </div>
          <div class="text-center">
            <KeyRound class="mx-auto h-5 w-5 text-[#fde9ff]" />
            <p class="mt-3 text-sm text-[#edfffe]">Two-factor authentication</p>
          </div>
          <div class="text-center">
            <Server class="mx-auto h-5 w-5 text-[#fde9ff]" />
            <p class="mt-3 text-sm text-[#edfffe]">Full audit log per tenant</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="mx-auto max-w-[1440px] px-8 pb-16">
      <div class="border-t border-[#707777]/30 pt-8 flex flex-col items-center justify-between gap-4 md:flex-row">
        <p class="text-xs uppercase tracking-[0.12em] text-[#707777]">Bloberry</p>
        <p class="text-xs text-[#707777]">Storage-agnostic object service — self-hosted</p>
      </div>
    </footer>
  </div>
</template>
