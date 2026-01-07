import { defineStore } from 'pinia'
import apiClient from '../services/api'
import { storeTokens, clearTokens, getAccessToken, getRefreshToken } from '../services/auth.service'

/**
 * Auth Store
 * Manages authentication state, tokens, and user information
 */
export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false
  }),

  getters: {
    /**
     * Get current user
     */
    currentUser: (state) => state.user,

    /**
     * Check if user is authenticated
     */
    isLoggedIn: (state) => state.isAuthenticated && !!state.accessToken
  },

  actions: {
    /**
     * Login action
     * @param {Object} credentials - User credentials
     * @param {string} credentials.username - Username
     * @param {string} credentials.password - Password
     */
    async login(credentials) {
      try {
        const response = await apiClient.post('/auth/login', credentials)
        const { user, access_token, refresh_token } = response.data

        this.setUser(user)
        this.setTokens(access_token, refresh_token)

        return response
      } catch (error) {
        throw error
      }
    },

    /**
     * Logout action
     */
    async logout() {
      try {
        await apiClient.post('/auth/logout')
      } catch (error) {
        // Continue with logout even if API call fails
        console.error('Logout API call failed:', error)
      } finally {
        this.clearAuth()
      }
    },

    /**
     * Set user information
     * @param {Object} user - User object
     */
    setUser(user) {
      this.user = user
      this.isAuthenticated = true
    },

    /**
     * Set authentication tokens
     * @param {string} accessToken - Access token
     * @param {string} refreshToken - Refresh token
     */
    setTokens(accessToken, refreshToken) {
      this.accessToken = accessToken
      this.refreshToken = refreshToken
      
      // Store tokens using auth service (which now uses secureStorage)
      storeTokens(accessToken, refreshToken)
      
      // Debug logging
      if (import.meta.env.DEV) {
        console.log('[AuthStore] Tokens set in store and storage')
      }
    },

    /**
     * Clear authentication state
     */
    clearAuth() {
      this.user = null
      this.accessToken = null
      this.refreshToken = null
      this.isAuthenticated = false
      
      // Clear tokens from storage
      clearTokens()
    },

    /**
     * Initialize auth state from storage
     */
    initializeAuth() {
      const accessToken = getAccessToken()
      const refreshToken = getRefreshToken()

      if (accessToken && refreshToken) {
        this.accessToken = accessToken
        this.refreshToken = refreshToken
        this.isAuthenticated = true
        
        if (import.meta.env.DEV) {
          console.log('[AuthStore] Auth state initialized from storage')
        }
      }
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'auth',
        storage: localStorage,
        paths: ['user', 'isAuthenticated']
      }
    ]
  }
})
