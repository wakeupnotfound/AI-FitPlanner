<template>
  <van-cell-group inset class="plan-card">
    <van-cell :border="false">
      <template #title>
        <div class="plan-header">
          <h3 class="plan-name">{{ plan.name }}</h3>
          <van-tag :type="statusType" size="medium">
            {{ t(`training.${plan.status}`) }}
          </van-tag>
        </div>
      </template>
      <template #label>
        <div class="plan-meta">
          <span class="meta-item">
            <van-icon name="calendar-o" />
            {{ formatDateRange(plan.start_date, plan.end_date) }}
          </span>
          <span class="meta-item">
            <van-icon name="clock-o" />
            {{ plan.total_weeks || plan.duration_weeks }} {{ t('training.week') }}
          </span>
        </div>
      </template>
    </van-cell>
    
    <van-grid :column-num="3" :border="false" class="plan-stats">
      <van-grid-item>
        <div class="stat-value">{{ t(`training.goals.${plan.goal}`) }}</div>
        <div class="stat-label">{{ t('training.goal') }}</div>
      </van-grid-item>
      <van-grid-item>
        <div class="stat-value">{{ t(`training.difficultyLevels.${plan.difficulty_level || plan.difficulty}`) }}</div>
        <div class="stat-label">{{ t('training.difficulty') }}</div>
      </van-grid-item>
      <van-grid-item>
        <div class="stat-value">{{ completedWeeks }}/{{ plan.total_weeks || plan.duration_weeks }}</div>
        <div class="stat-label">{{ t('training.completed') }}</div>
      </van-grid-item>
    </van-grid>

    <van-cell :border="false" v-if="showActions">
      <div class="plan-actions">
        <van-button 
          size="small" 
          type="primary" 
          plain
          @click="$emit('view', plan)"
        >
          {{ t('training.viewPlan') }}
        </van-button>
        <van-button 
          v-if="plan.status === 'active'"
          size="small" 
          type="primary"
          @click="$emit('start', plan)"
        >
          {{ t('training.startWorkout') }}
        </van-button>
      </div>
    </van-cell>
  </van-cell-group>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  plan: {
    type: Object,
    required: true
  },
  showActions: {
    type: Boolean,
    default: true
  }
})

defineEmits(['view', 'start'])

const statusType = computed(() => {
  switch (props.plan.status) {
    case 'active':
      return 'primary'
    case 'completed':
      return 'success'
    case 'paused':
      return 'warning'
    default:
      return 'default'
  }
})

const completedWeeks = computed(() => {
  if (!props.plan.start_date) return 0
  const start = new Date(props.plan.start_date)
  const now = new Date()
  const diffTime = now - start
  const diffWeeks = Math.floor(diffTime / (1000 * 60 * 60 * 24 * 7))
  const totalWeeks = props.plan.total_weeks || props.plan.duration_weeks || 0
  return Math.min(Math.max(0, diffWeeks), totalWeeks)
})

const formatDateRange = (startDate, endDate) => {
  if (!startDate) return ''
  const start = new Date(startDate).toLocaleDateString()
  const end = endDate ? new Date(endDate).toLocaleDateString() : ''
  return end ? `${start} - ${end}` : start
}
</script>

<style scoped>
.plan-card {
  margin: 12px 16px;
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.plan-name {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--van-text-color);
}

.plan-meta {
  display: flex;
  gap: 16px;
  margin-top: 8px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.plan-stats {
  padding: 8px 0;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--van-primary-color);
}

.stat-label {
  font-size: 12px;
  color: var(--van-text-color-3);
  margin-top: 4px;
}

.plan-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

:deep(.van-grid-item__content) {
  padding: 12px 8px;
}
</style>
