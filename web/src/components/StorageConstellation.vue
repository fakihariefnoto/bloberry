<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import {
  Folder, FileText, Cloud, HardDrive, Database, Box, Server, ShieldCheck, FileCode, Globe, UploadCloud,
} from 'lucide-vue-next'

// A gentle wave of storage icons drifting in the background: icons are laid
// out in columns and bob along sine-wave curves, connected to their neighbours
// by faint lines, with a few data packets flowing along the wave. Everything
// stays behind the content (content renders above, z-10). Direct DOM writes
// per frame — no Vue reactivity in the animation loop. Respects reduced motion.

const wrap = ref<HTMLDivElement | null>(null)
const canvas = ref<HTMLCanvasElement | null>(null)
let raf = 0
let reduced = false

const ICONS = [Folder, FileText, Cloud, HardDrive, Database, Box, Server, ShieldCheck, FileCode, Globe, UploadCloud]

interface WaveNode {
  id: number
  icon: number
  col: number // column index
  band: number // which wave band
  x: number // current px
  y: number
  phase: number
  size: number
  baseAlpha: number
  el: HTMLElement | null
}

// wave geometry
const COLS = 26
const BANDS = 3
const COL_GAP = 92

const mouse = { x: -9999, y: -9999, active: false }
const follow = { x: -9999, y: -9999 } // smoothed mouse position (wave follows it)
const bump = { amp: 0 } // eased bump amplitude

const nodes: WaveNode[] = []
for (let band = 0; band < BANDS; band++) {
  for (let col = 0; col < COLS; col++) {
    const id = band * COLS + col
    nodes.push({
      id,
      icon: (id * 3 + band) % ICONS.length,
      col,
      band,
      x: 0, y: 0,
      phase: col * 0.5 + band * 2.1,
      size: 12 + ((id % 3) * 2),
      baseAlpha: 0.2 + 0.12 * Math.sin(band * 2 + col),
      el: null,
    })
  }
}

interface Packet { a: number; b: number; t: number; speed: number }
const packets: Packet[] = Array.from({ length: 18 }, () => ({
  a: Math.floor(Math.random() * nodes.length),
  b: Math.floor(Math.random() * nodes.length),
  t: Math.random(),
  speed: 0.003 + Math.random() * 0.004,
}))

