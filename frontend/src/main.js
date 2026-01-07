import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import App from './App.vue'
import router from './router'
import i18n from './locales'
import { ErrorHandler } from './utils/errorHandler'
import { lazyLoadDirective } from './directives/lazyLoad'
import { sanitizeDirective } from './utils/sanitizer'
import { applyBrowserFixes, logBrowserInfo } from './utils/viewportFix'

// Import Vant components and styles
import Vant from 'vant'
import 'vant/lib/index.css'

// Import global styles
import './style.css'

// Apply browser-specific fixes on load
applyBrowserFixes()

// Log browser info in development
if (import.meta.env.DEV) {
  logBrowserInfo()
}

// Create Vue app
const app = createApp(App)

// Create Pinia store with persistence plugin
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

// Register global directives
app.directive('lazy', lazyLoadDirective)
app.directive('sanitize', sanitizeDirective)

// Use plugins
app.use(pinia)
app.use(router)
app.use(i18n)
app.use(Vant) // 注册Vant组件

// Initialize auth state from storage
import { useAuthStore } from './stores/auth'
const authStore = useAuthStore()
authStore.initializeAuth()

// Global error handler
app.config.errorHandler = (err, instance, info) => {
  console.error('Global error:', err, info)
  ErrorHandler.handle(err, { showMessage: true, logError: true })
}

// Mount app
app.mount('#app')

// Load security testing utilities in development
if (import.meta.env.DEV) {
  import('./utils/securityTest.js').then(module => {
    console.log('🔒 Security testing utilities loaded')
    console.log('Run window.securityTest.runAll() to test security measures')
  })
}

// Register service worker for PWA
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js').then(
      (registration) => {
        console.log('Service Worker registered:', registration.scope)
      },
      (error) => {
        console.error('Service Worker registration failed:', error)
      }
    )
  })
}
