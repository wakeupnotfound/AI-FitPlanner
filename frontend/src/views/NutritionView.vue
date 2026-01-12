<template>
  <div class="nutrition-view">
    <!-- Header -->
    <van-nav-bar :title="t('nutrition.title')" />

    <!-- Tabs -->
    <van-tabs v-model:active="activeTab" sticky>
      <van-tab :title="t('nutrition.todayMeals')" name="today" />
      <van-tab :title="t('nutrition.plans')" name="plans" />
      <van-tab :title="t('nutrition.history')" name="history" />
    </van-tabs>

    <!-- Tab Content -->
    <div class="tab-content">
      <!-- Today's Meals Tab -->
      <div v-show="activeTab === 'today'" class="tab-panel">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
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

          <!-- No Plan State -->
          <van-empty
            v-else-if="!currentPlan"
            image="search"
            :description="t('nutrition.noPlanHint')"
          >
            <van-button type="primary" @click="showPlanForm = true">
              {{ t('nutrition.generatePlan') }}
            </van-button>
          </van-empty>

          <!-- Today's Meals Content -->
          <div v-else-if="todayMeals">
            <!-- Nutrition Summary Card -->
            <NutritionSummaryCard
              :consumed-calories="todayTotalCalories"
              :target-calories="currentPlan?.daily_calories || 0"
              :consumed-protein="todayTotalProtein"
              :target-protein="dailyProteinTarget"
              :consumed-carbs="todayTotalCarbs"
              :target-carbs="dailyCarbsTarget"
              :consumed-fat="todayTotalFat"
              :target-fat="dailyFatTarget"
            />

            <!-- Meal Schedule -->
            <van-cell-group inset :title="t('nutrition.mealSchedule')" class="meal-list">
              <van-cell
                v-for="meal in todayMealsList"
                :key="meal.meal_type || meal.time"
                :title="meal.meal_type ? t(`nutrition.mealTypes.${meal.meal_type}`) : meal.time"
                :label="formatMealTime(meal.time || meal.meal_date)"
                is-link
                @click="viewMealDetails(meal)"
              >
                <template #value>
                  <span class="meal-calories">{{ getMealCalories(meal) }} {{ t('nutrition.kcal') }}</span>
                </template>
              </van-cell>
            </van-cell-group>

            <!-- Record Meal Button -->
            <div class="action-button">
              <van-button type="primary" block @click="showRecordForm = true">
                {{ t('nutrition.recordMeal') }}
              </van-button>
            </div>
          </div>

          <!-- No Meals Today -->
          <van-empty
            v-else
            image="search"
            :description="t('nutrition.noMealsToday')"
          >
            <van-button type="primary" @click="showRecordForm = true">
              {{ t('nutrition.recordMeal') }}
            </van-button>
          </van-empty>
        </van-pull-refresh>
      </div>

      <!-- Plans Tab -->
      <div v-show="activeTab === 'plans'" class="tab-panel">
        <van-pull-refresh v-model="refreshingPlans" @refresh="onRefreshPlans">
          <div v-if="loadingPlans" class="loading-container">
            <LoadingSpinner />
          </div>

          <!-- Plan Generation Progress -->
          <div v-else-if="isGeneratingDisplay" class="generating-container">
            <van-loading type="spinner" size="48" color="var(--van-primary-color)" />
            <h3>{{ t('nutrition.generating') }}</h3>
            <p>{{ t('nutrition.generatingHint') }}</p>
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
              :description="t('nutrition.planGenerateFailed')"
            >
              <div class="error-detail">{{ generationError }}</div>
              <van-button type="primary" plain @click="showPlanForm = true">
                {{ t('nutrition.generatePlan') }}
              </van-button>
            </van-empty>
          </div>

          <!-- No Plans -->
          <van-empty
            v-else-if="plans.length === 0"
            image="search"
            :description="t('nutrition.noPlanHint')"
          >
            <van-button type="primary" @click="showPlanForm = true">
              {{ t('nutrition.generatePlan') }}
            </van-button>
          </van-empty>

          <!-- Plans List -->
          <div v-else>
            <van-cell-group
              v-for="plan in plans"
              :key="plan.id"
              inset
              class="plan-card"
            >
              <van-cell :border="false">
                <template #title>
                  <div class="plan-header">
                  <span class="plan-name">{{ plan.plan_name || plan.name }}</span>
                  <van-tag v-if="plan.status === 'active'" type="success" size="small">
                    {{ t('training.active') }}
                  </van-tag>
                  </div>
                </template>
                <template #label>
                  <div class="plan-meta">
                    <span>{{ plan.daily_calories }} {{ t('nutrition.kcal') }}/{{ t('nutrition.dailyTarget') }}</span>
                  </div>
                </template>
              </van-cell>

              <!-- Dietary Restrictions -->
              <van-cell v-if="plan.plan_data?.restrictions?.length > 0" :border="false">
                <div class="restrictions-list">
                  <van-tag
                    v-for="restriction in plan.plan_data.restrictions"
                    :key="restriction"
                    type="warning"
                    size="small"
                    class="restriction-tag"
                  >
                    {{ t(`nutrition.restrictions.${restriction}`) }}
                  </van-tag>
                </div>
              </van-cell>

              <!-- Macro Ratios -->
              <van-cell :border="false">
                <div class="macro-ratios">
                  <div class="macro-item">
                    <span class="macro-label">{{ t('nutrition.protein') }}</span>
                    <span class="macro-value">{{ Math.round((plan.protein_ratio || 0) * 100) }}%</span>
                  </div>
                  <div class="macro-item">
                    <span class="macro-label">{{ t('nutrition.carbs') }}</span>
                    <span class="macro-value">{{ Math.round((plan.carb_ratio || 0) * 100) }}%</span>
                  </div>
                  <div class="macro-item">
                    <span class="macro-label">{{ t('nutrition.fat') }}</span>
                    <span class="macro-value">{{ Math.round((plan.fat_ratio || 0) * 100) }}%</span>
                  </div>
                </div>
              </van-cell>

              <van-cell :border="false">
                <van-button size="small" type="primary" plain @click="viewPlanDetails(plan)">
                  {{ t('nutrition.viewPlan') }}
                </van-button>
              </van-cell>
            </van-cell-group>

            <!-- Generate New Plan Button -->
            <div class="action-button">
              <van-button type="primary" plain block @click="showPlanForm = true">
                {{ t('nutrition.generatePlan') }}
              </van-button>
            </div>
          </div>
        </van-pull-refresh>
      </div>

      <!-- History Tab -->
      <div v-show="activeTab === 'history'" class="tab-panel">
        <van-pull-refresh v-model="refreshingHistory" @refresh="onRefreshHistory">
          <div v-if="loadingHistory" class="loading-container">
            <LoadingSpinner />
          </div>

          <!-- No History -->
          <van-empty
            v-else-if="Object.keys(historyGroupedByDate).length === 0"
            image="search"
            :description="t('nutrition.noHistory')"
          />

          <!-- History List Grouped by Date -->
          <div v-else class="history-list">
            <div
              v-for="(meals, date) in historyGroupedByDate"
              :key="date"
              class="history-day"
            >
              <div class="history-date-header">
                <span class="history-date">{{ formatDate(date) }}</span>
                <span class="history-day-total">
                  {{ calculateDayTotal(meals) }} {{ t('nutrition.kcal') }}
                </span>
              </div>

              <van-cell-group inset>
                <van-cell
                  v-for="meal in meals"
                  :key="meal.id"
                  :title="t(`nutrition.mealTypes.${meal.meal_type}`)"
                  :label="formatMealTime(meal.meal_date)"
                  is-link
                  @click="viewMealRecord(meal)"
                >
                  <template #value>
                    <span class="meal-calories">{{ getMealCalories(meal) }} {{ t('nutrition.kcal') }}</span>
                  </template>
                </van-cell>
              </van-cell-group>
            </div>
          </div>
        </van-pull-refresh>
      </div>
    </div>

    <!-- Plan Generation Form Popup -->
    <van-popup
      v-model:show="showPlanForm"
      position="bottom"
      round
      :style="{ height: '80%' }"
      closeable
    >
      <div class="plan-form-popup">
        <van-nav-bar :title="t('nutrition.generatePlan')" />
        <div class="plan-form-content">
          <van-form @submit="handleGeneratePlan">
            <van-cell-group inset>
              <van-field
                v-model="planForm.plan_name"
                :label="t('nutrition.planName')"
                :placeholder="t('nutrition.planName')"
                :rules="[{ required: true, message: t('nutrition.validation.planNameRequired') }]"
                autocapitalize="words"
                inputmode="text"
                enterkeyhint="next"
              />
              <van-field
                v-model.number="planForm.duration_days"
                type="digit"
                inputmode="numeric"
                enterkeyhint="next"
                :label="t('nutrition.durationDays')"
                :placeholder="t('nutrition.durationDays')"
                :rules="[{ required: true, message: t('nutrition.validation.durationRequired') }]"
              />
              <van-field
                v-model.number="planForm.daily_calories"
                type="number"
                inputmode="numeric"
                enterkeyhint="next"
                :label="t('nutrition.dailyCalories')"
                :placeholder="t('nutrition.dailyCalories')"
                :rules="[{ required: true }]"
              />
            </van-cell-group>

            <van-cell-group inset :title="t('nutrition.macroRatios')">
              <van-field
                v-model.number="planForm.protein_ratio"
                type="number"
                inputmode="numeric"
                enterkeyhint="next"
                :label="t('nutrition.protein')"
                :placeholder="'30'"
              >
                <template #button>%</template>
              </van-field>
              <van-field
                v-model.number="planForm.carb_ratio"
                type="number"
                inputmode="numeric"
                enterkeyhint="next"
                :label="t('nutrition.carbs')"
                :placeholder="'40'"
              >
                <template #button>%</template>
              </van-field>
              <van-field
                v-model.number="planForm.fat_ratio"
                type="number"
                inputmode="numeric"
                enterkeyhint="done"
                :label="t('nutrition.fat')"
                :placeholder="'30'"
              >
                <template #button>%</template>
              </van-field>
            </van-cell-group>

            <van-cell-group inset :title="t('nutrition.dietaryRestrictions')">
              <van-checkbox-group v-model="planForm.dietary_restrictions">
                <van-cell
                  v-for="restriction in restrictionOptions"
                  :key="restriction"
                  :title="t(`nutrition.restrictions.${restriction}`)"
                  clickable
                  @click="toggleRestriction(restriction)"
                >
                  <template #right-icon>
                    <van-checkbox :name="restriction" @click.stop />
                  </template>
                </van-cell>
              </van-checkbox-group>
            </van-cell-group>

            <div class="form-actions">
              <van-button type="primary" block native-type="submit" :loading="isGenerating">
                {{ t('nutrition.generatePlan') }}
              </van-button>
            </div>
          </van-form>
        </div>
      </div>
    </van-popup>

    <!-- Meal Record Form Popup -->
    <van-popup
      v-model:show="showRecordForm"
      position="bottom"
      round
      :style="{ height: '90%' }"
      closeable
    >
      <MealRecordForm
        :plan-meals="currentPlan?.plan_data?.meals"
        :loading="submittingRecord"
        @submit="handleRecordMeal"
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
            <van-cell :title="t('nutrition.dailyCalories')" :value="`${selectedPlan.daily_calories} ${t('nutrition.kcal')}`" />
            <van-cell :title="t('nutrition.protein')" :value="`${Math.round((selectedPlan.protein_ratio || 0) * 100)}%`" />
            <van-cell :title="t('nutrition.carbs')" :value="`${Math.round((selectedPlan.carb_ratio || 0) * 100)}%`" />
            <van-cell :title="t('nutrition.fat')" :value="`${Math.round((selectedPlan.fat_ratio || 0) * 100)}%`" />
          </van-cell-group>

          <!-- Dietary Restrictions -->
          <van-cell-group inset :title="t('nutrition.dietaryRestrictions')" v-if="selectedPlan.plan_data?.restrictions?.length > 0">
            <van-cell :border="false">
              <div class="restrictions-list">
                <van-tag
                  v-for="restriction in selectedPlan.plan_data.restrictions"
                  :key="restriction"
                  type="warning"
                  size="medium"
                  class="restriction-tag"
                >
                  {{ t(`nutrition.restrictions.${restriction}`) }}
                </van-tag>
              </div>
            </van-cell>
          </van-cell-group>
          <van-cell-group inset :title="t('nutrition.dietaryRestrictions')" v-else>
            <van-cell :title="t('nutrition.noDietaryRestrictions')" />
          </van-cell-group>

          <!-- Meal Schedule -->
          <h4 class="section-title">{{ t('nutrition.mealSchedule') }}</h4>
          <van-collapse v-model="activeMeals" v-if="selectedPlan.plan_data?.meals">
            <van-collapse-item
              v-for="(dayMeals, index) in selectedPlan.plan_data.meals"
              :key="index"
              :title="dayMeals.date || `${t('training.day')} ${index + 1}`"
              :name="index"
            >
              <div v-for="meal in dayMeals.meals" :key="meal.meal_type" class="meal-detail">
                <div class="meal-detail-header">
                  <span class="meal-type">{{ t(`nutrition.mealTypes.${meal.meal_type}`) }}</span>
                  <span class="meal-time">{{ meal.time }}</span>
                </div>
                <div class="meal-foods">
                  <div v-for="(food, i) in meal.foods" :key="i" class="food-item">
                    {{ food.name }} - {{ food.amount }}{{ food.unit }} ({{ food.calories }} {{ t('nutrition.kcal') }})
                  </div>
                </div>
              </div>
            </van-collapse-item>
          </van-collapse>
        </div>
      </div>
    </van-popup>

    <!-- Meal Record Details Popup -->
    <van-popup
      v-model:show="showRecordDetails"
      position="bottom"
      round
      :style="{ height: '80%' }"
      closeable
    >
      <div class="meal-record-details" v-if="selectedMealRecord">
        <van-nav-bar :title="t('nutrition.recordMeal')" />
        <div class="meal-record-content">
          <van-cell-group inset>
            <van-cell :title="t('nutrition.mealType')" :value="t(`nutrition.mealTypes.${selectedMealRecord.meal_type}`)" />
            <van-cell :title="t('nutrition.mealTime')" :value="formatMealTime(selectedMealRecord.meal_date)" />
            <van-cell :title="t('nutrition.calories')" :value="`${selectedMealRecord.calories || 0} ${t('nutrition.kcal')}`" />
          </van-cell-group>

          <van-cell-group inset :title="t('nutrition.foods')">
            <van-cell
              v-for="(food, idx) in getRecordFoods(selectedMealRecord)"
              :key="idx"
              :title="food.name"
              :label="`${food.amount}${food.unit} - ${food.calories} ${t('nutrition.kcal')}`"
            >
              <template #value>
                <span class="food-macros">{{ food.protein || 0 }}P / {{ food.carbs || 0 }}C / {{ food.fat || 0 }}F</span>
              </template>
            </van-cell>
          </van-cell-group>

          <van-cell-group inset v-if="selectedMealRecord.notes" :title="t('nutrition.notes')">
            <van-cell :border="false" :label="selectedMealRecord.notes" />
          </van-cell-group>
        </div>
      </div>
    </van-popup>

    <!-- Meal Schedule Details Popup -->
    <van-popup
      v-model:show="showMealDetailsPopup"
      position="bottom"
      round
      :style="{ height: '80%' }"
      closeable
    >
      <div class="meal-record-details" v-if="selectedMealDetail">
        <van-nav-bar :title="t('nutrition.mealSchedule')" />
        <div class="meal-record-content">
          <van-cell-group inset>
            <van-cell
              :title="t('nutrition.mealType')"
              :value="selectedMealDetail.meal_type ? t(`nutrition.mealTypes.${selectedMealDetail.meal_type}`) : selectedMealDetail.time"
            />
            <van-cell
              v-if="selectedMealDetail.time || selectedMealDetail.meal_date"
              :title="t('nutrition.mealTime')"
              :value="formatMealTime(selectedMealDetail.time || selectedMealDetail.meal_date)"
            />
            <van-cell
              :title="t('nutrition.calories')"
              :value="`${getMealCalories(selectedMealDetail)} ${t('nutrition.kcal')}`"
            />
          </van-cell-group>

          <van-cell-group inset :title="t('nutrition.foods')">
            <van-cell
              v-for="(food, idx) in getMealFoods(selectedMealDetail)"
              :key="idx"
              :title="food.name"
              :label="`${food.amount}${food.unit} - ${food.calories} ${t('nutrition.kcal')}`"
            >
              <template #value>
                <span class="food-macros">{{ food.protein || 0 }}P / {{ food.carbs || 0 }}C / {{ food.fat || 0 }}F</span>
              </template>
            </van-cell>
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
import { useNutritionStore } from '@/stores/nutrition'
import { useAIConfigStore } from '@/stores/aiConfig'
import { nutritionService } from '@/services/nutrition.service'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import NavigationBar from '@/components/common/NavigationBar.vue'

