import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/', name: 'landing', component: () => import('../pages/LandingPage.vue'), meta: { public: true } },
  { path: '/setup', name: 'setup', component: () => import('../pages/SetupPage.vue'), meta: { public: true } },
  { path: '/login', name: 'login', component: () => import('../pages/LoginPage.vue'), meta: { public: true } },
  { path: '/login/otp', name: 'otp-login', component: () => import('../pages/OtpLoginPage.vue'), meta: { public: true } },
  { path: '/activate', name: 'activate', component: () => import('../pages/ActivatePage.vue'), meta: { public: true } },
  { path: '/forgot-password', name: 'forgot-password', component: () => import('../pages/ForgotPasswordPage.vue'), meta: { public: true } },
  { path: '/reset-password', name: 'reset-password', component: () => import('../pages/ResetPasswordPage.vue'), meta: { public: true } },
  { path: '/invite/:token', name: 'accept-invitation', component: () => import('../pages/AcceptInvitationPage.vue'), meta: { public: true } },
  { path: '/s/:slug', name: 'link-expired', component: () => import('../pages/LinkExpiredPage.vue'), meta: { public: true } },

  {
    path: '/app',
    component: () => import('../components/AppShell.vue'),
    children: [
      { path: '', redirect: { name: 'files' } },
      { path: 'files/:folderId?', name: 'files', component: () => import('../pages/FilesPage.vue') },
      { path: 'files/detail/:fileId', name: 'file-detail', component: () => import('../pages/FileDetailPage.vue') },
      { path: 'shares', name: 'shares', component: () => import('../pages/SharesPage.vue') },
      { path: 'jobs', name: 'jobs', component: () => import('../pages/JobsPage.vue') },
      { path: 'transfers', name: 'transfers', component: () => import('../pages/TransfersPage.vue') },
      { path: 'applications', name: 'applications', component: () => import('../pages/ApplicationsPage.vue'), meta: { admin: true } },
      { path: 'api-keys', name: 'api-keys', component: () => import('../pages/ApiKeysPage.vue') },
      { path: 'applications/:appId', name: 'application-detail', component: () => import('../pages/ApplicationDetailPage.vue'), meta: { admin: true } },
      { path: 'members', name: 'members', component: () => import('../pages/MembersPage.vue'), meta: { admin: true } },
      { path: 'audit', name: 'audit', component: () => import('../pages/AuditPage.vue'), meta: { admin: true } },
      { path: 'usage', name: 'usage', component: () => import('../pages/UsagePage.vue'), meta: { admin: true } },
      { path: 'settings/tenant', name: 'tenant-settings', component: () => import('../pages/TenantSettingsPage.vue'), meta: { admin: true } },
      { path: 'profile', name: 'profile', component: () => import('../pages/ProfilePage.vue') },
      { path: 'settings/account', name: 'account-settings', component: () => import('../pages/AccountSettingsPage.vue') },
      { path: 'settings/pair', name: 'pair-device', component: () => import('../pages/PairDevicePage.vue') },
      { path: 'admin/tenants', name: 'admin-tenants', component: () => import('../pages/AdminTenantsPage.vue'), meta: { platformAdmin: true } },
      { path: 'admin/tenants/:tenantId', name: 'admin-tenant-detail', component: () => import('../pages/AdminTenantDetailPage.vue'), meta: { platformAdmin: true } },
      { path: 'admin/backends', name: 'admin-backends', component: () => import('../pages/AdminBackendsPage.vue'), meta: { platformAdmin: true } },
      { path: 'admin/backends/:backendId', name: 'admin-backend-detail', component: () => import('../pages/AdminBackendDetailPage.vue'), meta: { platformAdmin: true } },
      { path: 'admin/usage', name: 'admin-usage', component: () => import('../pages/AdminUsagePage.vue'), meta: { platformAdmin: true } },
      { path: 'suspended', name: 'tenant-suspended', component: () => import('../pages/TenantSuspendedPage.vue') },
      { path: '403', name: 'forbidden', component: () => import('../pages/ForbiddenPage.vue') },
    ],
  },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('../pages/NotFoundPage.vue'), meta: { public: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  // Already signed in? Never show the landing — go straight to the dashboard.
  if (to.name === 'landing' && auth.isAuthenticated) {
    return { name: 'files' }
  }
  if (!to.meta.public && !auth.isAuthenticated) {
    return { name: 'login', query: { next: to.fullPath } }
  }
  if (to.meta.platformAdmin && !auth.user?.platform_role) {
    return { name: 'forbidden' }
  }
  return true
})

export default router
