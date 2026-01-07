/**
 * Locale Composable
 * Provides locale-aware formatting functions integrated with Vue I18n
 */

import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  formatDate as _formatDate,
  formatDateShort as _formatDateShort,
  formatDateTime as _formatDateTime,
  formatRelativeTime as _formatRelativeTime,
  formatNumber as _formatNumber,
  formatDecimal as _formatDecimal,
  formatPercent as _formatPercent,
  formatCurrency as _formatCurrency,
  formatDuration as _formatDuration,
  formatWeight as _formatWeight,
  formatDistance as _formatDistance,
  getDateFormatPattern,
  getTimeFormatPattern
} from '@/utils/localeFormatter'

/**
 * Composable for locale-aware formatting
 * @returns {Object} Formatting functions and locale info
 */
export function useLocale() {
  const { locale } = useI18n()
  
  // Current locale value
  const currentLocale = computed(() => locale.value)
  
  // Date formatting functions
  const formatDate = (date, options) => _formatDate(date, currentLocale.value, options)
  const formatDateShort = (date) => _formatDateShort(date, currentLocale.value)
  const formatDateTime = (date) => _formatDateTime(date, currentLocale.value)
  const formatRelativeTime = (date) => _formatRelativeTime(date, currentLocale.value)
  
  // Number formatting functions
  const formatNumber = (value, options) => _formatNumber(value, currentLocale.value, options)
  const formatDecimal = (value, decimals) => _formatDecimal(value, currentLocale.value, decimals)
  const formatPercent = (value, isDecimal) => _formatPercent(value, currentLocale.value, isDecimal)
  const formatCurrency = (value, currency) => _formatCurrency(value, currentLocale.value, currency)
  
  // Specialized formatting functions
  const formatDuration = (minutes) => _formatDuration(minutes, currentLocale.value)
  const formatWeight = (value, unit) => _formatWeight(value, currentLocale.value, unit)
  const formatDistance = (value) => _formatDistance(value, currentLocale.value)
  
  // Format patterns
  const dateFormatPattern = computed(() => getDateFormatPattern(currentLocale.value))
  const timeFormatPattern = computed(() => getTimeFormatPattern(currentLocale.value))
  
  return {
    // Current locale
    currentLocale,
    
    // Date formatting
    formatDate,
    formatDateShort,
    formatDateTime,
    formatRelativeTime,
    
    // Number formatting
    formatNumber,
    formatDecimal,
    formatPercent,
    formatCurrency,
    
    // Specialized formatting
    formatDuration,
    formatWeight,
    formatDistance,
    
    // Format patterns
    dateFormatPattern,
    timeFormatPattern
  }
}
