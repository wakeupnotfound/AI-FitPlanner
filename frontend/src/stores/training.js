import { defineStore } from 'pinia'
import apiClient from '../services/api'

/**
 * Training Store
 * Manages training plans, workouts, and training history
 */
export const useTrainingStore = defineStore('training', {
  state: () => ({
    plans: [],
    currentPlan: null,
    todayWorkout: null,
    history: [],
    loading: false,
    error: null,
    isGenerating: false
  }),

  getters: {
    /**
     * Get all training plans
     */
    allPlans: (state) => state.plans,

    /**
     * Get active training plan
     */
    activePlan: (state) => {
      return state.currentPlan || state.plans.find(plan => plan.status === 'active')
    },

    /**
     * Get today's workout
     */
    getTodayWorkout: (state) => state.todayWorkout,

    /**
     * Get training history sorted by date (newest first)
     */
    trainingHistory: (state) => {
      const history = Array.isArray(state.history) ? state.history : []
      return [...history].sort((a, b) => 
        new Date(b.workout_date) - new Date(a.workout_date)
      )
    },

    /**
     * Check if plan is being generated
     */
    isGeneratingPlan: (state) => state.isGenerating,

    /**
     * Get completed workouts count
     */
    completedWorkoutsCount: (state) => state.history.length,

    /**
     * Get workouts by date range
     */
    getWorkoutsByDateRange: (state) => (startDate, endDate) => {
      return state.history.filter(workout => {
        const workoutDate = new Date(workout.workout_date)
        return workoutDate >= new Date(startDate) && workoutDate <= new Date(endDate)
      })
    }
  },

  actions: {
    /**
     * Generate new training plan
     * @param {Object} assessmentData - Assessment data for plan generation
     */
    async generatePlan(assessmentData) {
      this.isGenerating = true
      this.error = null

      try {
        const response = await apiClient.post('/training-plans/generate', assessmentData)
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.isGenerating = false
      }
    },

    /**
     * Fetch all training plans
     */
    async fetchPlans() {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/training-plans')
        const plans = response.data?.plans
        this.plans = Array.isArray(plans) ? plans : []
        
        // Set current plan to active plan if exists
        const activePlan = this.plans.find(plan => plan.status === 'active')
        if (activePlan) {
          this.currentPlan = activePlan
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
     * Fetch specific training plan
     * @param {number} planId - Plan ID
     */
    async fetchPlan(planId) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get(`/training-plans/${planId}`)
        const plan = response.data?.plan || response.data
        
        // Update plan in local state
        if (plan) {
          const planIndex = this.plans.findIndex(p => p.id === planId)
          if (planIndex !== -1) {
            this.plans[planIndex] = plan
          } else {
            this.plans.push(plan)
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
     * Fetch today's workout
     */
    async fetchTodayWorkout() {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/training-plans/today')
        this.todayWorkout = response.data?.schedule || null
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Record workout completion
     * @param {Object} workoutData - Workout record data
     * @param {string} workoutData.workout_date - Date of workout
     * @param {Array} workoutData.exercises - Array of exercise records
     * @param {number} workoutData.duration - Duration in minutes
     * @param {string} workoutData.notes - Optional notes
     */
    async recordWorkout(workoutData) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.post('/training-records', workoutData)
        const record = response.data?.record
        if (record) {
          const existingIndex = this.history.findIndex(item => item.id === record.id)
          if (existingIndex !== -1) {
            this.history.splice(existingIndex, 1)
          }
          this.history.unshift(record)
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
     * Fetch training history
     * @param {Object} params - Query parameters
     * @param {string} params.start_date - Start date
     * @param {string} params.end_date - End date
     */
    async fetchHistory(params = {}) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/training-records', { params })
        const records = response.data?.records
        this.history = Array.isArray(records) ? records : []
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Update training plan status
     * @param {number} planId - Plan ID
     * @param {string} status - New status (active, completed, paused)
     */
    async updatePlanStatus(planId, status) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.patch(`/training-plans/${planId}`, { status })
        
        // Update plan in local state
        const planIndex = this.plans.findIndex(p => p.id === planId)
        if (planIndex !== -1) {
          this.plans[planIndex].status = status
        }
        
        // Update current plan if it's the one being updated
        if (this.currentPlan?.id === planId) {
          this.currentPlan.status = status
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
     * Clear training data
     */
    clearTrainingData() {
      this.plans = []
      this.currentPlan = null
      this.todayWorkout = null
      this.history = []
      this.loading = false
      this.error = null
      this.isGenerating = false
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'training',
        storage: sessionStorage,
        paths: ['currentPlan', 'todayWorkout']
      }
    ]
  }
})
