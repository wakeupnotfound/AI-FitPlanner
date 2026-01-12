import { defineStore } from 'pinia'
import apiClient from '../services/api'

/**
 * AI Config Store
 * Manages AI API configurations
 */
export const useAIConfigStore = defineStore('aiConfig', {
  state: () => ({
    configs: [],
    defaultConfig: null,
    loading: false,
    error: null,
    testingConnection: false
  }),

  getters: {
    /**
     * Get all AI configurations
     */
    allConfigs: (state) => state.configs,

    /**
     * Get default AI configuration
     */
    getDefaultConfig: (state) => {
      if (state.defaultConfig) {
        return state.defaultConfig
      }
      // Find default from configs array
      return state.configs.find(config => config.is_default) || null
    },

    /**
     * Check if any AI config exists
     */
    hasConfig: (state) => state.configs.length > 0,

    /**
     * Check if default config is set
     */
    hasDefaultConfig: (state) => {
      return !!state.defaultConfig || state.configs.some(config => config.is_default)
    },

    /**
     * Get active configurations
     */
    activeConfigs: (state) => {
      return state.configs.filter(config => config.status === 'active')
    }
  },

  actions: {
    /**
     * Fetch all AI configurations
     */
    async fetchConfigs() {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/ai-apis')
        // Ensure response has data property and it's an array
        this.configs = response.data?.apis || response.data || []
        
        // Update default config
        if (this.configs && this.configs.length > 0) {
          const defaultCfg = this.configs.find(config => config.is_default)
          if (defaultCfg) {
            this.defaultConfig = defaultCfg
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
     * Add new AI configuration
     * @param {Object} configData - AI configuration data
     * @param {string} configData.provider - AI provider (openai, wenxin, tongyi)
     * @param {string} configData.api_key - API key
     * @param {string} configData.model_name - Model name
     */
    async addConfig(configData) {
      this.loading = true
      this.error = null
      
      try {
        const response = await apiClient.post('/ai-apis', configData)
        
        // Validate response before using it
        if (response.data) {
          this.configs.push(response.data)
          
          // If this is the first config, set it as default
          if (this.configs.length === 1) {
            this.defaultConfig = response.data
          }
        }
        
        return response
      } catch (error) {
        this.error = error
        console.error('AI Config add error:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Test AI API connection
     * @param {number} configId - Configuration ID
     */
     async testConnection(configId) {
      if (!configId) {
        throw new Error('配置ID无效')
      }
      
      this.testingConnection = true
      this.error = null
      
      try {
        const response = await apiClient.post(`/ai-apis/${configId}/test`)
        
        // Update config status in local state
        const configIndex = this.configs.findIndex(c => c.id === configId)
        if (configIndex !== -1) {
          const testStatus = response.data.test_result?.status
          if (testStatus) {
            this.configs[configIndex].status = testStatus === 'success'
          } else if (typeof response.data.status !== 'undefined') {
            this.configs[configIndex].status = response.data.status
          }
          this.configs[configIndex].last_test_at = new Date().toISOString()
        }
        
        // Return full response data for component to handle
        return response
      } catch (error) {
        this.error = error
        console.error('AI Config test error:', error)
        throw error
      } finally {
        this.testingConnection = false
      }
    },

    /**
     * Set default AI configuration
     * @param {number} configId - Configuration ID
     */
    async setDefault(configId) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.post(`/ai-apis/${configId}/set-default`)
        
        // Update is_default flag for all configs
        this.configs = this.configs.map(config => ({
          ...config,
          is_default: config.id === configId
        }))
        
        // Update default config
        this.defaultConfig = this.configs.find(config => config.id === configId)
        
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Delete AI configuration
     * @param {number} configId - Configuration ID
     */
    async deleteConfig(configId) {
      this.loading = true
      this.error = null

      try {
        await apiClient.delete(`/ai-apis/${configId}`)
        
        // Remove config from local state
        const deletedConfig = this.configs.find(c => c.id === configId)
        this.configs = this.configs.filter(config => config.id !== configId)
        
        // If deleted config was default, clear default
        if (deletedConfig?.is_default) {
          this.defaultConfig = null
          
          // Set first remaining config as default if any exist
          if (this.configs && this.configs.length > 0) {
            await this.setDefault(this.configs[0].id)
          }
        }
        
        return { success: true }
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Update AI configuration
     * @param {number} configId - Configuration ID
     * @param {Object} configData - Updated configuration data
     */
    async updateConfig(configId, configData) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.put(`/ai-apis/${configId}`, configData)
        
        // Update config in local state
        const configIndex = this.configs.findIndex(c => c.id === configId)
        if (configIndex !== -1) {
          this.configs[configIndex] = response.data
        }
        
        // Update default config if it was updated
        if (this.defaultConfig?.id === configId) {
          this.defaultConfig = response.data
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
     * Clear AI config data
     */
    clearConfigData() {
      this.configs = []
      this.defaultConfig = null
      this.loading = false
      this.error = null
      this.testingConnection = false
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'aiConfig',
        storage: localStorage,
        paths: ['configs', 'defaultConfig']
      }
    ]
  }
})
