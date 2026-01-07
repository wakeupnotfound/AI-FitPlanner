import { defineStore } from 'pinia'
import apiClient from '../services/api'

/**
 * User Store
 * Manages user profile, body data, and fitness goals
 */
export const useUserStore = defineStore('user', {
  state: () => ({
    profile: null,
    bodyData: [],
    goals: null,
    loading: false,
    error: null
  }),

  getters: {
    /**
     * Get user profile
     */
    userProfile: (state) => state.profile,

    /**
     * Get latest body data
     */
    latestBodyData: (state) => {
      if (!state.bodyData || state.bodyData.length === 0) {
        return null
      }
      return state.bodyData[0]
    },

    /**
     * Get body data history sorted by date (newest first)
     */
    bodyDataHistory: (state) => {
      return [...state.bodyData].sort((a, b) => 
        new Date(b.record_date) - new Date(a.record_date)
      )
    },

    /**
     * Get fitness goals
     */
    fitnessGoals: (state) => state.goals
  },

  actions: {
    /**
     * Fetch user profile
     */
    async fetchProfile() {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/user/profile')
        this.profile = response.data
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Update user profile
     * @param {Object} profileData - Profile data to update
     */
    async updateProfile(profileData) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.put('/user/profile', profileData)
        this.profile = response.data
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Add body data measurement
     * @param {Object} bodyDataEntry - Body data entry
     * @param {string} bodyDataEntry.record_date - Date of measurement
     * @param {number} bodyDataEntry.weight - Weight in kg
     * @param {number} bodyDataEntry.height - Height in cm
     * @param {number} bodyDataEntry.body_fat - Body fat percentage
     * @param {number} bodyDataEntry.muscle_mass - Muscle mass in kg
     */
    async addBodyData(bodyDataEntry) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.post('/user/body-data', bodyDataEntry)
        
        // Add new body data to the beginning of the array
        this.bodyData.unshift(response.data)
        
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch body data history
     */
    async fetchBodyData() {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/user/body-data')
        this.bodyData = response.data || []
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Set fitness goals
     * @param {Object} goalsData - Fitness goals data
     */
    async setGoals(goalsData) {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.post('/user/fitness-goals', goalsData)
        this.goals = response.data
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch fitness goals
     */
    async fetchGoals() {
      this.loading = true
      this.error = null

      try {
        const response = await apiClient.get('/user/fitness-goals')
        this.goals = response.data
        return response
      } catch (error) {
        this.error = error
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Clear user data
     */
    clearUserData() {
      this.profile = null
      this.bodyData = []
      this.goals = null
      this.loading = false
      this.error = null
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'user',
        storage: localStorage,
        paths: ['profile', 'bodyData', 'goals']
      }
    ]
  }
})
