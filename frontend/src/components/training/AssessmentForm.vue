<template>
  <div class="assessment-form">
    <!-- Header with progress -->
    <van-nav-bar
      :title="t('assessment.title')"
      left-arrow
      @click-left="handleBack"
    >
      <template #right>
        <span class="step-indicator">
          {{ t('assessment.step') }} {{ currentStep }} {{ t('assessment.of') }} {{ totalSteps }}
        </span>
      </template>
    </van-nav-bar>

    <!-- Progress bar -->
    <van-progress
      :percentage="progressPercentage"
      :show-pivot="false"
      color="var(--van-primary-color)"
      track-color="var(--van-gray-3)"
      stroke-width="4"
    />

    <!-- Step content -->
    <div class="step-content">
      <!-- Step 1: Experience Level -->
      <div v-show="currentStep === 1" class="step-panel">
        <h3 class="step-title">{{ t('assessment.experienceLevel') }}</h3>
        <p class="step-hint">{{ t('assessment.experienceLevelHint') }}</p>
        
        <van-radio-group v-model="formData.experience_level" direction="vertical">
          <van-cell-group inset>
            <van-cell clickable @click="formData.experience_level = 'beginner'">
              <template #title>
                <van-radio name="beginner">
                  <span class="option-title">{{ t('assessment.experienceLevels.beginner') }}</span>
                </van-radio>
              </template>
              <template #label>
                <span class="option-desc">{{ t('assessment.experienceLevels.beginnerDesc') }}</span>
              </template>
            </van-cell>
            <van-cell clickable @click="formData.experience_level = 'intermediate'">
              <template #title>
                <van-radio name="intermediate">
                  <span class="option-title">{{ t('assessment.experienceLevels.intermediate') }}</span>
                </van-radio>
              </template>
              <template #label>
                <span class="option-desc">{{ t('assessment.experienceLevels.intermediateDesc') }}</span>
              </template>
            </van-cell>
            <van-cell clickable @click="formData.experience_level = 'advanced'">
              <template #title>
                <van-radio name="advanced">
                  <span class="option-title">{{ t('assessment.experienceLevels.advanced') }}</span>
                </van-radio>
              </template>
              <template #label>
                <span class="option-desc">{{ t('assessment.experienceLevels.advancedDesc') }}</span>
              </template>
            </van-cell>
          </van-cell-group>
        </van-radio-group>
      </div>

      <!-- Step 2: Availability -->
      <div v-show="currentStep === 2" class="step-panel">
        <h3 class="step-title">{{ t('assessment.availability') }}</h3>
        <p class="step-hint">{{ t('assessment.availabilityHint') }}</p>
        
        <van-cell-group inset>
          <van-field
            v-model="formData.weekly_available_days"
            type="digit"
            :label="t('assessment.daysPerWeek')"
            :placeholder="'1-7'"
            :rules="[{ required: true, message: t('assessment.validation.daysRequired') }]"
          />
          <van-field
            v-model="formData.daily_available_minutes"
            type="digit"
            :label="t('assessment.minutesPerDay')"
            :placeholder="'30-120'"
            :rules="[{ required: true, message: t('assessment.validation.minutesRequired') }]"
          />
        </van-cell-group>

        <h4 class="subsection-title">{{ t('assessment.preferredDays') }}</h4>
        <van-checkbox-group v-model="formData.preferred_days" direction="horizontal" class="days-group">
          <van-checkbox
            v-for="day in weekDays"
            :key="day.value"
            :name="day.value"
            shape="square"
          >
            {{ day.label }}
          </van-checkbox>
        </van-checkbox-group>
      </div>

      <!-- Step 3: Activity Type & Equipment -->
      <div v-show="currentStep === 3" class="step-panel">
        <h3 class="step-title">{{ t('assessment.activityType') }}</h3>
        <p class="step-hint">{{ t('assessment.activityTypeHint') }}</p>
        
        <van-radio-group v-model="formData.activity_type" direction="vertical">
          <van-cell-group inset>
            <van-cell
              v-for="activity in activityTypes"
              :key="activity.value"
              clickable
              @click="formData.activity_type = activity.value"
            >
              <template #title>
                <van-radio :name="activity.value">{{ activity.label }}</van-radio>
              </template>
            </van-cell>
          </van-cell-group>
        </van-radio-group>

        <h4 class="subsection-title">{{ t('assessment.equipment') }}</h4>
        <p class="step-hint">{{ t('assessment.equipmentHint') }}</p>
        <van-checkbox-group v-model="formData.equipment_available" class="equipment-group">
          <van-cell-group inset>
            <van-cell
              v-for="equip in equipmentOptions"
              :key="equip.value"
              clickable
              @click="toggleEquipment(equip.value)"
            >
              <template #title>
                <van-checkbox :name="equip.value" shape="square">{{ equip.label }}</van-checkbox>
              </template>
            </van-cell>
          </van-cell-group>
        </van-checkbox-group>
      </div>

      <!-- Step 4: Health Information -->
      <div v-show="currentStep === 4" class="step-panel">
        <h3 class="step-title">{{ t('assessment.healthInfo') }}</h3>
        <p class="step-hint">{{ t('assessment.healthInfoHint') }}</p>
        
        <h4 class="subsection-title">{{ t('assessment.injuryHistory') }}</h4>
        <van-checkbox-group v-model="selectedInjuries" class="injury-group">
          <van-cell-group inset>
            <van-cell
              v-for="injury in commonInjuries"
              :key="injury.value"
              clickable
              @click="toggleInjury(injury.value)"
            >
              <template #title>
                <van-checkbox :name="injury.value" shape="square">{{ injury.label }}</van-checkbox>
              </template>
            </van-cell>
          </van-cell-group>
        </van-checkbox-group>

        <van-cell-group inset class="notes-group">
          <van-field
            v-model="formData.injury_history"
            type="textarea"
            :label="t('assessment.injuryHistory')"
            :placeholder="t('assessment.injuryHistoryPlaceholder')"
            rows="2"
            autosize
          />
          <van-field
            v-model="formData.health_conditions"
            type="textarea"
            :label="t('assessment.healthConditions')"
            :placeholder="t('assessment.healthConditionsPlaceholder')"
            rows="2"
            autosize
          />
        </van-cell-group>
      </div>
    </div>

    <!-- Navigation buttons -->
    <div class="nav-buttons">
      <van-button
        v-if="currentStep > 1"
        plain
        type="default"
        @click="previousStep"
      >
        {{ t('assessment.previous') }}
      </van-button>
      <van-button
        v-if="currentStep < totalSteps"
        type="primary"
        :disabled="!canProceed"
        @click="nextStep"
      >
        {{ t('assessment.next') }}
      </van-button>
      <van-button
        v-if="currentStep === totalSteps"
        type="primary"
        :loading="loading"
        :disabled="!canSubmit"
        @click="handleSubmit"
      >
        {{ loading ? t('assessment.submitting') : t('assessment.submit') }}
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { sanitizeInput } from '@/utils/sanitizer'

