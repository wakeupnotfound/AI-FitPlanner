<template>
  <div class="lazy-image-wrapper" :style="wrapperStyle">
    <img
      ref="imageRef"
      :class="imageClasses"
      :alt="alt"
      :width="width"
      :height="height"
      @load="onLoad"
      @error="onError"
    />
    <div v-if="loading" class="lazy-image-placeholder">
      <van-loading v-if="showLoader" size="24" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useLazyLoad, getPlaceholderImage } from '../../composables/useLazyLoad'

const props = defineProps({
  src: {
    type: String,
    required: true
  },
  alt: {
    type: String,
    default: ''
  },
  width: {
    type: [Number, String],
    default: null
  },
  height: {
    type: [Number, String],
    default: null
  },
  placeholder: {
    type: String,
    default: null
  },
  showLoader: {
    type: Boolean,
    default: true
  },
  aspectRatio: {
    type: String,
    default: null // e.g., '16/9', '4/3', '1/1'
  }
})

const emit = defineEmits(['load', 'error'])

const imageRef = ref(null)
const loading = ref(true)
const error = ref(false)
const loaded = ref(false)

const { observe, unobserve } = useLazyLoad()

const imageClasses = computed(() => ({
  'lazy-image': true,
  'lazy-loading': loading.value,
  'lazy-loaded': loaded.value,
  'lazy-error': error.value
}))

const wrapperStyle = computed(() => {
  const styles = {}
  
  if (props.aspectRatio) {
    styles.aspectRatio = props.aspectRatio
  }
  
  if (props.width) {
    styles.width = typeof props.width === 'number' ? `${props.width}px` : props.width
  }
  
  if (props.height && !props.aspectRatio) {
    styles.height = typeof props.height === 'number' ? `${props.height}px` : props.height
  }
  
  return styles
})

const onLoad = () => {
  loading.value = false
  loaded.value = true
  emit('load')
}

const onError = () => {
  loading.value = false
  error.value = true
  emit('error')
}

const setupLazyLoad = () => {
  if (!imageRef.value) return

  // Set placeholder
  const placeholderSrc = props.placeholder || getPlaceholderImage(
    props.width || 300,
    props.height || 200
  )
  imageRef.value.src = placeholderSrc

  // Set data-src for lazy loading
  imageRef.value.dataset.src = props.src

  // Start observing
  observe(imageRef.value)
}

watch(() => props.src, () => {
  if (imageRef.value) {
    loading.value = true
    loaded.value = false
    error.value = false
    setupLazyLoad()
  }
})

onMounted(() => {
  setupLazyLoad()
})

onUnmounted(() => {
  if (imageRef.value) {
    unobserve(imageRef.value)
  }
})
</script>

<style scoped>
.lazy-image-wrapper {
  position: relative;
  overflow: hidden;
  background-color: #f5f5f5;
  display: inline-block;
  width: 100%;
}

.lazy-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: opacity 0.3s ease-in-out, filter 0.3s ease-in-out;
}

.lazy-loading {
  opacity: 0.6;
  filter: blur(5px);
}

.lazy-loaded {
  opacity: 1;
  filter: blur(0);
}

.lazy-error {
  opacity: 0.5;
}

.lazy-image-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.8);
  pointer-events: none;
}
</style>
