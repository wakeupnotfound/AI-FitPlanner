<template>
  <div class="statistics-view">
    <!-- Header -->
    <van-nav-bar 
      :title="t('statistics.title')"
      left-arrow
      @click-left="goBack"
    />

    <!-- Pull to Refresh Container -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- Loading State -->
      <LoadingSpinner v-if="loading" />

      <!-- Main Content -->
      <div v-else class="statistics-content">
      <!-- Training Statistics Section -->
      <van-cell-group :title="t('statistics.trainingStats')" inset>
        <div class="stats-summary">
          <div class="stat-card">
            <div class="stat-icon">
              <van-icon name="fire-o" size="24" color="#ff6b6b" />
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ trainingStats.totalWorkouts }}</span>
              <span class="stat-label">{{ t('statistics.totalWorkouts') }}</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon">
              <van-icon name="clock-o" size="24" color="#1989fa" />
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ formatDuration(trainingStats.totalDuration) }}</span>
              <span class="stat-label">{{ t('statistics.totalDuration') }}</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon">
              <van-icon name="chart-trending-o" size="24" color="#07c160" />
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ trainingStats.avgDuration }}</span>
              <span class="stat-label">{{ t('statistics.avgDuration') }}</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon">
              <van-icon name="calendar-o" size="24" color="#ffa940" />
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ trainingStats.frequency }}</span>
              <span class="stat-label">{{ t('statistics.workoutFrequency') }}</span>
            </div>
          </div>
        </div>

        <!-- Training Chart -->
        <ChartWidget
          :title="t('statistics.trainingStats')"
          type="bar"
          :data="trainingChartData"
          :loading="chartLoading"
          color="#1989fa"
          :unit="t('statistics.minutes')"
          :default-range="selectedRange"
          @range-change="onRangeChange"
        />
      </van-cell-group>

      <!-- Body Trends Section -->
      <van-cell-group :title="t('statistics.bodyTrends')" inset>
        <!-- Metric Tabs -->
        <van-tabs v-model:active="activeMetric" shrink>
          <van-tab :title="t('statistics.weight')" name="weight" />
          <van-tab :title="t('statistics.bodyFat')" name="body_fat" />
          <van-tab :title="t('statistics.muscleMass')" name="muscle_mass" />
        </van-tabs>

        <!-- Body Trends Chart -->
        <div v-if="hasBodyData" class="body-chart-container">
          <ChartWidget
            :key="activeMetric"
            :title="getMetricTitle(activeMetric)"
            type="line"
            :data="bodyTrendData"
            :loading="chartLoading"
            :color="getMetricColor(activeMetric)"
            :unit="getMetricUnit(activeMetric)"
            show-legend
            :legend-label="getMetricTitle(activeMetric)"
            :default-range="selectedRange"
            @range-change="onRangeChange"
          />

          <!-- Current Value Display -->
          <div class="current-value-card">
            <div class="value-header">
              <span class="value-label">{{ t('statistics.current') }}</span>
              <span class="value-number">
                {{ getCurrentValue(activeMetric) }} {{ getMetricUnit(activeMetric) }}
              </span>
            </div>
            <div v-if="hasGoal(activeMetric)" class="value-target">
              <span class="target-label">{{ t('statistics.target') }}:</span>
              <span class="target-value">
                {{ getTargetValue(activeMetric) }} {{ getMetricUnit(activeMetric) }}
              </span>
            </div>
          </div>
        </div>

        <!-- No Body Data State -->
        <div v-else class="no-data-state">
          <van-empty 
            :description="t('statistics.noData')"
            image="search"
          >
            <template #description>
              <p class="empty-text">{{ t('statistics.noDataHint') }}</p>
            </template>
            <van-button 
              type="primary" 
              size="small" 
              @click="goToProfile"
            >
              {{ t('profile.addBodyData') }}
            </van-button>
          </van-empty>
        </div>
      </van-cell-group>

      <!-- Progress Section -->
      <van-cell-group :title="t('statistics.progress')" inset>
        <div v-if="hasGoals" class="progress-section">
          <ProgressCard :progress="progress" />
        </div>
        <div v-else class="no-goals-state">
          <van-empty 
            :description="t('dashboard.noGoalsSet')"
            image="search"
          >
            <van-button 
              type="primary" 
              size="small" 
              @click="goToProfile"
            >
              {{ t('dashboard.setGoals') }}
            </van-button>
          </van-empty>
        </div>
      </van-cell-group>
    </div>
    </van-pull-refresh>

    <!-- Navigation Bar -->
    <NavigationBar />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, defineAsyncComponent } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useStatisticsStore } from '../stores/statistics'
