import axios from 'axios'
import { offlineQueue } from '@/utils/offlineQueue'
import secureStorage from '@/utils/secureStorage'

/**
 * API Client Service
 * Handles all HTTP requests with authentication, error handling, and token refresh
 */
class APIClient {
  constructor() {
    this.client = axios.create({
      baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:9999/api/v1',
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    this.isRefreshing = false
    this.failedQueue = []

    this.setupInterceptors()
  }

  /**
   * Set up request and response interceptors
   */
  setupInterceptors() {
    // Request interceptor: Add authentication token
    this.client.interceptors.request.use(
      (config) => {
        const token = this.getAccessToken()
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // Response interceptor: Handle errors and token refresh
    this.client.interceptors.response.use(
      (response) => {
        return response.data
      },
      async (error) => {
        const originalRequest = error.config

        // Handle 401 Unauthorized - Token expired
        if (error.response?.status === 401 && !originalRequest._retry) {
          if (this.isRefreshing) {
            // Queue the request while token is being refreshed
            return new Promise((resolve, reject) => {
              this.failedQueue.push({ resolve, reject })
            })
              .then((token) => {
                originalRequest.headers.Authorization = `Bearer ${token}`
                return this.client(originalRequest)
              })
              .catch((err) => {
                return Promise.reject(err)
              })
          }

          originalRequest._retry = true
          this.isRefreshing = true

          try {
            const newToken = await this.handleTokenRefresh()
            this.isRefreshing = false
            this.processQueue(null, newToken)
            originalRequest.headers.Authorization = `Bearer ${newToken}`
            return this.client(originalRequest)
          } catch (refreshError) {
            this.isRefreshing = false
            this.processQueue(refreshError, null)
            this.clearTokens()
            // Redirect to login will be handled by router guard
            window.location.href = '/login'
            return Promise.reject(refreshError)
          }
        }

        // Handle other errors
        return Promise.reject(error)
      }
    )
  }

  /**
   * Process queued requests after token refresh
   */
  processQueue(error, token = null) {
    this.failedQueue.forEach((promise) => {
      if (error) {
        promise.reject(error)
      } else {
        promise.resolve(token)
      }
    })
    this.failedQueue = []
  }

  /**
   * Handle token refresh
   */
  async handleTokenRefresh() {
    const refreshToken = this.getRefreshToken()
    if (!refreshToken) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await axios.post(
        `${import.meta.env.VITE_API_BASE_URL}/auth/refresh`,
        { refresh_token: refreshToken },
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      )

      const { access_token, refresh_token } = response.data.data
      this.storeTokens(access_token, refresh_token)
      return access_token
    } catch (error) {
      throw new Error('Token refresh failed')
    }
  }

  /**
   * Get access token from secure storage
   */
  getAccessToken() {
    return secureStorage.getAccessToken()
  }

  /**
   * Get refresh token from secure storage
   */
  getRefreshToken() {
    return secureStorage.getRefreshToken()
  }

  /**
   * Store tokens in secure storage
   */
  storeTokens(accessToken, refreshToken) {
    secureStorage.setTokens(accessToken, refreshToken)
  }

  /**
   * Clear tokens from secure storage
   */
  clearTokens() {
    secureStorage.clearTokens()
  }

  /**
   * Make GET request with retry logic
   */
  async get(url, config = {}) {
    return this.requestWithRetry(() => this.client.get(url, config))
  }

  /**
   * Make POST request with retry logic
   */
  async post(url, data, config = {}) {
    // If offline, queue the operation
    if (!navigator.onLine) {
      return this.queueOfflineOperation('POST', url, data, config)
    }
    return this.requestWithRetry(() => this.client.post(url, data, config))
  }

  /**
   * Make PUT request with retry logic
   */
  async put(url, data, config = {}) {
    // If offline, queue the operation
    if (!navigator.onLine) {
      return this.queueOfflineOperation('PUT', url, data, config)
    }
    return this.requestWithRetry(() => this.client.put(url, data, config))
  }

  /**
   * Make DELETE request with retry logic
   */
  async delete(url, config = {}) {
    // If offline, queue the operation
    if (!navigator.onLine) {
      return this.queueOfflineOperation('DELETE', url, null, config)
    }
    return this.requestWithRetry(() => this.client.delete(url, config))
  }

  /**
   * Make PATCH request with retry logic
   */
  async patch(url, data, config = {}) {
    // If offline, queue the operation
    if (!navigator.onLine) {
      return this.queueOfflineOperation('PATCH', url, data, config)
    }
    return this.requestWithRetry(() => this.client.patch(url, data, config))
  }

  /**
   * Queue operation for offline sync
   */
  queueOfflineOperation(type, url, data, config) {
    const fullUrl = url.startsWith('http') ? url : `${this.client.defaults.baseURL}${url}`
    
    const operationId = offlineQueue.enqueue({
      type,
      url: fullUrl,
      data,
      headers: {
        ...this.client.defaults.headers,
        ...config.headers,
        Authorization: `Bearer ${this.getAccessToken()}`
      }
    })

    // Return a promise that resolves with queued operation info
    return Promise.resolve({
      queued: true,
      operationId,
      message: 'Operation queued for sync when online'
    })
  }

  /**
   * Execute request with retry logic for network failures
   */
  async requestWithRetry(requestFn, retries = 2) {
    try {
      return await requestFn()
    } catch (error) {
      // Only retry on network errors, not on 4xx or 5xx responses
      if (retries > 0 && this.isNetworkError(error)) {
        await this.delay(1000) // Wait 1 second before retry
        return this.requestWithRetry(requestFn, retries - 1)
      }
      throw error
    }
  }

  /**
   * Check if error is a network error (not a response error)
   */
  isNetworkError(error) {
    return !error.response && (error.code === 'ECONNABORTED' || error.message === 'Network Error')
  }

  /**
   * Delay helper for retry logic
   */
  delay(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms))
  }
}

// Create and export singleton instance
const apiClient = new APIClient()

export default apiClient

