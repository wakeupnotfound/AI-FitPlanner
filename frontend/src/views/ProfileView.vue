<template>
  <div class="profile-view">
    <!-- Header with user info -->
    <div class="profile-header">
      <van-image
        round
        width="80"
        height="80"
        :src="userAvatar"
        fit="cover"
        class="avatar"
      >
        <template #error>
          <van-icon name="user-o" size="40" />
        </template>
      </van-image>
      <div class="user-info">
        <h2 class="username">{{ profile?.nickname || profile?.username || t('profile.title') }}</h2>
        <p class="email">{{ profile?.email }}</p>
        <p class="member-since" v-if="profile?.created_at">
          {{ t('profile.memberSince') }}: {{ formatDate(profile.created_at) }}
        </p>
      </div>
      <van-button
        icon="edit"
        size="small"
        type="primary"
        plain
        class="touch-feedback"
        @click="showEditProfile = true"
      >
        {{ t('profile.editProfile') }}
      </van-button>
    </div>

    <!-- Pull to Refresh Container -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- Body Data Section -->
    <van-cell-group :title="t('profile.bodyData')" inset class="section">
      <template v-if="latestBodyData">
        <van-cell :title="t('profile.latestMeasurement')" :label="formatDate(latestBodyData.measurement_date)">
          <template #value>
            <van-button size="small" plain @click="showBodyDataHistory = true">
              {{ t('profile.viewHistory') }}
            </van-button>
          </template>
        </van-cell>
        <van-grid :column-num="4" :border="false">
          <van-grid-item>
            <div class="data-value">{{ latestBodyData.weight || '-' }}</div>
            <div class="data-label">{{ t('profile.weight') }} ({{ t('profile.kg') }})</div>
          </van-grid-item>
          <van-grid-item>
            <div class="data-value">{{ latestBodyData.height || '-' }}</div>
            <div class="data-label">{{ t('profile.height') }} ({{ t('profile.cm') }})</div>
          </van-grid-item>
          <van-grid-item>
            <div class="data-value">{{ latestBodyData.body_fat_percentage || '-' }}</div>
            <div class="data-label">{{ t('profile.bodyFat') }} ({{ t('profile.percent') }})</div>
          </van-grid-item>
          <van-grid-item>
            <div class="data-value">{{ latestBodyData.muscle_percentage || '-' }}</div>
            <div class="data-label">{{ t('profile.muscleMass') }} ({{ t('profile.percent') }})</div>
          </van-grid-item>
        </van-grid>
      </template>
      <van-empty v-else :description="t('profile.noBodyData')" image="search" />
      <van-cell>
        <van-button block type="primary" @click="showBodyDataForm = true">
          {{ t('profile.addBodyData') }}
        </van-button>
      </van-cell>
    </van-cell-group>

    <!-- Fitness Goals Section -->
    <van-cell-group :title="t('profile.goals')" inset class="section">
      <template v-if="goals">
        <van-cell :title="t('profile.goalType')" :value="formatGoalType(goals.goal_type)" />
        <van-cell v-if="goals.target_weight" :title="t('profile.targetWeight')" :value="`${goals.target_weight} ${t('profile.kg')}`" />
        <van-cell v-if="goals.target_body_fat" :title="t('profile.targetBodyFat')" :value="`${goals.target_body_fat}${t('profile.percent')}`" />
        <van-cell v-if="goals.target_date" :title="t('profile.targetDate')" :value="formatDate(goals.target_date)" />
        <van-cell v-if="goals.notes" :title="t('profile.notes')" :value="goals.notes" />
        <van-cell>
          <van-button block plain type="primary" @click="showGoalsForm = true">
            {{ t('profile.editGoals') }}
          </van-button>
        </van-cell>
      </template>
      <template v-else>
        <van-empty :description="t('profile.noGoals')" image="search" />
        <van-cell>
          <van-button block type="primary" @click="showGoalsForm = true">
            {{ t('profile.setGoals') }}
          </van-button>
        </van-cell>
      </template>
    </van-cell-group>

    <!-- Quick Actions -->
    <van-cell-group inset class="section">
      <van-cell :title="t('settings.language')" icon="globe-o" class="touch-target-list-item">
        <template #value>
          <LanguageSwitcher />
        </template>
      </van-cell>
      <van-cell :title="t('profile.aiConfig')" is-link to="/ai-config" icon="setting-o" class="touch-target-list-item" />
      <van-cell :title="t('profile.logout')" icon="revoke" @click="handleLogout" class="logout-cell touch-target-list-item" />
    </van-cell-group>
    </van-pull-refresh>

    <!-- Edit Profile Popup -->
    <van-popup v-model:show="showEditProfile" position="bottom" round :style="{ height: '60%' }">
      <div class="popup-content">
        <van-nav-bar :title="t('profile.editProfile')" left-arrow @click-left="showEditProfile = false">
          <template #right>
            <van-button type="primary" size="small" :loading="saving" @click="saveProfile">
              {{ t('app.save') }}
            </van-button>
          </template>
        </van-nav-bar>
        <van-form ref="profileFormRef">
          <van-cell-group inset>
            <van-cell :title="t('profile.avatar')">
              <div class="avatar-edit">
                <van-image
                  round
                  width="64"
                  height="64"
                  :src="editForm.avatar || userAvatar"
                  fit="cover"
                  class="avatar-preview"
                >
                  <template #error>
                    <van-icon name="user-o" size="32" />
                  </template>
                </van-image>
                <div class="avatar-actions">
                  <van-uploader
                    accept="image/*"
                    :max-count="1"
                    :after-read="onAvatarRead"
                  >
                    <van-button size="small" plain type="primary">
                      {{ t('profile.uploadAvatar') }}
                    </van-button>
                  </van-uploader>
                  <van-button
                    size="small"
                    plain
                    type="default"
                    @click="clearAvatar"
                  >
                    {{ t('profile.removeAvatar') }}
                  </van-button>
                </div>
              </div>
            </van-cell>
            <van-field
              v-model="editForm.nickname"
              :label="t('profile.nickname')"
              :placeholder="t('profile.nickname')"
            />
            <van-field
              v-model="editForm.phone"
              :label="t('profile.phone')"
              :placeholder="t('profile.phone')"
              type="tel"
            />
            <van-field
              v-model="editForm.avatar"
              :label="t('profile.avatar')"
              :placeholder="t('profile.avatar')"
            />
          </van-cell-group>
        </van-form>
      </div>
    </van-popup>

    <!-- Body Data Form Popup -->
    <van-popup v-model:show="showBodyDataForm" position="bottom" round :style="{ height: '70%' }">
      <BodyDataForm
        @submit="handleBodyDataSubmit"
        @cancel="showBodyDataForm = false"
        :loading="savingBodyData"
      />
    </van-popup>

    <!-- Body Data History Popup -->
    <van-popup v-model:show="showBodyDataHistory" position="bottom" round :style="{ height: '80%' }">
      <div class="popup-content">
        <van-nav-bar :title="t('profile.bodyDataHistory')" left-arrow @click-left="showBodyDataHistory = false" />
        <div class="history-list">
          <van-empty v-if="!bodyDataHistory.length" :description="t('profile.noBodyData')" />
          <van-cell-group v-else inset v-for="(item, index) in bodyDataHistory" :key="index">
            <van-cell :title="formatDate(item.measurement_date)" />
            <van-grid :column-num="4" :border="false">
              <van-grid-item>
                <div class="data-value">{{ item.weight || '-' }}</div>
                <div class="data-label">{{ t('profile.weight') }}</div>
              </van-grid-item>
              <van-grid-item>
                <div class="data-value">{{ item.height || '-' }}</div>
                <div class="data-label">{{ t('profile.height') }}</div>
              </van-grid-item>
              <van-grid-item>
                <div class="data-value">{{ item.body_fat_percentage || '-' }}</div>
                <div class="data-label">{{ t('profile.bodyFat') }}</div>
              </van-grid-item>
              <van-grid-item>
                <div class="data-value">{{ item.muscle_percentage || '-' }}</div>
                <div class="data-label">{{ t('profile.muscleMass') }}</div>
              </van-grid-item>
            </van-grid>
          </van-cell-group>
        </div>
      </div>
    </van-popup>

    <!-- Goals Form Popup -->
    <van-popup v-model:show="showGoalsForm" position="bottom" round :style="{ height: '70%' }">
      <div class="popup-content">
        <van-nav-bar :title="goals ? t('profile.editGoals') : t('profile.setGoals')" left-arrow @click-left="showGoalsForm = false">
          <template #right>
            <van-button type="primary" size="small" :loading="savingGoals" @click="saveGoals">
              {{ t('app.save') }}
            </van-button>
          </template>
        </van-nav-bar>
        <van-form ref="goalsFormRef">
          <van-cell-group inset>
            <van-field
              v-model="goalsForm.goal_type"
              is-link
              readonly
              :label="t('profile.goalType')"
              :placeholder="t('profile.goalType')"
              @click="showGoalTypePicker = true"
            />
            <van-field
              v-model="goalsForm.target_weight"
              :label="t('profile.targetWeight')"
              :placeholder="t('profile.targetWeight')"
              type="number"
            />
            <van-field
              v-model="goalsForm.target_body_fat"
              :label="t('profile.targetBodyFat')"
              :placeholder="t('profile.targetBodyFat')"
              type="number"
            />
            <van-field
              v-model="goalsForm.target_date"
              is-link
              readonly
              :label="t('profile.targetDate')"
              :placeholder="t('profile.targetDate')"
              @click="showDatePicker = true"
            />
            <van-field
              v-model="goalsForm.notes"
              :label="t('profile.notes')"
              :placeholder="t('profile.notes')"
              type="textarea"
              rows="2"
            />
          </van-cell-group>
        </van-form>
      </div>
    </van-popup>

    <!-- Goal Type Picker -->
    <van-popup v-model:show="showGoalTypePicker" position="bottom" round>
      <van-picker
        :columns="goalTypeOptions"
        @confirm="onGoalTypeConfirm"
        @cancel="showGoalTypePicker = false"
      />
    </van-popup>

    <!-- Date Picker -->
    <van-popup v-model:show="showDatePicker" position="bottom" round>
      <van-date-picker
        v-model="selectedDate"
        :min-date="minDate"
        :max-date="maxDate"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>

    <!-- Loading Overlay -->
    <van-overlay :show="loading" class="loading-overlay">
      <LoadingSpinner />
    </van-overlay>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { showToast, showConfirmDialog } from 'vant'
