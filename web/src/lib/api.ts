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

export function setToken(t: string | null) {
  token = t
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

export async function request<T = unknown>(
  method: string,
  path: string,
  body?: unknown,
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
