import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import { setUnauthenticatedHandler } from './lib/api'

const app = createApp(App)
app.use(createPinia())
app.use(router)

const auth = useAuthStore()
auth.bootstrap()

// Rotate a stale/expired access token up front so a page refresh never
// bounces through a 401 (the interceptor is the safety net for the rest).
auth.restoreSession()

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

app.mount('#app')
