import apiClient from './api'

/**
 * Statistics Service
 * Handles all statistics, body trends, and progress operations
 */
export const statisticsService = {
  /**
   * Fetch training statistics
   * @param {Object} [params] - Query parameters
   * @param {string} [params.period] - Period (week, month, quarter, year, all)
   * @returns {Promise<Object>} Response with training statistics
   */
  async fetchStatistics(params = { period: 'week' }) {
    return apiClient.get('/stats/training', { params })
  },

  /**
   * Fetch body data trends
   * @param {Object} [params] - Query parameters
   * @param {string} [params.start_date] - Start date filter (YYYY-MM-DD)
   * @param {string} [params.end_date] - End date filter (YYYY-MM-DD)
   * @param {string} [params.metric] - Metric type (weight, body_fat, muscle_mass)
   * @returns {Promise<Object>} Response with body trends data
   */
  async fetchBodyTrends(params = {}) {
    return apiClient.get('/stats/trends', { params })
  },

  /**
   * Calculate progress towards goals
   * @param {Object} currentData - Current body data
   * @param {Object} goalData - Goal data
   * @returns {Object} Progress calculation result
   */
  calculateProgress(currentData, goalData, history = []) {
    if (!currentData || !goalData) {
      return {
        percentage: 0,
        comparison: {}
      }
    }

    const historyItems = Array.isArray(history) ? history : []
    const goalCreatedAt = goalData.created_at ? new Date(goalData.created_at) : null

    const getInitialValue = (key, fallback) => {
      if (!historyItems.length) return fallback
      const filtered = goalCreatedAt
        ? historyItems.filter((item) => {
            if (!item.measurement_date) return false
            const itemDate = new Date(item.measurement_date)
            return itemDate <= goalCreatedAt
          })
        : historyItems
      const sorted = filtered
        .slice()
        .sort((a, b) => new Date(a.measurement_date) - new Date(b.measurement_date))
      for (const item of sorted) {
        const value = item[key]
        if (value !== null && value !== undefined && value !== '') {
          return Number(value)
        }
      }
      return fallback
    }

    const comparison = {}
    const percentages = []

    // Calculate weight progress
    if (goalData.target_weight && currentData.weight) {
      const initial = goalData.initial_weight || getInitialValue('weight', currentData.weight)
      const target = goalData.target_weight
      const current = currentData.weight
      
      const totalChange = target - initial
      const currentChange = current - initial
      const percentage = totalChange !== 0 ? (currentChange / totalChange) * 100 : 0
      
      comparison.weight = {
        initial,
        current,
        target,
        percentage: Math.min(100, Math.max(0, percentage)),
        remaining: Math.abs(target - current),
        direction: target > initial ? 'increase' : 'decrease'
      }
      percentages.push(comparison.weight.percentage)
    }

    // Calculate body fat progress
    const currentBodyFat = currentData.body_fat_percentage ?? currentData.body_fat
    if (goalData.target_body_fat && currentBodyFat !== null && currentBodyFat !== undefined) {
      const initial = goalData.initial_body_fat || getInitialValue('body_fat_percentage', currentBodyFat)
      const target = goalData.target_body_fat
      const current = currentBodyFat
      
      const totalChange = target - initial
      const currentChange = current - initial
      const percentage = totalChange !== 0 ? (currentChange / totalChange) * 100 : 0
      
      comparison.body_fat = {
        initial,
        current,
        target,
        percentage: Math.min(100, Math.max(0, percentage)),
        remaining: Math.abs(target - current),
        direction: target > initial ? 'increase' : 'decrease'
      }
      percentages.push(comparison.body_fat.percentage)
    }

    // Calculate muscle mass progress
    const currentMuscle = currentData.muscle_percentage ?? currentData.muscle_mass
    if (goalData.target_muscle_mass && currentMuscle !== null && currentMuscle !== undefined) {
      const initial = goalData.initial_muscle_mass || getInitialValue('muscle_percentage', currentMuscle)
      const target = goalData.target_muscle_mass
      const current = currentMuscle
      
      const totalChange = target - initial
      const currentChange = current - initial
      const percentage = totalChange !== 0 ? (currentChange / totalChange) * 100 : 0
      
      comparison.muscle_mass = {
        initial,
        current,
        target,
        percentage: Math.min(100, Math.max(0, percentage)),
        remaining: Math.abs(target - current),
        direction: target > initial ? 'increase' : 'decrease'
      }
      percentages.push(comparison.muscle_mass.percentage)
    }

    // Calculate overall progress percentage
    const overallPercentage = percentages.length > 0
      ? percentages.reduce((sum, p) => sum + p, 0) / percentages.length
      : 0

    return {
      percentage: Math.round(overallPercentage),
      comparison
    }
  },

  /**
   * Fetch dashboard summary data
   * Uses training stats as dashboard summary since backend doesn't have dedicated dashboard endpoint
   * @returns {Promise<Object>} Response with dashboard summary
   */
  async fetchDashboardSummary() {
    return apiClient.get('/stats/training', { params: { period: 'week' } })
  },

  /**
   * Fetch progress data from API
   * @returns {Promise<Object>} Response with progress data
   */
  async fetchProgress() {
    return apiClient.get('/stats/progress')
  },

  /**
   * Get date range for statistics queries
   * @param {string} range - Range type ('week', 'month', '3months', 'year')
   * @returns {Object} Object with start_date and end_date
   */
  getDateRange(range) {
    const end = new Date()
    const start = new Date()

    switch (range) {
      case 'week':
        start.setDate(end.getDate() - 7)
        break
      case 'month':
        start.setMonth(end.getMonth() - 1)
        break
      case '3months':
        start.setMonth(end.getMonth() - 3)
        break
      case 'year':
        start.setFullYear(end.getFullYear() - 1)
        break
      default:
        start.setMonth(end.getMonth() - 1)
    }

    return {
      start_date: start.toISOString().split('T')[0],
      end_date: end.toISOString().split('T')[0]
    }
  },

  /**
   * Format statistics data for chart display
   * @param {Array} data - Raw statistics data
   * @param {string} valueKey - Key for the value field
   * @returns {Array} Formatted data for charts
   */
  formatChartData(data, valueKey) {
    if (!Array.isArray(data)) {
      return []
    }

    return data.map(item => ({
      date: item.date || item.record_date,
      value: item[valueKey] || item.value || 0,
      label: item.label || ''
    }))
  }
}

export default statisticsService
