/**
 * Secure Storage Utility
 * Implements secure token storage best practices
 * Validates: Requirements 14.1
 */

/**
 * Storage keys with prefixes for organization
 */
const STORAGE_KEYS = {
  ACCESS_TOKEN: 'app_access_token',
  REFRESH_TOKEN: 'app_refresh_token',
  TOKEN_EXPIRY: 'app_token_expiry'
}

/**
 * Check if storage is available
 * @param {Storage} storage - Storage object (localStorage or sessionStorage)
 * @returns {boolean} True if storage is available
 */
function isStorageAvailable(storage) {
  try {
    const testKey = '__storage_test__'
    storage.setItem(testKey, 'test')
    storage.removeItem(testKey)
    return true
  } catch (e) {
    return false
  }
}

/**
 * Get storage object with fallback
 * Prefers sessionStorage for better security, falls back to localStorage
 * @returns {Storage} Storage object
 */
function getStorage() {
  // Try sessionStorage first (more secure - cleared on tab close)
  if (isStorageAvailable(sessionStorage)) {
    return sessionStorage
  }
  
  // Fallback to localStorage
  if (isStorageAvailable(localStorage)) {
    return localStorage
  }
  
  // If neither available, use in-memory storage
  return createMemoryStorage()
}

/**
 * Create in-memory storage fallback
 * @returns {Object} Memory storage object
 */
function createMemoryStorage() {
  const store = {}
  
  return {
    getItem: (key) => store[key] || null,
    setItem: (key, value) => {
      store[key] = value
    },
    removeItem: (key) => {
      delete store[key]
    },
    clear: () => {
      Object.keys(store).forEach(key => delete store[key])
    }
  }
}

/**
 * Simple XOR encryption for token obfuscation
 * Note: This is NOT cryptographically secure, but adds a layer of obfuscation
 * For production, consider using Web Crypto API or a proper encryption library
 * @param {string} text - Text to encrypt/decrypt
 * @param {string} key - Encryption key
 * @returns {string} Encrypted/decrypted text
 */
function xorEncrypt(text, key) {
  if (!text || !key) return text
  
  let result = ''
  for (let i = 0; i < text.length; i++) {
    result += String.fromCharCode(text.charCodeAt(i) ^ key.charCodeAt(i % key.length))
  }
  return btoa(result) // Base64 encode
}

/**
 * Decrypt XOR encrypted text
 * @param {string} encrypted - Encrypted text
 * @param {string} key - Encryption key
 * @returns {string} Decrypted text
 */
function xorDecrypt(encrypted, key) {
  if (!encrypted || !key) return encrypted
  
  try {
    const decoded = atob(encrypted) // Base64 decode
    let result = ''
    for (let i = 0; i < decoded.length; i++) {
      result += String.fromCharCode(decoded.charCodeAt(i) ^ key.charCodeAt(i % key.length))
    }
    return result
  } catch (e) {
    console.error('Decryption failed:', e)
    return null
  }
}

/**
 * Get encryption key from browser fingerprint
 * Creates a semi-unique key based on browser characteristics
 * @returns {string} Encryption key
 */
function getEncryptionKey() {
  // Create a fingerprint from browser characteristics
  const fingerprint = [
    navigator.userAgent,
    navigator.language,
    screen.width,
    screen.height,
    new Date().getTimezoneOffset()
  ].join('|')
  
  // Simple hash function
  let hash = 0
  for (let i = 0; i < fingerprint.length; i++) {
    const char = fingerprint.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash // Convert to 32-bit integer
  }
  
  return Math.abs(hash).toString(36)
}

/**
 * Secure Token Storage Class
 */
class SecureTokenStorage {
  constructor() {
    this.storage = getStorage()
    this.encryptionKey = getEncryptionKey()
    this.useEncryption = import.meta.env.PROD // Only encrypt in production
  }

  /**
   * Store access token securely
   * @param {string} token - Access token
   * @param {number} expiresIn - Token expiry time in seconds (optional)
   */
  setAccessToken(token, expiresIn = null) {
    if (!token) return
    
    const tokenToStore = this.useEncryption 
      ? xorEncrypt(token, this.encryptionKey)
      : token
    
    this.storage.setItem(STORAGE_KEYS.ACCESS_TOKEN, tokenToStore)
    
    // Store expiry time if provided
    if (expiresIn) {
      const expiryTime = Date.now() + (expiresIn * 1000)
      this.storage.setItem(STORAGE_KEYS.TOKEN_EXPIRY, expiryTime.toString())
    }
  }

  /**
   * Get access token
   * @returns {string|null} Access token or null if not found/expired
   */
  getAccessToken() {
    const stored = this.storage.getItem(STORAGE_KEYS.ACCESS_TOKEN)
    if (!stored) return null
    
    // Check if token is expired
    if (this.isTokenExpired()) {
      this.clearAccessToken()
      return null
    }
    
    return this.useEncryption 
      ? xorDecrypt(stored, this.encryptionKey)
      : stored
  }

  /**
   * Store refresh token securely
   * @param {string} token - Refresh token
   */
  setRefreshToken(token) {
    if (!token) return
    
    const tokenToStore = this.useEncryption 
      ? xorEncrypt(token, this.encryptionKey)
      : token
    
    this.storage.setItem(STORAGE_KEYS.REFRESH_TOKEN, tokenToStore)
  }

  /**
   * Get refresh token
   * @returns {string|null} Refresh token or null if not found
   */
  getRefreshToken() {
    const stored = this.storage.getItem(STORAGE_KEYS.REFRESH_TOKEN)
    if (!stored) return null
    
    return this.useEncryption 
      ? xorDecrypt(stored, this.encryptionKey)
      : stored
  }

  /**
   * Store both tokens
   * @param {string} accessToken - Access token
   * @param {string} refreshToken - Refresh token
   * @param {number} expiresIn - Token expiry time in seconds (optional)
   */
  setTokens(accessToken, refreshToken, expiresIn = null) {
    this.setAccessToken(accessToken, expiresIn)
    this.setRefreshToken(refreshToken)
  }

  /**
   * Clear access token
   */
  clearAccessToken() {
    this.storage.removeItem(STORAGE_KEYS.ACCESS_TOKEN)
    this.storage.removeItem(STORAGE_KEYS.TOKEN_EXPIRY)
  }

  /**
   * Clear refresh token
   */
  clearRefreshToken() {
    this.storage.removeItem(STORAGE_KEYS.REFRESH_TOKEN)
  }

  /**
   * Clear all tokens
   */
  clearTokens() {
    this.clearAccessToken()
    this.clearRefreshToken()
  }

  /**
   * Check if access token is expired
   * @returns {boolean} True if token is expired
   */
  isTokenExpired() {
    const expiryTime = this.storage.getItem(STORAGE_KEYS.TOKEN_EXPIRY)
    if (!expiryTime) return false
    
    return Date.now() >= parseInt(expiryTime)
  }

  /**
   * Get time until token expiry in seconds
   * @returns {number} Seconds until expiry, or 0 if expired/no expiry set
   */
  getTimeUntilExpiry() {
    const expiryTime = this.storage.getItem(STORAGE_KEYS.TOKEN_EXPIRY)
    if (!expiryTime) return 0
    
    const remaining = parseInt(expiryTime) - Date.now()
    return Math.max(0, Math.floor(remaining / 1000))
  }

  /**
   * Check if tokens exist
   * @returns {boolean} True if both tokens exist
   */
  hasTokens() {
    return !!(this.getAccessToken() && this.getRefreshToken())
  }
}

// Create and export singleton instance
export const secureStorage = new SecureTokenStorage()

export default secureStorage
