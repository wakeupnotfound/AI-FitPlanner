<template>
  <div class="training-view">
    <!-- Header -->
    <van-nav-bar :title="t('training.title')" />

    <!-- Tabs -->
    <van-tabs v-model:active="activeTab" sticky>
      <van-tab :title="t('training.todayWorkout')" name="today" />
      <van-tab :title="t('training.plans')" name="plans" />
      <van-tab :title="t('training.history')" name="history" />
    </van-tabs>

    <!-- Tab Content -->
    <div class="tab-content">
      <!-- Today's Workout Tab -->
      <div v-show="activeTab === 'today'" class="tab-panel">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
          <!-- Loading state -->
          <div v-if="loading" class="loading-container">
            <LoadingSpinner />
          </div>

          <!-- No AI Config Warning -->
          <van-empty
            v-else-if="!hasAIConfig"
            image="search"
            :description="t('ai.noConfigsHint')"
          >
            <van-button type="primary" size="small" to="/ai-config">
              {{ t('ai.addConfig') }}
            </van-button>
          </van-empty>

          <!-- Plan Generation Progress -->
          <div v-else-if="isGeneratingDisplay" class="generating-container">
            <van-loading type="spinner" size="48" color="var(--van-primary-color)" />
            <h3>{{ t('training.generating') }}</h3>
            <p>{{ t('training.generatingHint') }}</p>
            <van-progress
              :percentage="generationProgress"
              stroke-width="8"
              color="var(--van-primary-color)"
            />
          </div>

          <!-- No Plan State -->
          <van-empty
            v-else-if="!currentPlan"
            image="search"
            :description="t('training.noPlanHint')"
          >
            <van-button type="primary" @click="showAssessment = true">
              {{ t('training.generatePlan') }}
            </van-button>
          </van-empty>

          <!-- Rest Day -->
          <div v-else-if="isRestDay" class="rest-day">
            <van-empty
              image="https://fastly.jsdelivr.net/npm/@vant/assets/custom-empty-image.png"
              :description="t('training.restDayHint')"
            >
              <template #image>
                <van-icon name="smile-o" size="80" color="var(--van-primary-color)" />
              </template>
            </van-empty>
            <h3 class="rest-day-title">{{ t('training.restDay') }}</h3>
          </div>

          <!-- Today's Workout -->
          <div v-else-if="todayWorkout">
            <WorkoutCard
              :workout="todayWorkout"
              :show-exercises="true"
              @start="startWorkout"
              @view="viewWorkoutDetails"
            />

            <!-- Exercise List -->
            <van-cell-group inset :title="t('training.exercises')" class="exercise-list">
              <van-cell
                v-for="(exercise, index) in todayWorkout.exercises"
                :key="index"
                :class="{ 'completed-exercise': exercise.completed }"
                clickable
                @click="selectExercise(exercise, index)"
              >
                <template #title>
                  <div class="exercise-title">
                    <van-checkbox
                      v-model="exercise.completed"
                      shape="square"
                      @click.stop
                    />
                    <span :class="{ 'line-through': exercise.completed }">
                      {{ exercise.name }}
                    </span>
                  </div>
                </template>
                <template #label>
                  <div class="exercise-details">
                    <span>{{ exercise.sets }} {{ t('training.sets') }} × {{ exercise.reps }} {{ t('training.reps') }}</span>
                    <span v-if="exercise.weight">{{ exercise.weight }}</span>
                    <span v-if="exercise.rest">{{ t('training.rest') }}: {{ exercise.rest }}</span>
                  </div>
                  <div v-if="exercise.safety_notes" class="exercise-safety">
                    {{ exercise.safety_notes }}
                  </div>
                </template>
                <template #right-icon>
                  <van-icon name="arrow" />
                </template>
              </van-cell>
            </van-cell-group>

            <!-- Complete Workout Button -->
            <div class="action-button">
              <van-button
                type="primary"
                block
                :disabled="!canCompleteWorkout"
                @click="showRecordForm = true"
              >
                {{ t('training.completeWorkout') }}
              </van-button>
            </div>
          </div>

          <!-- No Workout Today -->
          <van-empty
            v-else
            image="search"
            :description="t('training.noWorkoutToday')"
          />
        </van-pull-refresh>
      </div>

      <!-- Plans Tab -->
      <div v-show="activeTab === 'plans'" class="tab-panel">
        <van-pull-refresh v-model="refreshingPlans" @refresh="onRefreshPlans">
          <!-- Loading state -->
          <div v-if="loadingPlans" class="loading-container">
            <LoadingSpinner />
          </div>

          <!-- Plan Generation Progress -->
          <div v-else-if="isGeneratingDisplay" class="generating-container">
            <van-loading type="spinner" size="48" color="var(--van-primary-color)" />
            <h3>{{ t('training.generating') }}</h3>
            <p>{{ t('training.generatingHint') }}</p>
            <van-progress
              :percentage="generationProgress"
              stroke-width="8"
              color="var(--van-primary-color)"
            />
          </div>

          <!-- Generation Error -->
          <div v-else-if="generationError" class="generation-error">
            <van-empty
              image="error"
              :description="t('training.planGenerateFailed')"
            >
              <div class="error-detail">{{ generationError }}</div>
              <van-button type="primary" plain @click="showAssessment = true">
                {{ t('training.generatePlan') }}
              </van-button>
            </van-empty>
          </div>

          <!-- No Plans -->
          <van-empty
            v-else-if="plans.length === 0"
            image="search"
            :description="t('training.noPlanHint')"
          >
            <van-button type="primary" @click="showAssessment = true">
              {{ t('training.generatePlan') }}
            </van-button>
          </van-empty>

          <!-- Plans List -->
          <div v-else>
            <TrainingPlanCard
              v-for="plan in plans"
              :key="plan.id"
              :plan="plan"
              :completed-count="getPlanCompletedCount(plan)"
              :total-count="getPlanTotalCount(plan)"
              @view="viewPlanDetails"
              @start="startPlanWorkout"
            />

            <!-- Generate New Plan Button -->
            <div class="action-button">
              <van-button
                type="primary"
                plain
                block
                @click="showAssessment = true"
              >
                {{ t('training.generatePlan') }}
              </van-button>
            </div>
          </div>
        </van-pull-refresh>
      </div>

      <!-- History Tab -->
      <div v-show="activeTab === 'history'" class="tab-panel">
        <van-pull-refresh v-model="refreshingHistory" @refresh="onRefreshHistory">
          <!-- Loading state -->
          <div v-if="loadingHistory" class="loading-container">
            <LoadingSpinner />
          </div>

          <!-- No History -->
          <van-empty
            v-else-if="history.length === 0"
            image="search"
            :description="t('app.noData')"
          />

          <!-- History Calendar -->
          <div v-else>
            <van-calendar
              v-model:show="showCalendar"
              type="single"
              :min-date="minCalendarDate"
              :max-date="maxCalendarDate"
              :formatter="calendarFormatter"
              @confirm="onCalendarConfirm"
            />

            <van-cell-group inset>
              <van-cell
                :title="t('training.history')"
                is-link
                @click="showCalendar = true"
              >
                <template #value>
                  {{ selectedDateDisplay }}
                </template>
              </van-cell>
            </van-cell-group>

            <!-- History List -->
            <div class="history-list">
              <van-cell-group
                v-for="record in filteredHistory"
                :key="record.id"
                inset
                class="history-item"
              >
                <van-cell :border="false" is-link @click="viewRecordDetails(record)">
                  <template #title>
                    <div class="history-header">
                      <span class="history-date">{{ formatDate(record.workout_date) }}</span>
                      <van-tag type="success" size="small">{{ t('training.completed') }}</van-tag>
                    </div>
                  </template>
                  <template #label>
                    <div class="history-meta">
                      <span>{{ record.duration_minutes }} {{ t('training.minutes') }}</span>
                      <span v-if="record.performance_data?.estimated_calories">
                        {{ record.performance_data.estimated_calories }} {{ t('training.kcal') }}
                      </span>
                      <span v-if="record.rating">
                        <van-rate v-model="record.rating" readonly size="12" />
                      </span>
                    </div>
                  </template>
                </van-cell>
                <van-cell
                  v-if="getRecordExercises(record).length > 0"
                  :border="false"
                >
                  <div class="history-exercises">
                    <span v-for="(ex, i) in getRecordExercises(record).slice(0, 3)" :key="i" class="exercise-tag">
                      {{ ex.exercise_name }}
                    </span>
                    <span v-if="getRecordExercises(record).length > 3" class="more-tag">
                      +{{ getRecordExercises(record).length - 3 }}
                    </span>
                  </div>
                </van-cell>
              </van-cell-group>
            </div>
          </div>
        </van-pull-refresh>
      </div>
    </div>

    <!-- Assessment Form Popup -->
    <van-popup
      v-model:show="showAssessment"
      position="bottom"
      round
      :style="{ height: '90%' }"
      closeable
    >
      <AssessmentForm
        :loading="submittingAssessment"
        @submit="handleAssessmentSubmit"
        @cancel="showAssessment = false"
      />
    </van-popup>

    <!-- Workout Record Form Popup -->
    <van-popup
      v-model:show="showRecordForm"
      position="bottom"
      round
      :style="{ height: '80%' }"
      closeable
    >
      <WorkoutRecordForm
        v-if="todayWorkout"
        :workout="todayWorkout"
        :plan-id="activePlanId"
        :loading="submittingRecord"
        @submit="handleRecordSubmit"
        @cancel="showRecordForm = false"
      />
    </van-popup>

    <!-- Plan Details Popup -->
    <van-popup
      v-model:show="showPlanDetails"
      position="bottom"
      round
      :style="{ height: '90%' }"
      closeable
    >
      <div class="plan-details-popup" v-if="selectedPlan">
        <van-nav-bar :title="selectedPlan.plan_name || selectedPlan.name" />
        <div class="plan-details-content">
          <van-cell-group inset>
            <van-cell :title="t('training.goal')" :value="formatPlanGoal(selectedPlan)" />
            <van-cell :title="t('training.difficulty')" :value="t(`training.difficultyLevels.${selectedPlan.difficulty_level || selectedPlan.difficulty}`)" />
            <van-cell :title="t('training.totalWeeks')" :value="selectedPlan.total_weeks || selectedPlan.duration_weeks" />
            <van-cell :title="t('training.status')" :value="t(`training.${selectedPlan.status}`)" />
          </van-cell-group>

          <!-- Weekly Schedule -->
          <h4 class="section-title">{{ t('training.weeklySchedule') }}</h4>
          <van-collapse v-model="activeWeeks" v-if="selectedPlan.plan_data?.weeks">
            <van-collapse-item
              v-for="week in selectedPlan.plan_data.weeks"
              :key="week.week"
              :title="`${t('training.week')} ${week.week}`"
              :name="week.week"
            >
              <div v-for="day in week.days" :key="day.day" class="day-schedule">
                <div class="day-header">
                  <span class="day-number">{{ t('training.day') }} {{ day.day }}</span>
                  <van-tag :type="day.type === 'rest' ? 'default' : 'primary'" size="small">
                    {{ t(`training.workoutTypes.${day.type}`) }}
                  </van-tag>
                </div>
                <div v-if="day.exercises" class="day-exercises">
                  <div v-for="(ex, i) in day.exercises" :key="i" class="mini-exercise">
                    {{ ex.name }} - {{ ex.sets }}×{{ ex.reps }}
                  </div>
                </div>
              </div>
            </van-collapse-item>
          </van-collapse>
        </div>
      </div>
    </van-popup>

    <!-- Record Details Popup -->
    <van-popup
      v-model:show="showRecordDetails"
      position="bottom"
      round
      :style="{ height: '80%' }"
      closeable
    >
      <div class="record-details-popup" v-if="selectedRecord">
        <van-nav-bar :title="t('training.recordWorkout')" />
        <div class="record-details-content">
          <van-cell-group inset>
            <van-cell :title="t('training.date')" :value="formatDate(selectedRecord.workout_date)" />
            <van-cell :title="t('training.workoutType')" :value="t(`training.workoutTypes.${selectedRecord.workout_type}`)" />
            <van-cell :title="t('training.duration')" :value="`${selectedRecord.duration_minutes || 0} ${t('training.minutes')}`" />
            <van-cell v-if="selectedRecord.rating" :title="t('training.rating')" :value="selectedRecord.rating" />
          </van-cell-group>

          <van-cell-group inset :title="t('training.exercises')">
            <van-cell
              v-for="(ex, idx) in getRecordExercises(selectedRecord)"
              :key="idx"
              :title="ex.exercise_name"
            >
              <template #label>
                <div class="record-exercise-meta">
                  <span>{{ ex.sets }} {{ t('training.sets') }}</span>
                  <span>{{ t('training.reps') }}: {{ ex.reps_per_set?.join(', ') || '-' }}</span>
                  <span>{{ t('training.weight') }}: {{ ex.weight_used?.join(', ') || '-' }}</span>
                </div>
                <div v-if="ex.notes" class="record-exercise-notes">
                  {{ ex.notes }}
                </div>
              </template>
            </van-cell>
          </van-cell-group>

          <van-cell-group inset v-if="selectedRecord.notes" :title="t('training.notes')">
            <van-cell :border="false" :label="selectedRecord.notes" />
          </van-cell-group>
          <van-cell-group inset v-if="selectedRecord.injury_report" :title="t('training.injuryReport')">
            <van-cell :border="false" :label="selectedRecord.injury_report" />
          </van-cell-group>
        </div>
      </div>
    </van-popup>

    <!-- Navigation Bar -->
    <NavigationBar />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, defineAsyncComponent } from 'vue'
