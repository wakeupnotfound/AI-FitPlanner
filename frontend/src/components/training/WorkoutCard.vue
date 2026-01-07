<template>
  <van-cell-group inset class="workout-card">
    <van-cell :border="false">
      <template #title>
        <div class="workout-header">
          <div class="workout-info">
            <h3 class="workout-title">{{ workout.focus_area ? t(`training.focusAreas.${workout.focus_area}`) : workout.type }}</h3>
            <van-tag :type="typeTagType" size="small">
              {{ t(`training.workoutTypes.${workout.type}`) }}
            </van-tag>
          </div>
          <van-tag v-if="workout.is_completed" type="success" size="medium">
            {{ t('training.completed') }}
          </van-tag>
        </div>
      </template>
      <template #label>
        <div class="workout-meta">
          <span class="meta-item">
            <van-icon name="clock-o" />
            {{ workout.duration }} {{ t('training.minutes') }}
          </span>
          <span class="meta-item" v-if="workout.estimated_calories">
            <van-icon name="fire-o" />
            {{ workout.estimated_calories }} {{ t('training.kcal') }}
          </span>
          <span class="meta-item" v-if="workout.exercises">
            <van-icon name="orders-o" />
            {{ workout.exercises.length }} {{ t('training.exercises') }}
          </span>
        </div>
      </template>
    </van-cell>

    <!-- Exercise preview -->
    <van-cell v-if="showExercises && workout.exercises && workout.exercises.length > 0" :border="false">
      <div class="exercise-preview">
        <div 
          v-for="(exercise, index) in previewExercises" 
          :key="index"
          class="exercise-item"
        >
          <span class="exercise-name">{{ exercise.name }}</span>
          <span class="exercise-detail">{{ exercise.sets }} × {{ exercise.reps }}</span>
        </div>
        <div v-if="workout.exercises.length > 3" class="more-exercises">
          +{{ workout.exercises.length - 3 }} {{ t('training.exercises') }}
        </div>
      </div>
    </van-cell>

    <!-- Safety alert -->
    <van-cell v-if="hasSafetyNotes" :border="false">
      <van-notice-bar
        left-icon="warning-o"
        :text="safetyNotesText"
        color="#ed6a0c"
        background="#fffbe8"
      />
    </van-cell>

    <!-- Actions -->
    <van-cell v-if="showActions" :border="false">
      <div class="workout-actions">
        <van-button 
          v-if="!workout.is_completed"
          size="small" 
          type="primary"
          @click="$emit('start', workout)"
        >
          {{ t('training.startWorkout') }}
        </van-button>
        <van-button 
          v-if="workout.is_completed"
          size="small" 
          plain
          @click="$emit('view', workout)"
        >
          {{ t('training.viewPlan') }}
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
  workout: {
    type: Object,
    required: true
  },
  showExercises: {
    type: Boolean,
    default: true
  },
  showActions: {
    type: Boolean,
    default: true
  }
})

defineEmits(['start', 'view'])

const typeTagType = computed(() => {
  switch (props.workout.type) {
    case 'strength':
      return 'primary'
    case 'cardio':
      return 'success'
    case 'hiit':
      return 'danger'
    case 'flexibility':
      return 'warning'
    case 'rest':
      return 'default'
    default:
      return 'default'
  }
})

const previewExercises = computed(() => {
  if (!props.workout.exercises) return []
  return props.workout.exercises.slice(0, 3)
})

const hasSafetyNotes = computed(() => {
  if (!props.workout.exercises) return false
  return props.workout.exercises.some(e => e.safety_notes)
})

const safetyNotesText = computed(() => {
  if (!props.workout.exercises) return ''
  const notes = props.workout.exercises
    .filter(e => e.safety_notes)
    .map(e => e.safety_notes)
  return notes.join(' | ')
})
</script>

<style scoped>
.workout-card {
  margin: 12px 16px;
}

.workout-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.workout-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.workout-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--van-text-color);
}

.workout-meta {
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

.exercise-preview {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.exercise-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: var(--van-background);
  border-radius: 6px;
}

.exercise-name {
  font-size: 14px;
  color: var(--van-text-color);
}

.exercise-detail {
  font-size: 12px;
  color: var(--van-text-color-2);
}

.more-exercises {
  font-size: 12px;
  color: var(--van-primary-color);
  text-align: center;
  padding: 4px;
}

.workout-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

:deep(.van-notice-bar) {
  border-radius: 6px;
}
</style>
