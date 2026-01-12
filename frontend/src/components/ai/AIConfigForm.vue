<template>
  <div class="ai-config-form">
    <!-- Header -->
    <div class="form-header">
      <h3 class="form-title">
        {{ config ? t('ai.editConfig') : t('ai.addConfig') }}
      </h3>
    </div>

    <!-- Form -->
    <van-form @submit="handleSubmit" class="form-content">
      <!-- Provider Selection -->
      <van-field
        v-model="formData.providerLabel"
        is-link
        readonly
        :label="t('ai.provider')"
        :placeholder="t('ai.validation.providerRequired')"
        :rules="[{ required: true, message: t('ai.validation.providerRequired') }]"
        @click="showProviderPicker = true"
      />

      <!-- Configuration Name -->
      <van-field
        v-model="formData.name"
        :label="t('ai.name')"
        :placeholder="t('ai.validation.nameRequired')"
        :rules="[{ required: true, message: t('ai.validation.nameRequired') }]"
        maxlength="50"
        show-word-limit
        autocapitalize="words"
        inputmode="text"
        enterkeyhint="next"
      />

      <!-- API Key with Masking -->
      <van-field
        v-model="formData.api_key"
        :type="showApiKey ? 'text' : 'password'"
        :label="t('ai.apiKey')"
        :placeholder="t('ai.apiKeyPlaceholder')"
        :rules="apiKeyRules"
        autocapitalize="none"
        autocorrect="off"
        spellcheck="false"
        inputmode="text"
        enterkeyhint="next"
      >
        <template #right-icon>
          <van-icon
            :name="showApiKey ? 'eye-o' : 'closed-eye'"
            class="touch-target-icon"
            @click="showApiKey = !showApiKey"
          />
        </template>
      </van-field>
      <div class="field-hint">{{ t('ai.apiKeyHint') }}</div>

      <!-- API Endpoint -->
      <van-field
        v-model="formData.api_endpoint"
        :label="t('ai.apiEndpoint')"
        placeholder="https://api.openai.com/v1"
        autocapitalize="none"
        autocorrect="off"
        inputmode="url"
        enterkeyhint="next"
      />

      <!-- Model Selection -->
      <van-field
        v-model="formData.modelLabel"
        is-link
        readonly
        :label="t('ai.model')"
        :placeholder="t('ai.validation.modelRequired')"
        :rules="[{ required: true, message: t('ai.validation.modelRequired') }]"
        @click="showModelPicker = true"
      />

      <!-- Max Tokens -->
      <van-field
        v-model="formData.max_tokens"
        type="digit"
        inputmode="numeric"
        enterkeyhint="next"
        :label="t('ai.maxTokens')"
        placeholder="2000"
      />

      <!-- Temperature -->
      <van-field
        v-model="formData.temperature"
        type="number"
        inputmode="decimal"
        enterkeyhint="done"
        :label="t('ai.temperature')"
        placeholder="0.7"
      />

      <!-- Set as Default -->
      <van-cell center :title="t('ai.isDefault')">
        <template #right-icon>
          <van-switch v-model="formData.is_default" />
        </template>
      </van-cell>

      <!-- Test Connection Button -->
      <div class="test-connection" v-if="config">
        <van-button
          plain
          type="primary"
          size="small"
          :loading="testing"
          :loading-text="t('ai.testing')"
          @click="handleTestConnection"
        >
          {{ t('ai.testConnection') }}
        </van-button>
      </div>

      <!-- Submit Buttons -->
      <div class="form-actions">
        <van-button
          block
          type="default"
          @click="$emit('cancel')"
        >
          {{ t('app.cancel') }}
        </van-button>
        <van-button
          block
          type="primary"
          native-type="submit"
          :loading="loading"
        >
          {{ config ? t('app.update') : t('app.add') }}
        </van-button>
      </div>
    </van-form>

    <!-- Provider Picker -->
    <van-popup :show="showProviderPicker" position="bottom" round @update:show="showProviderPicker = $event">
      <van-picker
        :columns="providerOptions"
        @confirm="onProviderConfirm"
        @cancel="showProviderPicker = false"
        :confirm-button-text="t('app.confirm')"
        :cancel-button-text="t('app.cancel')"
      />
    </van-popup>

    <!-- Model Picker -->
    <van-popup :show="showModelPicker" position="bottom" round @update:show="showModelPicker = $event">
      <van-picker
        :columns="modelOptions"
        @confirm="onModelConfirm"
        @cancel="showModelPicker = false"
        :confirm-button-text="t('app.confirm')"
        :cancel-button-text="t('app.cancel')"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  config: {
    type: Object,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'cancel', 'test'])

// State
const showApiKey = ref(false)
const showProviderPicker = ref(false)
const showModelPicker = ref(false)
const testing = ref(false)