import { useI18n } from 'vue-i18n'
import { showToast } from 'vant'
import { useTrainingStore } from '@/stores/training'
import { useAIConfigStore } from '@/stores/aiConfig'
import { trainingService } from '@/services/training.service'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import NavigationBar from '@/components/common/NavigationBar.vue'

// Lazy load heavy components
const AssessmentForm = defineAsyncComponent(() => import('@/components/training/AssessmentForm.vue'))
const TrainingPlanCard = defineAsyncComponent(() => import('@/components/training/TrainingPlanCard.vue'))
const WorkoutCard = defineAsyncComponent(() => import('@/components/training/WorkoutCard.vue'))
const WorkoutRecordForm = defineAsyncComponent(() => import('@/components/training/WorkoutRecordForm.vue'))

const { t } = useI18n()
const trainingStore = useTrainingStore()
const aiConfigStore = useAIConfigStore()

// State
const activeTab = ref('today')
const loading = ref(false)
const loadingPlans = ref(false)
const loadingHistory = ref(false)
const refreshing = ref(false)
const refreshingPlans = ref(false)
const refreshingHistory = ref(false)
const submittingAssessment = ref(false)
const submittingRecord = ref(false)

// Popups
const showAssessment = ref(false)
const showRecordForm = ref(false)
const showPlanDetails = ref(false)
const showCalendar = ref(false)
const showRecordDetails = ref(false)

