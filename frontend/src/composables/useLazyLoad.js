/**
 * Composable for lazy loading images
 * Provides utilities for implementing lazy loading with Intersection Observer
 */

import { ref, onMounted, onUnmounted } from 'vue'

/**
 * Use lazy loading for images
 * @param {Object} options - Configuration options
 * @param {string} options.rootMargin - Margin around root (default: '50px')
 * @param {number} options.threshold - Intersection threshold (default: 0.01)
 * @returns {Object} Lazy loading utilities
 */
export function useLazyLoad(options = {}) {
  const {
    rootMargin = '50px',
    threshold = 0.01
  } = options

  const observer = ref(null)
  const loadedImages = ref(new Set())

  /**
   * Initialize Intersection Observer
   */
  const initObserver = () => {
    if ('IntersectionObserver' in window) {
      observer.value = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              loadImage(entry.target)
              observer.value.unobserve(entry.target)
            }
          })
        },
        {
          rootMargin,
          threshold
        }
      )
    }
  }

  /**
   * Load image by setting src from data-src
   * @param {HTMLImageElement} img - Image element to load
   */
  const loadImage = (img) => {
    const src = img.dataset.src
    if (!src) return

    // Create a new image to preload
    const tempImg = new Image()
    tempImg.onload = () => {
      img.src = src
      img.classList.add('loaded')
      img.classList.remove('loading')
      loadedImages.value.add(src)
    }
    tempImg.onerror = () => {
      img.classList.add('error')
      img.classList.remove('loading')
    }
    tempImg.src = src
  }

  /**
   * Observe an image element
   * @param {HTMLImageElement} img - Image element to observe
   */
  const observe = (img) => {
    if (!img) return

    // Add loading class
    img.classList.add('loading')

    // If IntersectionObserver is not supported, load immediately
    if (!observer.value) {
      loadImage(img)
      return
    }

    observer.value.observe(img)
  }

  /**
   * Unobserve an image element
   * @param {HTMLImageElement} img - Image element to unobserve
   */
  const unobserve = (img) => {
    if (observer.value && img) {
      observer.value.unobserve(img)
    }
  }

  /**
   * Disconnect observer
   */
  const disconnect = () => {
    if (observer.value) {
      observer.value.disconnect()
    }
  }

  onMounted(() => {
    initObserver()
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    observe,
    unobserve,
    disconnect,
    loadedImages
  }
}

/**
 * Get placeholder image data URL
 * @param {number} width - Image width
 * @param {number} height - Image height
 * @param {string} color - Background color (default: '#f0f0f0')
 * @returns {string} Data URL for placeholder
 */
export function getPlaceholderImage(width = 300, height = 200, color = '#f0f0f0') {
  // Create a simple SVG placeholder
  const svg = `
    <svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">
      <rect width="100%" height="100%" fill="${color}"/>
      <text 
        x="50%" 
        y="50%" 
        font-family="Arial, sans-serif" 
        font-size="14" 
        fill="#999" 
        text-anchor="middle" 
        dominant-baseline="middle"
      >
        Loading...
      </text>
    </svg>
  `
  return `data:image/svg+xml;base64,${btoa(svg)}`
}

/**
 * Create a blur hash placeholder (simplified version)
 * @param {string} color - Primary color for blur effect
 * @returns {string} Data URL for blur placeholder
 */
export function getBlurPlaceholder(color = '#e0e0e0') {
  const svg = `
    <svg width="40" height="40" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <filter id="blur">
          <feGaussianBlur stdDeviation="10"/>
        </filter>
      </defs>
      <rect width="100%" height="100%" fill="${color}" filter="url(#blur)"/>
    </svg>
  `
  return `data:image/svg+xml;base64,${btoa(svg)}`
}