// Form data
const formData = reactive({
  provider: '',
  providerLabel: '',
  name: '',
  api_key: '',
  api_endpoint: '',
  model: '',
  modelLabel: '',
  max_tokens: '',
  temperature: '',
  is_default: false
})

// Provider options
const providers = [
  { value: 'openai', text: 'OpenAI', endpoint: 'https://api.openai.com/v1' },
  { value: 'wenxin', text: '文心一言（百度）', endpoint: 'https://aip.baidubce.com' },
  { value: 'tongyi', text: '通义千问（阿里）', endpoint: 'https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions' }
]

const providerOptions = computed(() => 
  providers.map(p => ({ text: p.text, value: p.value }))
)

// Model options based on provider
const modelsByProvider = {
  openai: [
    { value: 'gpt-4-turbo', text: 'GPT-4 Turbo' },
    { value: 'gpt-4', text: 'GPT-4' },
    { value: 'gpt-3.5-turbo', text: 'GPT-3.5 Turbo' }
  ],
  wenxin: [
    { value: 'ernie-bot-4', text: '文心大模型4.0' },
    { value: 'ernie-bot', text: '文心大模型' }
  ],
  tongyi: [
    { value: 'qwen-turbo', text: '通义千问Turbo' },
    { value: 'qwen-plus', text: '通义千问Plus' }
  ]
}

const modelOptions = computed(() => {
  const models = modelsByProvider[formData.provider] || []
  return models.map(m => ({ text: m.text, value: m.value }))
})

// API key validation rules
const apiKeyRules = computed(() => {
  // If editing and no new key provided, don't require it
  if (props.config && !formData.api_key) {
    return []
  }
  return [{ required: true, message: t('ai.validation.apiKeyRequired') }]
})

// Watch for config changes (edit mode)
watch(() => props.config, (newConfig) => {
  if (newConfig) {
    formData.provider = newConfig.provider || ''
    formData.providerLabel = providers.find(p => p.value === newConfig.provider)?.text || ''
    formData.name = newConfig.name || ''
    formData.api_key = '' // Don't show existing key
    formData.api_endpoint = newConfig.api_endpoint || ''
    formData.model = newConfig.model || ''
    formData.modelLabel = modelsByProvider[newConfig.provider]?.find(m => m.value === newConfig.model)?.text || ''
    formData.max_tokens = newConfig.max_tokens?.toString() || ''
    formData.temperature = newConfig.temperature?.toString() || ''
    formData.is_default = newConfig.is_default || false
  } else {
    resetForm()
  }
}, { immediate: true })

// Methods
function onProviderConfirm({ selectedOptions }) {
  const selected = selectedOptions[0]
  formData.provider = selected.value
  formData.providerLabel = selected.text
  
  // Set default endpoint
  const provider = providers.find(p => p.value === selected.value)
  if (provider && !formData.api_endpoint) {
    formData.api_endpoint = provider.endpoint
  }
  
  // Reset model when provider changes
  formData.model = ''
  formData.modelLabel = ''
  
  showProviderPicker.value = false
}

function onModelConfirm({ selectedOptions }) {
  const selected = selectedOptions[0]
  formData.model = selected.value
  formData.modelLabel = selected.text
  showModelPicker.value = false
}

function handleSubmit() {
  const submitData = {
    provider: formData.provider,
    name: formData.name,
    api_endpoint: formData.api_endpoint,
    model: formData.model,
    max_tokens: formData.max_tokens ? parseInt(formData.max_tokens) : 2000,
    temperature: formData.temperature ? parseFloat(formData.temperature) : 0.7,
    is_default: formData.is_default
  }
  
  // Only include api_key if provided
  if (formData.api_key) {
    submitData.api_key = formData.api_key
  }
  
  emit('submit', submitData)
}

function handleTestConnection() {
  if (props.config) {
    emit('test', props.config.id)
  }
}

function resetForm() {
  formData.provider = ''
  formData.providerLabel = ''
  formData.name = ''
  formData.api_key = ''
  formData.api_endpoint = ''
  formData.model = ''
  formData.modelLabel = ''
  formData.max_tokens = ''
  formData.temperature = ''
  formData.is_default = false
}
</script>

<style scoped>
.ai-config-form {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #fff;
}

.form-header {
  padding: 16px;
  border-bottom: 1px solid #ebedf0;
}

.form-title {
  font-size: 18px;
  font-weight: 600;
  color: #323233;
  margin: 0;
  text-align: center;
}

.form-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 100px;
}

.field-hint {
  padding: 4px 16px 12px;
  font-size: 12px;
  color: #969799;
}

.test-connection {
  padding: 12px 16px;
  display: flex;
  justify-content: center;
}

.form-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  background-color: #fff;
  border-top: 1px solid #ebedf0;
}

.form-actions .van-button {
  flex: 1;
}
</style>
