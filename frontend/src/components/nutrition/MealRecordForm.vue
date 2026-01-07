<template>
  <div class="meal-record-form">
    <van-nav-bar :title="t('nutrition.recordMeal')" />
    
    <div class="form-content">
      <van-form @submit="handleSubmit">
        <!-- Meal Type Selection -->
        <van-cell-group inset>
          <van-field
            v-model="form.meal_type"
            is-link
            readonly
            :label="t('nutrition.mealType')"
            :placeholder="t('nutrition.validation.mealTypeRequired')"
            :rules="[{ required: true, message: t('nutrition.validation.mealTypeRequired') }]"
            @click="showMealTypePicker = true"
          />
          <van-popup v-model:show="showMealTypePicker" position="bottom" round>
            <van-picker
              :columns="mealTypeColumns"
              @confirm="onMealTypeConfirm"
              @cancel="showMealTypePicker = false"
            />
          </van-popup>

          <van-field
            v-model="form.meal_date"
            is-link
            readonly
            :label="t('nutrition.mealTime')"
            :placeholder="t('nutrition.mealTime')"
            @click="showDatePicker = true"
          />
          <van-popup v-model:show="showDatePicker" position="bottom" round>
            <van-date-picker
              v-model="selectedDate"
              type="datetime"
              :min-date="minDate"
              :max-date="maxDate"
              @confirm="onDateConfirm"
              @cancel="showDatePicker = false"
            />
          </van-popup>
        </van-cell-group>

        <!-- Food Selection Tabs -->
        <van-tabs v-model:active="foodTab" class="food-tabs">
          <van-tab :title="t('nutrition.selectFromPlan')" name="plan" />
          <van-tab :title="t('nutrition.customFood')" name="custom" />
        </van-tabs>

        <!-- Select from Plan -->
        <div v-show="foodTab === 'plan'" class="food-section">
          <van-cell-group inset v-if="planMeals && planMeals.length > 0">
            <van-checkbox-group v-model="selectedPlanFoods">
              <van-cell
                v-for="(food, index) in availablePlanFoods"
                :key="index"
                :title="food.name"
                :label="`${food.amount}${food.unit} - ${food.calories} ${t('nutrition.kcal')}`"
                clickable
                @click="togglePlanFood(index)"
              >
                <template #right-icon>
                  <van-checkbox :name="index" @click.stop />
                </template>
              </van-cell>
            </van-checkbox-group>
          </van-cell-group>
          <van-empty v-else :description="t('nutrition.noPlan')" image="search" />
        </div>

        <!-- Custom Food Entry -->
        <div v-show="foodTab === 'custom'" class="food-section">
          <van-cell-group inset :title="t('nutrition.foods')">
            <div v-for="(food, index) in form.foods" :key="index" class="food-entry">
              <van-field
                v-model="food.name"
                :label="t('nutrition.foodName')"
                :placeholder="t('nutrition.validation.foodNameRequired')"
                :rules="[{ required: true, message: t('nutrition.validation.foodNameRequired') }]"
                autocapitalize="words"
                inputmode="text"
                enterkeyhint="next"
              />
              <div class="food-row">
                <van-field
                  v-model.number="food.amount"
                  type="number"
                  inputmode="decimal"
                  enterkeyhint="next"
                  :label="t('nutrition.amount')"
                  :placeholder="t('nutrition.validation.amountRequired')"
                  class="amount-field"
                />
                <van-field
                  v-model="food.unit"
                  is-link
                  readonly
                  :label="t('nutrition.unit')"
                  class="unit-field"
                  @click="showUnitPicker(index)"
                />
              </div>

              <div class="food-row">
                <van-field
                  v-model.number="food.calories"
                  type="number"
                  inputmode="numeric"
                  enterkeyhint="next"
                  :label="t('nutrition.calories')"
                  :placeholder="t('nutrition.validation.caloriesRequired')"
                  class="nutrition-field"
                />
                <van-field
                  v-model.number="food.protein"
                  type="number"
                  inputmode="decimal"
                  enterkeyhint="next"
                  :label="t('nutrition.protein')"
                  placeholder="0"
                  class="nutrition-field"
                >
                  <template #button>g</template>
                </van-field>
              </div>
              <div class="food-row">
                <van-field
                  v-model.number="food.carbs"
                  type="number"
                  inputmode="decimal"
                  enterkeyhint="next"
                  :label="t('nutrition.carbs')"
                  placeholder="0"
                  class="nutrition-field"
                >
                  <template #button>g</template>
                </van-field>
                <van-field
                  v-model.number="food.fat"
                  type="number"
                  inputmode="decimal"
                  enterkeyhint="done"
                  :label="t('nutrition.fat')"
                  placeholder="0"
                  class="nutrition-field"
                >
                  <template #button>g</template>
                </van-field>
              </div>
              <van-button
                v-if="form.foods.length > 1"
                type="danger"
                plain
                size="small"
                class="remove-food-btn"
                @click="removeFood(index)"
              >
                {{ t('app.delete') }}
              </van-button>
              <van-divider v-if="index < form.foods.length - 1" />
            </div>
          </van-cell-group>

          <div class="add-food-btn">
            <van-button type="primary" plain size="small" icon="plus" @click="addFood">
              {{ t('nutrition.addFood') }}
            </van-button>
          </div>
        </div>

        <!-- Nutrition Summary -->
        <van-cell-group inset :title="t('nutrition.nutritionSummary')">
          <van-cell :title="t('nutrition.totalCalories')" :value="`${totalCalories} ${t('nutrition.kcal')}`" />
          <van-cell :title="t('nutrition.protein')" :value="`${totalProtein} g`" />
          <van-cell :title="t('nutrition.carbs')" :value="`${totalCarbs} g`" />
          <van-cell :title="t('nutrition.fat')" :value="`${totalFat} g`" />
        </van-cell-group>

        <!-- Notes -->
        <van-cell-group inset>
          <van-field
            v-model="form.notes"
            type="textarea"
            :label="t('nutrition.notes')"
            :placeholder="t('nutrition.notes')"
            rows="2"
            autosize
          />
        </van-cell-group>

        <!-- Actions -->
        <div class="form-actions">
          <van-button type="default" block @click="$emit('cancel')">
            {{ t('app.cancel') }}
          </van-button>
          <van-button type="primary" block native-type="submit" :loading="loading">
            {{ t('app.save') }}
          </van-button>
        </div>
      </van-form>
    </div>

    <!-- Unit Picker Popup -->
    <van-popup v-model:show="showUnitPickerPopup" position="bottom" round>
      <van-picker
        :columns="unitColumns"
        @confirm="onUnitConfirm"
        @cancel="showUnitPickerPopup = false"
      />
    </van-popup>
  </div>
