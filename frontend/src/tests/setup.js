import { vi } from 'vitest'
import { config } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'

// Create a proper localStorage mock that stores values
const createStorageMock = () => {
  let store = {}
  return {
    getItem: vi.fn((key) => store[key] || null),
    setItem: vi.fn((key, value) => { store[key] = value }),
    removeItem: vi.fn((key) => { delete store[key] }),
    clear: vi.fn(() => { store = {} })
  }
}

global.localStorage = createStorageMock()

// Mock sessionStorage
global.sessionStorage = createStorageMock()

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  }))
})

// Mock IntersectionObserver
global.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  takeRecords() {
    return []
  }
  unobserve() {}
}

// Setup i18n for tests
const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      nav: {
        dashboard: 'Dashboard',
        training: 'Training',
        nutrition: 'Nutrition',
        profile: 'Profile'
      },
      common: {
        loading: 'Loading...',
        error: 'Error',
        success: 'Success'
      }
    }
  }
})

config.global.plugins = [i18n]