// Selected items
const selectedPlan = ref(null)
const selectedRecord = ref(null)
const selectedDate = ref(new Date())
const activeWeeks = ref([1])

// Generation state
const generationTaskId = ref(null)
const generationProgress = ref(0)
const generationInterval = ref(null)
const generationError = ref('')
const isGeneratingTask = ref(false)

// Computed
const hasAIConfig = computed(() => aiConfigStore.hasDefaultConfig)
const currentPlan = computed(() => trainingStore.activePlan)
const todayWorkout = computed(() => trainingStore.getTodayWorkout)
const plans = computed(() => trainingStore.allPlans)
const history = computed(() => trainingStore.trainingHistory)
const isGenerating = computed(() => trainingStore.isGeneratingPlan)
const isGeneratingDisplay = computed(() => isGenerating.value || isGeneratingTask.value)
const activePlanId = computed(() => currentPlan.value?.id || plans.value?.[0]?.id || null)

const isRestDay = computed(() => {
  return todayWorkout.value?.type === 'rest' || todayWorkout.value?.rest_day
})

const canCompleteWorkout = computed(() => {
  if (!todayWorkout.value?.exercises) return false
  if (isRecordedToday.value) return false
  return todayWorkout.value.exercises.some(e => e.completed)
})

const minCalendarDate = computed(() => {
  const date = new Date()
  date.setMonth(date.getMonth() - 6)
  return date
})