function tick(now: number) {
  const c = canvas.value
  const wr = wrap.value
  if (!c || !wr) return
  const ctx = c.getContext('2d')
  if (!ctx) return

  // geometry is relative to the hero wrapper, not the viewport — the wave
  // stays inside its container and never bleeds over lower cards
  const w = wr.clientWidth || window.innerWidth
  const h = wr.clientHeight || window.innerHeight
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  if (c.width !== w * dpr) c.width = w * dpr
  if (c.height !== h * dpr) c.height = h * dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  // three quiet bands, dimmed so cards below stay readable
  const bandY = [h * 0.24, h * 0.44, h * 0.64]
  const spacing = COL_GAP
  const cols = Math.floor(w / spacing) + 2
  const t = reduced ? 0 : now / 6000

  // wave follows the mouse: a smooth bump that glides toward the cursor
  if (!reduced) {
    if (mouse.active) {
      const rect = wr.getBoundingClientRect()
      const mx = mouse.x - rect.left
      const my = mouse.y - rect.top
      follow.x += (mx - follow.x) * 0.05
      follow.y += (my - follow.y) * 0.05
      bump.amp += (1 - bump.amp) * 0.04
    } else {
      bump.amp += (0 - bump.amp) * 0.02
    }
  }

  // wave surface: two sine components crossing for a richer ripple
  for (const n of nodes) {
    const x = (n.col % cols) * spacing + (n.col % 2 === 0 ? -14 : 14) + spacing / 2
    const amp = 16 + (n.band % 2) * 8
    const wave = Math.sin(x / 90 + t + n.phase) * amp
    const wave2 = Math.sin(x / 47 - t * 0.8 + n.band) * 5
    let y = bandY[n.band] + wave + wave2

    // gaussian bump under the cursor — a soft crest, never a strong push
    if (mouse.active && !reduced) {
      const dx = x - follow.x
      const sigma = 200
      const g = Math.exp(-(dx * dx) / (2 * sigma * sigma))
      const pull = (follow.y - y) * 0.35
      y += bump.amp * (g * 42 + pull * g)
    }

    n.x = x
    n.y = y

    if (n.el) {
      n.el.style.transform = `translate(${x - n.size / 2}px, ${y - n.size / 2}px)`
      n.el.style.opacity = String(n.baseAlpha)
    }
  }

  // connect horizontal neighbours within the same band (a smooth wave web)
  ctx.lineWidth = 0.8
  for (const n of nodes) {
    const nxt = n.col + 1
    const other = nodes.find((m) => m.col === nxt && m.band === n.band)
    if (!other) continue
    const dx = other.x - n.x
    const dy = other.y - n.y
    const d = Math.hypot(dx, dy)
    const alpha = 0.12 * (1 - Math.min(d / (spacing * 1.6), 1))
    ctx.strokeStyle = `rgba(139,125,235,${alpha})`
    ctx.beginPath()
    ctx.moveTo(n.x, n.y)
    ctx.lineTo(other.x, other.y)
    ctx.stroke()
  }

  // vertical links between adjacent bands (occasional)
  for (let band = 0; band < BANDS - 1; band++) {
    for (const n of nodes) {
      if (n.band !== band || n.col % 4 !== 0) continue
      const below = nodes.find((m) => m.band === band + 1 && m.col === n.col)
      if (!below) continue
      ctx.strokeStyle = 'rgba(139,125,235,0.06)'
      ctx.beginPath()
      ctx.moveTo(n.x, n.y)
      ctx.lineTo(below.x, below.y)
      ctx.stroke()
    }
  }

  // data packets along the wave
  for (const p of packets) {
    const a = nodes[p.a]
    const b = nodes[p.b]
    const pxx = a.x + (b.x - a.x) * p.t
    const pyy = a.y + (b.y - a.y) * p.t
    const head = Math.min(p.t + p.speed * 12, 1)
    const hx = a.x + (b.x - a.x) * head
    const hy = a.y + (b.y - a.y) * head
    ctx.strokeStyle = 'rgba(139,125,235,0.3)'
    ctx.lineWidth = 1.1
    ctx.beginPath()
    ctx.moveTo(pxx, pyy)
    ctx.lineTo(hx, hy)
    ctx.stroke()
    ctx.fillStyle = 'rgba(196,181,253,0.75)'
    ctx.beginPath()
    ctx.arc(hx, hy, 2, 0, Math.PI * 2)
    ctx.fill()
    p.t += p.speed
    if (p.t > 1) {
      p.a = p.b
      p.b = Math.floor(Math.random() * nodes.length)
      p.t = 0
    }
  }

  if (!reduced) raf = requestAnimationFrame(tick)
}

function onResize() {
  if (raf) cancelAnimationFrame(raf)
  raf = requestAnimationFrame(tick)
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
  reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  wrap.value?.querySelectorAll<HTMLElement>('[data-node]').forEach((el, i) => {
    if (nodes[i]) nodes[i].el = el
  })
  // initial placement (no top-left flash)
  const w = window.innerWidth
  const spacing = COL_GAP
  const cols = Math.floor(w / spacing) + 2
  const bandY = [w * 0, window.innerHeight * 0.22, window.innerHeight * 0.4, window.innerHeight * 0.58]
  for (const n of nodes) {
    n.x = ((n.col % cols) * spacing) + spacing / 2
    n.y = bandY[n.band + 1] || window.innerHeight * 0.4
  }
  window.addEventListener('resize', onResize)
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseleave', onMouseLeave)
  raf = requestAnimationFrame(tick)
})

onBeforeUnmount(() => {
  cancelAnimationFrame(raf)
  window.removeEventListener('resize', onResize)
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseleave', onMouseLeave)
})
</script>

<template>
  <div ref="wrap" class="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
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
