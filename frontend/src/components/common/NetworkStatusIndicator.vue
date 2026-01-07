<template>
  <div class="network-status-indicator">
    <!-- Offline Banner -->
    <van-notice-bar
      v-if="!isOnline"
      mode="closeable"
      color="#ed6a0c"
      background="#fffbe8"
      left-icon="warning-o"
      :scrollable="false"
    >
      <template #default>
        {{ t('network.offline') }}
        <span v-if="queueSize > 0" class="queue-info">
          ({{ t('network.queuedOperations', { count: queueSize }) }})
        </span>
      </template>
    </van-notice-bar>

    <!-- Syncing Banner -->
    <van-notice-bar
      v-if="isOnline && isSyncing"
      color="#1989fa"
      background="#ecf9ff"
      left-icon="clock-o"
      :scrollable="false"
    >
      {{ t('network.syncing') }}
    </van-notice-bar>

    <!-- Data Freshness Indicator (for cached data) -->
    <div v-if="showFreshnessIndicator && lastUpdateTime" class="freshness-indicator">
      <van-tag type="default" size="small">
        <van-icon name="clock-o" />
        {{ formatLastUpdate(lastUpdateTime) }}
      </van-tag>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNetworkStatus } from '@/composables/useNetworkStatus'

const props = defineProps({
  showFreshnessIndicator: {
    type: Boolean,
    default: false
  },
  lastUpdateTime: {
    type: [Date, String, Number],
    default: null
  }
})

const { t } = useI18n()
const { isOnline, queueSize, isSyncing } = useNetworkStatus()

/**
 * Format last update time
 */
const formatLastUpdate = (time) => {
  if (!time) return ''
  
  const date = time instanceof Date ? time : new Date(time)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  
  if (diffMins < 1) {
    return t('network.justNow')
  } else if (diffMins < 60) {
    return t('network.minutesAgo', { minutes: diffMins })
  } else if (diffMins < 1440) {
    const hours = Math.floor(diffMins / 60)
    return t('network.hoursAgo', { hours })
  } else {
    const days = Math.floor(diffMins / 1440)
    return t('network.daysAgo', { days })
  }
}
</script>

<style scoped>
.network-status-indicator {
  position: relative;
}

.queue-info {
  margin-left: 8px;
  font-weight: 500;
}

.freshness-indicator {
  padding: 8px 16px;
  background: #f7f8fa;
  border-bottom: 1px solid #ebedf0;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.freshness-indicator .van-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
