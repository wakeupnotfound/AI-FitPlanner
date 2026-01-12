import apiClient from './api'

/**
 * Nutrition Service
 * Handles all nutrition plan and meal operations
 */
export const nutritionService = {
  /**
   * Generate a new nutrition plan using AI
   * @param {Object} planData - Plan generation parameters
   * @param {string} planData.plan_name - Plan name
   * @param {number} planData.duration_days - Duration in days
   * @param {number} planData.daily_calories - Target daily calories
   * @param {number} planData.protein_ratio - Protein ratio (0-1)
   * @param {number} planData.carb_ratio - Carb ratio (0-1)
   * @param {number} planData.fat_ratio - Fat ratio (0-1)
   * @param {Array} planData.dietary_restrictions - Dietary restrictions
   * @param {number} [planData.ai_api_id] - AI API configuration ID
   * @returns {Promise<Object>} Response with task_id for async generation
   */
  async generatePlan(planData) {
    return apiClient.post('/nutrition-plans/generate', planData)
  },

  /**
   * Check nutrition plan generation task status
   * @param {string} taskId - Task ID from generatePlan
   * @returns {Promise<Object>} Response with task status and progress
   */
  async checkTaskStatus(taskId) {
    return apiClient.get(`/nutrition-plans/tasks/${taskId}`)
  },

  /**
   * Fetch all nutrition plans for the current user
   * @param {Object} [params] - Query parameters
   * @param {number} [params.page] - Page number
   * @param {number} [params.limit] - Items per page
   * @returns {Promise<Object>} Response with list of nutrition plans
   */
  async fetchPlans(params = {}) {
    return apiClient.get('/nutrition-plans', { params })
  },

  /**
   * Fetch a specific nutrition plan by ID
   * @param {number} planId - Plan ID
   * @returns {Promise<Object>} Response with plan details
   */
  async fetchPlan(planId) {
    return apiClient.get(`/nutrition-plans/${planId}`)
  },

  /**
   * Fetch today's meals schedule
   * @returns {Promise<Object>} Response with today's meal details
   */
  async fetchTodayMeals() {
    return apiClient.get('/nutrition-plans/today')
  },

  /**
   * Record a meal
   * @param {Object} mealData - Meal record data
   * @param {string} mealData.meal_date - Date and time of meal (ISO format)
   * @param {string} mealData.meal_type - Type of meal (breakfast, lunch, dinner, snack)
   * @param {Array} mealData.foods - Array of food items
   * @param {number} mealData.total_calories - Total calories
   * @param {number} [mealData.total_protein] - Total protein in grams
   * @param {number} [mealData.total_carbs] - Total carbs in grams
   * @param {number} [mealData.total_fat] - Total fat in grams
   * @param {string} [mealData.notes] - Optional notes
   * @returns {Promise<Object>} Response with created record
   */
  async recordMeal(mealData) {
    return apiClient.post('/nutrition-records', mealData)
  },

  /**
   * Fetch meal history/records
   * @param {Object} [params] - Query parameters
   * @param {string} [params.start_date] - Start date filter (YYYY-MM-DD)
   * @param {string} [params.end_date] - End date filter (YYYY-MM-DD)
   * @param {number} [params.page] - Page number
   * @param {number} [params.limit] - Items per page
   * @returns {Promise<Object>} Response with meal records
   */
  async fetchHistory(params = {}) {
    return apiClient.get('/nutrition-records', { params })
  },

  /**
   * Update nutrition plan status
   * @param {number} planId - Plan ID
   * @param {string} status - New status (active, completed, paused)
   * @returns {Promise<Object>} Response confirming update
   */
  async updatePlanStatus(planId, status) {
    return apiClient.patch(`/nutrition-plans/${planId}`, { status })
  },

  /**
   * Delete a meal record
   * @param {number} recordId - Meal record ID
   * @returns {Promise<Object>} Response confirming deletion
   */
  async deleteMealRecord(recordId) {
    return apiClient.delete(`/nutrition-records/${recordId}`)
  },

  /**
   * Get nutrition summary for a date range
   * @param {Object} params - Query parameters
   * @param {string} params.start_date - Start date (YYYY-MM-DD)
   * @param {string} params.end_date - End date (YYYY-MM-DD)
   * @returns {Promise<Object>} Response with nutrition summary
   */
  async getNutritionSummary(params) {
    return apiClient.get('/nutrition-records/daily-summary', { params })
  }
}

export default nutritionService