const maxCalendarDate = computed(() => new Date())

const selectedDateDisplay = computed(() => {
  return selectedDate.value.toLocaleDateString()
})

const getLocalDateKey = (date) => {
  const localDate = new Date(date)
  const offsetMs = localDate.getTimezoneOffset() * 60 * 1000
  return new Date(localDate.getTime() - offsetMs).toISOString().slice(0, 10)
}

const normalizeRecordDate = (dateValue) => {
  if (!dateValue) return ''
  if (typeof dateValue === 'string') {
    if (dateValue.length >= 10) {
      return dateValue.slice(0, 10)
    }
    return dateValue
  }
  return getLocalDateKey(dateValue)
}

const filteredHistory = computed(() => {
  if (!selectedDate.value) return history.value
  const selectedKey = getLocalDateKey(selectedDate.value)
  return history.value.filter(record => normalizeRecordDate(record.workout_date) === selectedKey)
})

const isRecordedToday = computed(() => {
  if (!todayWorkout.value) return false
  const workoutDate = normalizeRecordDate(todayWorkout.value.date || todayWorkout.value.workout_date)
  return history.value.some(record => {
    const recordDate = normalizeRecordDate(record.workout_date)
    if (recordDate !== workoutDate) return false
    if (!activePlanId.value) return true
    return record.plan_id === activePlanId.value
  })
})

