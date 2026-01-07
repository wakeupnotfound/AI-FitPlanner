<template>
  <div class="register-view">
    <div class="register-header">
      <h1 class="register-title">{{ t('auth.registerTitle') }}</h1>
      <p class="register-subtitle">{{ t('auth.registerSubtitle') }}</p>
    </div>

    <van-form @submit="handleSubmit" class="register-form">
      <van-cell-group inset>
        <!-- Username Field -->
        <van-field
          v-model="formData.username"
          name="username"
          :label="t('auth.username')"
          :placeholder="t('auth.username')"
          :rules="usernameRules"
          :error-message="fieldErrors.username"
          autocomplete="username"
          autocapitalize="none"
          autocorrect="off"
          spellcheck="false"
          inputmode="text"
          enterkeyhint="next"
          clearable
          left-icon="user-o"
        />

        <!-- Email Field -->
        <van-field
          v-model="formData.email"
          type="email"
          name="email"
          :label="t('auth.email')"
          :placeholder="t('auth.email')"
          :rules="emailRules"
          :error-message="fieldErrors.email"
          autocomplete="email"
          autocapitalize="none"
          autocorrect="off"
          inputmode="email"
          enterkeyhint="next"
          clearable
          left-icon="envelop-o"
        />

        <!-- Nickname Field -->
        <van-field
          v-model="formData.nickname"
          name="nickname"
          :label="t('auth.nickname')"
          :placeholder="t('auth.nickname')"
          :rules="nicknameRules"
          :error-message="fieldErrors.nickname"
          autocomplete="nickname"
          autocapitalize="words"
          inputmode="text"
          enterkeyhint="next"
          clearable
          left-icon="smile-o"
        />

        <!-- Password Field -->
        <van-field
          v-model="formData.password"
          type="password"
          name="password"
          :label="t('auth.password')"
          :placeholder="t('auth.password')"
          :rules="passwordRules"
          :error-message="fieldErrors.password"
          autocomplete="new-password"
          autocapitalize="none"
          autocorrect="off"
          inputmode="text"
          enterkeyhint="next"
          clearable
          left-icon="lock"
        >
          <template #extra>
            <span class="password-hint">至少8位，含大小写字母和数字</span>
          </template>
        </van-field>

        <!-- Confirm Password Field -->
        <van-field
          v-model="formData.confirmPassword"
          type="password"
          name="confirmPassword"
          :label="t('auth.confirmPassword')"
          :placeholder="t('auth.confirmPassword')"
          :rules="confirmPasswordRules"
          :error-message="fieldErrors.confirmPassword"
          autocomplete="new-password"
          autocapitalize="none"
          autocorrect="off"
          inputmode="text"
          enterkeyhint="done"
          clearable
          left-icon="lock"
        />
      </van-cell-group>

      <!-- Error Message -->
      <div v-if="error" class="error-message">
        <van-icon name="warning-o" />
        <span>{{ error }}</span>
      </div>

      <!-- Submit Button -->
      <div class="register-actions">
        <van-button
          round
          block
          type="primary"
          native-type="submit"
          :loading="loading"
          :disabled="loading"
          size="large"
        >
          {{ t('auth.register') }}
        </van-button>
      </div>
    </van-form>

    <!-- Login Link -->
    <div class="register-footer">
      <span>{{ t('auth.hasAccount') }}</span>
      <router-link to="/login" class="login-link">
        {{ t('auth.login') }}
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { showToast } from 'vant'
import { useAuth } from '../composables/useAuth'

const { t } = useI18n()
const { register, loading, error } = useAuth()

// Form data
const formData = reactive({
  username: '',
  email: '',
  nickname: '',
  password: '',
  confirmPassword: ''
})

// Field-level errors
const fieldErrors = reactive({
  username: '',
  email: '',
  nickname: '',
  password: '',
  confirmPassword: ''
})

// Email validation regex
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

// Password strength validation: 至少8个字符，包含大小写字母和数字
const passwordStrengthRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/

// Validation rules
const usernameRules = computed(() => [
  { required: true, message: t('auth.usernameRequired') }
])

