<template>
  <div class="body-data-form">
    <van-nav-bar
      :title="t('profile.addBodyData')"
      left-arrow
      @click-left="$emit('cancel')"
    >
      <template #right>
        <van-button
          type="primary"
          size="small"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ t('app.save') }}
        </van-button>
      </template>
    </van-nav-bar>

    <van-form ref="formRef" @submit="handleSubmit">
      <van-cell-group inset>
        <!-- Measurement Date -->
        <van-field
          v-model="formData.measurement_date"
          is-link
          readonly
          required
          :label="t('profile.measurementDate')"
          :placeholder="t('profile.measurementDate')"
          :rules="[{ required: true, message: t('profile.measurementDate') }]"
          @click="showDatePicker = true"
        />

        <!-- Weight -->
        <van-field
          v-model="formData.weight"
          required
          :label="t('profile.weight')"
          :placeholder="`${t('profile.weight')} (${t('profile.kg')})`"
          type="number"
          inputmode="decimal"
          enterkeyhint="next"
          :rules="weightRules"
        >
          <template #button>
            <span class="unit">{{ t('profile.kg') }}</span>
          </template>
        </van-field>

        <!-- Height -->
        <van-field
          v-model="formData.height"
          required
          :label="t('profile.height')"
          :placeholder="`${t('profile.height')} (${t('profile.cm')})`"
          type="number"
          inputmode="decimal"
          enterkeyhint="next"
          :rules="heightRules"
        >
          <template #button>
            <span class="unit">{{ t('profile.cm') }}</span>
          </template>
        </van-field>

        <!-- Age -->
        <van-field
          v-model="formData.age"
          :label="t('profile.age')"
          :placeholder="t('profile.age')"
          type="number"
          inputmode="numeric"
          enterkeyhint="next"
          :rules="ageRules"
        />

        <!-- Gender -->
        <van-field
          v-model="genderDisplay"
          is-link
          readonly
          :label="t('profile.gender')"
          :placeholder="t('profile.gender')"
          enterkeyhint="next"
          @click="showGenderPicker = true"
        />

        <!-- Body Fat Percentage -->
        <van-field
          v-model="formData.body_fat_percentage"
          :label="t('profile.bodyFat')"
          :placeholder="`${t('profile.bodyFat')} (${t('profile.percent')})`"
          type="number"
          inputmode="decimal"
          enterkeyhint="next"
          :rules="percentageRules"
        >
          <template #button>
            <span class="unit">{{ t('profile.percent') }}</span>
          </template>
        </van-field>

        <!-- Muscle Percentage -->
        <van-field
          v-model="formData.muscle_percentage"
          :label="t('profile.muscleMass')"
          :placeholder="`${t('profile.muscleMass')} (${t('profile.percent')})`"
          type="number"
          inputmode="decimal"
          enterkeyhint="done"
          :rules="percentageRules"
        >
          <template #button>
            <span class="unit">{{ t('profile.percent') }}</span>
          </template>
        </van-field>
      </van-cell-group>
    </van-form>

    <!-- Date Picker Popup -->
    <van-popup v-model:show="showDatePicker" position="bottom" round>
      <van-date-picker
        v-model="selectedDate"
        :min-date="minDate"
        :max-date="maxDate"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>

    <!-- Gender Picker Popup -->
    <van-popup v-model:show="showGenderPicker" position="bottom" round>
      <van-picker
        :columns="genderOptions"
        @confirm="onGenderConfirm"
        @cancel="showGenderPicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'cancel'])

// Form ref
const formRef = ref(null)

// Picker visibility
const showDatePicker = ref(false)
const showGenderPicker = ref(false)

// Form data
const formData = reactive({
  measurement_date: formatDateToString(new Date()),
  weight: '',
  height: '',
  age: '',
  gender: '',
  body_fat_percentage: '',
  muscle_percentage: ''
})

// Date picker state
const today = new Date()
const selectedDate = ref([
  today.getFullYear().toString(),
  (today.getMonth() + 1).toString().padStart(2, '0'),
  today.getDate().toString().padStart(2, '0')
])
const minDate = new Date(new Date().setFullYear(new Date().getFullYear() - 1))
const maxDate = new Date()

// Gender options
const genderOptions = computed(() => [
  { text: t('profile.male'), value: 'male' },
  { text: t('profile.female'), value: 'female' }
])

// Gender display text
const genderDisplay = computed(() => {
  if (!formData.gender) return ''
  return formData.gender === 'male' ? t('profile.male') : t('profile.female')
})

// Validation rules
const weightRules = [
  { required: true, message: t('profile.weight') },
  {
    validator: (val) => {
      const num = parseFloat(val)
      return !isNaN(num) && num > 0 && num < 500
    },
    message: '0-500'
  }
]

const heightRules = [
  { required: true, message: t('profile.height') },
  {
    validator: (val) => {
      const num = parseFloat(val)
      return !isNaN(num) && num > 0 && num < 300
    },
    message: '0-300'
  }
]

const ageRules = [
  {
    validator: (val) => {
      if (!val) return true
      const num = parseInt(val)
      return !isNaN(num) && num > 0 && num < 150
    },
    message: '1-150'
  }
]

const percentageRules = [
  {
    validator: (val) => {
      if (!val) return true
      const num = parseFloat(val)
      return !isNaN(num) && num >= 0 && num <= 100
    },
    message: '0-100'
  }
]

// Helper functions
function formatDateToString(date) {
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  return `${year}-${month}-${day}`
}

// Event handlers
const onDateConfirm = ({ selectedValues }) => {
  formData.measurement_date = selectedValues.join('-')
  showDatePicker.value = false
}

const onGenderConfirm = ({ selectedOptions }) => {
  formData.gender = selectedOptions[0]?.value || ''
  showGenderPicker.value = false
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    
    // Build submission data with proper types
    const submitData = {
      measurement_date: formData.measurement_date,
      weight: parseFloat(formData.weight),
      height: parseFloat(formData.height)
    }

    // Add optional fields if provided
    if (formData.age) {
      submitData.age = parseInt(formData.age)
    }
    if (formData.gender) {
      submitData.gender = formData.gender
    }
    if (formData.body_fat_percentage) {
      submitData.body_fat_percentage = parseFloat(formData.body_fat_percentage)
    }
    if (formData.muscle_percentage) {
      submitData.muscle_percentage = parseFloat(formData.muscle_percentage)
    }

    emit('submit', submitData)
  } catch (error) {
    // Validation failed
    console.log('Validation failed:', error)
  }
}
</script>

<style scoped>
.body-data-form {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--van-background);
}

.unit {
  color: var(--van-text-color-2);
  font-size: 14px;
}

:deep(.van-cell-group) {
  margin-top: 16px;
}

:deep(.van-field__label) {
  width: 80px;
}
</style>