import { useUserStore } from '../stores/user'
import { statisticsService } from '../services/statistics.service'
import NavigationBar from '../components/common/NavigationBar.vue'
import LoadingSpinner from '../components/common/LoadingSpinner.vue'

// Lazy load heavy components
const ChartWidget = defineAsyncComponent(() => import('../components/statistics/ChartWidget.vue'))
const ProgressCard = defineAsyncComponent(() => import('../components/statistics/ProgressCard.vue'))

const { t } = useI18n()
const router = useRouter()

// Stores
const statisticsStore = useStatisticsStore()
const userStore = useUserStore()

// State
const loading = ref(true)
const chartLoading = ref(false)
const refreshing = ref(false)
const selectedRange = ref('month')
const activeMetric = ref('weight')

// Computed
const trainingStats = computed(() => {
  const stats = statisticsStore.trainingStats || {}
  const startDate = stats.start_date ? new Date(stats.start_date) : null
  const endDate = stats.end_date ? new Date(stats.end_date) : null
  let frequency = stats.workout_frequency || 0

  if (!frequency && startDate && endDate) {
    const msPerDay = 24 * 60 * 60 * 1000
    const days = Math.max(1, Math.round((endDate - startDate) / msPerDay) + 1)
    const weeks = days / 7
    frequency = stats.total_workouts ? (stats.total_workouts / weeks) : 0
    frequency = Math.round(frequency * 10) / 10
  }

  return {
    totalWorkouts: stats.total_workouts || 0,
    totalDuration: stats.total_duration_minutes || stats.total_duration || 0,
    avgDuration: stats.average_duration_minutes || stats.average_duration || 0,
    frequency,
    caloriesBurned: stats.total_calories || 0
  }
})

const trainingChartData = computed(() => {
  const trends = statisticsStore.trainingTrends
  if (trends?.data_points?.length) {
    return trends.data_points.map(item => ({
      date: item.start_date,
      value: item.total_duration_minutes || 0,
      label: item.period_label || item.start_date
    }))
  }

  const stats = statisticsStore.trainingStats
  if (!stats?.total_duration_minutes) {
    return []
  }

  return [
    {
      date: stats.end_date || '',
      value: stats.total_duration_minutes,
      label: stats.end_date || t('statistics.totalDuration')
    }
  ]
})

const bodyTrendData = computed(() => {
  const items = userStore.bodyDataHistory || []
  if (!items.length) return []

  const keyMap = {
    weight: 'weight',
    body_fat: 'body_fat_percentage',
    muscle_mass: 'muscle_percentage'
  }
  const valueKey = keyMap[activeMetric.value]

  const filtered = items.filter(item => item[valueKey] !== null && item[valueKey] !== undefined)
  if (!filtered.length) {
    return []
  }

  return filtered
    .slice()
    .reverse()
    .map(item => ({
      date: item.measurement_date,
      value: Number(item[valueKey]) || 0,
      label: item.measurement_date
    }))
})

const hasBodyData = computed(() => {
  return bodyTrendData.value.length > 0 || userStore.bodyData?.length > 0
})

const hasGoals = computed(() => !!userStore.goals)

const progress = computed(() => {
  if (!hasGoals.value || !userStore.latestBodyData) {
    return { percentage: 0, comparison: {} }
  }
  
  return statisticsService.calculateProgress(
    userStore.latestBodyData,
    userStore.goals,
    userStore.bodyDataHistory
  )
})

