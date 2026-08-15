<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import {
  Folder, FileText, HardDrive, KeyRound, Cloud, Database, Box, Server, ShieldCheck, UploadCloud, FileCode, FileArchive,
} from 'lucide-vue-next'

// A globe made of small storage glyphs: icons sit on the surface of a rotating
// 3D sphere (fibonacci distribution), back-face icons dim + shrink, links and
// data packets run between nearby surface nodes. Moving the mouse tilts the
// globe toward the cursor and eases back — smooth, never destructive.
// Icon positions are written straight to the DOM each frame (no Vue reactivity).

const wrap = ref<HTMLDivElement | null>(null)
const canvas = ref<HTMLCanvasElement | null>(null)
let raf = 0
let reduced = false

const ICONS = [Folder, FileText, HardDrive, KeyRound, Cloud, Database, Box, Server, ShieldCheck, UploadCloud, FileCode, FileArchive]

interface IconNode {
  id: number
  icon: number
  // unit sphere coords (fibonacci)
  sx: number
  sy: number
  sz: number
  size: number
  el: HTMLElement | null
}

// 90 glyphs on a fibonacci sphere → a dense, satellite-like globe
const COUNT = 90
const nodes: IconNode[] = Array.from({ length: COUNT }, (_, i) => {
  const y = 1 - (i / (COUNT - 1)) * 2
  const r = Math.sqrt(Math.max(0, 1 - y * y))
  const theta = Math.PI * (3 - Math.sqrt(5)) * i
  return {
    id: i,
    icon: (i * 5 + 2) % ICONS.length,
    sx: Math.cos(theta) * r,
    sy: y,
    sz: Math.sin(theta) * r,
    size: 11 + (i % 4) * 2.5,
    el: null,
  }
})

interface Packet { from: number; to: number; t: number; speed: number; glow: boolean }
const packets: Packet[] = Array.from({ length: 20 }, () => ({
  from: Math.floor(Math.random() * COUNT),
  to: Math.floor(Math.random() * COUNT),
  t: Math.random(),
  speed: 0.005 + Math.random() * 0.009,
  glow: Math.random() > 0.6,
}))

const mouse = { x: 0, y: 0, active: false, tx: 0, ty: 0 }
const tilt = { x: 0, y: 0 } // current globe tilt (radians)
let rotY = 0

const SPRING = 0.06
const LINK_ANGLE = 0.34 // max angular distance for a link

