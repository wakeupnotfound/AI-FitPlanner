<template>
  <div class="dashboard-view">
    <!-- Header -->
    <div class="dashboard-header">
      <div class="welcome-section">
        <h1 class="welcome-title">{{ t('dashboard.welcome') }}, {{ userName }}</h1>
        <p class="welcome-date">{{ formattedDate }}</p>
      </div>
    </div>

    <!-- Pull to Refresh Container -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- Loading State -->
      <LoadingSpinner v-if="loading" />

      <!-- Main Content -->
      <div v-else class="dashboard-content">
      <!-- Today's Overview Section -->
      <van-cell-group :title="t('dashboard.todayOverview')" inset>
        <!-- Today's Workout Card -->
        <div class="overview-card">
          <div class="card-header">
            <van-icon name="fire-o" class="card-icon workout-icon" />
            <span class="card-title">{{ t('dashboard.todayWorkout') }}</span>
          </div>
          <div v-if="todayWorkout && !todayWorkout.rest_day" class="card-content">
            <p class="workout-name">{{ todayWorkout.name || t('training.todayWorkout') }}</p>
            <div class="workout-meta">
              <span v-if="todayWorkout.exercises">
                {{ todayWorkout.exercises.length }} {{ t('dashboard.exercises') }}
              </span>
              <span v-if="todayWorkout.duration">
                {{ todayWorkout.duration }} {{ t('dashboard.minutes') }}
              </span>
              <span v-if="todayWorkout.calories">
                {{ todayWorkout.calories }} {{ t('dashboard.kcal') }}
              </span>
            </div>
            <van-button 
              type="primary" 
              size="small" 
              block 
              class="touch-feedback"
              @click="goToTraining"
            >
              {{ t('dashboard.viewWorkout') }}
            </van-button>
          </div>
          <div v-else-if="todayWorkout?.rest_day" class="card-content rest-day">
            <van-icon name="smile-o" size="32" color="#52c41a" />
            <p class="rest-message">{{ t('dashboard.restDay') }}</p>
            <p class="rest-hint">{{ t('dashboard.restDayMessage') }}</p>
          </div>
          <div v-else class="card-content empty">
            <p class="empty-text">{{ t('dashboard.noWorkoutToday') }}</p>
            <van-button 
              type="primary" 
              size="small" 
              plain 
              class="touch-feedback"
              @click="goToTraining"
            >
              {{ t('dashboard.startWorkout') }}
            </van-button>
          </div>
        </div>

        <!-- Today's Meals Card -->
        <div class="overview-card">
          <div class="card-header">
            <van-icon name="coupon-o" class="card-icon meals-icon" />
            <span class="card-title">{{ t('dashboard.todayMeals') }}</span>
          </div>
          <div v-if="todayMeals && todayMeals.meals?.length > 0" class="card-content">
            <div class="meals-summary">
              <div class="calories-info">
                <span class="calories-consumed">{{ consumedCalories }}</span>
                <span class="calories-divider">/</span>
                <span class="calories-target">{{ targetCalories }} {{ t('dashboard.kcal') }}</span>
              </div>
              <van-progress 
                :percentage="caloriesPercentage" 
                :stroke-width="8"
                :color="caloriesPercentage > 100 ? '#ee0a24' : '#1989fa'"
              />
            </div>
            <van-button 
              type="primary" 
              size="small" 
              block 
              class="touch-feedback"
              @click="goToNutrition"
            >
              {{ t('dashboard.viewMeals') }}
            </van-button>
          </div>
          <div v-else class="card-content empty">
            <p class="empty-text">{{ t('dashboard.noMealsToday') }}</p>
            <van-button 
              type="primary" 
              size="small" 
              plain 
              class="touch-feedback"
              @click="goToNutrition"
            >
              {{ t('dashboard.recordMeal') }}
            </van-button>
          </div>
        </div>
      </van-cell-group>

      <!-- Weekly Stats Section -->
      <van-cell-group :title="t('dashboard.weeklyStats')" inset>
        <div class="stats-grid stats-grid-responsive">
          <div class="stat-item touch-feedback">
            <div class="stat-value">{{ weeklyStats.workouts }}</div>
            <div class="stat-label">{{ t('dashboard.workoutsCompleted') }}</div>
          </div>
          <div class="stat-item touch-feedback">
            <div class="stat-value">{{ weeklyStats.duration }}</div>
            <div class="stat-label">{{ t('dashboard.minutesTrained') }}</div>
          </div>
          <div class="stat-item touch-feedback">
            <div class="stat-value">{{ weeklyStats.calories }}</div>
            <div class="stat-label">{{ t('dashboard.caloriesBurned') }}</div>
          </div>
          <div class="stat-item touch-feedback">
            <div class="stat-value">{{ weeklyStats.meals }}</div>
            <div class="stat-label">{{ t('dashboard.mealsRecorded') }}</div>
          </div>
        </div>
      </van-cell-group>

      <!-- Progress Overview Section -->
      <van-cell-group :title="t('dashboard.progressOverview')" inset>
        <div v-if="hasGoals" class="progress-section">
          <ProgressCard 
            v-if="progress && progress.comparison"
            :progress="progress"
          />
          <van-button 
            type="default" 
            size="small" 
            block 
            @click="goToStatistics"
          >
            {{ t('dashboard.viewStatistics') }}
          </van-button>
        </div>
        <div v-else class="no-goals">
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
import { ref, computed, onMounted, defineAsyncComponent } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useTrainingStore } from '../stores/training'
import { useNutritionStore } from '../stores/nutrition'
import { useStatisticsStore } from '../stores/statistics'
import { useUserStore } from '../stores/user'
import { useAuthStore } from '../stores/auth'
import { useLocale } from '../composables/useLocale'
import NavigationBar from '../components/common/NavigationBar.vue'
import LoadingSpinner from '../components/common/LoadingSpinner.vue'

