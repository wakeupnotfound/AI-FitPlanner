import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { useAuthStore } from '../stores/auth'
import { authService, getAccessToken, getRefreshToken, hasTokens } from '../services/auth.service'
import { useI18n } from 'vue-i18n'

/**
 * useAuth Composable
 * Provides authentication functionality with automatic token refresh
 * Requirements: 1.4, 1.7
 */
export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()
  const { t } = useI18n()
  
  const loading = ref(false)
  const error = ref(null)

  /**
   * Login user with credentials
   * @param {Object} credentials - Login credentials
   * @param {string} credentials.username - Username
   * @param {string} credentials.password - Password
   * @returns {Promise<boolean>} True if login successful
   */
  const login = async (credentials) => {
    loading.value = true
    error.value = null

    try {
      const response = await authService.login(credentials)
      
      // Update store with user data and tokens
      if (response.data) {
        authStore.setUser(response.data.user)
        authStore.setTokens(response.data.access_token, response.data.refresh_token)
        
        // Debug logging
        if (import.meta.env.DEV) {
          console.log('[useAuth] Login successful, tokens stored')
          console.log('[useAuth] Token exists:', !!getAccessToken())
        }
      }

      // Small delay to ensure tokens are stored
      await new Promise(resolve => setTimeout(resolve, 100))

      // Navigate to dashboard or redirect URL
      const redirectUrl = router.currentRoute.value.query.redirect || '/dashboard'
      if (import.meta.env.DEV) {
        console.log('[useAuth] Navigating to:', redirectUrl)
      }
      await router.push(redirectUrl)
      
      return true
    } catch (err) {
      error.value = err.response?.data?.message || err.message || t('error.unknown')
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Register new user and auto-login
   * @param {Object} userData - Registration data
   * @returns {Promise<boolean>} True if registration successful
   */
  const register = async (userData) => {
    loading.value = true
    error.value = null

    try {
      // Register user
      await authService.register(userData)
      
      // Auto-login after successful registration
      await login({
        username: userData.username,
        password: userData.password
      })
      
      return true
    } catch (err) {
      error.value = err.response?.data?.message || err.message || t('error.unknown')
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Logout current user
   */
  const logout = async () => {
    loading.value = true

    try {
      await authService.logout()
      authStore.clearAuth()
      await router.push('/login')
    } catch (err) {
      // Still clear auth and redirect even if API fails
      authStore.clearAuth()
      await router.push('/login')
    } finally {
      loading.value = false
    }
  }

  /**
   * Refresh access token
   * @returns {Promise<boolean>} True if refresh successful
   */
  const refreshToken = async () => {
    try {
      const response = await authService.refreshToken()
      
      if (response.data) {
        authStore.setTokens(response.data.access_token, response.data.refresh_token)
      }
      
      return true
    } catch (err) {
      // If refresh fails, logout user
      await logout()
      return false
    }
  }

  /**
   * Initialize authentication state from localStorage
   * Called on app startup
   */
  const initAuth = () => {
    if (hasTokens()) {
      const accessToken = getAccessToken()
      const refreshToken = getRefreshToken()
      authStore.setTokens(accessToken, refreshToken)
      authStore.isAuthenticated = true
    }
  }

  /**
   * Check if current session is valid
   * Attempts token refresh if needed
   * @returns {Promise<boolean>} True if session is valid
   */
  const checkSession = async () => {
    if (!hasTokens()) {
      return false
    }

    // Try to validate session by refreshing token
    try {
      await refreshToken()
      return true
    } catch {
      return false
    }
  }

  // Computed properties
  const isAuthenticated = computed(() => authStore.isAuthenticated && hasTokens())
  const user = computed(() => authStore.user)
  const currentUser = computed(() => authStore.currentUser)

  return {
    // State
    loading,
    error,
    
    // Computed
    isAuthenticated,
    user,
    currentUser,
    
    // Methods
    login,
    register,
    logout,
    refreshToken,
    initAuth,
    checkSession
  }
}

export default useAuth
