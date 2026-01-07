<template>
  <van-cell-group inset class="nutrition-summary-card">
    <van-cell :border="false">
      <template #title>
        <div class="summary-header">
          <span class="summary-title">{{ t('nutrition.nutritionSummary') }}</span>
        </div>
      </template>
    </van-cell>

    <!-- Calories Progress -->
    <van-cell :border="false">
      <div class="progress-section">
        <div class="progress-header">
          <span class="progress-label">{{ t('nutrition.calories') }}</span>
          <span class="progress-value" :class="{ exceeded: caloriesExceeded }">
            {{ consumedCalories }} / {{ targetCalories }} {{ t('nutrition.kcal') }}
          </span>
        </div>
        <van-progress
          :percentage="caloriesPercentage"
          :stroke-width="12"
          :color="caloriesExceeded ? '#ee0a24' : '#07c160'"
          track-color="#ebedf0"
        />
        <div class="progress-footer">
          <span v-if="!caloriesExceeded" class="remaining">
            {{ t('nutrition.remaining') }}: {{ remainingCalories }} {{ t('nutrition.kcal') }}
          </span>
          <span v-else class="exceeded-text">
            {{ t('nutrition.exceeded') }}: {{ Math.abs(remainingCalories) }} {{ t('nutrition.kcal') }}
          </span>
        </div>
      </div>
    </van-cell>

    <!-- Macros Progress -->
    <van-cell :border="false">
      <div class="macros-section">
        <div class="macro-progress">
          <div class="macro-header">
            <span class="macro-label">{{ t('nutrition.protein') }}</span>
            <span class="macro-value">{{ consumedProtein }}g / {{ targetProtein }}g</span>
          </div>
          <van-progress
            :percentage="proteinPercentage"
            :stroke-width="8"
            color="#1989fa"
            track-color="#ebedf0"
          />
        </div>

        <div class="macro-progress">
          <div class="macro-header">
            <span class="macro-label">{{ t('nutrition.carbs') }}</span>
            <span class="macro-value">{{ consumedCarbs }}g / {{ targetCarbs }}g</span>
          </div>
          <van-progress
            :percentage="carbsPercentage"
            :stroke-width="8"
            color="#ff976a"
            track-color="#ebedf0"
          />
        </div>

        <div class="macro-progress">
          <div class="macro-header">
            <span class="macro-label">{{ t('nutrition.fat') }}</span>
            <span class="macro-value">{{ consumedFat }}g / {{ targetFat }}g</span>
          </div>
          <van-progress
            :percentage="fatPercentage"
            :stroke-width="8"
            color="#7232dd"
            track-color="#ebedf0"
          />
        </div>
      </div>
    </van-cell>
  </van-cell-group>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  consumedCalories: {
    type: Number,
    default: 0
  },
  targetCalories: {
    type: Number,
    default: 2000
  },
  consumedProtein: {
    type: Number,
    default: 0
  },
  targetProtein: {
    type: Number,
    default: 150
  },
  consumedCarbs: {
    type: Number,
    default: 0
  },
  targetCarbs: {
    type: Number,
    default: 200
  },
  consumedFat: {
    type: Number,
    default: 0
  },
  targetFat: {
    type: Number,
    default: 65
  }
})

// Computed
const caloriesPercentage = computed(() => {
  if (props.targetCalories === 0) return 0
  return Math.min(100, Math.round((props.consumedCalories / props.targetCalories) * 100))
})

const caloriesExceeded = computed(() => props.consumedCalories > props.targetCalories)

const remainingCalories = computed(() => props.targetCalories - props.consumedCalories)

const proteinPercentage = computed(() => {
  if (props.targetProtein === 0) return 0
  return Math.min(100, Math.round((props.consumedProtein / props.targetProtein) * 100))
})

const carbsPercentage = computed(() => {
  if (props.targetCarbs === 0) return 0
  return Math.min(100, Math.round((props.consumedCarbs / props.targetCarbs) * 100))
})

const fatPercentage = computed(() => {
  if (props.targetFat === 0) return 0
  return Math.min(100, Math.round((props.consumedFat / props.targetFat) * 100))
})
</script>

<style scoped>
.nutrition-summary-card {
  margin: 16px;
}

.summary-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--van-text-color);
}

.progress-section {
  padding: 8px 0;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.progress-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--van-text-color);
}

.progress-value {
  font-size: 14px;
  color: var(--van-text-color-2);
}

.progress-value.exceeded {
  color: #ee0a24;
}

.progress-footer {
  margin-top: 8px;
  text-align: right;
}

.remaining {
  font-size: 12px;
  color: #07c160;
}

.exceeded-text {
  font-size: 12px;
  color: #ee0a24;
}

.macros-section {
  padding: 8px 0;
}

.macro-progress {
  margin-bottom: 16px;
}

.macro-progress:last-child {
  margin-bottom: 0;
}

.macro-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.macro-label {
  font-size: 13px;
  color: var(--van-text-color);
}

.macro-value {
  font-size: 12px;
  color: var(--van-text-color-2);
}
</style>
