import { defineStore } from 'pinia'
import apiClient from '../services/api'

/**
 * Nutrition Store
 * Manages nutrition plans, meals, and meal history
 */
export const useNutritionStore = defineStore('nutrition', {
  state: () => ({
    plans: [],
    currentPlan: null,
    todayMeals: null,
    history: [],
    loading: false,
    error: null,
    isGenerating: false
  }),

  getters: {
    /**
     * Get all nutrition plans
     */
    allPlans: (state) => state.plans,

    /**
     * Get active nutrition plan
     */
    activePlan: (state) => {
      return state.currentPlan || state.plans.find(plan => plan.status === 'active')
    },

    /**
     * Get today's meals
     */
    getTodayMeals: (state) => state.todayMeals,

    /**
     * Get meal history sorted by date (newest first)
     */
    mealHistory: (state) => {
      return [...state.history].sort((a, b) => 
        new Date(b.meal_date) - new Date(a.meal_date)
      )
    },

    /**
     * Get meal history grouped by date
     */
    mealHistoryGroupedByDate: (state) => {
      const grouped = {}
      
      state.history.forEach(meal => {
        const date = meal.meal_date.split('T')[0] // Get date part only
        if (!grouped[date]) {
          grouped[date] = []
        }
        grouped[date].push(meal)
      })
      
      return grouped
    },

    /**
     * Check if plan is being generated
     */
    isGeneratingPlan: (state) => state.isGenerating,

    /**
     * Get total calories for today
     */
    todayTotalCalories: (state) => {
      if (!state.todayMeals || !state.todayMeals.meals) {
        return 0
      }
      return state.todayMeals.meals.reduce((total, meal) => {
        return total + (meal.total_calories || 0)
      }, 0)
    },

    /**
     * Get meals by date range
     */
    getMealsByDateRange: (state) => (startDate, endDate) => {
      return state.history.filter(meal => {
        const mealDate = new Date(meal.meal_date)
        return mealDate >= new Date(startDate) && mealDate <= new Date(endDate)
      })
    }
  },

  actions: {
    /**
     * Generate new nutrition plan
     * @param {Object} planData - Plan generation data
     * @param {Array} planData.dietary_restrictions - Dietary restrictions
     * @param {number} planData.daily_calories - Target daily calories
     * @param {Object} planData.macro_ratios - Macro nutrient ratios
     */
    async generatePlan(planData) {
      this.isGenerating = true
      this.error = null

      try {
        const response = await apiClient.post('/nutrition-plans/generate', planData)
        const newPlan = response.data
        
        this.plans.push(newPlan)
        this.currentPlan = newPlan
        
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.isGenerating = false
      }
    },

    /**
     * Fetch all nutrition plans
     */
    async fetchPlans() {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/nutrition-plans')
        this.plans = response.data || []
        
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
     * Fetch specific nutrition plan
     * @param {number} planId - Plan ID
     */
    async fetchPlan(planId) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get(`/nutrition-plans/${planId}`)
        
        // Update plan in local state
        const planIndex = this.plans.findIndex(p => p.id === planId)
        if (planIndex !== -1) {
          this.plans[planIndex] = response.data
        } else {
          this.plans.push(response.data)
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
     * Fetch today's meals
     */
    async fetchTodayMeals() {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/nutrition-plans/today')
        this.todayMeals = response.data
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Record meal
     * @param {Object} mealData - Meal record data
     * @param {string} mealData.meal_date - Date and time of meal
     * @param {string} mealData.meal_type - Type of meal (breakfast, lunch, dinner, snack)
     * @param {Array} mealData.foods - Array of food items
     * @param {number} mealData.total_calories - Total calories
     */
    async recordMeal(mealData) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.post('/nutrition-records', mealData)
        
        // Add to history
        this.history.unshift(response.data)
        
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch meal history
     * @param {Object} params - Query parameters
     * @param {string} params.start_date - Start date
     * @param {string} params.end_date - End date
     */
    async fetchHistory(params = {}) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/nutrition-records', { params })
        this.history = response.data || []
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Update nutrition plan status
     * @param {number} planId - Plan ID
     * @param {string} status - New status (active, completed, paused)
     */
    async updatePlanStatus(planId, status) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.patch(`/nutrition-plans/${planId}`, { status })
        
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
     * Clear nutrition data
     */
    clearNutritionData() {
      this.plans = []
      this.currentPlan = null
      this.todayMeals = null
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
        key: 'nutrition',
        storage: sessionStorage,
        paths: ['currentPlan', 'todayMeals']
      }
    ]
  }
})