function tick() {
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

  // globe center sits in the hero zone
  const cx = w * 0.5
  const cy = h * 0.42
  const radius = Math.min(w, h) * 0.32

  // smooth mouse parallax → tilt targets, eased back
  const tx = reduced ? 0 : mouse.active ? mouse.ty * 0.00002 * Math.PI * 2 * 2 : 0
  const ty = reduced ? 0 : mouse.active ? -mouse.tx * 0.00002 * Math.PI * 2 * 2 : 0
  tilt.x += (tx - tilt.x) * SPRING
  tilt.y += (ty - tilt.y) * SPRING

  if (!reduced) rotY += 0.0022

  const sinX = Math.sin(tilt.x)
  const cosX = Math.cos(tilt.x)
  const sinY = Math.sin(tilt.y)
  const cosY = Math.cos(tilt.y)
  const sinR = Math.sin(rotY)
  const cosR = Math.cos(rotY)

  // projected positions (rotate Y then tilt X/Y, orthographic)
  const pts: { x: number; y: number; z: number }[] = []
  for (const n of nodes) {
    let x = n.sx * cosR - n.sz * sinR
    let z = n.sx * sinR + n.sz * cosR
    let y = n.sy
    // tilt around X
    let y2 = y * cosX - z * sinX
    let z2 = y * sinX + z * cosX
    // tilt around Y
    let x2 = x * cosY + z2 * sinY
    let z3 = -x * sinY + z2 * cosY

    const px = cx + x2 * radius
    const py = cy + y2 * radius
    pts.push({ x: px, y: py, z: z3 })
  }

  // depth sort: back face first
  const order = nodes.map((_n, i) => ({ i, z: pts[i].z })).sort((a, b) => a.z - b.z)

  // links between nearby surface nodes (3D angular distance)
  const linkSq = LINK_ANGLE * LINK_ANGLE
  ctx.lineWidth = 1
  for (let i = 0; i < COUNT; i++) {
    for (let j = i + 1; j < COUNT; j++) {
      const dx = nodes[i].sx - nodes[j].sx
      const dy = nodes[i].sy - nodes[j].sy
      const dz = nodes[i].sz - nodes[j].sz
      const d2 = dx * dx + dy * dy + dz * dz
      if (d2 > linkSq) continue
      // only draw links where at least one endpoint is on the front
      if (pts[i].z < 0 && pts[j].z < 0) continue
      const alpha = 0.10 * (1 - Math.sqrt(d2) / LINK_ANGLE)
      ctx.strokeStyle = `rgba(139,125,235,${alpha})`
      ctx.beginPath()
      ctx.moveTo(pts[i].x, pts[i].y)
      ctx.lineTo(pts[j].x, pts[j].y)
      ctx.stroke()
    }
  }

  // data packets along links
  for (const p of packets) {
    const a = pts[p.from]
    const b = pts[p.to]
    const d2 = (a.x - b.x) ** 2 + (a.y - b.y) ** 2
    if (d2 > (radius * LINK_ANGLE) ** 2 * 4) continue
    const pxx = a.x + (b.x - a.x) * p.t
    const pyy = a.y + (b.y - a.y) * p.t
    const head = Math.min(p.t + p.speed * 14, 1)
    const hx = a.x + (b.x - a.x) * head
    const hy = a.y + (b.y - a.y) * head

    ctx.strokeStyle = p.glow ? 'rgba(139,125,235,0.5)' : 'rgba(167,139,250,0.35)'
    ctx.lineWidth = 1.3
    ctx.beginPath()
    ctx.moveTo(pxx, pyy)
    ctx.lineTo(hx, hy)
    ctx.stroke()

    const g = ctx.createRadialGradient(hx, hy, 0, hx, hy, 3)
    g.addColorStop(0, p.glow ? 'rgba(196,181,253,0.9)' : 'rgba(139,125,235,0.65)')
    g.addColorStop(1, 'rgba(139,125,235,0)')
    ctx.fillStyle = g
    ctx.beginPath()
    ctx.arc(hx, hy, 3, 0, Math.PI * 2)
    ctx.fill()

    p.t += p.speed
    if (p.t > 1) {
      p.from = p.to
      p.to = Math.floor(Math.random() * COUNT)
      p.t = 0
    }
  }

  // icons: front face bright + full size, back face dim + small
  for (const o of order) {
    const nn = nodes[o.i]
    const p = pts[o.i]
    if (!nn.el) continue
    const front = p.z > 0
    const depth = front ? p.z : 0
    const scale = 0.55 + 0.45 * (p.z + 1) / 2
    const s = nn.size * (0.6 + 0.4 * ((p.z + 1) / 2))
    const op = front ? 0.5 + 0.5 * depth : 0.14 + 0.12 * ((p.z + 1) / 2)
    nn.el.style.transform = `translate(${p.x - s / 2}px, ${p.y - s / 2}px) scale(${scale})`
    nn.el.style.opacity = String(op)
  }

  if (!reduced) raf = requestAnimationFrame(tick)
}

function onMouseMove(e: MouseEvent) {
  mouse.tx = e.clientX - window.innerWidth / 2
  mouse.ty = e.clientY - window.innerHeight / 2
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
  wrap.value?.querySelectorAll<HTMLElement>('[data-node]').forEach((el, i) => {
    if (nodes[i]) nodes[i].el = el
  })
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
      :style="{ width: `${n.size}px`, height: `${n.size}px` }"
    >
      <component :is="ICONS[n.icon]" class="h-full w-full" stroke-width="1.6" />
    </div>
  </div>
</template>
