/**
 * Viewport Fix Utility
 * 
 * Handles browser-specific viewport issues, particularly for iOS Safari
 * where the address bar affects viewport height calculations.
 */

/**
 * Set CSS custom property for actual viewport height
 * This fixes the iOS Safari 100vh issue where the address bar is included
 */
export function setViewportHeight() {
  // Calculate actual viewport height
  const vh = window.innerHeight * 0.01
  
  // Set CSS custom property
  document.documentElement.style.setProperty('--vh', `${vh}px`)
}

/**
 * Initialize viewport height fix
 * Sets up event listeners for resize and orientation change
 */
export function initViewportFix() {
  // Set initial value
  setViewportHeight()
  
  // Update on resize
  window.addEventListener('resize', setViewportHeight)
  
  // Update on orientation change (mobile devices)
  window.addEventListener('orientationchange', () => {
    // Delay to ensure new dimensions are available
    setTimeout(setViewportHeight, 100)
  })
  
  // Update on visual viewport resize (iOS Safari address bar)
  if (window.visualViewport) {
    window.visualViewport.addEventListener('resize', setViewportHeight)
  }
}

/**
 * Cleanup viewport fix event listeners
 */
export function cleanupViewportFix() {
  window.removeEventListener('resize', setViewportHeight)
  window.removeEventListener('orientationchange', setViewportHeight)
  
  if (window.visualViewport) {
    window.visualViewport.removeEventListener('resize', setViewportHeight)
  }
}

/**
 * Detect if running on iOS
 */
export function isIOS() {
  return /iPad|iPhone|iPod/.test(navigator.userAgent) && !window.MSStream
}

/**
 * Detect if running on iOS Safari
 */
export function isIOSSafari() {
  const ua = navigator.userAgent
  const iOS = /iPad|iPhone|iPod/.test(ua) && !window.MSStream
  const webkit = /WebKit/.test(ua)
  const notChrome = !(/CriOS/.test(ua))
  const notFirefox = !(/FxiOS/.test(ua))
  
  return iOS && webkit && notChrome && notFirefox
}

/**
 * Detect if running on Android
 */
export function isAndroid() {
  return /Android/.test(navigator.userAgent)
}

/**
 * Detect if running on Android Chrome
 */
export function isAndroidChrome() {
  const ua = navigator.userAgent
  return /Android/.test(ua) && /Chrome/.test(ua) && !/Edge/.test(ua)
}

/**
 * Detect browser type
 */
export function detectBrowser() {
  const ua = navigator.userAgent
  
  if (/Firefox/.test(ua)) {
    return 'firefox'
  } else if (/Chrome/.test(ua) && !/Edge/.test(ua)) {
    return 'chrome'
  } else if (/Safari/.test(ua) && !/Chrome/.test(ua)) {
    return 'safari'
  } else if (/Edge/.test(ua)) {
    return 'edge'
  } else {
    return 'unknown'
  }
}

/**
 * Get browser version
 */
export function getBrowserVersion() {
  const ua = navigator.userAgent
  let match
  
  if ((match = ua.match(/Firefox\/(\d+)/))) {
    return { browser: 'firefox', version: parseInt(match[1]) }
  } else if ((match = ua.match(/Chrome\/(\d+)/))) {
    return { browser: 'chrome', version: parseInt(match[1]) }
  } else if ((match = ua.match(/Version\/(\d+).*Safari/))) {
    return { browser: 'safari', version: parseInt(match[1]) }
  } else if ((match = ua.match(/Edge\/(\d+)/))) {
    return { browser: 'edge', version: parseInt(match[1]) }
  }
  
  return { browser: 'unknown', version: 0 }
}

/**
 * Check if browser supports a specific feature
 */
