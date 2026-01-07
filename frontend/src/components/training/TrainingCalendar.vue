<template>
  <div class="training-calendar">
    <!-- Month navigation -->
    <div class="calendar-header">
      <van-icon name="arrow-left" @click="previousMonth" />
      <span class="month-title">{{ monthTitle }}</span>
      <van-icon name="arrow" @click="nextMonth" />
    </div>

    <!-- Weekday headers -->
    <div class="weekday-row">
      <span v-for="day in weekdays" :key="day" class="weekday">{{ day }}</span>
    </div>

    <!-- Calendar grid -->
    <div class="calendar-grid">
      <div
        v-for="(day, index) in calendarDays"
        :key="index"
        class="calendar-day"
        :class="{
          'other-month': !day.isCurrentMonth,
          'today': day.isToday,
          'selected': isSelected(day),
          'has-workout': day.hasWorkout
        }"
        @click="selectDay(day)"
      >
        <span class="day-number">{{ day.date.getDate() }}</span>
        <span v-if="day.hasWorkout" class="workout-indicator">●</span>
      </div>
    </div>

    <!-- Selected day details -->
    <div v-if="selectedDay && selectedDayWorkouts.length > 0" class="selected-day-details">
      <h4 class="details-title">{{ formatSelectedDate }}</h4>
      <van-cell-group inset>
        <van-cell
          v-for="workout in selectedDayWorkouts"
          :key="workout.id"
          clickable
          @click="$emit('view-workout', workout)"
        >
          <template #title>
            <div class="workout-summary">
              <span class="workout-type">{{ workout.workout_type }}</span>
              <van-tag type="success" size="small">{{ t('training.completed') }}</van-tag>
            </div>
          </template>
          <template #label>
            <div class="workout-meta">
              <span>{{ workout.duration_minutes }} {{ t('training.minutes') }}</span>
              <span v-if="workout.exercises">{{ workout.exercises.length }} {{ t('training.exercises') }}</span>
              <van-rate v-if="workout.rating" v-model="workout.rating" readonly size="12" />
            </div>
          </template>
        </van-cell>
      </van-cell-group>
    </div>

    <!-- No workout message -->
    <div v-else-if="selectedDay" class="no-workout-message">
      <van-empty
        image="search"
        :description="t('training.noWorkoutToday')"
        :image-size="60"
      />
    </div>

    <!-- Stats summary -->
    <van-cell-group inset class="month-stats" v-if="monthStats">
      <van-grid :column-num="3" :border="false">
        <van-grid-item>
          <div class="stat-value">{{ monthStats.totalWorkouts }}</div>
          <div class="stat-label">{{ t('training.completed') }}</div>
        </van-grid-item>
        <van-grid-item>
          <div class="stat-value">{{ monthStats.totalMinutes }}</div>
          <div class="stat-label">{{ t('training.minutes') }}</div>
        </van-grid-item>
        <van-grid-item>
          <div class="stat-value">{{ monthStats.avgRating.toFixed(1) }}</div>
          <div class="stat-label">{{ t('training.rating') }}</div>
        </van-grid-item>
      </van-grid>
    </van-cell-group>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

const props = defineProps({
  workouts: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['select-date', 'view-workout'])

// State
const currentDate = ref(new Date())
const selectedDay = ref(null)

// Computed
const weekdays = computed(() => {
  const days = locale.value === 'zh' 
    ? ['日', '一', '二', '三', '四', '五', '六']
    : ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
  return days
})

const monthTitle = computed(() => {
  const options = { year: 'numeric', month: 'long' }
  return currentDate.value.toLocaleDateString(locale.value === 'zh' ? 'zh-CN' : 'en-US', options)
})

const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  
  const days = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  // Add days from previous month
  const startDayOfWeek = firstDay.getDay()
  for (let i = startDayOfWeek - 1; i >= 0; i--) {
    const date = new Date(year, month, -i)
    days.push({
      date,
      isCurrentMonth: false,
      isToday: false,
      hasWorkout: hasWorkoutOnDate(date)
    })
  }
  
  // Add days of current month
  for (let i = 1; i <= lastDay.getDate(); i++) {
    const date = new Date(year, month, i)
    days.push({
      date,
      isCurrentMonth: true,
      isToday: date.getTime() === today.getTime(),
      hasWorkout: hasWorkoutOnDate(date)
    })
  }
  
  // Add days from next month to complete the grid
  const remainingDays = 42 - days.length // 6 rows × 7 days
  for (let i = 1; i <= remainingDays; i++) {
    const date = new Date(year, month + 1, i)
    days.push({
      date,
      isCurrentMonth: false,
      isToday: false,
      hasWorkout: hasWorkoutOnDate(date)
    })
  }
  
  return days
})

