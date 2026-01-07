/**
 * Locale Formatting Utilities
 * Provides locale-aware formatting for dates, numbers, and currencies
 */

/**
 * Format a date according to the current locale
 * @param {Date|string|number} date - The date to format
 * @param {string} locale - The locale to use (e.g., 'en', 'zh')
 * @param {Object} options - Intl.DateTimeFormat options
 * @returns {string} Formatted date string
 */
export function formatDate(date, locale = 'en', options = {}) {
  if (!date) return ''
  
  const dateObj = date instanceof Date ? date : new Date(date)
  
  if (isNaN(dateObj.getTime())) {
    return ''
  }
  
  const defaultOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    ...options
  }
  
  try {
    return new Intl.DateTimeFormat(locale, defaultOptions).format(dateObj)
  } catch (error) {
    console.error('Date formatting error:', error)
    return dateObj.toLocaleDateString()
  }
}

/**
 * Format a date as short format (e.g., 2024-01-15)
 * @param {Date|string|number} date - The date to format
 * @param {string} locale - The locale to use
 * @returns {string} Formatted date string
 */
export function formatDateShort(date, locale = 'en') {
  return formatDate(date, locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

/**
 * Format a date with time
 * @param {Date|string|number} date - The date to format
 * @param {string} locale - The locale to use
 * @returns {string} Formatted date and time string
 */
export function formatDateTime(date, locale = 'en') {
  return formatDate(date, locale, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * Format a relative time (e.g., "2 days ago", "in 3 hours")
 * @param {Date|string|number} date - The date to format
 * @param {string} locale - The locale to use
 * @returns {string} Formatted relative time string
 */
export function formatRelativeTime(date, locale = 'en') {
  if (!date) return ''
  
  const dateObj = date instanceof Date ? date : new Date(date)
  
  if (isNaN(dateObj.getTime())) {
    return ''
  }
  
  const now = new Date()
  const diffMs = dateObj.getTime() - now.getTime()
  const diffSec = Math.round(diffMs / 1000)
  const diffMin = Math.round(diffSec / 60)
  const diffHour = Math.round(diffMin / 60)
  const diffDay = Math.round(diffHour / 24)
  
  try {
    const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })
    
    if (Math.abs(diffSec) < 60) {
      return rtf.format(diffSec, 'second')
    } else if (Math.abs(diffMin) < 60) {
      return rtf.format(diffMin, 'minute')
    } else if (Math.abs(diffHour) < 24) {
      return rtf.format(diffHour, 'hour')
    } else if (Math.abs(diffDay) < 30) {
      return rtf.format(diffDay, 'day')
    } else {
      return formatDate(date, locale)
    }
  } catch (error) {
    console.error('Relative time formatting error:', error)
    return formatDate(date, locale)
  }
}

/**
 * Format a number according to the current locale
 * @param {number} value - The number to format
 * @param {string} locale - The locale to use
 * @param {Object} options - Intl.NumberFormat options
 * @returns {string} Formatted number string
 */
export function formatNumber(value, locale = 'en', options = {}) {
  if (value === null || value === undefined || isNaN(value)) {
    return ''
  }
  
  try {
    return new Intl.NumberFormat(locale, options).format(value)
  } catch (error) {
    console.error('Number formatting error:', error)
    return value.toString()
  }
}

/**
 * Format a number with decimal places
 * @param {number} value - The number to format
 * @param {string} locale - The locale to use
 * @param {number} decimals - Number of decimal places
 * @returns {string} Formatted number string
 */
export function formatDecimal(value, locale = 'en', decimals = 2) {
  return formatNumber(value, locale, {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals
  })
}

/**
 * Format a percentage
 * @param {number} value - The value to format (0-1 or 0-100)
 * @param {string} locale - The locale to use
 * @param {boolean} isDecimal - Whether the value is in decimal form (0-1)
 * @returns {string} Formatted percentage string
 */
export function formatPercent(value, locale = 'en', isDecimal = false) {
  if (value === null || value === undefined || isNaN(value)) {
    return ''
  }
  
  const percentValue = isDecimal ? value : value / 100
  
  try {
    return new Intl.NumberFormat(locale, {
      style: 'percent',
      minimumFractionDigits: 0,
      maximumFractionDigits: 1
    }).format(percentValue)
  } catch (error) {
    console.error('Percent formatting error:', error)
    return `${value}%`
  }
}

/**
 * Format a currency value
 * @param {number} value - The value to format
 * @param {string} locale - The locale to use
 * @param {string} currency - The currency code (e.g., 'USD', 'CNY')
 * @returns {string} Formatted currency string
 */
export function formatCurrency(value, locale = 'en', currency = 'USD') {
  if (value === null || value === undefined || isNaN(value)) {
    return ''
  }
  
  try {
    return new Intl.NumberFormat(locale, {
      style: 'currency',
      currency: currency
    }).format(value)
  } catch (error) {
    console.error('Currency formatting error:', error)
    return `${currency} ${value}`
  }
}

/**
 * Format a duration in minutes to a human-readable string
 * @param {number} minutes - Duration in minutes
 * @param {string} locale - The locale to use
 * @returns {string} Formatted duration string
 */
export function formatDuration(minutes, locale = 'en') {
  if (!minutes || isNaN(minutes)) {
    return ''
  }
  
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  
  if (hours > 0) {
    return locale === 'zh' 
      ? `${hours}小时${mins > 0 ? mins + '分钟' : ''}`
      : `${hours}h ${mins > 0 ? mins + 'm' : ''}`
  }
  
  return locale === 'zh' ? `${mins}分钟` : `${mins}m`
}

/**
 * Format a weight value with appropriate unit
 * @param {number} value - Weight value
 * @param {string} locale - The locale to use
 * @param {string} unit - Unit ('kg' or 'lbs')
 * @returns {string} Formatted weight string
 */
export function formatWeight(value, locale = 'en', unit = 'kg') {
  if (value === null || value === undefined || isNaN(value)) {
    return ''
  }
  
  const formattedValue = formatDecimal(value, locale, 1)
  return `${formattedValue} ${unit}`
}

/**
 * Format a distance value with appropriate unit
 * @param {number} value - Distance value in meters
 * @param {string} locale - The locale to use
 * @returns {string} Formatted distance string
 */
export function formatDistance(value, locale = 'en') {
  if (value === null || value === undefined || isNaN(value)) {
    return ''
  }
  
  if (value >= 1000) {
    const km = value / 1000
    const formattedValue = formatDecimal(km, locale, 2)
    return locale === 'zh' ? `${formattedValue}公里` : `${formattedValue} km`
  }
  
  const formattedValue = formatNumber(value, locale)
  return locale === 'zh' ? `${formattedValue}米` : `${formattedValue} m`
}

/**
 * Get locale-specific date format pattern
 * @param {string} locale - The locale to use
 * @returns {string} Date format pattern
 */
export function getDateFormatPattern(locale = 'en') {
  const patterns = {
    'en': 'MM/DD/YYYY',
    'zh': 'YYYY年MM月DD日'
  }
  
  return patterns[locale] || patterns['en']
}

/**
 * Get locale-specific time format pattern
 * @param {string} locale - The locale to use
 * @returns {string} Time format pattern
 */
export function getTimeFormatPattern(locale = 'en') {
  const patterns = {
    'en': 'hh:mm A',
    'zh': 'HH:mm'
  }
  
  return patterns[locale] || patterns['en']
}