export function supportsFeature(feature) {
  const features = {
    // CSS features
    'css-grid': CSS.supports('display', 'grid'),
    'css-flexbox': CSS.supports('display', 'flex'),
    'css-gap': CSS.supports('gap', '1rem'),
    'css-aspect-ratio': CSS.supports('aspect-ratio', '16/9'),
    'css-backdrop-filter': CSS.supports('backdrop-filter', 'blur(10px)') || 
                          CSS.supports('-webkit-backdrop-filter', 'blur(10px)'),
    
    // JavaScript APIs
    'service-worker': 'serviceWorker' in navigator,
    'intersection-observer': 'IntersectionObserver' in window,
    'resize-observer': 'ResizeObserver' in window,
    'mutation-observer': 'MutationObserver' in window,
    'local-storage': 'localStorage' in window,
    'session-storage': 'sessionStorage' in window,
    'geolocation': 'geolocation' in navigator,
    'notification': 'Notification' in window,
    'clipboard': 'clipboard' in navigator,
    
    // Touch and pointer events
    'touch-events': 'ontouchstart' in window,
    'pointer-events': 'PointerEvent' in window,
    
    // Network APIs
    'online-status': 'onLine' in navigator,
    'network-information': 'connection' in navigator || 
                          'mozConnection' in navigator || 
                          'webkitConnection' in navigator,
    
    // Performance APIs
    'performance': 'performance' in window,
    'performance-observer': 'PerformanceObserver' in window,
    
    // Media APIs
    'webgl': (() => {
      try {
        const canvas = document.createElement('canvas')
        return !!(canvas.getContext('webgl') || canvas.getContext('experimental-webgl'))
      } catch (e) {
        return false
      }
    })(),
    
    // Modern JavaScript features
    'async-await': (() => {
      try {
        eval('(async () => {})')
        return true
      } catch (e) {
        return false
      }
    })(),
    'promises': 'Promise' in window,
    'fetch': 'fetch' in window,
    'arrow-functions': (() => {
      try {
        eval('() => {}')
        return true
      } catch (e) {
        return false
      }
    })()
  }
  
  return features[feature] !== undefined ? features[feature] : false
}

/**
 * Apply browser-specific fixes
 */
export function applyBrowserFixes() {
  const browser = detectBrowser()
  const body = document.body
  
  // Add browser class to body for CSS targeting
  body.classList.add(`browser-${browser}`)
  
  // Add iOS class
  if (isIOS()) {
    body.classList.add('is-ios')
  }
  
  // Add Android class
  if (isAndroid()) {
    body.classList.add('is-android')
  }
  
  // Add touch device class
  if (supportsFeature('touch-events')) {
    body.classList.add('is-touch')
  }
  
  // Initialize viewport fix for iOS
  if (isIOS()) {
    initViewportFix()
  }
  
  // Prevent zoom on input focus for iOS
  if (isIOS()) {
    const meta = document.querySelector('meta[name="viewport"]')
    if (meta) {
      const content = meta.getAttribute('content')
      if (!content.includes('maximum-scale')) {
        meta.setAttribute('content', content + ', maximum-scale=1.0')
      }
    }
  }
}

/**
 * Log browser information for debugging
 */
export function logBrowserInfo() {
  const info = {
    userAgent: navigator.userAgent,
    browser: detectBrowser(),
    version: getBrowserVersion(),
    isIOS: isIOS(),
    isIOSSafari: isIOSSafari(),
    isAndroid: isAndroid(),
    isAndroidChrome: isAndroidChrome(),
    viewport: {
      width: window.innerWidth,
      height: window.innerHeight,
      devicePixelRatio: window.devicePixelRatio
    },
    features: {
      serviceWorker: supportsFeature('service-worker'),
      intersectionObserver: supportsFeature('intersection-observer'),
      touchEvents: supportsFeature('touch-events'),
      cssGrid: supportsFeature('css-grid'),
      cssGap: supportsFeature('css-gap')
    }
  }
  
  console.log('Browser Information:', info)
  return info
}

export default {
  setViewportHeight,
  initViewportFix,
  cleanupViewportFix,
  isIOS,
  isIOSSafari,
  isAndroid,
  isAndroidChrome,
  detectBrowser,
  getBrowserVersion,
  supportsFeature,
  applyBrowserFixes,
  logBrowserInfo
}
