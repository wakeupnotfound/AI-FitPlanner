<template>
  <div class="workout-record-form">
    <!-- Header -->
    <van-nav-bar
      :title="t('training.recordWorkout')"
      left-arrow
      @click-left="$emit('cancel')"
    >
      <template #right>
        <van-button
          type="primary"
          size="small"
          :loading="loading"
          :disabled="!canSubmit"
          @click="handleSubmit"
        >
          {{ t('app.save') }}
        </van-button>
      </template>
    </van-nav-bar>

    <!-- Form content -->
    <div class="form-content">
      <!-- Workout summary -->
      <van-cell-group inset :title="t('training.planOverview')">
        <van-cell :title="t('training.duration')">
          <template #value>
            <van-stepper
              v-model="formData.duration_minutes"
              :min="1"
              :max="300"
              :step="5"
            />
            <span class="unit">{{ t('training.minutes') }}</span>
          </template>
        </van-cell>
        <van-cell :title="t('training.rating')">
          <template #value>
            <van-rate v-model="formData.rating" allow-half />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Exercise records -->
      <van-cell-group inset :title="t('training.exercises')">
        <van-collapse v-model="activeExercises">
          <van-collapse-item
            v-for="(exercise, index) in formData.exercises"
            :key="index"
            :name="index"
          >
            <template #title>
              <div class="exercise-collapse-title">
                <van-checkbox
                  v-model="exercise.completed"
                  shape="square"
                  @click.stop
                />
                <span :class="{ 'completed-text': exercise.completed }">
                  {{ exercise.exercise_name }}
                </span>
              </div>
            </template>

            <div class="exercise-record">
              <!-- Sets -->
              <div class="sets-header">
                <span>{{ t('training.sets') }}</span>
                <van-button
                  size="mini"
                  plain
                  type="primary"
                  @click="addSet(index)"
                >
                  + {{ t('app.add') }}
                </van-button>
              </div>

              <div class="sets-list">
                <div
                  v-for="(set, setIndex) in exercise.sets_data"
                  :key="setIndex"
                  class="set-row"
                >
                  <span class="set-number">{{ setIndex + 1 }}</span>
                  <van-field
                    v-model="set.reps"
                    type="digit"
                    inputmode="numeric"
                    enterkeyhint="next"
                    :placeholder="t('training.reps')"
                    class="set-input"
                  />
                  <span class="input-separator">×</span>
                  <van-field
                    v-model="set.weight"
                    type="number"
                    inputmode="decimal"
                    enterkeyhint="next"
                    :placeholder="t('training.weight')"
                    class="set-input"
                  />
                  <span class="weight-unit">{{ t('training.kg') }}</span>
                  <van-icon
                    name="delete-o"
                    class="delete-set touch-target-icon"
                    @click="removeSet(index, setIndex)"
                  />
                </div>
              </div>

              <!-- Exercise notes -->
              <van-field
                v-model="exercise.notes"
                type="textarea"
                :placeholder="t('training.notes')"
                rows="1"
                autosize
                class="exercise-notes"
              />
            </div>
          </van-collapse-item>
        </van-collapse>
      </van-cell-group>

      <!-- Injury report -->
      <van-cell-group inset :title="t('training.injuryReport')">
        <van-cell>
          <van-checkbox v-model="hasInjury" shape="square">
            {{ t('training.reportInjury') }}
          </van-checkbox>
        </van-cell>
        <van-field
          v-if="hasInjury"
          v-model="formData.injury_report"
          type="textarea"
          :placeholder="t('assessment.injuryHistoryPlaceholder')"
          rows="2"
          autosize
        />
        <van-cell v-if="!hasInjury" class="no-injury">
          <van-icon name="passed" color="var(--van-success-color)" />
          <span>{{ t('training.noInjury') }}</span>
        </van-cell>
      </van-cell-group>

      <!-- General notes -->
      <van-cell-group inset :title="t('training.notes')">
        <van-field
          v-model="formData.notes"
          type="textarea"
          :placeholder="t('training.notes')"
          rows="2"
          autosize
        />
      </van-cell-group>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  workout: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  planId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['submit', 'cancel'])

// State
const hasInjury = ref(false)
const activeExercises = ref([0])

// Initialize form data from workout
const initializeExercises = () => {
  if (!props.workout?.exercises) return []
  
  return props.workout.exercises.map(ex => ({
    exercise_name: ex.name,
    completed: ex.completed || false,
    sets_data: initializeSets(ex.sets, ex.reps, ex.weight),
    notes: ''
  }))
}

