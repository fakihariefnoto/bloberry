<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import {
  Folder, FileText, HardDrive, KeyRound, Cloud, Database, Box, Server, ShieldCheck, UploadCloud, FileCode, FileArchive,
} from 'lucide-vue-next'

// Abstract network of storage glyphs spanning the whole viewport. Nodes drift
// lazily; moving the mouse near one pushes it away smoothly (spring physics)
// and it eases back. Links are degree-limited to keep the web clean, not a
// tangled blob. Icon positions are written directly to the DOM each frame
// (no Vue reactivity) so it stays at 60fps.

const wrap = ref<HTMLDivElement | null>(null)
const canvas = ref<HTMLCanvasElement | null>(null)
let raf = 0
let reduced = false

const ICONS = [Folder, FileText, HardDrive, KeyRound, Cloud, Database, Box, Server, ShieldCheck, UploadCloud, FileCode, FileArchive]

interface IconNode {
  id: number
  icon: number
  bx: number // base fraction (of viewport)
  by: number
  x: number // current px
  y: number
  vx: number
  vy: number
  phase: number
  size: number
  baseOpacity: number
  el: HTMLElement | null
}

const COUNT = 26
const nodes: IconNode[] = Array.from({ length: COUNT }, (_, i) => ({
  id: i,
  icon: (i * 5 + 2) % ICONS.length,
  bx: 0.04 + 0.92 * ((i * 0.618033988749895) % 1),
  by: 0.08 + 0.84 * ((i * 0.381966011250105) % 1),
  x: 0, y: 0, vx: 0, vy: 0,
  phase: Math.random() * Math.PI * 2,
  size: 15 + (i % 4) * 3,
  baseOpacity: 0.5 + 0.25 * Math.random(),
  el: null,
}))

interface Packet { from: number; to: number; t: number; speed: number; glow: boolean }
const packets: Packet[] = Array.from({ length: 14 }, () => ({
  from: Math.floor(Math.random() * COUNT),
  to: Math.floor(Math.random() * COUNT),
  t: Math.random(),
  speed: 0.004 + Math.random() * 0.008,
  glow: Math.random() > 0.65,
}))

const mouse = { x: -9999, y: -9999, active: false }
const REPEL = 140
const SPRING = 0.055
const DAMP = 0.84
const MAX_LINKS = 3 // per-node link cap → clean web

