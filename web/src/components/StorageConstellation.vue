<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

// Abstract orb in the background: hundreds of tiny particles distributed on a
// sphere with per-particle radius jitter (so it reads as an organic nebula,
// not a hard sphere), gently rotating. Near particles connect with faint
// lines and a few data packets flow between them. Canvas-only — nothing floats
// in front of the content. Subtle, low-opacity, honours prefers-reduced-motion.

const canvas = ref<HTMLCanvasElement | null>(null)
let raf = 0
let reduced = false

const COUNT = 320

interface P {
  // base unit-sphere coords
  sx: number; sy: number; sz: number
  // radial distortion (abstractness)
  rj: number // radius jitter 0..1
  spike: number // some particles get pushed outward
  phase: number
  size: number
  baseAlpha: number
}
const parts: P[] = Array.from({ length: COUNT }, (_, i) => {
  const y = 1 - (i / (COUNT - 1)) * 2
  const r = Math.sqrt(Math.max(0, 1 - y * y))
  const theta = Math.PI * (3 - Math.sqrt(5)) * i
  const spike = Math.random() > 0.86 ? 0.25 + Math.random() * 0.35 : 0
  return {
    sx: Math.cos(theta) * r,
    sy: y,
    sz: Math.sin(theta) * r,
    rj: (Math.random() - 0.5) * 0.14,
    spike,
    phase: Math.random() * Math.PI * 2,
    size: 1 + Math.random() * 1.6,
    baseAlpha: 0.25 + Math.random() * 0.4,
  }
})

interface Packet { from: number; to: number; t: number; speed: number }
const packets: Packet[] = Array.from({ length: 12 }, () => ({
  from: Math.floor(Math.random() * COUNT),
  to: Math.floor(Math.random() * COUNT),
  t: Math.random(),
  speed: 0.004 + Math.random() * 0.008,
}))

const mouse = { active: false, tx: 0, ty: 0 }
const tilt = { x: 0, y: 0 }
const SPRING = 0.05
let rotY = 0
const LINK_D2 = 0.012

function tick(now: number) {
  const c = canvas.value
  if (!c) return
  const ctx = c.getContext('2d')
  if (!ctx) return

  const w = window.innerWidth
  const h = window.innerHeight
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  if (c.width !== w * dpr) c.width = w * dpr
  if (c.height !== h * dpr) c.height = h * dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const cx = w * 0.5
  const cy = h * 0.42
  const radius = Math.min(w, h) * 0.30

  // gentle pulsing + mouse parallax tilt
  const pulse = reduced ? 1 : 1 + 0.03 * Math.sin(now / 1600)
  const tx = reduced ? 0 : mouse.active ? mouse.ty * 0.000018 : 0
  const ty = reduced ? 0 : mouse.active ? -mouse.tx * 0.000018 : 0
  tilt.x += (tx - tilt.x) * SPRING
  tilt.y += (ty - tilt.y) * SPRING
  if (!reduced) rotY += 0.0016

  const sinX = Math.sin(tilt.x)
  const cosX = Math.cos(tilt.x)
  const sinY = Math.sin(tilt.y)
  const cosY = Math.cos(tilt.y)
  const sinR = Math.sin(rotY)
  const cosR = Math.cos(rotY)

  const pts: { x: number; y: number; z: number }[] = new Array(COUNT)
  for (let i = 0; i < COUNT; i++) {
    const p = parts[i]
    // radial length: base 1 + jitter + occasional spike (abstract shape)
    const grow = reduced ? 1 : 1 + 0.04 * Math.sin(now / 2400 + p.phase)
    const rl = (1 + p.rj + p.spike * grow) * pulse

    let x = p.sx * rl * cosR - p.sz * rl * sinR
    let z = p.sx * rl * sinR + p.sz * rl * cosR
    let y = p.sy * rl
    const y2 = y * cosX - z * sinX
    const z2 = y * sinX + z * cosX
    const x2 = x * cosY + z2 * sinY
    const z3 = -x * sinY + z2 * cosY

    pts[i] = { x: cx + x2 * radius, y: cy + y2 * radius, z: z3 }
  }

  // faint connecting lines between nearby particles (back-face dim)
  ctx.lineWidth = 0.6
  for (let i = 0; i < COUNT; i++) {
    for (let j = i + 1; j < COUNT; j++) {
      const dx = parts[i].sx - parts[j].sx
      const dy = parts[i].sy - parts[j].sy
      const dz = parts[i].sz - parts[j].sz
      const d2 = dx * dx + dy * dy + dz * dz
      if (d2 > LINK_D2) continue
      if (pts[i].z < 0 && pts[j].z < 0) continue
      const alpha = 0.05 * (1 - Math.sqrt(d2) / Math.sqrt(LINK_D2))
      ctx.strokeStyle = `rgba(139,125,235,${alpha})`
      ctx.beginPath()
      ctx.moveTo(pts[i].x, pts[i].y)
      ctx.lineTo(pts[j].x, pts[j].y)
      ctx.stroke()
    }
  }

  // data packets
  for (const pk of packets) {
    const a = pts[pk.from]
    const b = pts[pk.to]
    const pxx = a.x + (b.x - a.x) * pk.t
    const pyy = a.y + (b.y - a.y) * pk.t
    const head = Math.min(pk.t + pk.speed * 14, 1)
    const hx = a.x + (b.x - a.x) * head
    const hy = a.y + (b.y - a.y) * head
    ctx.strokeStyle = 'rgba(139,125,235,0.28)'
    ctx.lineWidth = 1
    ctx.beginPath()
    ctx.moveTo(pxx, pyy)
    ctx.lineTo(hx, hy)
    ctx.stroke()
    ctx.fillStyle = 'rgba(196,181,253,0.7)'
    ctx.beginPath()
    ctx.arc(hx, hy, 1.8, 0, Math.PI * 2)
    ctx.fill()
    pk.t += pk.speed
    if (pk.t > 1) {
      pk.from = pk.to
      pk.to = Math.floor(Math.random() * COUNT)
      pk.t = 0
    }
  }

  // particles — front bright, back dim
  for (let i = 0; i < COUNT; i++) {
    const p = parts[i]
    const pt = pts[i]
    const front = pt.z > 0
    const depth = (pt.z + 1) / 2 // 0..1
    const alpha = p.baseAlpha * (front ? 0.4 + 0.6 * depth : 0.12 + 0.2 * depth)
    ctx.fillStyle = `rgba(139,125,235,${alpha})`
    ctx.beginPath()
    ctx.arc(pt.x, pt.y, p.size, 0, Math.PI * 2)
    ctx.fill()
  }

  if (!reduced) raf = requestAnimationFrame(tick)
}

function onMouseMove(e: MouseEvent) {
  mouse.tx = e.clientX - window.innerWidth / 2
  mouse.ty = e.clientY - window.innerHeight / 2
  mouse.active = true
}
function onMouseLeave() { mouse.active = false }
function onResize() {
  if (raf) cancelAnimationFrame(raf)
  raf = requestAnimationFrame(tick)
}

onMounted(() => {
  reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
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
  <canvas ref="canvas" class="pointer-events-none fixed inset-0 h-full w-full" aria-hidden="true" />
</template>