</template>


<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { sanitizeInput } from '@/utils/sanitizer'

const { t } = useI18n()

const props = defineProps({
  planMeals: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'cancel'])

// State
const showMealTypePicker = ref(false)
const showDatePicker = ref(false)
const showUnitPickerPopup = ref(false)
const currentFoodIndex = ref(0)
const foodTab = ref('custom')
const selectedPlanFoods = ref([])

const minDate = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
const maxDate = new Date()
const selectedDate = ref([
  new Date().getFullYear().toString(),
  (new Date().getMonth() + 1).toString().padStart(2, '0'),
  new Date().getDate().toString().padStart(2, '0'),
  new Date().getHours().toString().padStart(2, '0'),
  new Date().getMinutes().toString().padStart(2, '0')
])

const form = ref({
  meal_type: '',
  meal_date: new Date().toISOString(),
  foods: [createEmptyFood()],
  notes: ''
})

const mealTypeColumns = [
  { text: t('nutrition.mealTypes.breakfast'), value: 'breakfast' },
  { text: t('nutrition.mealTypes.lunch'), value: 'lunch' },
  { text: t('nutrition.mealTypes.dinner'), value: 'dinner' },
  { text: t('nutrition.mealTypes.snack'), value: 'snack' }
]

const unitColumns = [
  { text: t('nutrition.units.g'), value: 'g' },
  { text: t('nutrition.units.ml'), value: 'ml' },
  { text: t('nutrition.units.cup'), value: 'cup' },
  { text: t('nutrition.units.tbsp'), value: 'tbsp' },
  { text: t('nutrition.units.tsp'), value: 'tsp' },
  { text: t('nutrition.units.piece'), value: 'piece' },
  { text: t('nutrition.units.serving'), value: 'serving' }
]

