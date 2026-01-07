/**
 * Composable for debouncing and throttling
 * Provides utilities for optimizing frequent function calls
 */

import { ref, watch, onUnmounted } from 'vue'

/**
 * Debounce a function call
 * Delays execution until after a specified time has elapsed since the last call
 * @param {Function} fn - Function to debounce
 * @param {number} delay - Delay in milliseconds (default: 300)
 * @returns {Function} Debounced function
 */
export function debounce(fn, delay = 300) {
  let timeoutId = null

  const debounced = function (...args) {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    timeoutId = setTimeout(() => {
      fn.apply(this, args)
      timeoutId = null
    }, delay)
  }

  debounced.cancel = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
  }

  return debounced
}

/**
 * Throttle a function call
 * Ensures function is called at most once per specified time period
 * @param {Function} fn - Function to throttle
 * @param {number} limit - Time limit in milliseconds (default: 300)
 * @returns {Function} Throttled function
 */
export function throttle(fn, limit = 300) {
  let inThrottle = false
  let lastResult = null

  const throttled = function (...args) {
    if (!inThrottle) {
      lastResult = fn.apply(this, args)
      inThrottle = true

      setTimeout(() => {
        inThrottle = false
      }, limit)
    }

    return lastResult
  }

  return throttled
}

/**
 * Use debounced value
 * Returns a debounced version of the input value
 * @param {Ref} value - Reactive value to debounce
 * @param {number} delay - Delay in milliseconds (default: 300)
 * @returns {Ref} Debounced value
 */
export function useDebouncedValue(value, delay = 300) {
  const debouncedValue = ref(value.value)
  let timeoutId = null

  watch(value, (newValue) => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    timeoutId = setTimeout(() => {
      debouncedValue.value = newValue
      timeoutId = null
    }, delay)
  })

  onUnmounted(() => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
  })

  return debouncedValue
}

/**
 * Use debounced function
 * Returns a debounced version of the input function with reactive state
 * @param {Function} fn - Function to debounce
 * @param {number} delay - Delay in milliseconds (default: 300)
 * @returns {Object} Debounced function and state
 */
export function useDebouncedFn(fn, delay = 300) {
  const isPending = ref(false)
  let timeoutId = null

  const debouncedFn = (...args) => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    isPending.value = true

    return new Promise((resolve, reject) => {
      timeoutId = setTimeout(async () => {
        try {
          const result = await fn(...args)
          isPending.value = false
          timeoutId = null
          resolve(result)
        } catch (error) {
          isPending.value = false
          timeoutId = null
          reject(error)
        }
      }, delay)
    })
  }

  const cancel = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
      isPending.value = false
    }
  }

  onUnmounted(() => {
    cancel()
  })

  return {
    debouncedFn,
    isPending,
    cancel
  }
}

/**
 * Use throttled function
 * Returns a throttled version of the input function with reactive state
 * @param {Function} fn - Function to throttle
 * @param {number} limit - Time limit in milliseconds (default: 300)
 * @returns {Object} Throttled function and state
 */
export function useThrottledFn(fn, limit = 300) {
  const isThrottled = ref(false)
  let timeoutId = null

  const throttledFn = async (...args) => {
    if (isThrottled.value) {
      return
    }

    isThrottled.value = true

    try {
      const result = await fn(...args)
      
      timeoutId = setTimeout(() => {
        isThrottled.value = false
        timeoutId = null
      }, limit)

      return result
    } catch (error) {
      isThrottled.value = false
      throw error
    }
  }

  const reset = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
    isThrottled.value = false
  }

  onUnmounted(() => {
    reset()
  })

  return {
    throttledFn,
    isThrottled,
    reset
  }
}

/**
 * Use debounced search
 * Specialized composable for search input with debouncing
 * @param {Function} searchFn - Search function to execute
 * @param {Object} options - Configuration options
 * @returns {Object} Search utilities
 */
export function useDebouncedSearch(searchFn, options = {}) {
  const {
    delay = 300,
    minLength = 0,
    immediate = false
  } = options

  const searchQuery = ref('')
  const searchResults = ref([])
  const isSearching = ref(false)
  const error = ref(null)

  let timeoutId = null

  const executeSearch = async (query) => {
    if (query.length < minLength) {
      searchResults.value = []
      return
    }

    isSearching.value = true
    error.value = null

    try {
      const results = await searchFn(query)
      searchResults.value = results
    } catch (err) {
      error.value = err
      searchResults.value = []
    } finally {
      isSearching.value = false
    }
  }

  const debouncedSearch = (query) => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    searchQuery.value = query

    if (query.length < minLength) {
      searchResults.value = []
      isSearching.value = false
      return
    }

    isSearching.value = true

    timeoutId = setTimeout(() => {
      executeSearch(query)
      timeoutId = null
    }, delay)
  }

  const clearSearch = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
    searchQuery.value = ''
    searchResults.value = []
    isSearching.value = false
    error.value = null
  }

  // Immediate search if enabled
  if (immediate && searchQuery.value) {
    executeSearch(searchQuery.value)
  }

  onUnmounted(() => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
  })

  return {
    searchQuery,
    searchResults,
    isSearching,
    error,
    search: debouncedSearch,
    clearSearch
  }
}

/**
 * Default debounce delays for different use cases
 */
export const DEBOUNCE_DELAYS = {
  SEARCH: 300,        // Search input
  INPUT: 500,         // General input
  RESIZE: 150,        // Window resize
  SCROLL: 100,        // Scroll events
  API_CALL: 500,      // API requests
  AUTOCOMPLETE: 200   // Autocomplete suggestions
}

/**
 * Default throttle limits for different use cases
 */
export const THROTTLE_LIMITS = {
  SCROLL: 100,        // Scroll events
  RESIZE: 150,        // Window resize
  MOUSE_MOVE: 50,     // Mouse movement
  API_CALL: 1000,     // API requests
  BUTTON_CLICK: 500   // Button clicks
}
