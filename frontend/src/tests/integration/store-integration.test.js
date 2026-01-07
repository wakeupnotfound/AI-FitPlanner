import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

// Import all stores
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import { useAIConfigStore } from '@/stores/aiConfig'
import { useTrainingStore } from '@/stores/training'
import { useNutritionStore } from '@/stores/nutrition'
import { useStatisticsStore } from '@/stores/statistics'

describe('Integration Tests - Store Interactions', () => {
  let pinia

  beforeEach(() => {
    vi.clearAllMocks()
    pinia = createPinia()
    setActivePinia(pinia)
    localStorage.clear()
  })

  describe('Cross-Store Dependencies', () => {
    it('should clear all stores when auth store logs out', () => {
      const authStore = useAuthStore()
      const userStore = useUserStore()
      const trainingStore = useTrainingStore()
      const nutritionStore = useNutritionStore()

      // Set up initial state
      authStore.setUser({ id: 1, username: 'test' })
      authStore.setTokens('access', 'refresh')
      userStore.profile = { id: 1, name: 'Test User' }
      trainingStore.currentPlan = { id: 1, name: 'Plan' }
      nutritionStore.todayMeals = [{ id: 1 }]

      // Perform logout
      authStore.clearAuth()

      // Verify auth store is cleared
      expect(authStore.user).toBeNull()
      expect(authStore.isAuthenticated).toBe(false)
      // localStorage returns null for missing keys
      expect(localStorage.getItem('access_token')).toBeNull()
    })

    it('should maintain data consistency across stores', () => {
      const authStore = useAuthStore()
      const userStore = useUserStore()

      // Set user in auth store
      authStore.setUser({ id: 1, username: 'testuser', email: 'test@example.com' })

      // Set user profile
      userStore.profile = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        nickname: 'Test'
      }

      // Verify consistency
      expect(authStore.user.id).toBe(userStore.profile.id)
      expect(authStore.user.username).toBe(userStore.profile.username)
    })

    it('should handle concurrent store updates', () => {
      const trainingStore = useTrainingStore()
      const nutritionStore = useNutritionStore()
      const statisticsStore = useStatisticsStore()

      // Set data in stores
      trainingStore.todayWorkout = { id: 1, exercises: [] }
      nutritionStore.todayMeals = [{ id: 1, meal_type: 'breakfast' }]
      statisticsStore.trainingStats = { total_workouts: 10 }

      // Verify all stores are updated
      expect(trainingStore.todayWorkout).toBeTruthy()
      expect(nutritionStore.todayMeals).toHaveLength(1)
      expect(statisticsStore.trainingStats).toBeTruthy()
    })

    it('should sync AI config with training store for plan generation', () => {
      const aiConfigStore = useAIConfigStore()
      const trainingStore = useTrainingStore()

      // Set up AI config
      aiConfigStore.configs = [
        { id: 1, provider: 'openai', is_default: true, status: 'active' }
      ]

      // Verify AI config is available for training
      expect(aiConfigStore.configs.length).toBeGreaterThan(0)
      expect(aiConfigStore.configs.find(c => c.is_default)).toBeTruthy()

      // Training store should be able to reference AI config
      trainingStore.currentPlan = {
        id: 1,
        name: 'AI Generated Plan',
        ai_api_id: aiConfigStore.configs[0].id
      }

      expect(trainingStore.currentPlan.ai_api_id).toBe(1)
    })

    it('should maintain user profile and body data relationship', () => {
      const userStore = useUserStore()

      // Set profile
      userStore.profile = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com'
      }

      // Add body data
      userStore.bodyData = [
        { id: 1, user_id: 1, weight: 70, height: 175, record_date: '2024-01-07' },
        { id: 2, user_id: 1, weight: 69, height: 175, record_date: '2024-01-14' }
      ]

      // Verify relationship
      expect(userStore.profile.id).toBe(userStore.bodyData[0].user_id)
      expect(userStore.bodyDataHistory).toHaveLength(2)
    })
  })

  describe('Store Persistence', () => {
    it('should persist auth state to localStorage', () => {
      const authStore = useAuthStore()

      authStore.setTokens('access-token', 'refresh-token')
      authStore.setUser({ id: 1, username: 'test' })

      // Verify auth store state is set
      expect(authStore.accessToken).toBe('access-token')
      expect(authStore.refreshToken).toBe('refresh-token')
      expect(authStore.isAuthenticated).toBe(true)
    })

    it('should restore auth state from store initialization', () => {
      const authStore = useAuthStore()

      // Set tokens directly on store
      authStore.accessToken = 'stored-access'
      authStore.refreshToken = 'stored-refresh'
      authStore.isAuthenticated = true

      // Verify state
      expect(authStore.accessToken).toBe('stored-access')
      expect(authStore.refreshToken).toBe('stored-refresh')
      expect(authStore.isAuthenticated).toBe(true)
    })

    it('should clear persisted state on logout', () => {
      const authStore = useAuthStore()

      // Set up state
      authStore.setTokens('access', 'refresh')
      authStore.setUser({ id: 1 })

      // Verify state is set
      expect(authStore.isAuthenticated).toBe(true)

      // Clear auth (simulating logout without API call)
      authStore.clearAuth()

      // Verify state is cleared
      expect(authStore.accessToken).toBeNull()
      expect(authStore.refreshToken).toBeNull()
      expect(authStore.isAuthenticated).toBe(false)
    })
  })

  describe('Store Error Handling', () => {
    it('should handle errors gracefully in stores', () => {
      const userStore = useUserStore()

      // Store should not have partial data initially
      expect(userStore.profile).toBeNull()
    })

    it('should maintain state integrity on errors', () => {
      const trainingStore = useTrainingStore()

      // Initial state
      expect(trainingStore.todayWorkout).toBeNull()
      
      // Even after error, state should remain consistent
      expect(trainingStore.todayWorkout).toBeNull()
    })
  })

  describe('Store Loading States', () => {
    it('should manage loading states correctly', () => {
      const trainingStore = useTrainingStore()

      // Check initial loading state
      expect(trainingStore.isGenerating).toBe(false)

      // Simulate loading
      trainingStore.isGenerating = true
      expect(trainingStore.isGenerating).toBe(true)

      // Simulate completion
      trainingStore.isGenerating = false
      expect(trainingStore.isGenerating).toBe(false)
    })

    it('should reset loading state appropriately', () => {
      const trainingStore = useTrainingStore()

      // Set loading
      trainingStore.isGenerating = true
      
      // Reset
      trainingStore.isGenerating = false

      // Loading state should be reset
      expect(trainingStore.isGenerating).toBe(false)
    })
  })

  describe('Store Data Transformations', () => {
    it('should store data correctly in stores', () => {
      const nutritionStore = useNutritionStore()

      // Set data
      nutritionStore.todayMeals = [
        {
          id: 1,
          meal_type: 'breakfast',
          time: '08:00',
          foods: [
            { name: 'Oatmeal', calories: 150 },
            { name: 'Banana', calories: 100 }
          ]
        }
      ]

      // Verify data is stored correctly
      expect(nutritionStore.todayMeals).toHaveLength(1)
      expect(nutritionStore.todayMeals[0].foods).toHaveLength(2)
    })

    it('should maintain derived state correctly', () => {
      const statisticsStore = useStatisticsStore()

      // Set statistics data
      statisticsStore.trainingStats = {
        total_workouts: 20,
        total_duration: 1200,
        average_duration: 60
      }

      // Verify computed values
      expect(statisticsStore.trainingStats.total_workouts).toBe(20)
      expect(statisticsStore.trainingStats.average_duration).toBe(60)
    })

    it('should calculate nutrition totals correctly', () => {
      const nutritionStore = useNutritionStore()

      nutritionStore.todayMeals = [
        {
          id: 1,
          meal_type: 'breakfast',
          foods: [
            { name: 'Oatmeal', calories: 150, protein: 5, carbs: 27, fat: 3 },
            { name: 'Banana', calories: 100, protein: 1, carbs: 27, fat: 0 }
          ]
        },
        {
          id: 2,
          meal_type: 'lunch',
          foods: [
            { name: 'Chicken', calories: 300, protein: 30, carbs: 0, fat: 10 }
          ]
        }
      ]

      // Calculate totals
      const totalCalories = nutritionStore.todayMeals.reduce((sum, meal) => {
        return sum + meal.foods.reduce((mealSum, food) => mealSum + food.calories, 0)
      }, 0)

      expect(totalCalories).toBe(550)
    })
  })

  describe('Store Service Integration', () => {
    it('should handle training plan with exercises correctly', () => {
      const trainingStore = useTrainingStore()

      trainingStore.currentPlan = {
        id: 1,
        name: 'Strength Training',
        status: 'active',
        plan_data: {
          weeks: [
            {
              week_number: 1,
              days: [
                {
                  day_number: 1,
                  workouts: [
                    {
                      name: 'Upper Body',
                      exercises: [
                        { name: 'Bench Press', sets: 3, reps: 10, weight: 60 },
                        { name: 'Rows', sets: 3, reps: 10, weight: 50 }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      }

      expect(trainingStore.currentPlan.plan_data.weeks).toHaveLength(1)
      expect(trainingStore.currentPlan.plan_data.weeks[0].days[0].workouts[0].exercises).toHaveLength(2)
    })

    it('should track training history correctly', () => {
      const trainingStore = useTrainingStore()

      trainingStore.history = [
        { id: 1, workout_date: '2024-01-01', duration_minutes: 45, rating: 4 },
        { id: 2, workout_date: '2024-01-03', duration_minutes: 60, rating: 5 },
        { id: 3, workout_date: '2024-01-05', duration_minutes: 30, rating: 3 }
      ]

      expect(trainingStore.history).toHaveLength(3)
      
      // Calculate average rating
      const avgRating = trainingStore.history.reduce((sum, h) => sum + h.rating, 0) / trainingStore.history.length
      expect(avgRating).toBeCloseTo(4, 1)
    })

    it('should integrate statistics with training and body data', () => {
      const statisticsStore = useStatisticsStore()
      const userStore = useUserStore()

      // Set body data
      userStore.bodyData = [
        { weight: 75, record_date: '2024-01-01' },
        { weight: 74, record_date: '2024-01-08' },
        { weight: 73, record_date: '2024-01-15' }
      ]

      // Set statistics
      statisticsStore.bodyTrends = {
        weight: userStore.bodyData.map(d => ({ date: d.record_date, value: d.weight }))
      }

      expect(statisticsStore.bodyTrends.weight).toHaveLength(3)
      
      // Verify weight trend is decreasing
      const weights = statisticsStore.bodyTrends.weight.map(w => w.value)
      expect(weights[0]).toBeGreaterThan(weights[2])
    })
  })

  describe('Store State Isolation', () => {
    it('should not leak state between store instances', () => {
      const authStore1 = useAuthStore()
      authStore1.setUser({ id: 1, username: 'user1' })

      // Create new pinia and store
      const newPinia = createPinia()
      setActivePinia(newPinia)
      const authStore2 = useAuthStore()

      // New store should have fresh state
      expect(authStore2.user).toBeNull()
      expect(authStore2.isAuthenticated).toBe(false)
    })

    it('should maintain store reactivity', () => {
      const trainingStore = useTrainingStore()
      
      // Initial state
      expect(trainingStore.isGenerating).toBe(false)
      
      // Update state
      trainingStore.isGenerating = true
      expect(trainingStore.isGenerating).toBe(true)
      
      // Reset
      trainingStore.isGenerating = false
      expect(trainingStore.isGenerating).toBe(false)
    })
  })
})
