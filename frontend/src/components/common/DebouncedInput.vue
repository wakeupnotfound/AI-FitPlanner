<template>
  <div class="debounced-input">
    <van-field
      v-model="localValue"
      v-bind="$attrs"
      :placeholder="placeholder"
      :clearable="clearable"
      :left-icon="leftIcon"
      :right-icon="rightIcon"
      @update:model-value="handleInput"
      @clear="handleClear"
    >
      <template v-if="$slots.leftIcon" #left-icon>
        <slot name="leftIcon" />
      </template>
      <template v-if="$slots.rightIcon" #right-icon>
        <slot name="rightIcon" />
      </template>
      <template v-if="showPending && isPending" #button>
        <van-loading size="16" />
      </template>
    </van-field>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useDebouncedFn } from '../../composables/useDebounce'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  delay: {
    type: Number,
    default: 300
  },
  placeholder: {
    type: String,
    default: ''
  },
  clearable: {
    type: Boolean,
    default: true
  },
  leftIcon: {
    type: String,
    default: ''
  },
  rightIcon: {
    type: String,
    default: ''
  },
  showPending: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'clear'])

const localValue = ref(props.modelValue)

// Create debounced function
const { debouncedFn, isPending } = useDebouncedFn(
  (value) => {
    emit('update:modelValue', value)
    emit('change', value)
  },
  props.delay
)

// Handle input
const handleInput = (value) => {
  localValue.value = value
  debouncedFn(value)
}

// Handle clear
const handleClear = () => {
  localValue.value = ''
  emit('update:modelValue', '')
  emit('clear')
}

// Watch for external changes
watch(() => props.modelValue, (newValue) => {
  if (newValue !== localValue.value) {
    localValue.value = newValue
  }
})
</script>

<style scoped>
.debounced-input {
  width: 100%;
}
</style>