const emailRules = computed(() => [
  { required: true, message: t('auth.emailRequired') },
  { 
    validator: (val) => emailRegex.test(val),
    message: t('auth.emailInvalid')
  }
])

const nicknameRules = computed(() => [
  { required: true, message: t('auth.nicknameRequired') }
])

const passwordRules = computed(() => [
  { required: true, message: t('auth.passwordRequired') },
  { 
    validator: (val) => val.length >= 8,
    message: '密码至少需要8个字符'
  },
  { 
    validator: (val) => passwordStrengthRegex.test(val),
    message: '密码必须包含大小写字母和数字'
  }
])

const confirmPasswordRules = computed(() => [
  { required: true, message: t('auth.passwordRequired') },
  { 
    validator: (val) => val === formData.password,
    message: t('auth.passwordMismatch')
  }
])

/**
 * Validate form before submission
 * @returns {boolean} True if form is valid
 */
const validateForm = () => {
  // Reset field errors
  Object.keys(fieldErrors).forEach(key => {
    fieldErrors[key] = ''
  })

  let isValid = true

  if (!formData.username.trim()) {
    fieldErrors.username = t('auth.usernameRequired')
    isValid = false
  }

  if (!formData.email.trim()) {
    fieldErrors.email = t('auth.emailRequired')
    isValid = false
  } else if (!emailRegex.test(formData.email)) {
    fieldErrors.email = t('auth.emailInvalid')
    isValid = false
  }

  if (!formData.nickname.trim()) {
    fieldErrors.nickname = t('auth.nicknameRequired')
    isValid = false
  }

  if (!formData.password) {
    fieldErrors.password = t('auth.passwordRequired')
    isValid = false
  } else if (formData.password.length < 8) {
    fieldErrors.password = '密码至少需要8个字符'
    isValid = false
  } else if (!passwordStrengthRegex.test(formData.password)) {
    fieldErrors.password = '密码必须包含大小写字母和数字'
    isValid = false
  }

  if (!formData.confirmPassword) {
    fieldErrors.confirmPassword = t('auth.passwordRequired')
    isValid = false
  } else if (formData.confirmPassword !== formData.password) {
    fieldErrors.confirmPassword = t('auth.passwordMismatch')
    isValid = false
  }

  return isValid
}

/**
 * Handle form submission
 */
const handleSubmit = async () => {
  // Validate form first
  if (!validateForm()) {
    return
  }

  try {
    await register({
      username: formData.username.trim(),
      email: formData.email.trim(),
      nickname: formData.nickname.trim(),
      password: formData.password,
      confirm_password: formData.confirmPassword // 添加确认密码字段
    })
    
    showToast({
      message: t('auth.registerSuccess'),
      type: 'success'
    })
  } catch (err) {
    // Error is already set in useAuth composable
    showToast({
      message: err.response?.data?.message || t('error.registerFailed'),
      type: 'fail'
    })
  }
}
</script>

<style scoped>
.register-view {
  min-height: 100vh;
  padding: 40px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
  color: #fff;
}

.register-title {
  font-size: 32px;
  font-weight: 600;
  margin: 0 0 12px;
}

.register-subtitle {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
}

.register-form {
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
  background: #fff;
  border-radius: 16px;
  padding: 24px 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.register-form :deep(.van-cell-group--inset) {
  margin: 0;
}

.register-form :deep(.van-field) {
  padding: 14px 16px;
  font-size: 16px;
}

.register-form :deep(.van-field__label) {
  font-size: 16px;
}

.register-form :deep(.van-field__left-icon) {
  margin-right: 8px;
  color: #667eea;
  font-size: 20px;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  margin: 16px 0;
  background: #fff0f0;
  border-radius: 8px;
  color: #ee0a24;
  font-size: 14px;
}

.register-actions {
  margin-top: 24px;
  padding: 0 16px;
}

.register-actions :deep(.van-button) {
  height: 48px;
  font-size: 18px;
}

.register-actions :deep(.van-button--primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.register-footer {
  text-align: center;
  margin-top: 24px;
  color: #fff;
  font-size: 16px;
}

.login-link {
  color: #fff;
  font-weight: 600;
  text-decoration: underline;
  margin-left: 4px;
}

.password-hint {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
}
</style>