const selectedDayWorkouts = computed(() => {
  if (!selectedDay.value) return []
  const dateStr = formatDateStr(selectedDay.value.date)
  return props.workouts.filter(w => w.workout_date === dateStr)
})

const formatSelectedDate = computed(() => {
  if (!selectedDay.value) return ''
  return selectedDay.value.date.toLocaleDateString(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const monthStats = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  const monthWorkouts = props.workouts.filter(w => {
    const date = new Date(w.workout_date)
    return date.getFullYear() === year && date.getMonth() === month
  })
  
  if (monthWorkouts.length === 0) return null
  
  const totalMinutes = monthWorkouts.reduce((sum, w) => sum + (w.duration_minutes || 0), 0)
  const ratings = monthWorkouts.filter(w => w.rating).map(w => w.rating)
  const avgRating = ratings.length > 0 ? ratings.reduce((a, b) => a + b, 0) / ratings.length : 0
  
  return {
    totalWorkouts: monthWorkouts.length,
    totalMinutes,
    avgRating
  }
})

// Methods
const formatDateStr = (date) => {
  return date.toISOString().split('T')[0]
}

const hasWorkoutOnDate = (date) => {
  const dateStr = formatDateStr(date)
  return props.workouts.some(w => w.workout_date === dateStr)
}

const isSelected = (day) => {
  if (!selectedDay.value) return false
  return day.date.getTime() === selectedDay.value.date.getTime()
}

const selectDay = (day) => {
  selectedDay.value = day
  emit('select-date', day.date)
}

const previousMonth = () => {
  currentDate.value = new Date(
    currentDate.value.getFullYear(),
    currentDate.value.getMonth() - 1,
    1
  )
}

const nextMonth = () => {
  currentDate.value = new Date(
    currentDate.value.getFullYear(),
    currentDate.value.getMonth() + 1,
    1
  )
}

// Initialize with today selected
selectedDay.value = calendarDays.value.find(d => d.isToday) || null
</script>

<style scoped>
.training-calendar {
  background: var(--van-background-2);
  border-radius: 12px;
  padding: 16px;
  margin: 16px;
}

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0 16px;
}

.month-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--van-text-color);
}

.calendar-header .van-icon {
  font-size: 20px;
  color: var(--van-primary-color);
  cursor: pointer;
  padding: 8px;
}

.weekday-row {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  text-align: center;
  margin-bottom: 8px;
}

.weekday {
  font-size: 12px;
  color: var(--van-text-color-2);
  padding: 8px 0;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.calendar-day {
  aspect-ratio: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  transition: all 0.2s ease;
}

.calendar-day:hover {
  background: var(--van-gray-2);
}

.calendar-day.other-month {
  opacity: 0.3;
}

.calendar-day.today {
  background: var(--van-primary-color-light);
}

.calendar-day.today .day-number {
  color: var(--van-primary-color);
  font-weight: 600;
}

.calendar-day.selected {
  background: var(--van-primary-color);
}

.calendar-day.selected .day-number {
  color: white;
}

.calendar-day.has-workout .workout-indicator {
  color: var(--van-success-color);
}

.calendar-day.selected .workout-indicator {
  color: white;
}

.day-number {
  font-size: 14px;
  color: var(--van-text-color);
}

.workout-indicator {
  font-size: 8px;
  position: absolute;
  bottom: 4px;
}

.selected-day-details {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--van-border-color);
}

.details-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--van-text-color);
  margin: 0 0 12px 0;
}

.workout-summary {
  display: flex;
  align-items: center;
  gap: 8px;
}

.workout-type {
  font-weight: 500;
  text-transform: capitalize;
}

.workout-meta {
  display: flex;
  gap: 12px;
  margin-top: 4px;
  font-size: 12px;
  color: var(--van-text-color-2);
}

.no-workout-message {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--van-border-color);
}

.month-stats {
  margin: 16px 0 0 0;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: var(--van-primary-color);
}

.stat-label {
  font-size: 12px;
  color: var(--van-text-color-3);
  margin-top: 4px;
}

:deep(.van-grid-item__content) {
  padding: 12px 8px;
}

:deep(.van-cell-group--inset) {
  margin: 0;
}
</style>
