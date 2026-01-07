import { defineStore } from 'pinia'
import apiClient from '../services/api'

/**
 * Statistics Store
 * Manages training statistics, body data trends, and progress tracking
 */
export const useStatisticsStore = defineStore('statistics', {
  state: () => ({
    trainingStats: null,
    bodyTrends: null,
    progress: null,
    loading: false,
    error: null,
    dateRange: {
      start: null,
      end: null
    }
  }),

  getters: {
    /**
     * Get training statistics
     */
    getTrainingStats: (state) => state.trainingStats,

    /**
     * Get body data trends
     */
    getBodyTrends: (state) => state.bodyTrends,

    /**
     * Get progress data
     */
    getProgress: (state) => state.progress,

    /**
     * Check if statistics data is available
     */
    hasStatistics: (state) => {
      return !!(state.trainingStats || state.bodyTrends || state.progress)
    },

    /**
     * Get workout frequency (workouts per week)
     */
    workoutFrequency: (state) => {
      if (!state.trainingStats?.workout_frequency) {
        return 0
      }
      return state.trainingStats.workout_frequency
    },

    /**
     * Get total workouts completed
     */
    totalWorkouts: (state) => {
      if (!state.trainingStats?.total_workouts) {
        return 0
      }
      return state.trainingStats.total_workouts
    },

    /**
     * Get average workout duration
     */
    averageDuration: (state) => {
      if (!state.trainingStats?.average_duration) {
        return 0
      }
      return state.trainingStats.average_duration
    },

    /**
     * Get weight trend data
     */
    weightTrend: (state) => {
      if (!state.bodyTrends?.weight) {
        return []
      }
      return state.bodyTrends.weight
    },

    /**
     * Get body fat trend data
     */
    bodyFatTrend: (state) => {
      if (!state.bodyTrends?.body_fat) {
        return []
      }
      return state.bodyTrends.body_fat
    },

    /**
     * Get muscle mass trend data
     */
    muscleMassTrend: (state) => {
      if (!state.bodyTrends?.muscle_mass) {
        return []
      }
      return state.bodyTrends.muscle_mass
    },

    /**
     * Get progress percentage
     */
    progressPercentage: (state) => {
      if (!state.progress?.percentage) {
        return 0
      }
      return state.progress.percentage
    },

    /**
     * Get goal comparison data
     */
    goalComparison: (state) => {
      if (!state.progress?.comparison) {
        return null
      }
      return state.progress.comparison
    }
  },

  actions: {
    /**
     * Fetch training statistics
     * @param {Object} params - Query parameters
     * @param {string} params.period - Period (week, month, quarter, year, all)
     */
    async fetchStatistics(params = { period: 'week' }) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/stats/training', { params })
        this.trainingStats = response.data
        
        // Update date range if provided
        if (params.start_date && params.end_date) {
          this.dateRange = {
            start: params.start_date,
            end: params.end_date
          }
        }
        
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch body data trends
     * @param {Object} params - Query parameters
     * @param {string} params.start_date - Start date
     * @param {string} params.end_date - End date
     * @param {string} params.metric - Metric type (weight, body_fat, muscle_mass)
     */
    async fetchTrends(params = {}) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/stats/trends', { params })
        this.bodyTrends = response.data
        
        // Update date range if provided
        if (params.start_date && params.end_date) {
          this.dateRange = {
            start: params.start_date,
            end: params.end_date
          }
        }
        
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Calculate progress towards goals
     * @param {Object} currentData - Current body data
     * @param {Object} goalData - Goal data
     */
    async calculateProgress(currentData = null, goalData = null) {
      this.loading = true
      this.error = null

      try {
        let response
        
        if (currentData && goalData) {
          // Calculate locally if data is provided
          response = this.calculateProgressLocally(currentData, goalData)
        } else {
          // Fetch from API
          response = await apiClient.get('/stats/progress')
          this.progress = response.data
        }
        
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Calculate progress locally
     * @param {Object} currentData - Current body data
     * @param {Object} goalData - Goal data
     */
    calculateProgressLocally(currentData, goalData) {
      const comparison = {}
      const percentages = []

      // Calculate for each metric
      if (goalData.target_weight && currentData.weight) {
        const initial = goalData.initial_weight || currentData.weight
        const target = goalData.target_weight
        const current = currentData.weight
        
        const totalChange = target - initial
        const currentChange = current - initial
        const percentage = totalChange !== 0 ? (currentChange / totalChange) * 100 : 0
        
        comparison.weight = {
          initial,
          current,
          target,
          percentage: Math.min(100, Math.max(0, percentage))
        }
        percentages.push(comparison.weight.percentage)
      }

      if (goalData.target_body_fat && currentData.body_fat) {
        const initial = goalData.initial_body_fat || currentData.body_fat
        const target = goalData.target_body_fat
        const current = currentData.body_fat
        
        const totalChange = target - initial
        const currentChange = current - initial
        const percentage = totalChange !== 0 ? (currentChange / totalChange) * 100 : 0
        
        comparison.body_fat = {
          initial,
          current,
          target,
          percentage: Math.min(100, Math.max(0, percentage))
        }
        percentages.push(comparison.body_fat.percentage)
      }

      if (goalData.target_muscle_mass && currentData.muscle_mass) {
        const initial = goalData.initial_muscle_mass || currentData.muscle_mass
        const target = goalData.target_muscle_mass
        const current = currentData.muscle_mass
        
        const totalChange = target - initial
        const currentChange = current - initial
        const percentage = totalChange !== 0 ? (currentChange / totalChange) * 100 : 0
        
        comparison.muscle_mass = {
          initial,
          current,
          target,
          percentage: Math.min(100, Math.max(0, percentage))
        }
        percentages.push(comparison.muscle_mass.percentage)
      }

      // Calculate overall progress percentage
      const overallPercentage = percentages.length > 0
        ? percentages.reduce((sum, p) => sum + p, 0) / percentages.length
        : 0

      this.progress = {
        percentage: Math.round(overallPercentage),
        comparison
      }

      return { data: this.progress }
    },

    /**
     * Fetch dashboard summary
     */
    async fetchDashboardSummary() {
      this.loading = true
      this.error = null

      try {
        // Use training stats as dashboard summary with default period
        const response = await apiClient.get('/stats/training', { 
          params: { period: 'week' } 
        })
        
        // Update training stats from response
        if (response.data) {
          this.trainingStats = response.data
        }
        
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Set date range for statistics
     * @param {string} startDate - Start date
     * @param {string} endDate - End date
     */
    setDateRange(startDate, endDate) {
      this.dateRange = {
        start: startDate,
        end: endDate
      }
    },

    /**
     * Clear statistics data
     */
    clearStatisticsData() {
      this.trainingStats = null
      this.bodyTrends = null
      this.progress = null
      this.loading = false
      this.error = null
      this.dateRange = {
        start: null,
        end: null
      }
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'statistics',
        storage: sessionStorage,
        paths: ['trainingStats', 'bodyTrends', 'progress', 'dateRange']
      }
    ]
  }
})
