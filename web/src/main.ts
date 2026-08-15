import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import { setUnauthenticatedHandler, setTokenChangedHandler } from './lib/api'

const app = createApp(App)
app.use(createPinia())
app.use(router)

const auth = useAuthStore()
auth.bootstrap()

// Keep the store in sync whenever the interceptor rotates tokens.
setTokenChangedHandler(() => auth.syncFromStorage())

// When a refresh attempt fails (expired/revoked session), clear local state
// and send the user to login — without this, a page refresh with a stale
// access token silently strands the user on a blank/unauthenticated dashboard.
setUnauthenticatedHandler(() => {
  auth.clear()
  const current = router.currentRoute.value
  if (current.name !== 'login') {
    router.push({ name: 'login', query: { next: current.fullPath } })
  }
})

// Rotate a stale/expired access token BEFORE the first render, so a page
// refresh or direct deep-link never flashes an unauthenticated state and never
// races the interceptor's own refresh.
auth.restoreSession().then(() => {
  app.mount('#app')
})

