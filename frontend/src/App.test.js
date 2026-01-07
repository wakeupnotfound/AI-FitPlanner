import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'

// Helper to create a fresh router for each test
const createTestRouter = () => {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', name: 'dashboard', component: { template: '<div>Dashboard</div>' } },
      { path: '/login', name: 'login', component: { template: '<div>Login</div>' } },
      { path: '/register', name: 'register', component: { template: '<div>Register</div>' } }
    ]
  })
}

describe('App.vue', () => {
  it('renders without crashing', async () => {
    const router = createTestRouter()
    router.push('/')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router],
        stubs: {
          NavigationBar: true
        }
      }
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('shows navigation bar on dashboard route', async () => {
    const router = createTestRouter()
    router.push('/')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router],
        stubs: {
          NavigationBar: true
        }
      }
    })

    expect(wrapper.findComponent({ name: 'NavigationBar' }).exists()).toBe(true)
  })

  it('hides navigation bar on login route', async () => {
    const router = createTestRouter()
    router.push('/login')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router],
        stubs: {
          NavigationBar: true
        }
      }
    })

    expect(wrapper.findComponent({ name: 'NavigationBar' }).exists()).toBe(false)
  })
})