const { t } = useI18n()

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  initialData: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['submit', 'cancel'])

// State
const currentStep = ref(1)
const totalSteps = 4
const selectedInjuries = ref([])

// Form data
const formData = reactive({
  experience_level: props.initialData.experience_level || '',
  weekly_available_days: props.initialData.weekly_available_days?.toString() || '',
  daily_available_minutes: props.initialData.daily_available_minutes?.toString() || '',
  activity_type: props.initialData.activity_type || '',
  injury_history: props.initialData.injury_history || '',
  health_conditions: props.initialData.health_conditions || '',
  preferred_days: props.initialData.preferred_days || [],
  equipment_available: props.initialData.equipment_available || []
})

// Computed
const progressPercentage = computed(() => {
  return Math.round((currentStep.value / totalSteps) * 100)
})

const canProceed = computed(() => {
  switch (currentStep.value) {
    case 1:
      return !!formData.experience_level
    case 2:
      return formData.weekly_available_days && formData.daily_available_minutes
    case 3:
      return !!formData.activity_type
    case 4:
      return true
    default:
      return false
  }
})

const canSubmit = computed(() => {
  return (
    formData.experience_level &&
    formData.weekly_available_days &&
    formData.daily_available_minutes &&
    formData.activity_type
  )
})

// Options
const weekDays = computed(() => [
  { value: 'monday', label: t('assessment.days.monday') },
  { value: 'tuesday', label: t('assessment.days.tuesday') },
  { value: 'wednesday', label: t('assessment.days.wednesday') },
  { value: 'thursday', label: t('assessment.days.thursday') },
  { value: 'friday', label: t('assessment.days.friday') },
  { value: 'saturday', label: t('assessment.days.saturday') },
  { value: 'sunday', label: t('assessment.days.sunday') }
])

const activityTypes = computed(() => [
  { value: 'strength_training', label: t('assessment.activityTypes.strength_training') },
  { value: 'cardio', label: t('assessment.activityTypes.cardio') },
  { value: 'mixed', label: t('assessment.activityTypes.mixed') },
  { value: 'flexibility', label: t('assessment.activityTypes.flexibility') }
])

const equipmentOptions = computed(() => [
  { value: 'none', label: t('assessment.equipmentOptions.none') },
  { value: 'dumbbells', label: t('assessment.equipmentOptions.dumbbells') },
  { value: 'barbell', label: t('assessment.equipmentOptions.barbell') },
  { value: 'kettlebell', label: t('assessment.equipmentOptions.kettlebell') },
  { value: 'resistance_bands', label: t('assessment.equipmentOptions.resistance_bands') },
  { value: 'pull_up_bar', label: t('assessment.equipmentOptions.pull_up_bar') },
  { value: 'bench', label: t('assessment.equipmentOptions.bench') },
  { value: 'cable_machine', label: t('assessment.equipmentOptions.cable_machine') },
  { value: 'full_gym', label: t('assessment.equipmentOptions.full_gym') }
])