import { useUserStore } from '@/stores/user'
import { useAuth } from '@/composables/useAuth'
import { useLocale } from '@/composables/useLocale'
import { userService } from '@/services/user.service'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import BodyDataForm from '@/components/fitness/BodyDataForm.vue'
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue'

const { t } = useI18n()
const { formatDate } = useLocale()
const userStore = useUserStore()
const { logout } = useAuth()

const formatGoalType = (goalType) => {
  if (!goalType) {
    return t('profile.goalTypeUnknown')
  }
  return t(`profile.goalTypes.${goalType}`)
}

// State
const loading = ref(false)
const refreshing = ref(false)
const saving = ref(false)
const savingBodyData = ref(false)
const savingGoals = ref(false)

// Popup visibility
const showEditProfile = ref(false)
const showBodyDataForm = ref(false)
const showBodyDataHistory = ref(false)
const showGoalsForm = ref(false)
const showGoalTypePicker = ref(false)
const showDatePicker = ref(false)

// Form refs
const profileFormRef = ref(null)
const goalsFormRef = ref(null)

// Computed
const profile = computed(() => userStore.profile)
const latestBodyData = computed(() => userStore.latestBodyData)
const bodyDataHistory = computed(() => userStore.bodyDataHistory)
const goals = computed(() => userStore.goals)