const syncTodayWorkoutCompletion = () => {
  if (!todayWorkout.value) return
  const exercises = todayWorkout.value.exercises || []
  todayWorkout.value.is_completed = isRecordedToday.value
  todayWorkout.value.completed_exercises = isRecordedToday.value ? exercises.length : 0
  todayWorkout.value.exercises = exercises.map(ex => ({
    ...ex,
    completed: isRecordedToday.value ? true : ex.completed
  }))
}

// Methods
const loadData = async () => {
  loading.value = true
  try {
    await Promise.all([
      trainingStore.fetchPlans(),
      trainingStore.fetchTodayWorkout(),
      trainingStore.fetchHistory(),
      aiConfigStore.fetchConfigs()
    ])
    syncTodayWorkoutCompletion()
  } catch (error) {
    console.error('Failed to load training data:', error)
  } finally {
    loading.value = false
  }
}

const loadPlans = async () => {
  loadingPlans.value = true
  try {
    await trainingStore.fetchPlans()
  } catch (error) {
    console.error('Failed to load plans:', error)
  } finally {
    loadingPlans.value = false
  }
}

const loadHistory = async () => {
  loadingHistory.value = true
  try {
    await trainingStore.fetchHistory()
    syncTodayWorkoutCompletion()
  } catch (error) {
    console.error('Failed to load history:', error)
  } finally {
    loadingHistory.value = false
  }
}

const onRefresh = async () => {
  await loadData()
  refreshing.value = false
}

const onRefreshPlans = async () => {
  await loadPlans()
  refreshingPlans.value = false
}

const onRefreshHistory = async () => {
  await loadHistory()
  refreshingHistory.value = false
}

const handleAssessmentSubmit = async (assessmentData) => {
  submittingAssessment.value = true
  try {
    // Submit assessment
    await trainingService.submitAssessment(assessmentData)
    showToast({ type: 'success', message: t('assessment.success') })
    
    // Generate plan
    const defaultConfig = aiConfigStore.defaultConfig
    const planData = {
      plan_name: `${t('training.plans')} - ${new Date().toLocaleDateString()}`,
      duration_weeks: 12,
      goal: 'general_fitness',
      difficulty_level: assessmentData.experience_level === 'beginner' ? 'easy' : 
                        assessmentData.experience_level === 'advanced' ? 'hard' : 'medium',
      ai_api_id: defaultConfig?.id
    }
    
    const response = await trainingStore.generatePlan(planData)
    
    if (response?.data?.task_id) {
      generationTaskId.value = response.data.task_id
      startPollingTaskStatus()
    }
    
    showAssessment.value = false
  } catch (error) {
    showToast({ type: 'fail', message: t('assessment.failed') })
  } finally {
    submittingAssessment.value = false
  }
}

