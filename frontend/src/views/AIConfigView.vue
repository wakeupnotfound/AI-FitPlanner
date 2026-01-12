<template>
  <div class="ai-config-view">
    <!-- Header -->
    <van-nav-bar
      :title="t('ai.title')"
      left-arrow
      @click-left="router.back()"
    >
      <template #right>
        <van-icon name="add-o" size="20" @click="showAddForm = true" />
      </template>
    </van-nav-bar>

    <!-- Loading State -->
    <LoadingSpinner v-if="loading" :loading="loading" centered />

    <!-- Error State -->
    <ErrorMessage
      v-else-if="error"
      :message="error"
      inline
      retryable
      @retry="fetchConfigs"
    />

    <!-- Empty State -->
    <van-empty
      v-else-if="configs.length === 0"
      :description="t('ai.noConfigs')"
      image="search"
    >
      <template #description>
        <p class="empty-description">{{ t('ai.noConfigs') }}</p>
        <p class="empty-hint">{{ t('ai.noConfigsHint') }}</p>
      </template>
      <van-button type="primary" size="small" @click="showAddForm = true">
        {{ t('ai.addConfig') }}
      </van-button>
    </van-empty>

    <!-- Config List -->
    <div v-else class="config-list">
      <AIConfigCard
        v-for="config in configs"
        :key="config.id"
        :config="config"
        :testing="testingId === config.id"
        @test="handleTestConnection"
        @set-default="handleSetDefault"
        @delete="handleDeleteConfirm"
        @edit="handleEdit"
      />
    </div>

    <!-- Add/Edit Form Popup -->
    <van-popup
      :show="showAddForm"
      position="bottom"
      round
      :style="{ height: '85%' }"
      closeable
      @update:show="showAddForm = $event"
      @close="resetForm"
    >
      <AIConfigForm
        :config="editingConfig"
        :loading="formLoading"
        @submit="handleSubmit"
        @cancel="showAddForm = false"
      />
    </van-popup>

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      v-model="showDeleteDialog"
      :title="t('app.confirm')"
      :message="t('ai.deleteConfirm')"
      danger
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { showToast, showSuccessToast } from 'vant'
import { useAIConfigStore } from '@/stores/aiConfig'
import { LoadingSpinner, ErrorMessage, ConfirmDialog } from '@/components/common'
import AIConfigCard from '@/components/ai/AIConfigCard.vue'
import AIConfigForm from '@/components/ai/AIConfigForm.vue'

const router = useRouter()
const { t } = useI18n()
const aiConfigStore = useAIConfigStore()

// State
const loading = ref(false)
const error = ref(null)
const showAddForm = ref(false)
const showDeleteDialog = ref(false)
const editingConfig = ref(null)
const deletingConfigId = ref(null)
const testingId = ref(null)
const formLoading = ref(false)

// Computed
const configs = computed(() => aiConfigStore.allConfigs)

// Methods
async function fetchConfigs() {
  loading.value = true
  error.value = null
  try {
    await aiConfigStore.fetchConfigs()
  } catch (err) {
    error.value = err.message || t('error.unknown')
  } finally {
    loading.value = false
  }
}

async function handleSubmit(formData) {
  formLoading.value = true
  try {
    if (editingConfig.value) {
      await aiConfigStore.updateConfig(editingConfig.value.id, formData)
      showSuccessToast(t('ai.updateSuccess'))
    } else {
      await aiConfigStore.addConfig(formData)
      showSuccessToast(t('ai.addSuccess'))
    }
    showAddForm.value = false
    resetForm()
  } catch (err) {
    showToast(err.message || t('error.unknown'))
  } finally {
    formLoading.value = false
  }
}

async function handleTestConnection(configId) {
  if (!configId) {
    showToast('配置ID无效')
    return
  }
  
  testingId.value = configId
  try {
        const result = await aiConfigStore.testConnection(configId)
        if (result?.data?.test_result?.status === 'success') {
          showSuccessToast(t('ai.testSuccess'))
        } else {
          const errorMsg = result?.test_result?.message || 
                           result?.data?.test_result?.message || 
                           t('ai.testFailed')
          showToast(errorMsg)
        }
      } catch (err) {
        showToast(t('ai.testFailed'))
      } finally {
        testingId.value = null
      }
}

async function handleSetDefault(configId) {
  try {
    await aiConfigStore.setDefault(configId)
    showSuccessToast(t('ai.setDefaultSuccess'))
  } catch (err) {
    showToast(err.message || t('error.unknown'))
  }
}

function handleDeleteConfirm(configId) {
  deletingConfigId.value = configId
  showDeleteDialog.value = true
}

async function handleDelete() {
  if (!deletingConfigId.value) return
  
  try {
    await aiConfigStore.deleteConfig(deletingConfigId.value)
    showSuccessToast(t('ai.deleteSuccess'))
  } catch (err) {
    showToast(err.message || t('error.unknown'))
  } finally {
    deletingConfigId.value = null
    showDeleteDialog.value = false
  }
}

function handleEdit(config) {
  editingConfig.value = config
  showAddForm.value = true
}

function resetForm() {
  editingConfig.value = null
}

// Lifecycle
onMounted(() => {
  fetchConfigs()
})
</script>

<style scoped>
.ai-config-view {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 60px;
}

.config-list {
  padding: 12px;
}

.empty-description {
  font-size: 14px;
  color: #969799;
  margin-bottom: 4px;
}

.empty-hint {
  font-size: 12px;
  color: #c8c9cc;
  margin-bottom: 16px;
}
</style>