// Computed
const availablePlanFoods = computed(() => {
  if (!props.planMeals || props.planMeals.length === 0) return []
  const allFoods = []
  props.planMeals.forEach(dayMeals => {
    if (dayMeals.meals) {
      dayMeals.meals.forEach(meal => {
        if (meal.foods) {
          allFoods.push(...meal.foods)
        }
      })
    }
  })
  return allFoods
})

const totalCalories = computed(() => {
  let total = 0
  if (foodTab.value === 'plan') {
    selectedPlanFoods.value.forEach(index => {
      const food = availablePlanFoods.value[index]
      if (food) total += food.calories || 0
    })
  } else {
    form.value.foods.forEach(food => {
      total += food.calories || 0
    })
  }
  return total
})

const totalProtein = computed(() => {
  let total = 0
  if (foodTab.value === 'plan') {
    selectedPlanFoods.value.forEach(index => {
      const food = availablePlanFoods.value[index]
      if (food) total += food.protein || 0
    })
  } else {
    form.value.foods.forEach(food => {
      total += food.protein || 0
    })
  }
  return total
})

const totalCarbs = computed(() => {
  let total = 0
  if (foodTab.value === 'plan') {
    selectedPlanFoods.value.forEach(index => {
      const food = availablePlanFoods.value[index]
      if (food) total += food.carbs || 0
    })
  } else {
    form.value.foods.forEach(food => {
      total += food.carbs || 0
    })
  }
  return total
})

const totalFat = computed(() => {
  let total = 0
  if (foodTab.value === 'plan') {
    selectedPlanFoods.value.forEach(index => {
      const food = availablePlanFoods.value[index]
      if (food) total += food.fat || 0
    })
  } else {
    form.value.foods.forEach(food => {
      total += food.fat || 0
    })
  }
  return total
})

// Methods
function createEmptyFood() {
  return {
    name: '',
    amount: null,
    unit: 'g',
    calories: null,
    protein: 0,
    carbs: 0,
    fat: 0
  }
}

const onMealTypeConfirm = ({ selectedOptions }) => {
  form.value.meal_type = selectedOptions[0].value
  showMealTypePicker.value = false
}

const onDateConfirm = ({ selectedValues }) => {
  const [year, month, day, hour, minute] = selectedValues
  const date = new Date(year, month - 1, day, hour, minute)
  form.value.meal_date = date.toISOString()
  showDatePicker.value = false
}

const showUnitPicker = (index) => {
  currentFoodIndex.value = index
  showUnitPickerPopup.value = true
}

const onUnitConfirm = ({ selectedOptions }) => {
  form.value.foods[currentFoodIndex.value].unit = selectedOptions[0].value
  showUnitPickerPopup.value = false
}

const addFood = () => {
  form.value.foods.push(createEmptyFood())
}

const removeFood = (index) => {
  form.value.foods.splice(index, 1)
}

const togglePlanFood = (index) => {
  const idx = selectedPlanFoods.value.indexOf(index)
  if (idx === -1) {
    selectedPlanFoods.value.push(index)
  } else {
    selectedPlanFoods.value.splice(idx, 1)
  }
}

const handleSubmit = () => {
  let foods = []
  
  if (foodTab.value === 'plan') {
    foods = selectedPlanFoods.value.map(index => availablePlanFoods.value[index])
  } else {
    foods = form.value.foods.filter(f => f.name && f.calories)
  }

  if (foods.length === 0) {
    return
  }

  const mealData = {
    meal_type: form.value.meal_type,
    meal_date: form.value.meal_date,
    foods,
    total_calories: totalCalories.value,
    total_protein: totalProtein.value,
    total_carbs: totalCarbs.value,
    total_fat: totalFat.value,
    notes: sanitizeInput(form.value.notes)
  }

  emit('submit', mealData)
}
</script>


<style scoped>
.meal-record-form {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.form-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 16px;
}

.food-tabs {
  margin: 16px 0;
}

.food-section {
  min-height: 200px;
}

.food-entry {
  padding: 8px 0;
}

.food-row {
  display: flex;
  gap: 8px;
}

.amount-field {
  flex: 1;
}

.unit-field {
  width: 100px;
}

.nutrition-field {
  flex: 1;
}

.remove-food-btn {
  margin: 8px 16px;
}

.add-food-btn {
  padding: 16px;
  text-align: center;
}

.form-actions {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: var(--van-background-2);
  border-top: 1px solid var(--van-border-color);
}

.form-actions .van-button {
  flex: 1;
}
</style>