// Lazy load heavy components
const NutritionSummaryCard = defineAsyncComponent(() => import('@/components/nutrition/NutritionSummaryCard.vue'))
const MealRecordForm = defineAsyncComponent(() => import('@/components/nutrition/MealRecordForm.vue'))

const { t } = useI18n()
const nutritionStore = useNutritionStore()
const aiConfigStore = useAIConfigStore()

// State
const activeTab = ref('today')
const loading = ref(false)
const loadingPlans = ref(false)
const loadingHistory = ref(false)
const refreshing = ref(false)
const refreshingPlans = ref(false)
const refreshingHistory = ref(false)
const submittingRecord = ref(false)

// Popups
const showPlanForm = ref(false)
const showRecordForm = ref(false)
const showPlanDetails = ref(false)
const showRecordDetails = ref(false)
const showMealDetailsPopup = ref(false)
const generationTaskId = ref(null)
const generationProgress = ref(0)
const generationInterval = ref(null)
const generationError = ref('')
const isGeneratingTask = ref(false)

// Selected items
const selectedPlan = ref(null)
const selectedMealRecord = ref(null)
const selectedMealDetail = ref(null)
const activeMeals = ref([0])

// Plan form
const planForm = ref({
  plan_name: '',
  duration_days: 7,
  daily_calories: 2000,
  protein_ratio: 30,
  carb_ratio: 40,
  fat_ratio: 30,
  dietary_restrictions: []
})