const userAvatar = computed(() => {
  return profile.value?.avatar || ''
})

// Edit form
const editForm = reactive({
  nickname: '',
  phone: '',
  avatar: ''
})

// Goals form
const goalsForm = reactive({
  goal_type: '',
  target_weight: '',
  target_body_fat: '',
  target_date: '',
  notes: ''
})

// Date picker
const selectedDate = ref([])
const minDate = new Date()
const maxDate = new Date(new Date().setFullYear(new Date().getFullYear() + 2))

// Goal type options
const goalTypeOptions = computed(() => [
  { text: t('profile.goalTypes.weight_loss'), value: 'weight_loss' },
  { text: t('profile.goalTypes.muscle_gain'), value: 'muscle_gain' },
  { text: t('profile.goalTypes.endurance'), value: 'endurance' },
  { text: t('profile.goalTypes.general_fitness'), value: 'general_fitness' }
])

// Methods
const onRefresh = async () => {
  try {
    await loadData()
  } finally {
    refreshing.value = false
  }
}

const loadData = async () => {
  loading.value = true
  try {
    await Promise.all([
      userStore.fetchProfile(),
      userStore.fetchBodyData(),
      userStore.fetchGoals()
    ])
  } catch (error) {
    console.error('Failed to load profile data:', error)
  } finally {
    loading.value = false
  }
}

const initEditForm = () => {
  editForm.nickname = profile.value?.nickname || ''
  editForm.phone = profile.value?.phone || ''
  editForm.avatar = profile.value?.avatar || ''
}

const initGoalsForm = () => {
  if (goals.value) {
    goalsForm.goal_type = goals.value.goal_type || ''
    goalsForm.target_weight = goals.value.target_weight?.toString() || ''
    goalsForm.target_body_fat = goals.value.target_body_fat?.toString() || ''
    goalsForm.target_date = goals.value.target_date || ''
    goalsForm.notes = goals.value.notes || ''
  } else {
    goalsForm.goal_type = ''
    goalsForm.target_weight = ''
    goalsForm.target_body_fat = ''
    goalsForm.target_date = ''
    goalsForm.notes = ''
  }
}

const saveProfile = async () => {
  saving.value = true
  try {
    await userStore.updateProfile(editForm)
    showToast({ type: 'success', message: t('profile.updateSuccess') })
    showEditProfile.value = false
  } catch (error) {
    showToast({ type: 'fail', message: t('profile.updateFailed') })
  } finally {
    saving.value = false
  }
}

