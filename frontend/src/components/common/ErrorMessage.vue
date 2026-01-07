<template>
  <!-- Toast notification variant (uses Vant's showNotify) -->
  <div v-if="toast && visible" class="error-toast-trigger" />

  <!-- Inline error display variant -->
  <div 
    v-else-if="inline && message" 
    class="error-inline"
    :class="[`error-inline--${type}`]"
  >
    <van-icon 
      :name="iconName" 
      :color="iconColor" 
      class="error-inline__icon"
    />
    <div class="error-inline__content">
      <span class="error-inline__message">{{ message }}</span>
      <van-button 
        v-if="retryable" 
        type="primary" 
        size="small" 
        plain
        class="error-inline__retry"
        @click="handleRetry"
      >
        {{ retryText || t('app.retry') }}
      </van-button>
    </div>
    <van-icon 
      v-if="closable" 
      name="cross" 
      class="error-inline__close"
      @click="handleClose"
    />
  </div>

  <!-- Empty state variant -->
  <van-empty 
    v-else-if="empty && message"
    :image="emptyImage"
    :description="message"
  >
    <van-button 
      v-if="retryable"
      type="primary" 
      size="small"
      @click="handleRetry"
    >
      {{ retryText || t('app.retry') }}
    </van-button>
  </van-empty>
</template>

<script setup>
import { computed, watch, onMounted } from 'vue'
import { showNotify, closeNotify } from 'vant'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  // Common props
  message: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'error',
    validator: (value) => ['error', 'warning', 'info', 'success'].includes(value)
  },
  closable: {
    type: Boolean,
    default: false
  },
  retryable: {
    type: Boolean,
    default: false
  },
  retryText: {
    type: String,
    default: ''
  },

  // Toast variant props
  toast: {
    type: Boolean,
    default: false
  },
  visible: {
    type: Boolean,
    default: false
  },
  duration: {
    type: Number,
    default: 3000
  },
  position: {
    type: String,
    default: 'top',
    validator: (value) => ['top', 'bottom'].includes(value)
  },

  // Inline variant props
  inline: {
    type: Boolean,
    default: true
  },

  // Empty state variant props
  empty: {
    type: Boolean,
    default: false
  },
  emptyImage: {
    type: String,
    default: 'error'
  }
})

const emit = defineEmits(['close', 'retry', 'update:visible'])

// Computed icon based on type
const iconName = computed(() => {
  const icons = {
    error: 'warning-o',
    warning: 'info-o',
    info: 'info-o',
    success: 'checked'
  }
  return icons[props.type] || 'warning-o'
})

// Computed icon color based on type
const iconColor = computed(() => {
  const colors = {
    error: '#ee0a24',
    warning: '#ff976a',
    info: '#1989fa',
    success: '#07c160'
  }
  return colors[props.type] || '#ee0a24'
})

// Computed notify type for toast
const notifyType = computed(() => {
  const types = {
    error: 'danger',
    warning: 'warning',
    info: 'primary',
    success: 'success'
  }
  return types[props.type] || 'danger'
})

// Show toast notification
function showToastNotification() {
  if (props.toast && props.visible && props.message) {
    showNotify({
      type: notifyType.value,
      message: props.message,
      duration: props.duration,
      position: props.position,
      onClose: () => {
        emit('update:visible', false)
        emit('close')
      }
    })
  }
}

// Handle close action
function handleClose() {
  if (props.toast) {
    closeNotify()
  }
  emit('update:visible', false)
  emit('close')
}

// Handle retry action
function handleRetry() {
  emit('retry')
}

// Watch for visibility changes to show/hide toast
watch(
  () => props.visible,
  (newVal) => {
    if (newVal && props.toast) {
      showToastNotification()
    }
  }
)

// Show toast on mount if visible
onMounted(() => {
  if (props.toast && props.visible) {
    showToastNotification()
  }
})
</script>

<style scoped>
.error-toast-trigger {
  display: none;
}

.error-inline {
  display: flex;
  align-items: flex-start;
  padding: 12px 16px;
  border-radius: 8px;
  margin: 8px 0;
}

.error-inline--error {
  background-color: #fff0f0;
  border: 1px solid #ffccc7;
}

.error-inline--warning {
  background-color: #fffbe6;
  border: 1px solid #ffe58f;
}

.error-inline--info {
  background-color: #e6f7ff;
  border: 1px solid #91d5ff;
}

.error-inline--success {
  background-color: #f6ffed;
  border: 1px solid #b7eb8f;
}

.error-inline__icon {
  flex-shrink: 0;
  font-size: 18px;
  margin-right: 8px;
  margin-top: 2px;
}

.error-inline__content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.error-inline__message {
  font-size: 14px;
  line-height: 1.5;
  color: #323233;
  word-break: break-word;
}

.error-inline__retry {
  align-self: flex-start;
}

.error-inline__close {
  flex-shrink: 0;
  font-size: 16px;
  color: #969799;
  cursor: pointer;
  margin-left: 8px;
  padding: 4px;
}

.error-inline__close:active {
  opacity: 0.7;
}
</style>
