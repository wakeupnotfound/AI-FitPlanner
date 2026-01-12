<template>
  <van-cell-group class="ai-config-card" inset>
    <!-- Header with provider and status -->
    <div class="card-header">
      <div class="provider-info">
        <van-icon :name="providerIcon" class="provider-icon" :style="{ color: providerColor }" />
        <span class="provider-name">{{ providerLabel }}</span>
        <van-tag v-if="config.is_default" type="primary" size="small" class="default-tag">
          {{ t('ai.default') }}
        </van-tag>
      </div>
      <van-tag :type="statusType" size="small">
        {{ statusLabel }}
      </van-tag>
    </div>

    <!-- Config Details -->
    <van-cell :title="t('ai.name')" :value="config.name" />
    <van-cell :title="t('ai.model')" :value="config.model" />
    <van-cell :title="t('ai.apiKey')" :value="maskedApiKey" />
    
    <!-- Last Test Info -->
    <van-cell 
      v-if="config.last_test_at" 
      :title="t('ai.lastTested')" 
      :value="formatDate(config.last_test_at)"
    />

    <!-- Actions -->
    <div class="card-actions">
      <van-button
        size="small"
        plain
        type="primary"
        :loading="testing"
        :loading-text="t('ai.testing')"
        @click="emit('test', config.id)"
      >
        <van-icon name="play-circle-o" />
        {{ t('ai.testConnection') }}
      </van-button>
      
      <van-button
        v-if="!config.is_default"
        size="small"
        plain
        type="success"
        @click="$emit('set-default', config.id)"
      >
        <van-icon name="star-o" />
        {{ t('ai.setDefault') }}
      </van-button>
      
      <van-button
        size="small"
        plain
        @click="$emit('edit', config)"
      >
        <van-icon name="edit" />
        {{ t('app.edit') }}
      </van-button>
      
      <van-button
        size="small"
        plain
        type="danger"
        @click="$emit('delete', config.id)"
      >
        <van-icon name="delete-o" />
        {{ t('app.delete') }}
      </van-button>
    </div>
  </van-cell-group>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const emit = defineEmits(['test', 'set-default', 'delete', 'edit'])

const props = defineProps({
  config: {
    type: Object,
    required: true
  },
  testing: {
    type: Boolean,
    default: false
  }
})

// Provider display info
const providerInfo = {
  openai: { icon: 'chat-o', color: '#10a37f', label: 'OpenAI' },
  wenxin: { icon: 'comment-o', color: '#2932e1', label: '文心一言' },
  tongyi: { icon: 'service-o', color: '#ff6a00', label: '通义千问' }
}

const providerIcon = computed(() => 
  providerInfo[props.config.provider]?.icon || 'setting-o'
)

const providerColor = computed(() => 
  providerInfo[props.config.provider]?.color || '#1989fa'
)

const providerLabel = computed(() => 
  providerInfo[props.config.provider]?.label || props.config.provider
)

// Status display
const statusType = computed(() => {
  if (props.config.status === 'active' || props.config.status === 1 || props.config.status === true || props.config.status === 'success') {
    return 'success'
  }
  return 'warning'
})

const statusLabel = computed(() => {
  if (props.config.status === 'active' || props.config.status === 1 || props.config.status === true || props.config.status === 'success') {
    return t('ai.active')
  }
  return t('ai.inactive')
})

// Mask API key for security (show only first 4 and last 4 characters)
const maskedApiKey = computed(() => {
  const key = props.config.api_key_masked || props.config.api_key || ''
  if (key.length <= 8) {
    return '••••••••'
  }
  return `${key.slice(0, 4)}${'•'.repeat(Math.min(key.length - 8, 20))}${key.slice(-4)}`
})

// Format date
function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString()
}
</script>

<style scoped>
.ai-config-card {
  margin-bottom: 12px;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
  border-bottom: 1px solid #ebedf0;
}

.provider-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.provider-icon {
  font-size: 20px;
}

.provider-name {
  font-size: 14px;
  font-weight: 600;
  color: #323233;
}

.default-tag {
  margin-left: 4px;
}

.card-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid #ebedf0;
  background-color: #fafafa;
}

.card-actions .van-button {
  display: flex;
  align-items: center;
  gap: 4px;
}

.card-actions .van-button .van-icon {
  font-size: 14px;
}
</style>
