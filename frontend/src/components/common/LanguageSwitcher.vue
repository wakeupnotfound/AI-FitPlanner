<template>
  <van-dropdown-menu>
    <van-dropdown-item v-model="currentLanguage" :options="languageOptions" @change="handleLanguageChange">
      <template #title>
        <van-icon name="globe-o" />
        <span class="language-label">{{ currentLanguageLabel }}</span>
      </template>
    </van-dropdown-item>
  </van-dropdown-menu>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { showSuccessToast } from 'vant'

const { locale, t } = useI18n()

const languageOptions = [
  { text: 'English', value: 'en' },
  { text: '中文', value: 'zh' }
]

const currentLanguage = ref(locale.value)

const currentLanguageLabel = computed(() => {
  const option = languageOptions.find(opt => opt.value === currentLanguage.value)
  return option ? option.text : 'English'
})

const handleLanguageChange = (value) => {
  locale.value = value
  localStorage.setItem('language', value)
  showSuccessToast(t('settings.languageChanged'))
}
</script>

<style scoped>
.language-label {
  margin-left: 4px;
}
</style>