// Methods
const goBack = () => router.back()
const goToProfile = () => router.push('/profile')

const onRefresh = async () => {
  try {
    await loadInitialData()
  } finally {
    refreshing.value = false
  }
}

const formatDuration = (minutes) => {
  if (!minutes) return '0'
  if (minutes < 60) return `${minutes} ${t('statistics.minutes')}`
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return `${hours}${t('statistics.hours')} ${mins}${t('statistics.minutes')}`
}

const getMetricTitle = (metric) => {
  const titles = {
    weight: t('statistics.weight'),
    body_fat: t('statistics.bodyFat'),
    muscle_mass: t('statistics.muscleMass')
  }
  return titles[metric] || metric
}

const getMetricColor = (metric) => {
  const colors = {
    weight: '#1989fa',
    body_fat: '#ff6b6b',
    muscle_mass: '#07c160'
  }
  return colors[metric] || '#1989fa'
}

const getMetricUnit = (metric) => {
  const units = {
    weight: t('statistics.kg'),
    body_fat: t('statistics.percent'),
    muscle_mass: t('statistics.percent')
  }
  return units[metric] || ''
}

const getCurrentValue = (metric) => {
  const latestData = userStore.latestBodyData
  if (!latestData) return '-'
  
  const keyMap = {
    weight: 'weight',
    body_fat: 'body_fat_percentage',
    muscle_mass: 'muscle_percentage'
  }
  
  return latestData[keyMap[metric]] || '-'
}

const hasGoal = (metric) => {
  if (!userStore.goals) return false
  const goalKey = `target_${metric}`
  return !!userStore.goals[goalKey]
}

const getTargetValue = (metric) => {
  if (!userStore.goals) return '-'
  const goalKey = `target_${metric}`
  return userStore.goals[goalKey] || '-'
}

const onRangeChange = async (range) => {
  selectedRange.value = range
  await loadChartData(range)
}

const loadChartData = async (range) => {
  chartLoading.value = true
  try {
    const dateRange = statisticsService.getDateRange(range)
    await Promise.all([
      statisticsStore.fetchStatistics(dateRange),
      statisticsStore.fetchTrends(dateRange)
    ])
  } catch (error) {
    console.error('Failed to load chart data:', error)
  } finally {
    chartLoading.value = false
  }
}

const loadInitialData = async () => {
  loading.value = true
  try {
    const dateRange = statisticsService.getDateRange(selectedRange.value)
    await Promise.all([
      statisticsStore.fetchStatistics(dateRange),
      statisticsStore.fetchTrends(dateRange),
      userStore.fetchBodyData().catch(() => {}),
      userStore.fetchGoals().catch(() => {})
    ])
  } catch (error) {
    console.error('Failed to load statistics:', error)
  } finally {
    loading.value = false
  }
}

// Watch for metric changes to reload data if needed
watch(activeMetric, () => {
  // Data is already loaded, just switch the view
})

// Lifecycle
onMounted(() => {
  loadInitialData()
})
</script>

<style scoped>
.statistics-view {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 60px;
}

.statistics-content {
  padding: 12px 0;
}

.stats-summary {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.stat-icon {
  margin-right: 12px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #323233;
}

.stat-label {
  font-size: 11px;
  color: #969799;
  margin-top: 2px;
}

.body-chart-container {
  padding: 16px;
}

.current-value-card {
  margin-top: 16px;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.value-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.value-label {
  font-size: 14px;
  color: #969799;
}

.value-number {
  font-size: 20px;
  font-weight: 600;
  color: #1989fa;
}

.value-target {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #ebedf0;
}

.target-label {
  font-size: 12px;
  color: #969799;
}

.target-value {
  font-size: 14px;
  font-weight: 500;
  color: #323233;
}

.no-data-state {
  padding: 20px;
}

.empty-text {
  font-size: 12px;
  color: #969799;
  margin: 8px 0;
}

.progress-section {
  padding: 16px;
}

.no-goals-state {
  padding: 20px;
}
</style>
