import { createI18n } from 'vue-i18n'
import en from './en.json'
import zh from './zh.json'

const messages = {
  en,
  zh
}

// Detect browser language
const getBrowserLanguage = () => {
  const lang = navigator.language || navigator.userLanguage
  return lang.startsWith('zh') ? 'zh' : 'en'
}

// Get stored language or use browser language
const getInitialLanguage = () => {
  return localStorage.getItem('language') || getBrowserLanguage()
}

const i18n = createI18n({
  legacy: false,
  locale: getInitialLanguage(),
  fallbackLocale: 'en',
  messages,
  globalInjection: true
})

export default i18n
