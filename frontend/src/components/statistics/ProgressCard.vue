<template>
  <div class="progress-card">
    <!-- Overall Progress -->
    <div class="overall-progress">
      <div class="progress-circle">
        <van-circle
          v-model:current-rate="animatedPercentage"
          :rate="progress.percentage"
          :speed="100"
          :stroke-width="60"
          :color="progressColor"
          :text="progressText"
          size="120px"
        />
      </div>
      <div class="progress-label">
        <span class="label-text">{{ t('dashboard.goalProgress') }}</span>
      </div>
    </div>

    <!-- Individual Metrics -->
    <div v-if="hasComparison" class="metrics-list">
      <!-- Weight Progress -->
      <div v-if="comparison.weight" class="metric-item">
        <div class="metric-header">
          <span class="metric-name">{{ t('statistics.weight') }}</span>
          <span class="metric-percentage">{{ Math.round(comparison.weight.percentage) }}%</span>
        </div>
        <van-progress 
          :percentage="comparison.weight.percentage" 
          :stroke-width="6"
          :color="getProgressColor(comparison.weight.percentage)"
          :show-pivot="false"
        />
        <div class="metric-values">
          <span class="value-item">
            <span class="value-label">{{ t('statistics.current') }}:</span>
            <span class="value-number">{{ comparison.weight.current }} {{ t('statistics.kg') }}</span>
          </span>
          <span class="value-item">
            <span class="value-label">{{ t('statistics.target') }}:</span>
            <span class="value-number">{{ comparison.weight.target }} {{ t('statistics.kg') }}</span>
          </span>
        </div>
      </div>

      <!-- Body Fat Progress -->
      <div v-if="comparison.body_fat" class="metric-item">
        <div class="metric-header">
          <span class="metric-name">{{ t('statistics.bodyFat') }}</span>
          <span class="metric-percentage">{{ Math.round(comparison.body_fat.percentage) }}%</span>
        </div>
        <van-progress 
          :percentage="comparison.body_fat.percentage" 
          :stroke-width="6"
          :color="getProgressColor(comparison.body_fat.percentage)"
          :show-pivot="false"
        />
        <div class="metric-values">
          <span class="value-item">
            <span class="value-label">{{ t('statistics.current') }}:</span>
            <span class="value-number">{{ comparison.body_fat.current }}{{ t('statistics.percent') }}</span>
          </span>
          <span class="value-item">
            <span class="value-label">{{ t('statistics.target') }}:</span>
            <span class="value-number">{{ comparison.body_fat.target }}{{ t('statistics.percent') }}</span>
          </span>
        </div>
      </div>

      <!-- Muscle Mass Progress -->
      <div v-if="comparison.muscle_mass" class="metric-item">
        <div class="metric-header">
          <span class="metric-name">{{ t('statistics.muscleMass') }}</span>
          <span class="metric-percentage">{{ Math.round(comparison.muscle_mass.percentage) }}%</span>
        </div>
        <van-progress 
          :percentage="comparison.muscle_mass.percentage" 
          :stroke-width="6"
          :color="getProgressColor(comparison.muscle_mass.percentage)"
          :show-pivot="false"
        />
        <div class="metric-values">
          <span class="value-item">
            <span class="value-label">{{ t('statistics.current') }}:</span>
            <span class="value-number">{{ comparison.muscle_mass.current }} {{ t('statistics.kg') }}</span>
          </span>
          <span class="value-item">
            <span class="value-label">{{ t('statistics.target') }}:</span>
            <span class="value-number">{{ comparison.muscle_mass.target }} {{ t('statistics.kg') }}</span>
          </span>
        </div>
      </div>
    </div>

    <!-- No Comparison Data -->
    <div v-else class="no-comparison">
      <p class="no-data-text">{{ t('statistics.insufficientData') }}</p>
      <p class="no-data-hint">{{ t('statistics.insufficientDataHint') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  progress: {
    type: Object,
    required: true,
    default: () => ({
      percentage: 0,
      comparison: {}
    })
  }
})

// Animated percentage for circle
const animatedPercentage = ref(0)

// Watch for progress changes to animate
watch(() => props.progress.percentage, (newVal) => {
  animatedPercentage.value = newVal
}, { immediate: true })

// Computed
const comparison = computed(() => props.progress.comparison || {})

const hasComparison = computed(() => {
  return comparison.value && Object.keys(comparison.value).length > 0
})

const progressColor = computed(() => {
  const percentage = props.progress.percentage
  if (percentage >= 80) return '#52c41a'
  if (percentage >= 50) return '#1989fa'
  if (percentage >= 25) return '#ffa940'
  return '#ff6b6b'
})

const progressText = computed(() => {
  return `${props.progress.percentage}%`
})

// Methods
const getProgressColor = (percentage) => {
  if (percentage >= 80) return '#52c41a'
  if (percentage >= 50) return '#1989fa'
  if (percentage >= 25) return '#ffa940'
  return '#ff6b6b'
}
</script>

<style scoped>
.progress-card {
  background: white;
  border-radius: 8px;
}

.overall-progress {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 0;
}

.progress-circle {
  margin-bottom: 8px;
}

.progress-label {
  text-align: center;
}

.label-text {
  font-size: 14px;
  color: #969799;
}

.metrics-list {
  padding: 0 16px 16px;
}

.metric-item {
  padding: 12px 0;
  border-top: 1px solid #ebedf0;
}

.metric-item:first-child {
  border-top: none;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.metric-name {
  font-size: 14px;
  font-weight: 500;
  color: #323233;
}

.metric-percentage {
  font-size: 14px;
  font-weight: 600;
  color: #1989fa;
}

.metric-values {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
}

.value-item {
  font-size: 12px;
}

.value-label {
  color: #969799;
  margin-right: 4px;
}

.value-number {
  color: #323233;
  font-weight: 500;
}

.no-comparison {
  text-align: center;
  padding: 16px;
}

.no-data-text {
  font-size: 14px;
  color: #969799;
  margin: 0 0 4px 0;
}

.no-data-hint {
  font-size: 12px;
  color: #c8c9cc;
  margin: 0;
}
</style>