const commonInjuries = computed(() => [
  { value: 'none', label: t('assessment.commonInjuries.none') },
  { value: 'knee', label: t('assessment.commonInjuries.knee') },
  { value: 'back', label: t('assessment.commonInjuries.back') },
  { value: 'shoulder', label: t('assessment.commonInjuries.shoulder') },
  { value: 'wrist', label: t('assessment.commonInjuries.wrist') },
  { value: 'ankle', label: t('assessment.commonInjuries.ankle') },
  { value: 'neck', label: t('assessment.commonInjuries.neck') },
  { value: 'hip', label: t('assessment.commonInjuries.hip') }
])

// Methods
const nextStep = () => {
  if (currentStep.value < totalSteps && canProceed.value) {
    currentStep.value++
  }
}

const previousStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

const handleBack = () => {
  if (currentStep.value > 1) {
    previousStep()
  } else {
    emit('cancel')
  }
}

const toggleEquipment = (value) => {
  const index = formData.equipment_available.indexOf(value)
  if (index === -1) {
    // If selecting 'none', clear other selections
    if (value === 'none') {
      formData.equipment_available = ['none']
    } else {
      // Remove 'none' if selecting other equipment
      const noneIndex = formData.equipment_available.indexOf('none')
      if (noneIndex !== -1) {
        formData.equipment_available.splice(noneIndex, 1)
      }
      formData.equipment_available.push(value)
    }
  } else {
    formData.equipment_available.splice(index, 1)
  }
}

const toggleInjury = (value) => {
  const index = selectedInjuries.value.indexOf(value)
  if (index === -1) {
    // If selecting 'none', clear other selections
    if (value === 'none') {
      selectedInjuries.value = ['none']
    } else {
      // Remove 'none' if selecting other injuries
      const noneIndex = selectedInjuries.value.indexOf('none')
      if (noneIndex !== -1) {
        selectedInjuries.value.splice(noneIndex, 1)
      }
      selectedInjuries.value.push(value)
    }
  } else {
    selectedInjuries.value.splice(index, 1)
  }
}

const handleSubmit = () => {
  if (!canSubmit.value) return

  // Build injury history from selected injuries
  let injuryText = sanitizeInput(formData.injury_history)
  if (selectedInjuries.value.length > 0 && !selectedInjuries.value.includes('none')) {
    const injuryLabels = selectedInjuries.value.map(v => {
      const injury = commonInjuries.value.find(i => i.value === v)
      return injury ? injury.label : v
    })
    injuryText = injuryLabels.join(', ') + (injuryText ? '. ' + injuryText : '')
  }

  const submitData = {
    experience_level: formData.experience_level,
    weekly_available_days: parseInt(formData.weekly_available_days),
    daily_available_minutes: parseInt(formData.daily_available_minutes),
    activity_type: formData.activity_type,
    assessment_date: new Date().toISOString().slice(0, 10),
    injury_history: injuryText || null,
    health_conditions: sanitizeInput(formData.health_conditions) || null,
    preferred_days: formData.preferred_days.length > 0 ? formData.preferred_days : null,
    equipment_available: formData.equipment_available.length > 0 ? formData.equipment_available : null
  }

  emit('submit', submitData)
}
</script>

<style scoped>
.assessment-form {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--van-background);
}

.step-indicator {
  font-size: 12px;
  color: var(--van-text-color-2);
}

.step-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.step-panel {
  padding: 0 16px;
}

.step-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--van-text-color);
}

.step-hint {
  font-size: 14px;
  color: var(--van-text-color-2);
  margin: 0 0 16px 0;
}

.subsection-title {
  font-size: 16px;
  font-weight: 500;
  margin: 24px 0 12px 0;
  color: var(--van-text-color);
}

.option-title {
  font-weight: 500;
}

.option-desc {
  font-size: 12px;
  color: var(--van-text-color-3);
}

.days-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px;
  background: var(--van-background-2);
  border-radius: 8px;
}

.days-group :deep(.van-checkbox) {
  margin: 0;
}

.equipment-group,
.injury-group {
  margin-top: 8px;
}

.notes-group {
  margin-top: 16px;
}

.nav-buttons {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: var(--van-background-2);
  border-top: 1px solid var(--van-border-color);
}

.nav-buttons .van-button {
  flex: 1;
}

:deep(.van-progress) {
  border-radius: 0;
}

:deep(.van-cell-group--inset) {
  margin: 0;
}

:deep(.van-radio__label) {
  display: flex;
  flex-direction: column;
}
</style>
