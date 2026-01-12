import { vi } from 'vitest'
import { config } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { createRouter, createWebHistory } from 'vue-router'

// Import Vant UI components for manual registration
import { 
  NoticeBar,
  Icon,
  Tag,
  NavBar,
  Button,
  Form,
  Field,
  Cell,
  CellGroup,
  Tabbar,
  TabbarItem,
  Swipe,
  SwipeItem,
  Card,
  Grid,
  GridItem,
  List,
  PullRefresh,
  Loading,
  Empty,
  Image as VanImage,
  Dialog,
  Toast,
  Popup,
  ActionSheet,
  DatetimePicker,
  Picker,
  Progress,
  Step,
  Steps,
  Divider,
  Collapse,
  CollapseItem,
  Switch,
  Overlay,
  Skeleton
} from 'vant'

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

// Create mock router for tests
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/dashboard', name: 'Dashboard', component: { template: '<div>Dashboard</div>' } },
    { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
    { path: '/register', name: 'Register', component: { template: '<div>Register</div>' } },
    { path: '/profile', name: 'Profile', component: { template: '<div>Profile</div>' } },
    { path: '/training', name: 'Training', component: { template: '<div>Training</div>' } },
    { path: '/nutrition', name: 'Nutrition', component: { template: '<div>Nutrition</div>' } },
    { path: '/statistics', name: 'Statistics', component: { template: '<div>Statistics</div>' } },
    { path: '/ai-config', name: 'AIConfig', component: { template: '<div>AI Config</div>' } }
  ]
})

// Setup i18n for tests with missing translation keys
const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      nav: {
        dashboard: 'Dashboard',
        training: 'Training',
        nutrition: 'Nutrition',
        profile: 'Profile',
        statistics: 'Statistics'
      },
      common: {
        loading: 'Loading...',
        error: 'Error',
        success: 'Success',
        confirm: 'Confirm',
        cancel: 'Cancel'
      },
      app: {
        confirm: 'Confirm',
        cancel: 'Cancel',
        save: 'Save',
        add: 'Add',
        loading: 'Loading...'
      },
      dashboard: {
        goalProgress: 'Goal Progress'
      },
      statistics: {
        insufficientData: 'Insufficient Data'
      },
      auth: {
        loginTitle: 'Login',
        loginSubtitle: 'Welcome back',
        noAccount: 'No account?',
        register: 'Register',
        username: 'Username',
        password: 'Password',
        login: 'Login',
        usernameRequired: 'Username is required',
        passwordRequired: 'Password is required'
      },
      ai: {
        addConfig: 'Add AI Configuration',
        provider: 'Provider',
        name: 'Name',
        apiKey: 'API Key',
        apiKeyPlaceholder: 'Enter your API key',
        apiKeyHint: 'Your API key will be encrypted and stored securely',
        apiEndpoint: 'API Endpoint',
        model: 'Model',
        maxTokens: 'Max Tokens',
        temperature: 'Temperature',
        isDefault: 'Set as Default',
        validation: {
          providerRequired: 'Provider is required',
          nameRequired: 'Name is required',
          apiKeyRequired: 'API key is required',
          modelRequired: 'Model is required'
        }
      },
      profile: {
        weight: 'Weight',
        height: 'Height',
        age: 'Age',
        gender: 'Gender',
        bodyFat: 'Body Fat',
        muscleMass: 'Muscle Mass',
        kg: 'kg',
        cm: 'cm',
        percent: '%',
        measurementDate: 'Measurement Date',
        addBodyData: 'Add Body Data'
      }
    }
  }
})

// Configure Vue Test Utils with all required mocks
config.global.plugins = [i18n]

// Manually register Vant components globally
config.global.components = {
  'van-notice-bar': NoticeBar,
  'van-icon': Icon,
  'van-tag': Tag,
  'van-nav-bar': NavBar,
  'van-button': Button,
  'van-form': Form,
  'van-field': Field,
  'van-cell': Cell,
  'van-cell-group': CellGroup,
  'van-tabbar': Tabbar,
  'van-tabbar-item': TabbarItem,
  'van-swipe': Swipe,
  'van-swipe-item': SwipeItem,
  'van-card': Card,
  'van-grid': Grid,
  'van-grid-item': GridItem,
  'van-list': List,
  'van-pull-refresh': PullRefresh,
  'van-loading': Loading,
  'van-empty': Empty,
  'van-image': VanImage,
  'van-dialog': Dialog,
  'van-toast': Toast,
  'van-popup': Popup,
  'van-action-sheet': ActionSheet,
  'van-datetime-picker': DatetimePicker,
  'van-picker': Picker,
  'van-progress': Progress,
  'van-step': Step,
  'van-steps': Steps,
  'van-divider': Divider,
  'van-collapse': Collapse,
  'van-collapse-item': CollapseItem,
  'van-switch': Switch,
  'van-overlay': Overlay,
  'van-skeleton': Skeleton
}

config.global.mocks = {
  $router: router
}
config.global.provide = {
  router
}