const startPollingTaskStatus = () => {
  generationProgress.value = 0
  generationError.value = ''
  isGeneratingTask.value = true
  generationInterval.value = setInterval(async () => {
    try {
      const response = await trainingService.checkTaskStatus(generationTaskId.value)
      const task = response?.data?.task || response?.data
      
      if (task) {
        generationProgress.value = task.progress || 0
        
        if (task.status === 'completed') {
          clearInterval(generationInterval.value)
          isGeneratingTask.value = false
          showToast({ type: 'success', message: t('training.planReady') })
          await loadPlans()
          await trainingStore.fetchTodayWorkout()
        } else if (task.status === 'failed') {
          clearInterval(generationInterval.value)
          isGeneratingTask.value = false
          generationError.value = task.error_message || t('error.unknown')
          showToast({ type: 'fail', message: task.error_message || t('error.unknown') })
        }
      }
    } catch (error) {
      clearInterval(generationInterval.value)
      isGeneratingTask.value = false
      console.error('Failed to check task status:', error)
    }
  }, 2000)
}

const handleRecordSubmit = async (recordData) => {
  submittingRecord.value = true
  try {
    if (isRecordedToday.value) {
      showToast({ type: 'fail', message: t('training.noWorkoutToday') })
      return
    }
    await trainingStore.recordWorkout(recordData)
    showToast({ type: 'success', message: t('app.success') })
    showRecordForm.value = false
    await trainingStore.fetchHistory()
    syncTodayWorkoutCompletion()
  } catch (error) {
    showToast({ type: 'fail', message: t('error.unknown') })
  } finally {
    submittingRecord.value = false
  }
}

const startWorkout = (workout) => {
  if (!workout || !workout.exercises || workout.exercises.length === 0) {
    showToast({ type: 'fail', message: t('training.noWorkoutToday') })
    return
  }
  if (isRecordedToday.value) {
    showToast({ type: 'success', message: t('training.completed') })
    return
  }
  if (workout.type === 'rest') {
    showToast({ type: 'success', message: t('training.restDay') })
    return
  }
  showRecordForm.value = true
}

const viewWorkoutDetails = async () => {
  if (currentPlan.value) {
    await viewPlanDetails(currentPlan.value)
    return
  }
  await loadPlans()
  if (plans.value.length > 0) {
    await viewPlanDetails(plans.value[0])
    return
  }
  showToast({ type: 'fail', message: t('training.noPlanHint') })
}

const viewPlanDetails = async (plan) => {
  try {
    if (!plan?.plan_data) {
      const response = await trainingStore.fetchPlan(plan.id)
      selectedPlan.value = response?.data?.plan || response?.data || plan
    } else {
      selectedPlan.value = plan
    }
    showPlanDetails.value = true
  } catch (error) {
    showToast({ type: 'fail', message: t('error.unknown') })
  }
}

const viewRecordDetails = (record) => {
  selectedRecord.value = record
  showRecordDetails.value = true
}

const startPlanWorkout = async (plan) => {
  if (plan) {
    trainingStore.currentPlan = plan
  }
  activeTab.value = 'today'
  loading.value = true
  try {
    await trainingStore.fetchTodayWorkout()
    syncTodayWorkoutCompletion()
    const workout = trainingStore.getTodayWorkout
    if (!workout || !workout.exercises || workout.exercises.length === 0) {
      showToast({ type: 'fail', message: t('training.noWorkoutToday') })
      return
    }
    if (isRecordedToday.value) {
      showToast({ type: 'success', message: t('training.completed') })
      return
    }
    if (workout.type === 'rest') {
      showToast({ type: 'success', message: t('training.restDay') })
      return
    }
    showRecordForm.value = true
  } catch (error) {
    showToast({ type: 'fail', message: t('error.unknown') })
  } finally {
    loading.value = false
  }
}