const handleBodyDataSubmit = async (bodyData) => {
  savingBodyData.value = true
  try {
    await userStore.addBodyData(bodyData)
    showToast({ type: 'success', message: t('profile.bodyDataSuccess') })
    showBodyDataForm.value = false
  } catch (error) {
    showToast({ type: 'fail', message: t('profile.bodyDataFailed') })
  } finally {
    savingBodyData.value = false
  }
}

const saveGoals = async () => {
  savingGoals.value = true
  try {
    const goalsData = {
      goal_type: goalsForm.goal_type,
      target_weight: goalsForm.target_weight ? parseFloat(goalsForm.target_weight) : null,
      target_body_fat: goalsForm.target_body_fat ? parseFloat(goalsForm.target_body_fat) : null,
      target_date: goalsForm.target_date || null,
      notes: goalsForm.notes || null
    }
    await userStore.setGoals(goalsData)
    showToast({ type: 'success', message: t('profile.goalsSuccess') })
    showGoalsForm.value = false
  } catch (error) {
    showToast({ type: 'fail', message: t('profile.goalsFailed') })
  } finally {
    savingGoals.value = false
  }
}

const onGoalTypeConfirm = ({ selectedOptions }) => {
  goalsForm.goal_type = selectedOptions[0]?.value || ''
  showGoalTypePicker.value = false
}

const onDateConfirm = ({ selectedValues }) => {
  goalsForm.target_date = selectedValues.join('-')
  showDatePicker.value = false
}

const readFileAsDataURL = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

const loadImage = (src) => {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => resolve(img)
    img.onerror = reject
    img.src = src
  })
}

const compressAvatar = async (payload) => {
  const sourceFile = payload?.file
  let dataUrl = payload?.content

  if (!dataUrl && sourceFile) {
    dataUrl = await readFileAsDataURL(sourceFile)
  }

  if (!dataUrl) {
    throw new Error('No image data')
  }

  const img = await loadImage(dataUrl)
  const maxSize = 512
  const scale = Math.min(1, maxSize / Math.max(img.width, img.height))
  const targetWidth = Math.round(img.width * scale)
  const targetHeight = Math.round(img.height * scale)

  const canvas = document.createElement('canvas')
  canvas.width = targetWidth
  canvas.height = targetHeight
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    throw new Error('No canvas context')
  }
  ctx.drawImage(img, 0, 0, targetWidth, targetHeight)

  return canvas.toDataURL('image/jpeg', 0.8)
}

const onAvatarRead = async (file) => {
  const payload = Array.isArray(file) ? file[0] : file
  try {
    editForm.avatar = await compressAvatar(payload)
  } catch (error) {
    showToast({ type: 'fail', message: t('profile.avatarCompressFailed') })
  }
}

const clearAvatar = () => {
  editForm.avatar = ''
}

const handleLogout = async () => {
  try {
    await showConfirmDialog({
      title: t('profile.logout'),
      message: t('app.confirm') + '?'
    })
    await logout()
  } catch {
    // User cancelled
  }
}

// Watch for popup open to init forms
const onEditProfileOpen = () => {
  initEditForm()
}

const onGoalsFormOpen = () => {
  initGoalsForm()
}

// Lifecycle
onMounted(() => {
  loadData()
})

// Watch popup visibility to init forms
import { watch } from 'vue'
watch(showEditProfile, (val) => {
  if (val) onEditProfileOpen()
})
watch(showGoalsForm, (val) => {
  if (val) onGoalsFormOpen()
})
</script>

<style scoped>
.profile-view {
  padding-bottom: 80px;
  background-color: var(--van-background);
  min-height: 100vh;
}

.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 16px;
  background: linear-gradient(135deg, #1f4d7a, #17304d);
  color: white;
  gap: 12px;
}

.avatar {
  border: 3px solid white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.user-info {
  text-align: center;
}

.username {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.email {
  margin: 4px 0;
  font-size: 14px;
  opacity: 0.95;
}

.member-since {
  margin: 0;
  font-size: 12px;
  opacity: 0.85;
}

.section {
  margin-top: 16px;
}

.data-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--van-primary-color);
}

.data-label {
  font-size: 12px;
  color: var(--van-text-color);
  margin-top: 4px;
}

.avatar-edit {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.avatar-preview {
  border: 2px solid rgba(255, 255, 255, 0.6);
}

.avatar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.popup-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.history-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.logout-cell {
  color: var(--van-danger-color);
}

.loading-overlay {
  display: flex;
  align-items: center;
  justify-content: center;
}

:deep(.van-grid-item__content) {
  padding: 12px 8px;
}

:deep(.van-cell-group__title) {
  padding-left: 16px;
}
</style>
