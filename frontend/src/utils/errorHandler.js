import { showToast, showDialog } from 'vant'

/**
 * Error Handler Utility
 * Classifies errors and provides user-friendly messages
 */
export class ErrorHandler {
  /**
   * Main error handling entry point
   * @param {Error} error - The error object
   * @param {Object} options - Additional options for error handling
   * @returns {Object} Processed error information
   */
  static handle(error, options = {}) {
    const { showMessage = true, logError = true } = options

    let errorInfo

    if (error.response) {
      // API error - server responded with error status
      errorInfo = this.handleAPIError(error.response)
    } else if (error.request) {
      // Network error - request was made but no response received
      errorInfo = this.handleNetworkError(error)
    } else {
      // Client error - error in request setup or other client-side issue
      errorInfo = this.handleClientError(error)
    }

    // Log error with sanitization
    if (logError) {
      this.logError(errorInfo)
    }

    // Show user-friendly message
    if (showMessage) {
      this.showErrorMessage(errorInfo)
    }

    return errorInfo
  }

  /**
   * Handle API errors (4xx, 5xx responses)
   * @param {Object} response - Axios response object
   * @returns {Object} Error information
   */
  static handleAPIError(response) {
    const { status, data } = response
    const errorCode = data?.code || status
    const serverMessage = data?.message || data?.error

    let userMessage = this.getAPIErrorMessage(status, serverMessage)
    let errorType = 'api'
    let severity = 'error'

    // Classify error severity
    if (status >= 500) {
      severity = 'critical'
    } else if (status === 429) {
      severity = 'warning'
    } else if (status >= 400 && status < 500) {
      severity = 'error'
    }

    return {
      type: errorType,
      severity,
      status,
      code: errorCode,
      message: userMessage,
      originalMessage: serverMessage,
      timestamp: new Date().toISOString()
    }
  }

  /**
   * Handle network errors (timeout, connection issues)
   * @param {Error} error - Axios error object
   * @returns {Object} Error information
   */
  static handleNetworkError(error) {
    let userMessage = 'Network error. Please check your connection.'
    let errorType = 'network'

    if (error.code === 'ECONNABORTED') {
      userMessage = 'Request timeout. Please try again.'
      errorType = 'timeout'
    } else if (error.message === 'Network Error') {
      userMessage = 'Unable to connect to server. Please check your internet connection.'
      errorType = 'connection'
    }

    return {
      type: errorType,
      severity: 'warning',
      code: error.code,
      message: userMessage,
      originalMessage: error.message,
      timestamp: new Date().toISOString()
    }
  }

  /**
   * Handle client-side errors
   * @param {Error} error - Error object
   * @returns {Object} Error information
   */
  static handleClientError(error) {
    return {
      type: 'client',
      severity: 'error',
      message: 'An unexpected error occurred. Please try again.',
      originalMessage: error.message,
      stack: error.stack,
      timestamp: new Date().toISOString()
    }
  }

  /**
   * Get user-friendly message for API errors
   * @param {number} status - HTTP status code
   * @param {string} serverMessage - Message from server
   * @returns {string} User-friendly error message
   */
  static getAPIErrorMessage(status, serverMessage) {
    const errorMessages = {
      400: serverMessage || 'Invalid request. Please check your input.',
      401: 'Authentication required. Please log in again.',
      403: 'Access denied. You do not have permission to perform this action.',
      404: 'Resource not found. The requested item does not exist.',
      409: serverMessage || 'Conflict. The resource already exists or is in use.',
      422: serverMessage || 'Validation failed. Please check your input.',
      429: 'Too many requests. Please wait a moment and try again.',
      500: 'Server error. Please try again later.',
      502: 'Bad gateway. The server is temporarily unavailable.',
      503: 'Service unavailable. Please try again later.',
      504: 'Gateway timeout. The server took too long to respond.'
    }

    return errorMessages[status] || serverMessage || 'An error occurred. Please try again.'
  }

  /**
   * Show error message to user
   * @param {Object} errorInfo - Error information object
   */
  static showErrorMessage(errorInfo) {
    const { severity, message } = errorInfo

    if (severity === 'critical') {
      // Show dialog for critical errors
      showDialog({
        title: 'Error',
        message: message,
        confirmButtonText: 'OK'
      })
    } else {
      // Show toast for other errors
      showToast({
        message: message,
        position: 'top',
        duration: 3000
      })
    }
  }

  /**
   * Log error with sensitive data sanitization
   * @param {Object} errorInfo - Error information object
   */
  static logError(errorInfo) {
    // Sanitize error info before logging
    const sanitizedError = this.sanitizeErrorData(errorInfo)

    // Log to console in development
    if (import.meta.env.DEV) {
      console.error('[Error Handler]', sanitizedError)
    }

    // In production, you would send to error tracking service
    // Example: Sentry, LogRocket, etc.
    // this.sendToErrorTracking(sanitizedError)
  }

  /**
   * Sanitize error data to remove sensitive information
   * @param {Object} errorInfo - Error information object
   * @returns {Object} Sanitized error information
   */
  static sanitizeErrorData(errorInfo) {
    const sanitized = { ...errorInfo }

    // Remove sensitive patterns from messages
    const sensitivePatterns = [
      /Bearer\s+[\w-]+\.[\w-]+\.[\w-]+/gi, // JWT tokens
      /api[_-]?key[:\s=]+[\w-]+/gi, // API keys
      /password[:\s=]+\S+/gi, // Passwords
      /token[:\s=]+[\w-]+/gi, // Generic tokens
      /authorization[:\s=]+\S+/gi, // Authorization headers
      /\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b/g // Email addresses (partial)
    ]

    // Sanitize message
    if (sanitized.message) {
      sensitivePatterns.forEach((pattern) => {
        sanitized.message = sanitized.message.replace(pattern, '[REDACTED]')
      })
    }

    // Sanitize original message
    if (sanitized.originalMessage) {
      sensitivePatterns.forEach((pattern) => {
        sanitized.originalMessage = sanitized.originalMessage.replace(pattern, '[REDACTED]')
      })
    }

    // Sanitize stack trace
    if (sanitized.stack) {
      sensitivePatterns.forEach((pattern) => {
        sanitized.stack = sanitized.stack.replace(pattern, '[REDACTED]')
      })
    }

    return sanitized
  }

  /**
   * Create a custom error with additional context
   * @param {string} message - Error message
   * @param {Object} context - Additional context
   * @returns {Error} Custom error object
   */
  static createError(message, context = {}) {
    const error = new Error(message)
    error.context = context
    return error
  }
}

/**
 * Convenience function to handle errors
 * @param {Error} error - The error to handle
 * @param {Object} options - Error handling options
 */
export function handleError(error, options = {}) {
  return ErrorHandler.handle(error, options)
}

export default ErrorHandler

