<template>
  <div 
    ref="containerRef"
    class="virtual-list-container"
    :style="containerStyle"
    @scroll="handleScroll"
  >
    <div 
      class="virtual-list-phantom"
      :style="{ height: `${totalHeight}px` }"
    />
    <div 
      class="virtual-list-content"
      :style="{ transform: `translateY(${offsetY}px)` }"
    >
      <div
        v-for="item in visibleItems"
        :key="getItemKey(item)"
        class="virtual-list-item"
        :style="{ height: `${itemHeight}px` }"
      >
        <slot :item="item" :index="item._index" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, toRef } from 'vue'
import { useVirtualScroll } from '../../composables/useVirtualScroll'

const props = defineProps({
  items: {
    type: Array,
    required: true
  },
  itemHeight: {
    type: Number,
    default: 60
  },
  height: {
    type: [Number, String],
    default: 600
  },
  bufferSize: {
    type: Number,
    default: 5
  },
  itemKey: {
    type: [String, Function],
    default: 'id'
  }
})

const emit = defineEmits(['scroll', 'reach-bottom'])

// Convert height to number
const containerHeight = computed(() => {
  if (typeof props.height === 'number') {
    return props.height
  }
  return parseInt(props.height) || 600
})

const containerStyle = computed(() => ({
  height: typeof props.height === 'number' ? `${props.height}px` : props.height
}))

// Use virtual scroll composable
const {
  containerRef,
  scrollTop,
  totalHeight,
  visibleItems,
  offsetY,
  handleScroll: onScroll,
  scrollToIndex,
  scrollToTop,
  scrollToBottom
} = useVirtualScroll({
  items: toRef(props, 'items'),
  itemHeight: props.itemHeight,
  containerHeight: containerHeight.value,
  bufferSize: props.bufferSize
})

// Get item key
const getItemKey = (item) => {
  if (typeof props.itemKey === 'function') {
    return props.itemKey(item)
  }
  return item[props.itemKey] || item._index
}

// Handle scroll with emit
const handleScroll = (event) => {
  onScroll(event)
  emit('scroll', {
    scrollTop: scrollTop.value,
    scrollHeight: totalHeight.value,
    clientHeight: containerHeight.value
  })

  // Check if reached bottom
  const isBottom = scrollTop.value + containerHeight.value >= totalHeight.value - 50
  if (isBottom) {
    emit('reach-bottom')
  }
}

// Expose methods
defineExpose({
  scrollToIndex,
  scrollToTop,
  scrollToBottom
})
</script>

<style scoped>
.virtual-list-container {
  position: relative;
  overflow-y: auto;
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
}

.virtual-list-phantom {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  z-index: -1;
}

.virtual-list-content {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  will-change: transform;
}

.virtual-list-item {
  overflow: hidden;
}
</style>
