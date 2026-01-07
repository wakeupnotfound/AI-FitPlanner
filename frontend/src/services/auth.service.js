import apiClient from './api'
import secureStorage from '@/utils/secureStorage'

/**
 * Authentication Service
 * Handles user authentication operations: register, login, logout, token refresh
 * Requirements: 1.1, 1.2, 1.3, 1.4
 */
export const authService = {
  /**
   * Register a new user
   * @param {Object} userData - User registration data
   * @param {string} userData.username - Username
   * @param {string} userData.email - Email address
   * @param {string} userData.password - Password
   * @param {string} userData.nickname - Optional nickname
   * @returns {Promise<Object>} Registration response with user data
   * Requirements: 1.2
   */
  async register(userData) {
    const response = await apiClient.post('/auth/register', userData)
    return response
  },

  /**
   * Login user with credentials
   * @param {Object} credentials - Login credentials
   * @param {string} credentials.username - Username
   * @param {string} credentials.password - Password
   * @returns {Promise<Object>} Login response with tokens and user data
   * Requirements: 1.1, 1.4
   */
  async login(credentials) {
    const response = await apiClient.post('/auth/login', credentials)
    
    // Store tokens if login successful
    if (response.data?.access_token && response.data?.refresh_token) {
      secureStorage.setTokens(response.data.access_token, response.data.refresh_token, response.data.expires_in)
    }
    
    return response
  },

  /**
   * Logout current user
   * Clears tokens regardless of API call success
   * @returns {Promise<void>}
   * Requirements: 1.1
   */
  async logout() {
    try {
      await apiClient.post('/auth/logout')
    } catch (error) {
      // Continue with logout even if API call fails
      console.warn('Logout API call failed:', error.message)
    } finally {
      secureStorage.clearTokens()
    }
  },

  /**
   * Refresh access token using refresh token
   * @returns {Promise<Object>} New tokens
   * Requirements: 1.4
   */
  async refreshToken() {
    const refreshToken = secureStorage.getRefreshToken()
    
    if (!refreshToken) {
      throw new Error('No refresh token available')
    }

    const response = await apiClient.post('/auth/refresh', {
      refresh_token: refreshToken
    })

    // Store new tokens
    if (response.data?.access_token && response.data?.refresh_token) {
      secureStorage.setTokens(response.data.access_token, response.data.refresh_token, response.data.expires_in)
    }

    return response
  },

  /**
   * Check if user is currently authenticated
   * @returns {boolean} True if access token exists
   */
  isAuthenticated() {
    return secureStorage.hasTokens()
  }
}

// Token storage utilities (using secureStorage)

/**
 * Get access token from secure storage
 * @returns {string|null} Access token or null
 */
export function getAccessToken() {
  return secureStorage.getAccessToken()
}

/**
 * Get refresh token from secure storage
 * @returns {string|null} Refresh token or null
 */
export function getRefreshToken() {
  return secureStorage.getRefreshToken()
}

/**
 * Store tokens in secure storage
 * @param {string} accessToken - Access token
 * @param {string} refreshToken - Refresh token
 * @param {number} expiresIn - Token expiry time in seconds
 */
export function storeTokens(accessToken, refreshToken, expiresIn = null) {
  secureStorage.setTokens(accessToken, refreshToken, expiresIn)
}

/**
 * Clear all tokens from secure storage
 */
export function clearTokens() {
  secureStorage.clearTokens()
}

/**
 * Check if tokens exist in secure storage
 * @returns {boolean} True if both tokens exist
 */
export function hasTokens() {
  return secureStorage.hasTokens()
}

export default authService
