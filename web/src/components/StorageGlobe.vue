<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

// Storage globe — a rotating wireframe globe of storage nodes with data flows
// travelling between them. Pure canvas 2D + requestAnimationFrame. Respects
// OS "reduce motion" (style-guide §Motion): static render, no rotation.
const canvas = ref<HTMLCanvasElement | null>(null)
let raf = 0
let reduced = false

interface Node {
  x: number
  y: number
  z: number
  pulse: number
  phase: number
  ring: number // dot size bucket for depth
}

interface Flow {
  from: number
  to: number
  t: number
  speed: number
  color: string
}

function makeSphere(count: number): Node[] {
  const nodes: Node[] = []
  for (let i = 0; i < count; i++) {
    // fibonacci sphere for even distribution
    const y = 1 - (i / (count - 1)) * 2
    const r = Math.sqrt(1 - y * y)
    const theta = Math.PI * (3 - Math.sqrt(5)) * i
    nodes.push({
      x: Math.cos(theta) * r,
      y,
      z: Math.sin(theta) * r,
      pulse: 0,
      phase: Math.random() * Math.PI * 2,
      ring: Math.random(),
    })
  }
  return nodes
}

const NODES = 220
const FLOWS = 26
const nodes = makeSphere(NODES)

function buildFlows(): Flow[] {
  const flows: Flow[] = []
  for (let i = 0; i < FLOWS; i++) {
    flows.push({
      from: Math.floor(Math.random() * NODES),
      to: Math.floor(Math.random() * NODES),
      t: Math.random(),
      speed: 0.003 + Math.random() * 0.006,
      color: Math.random() > 0.6 ? '#8B7DEB' : '#6F5CE0',
    })
  }
  return flows
}
const flows = buildFlows()

// helper types for the canvas animation loop
type P3 = { x: number; y: number; z: number }

function project(p: P3, rot: number, w: number, h: number, radius: number): { x: number; y: number; z: number } {
  const cos = Math.cos(rot)
  const sin = Math.sin(rot)
  const x = p.x * cos - p.z * sin
  const z = p.x * sin + p.z * cos
  const y = p.y
  const scale = radius / (radius + z * radius * 0.6)
  return {
    x: w / 2 + x * radius * scale,
    y: h / 2 + y * radius * scale,
    z,
  }
}

let angle = 0
let running = true

function draw() {
  const c = canvas.value
  if (!c) return
  const ctx = c.getContext('2d')
  if (!ctx) return

  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  const w = c.clientWidth
  const h = c.clientHeight
  if (c.width !== w * dpr) c.width = w * dpr
  if (c.height !== h * dpr) c.height = h * dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)

  ctx.clearRect(0, 0, w, h)

  const radius = Math.min(w, h) * 0.32
  const rot = reduced ? 0 : angle

  // projected positions
  const pts = nodes.map((n) => project(n, rot, w, h, radius))

  // glow behind the sphere
  const glow = ctx.createRadialGradient(w / 2, h / 2, 0, w / 2, h / 2, radius * 1.6)
  glow.addColorStop(0, 'rgba(139,125,235,0.16)')
  glow.addColorStop(0.5, 'rgba(139,125,235,0.05)')
  glow.addColorStop(1, 'rgba(139,125,235,0)')
  ctx.fillStyle = glow
  ctx.beginPath()
  ctx.arc(w / 2, h / 2, radius * 1.6, 0, Math.PI * 2)
  ctx.fill()

  // back-face nodes (behind the sphere) — dimmer, drawn first
  ctx.fillStyle = 'rgba(139,125,235,0.25)'
  for (let i = 0; i < pts.length; i++) {
    const p = pts[i]
    if (p.z > 0) continue
    ctx.beginPath()
    ctx.arc(p.x, p.y, Math.max(1, (nodes[i].ring + 0.5) * 1.6), 0, Math.PI * 2)
    ctx.fill()
  }

  // connections between near front-face nodes
  ctx.lineWidth = 0.6
  for (let i = 0; i < pts.length; i++) {
    const a = pts[i]
    if (a.z > 0) continue
    for (let j = i + 1; j < pts.length; j++) {
      const b = pts[j]
      if (b.z > 0) continue
      const dx = a.x - b.x
      const dy = a.y - b.y
      const d2 = dx * dx + dy * dy
      if (d2 < radius * radius * 0.22) {
        const alpha = 0.10 * (1 - Math.sqrt(d2) / (radius * 0.47))
        ctx.strokeStyle = `rgba(139,125,235,${alpha})`
        ctx.beginPath()
        ctx.moveTo(a.x, a.y)
        ctx.lineTo(b.x, b.y)
        ctx.stroke()
      }
    }
  }

  // data flows travelling node→node
  for (const f of flows) {
    const a = pts[f.from]
    const b = pts[f.to]
    if (a.z > 0 || b.z > 0) continue
    const px = a.x + (b.x - a.x) * f.t
    const py = a.y + (b.y - a.y) * f.t
    const head = f.t + f.speed * 10
    const hx = a.x + (b.x - a.x) * Math.min(head, 1)
    const hy = a.y + (b.y - a.y) * Math.min(head, 1)

    ctx.strokeStyle = 'rgba(139,125,235,0.35)'
    ctx.lineWidth = 0.8
    ctx.beginPath()
    ctx.moveTo(px, py)
    ctx.lineTo(hx, hy)
    ctx.stroke()

    const comet = ctx.createRadialGradient(hx, hy, 0, hx, hy, 3)
    comet.addColorStop(0, 'rgba(196,181,253,0.95)')
    comet.addColorStop(1, 'rgba(196,181,253,0)')
    ctx.fillStyle = comet
    ctx.beginPath()
    ctx.arc(hx, hy, 3, 0, Math.PI * 2)
    ctx.fill()

    f.t += f.speed
    if (f.t > 1) {
      f.from = f.to
      f.to = Math.floor(Math.random() * NODES)
      f.t = 0
    }
  }

  // front-face nodes — brighter, with soft halo on some
  for (let i = 0; i < pts.length; i++) {
    const p = pts[i]
    const n = nodes[i]
    if (p.z <= 0) continue
    const size = Math.max(1.4, (n.ring + 0.6) * 2.2)
    const isPrimary = (i * 7) % 11 === 0
    ctx.fillStyle = isPrimary ? 'rgba(196,181,253,0.95)' : 'rgba(139,125,235,0.7)'
    ctx.beginPath()
    ctx.arc(p.x, p.y, size, 0, Math.PI * 2)
    ctx.fill()
    if (isPrimary) {
      ctx.fillStyle = 'rgba(139,125,235,0.25)'
      ctx.beginPath()
      ctx.arc(p.x, p.y, size * 2.4, 0, Math.PI * 2)
      ctx.fill()
    }
  }

  if (!reduced) {
    angle += 0.0018
    raf = requestAnimationFrame(draw)
  }
}

function resize() {
  if (running && !reduced) return
  draw()
}

onMounted(() => {
  const mq = window.matchMedia('(prefers-reduced-motion: reduce)')
  reduced = mq.matches
  draw()
  if (!reduced) {
    window.addEventListener('resize', resize)
  }
})

onBeforeUnmount(() => {
  running = false
  cancelAnimationFrame(raf)
})
</script>

<template>
  <canvas ref="canvas" class="h-full w-full" aria-hidden="true" />
</template>
