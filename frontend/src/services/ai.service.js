import apiClient from './api'

/**
 * AI Configuration Service
 * Handles all AI API configuration operations
 */
export const aiService = {
  /**
   * Fetch all AI configurations for the current user
   * @returns {Promise<Object>} Response with list of AI configs
   */
  async fetchConfigs() {
    return apiClient.get('/ai-apis')
  },

  /**
   * Add a new AI API configuration
   * @param {Object} configData - AI configuration data
   * @param {string} configData.provider - AI provider (openai, wenxin, tongyi)
   * @param {string} configData.name - Configuration name
   * @param {string} configData.api_endpoint - API endpoint URL
   * @param {string} configData.api_key - API key
   * @param {string} configData.model - Model name
   * @param {number} [configData.max_tokens] - Maximum tokens
   * @param {number} [configData.temperature] - Temperature setting
   * @param {boolean} [configData.is_default] - Set as default
   * @returns {Promise<Object>} Response with created config
   */
  async addConfig(configData) {
    return apiClient.post('/ai-apis', configData)
  },

  /**
   * Update an existing AI API configuration
   * @param {number} configId - Configuration ID
   * @param {Object} configData - Updated configuration data
   * @returns {Promise<Object>} Response with updated config
   */
  async updateConfig(configId, configData) {
    return apiClient.put(`/ai-apis/${configId}`, configData)
  },

  /**
   * Test AI API connection
   * @param {number} configId - Configuration ID to test
   * @returns {Promise<Object>} Response with test results
   */
  async testConnection(configId) {
    return apiClient.post(`/ai-apis/${configId}/test`)
  },

  /**
   * Set an AI configuration as default
   * @param {number} configId - Configuration ID to set as default
   * @returns {Promise<Object>} Response confirming default set
   */
  async setDefault(configId) {
    return apiClient.put(`/ai-apis/${configId}/default`)
  },

  /**
   * Delete an AI API configuration
   * @param {number} configId - Configuration ID to delete
   * @returns {Promise<Object>} Response confirming deletion
   */
  async deleteConfig(configId) {
    return apiClient.delete(`/ai-apis/${configId}`)
  },

  /**
   * Get a single AI configuration by ID
   * @param {number} configId - Configuration ID
   * @returns {Promise<Object>} Response with config details
   */
  async getConfig(configId) {
    return apiClient.get(`/ai-apis/${configId}`)
  }
}

export default aiService
