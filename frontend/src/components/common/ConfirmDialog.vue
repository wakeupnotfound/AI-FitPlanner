<template>
  <van-dialog
    :show="dialogVisible"
    :title="title"
    :message="message"
    :show-cancel-button="showCancelButton"
    :show-confirm-button="showConfirmButton"
    :confirm-button-text="confirmText || t('app.confirm')"
    :cancel-button-text="cancelText || t('app.cancel')"
    :confirm-button-color="computedConfirmColor"
    :cancel-button-color="cancelButtonColor"
    :confirm-button-disabled="confirmDisabled"
    :close-on-click-overlay="closeOnClickOverlay"
    :before-close="handleBeforeClose"
    :theme="theme"
    :width="width"
    class="confirm-dialog"
    @update:show="updateVisible"
    @confirm="handleConfirm"
    @cancel="handleCancel"
    @close="handleClose"
  >
    <!-- Custom content slot -->
    <template v-if="$slots.default" #default>
      <div class="confirm-dialog__content">
        <slot />
      </div>
    </template>

    <!-- Custom title slot -->
    <template v-if="$slots.title" #title>
      <slot name="title" />
    </template>

    <!-- Custom footer slot -->
    <template v-if="$slots.footer" #footer>
      <slot name="footer" :confirm="handleConfirm" :cancel="handleCancel" />
    </template>
  </van-dialog>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  // v-model for visibility
  modelValue: {
    type: Boolean,
    default: false
  },
  // Dialog title
  title: {
    type: String,
    default: ''
  },
  // Dialog message
  message: {
    type: String,
    default: ''
  },
  // Confirm button text
  confirmText: {
    type: String,
    default: ''
  },
  // Cancel button text
  cancelText: {
    type: String,
    default: ''
  },
  // Show cancel button
  showCancelButton: {
    type: Boolean,
    default: true
  },
  // Show confirm button
  showConfirmButton: {
    type: Boolean,
    default: true
  },
  // Confirm button color
  confirmButtonColor: {
    type: String,
    default: '#1989fa'
  },
  // Cancel button color
  cancelButtonColor: {
    type: String,
    default: '#323233'
  },
  // Disable confirm button
  confirmDisabled: {
    type: Boolean,
    default: false
  },
  // Close on overlay click
  closeOnClickOverlay: {
    type: Boolean,
    default: false
  },
  // Dialog theme
  theme: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'round-button'].includes(value)
  },
  // Dialog width
  width: {
    type: [String, Number],
    default: '320px'
  },
  // Loading state for async confirm
  loading: {
    type: Boolean,
    default: false
  },
  // Danger mode (red confirm button)
  danger: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel', 'close'])

// Computed visibility with v-model support
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// Override confirm button color for danger mode
const computedConfirmColor = computed(() => {
  if (props.danger) {
    return '#ee0a24'
  }
  return props.confirmButtonColor
})

// Handle before close for async operations
async function handleBeforeClose(action) {
  if (action === 'confirm' && props.loading) {
    // Keep dialog open during loading
    return false
  }
  return true
}

// Handle confirm action
function handleConfirm() {
  emit('confirm')
  if (!props.loading) {
    dialogVisible.value = false
  }
}

// Handle cancel action
function handleCancel() {
  emit('cancel')
  dialogVisible.value = false
}

// Handle close action
function handleClose() {
  emit('close')
}

// Handle visibility update from dialog
function updateVisible(value) {
  dialogVisible.value = value
}
</script>

<style scoped>
.confirm-dialog__content {
  padding: 16px 24px;
  font-size: 14px;
  line-height: 1.6;
  color: #646566;
}

:deep(.van-dialog__header) {
  font-weight: 600;
}

:deep(.van-dialog__message) {
  font-size: 14px;
  line-height: 1.6;
}
</style>