const restrictionOptions = [
  'vegetarian', 'vegan', 'gluten_free', 'dairy_free',
  'nut_free', 'low_carb', 'keto', 'halal', 'kosher'
]

// Computed
const hasAIConfig = computed(() => aiConfigStore.hasDefaultConfig)
const currentPlan = computed(() => nutritionStore.activePlan)
const todayMeals = computed(() => nutritionStore.getTodayMeals)
const plans = computed(() => nutritionStore.allPlans)
const history = computed(() => nutritionStore.mealHistory)
const historyGroupedByDate = computed(() => nutritionStore.mealHistoryGroupedByDate)
const isGenerating = computed(() => nutritionStore.isGeneratingPlan)
const isGeneratingDisplay = computed(() => isGenerating.value || isGeneratingTask.value)

const todayTotalCalories = computed(() => nutritionStore.todayTotalCalories)
const todayKey = computed(() => getLocalDateKey(new Date()))
const todayMealRecords = computed(() => {
  return history.value.filter(record => normalizeRecordDate(record.meal_date) === todayKey.value)
})
const todayMealsList = computed(() => {
  if (todayMealRecords.value.length > 0) {
    return todayMealRecords.value.map(record => ({
      ...record,
      total_calories: record.calories,
      time: record.meal_date
    }))
  }
  const meals = Object.values(todayMeals.value?.meals || {})
  return meals
})

