<template>
  <van-cell-group inset class="exercise-item" :class="{ 'is-completed': exercise.completed, 'is-current': isCurrent }">
    <van-cell :border="false" clickable @click="$emit('click', exercise)">
      <template #title>
        <div class="exercise-header">
          <van-checkbox
            v-model="localCompleted"
            shape="square"
            @click.stop
            @change="handleCompletedChange"
          />
          <div class="exercise-info">
            <h4 class="exercise-name" :class="{ 'line-through': exercise.completed }">
              {{ exercise.name }}
            </h4>
            <div class="exercise-meta">
              <span class="meta-badge sets">
                {{ exercise.sets }} {{ t('training.sets') }}
              </span>
              <span class="meta-badge reps">
                {{ exercise.reps }} {{ t('training.reps') }}
              </span>
              <span v-if="exercise.weight" class="meta-badge weight">
                {{ exercise.weight }}
              </span>
            </div>
          </div>
        </div>
      </template>
      <template #right-icon>
        <van-icon name="arrow" v-if="!exercise.completed" />
        <van-icon name="success" color="var(--van-success-color)" v-else />
      </template>
    </van-cell>

    <!-- Expanded details -->
    <van-cell v-if="showDetails" :border="false" class="exercise-details">
      <div class="details-grid">
        <div class="detail-item" v-if="exercise.rest">
          <span class="detail-label">{{ t('training.rest') }}</span>
          <span class="detail-value">{{ exercise.rest }}</span>
        </div>
        <div class="detail-item" v-if="exercise.difficulty">
          <span class="detail-label">{{ t('training.difficulty') }}</span>
          <span class="detail-value">{{ t(`training.difficultyLevels.${exercise.difficulty}`) }}</span>
        </div>
      </div>

      <!-- Safety notes -->
      <van-notice-bar
        v-if="exercise.safety_notes"
        left-icon="warning-o"
        :text="exercise.safety_notes"
        color="#ed6a0c"
        background="#fffbe8"
        class="safety-notice"
      />

      <!-- Notes -->
      <div v-if="exercise.notes" class="exercise-notes">
        <span class="notes-label">{{ t('training.notes') }}:</span>
        <span class="notes-text">{{ exercise.notes }}</span>
      </div>
    </van-cell>

    <!-- Action buttons -->
    <van-cell v-if="showActions && !exercise.completed" :border="false">
      <div class="exercise-actions">
        <van-button
          size="small"
          type="primary"
          plain
          @click="$emit('record', exercise)"
        >
          {{ t('training.recordWorkout') }}
        </van-button>
      </div>
    </van-cell>
  </van-cell-group>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  exercise: {
    type: Object,
    required: true
  },
  isCurrent: {
    type: Boolean,
    default: false
  },
  showDetails: {
    type: Boolean,
    default: false
  },
  showActions: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click', 'record', 'complete'])

const localCompleted = ref(props.exercise.completed || false)

watch(() => props.exercise.completed, (newVal) => {
  localCompleted.value = newVal
})

const handleCompletedChange = (value) => {
  emit('complete', { ...props.exercise, completed: value })
}
</script>

<style scoped>
.exercise-item {
  margin: 8px 16px;
  transition: all 0.3s ease;
}

.exercise-item.is-completed {
  opacity: 0.7;
}

.exercise-item.is-current {
  border-left: 3px solid var(--van-primary-color);
}

.exercise-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.exercise-info {
  flex: 1;
}

.exercise-name {
  margin: 0 0 8px 0;
  font-size: 15px;
  font-weight: 500;
  color: var(--van-text-color);
}

.exercise-name.line-through {
  text-decoration: line-through;
  color: var(--van-text-color-2);
}

.exercise-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.meta-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.meta-badge.sets {
  background: var(--van-primary-color-light);
  color: var(--van-primary-color);
}

.meta-badge.reps {
  background: var(--van-success-color-light);
  color: var(--van-success-color);
}

.meta-badge.weight {
  background: var(--van-warning-color-light);
  color: var(--van-warning-color);
}

.exercise-details {
  padding-top: 0;
}

.details-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 12px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  color: var(--van-text-color-3);
}

.detail-value {
  font-size: 14px;
  font-weight: 500;
  color: var(--van-text-color);
}

.safety-notice {
  border-radius: 6px;
  margin-bottom: 12px;
}

.exercise-notes {
  padding: 8px 12px;
  background: var(--van-background);
  border-radius: 6px;
  font-size: 13px;
}

.notes-label {
  color: var(--van-text-color-2);
  margin-right: 4px;
}

.notes-text {
  color: var(--van-text-color);
}

.exercise-actions {
  display: flex;
  justify-content: flex-end;
}

:deep(.van-checkbox) {
  margin-top: 2px;
}
</style>
