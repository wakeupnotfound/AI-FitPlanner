/**
 * Composable for virtual scrolling
 * Renders only visible items in a long list for better performance
 */

import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

/**
 * Use virtual scrolling for large lists
 * @param {Object} options - Configuration options
 * @param {Array} options.items - Array of items to render
 * @param {number} options.itemHeight - Height of each item in pixels
 * @param {number} options.containerHeight - Height of the scroll container
 * @param {number} options.bufferSize - Number of items to render outside viewport (default: 5)
 * @returns {Object} Virtual scroll utilities
 */
export function useVirtualScroll(options) {
  const {
    items = ref([]),
    itemHeight = 60,
    containerHeight = 600,
    bufferSize = 5
  } = options

  // State
  const scrollTop = ref(0)
  const containerRef = ref(null)

  // Computed
  const totalHeight = computed(() => {
    return items.value.length * itemHeight
  })

  const visibleCount = computed(() => {
    return Math.ceil(containerHeight / itemHeight)
  })

  const startIndex = computed(() => {
    const index = Math.floor(scrollTop.value / itemHeight) - bufferSize
    return Math.max(0, index)
  })

  const endIndex = computed(() => {
    const index = startIndex.value + visibleCount.value + bufferSize * 2
    return Math.min(items.value.length, index)
  })

  const visibleItems = computed(() => {
    return items.value.slice(startIndex.value, endIndex.value).map((item, index) => ({
      ...item,
      _index: startIndex.value + index,
      _offsetTop: (startIndex.value + index) * itemHeight
    }))
  })

  const offsetY = computed(() => {
    return startIndex.value * itemHeight
  })

  // Methods
  const handleScroll = (event) => {
    scrollTop.value = event.target.scrollTop
  }

  const scrollToIndex = (index) => {
    if (!containerRef.value) return
    const targetScrollTop = index * itemHeight
    containerRef.value.scrollTop = targetScrollTop
  }

  const scrollToTop = () => {
    scrollToIndex(0)
  }

  const scrollToBottom = () => {
    scrollToIndex(items.value.length - 1)
  }

  return {
    // Refs
    containerRef,
    scrollTop,
    
    // Computed
    totalHeight,
    visibleCount,
    startIndex,
    endIndex,
    visibleItems,
    offsetY,
    
    // Methods
    handleScroll,
    scrollToIndex,
    scrollToTop,
    scrollToBottom
  }
}

/**
 * Use virtual scrolling with dynamic item heights
 * More complex but supports variable height items
 * @param {Object} options - Configuration options
 * @returns {Object} Virtual scroll utilities
 */
export function useVirtualScrollDynamic(options) {
  const {
    items = ref([]),
    estimatedItemHeight = 60,
    containerHeight = 600,
    bufferSize = 5
  } = options

  // State
  const scrollTop = ref(0)
  const containerRef = ref(null)
  const itemHeights = ref(new Map())
  const itemOffsets = ref(new Map())

  // Calculate offsets based on measured heights
  const calculateOffsets = () => {
    let offset = 0
    items.value.forEach((item, index) => {
      itemOffsets.value.set(index, offset)
      const height = itemHeights.value.get(index) || estimatedItemHeight
      offset += height
    })
  }

  // Computed
  const totalHeight = computed(() => {
    if (itemOffsets.value.size === 0) {
      return items.value.length * estimatedItemHeight
    }
    const lastIndex = items.value.length - 1
    const lastOffset = itemOffsets.value.get(lastIndex) || 0
    const lastHeight = itemHeights.value.get(lastIndex) || estimatedItemHeight
    return lastOffset + lastHeight
  })

  const startIndex = computed(() => {
    let start = 0
    for (let i = 0; i < items.value.length; i++) {
      const offset = itemOffsets.value.get(i) || i * estimatedItemHeight
      if (offset >= scrollTop.value) {
        start = Math.max(0, i - bufferSize)
        break
      }
    }
    return start
  })

  const endIndex = computed(() => {
    const viewportBottom = scrollTop.value + containerHeight
    let end = items.value.length
    for (let i = startIndex.value; i < items.value.length; i++) {
      const offset = itemOffsets.value.get(i) || i * estimatedItemHeight
      if (offset > viewportBottom) {
        end = Math.min(items.value.length, i + bufferSize)
        break
      }
    }
    return end
  })

  const visibleItems = computed(() => {
    return items.value.slice(startIndex.value, endIndex.value).map((item, index) => {
      const actualIndex = startIndex.value + index
      return {
        ...item,
        _index: actualIndex,
        _offsetTop: itemOffsets.value.get(actualIndex) || actualIndex * estimatedItemHeight
      }
    })
  })

  // Methods
  const handleScroll = (event) => {
    scrollTop.value = event.target.scrollTop
  }

  const setItemHeight = (index, height) => {
    itemHeights.value.set(index, height)
    calculateOffsets()
  }

  const scrollToIndex = (index) => {
    if (!containerRef.value) return
    const targetScrollTop = itemOffsets.value.get(index) || index * estimatedItemHeight
    containerRef.value.scrollTop = targetScrollTop
  }

  // Watch items changes
  watch(() => items.value.length, () => {
    calculateOffsets()
  })

  return {
    // Refs
    containerRef,
    scrollTop,
    
    // Computed
    totalHeight,
    startIndex,
    endIndex,
    visibleItems,
    
    // Methods
    handleScroll,
    setItemHeight,
    scrollToIndex
  }
}