const todayTotalProtein = computed(() => {
  if (todayMealRecords.value.length > 0) {
    return todayMealRecords.value.reduce((total, meal) => total + (meal.protein || 0), 0)
  }
  if (!todayMeals.value?.meals) return 0
  const meals = Object.values(todayMeals.value.meals || {})
  return meals.reduce((total, meal) => total + (meal?.total_protein || 0), 0)
})

const todayTotalCarbs = computed(() => {
  if (todayMealRecords.value.length > 0) {
    return todayMealRecords.value.reduce((total, meal) => total + (meal.carbs || 0), 0)
  }
  if (!todayMeals.value?.meals) return 0
  const meals = Object.values(todayMeals.value.meals || {})
  return meals.reduce((total, meal) => total + (meal?.total_carbs || 0), 0)
})

const todayTotalFat = computed(() => {
  if (todayMealRecords.value.length > 0) {
    return todayMealRecords.value.reduce((total, meal) => total + (meal.fat || 0), 0)
  }
  if (!todayMeals.value?.meals) return 0
  const meals = Object.values(todayMeals.value.meals || {})
  return meals.reduce((total, meal) => total + (meal?.total_fat || 0), 0)
})

const dailyProteinTarget = computed(() => {
  if (!currentPlan.value) return 0
  const ratio = currentPlan.value.protein_ratio || 0.3
  return Math.round((currentPlan.value.daily_calories * ratio) / 4)
})

