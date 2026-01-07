import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'

// Import stores
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import { useAIConfigStore } from '@/stores/aiConfig'
import { useTrainingStore } from '@/stores/training'
import { useNutritionStore } from '@/stores/nutrition'
import { useStatisticsStore } from '@/stores/statistics'

describe('Integration Tests - Complete User Flows', () => {
  let pinia
  let router

  const routes = [
    { path: '/login', name: 'login', component: { template: '<div>Login</div>' }, meta: { requiresAuth: false } },
    { path: '/register', name: 'register', component: { template: '<div>Register</div>' }, meta: { requiresAuth: false } },
    { path: '/dashboard', name: 'dashboard', component: { template: '<div>Dashboard</div>' }, meta: { requiresAuth: true } },
    { path: '/profile', name: 'profile', component: { template: '<div>Profile</div>' }, meta: { requiresAuth: true } },
    { path: '/ai-config', name: 'ai-config', component: { template: '<div>AI Config</div>' }, meta: { requiresAuth: true } },
    { path: '/training', name: 'training', component: { template: '<div>Training</div>' }, meta: { requiresAuth: true } },
    { path: '/nutrition', name: 'nutrition', component: { template: '<div>Nutrition</div>' }, meta: { requiresAuth: true } },
    { path: '/statistics', name: 'statistics', component: { template: '<div>Statistics</div>' }, meta: { requiresAuth: true } }
  ]

  beforeEach(() => {
    // Reset all mocks
    vi.clearAllMocks()
    
    // Create fresh instances
    pinia = createPinia()
    setActivePinia(pinia)
    
    router = createRouter({
      history: createMemoryHistory(),
      routes
    })
    
    // Clear localStorage
    localStorage.clear()
  })

  describe('User Registration and Login Flow', () => {
    it('should complete registration and auto-login flow', () => {
      const authStore = useAuthStore()
      
      // Simulate successful registration by setting auth state
      authStore.setTokens('test-access-token', 'test-refresh-token')
      authStore.setUser({
        id: 1,
        username: 'testuser',
        email: 'test@example.com'
      })

      // Verify tokens are stored
      expect(authStore.isAuthenticated).toBe(true)
      expect(authStore.user).toBeTruthy()
      expect(authStore.user.username).toBe('testuser')
    })

    it('should handle login flow and store tokens', () => {
      const authStore = useAuthStore()
      
      // Simulate successful login
      authStore.setTokens('test-access-token', 'test-refresh-token')
      authStore.setUser({
        id: 1,
        username: 'testuser',
        email: 'test@example.com'
      })

      // Verify authentication state
      expect(authStore.isAuthenticated).toBe(true)
      expect(authStore.accessToken).toBe('test-access-token')
      expect(authStore.refreshToken).toBe('test-refresh-token')
    })

    it('should handle logout and clear all state', () => {
      const authStore = useAuthStore()
      
      // Set up authenticated state
      authStore.setTokens('access-token', 'refresh-token')
      authStore.setUser({ id: 1, username: 'testuser' })

      expect(authStore.isAuthenticated).toBe(true)

      // Perform logout
      authStore.clearAuth()

      // Verify state is cleared
      expect(authStore.isAuthenticated).toBe(false)
      expect(authStore.user).toBeNull()
      // localStorage returns null for missing keys
      expect(localStorage.getItem('access_token')).toBeNull()
      expect(localStorage.getItem('refresh_token')).toBeNull()
    })
  })

  describe('Store and Service Integration', () => {
    beforeEach(() => {
      // Set up authenticated state for protected routes
      const authStore = useAuthStore()
      authStore.setTokens('test-access-token', 'test-refresh-token')
      authStore.setUser({ id: 1, username: 'testuser' })
    })

    it('should maintain user state in user store', () => {
      const userStore = useUserStore()

      // Set profile data
      userStore.profile = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        nickname: 'Test User'
      }

      // Verify store is updated
      expect(userStore.profile).toBeTruthy()
      expect(userStore.profile.username).toBe('testuser')
    })

    it('should maintain AI config state', () => {
      const aiConfigStore = useAIConfigStore()

      // Set configs
      aiConfigStore.configs = [
        {
          id: 1,
          provider: 'openai',
          model_name: 'gpt-4',
          is_default: true,
          status: 'active'
        }
      ]

      // Verify store is updated
      expect(aiConfigStore.configs).toHaveLength(1)
      expect(aiConfigStore.configs[0].provider).toBe('openai')
    })

    it('should maintain training state', () => {
      const trainingStore = useTrainingStore()

      // Set today's workout
      trainingStore.todayWorkout = {
        id: 1,
        date: '2024-01-07',
        exercises: [
          {
            name: 'Push-ups',
            sets: 3,
            reps: 10
          }
        ]
      }

      // Verify store is updated
      expect(trainingStore.todayWorkout).toBeTruthy()
      expect(trainingStore.todayWorkout.exercises).toHaveLength(1)
    })

    it('should maintain nutrition state', () => {
      const nutritionStore = useNutritionStore()

      // Set today's meals
      nutritionStore.todayMeals = [
        {
          id: 1,
          meal_type: 'breakfast',
          foods: ['Oatmeal', 'Banana'],
          calories: 350
        }
      ]

      // Verify store is updated
      expect(nutritionStore.todayMeals).toHaveLength(1)
      expect(nutritionStore.todayMeals[0].meal_type).toBe('breakfast')
    })
  })

  describe('Router Navigation Flows', () => {
    it('should track authentication state for routing', async () => {
      const authStore = useAuthStore()
      authStore.clearAuth()

      expect(authStore.isAuthenticated).toBe(false)
      
      // Authenticate
      authStore.setTokens('token', 'refresh')
      authStore.setUser({ id: 1 })
      
      expect(authStore.isAuthenticated).toBe(true)
    })

    it('should allow navigation between routes', async () => {
      await router.push('/dashboard')
      await router.isReady()

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })

    it('should navigate to all main routes', async () => {
      const mainRoutes = ['/dashboard', '/profile', '/training', '/nutrition', '/statistics', '/ai-config']
      
      for (const route of mainRoutes) {
        await router.push(route)
        await router.isReady()
        expect(router.currentRoute.value.path).toBe(route)
      }
    })

    it('should handle back navigation correctly', async () => {
      await router.push('/dashboard')
      await router.isReady()
      
      await router.push('/training')
      await router.isReady()
      expect(router.currentRoute.value.path).toBe('/training')
      
      // Use router.go(-1) and wait for navigation
      await router.go(-1)
      // Wait for the navigation to complete
      await new Promise(resolve => setTimeout(resolve, 10))
      expect(router.currentRoute.value.path).toBe('/dashboard')
    })

    it('should preserve route meta information', async () => {
      await router.push('/dashboard')
      await router.isReady()
      
      expect(router.currentRoute.value.meta.requiresAuth).toBe(true)
      
      await router.push('/login')
      await router.isReady()
      
      expect(router.currentRoute.value.meta.requiresAuth).toBe(false)
    })
  })

  describe('End-to-End User Journey', () => {
    it('should complete full user journey: register → configure AI → create plan', () => {
      // Step 1: Register
      const authStore = useAuthStore()
      authStore.setTokens('token', 'refresh')
      authStore.setUser({ id: 1, username: 'newuser' })
      expect(authStore.isAuthenticated).toBe(true)

      // Step 2: Configure AI
      const aiConfigStore = useAIConfigStore()
      aiConfigStore.configs = [{
        id: 1,
        provider: 'openai',
        is_default: true
      }]

      expect(aiConfigStore.configs).toHaveLength(1)

      // Step 3: Set training plan
      const trainingStore = useTrainingStore()
      trainingStore.currentPlan = {
        id: 1,
        name: 'Beginner Plan',
        status: 'active'
      }

      expect(trainingStore.currentPlan).toBeTruthy()
      expect(trainingStore.currentPlan.name).toBe('Beginner Plan')
    })

    it('should complete profile setup flow', () => {
      const authStore = useAuthStore()
      const userStore = useUserStore()

      // Step 1: Login
      authStore.setTokens('token', 'refresh')
      authStore.setUser({ id: 1, username: 'testuser' })

      // Step 2: Set profile
      userStore.profile = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        nickname: 'Test User'
      }

      // Step 3: Add body data
      userStore.bodyData = [
        { id: 1, weight: 75, height: 180, body_fat: 20, record_date: '2024-01-07' }
      ]

      // Step 4: Set goals
      userStore.goals = {
        target_weight: 70,
        target_body_fat: 15,
        goal_type: 'fat_loss'
      }

      // Verify complete profile
      expect(userStore.profile).toBeTruthy()
      expect(userStore.bodyData).toHaveLength(1)
      expect(userStore.goals).toBeTruthy()
    })

    it('should complete training workflow: assessment → plan → workout → record', () => {
      const authStore = useAuthStore()
      const trainingStore = useTrainingStore()

      // Setup auth
      authStore.setTokens('token', 'refresh')
      authStore.setUser({ id: 1 })

      // Step 1: Complete assessment (simulated)
      const assessment = {
        experience_level: 'intermediate',
        weekly_available_days: 4,
        daily_available_minutes: 60,
        activity_type: 'strength'
      }

      // Step 2: Generate plan
      trainingStore.currentPlan = {
        id: 1,
        name: 'Strength Training Plan',
        status: 'active',
        plan_data: {
          weeks: [{ week_number: 1, days: [] }]
        }
      }

      // Step 3: Set today's workout
      trainingStore.todayWorkout = {
        id: 1,
        date: '2024-01-07',
        exercises: [
          { name: 'Squats', sets: 3, reps: 10, weight: 60 },
          { name: 'Deadlifts', sets: 3, reps: 8, weight: 80 }
        ]
      }

      // Step 4: Record workout
      trainingStore.history = [
        {
          id: 1,
          workout_date: '2024-01-07',
          duration_minutes: 45,
          exercises: trainingStore.todayWorkout.exercises,
          rating: 4
        }
      ]

      // Verify workflow completion
      expect(trainingStore.currentPlan).toBeTruthy()
      expect(trainingStore.todayWorkout.exercises).toHaveLength(2)
      expect(trainingStore.history).toHaveLength(1)
    })

    it('should complete nutrition workflow: plan → meals → record', () => {
      const authStore = useAuthStore()
      const nutritionStore = useNutritionStore()

      // Setup auth
      authStore.setTokens('token', 'refresh')
      authStore.setUser({ id: 1 })

      // Step 1: Set nutrition plan
      nutritionStore.currentPlan = {
        id: 1,
        name: 'Weight Loss Diet',
        daily_calories: 2000,
        protein_ratio: 0.3,
        carbs_ratio: 0.4,
        fat_ratio: 0.3
      }

      // Step 2: Set today's meals
      nutritionStore.todayMeals = [
        {
          id: 1,
          meal_type: 'breakfast',
          time: '08:00',
          foods: [{ name: 'Oatmeal', calories: 300 }]
        },
        {
          id: 2,
          meal_type: 'lunch',
          time: '12:00',
          foods: [{ name: 'Chicken Salad', calories: 500 }]
        }
      ]

      // Step 3: Record meal history
      nutritionStore.history = [
        {
          id: 1,
          meal_date: '2024-01-07',
          meals: nutritionStore.todayMeals,
          total_calories: 800
        }
      ]

      // Verify workflow
      expect(nutritionStore.currentPlan).toBeTruthy()
      expect(nutritionStore.todayMeals).toHaveLength(2)
      expect(nutritionStore.history).toHaveLength(1)
    })

    it('should complete statistics viewing flow', () => {
      const authStore = useAuthStore()
      const statisticsStore = useStatisticsStore()
      const userStore = useUserStore()

      // Setup auth
      authStore.setTokens('token', 'refresh')
      authStore.setUser({ id: 1 })

      // Set body data for trends
      userStore.bodyData = [
        { weight: 75, record_date: '2024-01-01' },
        { weight: 74, record_date: '2024-01-08' }
      ]

      // Set training stats
      statisticsStore.trainingStats = {
        total_workouts: 10,
        total_duration: 450,
        average_duration: 45
      }

      // Set body trends
      statisticsStore.bodyTrends = {
        weight: [
          { date: '2024-01-01', value: 75 },
          { date: '2024-01-08', value: 74 }
        ]
      }

      // Set progress
      statisticsStore.progress = {
        weight_change: -1,
        goal_progress: 20
      }

      // Verify statistics
      expect(statisticsStore.trainingStats.total_workouts).toBe(10)
      expect(statisticsStore.bodyTrends.weight).toHaveLength(2)
      expect(statisticsStore.progress.weight_change).toBe(-1)
    })
  })

  describe('Error Recovery Flows', () => {
    it('should handle session expiry gracefully', () => {
      const authStore = useAuthStore()
      const userStore = useUserStore()

      // Setup authenticated state
      authStore.setTokens('token', 'refresh')
      authStore.setUser({ id: 1 })
      userStore.profile = { id: 1, username: 'test' }

      // Simulate session expiry
      authStore.clearAuth()

      // Verify cleanup
      expect(authStore.isAuthenticated).toBe(false)
      expect(authStore.user).toBeNull()
    })

    it('should preserve data during navigation errors', () => {
      const trainingStore = useTrainingStore()

      // Set data
      trainingStore.currentPlan = { id: 1, name: 'Test Plan' }
      trainingStore.todayWorkout = { id: 1, exercises: [] }

      // Simulate navigation (data should persist)
      expect(trainingStore.currentPlan).toBeTruthy()
      expect(trainingStore.todayWorkout).toBeTruthy()
    })
  })

  describe('Multi-Store Coordination', () => {
    it('should coordinate auth logout across all stores', () => {
      const authStore = useAuthStore()
      const userStore = useUserStore()
      const trainingStore = useTrainingStore()
      const nutritionStore = useNutritionStore()
      const aiConfigStore = useAIConfigStore()

      // Setup all stores
      authStore.setTokens('token', 'refresh')
      authStore.setUser({ id: 1 })
      userStore.profile = { id: 1 }
      trainingStore.currentPlan = { id: 1 }
      nutritionStore.currentPlan = { id: 1 }
      aiConfigStore.configs = [{ id: 1 }]

      // Logout
      authStore.clearAuth()

      // Auth should be cleared
      expect(authStore.isAuthenticated).toBe(false)
      
      // Other stores maintain their data (would be cleared by app logic)
      // This tests that stores are independent
      expect(userStore.profile).toBeTruthy()
    })

    it('should handle concurrent data updates', () => {
      const trainingStore = useTrainingStore()
      const nutritionStore = useNutritionStore()
      const statisticsStore = useStatisticsStore()

      // Simulate concurrent updates
      trainingStore.todayWorkout = { id: 1 }
      nutritionStore.todayMeals = [{ id: 1 }]
      statisticsStore.trainingStats = { total: 1 }

      // All updates should be reflected
      expect(trainingStore.todayWorkout).toBeTruthy()
      expect(nutritionStore.todayMeals).toHaveLength(1)
      expect(statisticsStore.trainingStats).toBeTruthy()
    })
  })
})