// Lazy load heavy components
const ProgressCard = defineAsyncComponent(() => import('../components/statistics/ProgressCard.vue'))

const { t, locale } = useI18n()
const { formatDate } = useLocale()
const router = useRouter()

// Stores
const trainingStore = useTrainingStore()
const nutritionStore = useNutritionStore()
const statisticsStore = useStatisticsStore()
const userStore = useUserStore()
const authStore = useAuthStore()

// State
const loading = ref(true)
const refreshing = ref(false)

// Computed
const userName = computed(() => {
  return userStore.profile?.nickname || 
         userStore.profile?.username || 
         authStore.user?.nickname ||
         authStore.user?.username ||
         t('dashboard.welcome')
})

const formattedDate = computed(() => {
  return formatDate(new Date(), {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const todayWorkout = computed(() => trainingStore.todayWorkout)
const todayMeals = computed(() => nutritionStore.todayMeals)
const progress = computed(() => statisticsStore.progress)
const hasGoals = computed(() => !!userStore.goals)

const consumedCalories = computed(() => {
  if (!todayMeals.value?.meals) return 0
  return todayMeals.value.meals.reduce((sum, meal) => sum + (meal.total_calories || 0), 0)
})

const targetCalories = computed(() => {
  return nutritionStore.currentPlan?.daily_calories || 2000
})

const caloriesPercentage = computed(() => {
  if (targetCalories.value === 0) return 0
  return Math.round((consumedCalories.value / targetCalories.value) * 100)
})

const weeklyStats = computed(() => {
  const stats = statisticsStore.trainingStats
  return {
    workouts: stats?.total_workouts || 0,
    duration: stats?.total_duration_minutes || 0,
    calories: stats?.total_calories || 0,
    meals: 0 // This would need to come from nutrition stats
  }
})

// Methods
const goToTraining = () => router.push('/training')
const goToNutrition = () => router.push('/nutrition')
const goToStatistics = () => router.push('/statistics')
const goToProfile = () => router.push('/profile')

const onRefresh = async () => {
  try {
    await loadDashboardData()
  } finally {
    refreshing.value = false
  }
}

const loadDashboardData = async () => {
  loading.value = true
  try {
    await Promise.all([
      trainingStore.fetchTodayWorkout().catch(() => {}),
      nutritionStore.fetchTodayMeals().catch(() => {}),
      statisticsStore.fetchDashboardSummary().catch(() => {}),
      userStore.fetchGoals().catch(() => {})
    ])
  } catch (error) {
    console.error('Failed to load dashboard data:', error)
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.dashboard-view {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 60px;
}

.dashboard-header {
  background: linear-gradient(135deg, #1989fa 0%, #07c160 100%);
  padding: 20px 16px;
  color: white;
}

.welcome-section {
  padding-top: 10px;
}

.welcome-title {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 4px 0;
}

.welcome-date {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

.dashboard-content {
  padding: 12px 0;
}

.overview-card {
  padding: 16px;
  border-bottom: 1px solid #ebedf0;
}

.overview-card:last-child {
  border-bottom: none;
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.card-icon {
  font-size: 24px;
  margin-right: 8px;
}

.workout-icon {
  color: #ff6b6b;
}

.meals-icon {
  color: #ffa940;
}

.card-title {
  font-size: 16px;
  font-weight: 500;
  color: #323233;
}

.card-content {
  padding: 8px 0;
}

.card-content.empty {
  text-align: center;
  padding: 16px 0;
}

.card-content.rest-day {
  text-align: center;
  padding: 16px 0;
}

.empty-text {
  color: #969799;
  font-size: 14px;
  margin-bottom: 12px;
}

.rest-message {
  font-size: 16px;
  font-weight: 500;
  color: #52c41a;
  margin: 8px 0 4px;
}

.rest-hint {
  font-size: 14px;
  color: #969799;
  margin: 0;
}

.workout-name {
  font-size: 15px;
  font-weight: 500;
  color: #323233;
  margin: 0 0 8px 0;
}

.workout-meta {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: #969799;
  margin-bottom: 12px;
}

.meals-summary {
  margin-bottom: 12px;
}

.calories-info {
  display: flex;
  align-items: baseline;
  margin-bottom: 8px;
}

.calories-consumed {
  font-size: 24px;
  font-weight: 600;
  color: #1989fa;
}

.calories-divider {
  margin: 0 4px;
  color: #969799;
}

.calories-target {
  font-size: 14px;
  color: #969799;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  padding: 16px;
}

.stat-item {
  text-align: center;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1989fa;
}

.stat-label {
  font-size: 12px;
  color: #969799;
  margin-top: 4px;
}

.progress-section {
  padding: 16px;
}

.no-goals {
  padding: 16px;
}
</style>
