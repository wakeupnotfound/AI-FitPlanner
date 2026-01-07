import { createRouter, createWebHistory } from 'vue-router'
import secureStorage from '@/utils/secureStorage'

/**
 * Vue Router Configuration
 * Defines application routes with authentication requirements
 */
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: (to) => {
        // Redirect to login if not authenticated, otherwise to dashboard
        const token = secureStorage.getAccessToken()
        return token ? '/dashboard' : '/login'
      }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { 
        requiresAuth: false,
        title: 'Login'
      }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue'),
      meta: { 
        requiresAuth: false,
        title: 'Register'
      }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { 
        requiresAuth: true,
        title: 'Dashboard'
      }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/ProfileView.vue'),
      meta: { 
        requiresAuth: true,
        title: 'Profile'
      }
    },
    {
      path: '/ai-config',
      name: 'ai-config',
      component: () => import('@/views/AIConfigView.vue'),
      meta: { 
        requiresAuth: true,
        title: 'AI Configuration'
      }
    },
    {
      path: '/training',
      name: 'training',
      component: () => import('@/views/TrainingView.vue'),
      meta: { 
        requiresAuth: true,
        title: 'Training'
      }
    },
    {
      path: '/nutrition',
      name: 'nutrition',
      component: () => import('@/views/NutritionView.vue'),
      meta: { 
        requiresAuth: true,
        title: 'Nutrition'
      }
    },
    {
      path: '/statistics',
      name: 'statistics',
      component: () => import('@/views/StatisticsView.vue'),
      meta: { 
        requiresAuth: true,
        title: 'Statistics'
      }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      redirect: '/dashboard'
    }
  ],
  scrollBehavior(to, from, savedPosition) {
    // Restore scroll position when using browser back/forward
    if (savedPosition) {
      return savedPosition
    }
    
    // Scroll to top for new routes
    if (to.path !== from.path) {
      return { top: 0, behavior: 'smooth' }
    }
    
    // Preserve scroll position for same route
    return false
  }
})

/**
 * Check if user is authenticated
 * @returns {boolean} True if user has valid access token
 */
function isAuthenticated() {
  const token = secureStorage.getAccessToken()
  return !!token
}

/**
 * Navigation guard for authentication
 * Redirects unauthenticated users to login page
 * Redirects authenticated users away from login/register pages
 */
router.beforeEach((to, from, next) => {
  const requiresAuth = to.meta.requiresAuth
  const authenticated = isAuthenticated()

  // Debug logging in development
  if (import.meta.env.DEV) {
    console.log(`[Router] Navigation: ${from.path} → ${to.path}`)
    console.log(`[Router] Requires auth: ${requiresAuth}, Authenticated: ${authenticated}`)
    console.log(`[Router] Token exists: ${!!secureStorage.getAccessToken()}`)
  }

  // Set page title
  if (to.meta.title) {
    document.title = `${to.meta.title} - ${import.meta.env.VITE_APP_TITLE || 'AI Fitness Planner'}`
  }

  // Route requires authentication
  if (requiresAuth && !authenticated) {
    if (import.meta.env.DEV) {
      console.log(`[Router] Redirecting to login: route requires auth but user not authenticated`)
    }
    // Redirect to login with return URL
    next({
      name: 'login',
      query: { redirect: to.fullPath }
    })
    return
  }

  // User is authenticated but trying to access login/register
  if (!requiresAuth && authenticated && (to.name === 'login' || to.name === 'register')) {
    if (import.meta.env.DEV) {
      console.log(`[Router] Redirecting to dashboard: user authenticated but accessing auth page`)
    }
    // Redirect to dashboard
    next({ name: 'dashboard' })
    return
  }

  // Allow navigation
  if (import.meta.env.DEV) {
    console.log(`[Router] Navigation allowed`)
  }
  next()
})

/**
 * Global after hook for navigation tracking
 */
router.afterEach((to, from) => {
  // Track page views (can be used for analytics)
  if (import.meta.env.DEV) {
    console.log(`[Router] Navigated from ${from.path} to ${to.path}`)
  }
})

export default router
