/**
 * Input Sanitization Utility
 * Provides functions to sanitize user input to prevent XSS attacks
 * Validates: Requirements 14.4
 */

/**
 * Escape HTML special characters to prevent XSS
 * @param {string} str - Input string to sanitize
 * @returns {string} Sanitized string
 */
export function escapeHtml(str) {
  if (typeof str !== 'string') {
    return str
  }

  const htmlEscapeMap = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#x27;',
    '/': '&#x2F;'
  }

  return str.replace(/[&<>"'/]/g, (char) => htmlEscapeMap[char])
}

/**
 * Remove all HTML tags from string
 * @param {string} str - Input string
 * @returns {string} String without HTML tags
 */
export function stripHtmlTags(str) {
  if (typeof str !== 'string') {
    return str
  }

  return str.replace(/<[^>]*>/g, '')
}

/**
 * Remove script tags and event handlers
 * @param {string} str - Input string
 * @returns {string} Sanitized string
 */
export function removeScripts(str) {
  if (typeof str !== 'string') {
    return str
  }

  // Remove script tags and their content
  let sanitized = str.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')

  // Remove event handlers (onclick, onerror, etc.)
  sanitized = sanitized.replace(/\s*on\w+\s*=\s*["'][^"']*["']/gi, '')
  sanitized = sanitized.replace(/\s*on\w+\s*=\s*[^\s>]*/gi, '')

  // Remove javascript: protocol
  sanitized = sanitized.replace(/javascript:/gi, '')

  // Remove data: protocol (can be used for XSS)
  sanitized = sanitized.replace(/data:text\/html/gi, '')

  return sanitized
}

/**
 * Sanitize user input for safe display
 * This is the main function to use for general user input
 * @param {string} input - User input string
 * @param {Object} options - Sanitization options
 * @param {boolean} options.allowHtml - Allow safe HTML tags (default: false)
 * @param {boolean} options.stripTags - Strip all HTML tags (default: true)
 * @returns {string} Sanitized input
 */
export function sanitizeInput(input, options = {}) {
  if (typeof input !== 'string') {
    return input
  }

  const {
    allowHtml = false,
    stripTags = true
  } = options

  let sanitized = input

  // Always remove scripts and dangerous content
  sanitized = removeScripts(sanitized)

  if (stripTags && !allowHtml) {
    // Remove all HTML tags
    sanitized = stripHtmlTags(sanitized)
  } else if (!allowHtml) {
    // Escape HTML characters
    sanitized = escapeHtml(sanitized)
  }

  // Trim whitespace
  sanitized = sanitized.trim()

  return sanitized
}

/**
 * Sanitize object properties recursively
 * @param {Object} obj - Object to sanitize
 * @param {Array<string>} excludeKeys - Keys to exclude from sanitization
 * @returns {Object} Sanitized object
 */
export function sanitizeObject(obj, excludeKeys = []) {
  if (typeof obj !== 'object' || obj === null) {
    return obj
  }

  if (Array.isArray(obj)) {
    return obj.map((item) => sanitizeObject(item, excludeKeys))
  }

  const sanitized = {}

  for (const [key, value] of Object.entries(obj)) {
    if (excludeKeys.includes(key)) {
      sanitized[key] = value
    } else if (typeof value === 'string') {
      sanitized[key] = sanitizeInput(value)
    } else if (typeof value === 'object' && value !== null) {
      sanitized[key] = sanitizeObject(value, excludeKeys)
    } else {
      sanitized[key] = value
    }
  }

  return sanitized
}

/**
 * Sanitize URL to prevent javascript: and data: protocols
 * @param {string} url - URL to sanitize
 * @returns {string} Sanitized URL or empty string if dangerous
 */
export function sanitizeUrl(url) {
  if (typeof url !== 'string') {
    return ''
  }

  const trimmedUrl = url.trim().toLowerCase()

  // Block dangerous protocols
  const dangerousProtocols = ['javascript:', 'data:', 'vbscript:', 'file:']

  for (const protocol of dangerousProtocols) {
    if (trimmedUrl.startsWith(protocol)) {
      return ''
    }
  }

  return url.trim()
}

/**
 * Sanitize filename to prevent path traversal
 * @param {string} filename - Filename to sanitize
 * @returns {string} Sanitized filename
 */
export function sanitizeFilename(filename) {
  if (typeof filename !== 'string') {
    return ''
  }

  // Remove path traversal attempts
  let sanitized = filename.replace(/\.\./g, '')

  // Remove path separators
  sanitized = sanitized.replace(/[/\\]/g, '')

  // Remove null bytes
  sanitized = sanitized.replace(/\0/g, '')

  return sanitized.trim()
}

/**
 * Create a Vue directive for automatic input sanitization
 * Usage: v-sanitize or v-sanitize="{ allowHtml: true }"
 */
export const sanitizeDirective = {
  mounted(el, binding) {
    const options = binding.value || {}

    const sanitizeElement = (element) => {
      if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
        element.addEventListener('blur', (event) => {
          event.target.value = sanitizeInput(event.target.value, options)
        })
      } else {
        // For other elements, sanitize text content
        if (element.textContent) {
          element.textContent = sanitizeInput(element.textContent, options)
        }
      }
    }

    sanitizeElement(el)
  }
}

export default {
  escapeHtml,
  stripHtmlTags,
  removeScripts,
  sanitizeInput,
  sanitizeObject,
  sanitizeUrl,
  sanitizeFilename,
  sanitizeDirective
}
