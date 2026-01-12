import apiClient from './api'

/**
 * Training Service
 * Handles all training plan and workout operations
 */
export const trainingService = {
  /**
   * Generate a new training plan using AI
   * @param {Object} planData - Plan generation parameters
   * @param {string} planData.plan_name - Plan name
   * @param {number} planData.duration_weeks - Duration in weeks
   * @param {string} planData.goal - Goal type (muscle_gain, fat_loss, endurance)
   * @param {string} planData.difficulty_level - Difficulty level (easy, medium, hard, extreme)
   * @param {number} [planData.ai_api_id] - AI API configuration ID
   * @returns {Promise<Object>} Response with task_id for async generation
   */
  async generatePlan(planData) {
    return apiClient.post('/training-plans/generate', planData)
  },

  /**
   * Check plan generation task status
   * @param {string} taskId - Task ID from generatePlan
   * @returns {Promise<Object>} Response with task status and progress
   */
  async checkTaskStatus(taskId) {
    return apiClient.get(`/training-plans/tasks/${taskId}`)
  },

  /**
   * Fetch all training plans for the current user
   * @param {Object} [params] - Query parameters
   * @param {number} [params.page] - Page number
   * @param {number} [params.limit] - Items per page
   * @returns {Promise<Object>} Response with list of training plans
   */
  async fetchPlans(params = {}) {
    return apiClient.get('/training-plans', { params })
  },

  /**
   * Fetch a specific training plan by ID
   * @param {number} planId - Plan ID
   * @returns {Promise<Object>} Response with plan details
   */
  async fetchPlan(planId) {
    return apiClient.get(`/training-plans/${planId}`)
  },

  /**
   * Fetch today's workout schedule
   * @returns {Promise<Object>} Response with today's workout details
   */
  async fetchTodayWorkout() {
    return apiClient.get('/training-plans/today')
  },

  /**
   * Record a completed workout
   * @param {Object} workoutData - Workout record data
   * @param {number} workoutData.plan_id - Training plan ID
   * @param {string} workoutData.workout_date - Date of workout (YYYY-MM-DD)
   * @param {string} workoutData.workout_type - Type of workout (strength, cardio, etc.)
   * @param {number} workoutData.duration_minutes - Duration in minutes
   * @param {Array} workoutData.exercises - Array of exercise records
   * @param {Object} [workoutData.performance_data] - Performance metrics
   * @param {string} [workoutData.notes] - Optional notes
   * @param {number} [workoutData.rating] - Workout rating (1-5)
   * @param {string} [workoutData.injury_report] - Injury report if any
   * @returns {Promise<Object>} Response with created record
   */
  async recordWorkout(workoutData) {
    return apiClient.post('/training-records', workoutData)
  },

  /**
   * Fetch training history/records
   * @param {Object} [params] - Query parameters
   * @param {string} [params.start_date] - Start date filter
   * @param {string} [params.end_date] - End date filter
   * @param {number} [params.page] - Page number
   * @param {number} [params.limit] - Items per page
   * @returns {Promise<Object>} Response with training records
   */
  async fetchHistory(params = {}) {
    return apiClient.get('/training-records', { params })
  },

  /**
   * Update training plan status
   * @param {number} planId - Plan ID
   * @param {string} status - New status (active, completed, paused)
   * @returns {Promise<Object>} Response confirming update
   */
  async updatePlanStatus(planId, status) {
    return apiClient.patch(`/training-plans/${planId}`, { status })
  },

  /**
   * Submit fitness assessment
   * @param {Object} assessmentData - Assessment data
   * @param {string} assessmentData.experience_level - Experience level (beginner, intermediate, advanced)
   * @param {number} assessmentData.weekly_available_days - Days available per week
   * @param {number} assessmentData.daily_available_minutes - Minutes available per day
   * @param {string} assessmentData.activity_type - Preferred activity type
   * @param {string} [assessmentData.injury_history] - Injury history description
   * @param {string} [assessmentData.health_conditions] - Health conditions
   * @param {Array} [assessmentData.preferred_days] - Preferred workout days
   * @param {Array} [assessmentData.equipment_available] - Available equipment
   * @returns {Promise<Object>} Response with assessment result
   */
  async submitAssessment(assessmentData) {
    return apiClient.post('/assessments', assessmentData)
  },

  /**
   * Fetch user's latest assessment
   * @returns {Promise<Object>} Response with assessment data
   */
  async fetchAssessment() {
    return apiClient.get('/assessments/latest')
  }
}

export default trainingService
