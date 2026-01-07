import apiClient from './api'

/**
 * User Service
 * Handles user profile, body data, and fitness goals operations
 * Requirements: 2.1, 2.2, 2.4, 2.6
 */
export const userService = {
  /**
   * Fetch user profile
   * @returns {Promise<Object>} User profile data
   * Requirements: 2.1
   */
  async fetchProfile() {
    const response = await apiClient.get('/user/profile')
    return response
  },

  /**
   * Update user profile
   * @param {Object} profileData - Profile data to update
   * @param {string} [profileData.username] - Username
   * @param {string} [profileData.phone] - Phone number
   * @param {string} [profileData.avatar] - Avatar URL
   * @param {string} [profileData.nickname] - Nickname
   * @returns {Promise<Object>} Updated profile data
   * Requirements: 2.2
   */
  async updateProfile(profileData) {
    const response = await apiClient.put('/user/profile', profileData)
    return response
  },

  /**
   * Add body data measurement
   * @param {Object} bodyData - Body measurement data
   * @param {number} [bodyData.age] - Age
   * @param {string} [bodyData.gender] - Gender (male/female)
   * @param {number} bodyData.height - Height in cm
   * @param {number} bodyData.weight - Weight in kg
   * @param {number} [bodyData.body_fat_percentage] - Body fat percentage
   * @param {number} [bodyData.muscle_percentage] - Muscle percentage
   * @param {string} bodyData.measurement_date - Date of measurement (YYYY-MM-DD)
   * @returns {Promise<Object>} Created body data record
   * Requirements: 2.4
   */
  async addBodyData(bodyData) {
    const response = await apiClient.post('/user/body-data', bodyData)
    return response
  },

  /**
   * Fetch body data history
   * @param {Object} [params] - Query parameters
   * @param {string} [params.start_date] - Start date filter
   * @param {string} [params.end_date] - End date filter
   * @param {number} [params.limit] - Number of records to fetch
   * @returns {Promise<Object>} Body data history
   * Requirements: 2.5
   */
  async fetchBodyData(params = {}) {
    const response = await apiClient.get('/user/body-data', { params })
    return response
  },

  /**
   * Set fitness goals
   * @param {Object} goalsData - Fitness goals data
   * @param {string} goalsData.goal_type - Goal type (weight_loss/muscle_gain/endurance/general_fitness)
   * @param {number} [goalsData.target_weight] - Target weight in kg
   * @param {number} [goalsData.target_body_fat] - Target body fat percentage
   * @param {string} [goalsData.target_date] - Target date to achieve goal
   * @param {string} [goalsData.notes] - Additional notes
   * @returns {Promise<Object>} Created/updated goals
   * Requirements: 2.6
   */
  async setGoals(goalsData) {
    const response = await apiClient.post('/user/fitness-goals', goalsData)
    return response
  },

  /**
   * Fetch fitness goals
   * @returns {Promise<Object>} Current fitness goals
   * Requirements: 2.6
   */
  async fetchGoals() {
    const response = await apiClient.get('/user/fitness-goals')
    return response
  },

  /**
   * Update fitness goals
   * @param {Object} goalsData - Updated goals data
   * @returns {Promise<Object>} Updated goals
   * Requirements: 2.6
   */
  async updateGoals(goalsData) {
    const response = await apiClient.put('/user/fitness-goals', goalsData)
    return response
  }
}

export default userService
