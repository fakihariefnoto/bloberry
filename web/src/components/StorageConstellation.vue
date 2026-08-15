<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import {
  Folder, FileText, HardDrive, KeyRound, Cloud, Database, Box, Server, ShieldCheck, UploadCloud, FileCode, FileArchive,
} from 'lucide-vue-next'

// Abstract network of storage glyphs spanning the whole viewport. Nodes
// drift lazily; when the mouse moves near one it is pushed away smoothly
// (spring physics) and eases back. Lines + flowing data packets on canvas.

const wrap = ref<HTMLDivElement | null>(null)
const canvas = ref<HTMLCanvasElement | null>(null)
let raf = 0
let reduced = false

const ICONS = [Folder, FileText, HardDrive, KeyRound, Cloud, Database, Box, Server, ShieldCheck, UploadCloud, FileCode, FileArchive]

interface IconNode {
  id: number
  icon: number
  // base position in fractions of viewport
  bx: number
  by: number
  // current pixel pos + velocity
  x: number
  y: number
  vx: number
  vy: number
  phase: number
  size: number
  baseOpacity: number
}

const COUNT = 30
const nodes: IconNode[] = Array.from({ length: COUNT }, (_, i) => {
  // even-ish scatter across the full viewport using a halton-like spread
  const bx = 0.03 + 0.94 * ((i * 0.618033988749895) % 1)
  const by = 0.06 + 0.88 * ((i * 0.381966011250105) % 1)
  return {
    id: i,
    icon: (i * 5 + 2) % ICONS.length,
    bx, by,
    x: 0, y: 0, vx: 0, vy: 0,
    phase: Math.random() * Math.PI * 2,
    size: 15 + (i % 4) * 3,
    baseOpacity: 0.45 + 0.25 * Math.random(),
  }
})

interface Packet {
  from: number
  to: number
  t: number
  speed: number
  glow: boolean
}

function makePackets(n: number): Packet[] {
  const out: Packet[] = []
  for (let i = 0; i < n; i++) {
    out.push({
      from: Math.floor(Math.random() * nodes.length),
      to: Math.floor(Math.random() * nodes.length),
      t: Math.random(),
      speed: 0.004 + Math.random() * 0.008,
      glow: Math.random() > 0.65,
    })
  }
  return out
}
const packets = makePackets(16)

// reactive positions the icon DOM binds to each frame
interface P { x: number; y: number; o: number }
const positions = ref<P[]>(nodes.map(() => ({ x: 0, y: 0, o: 0.6 })))

const mouse = { x: -9999, y: -9999, active: false }

const REPEL = 150 // push radius px
const STRENGTH = 0.9
const SPRING = 0.06 // ease back to base
const DAMP = 0.82

function tick(now: number) {
  if (!wrap.value || !canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return

  const w = window.innerWidth
  const h = window.innerHeight
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  if (canvas.value.width !== w * dpr) canvas.value.width = w * dpr
  if (canvas.value.height !== h * dpr) canvas.value.height = h * dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const drift = reduced ? 0 : (Math.sin((now / 2200) * 0.5) + 1) / 2

  // physics: spring home + mouse repulsion
  const px = positions.value
  for (const n of nodes) {
    const tx = n.bx * w
    const ty = n.by * h

    if (mouse.active && !reduced) {
      const dx = n.x - mouse.x
      const dy = n.y - mouse.y
      const d = Math.sqrt(dx * dx + dy * dy)
      if (d < REPEL && d > 0.01) {
        const f = (1 - d / REPEL) * STRENGTH
        n.vx += (dx / d) * f * 3
        n.vy += (dy / d) * f * 3
      }
    }

    // ease back toward base (with lazy drift)
    n.vx += (tx + n.bx * 14 * drift - n.x) * SPRING
    n.vy += (ty + n.by * 10 * drift - n.y) * SPRING
    n.vx *= DAMP
    n.vy *= DAMP
    n.x += n.vx
    n.y += n.vy

    const dpx = mouse.active ? Math.hypot(n.x - mouse.x, n.y - mouse.y) : 9999
    const near = mouse.active && dpx < REPEL ? 0.25 + (dpx / REPEL) * 0.35 : n.baseOpacity
    px[n.id] = { x: n.x, y: n.y, o: near }
  }

  const linkDist = Math.min(w, h) * 0.28

  // links
  for (let i = 0; i < nodes.length; i++) {
    for (let j = i + 1; j < nodes.length; j++) {
      const a = px[i]
      const b = px[j]
      const d2 = (a.x - b.x) ** 2 + (a.y - b.y) ** 2
      if (d2 < linkDist * linkDist) {
        const d = Math.sqrt(d2)
        const alpha = 0.13 * (1 - d / linkDist)
        ctx.strokeStyle = `rgba(139,125,235,${alpha})`
        ctx.lineWidth = 1
        ctx.beginPath()
        ctx.moveTo(a.x, a.y)
        ctx.lineTo(b.x, b.y)
        ctx.stroke()
      }
    }
  }

  // data packets
  for (const p of packets) {
    const a = px[p.from]
    const b = px[p.to]
    const d2 = (a.x - b.x) ** 2 + (a.y - b.y) ** 2
    if (d2 > linkDist * linkDist) continue
    const pxx = a.x + (b.x - a.x) * p.t
    const pyy = a.y + (b.y - a.y) * p.t
    const head = Math.min(p.t + p.speed * 14, 1)
    const hx = a.x + (b.x - a.x) * head
    const hy = a.y + (b.y - a.y) * head

    ctx.strokeStyle = p.glow ? 'rgba(139,125,235,0.55)' : 'rgba(167,139,250,0.4)'
    ctx.lineWidth = 1.4
    ctx.beginPath()
    ctx.moveTo(pxx, pyy)
    ctx.lineTo(hx, hy)
    ctx.stroke()

    const g = ctx.createRadialGradient(hx, hy, 0, hx, hy, 3.5)
    g.addColorStop(0, p.glow ? 'rgba(196,181,253,0.95)' : 'rgba(139,125,235,0.7)')
    g.addColorStop(1, 'rgba(139,125,235,0)')
    ctx.fillStyle = g
    ctx.beginPath()
    ctx.arc(hx, hy, 3.5, 0, Math.PI * 2)
    ctx.fill()

    p.t += p.speed
    if (p.t > 1) {
      p.from = p.to
      p.to = Math.floor(Math.random() * nodes.length)
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

onMounted(() => {
  const mq = window.matchMedia('(prefers-reduced-motion: reduce)')
  reduced = mq.matches
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseleave', onMouseLeave)
  raf = requestAnimationFrame(tick)
})

onBeforeUnmount(() => {
  cancelAnimationFrame(raf)
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseleave', onMouseLeave)
})
</script>

<template>
  <div ref="wrap" class="pointer-events-none fixed inset-0 overflow-hidden" aria-hidden="true">
    <canvas ref="canvas" class="absolute inset-0 h-full w-full" />
    <div
      v-for="n in nodes"
      :key="n.id"
      class="absolute text-[var(--color-primary)] transition-opacity duration-300"
      :style="{
        left: `${positions[n.id].x}px`,
        top: `${positions[n.id].y}px`,
        transform: 'translate(-50%, -50%)',
        opacity: positions[n.id].o,
      }"
    >
      <component :is="ICONS[n.icon]" :style="{ width: `${n.size}px`, height: `${n.size}px` }" stroke-width="1.6" />
    </div>
  </div>
</template>