function tick(now: number) {
  if (!canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return

  const w = window.innerWidth
  const h = window.innerHeight
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  if (canvas.value.width !== w * dpr) canvas.value.width = w * dpr
  if (canvas.value.height !== h * dpr) canvas.value.height = h * dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const drift = reduced ? 0 : (Math.sin(now / 2600) + 1) / 2

  // physics + direct DOM writes (no Vue reactivity per frame)
  for (const n of nodes) {
    const tx = n.bx * w
    const ty = n.by * h

    if (mouse.active && !reduced) {
      const dx = n.x - mouse.x
      const dy = n.y - mouse.y
      const d = Math.hypot(dx, dy)
      if (d < REPEL && d > 0.01) {
        const f = (1 - d / REPEL) * 1.6
        n.vx += (dx / d) * f
        n.vy += (dy / d) * f
      }
    }

    n.vx += (tx + n.bx * 16 * drift - n.x) * SPRING
    n.vy += (ty + n.by * 12 * drift - n.y) * SPRING
    n.vx *= DAMP
    n.vy *= DAMP
    n.x += n.vx
    n.y += n.vy

    if (n.el) {
      n.el.style.transform = `translate(${n.x - n.size / 2}px, ${n.y - n.size / 2}px)`
      const dist = mouse.active ? Math.hypot(n.x - mouse.x, n.y - mouse.y) : REPEL + 1
      const o = dist < REPEL ? 0.3 + (dist / REPEL) * 0.45 : n.baseOpacity
      n.el.style.opacity = String(o)
    }
  }

  // degree-limited links
  const linkDist = Math.min(w, h) * 0.24
  const degree = new Array(COUNT).fill(0)
  for (let i = 0; i < COUNT; i++) {
    for (let j = i + 1; j < COUNT; j++) {
      if (degree[i] >= MAX_LINKS || degree[j] >= MAX_LINKS) continue
      const a = nodes[i]
      const b = nodes[j]
      const d2 = (a.x - b.x) ** 2 + (a.y - b.y) ** 2
      if (d2 > linkDist * linkDist) continue
      const d = Math.sqrt(d2)
      const alpha = 0.16 * (1 - d / linkDist)
      ctx.strokeStyle = `rgba(139,125,235,${alpha})`
      ctx.lineWidth = 1
      ctx.beginPath()
      ctx.moveTo(a.x, a.y)
      ctx.lineTo(b.x, b.y)
      ctx.stroke()
      degree[i]++
      degree[j]++
    }
  }

  // data packets
  for (const p of packets) {
    const a = nodes[p.from]
    const b = nodes[p.to]
    const d2 = (a.x - b.x) ** 2 + (a.y - b.y) ** 2
    if (d2 > linkDist * linkDist) continue
    const pxx = a.x + (b.x - a.x) * p.t
    const pyy = a.y + (b.y - a.y) * p.t
    const head = Math.min(p.t + p.speed * 14, 1)
    const hx = a.x + (b.x - a.x) * head
    const hy = a.y + (b.y - a.y) * head

    ctx.strokeStyle = p.glow ? 'rgba(139,125,235,0.5)' : 'rgba(167,139,250,0.35)'
    ctx.lineWidth = 1.4
    ctx.beginPath()
    ctx.moveTo(pxx, pyy)
    ctx.lineTo(hx, hy)
    ctx.stroke()

    const g = ctx.createRadialGradient(hx, hy, 0, hx, hy, 3.5)
    g.addColorStop(0, p.glow ? 'rgba(196,181,253,0.9)' : 'rgba(139,125,235,0.65)')
    g.addColorStop(1, 'rgba(139,125,235,0)')
    ctx.fillStyle = g
    ctx.beginPath()
    ctx.arc(hx, hy, 3.5, 0, Math.PI * 2)
    ctx.fill()

    p.t += p.speed
    if (p.t > 1) {
      p.from = p.to
      p.to = Math.floor(Math.random() * COUNT)
      p.t = 0
    }
  }

  if (!reduced) raf = requestAnimationFrame(tick)
}

function onMouseMove(e: MouseEvent) {
  mouse.x = e.clientX
  mouse.y = e.clientY
  mouse.active = true
}
function onMouseLeave() {
  mouse.active = false
}

function onResize() {
  if (raf) cancelAnimationFrame(raf)
  raf = requestAnimationFrame(tick)
}

onMounted(() => {
  const mq = window.matchMedia('(prefers-reduced-motion: reduce)')
  reduced = mq.matches
  // collect icon elements once
  wrap.value?.querySelectorAll<HTMLElement>('[data-node]').forEach((el, i) => {
    if (nodes[i]) nodes[i].el = el
  })
  // set initial positions (avoid top-left flash)
  for (const n of nodes) {
    n.x = n.bx * window.innerWidth
    n.y = n.by * window.innerHeight
  }
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseleave', onMouseLeave)
  window.addEventListener('resize', onResize)
  raf = requestAnimationFrame(tick)
})

onBeforeUnmount(() => {
  cancelAnimationFrame(raf)
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseleave', onMouseLeave)
  window.removeEventListener('resize', onResize)
})
</script>

<template>
  <div ref="wrap" class="pointer-events-none fixed inset-0 overflow-hidden" aria-hidden="true">
    <canvas ref="canvas" class="absolute inset-0 h-full w-full" />
    <div
      v-for="n in nodes"
      :key="n.id"
      data-node
      class="absolute left-0 top-0 will-change-transform text-[var(--color-primary)]"
      :style="{
        width: `${n.size}px`,
        height: `${n.size}px`,
      }"
    >
      <component :is="ICONS[n.icon]" class="h-full w-full" stroke-width="1.6" />
    </div>
  </div>
</template>