const dailyCarbsTarget = computed(() => {
  if (!currentPlan.value) return 0
  const ratio = currentPlan.value.carb_ratio || 0.4
  return Math.round((currentPlan.value.daily_calories * ratio) / 4)
})

const dailyFatTarget = computed(() => {
  if (!currentPlan.value) return 0
  const ratio = currentPlan.value.fat_ratio || 0.3
  return Math.round((currentPlan.value.daily_calories * ratio) / 9)
})

// Methods
const loadData = async () => {
  loading.value = true
  try {
    await Promise.all([
      nutritionStore.fetchPlans(),
      nutritionStore.fetchTodayMeals(),
      aiConfigStore.fetchConfigs()
    ])
  } catch (error) {
    console.error('Failed to load nutrition data:', error)
  } finally {
    loading.value = false
  }
}

const loadPlans = async () => {
  loadingPlans.value = true
  try {
    await nutritionStore.fetchPlans()
  } catch (error) {
    console.error('Failed to load plans:', error)
  } finally {
    loadingPlans.value = false
  }
}

const loadHistory = async () => {
  loadingHistory.value = true
  try {
    await nutritionStore.fetchHistory()
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

const handleGeneratePlan = async () => {
  try {
    const defaultConfig = aiConfigStore.defaultConfig
    const proteinRatio = planForm.value.protein_ratio > 1 ? planForm.value.protein_ratio / 100 : planForm.value.protein_ratio
    const carbRatio = planForm.value.carb_ratio > 1 ? planForm.value.carb_ratio / 100 : planForm.value.carb_ratio
    const fatRatio = planForm.value.fat_ratio > 1 ? planForm.value.fat_ratio / 100 : planForm.value.fat_ratio
    const planData = {
      ...planForm.value,
      protein_ratio: proteinRatio,
      carb_ratio: carbRatio,
      fat_ratio: fatRatio,
      ai_api_id: defaultConfig?.id
    }
    
    const response = await nutritionStore.generatePlan(planData)
    const taskId = response?.data?.task_id
    if (taskId) {
      generationTaskId.value = taskId
      startPollingTaskStatus()
      showPlanForm.value = false
    } else {
      showToast({ type: 'success', message: t('nutrition.planReady') })
      showPlanForm.value = false
    }
    
    // Reset form
    planForm.value = {
      plan_name: '',
      duration_days: 7,
      daily_calories: 2000,
      protein_ratio: 30,
      carb_ratio: 40,
      fat_ratio: 30,
      dietary_restrictions: []
    }
  } catch (error) {
    showToast({ type: 'fail', message: t('error.unknown') })
  }
}

const startPollingTaskStatus = () => {
  generationProgress.value = 0
  generationError.value = ''
  isGeneratingTask.value = true
  generationInterval.value = setInterval(async () => {
    try {
      const response = await nutritionService.checkTaskStatus(generationTaskId.value)
      const task = response?.data?.task || response?.data
      if (task) {
        generationProgress.value = task.progress || 0
        if (task.status === 'completed') {
          clearInterval(generationInterval.value)
          isGeneratingTask.value = false
          showToast({ type: 'success', message: t('nutrition.planReady') })
          await loadPlans()
          await nutritionStore.fetchTodayMeals()
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

const handleRecordMeal = async (mealData) => {
  submittingRecord.value = true
  try {
    await nutritionStore.recordMeal(mealData)
    showToast({ type: 'success', message: t('nutrition.recordSuccess') })
    showRecordForm.value = false
    await nutritionStore.fetchHistory()
    await nutritionStore.fetchTodayMeals()
  } catch (error) {
    showToast({ type: 'fail', message: t('nutrition.recordFailed') })
  } finally {
    submittingRecord.value = false
  }
}

const toggleRestriction = (restriction) => {
  const index = planForm.value.dietary_restrictions.indexOf(restriction)
  if (index === -1) {
    planForm.value.dietary_restrictions.push(restriction)
  } else {
    planForm.value.dietary_restrictions.splice(index, 1)
  }
}

const viewPlanDetails = async (plan) => {
  try {
    if (!plan?.plan_data) {
      const response = await nutritionStore.fetchPlan(plan.id)
      selectedPlan.value = response?.data?.plan || response?.data || plan
    } else {
      selectedPlan.value = plan
    }
    showPlanDetails.value = true
  } catch (error) {
    showToast({ type: 'fail', message: t('error.unknown') })
  }
}

const viewMealDetails = (meal) => {
  selectedMealDetail.value = meal
  showMealDetailsPopup.value = true
}

const viewMealRecord = (meal) => {
  selectedMealRecord.value = meal
  showRecordDetails.value = true
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString()
}

const normalizeRecordDate = (dateValue) => {
  if (!dateValue) return ''
  if (typeof dateValue === 'string') {
    return dateValue.slice(0, 10)
  }
  return getLocalDateKey(dateValue)
}

const getLocalDateKey = (date) => {
  const localDate = new Date(date)
  const offsetMs = localDate.getTimezoneOffset() * 60 * 1000
  return new Date(localDate.getTime() - offsetMs).toISOString().slice(0, 10)
}

const formatMealTime = (timeStr) => {
  if (!timeStr) return ''
  if (timeStr.includes('T')) {
    return new Date(timeStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
  return timeStr
}

const getMealCalories = (meal) => {
  return meal?.total_calories ?? meal?.calories ?? 0
}

const calculateDayTotal = (meals) => {
  return meals.reduce((total, meal) => total + getMealCalories(meal), 0)
}

const getRecordFoods = (record) => {
  const foods = record?.foods
  if (Array.isArray(foods)) {
    return foods
  }
  if (foods && Array.isArray(foods.items)) {
    return foods.items
  }
  return []
}

const getMealFoods = (meal) => {
  const foods = meal?.foods
  if (Array.isArray(foods)) {
    return foods
  }
  if (foods && Array.isArray(foods.items)) {
    return foods.items
  }
  return []
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
.nutrition-view {
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
  margin: 0;
  color: var(--van-text-color-2);
  font-size: 14px;
}

.generating-container .van-progress {
  width: 80%;
  margin-top: 16px;
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

.meal-list {
  margin: 16px;
}

.meal-calories {
  color: var(--van-primary-color);
  font-weight: 500;
}

.action-button {
  padding: 16px;
}

.plan-card {
  margin: 8px 16px;
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.plan-name {
  font-weight: 500;
}

.plan-meta {
  margin-top: 4px;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.restrictions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.restriction-tag {
  margin: 2px;
}

.macro-ratios {
  display: flex;
  justify-content: space-around;
}

.macro-item {
  text-align: center;
}

.macro-label {
  display: block;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.macro-value {
  font-weight: 500;
  color: var(--van-primary-color);
}

.history-list {
  padding: 8px 0;
}

.history-day {
  margin-bottom: 16px;
}

.history-date-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
}

.history-date {
  font-weight: 500;
  color: var(--van-text-color);
}

.history-day-total {
  font-size: 14px;
  color: var(--van-primary-color);
}

.plan-form-popup,
.plan-details-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.plan-form-content,
.plan-details-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.meal-record-details {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.meal-record-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.food-macros {
  font-size: 12px;
  color: var(--van-text-color-2);
}

.form-actions {
  padding: 16px;
}

.section-title {
  padding: 16px;
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: var(--van-text-color);
}

.meal-detail {
  padding: 12px 0;
  border-bottom: 1px solid var(--van-border-color);
}

.meal-detail:last-child {
  border-bottom: none;
}

.meal-detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.meal-type {
  font-weight: 500;
}

.meal-time {
  font-size: 12px;
  color: var(--van-text-color-2);
}

.meal-foods {
  padding-left: 8px;
}

.food-item {
  font-size: 12px;
  color: var(--van-text-color-2);
  padding: 2px 0;
}
</style>
