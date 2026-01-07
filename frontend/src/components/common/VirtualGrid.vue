<template>
  <div 
    ref="containerRef"
    class="virtual-grid-container"
    :style="containerStyle"
    @scroll="handleScroll"
  >
    <div 
      class="virtual-grid-phantom"
      :style="{ height: `${totalHeight}px` }"
    />
    <div 
      class="virtual-grid-content"
      :style="contentStyle"
    >
      <div
        v-for="item in visibleItems"
        :key="getItemKey(item)"
        class="virtual-grid-item"
        :style="itemStyle"
      >
        <slot :item="item" :index="item._index" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    required: true
  },
  itemHeight: {
    type: Number,
    default: 200
  },
  columns: {
    type: Number,
    default: 2
  },
  gap: {
    type: Number,
    default: 16
  },
  height: {
    type: [Number, String],
    default: 600
  },
  bufferRows: {
    type: Number,
    default: 2
  },
  itemKey: {
    type: [String, Function],
    default: 'id'
  }
})

const emit = defineEmits(['scroll', 'reach-bottom'])

const containerRef = ref(null)
const scrollTop = ref(0)

// Computed
const containerHeight = computed(() => {
  if (typeof props.height === 'number') {
    return props.height
  }
  return parseInt(props.height) || 600
})

const containerStyle = computed(() => ({
  height: typeof props.height === 'number' ? `${props.height}px` : props.height
}))

const rowHeight = computed(() => props.itemHeight + props.gap)

const totalRows = computed(() => Math.ceil(props.items.length / props.columns))

const totalHeight = computed(() => {
  return totalRows.value * rowHeight.value - props.gap
})

const visibleRows = computed(() => Math.ceil(containerHeight.value / rowHeight.value))

const startRow = computed(() => {
  const row = Math.floor(scrollTop.value / rowHeight.value) - props.bufferRows
  return Math.max(0, row)
})

const endRow = computed(() => {
  const row = startRow.value + visibleRows.value + props.bufferRows * 2
  return Math.min(totalRows.value, row)
})

const startIndex = computed(() => startRow.value * props.columns)
const endIndex = computed(() => Math.min(props.items.length, endRow.value * props.columns))

const visibleItems = computed(() => {
  return props.items.slice(startIndex.value, endIndex.value).map((item, index) => ({
    ...item,
    _index: startIndex.value + index
  }))
})

const offsetY = computed(() => startRow.value * rowHeight.value)

const contentStyle = computed(() => ({
  transform: `translateY(${offsetY.value}px)`,
  display: 'grid',
  gridTemplateColumns: `repeat(${props.columns}, 1fr)`,
  gap: `${props.gap}px`
}))

const itemStyle = computed(() => ({
  height: `${props.itemHeight}px`
}))

// Methods
const getItemKey = (item) => {
  if (typeof props.itemKey === 'function') {
    return props.itemKey(item)
  }
  return item[props.itemKey] || item._index
}

const handleScroll = (event) => {
  scrollTop.value = event.target.scrollTop
  
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

const scrollToIndex = (index) => {
  if (!containerRef.value) return
  const row = Math.floor(index / props.columns)
  const targetScrollTop = row * rowHeight.value
  containerRef.value.scrollTop = targetScrollTop
}

const scrollToTop = () => {
  if (containerRef.value) {
    containerRef.value.scrollTop = 0
  }
}

const scrollToBottom = () => {
  if (containerRef.value) {
    containerRef.value.scrollTop = totalHeight.value
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
.virtual-grid-container {
  position: relative;
  overflow-y: auto;
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
}

.virtual-grid-phantom {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  z-index: -1;
}

.virtual-grid-content {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  will-change: transform;
  padding: 0 16px;
}

.virtual-grid-item {
  overflow: hidden;
}
</style>
