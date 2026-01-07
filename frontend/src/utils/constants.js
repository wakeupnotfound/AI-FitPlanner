// API Configuration
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
export const API_TIMEOUT = 10000

// Storage Keys
export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'access_token',
  REFRESH_TOKEN: 'refresh_token',
  USER_INFO: 'user_info',
  LANGUAGE: 'language',
  THEME: 'theme'
}

// Route Names
export const ROUTE_NAMES = {
  LOGIN: 'login',
  REGISTER: 'register',
  DASHBOARD: 'dashboard',
  PROFILE: 'profile',
  AI_CONFIG: 'ai-config',
  TRAINING: 'training',
  NUTRITION: 'nutrition',
  STATISTICS: 'statistics'
}

// HTTP Status Codes
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  TOO_MANY_REQUESTS: 429,
  INTERNAL_SERVER_ERROR: 500
}

// AI Providers
export const AI_PROVIDERS = {
  OPENAI: 'openai',
  WENXIN: 'wenxin',
  TONGYI: 'tongyi'
}

// Meal Types
export const MEAL_TYPES = {
  BREAKFAST: 'breakfast',
  LUNCH: 'lunch',
  DINNER: 'dinner',
  SNACK: 'snack'
}

// Plan Status
export const PLAN_STATUS = {
  ACTIVE: 'active',
  COMPLETED: 'completed',
  PAUSED: 'paused'
}

// Touch Target Minimum Size (px)
export const MIN_TOUCH_TARGET_SIZE = 44

// Debounce Delay (ms)
export const DEBOUNCE_DELAY = 300

// Request Retry Configuration
export const RETRY_CONFIG = {
  MAX_RETRIES: 3,
  RETRY_DELAY: 1000
}
