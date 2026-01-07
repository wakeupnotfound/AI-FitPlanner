<template>
  <div class="login-view">
    <div class="login-header">
      <h1 class="login-title">{{ t('auth.loginTitle') }}</h1>
      <p class="login-subtitle">{{ t('auth.loginSubtitle') }}</p>
    </div>

    <div class="login-form-container">
      <van-form @submit="handleSubmit" class="login-form">
        <!-- Username Field -->
        <van-cell-group inset>
          <van-field
            v-model="formData.username"
            name="username"
            :label="t('auth.username')"
            :placeholder="t('auth.username')"
            autocomplete="username"
            clearable
            left-icon="user-o"
          />

          <!-- Password Field -->
          <van-field
            v-model="formData.password"
            type="password"
            name="password"
            :label="t('auth.password')"
            :placeholder="t('auth.password')"
            autocomplete="current-password"
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
        <div class="login-actions">
          <van-button
            round
            block
            type="primary"
            native-type="submit"
            :loading="loading"
            :disabled="loading"
            size="large"
          >
            {{ t('auth.login') }}
          </van-button>
        </div>
      </van-form>

      <!-- Register Link -->
      <div class="login-footer">
        <span>{{ t('auth.noAccount') }}</span>
        <router-link to="/register" class="register-link">
          {{ t('auth.register') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { showToast } from 'vant'
import { useAuth } from '../composables/useAuth'

const { t } = useI18n()
const { login, loading, error } = useAuth()

// Form data
const formData = reactive({
  username: '',
  password: ''
})

/**
 * Validate form before submission
 * @returns {boolean} True if form is valid
 */
const validateForm = () => {
  if (!formData.username.trim()) {
    showToast({
      message: t('auth.usernameRequired'),
      type: 'fail'
    })
    return false
  }

  if (!formData.password) {
    showToast({
      message: t('auth.passwordRequired'),
      type: 'fail'
    })
    return false
  }

  if (formData.password.length < 6) {
    showToast({
      message: t('auth.passwordMinLength'),
      type: 'fail'
    })
    return false
  }

  return true
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
    await login({
      username: formData.username.trim(),
      password: formData.password
    })
    
    showToast({
      message: t('auth.loginSuccess'),
      type: 'success'
    })
  } catch (err) {
    // Error is already set in useAuth composable
    showToast({
      message: err.response?.data?.message || t('error.loginFailed'),
      type: 'fail'
    })
  }
}
</script>

<style scoped>
.login-view {
  min-height: 100vh;
  padding: 40px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
  color: #fff;
}

.login-title {
  font-size: 32px;
  font-weight: 600;
  margin: 0 0 12px;
}

.login-subtitle {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
}

.login-form-container {
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
}

.login-form {
  background: #fff;
  border-radius: 16px;
  padding: 24px 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.login-form :deep(.van-cell-group--inset) {
  margin: 0;
}

.login-form :deep(.van-field) {
  padding: 14px 16px;
  font-size: 16px;
}

.login-form :deep(.van-field__label) {
  font-size: 16px;
}

.login-form :deep(.van-field__left-icon) {
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

.login-actions {
  margin-top: 24px;
  padding: 0 16px;
}

.login-actions :deep(.van-button) {
  height: 48px;
  font-size: 18px;
}

.login-actions :deep(.van-button--primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  color: #fff;
  font-size: 16px;
}

.register-link {
  color: #fff;
  font-weight: 600;
  text-decoration: underline;
  margin-left: 4px;
}
</style>