const selectExercise = (exercise, index) => {
  console.log('Selected exercise:', exercise, index)
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}

const calendarFormatter = (day) => {
  // Mark days with workouts
  const dateKey = getLocalDateKey(day.date)
  const hasWorkout = history.value.some(h => h.workout_date === dateKey)
  if (hasWorkout) {
    day.bottomInfo = '●'
  }
  return day
}

const getRecordExercises = (record) => {
  const exercises = record?.exercises
  if (Array.isArray(exercises)) {
    return exercises
  }
  if (exercises && Array.isArray(exercises.items)) {
    return exercises.items
  }
  return []
}

const getPlanCompletedCount = (plan) => {
  if (!plan?.id) return 0
  return history.value.filter(record => record.plan_id === plan.id).length
}

const getPlanTotalCount = (plan) => {
  return plan?.total_weeks || plan?.duration_weeks || 0
}

const formatPlanGoal = (plan) => {
  const goal = plan?.training_purpose || plan?.goal
  if (!goal) {
    return t('training.goalUnknown')
  }
  return t(`training.goals.${goal}`)
}

const onCalendarConfirm = (date) => {
  selectedDate.value = date
  showCalendar.value = false
}

// Watch tab changes to load data
watch(activeTab, (newTab) => {
  if (newTab === 'plans' && plans.value.length === 0) {
    loadPlans()
  } else if (newTab === 'history' && history.value.length === 0) {
    loadHistory()
  }
})

// Lifecycle
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.training-view {
  min-height: 100vh;
  background-color: var(--van-background);
  padding-bottom: 80px;
}

.tab-content {
  min-height: calc(100vh - 150px);
}

.tab-panel {
  padding: 16px 0;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.generating-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.generating-container h3 {
  margin: 16px 0 8px;
  color: var(--van-text-color);
}

.generating-container p {
  margin: 0 0 24px;
  color: var(--van-text-color-2);
  font-size: 14px;
}

.generating-container .van-progress {
  width: 80%;
}

.generation-error {
  padding: 16px;
}

.generation-error .error-detail {
  margin: 8px auto 16px;
  max-width: 520px;
  padding: 10px 12px;
  border-radius: 10px;
  background: #fff7e6;
  border: 1px solid #ffd591;
  color: #8c4a0d;
  font-size: 12px;
  line-height: 1.4;
}

.rest-day {
  text-align: center;
  padding: 40px 20px;
}

.rest-day-title {
  color: var(--van-primary-color);
  margin-top: 16px;
}

.exercise-list {
  margin: 16px;
}

.exercise-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.exercise-details {
  display: flex;
  gap: 16px;
  margin-top: 4px;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.exercise-safety {
  margin-top: 6px;
  padding: 6px 8px;
  border-radius: 6px;
  background: #f6f8fa;
  color: var(--van-text-color-2);
  font-size: 12px;
  line-height: 1.4;
}

.completed-exercise {
  opacity: 0.6;
}

.line-through {
  text-decoration: line-through;
}

.action-button {
  padding: 16px;
}

.history-list {
  padding: 8px 0;
}

.history-item {
  margin: 8px 16px;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.history-date {
  font-weight: 500;
}

.history-meta {
  display: flex;
  gap: 16px;
  margin-top: 4px;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.history-exercises {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.exercise-tag {
  padding: 2px 8px;
  background: var(--van-background);
  border-radius: 4px;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.more-tag {
  padding: 2px 8px;
  color: var(--van-primary-color);
  font-size: 12px;
}

.plan-details-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.plan-details-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.record-details-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.record-details-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.record-exercise-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 4px;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.record-exercise-notes {
  margin-top: 6px;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.section-title {
  padding: 16px;
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: var(--van-text-color);
}

.day-schedule {
  padding: 12px 0;
  border-bottom: 1px solid var(--van-border-color);
}

.day-schedule:last-child {
  border-bottom: none;
}

.day-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.day-number {
  font-weight: 500;
}

.day-exercises {
  padding-left: 8px;
}

.mini-exercise {
  font-size: 12px;
  color: var(--van-text-color-2);
  padding: 2px 0;
}
</style>
