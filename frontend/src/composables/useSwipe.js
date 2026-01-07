/**
 * Composable for handling swipe gestures
 * Provides touch-based swipe detection for mobile interactions
 */
import { ref, onMounted, onUnmounted } from 'vue'

/**
 * Swipe direction constants
 */
export const SWIPE_DIRECTION = {
  LEFT: 'left',
  RIGHT: 'right',
  UP: 'up',
  DOWN: 'down'
}

/**
 * useSwipe composable
 * @param {Object} options - Configuration options
 * @param {number} options.threshold - Minimum distance for swipe detection (default: 50px)
 * @param {number} options.timeout - Maximum time for swipe gesture (default: 300ms)
 * @param {Function} options.onSwipeLeft - Callback for left swipe
 * @param {Function} options.onSwipeRight - Callback for right swipe
 * @param {Function} options.onSwipeUp - Callback for up swipe
 * @param {Function} options.onSwipeDown - Callback for down swipe
 * @returns {Object} Swipe state and handlers
 */
export function useSwipe(options = {}) {
  const {
    threshold = 50,
    timeout = 300,
    onSwipeLeft,
    onSwipeRight,
    onSwipeUp,
    onSwipeDown
  } = options

  // State
  const isSwiping = ref(false)
  const swipeDirection = ref(null)
  const swipeDistance = ref({ x: 0, y: 0 })

  // Touch tracking
  let startX = 0
  let startY = 0
  let startTime = 0
  let elementRef = null

  /**
   * Handle touch start
   */
  const handleTouchStart = (event) => {
    const touch = event.touches[0]
    startX = touch.clientX
    startY = touch.clientY
    startTime = Date.now()
    isSwiping.value = true
    swipeDirection.value = null
    swipeDistance.value = { x: 0, y: 0 }
  }

  /**
   * Handle touch move
   */
  const handleTouchMove = (event) => {
    if (!isSwiping.value) return

    const touch = event.touches[0]
    const deltaX = touch.clientX - startX
    const deltaY = touch.clientY - startY

    swipeDistance.value = { x: deltaX, y: deltaY }

    // Determine direction based on dominant axis
    if (Math.abs(deltaX) > Math.abs(deltaY)) {
      swipeDirection.value = deltaX > 0 ? SWIPE_DIRECTION.RIGHT : SWIPE_DIRECTION.LEFT
    } else {
      swipeDirection.value = deltaY > 0 ? SWIPE_DIRECTION.DOWN : SWIPE_DIRECTION.UP
    }
  }

  /**
   * Handle touch end
   */
  const handleTouchEnd = () => {
    if (!isSwiping.value) return

    const elapsed = Date.now() - startTime
    const { x: deltaX, y: deltaY } = swipeDistance.value

    // Check if swipe meets threshold and timeout requirements
    if (elapsed <= timeout) {
      const absX = Math.abs(deltaX)
      const absY = Math.abs(deltaY)

      if (absX >= threshold && absX > absY) {
        // Horizontal swipe
        if (deltaX > 0 && onSwipeRight) {
          onSwipeRight()
        } else if (deltaX < 0 && onSwipeLeft) {
          onSwipeLeft()
        }
      } else if (absY >= threshold && absY > absX) {
        // Vertical swipe
        if (deltaY > 0 && onSwipeDown) {
          onSwipeDown()
        } else if (deltaY < 0 && onSwipeUp) {
          onSwipeUp()
        }
      }
    }

    // Reset state
    isSwiping.value = false
    swipeDirection.value = null
    swipeDistance.value = { x: 0, y: 0 }
  }

  /**
   * Bind swipe handlers to an element
   * @param {HTMLElement} element - Element to bind handlers to
   */
  const bindSwipeHandlers = (element) => {
    if (!element) return

    elementRef = element
    element.addEventListener('touchstart', handleTouchStart, { passive: true })
    element.addEventListener('touchmove', handleTouchMove, { passive: true })
    element.addEventListener('touchend', handleTouchEnd, { passive: true })
    element.addEventListener('touchcancel', handleTouchEnd, { passive: true })
  }

  /**
   * Unbind swipe handlers from element
   */
  const unbindSwipeHandlers = () => {
    if (!elementRef) return

    elementRef.removeEventListener('touchstart', handleTouchStart)
    elementRef.removeEventListener('touchmove', handleTouchMove)
    elementRef.removeEventListener('touchend', handleTouchEnd)
    elementRef.removeEventListener('touchcancel', handleTouchEnd)
    elementRef = null
  }

  // Cleanup on unmount
  onUnmounted(() => {
    unbindSwipeHandlers()
  })

  return {
    isSwiping,
    swipeDirection,
    swipeDistance,
    bindSwipeHandlers,
    unbindSwipeHandlers,
    SWIPE_DIRECTION
  }
}

/**
 * useSwipeToDelete composable
 * Specialized swipe handler for swipe-to-delete functionality
 * @param {Object} options - Configuration options
 * @param {Function} options.onDelete - Callback when delete is triggered
 * @param {number} options.deleteThreshold - Distance to trigger delete (default: 100px)
 * @returns {Object} Swipe state and handlers
 */
export function useSwipeToDelete(options = {}) {
  const { onDelete, deleteThreshold = 100 } = options

  const swipeOffset = ref(0)
  const isDeleting = ref(false)

  let startX = 0
  let currentElement = null

  const handleTouchStart = (event) => {
    const touch = event.touches[0]
    startX = touch.clientX
    isDeleting.value = false
  }

  const handleTouchMove = (event) => {
    const touch = event.touches[0]
    const deltaX = touch.clientX - startX

    // Only allow left swipe (negative delta)
    if (deltaX < 0) {
      swipeOffset.value = Math.max(deltaX, -deleteThreshold * 1.5)
    }
  }

  const handleTouchEnd = () => {
    if (Math.abs(swipeOffset.value) >= deleteThreshold) {
      isDeleting.value = true
      if (onDelete) {
        onDelete()
      }
    }
    
    // Reset offset with animation
    swipeOffset.value = 0
  }

  const bindToElement = (element) => {
    if (!element) return

    currentElement = element
    element.addEventListener('touchstart', handleTouchStart, { passive: true })
    element.addEventListener('touchmove', handleTouchMove, { passive: false })
    element.addEventListener('touchend', handleTouchEnd, { passive: true })
  }

  const unbindFromElement = () => {
    if (!currentElement) return

    currentElement.removeEventListener('touchstart', handleTouchStart)
    currentElement.removeEventListener('touchmove', handleTouchMove)
    currentElement.removeEventListener('touchend', handleTouchEnd)
    currentElement = null
  }

  onUnmounted(() => {
    unbindFromElement()
  })

  return {
    swipeOffset,
    isDeleting,
    bindToElement,
    unbindFromElement
  }
}

export default useSwipe