const initializeSets = (setsCount, reps, weight) => {
  const sets = []
  const numSets = parseInt(setsCount) || 3
  const defaultReps = typeof reps === 'string' ? reps.split('-')[0] : reps
  const defaultWeight = weight ? weight.replace(/[^0-9.]/g, '') : ''
  
  for (let i = 0; i < numSets; i++) {
    sets.push({
      reps: defaultReps?.toString() || '',
      weight: defaultWeight || ''
    })
  }
  return sets
}

// Form data
const formData = reactive({
  plan_id: props.planId,
  workout_date: new Date().toISOString().split('T')[0],
  workout_type: props.workout?.type || 'strength',
  duration_minutes: props.workout?.duration || 60,
  exercises: initializeExercises(),
  performance_data: {
    total_volume: 0,
    estimated_calories: props.workout?.estimated_calories || 0
  },
  notes: '',
  rating: 4,
  injury_report: null
})

// Computed
const canSubmit = computed(() => {
  return formData.exercises.some(ex => ex.completed)
})

// Methods
const addSet = (exerciseIndex) => {
  const exercise = formData.exercises[exerciseIndex]
  const lastSet = exercise.sets_data[exercise.sets_data.length - 1]
  exercise.sets_data.push({
    reps: lastSet?.reps || '',
    weight: lastSet?.weight || ''
  })
}

const removeSet = (exerciseIndex, setIndex) => {
  const exercise = formData.exercises[exerciseIndex]
  if (exercise.sets_data.length > 1) {
    exercise.sets_data.splice(setIndex, 1)
  }
}

const calculateTotalVolume = () => {
  let total = 0
  formData.exercises.forEach(ex => {
    if (ex.completed) {
      ex.sets_data.forEach(set => {
        const reps = parseInt(set.reps) || 0
        const weight = parseFloat(set.weight) || 0
        total += reps * weight
      })
    }
  })
  return total
}

const handleSubmit = () => {
  // Calculate performance data
  formData.performance_data.total_volume = calculateTotalVolume()
  
  // Format exercises for API
  const exercises = formData.exercises
    .filter(ex => ex.completed)
    .map(ex => ({
      exercise_name: ex.exercise_name,
      sets: ex.sets_data.length,
      reps_per_set: ex.sets_data.map(s => parseInt(s.reps) || 0),
      weight_used: ex.sets_data.map(s => parseFloat(s.weight) || 0),
      notes: ex.notes || null
    }))

  const submitData = {
    plan_id: formData.plan_id,
    workout_date: formData.workout_date,
    workout_type: formData.workout_type,
    duration_minutes: formData.duration_minutes,
    exercises,
    performance_data: formData.performance_data,
    notes: formData.notes || null,
    rating: formData.rating,
    injury_report: hasInjury.value ? formData.injury_report : null
  }

  emit('submit', submitData)
}

// Watch for injury checkbox
watch(hasInjury, (value) => {
  if (!value) {
    formData.injury_report = null
  }
})

// Watch workout prop changes
watch(() => props.workout, () => {
  formData.exercises = initializeExercises()
}, { deep: true })
</script>

<style scoped>
.workout-record-form {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--van-background);
}

.form-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.unit {
  margin-left: 8px;
  font-size: 14px;
  color: var(--van-text-color-2);
}

.exercise-collapse-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.completed-text {
  text-decoration: line-through;
  color: var(--van-text-color-2);
}

.exercise-record {
  padding: 12px 0;
}

.sets-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 14px;
  font-weight: 500;
}

.sets-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.set-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.set-number {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--van-primary-color);
  color: white;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 500;
}

.set-input {
  flex: 1;
  padding: 0;
}

.set-input :deep(.van-field__control) {
  text-align: center;
}

.input-separator {
  color: var(--van-text-color-2);
}

.weight-unit {
  font-size: 12px;
  color: var(--van-text-color-2);
  min-width: 24px;
}

.delete-set {
  color: var(--van-danger-color);
  font-size: 18px;
  cursor: pointer;
}

.exercise-notes {
  margin-top: 12px;
}

.no-injury {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--van-success-color);
}

:deep(.van-cell-group__title) {
  padding-left: 16px;
}

:deep(.van-collapse-item__content) {
  padding: 0 16px 16px;
}

:deep(.van-stepper) {
  display: inline-flex;
}
</style>
