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
              :target-calories="currentPlan.daily_calories"
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
                v-for="meal in todayMeals.meals"
                :key="meal.meal_type"
                :title="t(`nutrition.mealTypes.${meal.meal_type}`)"
                :label="formatMealTime(meal.time)"
                is-link
                @click="viewMealDetails(meal)"
              >
                <template #value>
                  <span class="meal-calories">{{ meal.total_calories }} {{ t('nutrition.kcal') }}</span>
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
          <div v-else-if="isGenerating" class="generating-container">
            <van-loading type="spinner" size="48" color="var(--van-primary-color)" />
            <h3>{{ t('nutrition.generating') }}</h3>
            <p>{{ t('nutrition.generatingHint') }}</p>
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
                    <span class="plan-name">{{ plan.name }}</span>
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
                    <span class="macro-value">{{ plan.protein_ratio || 30 }}%</span>
                  </div>
                  <div class="macro-item">
                    <span class="macro-label">{{ t('nutrition.carbs') }}</span>
                    <span class="macro-value">{{ plan.carbs_ratio || 40 }}%</span>
                  </div>
                  <div class="macro-item">
                    <span class="macro-label">{{ t('nutrition.fat') }}</span>
                    <span class="macro-value">{{ plan.fat_ratio || 30 }}%</span>
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
                    <span class="meal-calories">{{ meal.total_calories }} {{ t('nutrition.kcal') }}</span>
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
                v-model="planForm.name"
                :label="t('ai.name')"
                :placeholder="t('ai.name')"
                :rules="[{ required: true }]"
                autocapitalize="words"
                inputmode="text"
                enterkeyhint="next"
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
                v-model.number="planForm.carbs_ratio"
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
        <van-nav-bar :title="selectedPlan.name" />
        <div class="plan-details-content">
          <van-cell-group inset>
            <van-cell :title="t('nutrition.dailyCalories')" :value="`${selectedPlan.daily_calories} ${t('nutrition.kcal')}`" />
            <van-cell :title="t('nutrition.protein')" :value="`${selectedPlan.protein_ratio || 30}%`" />
            <van-cell :title="t('nutrition.carbs')" :value="`${selectedPlan.carbs_ratio || 40}%`" />
            <van-cell :title="t('nutrition.fat')" :value="`${selectedPlan.fat_ratio || 30}%`" />
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

// Selected items
const selectedPlan = ref(null)
const activeMeals = ref([0])

// Plan form
const planForm = ref({
  name: '',
  daily_calories: 2000,
  protein_ratio: 30,
  carbs_ratio: 40,
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

const todayTotalCalories = computed(() => nutritionStore.todayTotalCalories)

const todayTotalProtein = computed(() => {
  if (!todayMeals.value?.meals) return 0
  return todayMeals.value.meals.reduce((total, meal) => total + (meal.total_protein || 0), 0)
})

const todayTotalCarbs = computed(() => {
  if (!todayMeals.value?.meals) return 0
  return todayMeals.value.meals.reduce((total, meal) => total + (meal.total_carbs || 0), 0)
})

const todayTotalFat = computed(() => {
  if (!todayMeals.value?.meals) return 0
  return todayMeals.value.meals.reduce((total, meal) => total + (meal.total_fat || 0), 0)
})

const dailyProteinTarget = computed(() => {
  if (!currentPlan.value) return 0
  const ratio = currentPlan.value.protein_ratio || 30
  return Math.round((currentPlan.value.daily_calories * ratio / 100) / 4)
})

const dailyCarbsTarget = computed(() => {
  if (!currentPlan.value) return 0
  const ratio = currentPlan.value.carbs_ratio || 40
  return Math.round((currentPlan.value.daily_calories * ratio / 100) / 4)
})

const dailyFatTarget = computed(() => {
  if (!currentPlan.value) return 0
  const ratio = currentPlan.value.fat_ratio || 30
  return Math.round((currentPlan.value.daily_calories * ratio / 100) / 9)
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
    const planData = {
      ...planForm.value,
      ai_api_id: defaultConfig?.id
    }
    
    await nutritionStore.generatePlan(planData)
    showToast({ type: 'success', message: t('nutrition.planReady') })
    showPlanForm.value = false
    
    // Reset form
    planForm.value = {
      name: '',
      daily_calories: 2000,
      protein_ratio: 30,
      carbs_ratio: 40,
      fat_ratio: 30,
      dietary_restrictions: []
    }
  } catch (error) {
    showToast({ type: 'fail', message: t('error.unknown') })
  }
}

const handleRecordMeal = async (mealData) => {
  submittingRecord.value = true
  try {
    await nutritionStore.recordMeal(mealData)
    showToast({ type: 'success', message: t('nutrition.recordSuccess') })
    showRecordForm.value = false
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

const viewPlanDetails = (plan) => {
  selectedPlan.value = plan
  showPlanDetails.value = true
}

const viewMealDetails = (meal) => {
  console.log('View meal details:', meal)
}

const viewMealRecord = (meal) => {
  console.log('View meal record:', meal)
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString()
}

const formatMealTime = (timeStr) => {
  if (!timeStr) return ''
  if (timeStr.includes('T')) {
    return new Date(timeStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
  return timeStr
}

const calculateDayTotal = (meals) => {
  return meals.reduce((total, meal) => total + (meal.total_calories || 0), 0)
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
