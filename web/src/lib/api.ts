// API client — talks to the Go backend over the standard envelope
// {data?, messages?: [{code, content}]}. snake_case on the wire; the API
// surface mirrors api/openapi.yaml (ADR-11).

export interface Message {
  code: string
  content?: string
}

export interface Envelope<T = unknown> {
  data?: T
  messages?: Message[]
}

export class ApiError extends Error {
  status: number
  code: string
  content?: string

  constructor(status: number, code: string, content?: string) {
    super(content || code)
    this.status = status
    this.code = code
    this.content = content
  }
}

let baseURL = import.meta.env.VITE_API_BASE || ''
let token: string | null = null

const ACCESS_KEY = 'bloberry.access'
const REFRESH_KEY = 'bloberry.refresh'

export function setToken(t: string | null) {
  token = t
  if (t === null) {
    localStorage.removeItem(ACCESS_KEY)
    localStorage.removeItem(REFRESH_KEY)
  } else {
    localStorage.setItem(ACCESS_KEY, t)
  }
}

// Runtime override for the Wails desktop shell, which serves the same build
// from a native origin and must point at its configured server (05-core-infra).
export function setApiBase(url: string) {
  baseURL = url
}

function apiBase(): string {
  // Desktop: a global hook set by the Wails host.
  const w = window as unknown as { __bloberry_api_base?: string }
  if (w.__bloberry_api_base) return w.__bloberry_api_base
  return baseURL
}

// onUnauthenticated is invoked when a refresh attempt fails (expired or
// revoked session) so the app can clear state and redirect to login.
let onUnauthenticated: (() => void) | null = null
export function setUnauthenticatedHandler(fn: () => void) {
  onUnauthenticated = fn
}

// Single-flight refresh: concurrent 401s share one refresh promise.
let refreshPromise: Promise<boolean> | null = null

async function tryRefresh(): Promise<boolean> {
  if (refreshPromise) return refreshPromise
  refreshPromise = (async () => {
    const refreshToken = localStorage.getItem(REFRESH_KEY)
    if (!refreshToken) return false
    try {
      const res = await fetch(`${apiBase()}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      })
      if (!res.ok) return false
      const env = (await res.json()) as Envelope<{ access_token: string; refresh_token: string }>
      const data = env.data
      if (!data?.access_token) return false
      token = data.access_token
      localStorage.setItem(ACCESS_KEY, data.access_token)
      if (data.refresh_token) {
        localStorage.setItem(REFRESH_KEY, data.refresh_token)
      }
      return true
    } catch {
      return false
    } finally {
      refreshPromise = null
    }
  })()
  return refreshPromise
}

export async function request<T = unknown>(
  method: string,
  path: string,
  body?: unknown,
  _retried = false,
): Promise<T> {
  const headers: Record<string, string> = {}
  if (body !== undefined) headers['Content-Type'] = 'application/json'
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${apiBase()}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })

  if (res.status === 204) return undefined as T

  let env: Envelope<T> | null = null
  try {
    env = (await res.json()) as Envelope<T>
  } catch {
    env = null
  }

  // 401 + a refresh token available → rotate and retry once.
  if (res.status === 401 && !_retried && localStorage.getItem(REFRESH_KEY)) {
    const ok = await tryRefresh()
    if (ok) {
      return request<T>(method, path, body, true)
    }
    onUnauthenticated?.()
  }

  if (!res.ok) {
    const msg = env?.messages?.[0]
    throw new ApiError(res.status, msg?.code || 'unknown_error', msg?.content)
  }
  return env?.data as T
}

export const api = {
  get: <T>(path: string) => request<T>('GET', path),
  post: <T>(path: string, body?: unknown) => request<T>('POST', path, body),
  patch: <T>(path: string, body?: unknown) => request<T>('PATCH', path, body),
  delete: <T>(path: string) => request<T>('DELETE', path),
}

// downloadUrl resolves a protected download endpoint to the actual URL the
// browser can navigate to (a presigned/raw URL that is self-authorizing), so
// plain anchor navigation — which carries no Authorization header — works.
export async function downloadUrl(path: string): Promise<string> {
  const res = await fetch(`${apiBase()}${path}`, {
    method: 'GET',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    redirect: 'manual',
  })
  if (res.status >= 300 && res.status < 400) {
    const loc = res.headers.get('Location')
    if (loc) return loc.startsWith('http') ? loc : `${apiBase()}${loc}`
  }
  if (res.ok) {
    // proxy path (no redirect) — return the endpoint itself; caller opens it
    return `${apiBase()}${path}`
  }
  throw new ApiError(res.status, 'download_failed', `Download failed (${res.status})`)
}
