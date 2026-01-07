<template>
  <div id="app">
    <!-- Network Status Indicator -->
    <NetworkStatusIndicator />
    
    <div class="app-content" :class="{ 'has-tabbar': showNavigation }">
      <router-view />
    </div>
    <NavigationBar v-if="showNavigation" />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import NavigationBar from '@/components/common/NavigationBar.vue'
import NetworkStatusIndicator from '@/components/common/NetworkStatusIndicator.vue'

const route = useRoute()

// Show navigation bar only for authenticated routes
const showNavigation = computed(() => {
  const noNavRoutes = ['login', 'register']
  return !noNavRoutes.includes(route.name)
})

onMounted(() => {
  // Set viewport height for mobile browsers
  const setViewportHeight = () => {
    const vh = window.innerHeight * 0.01
    document.documentElement.style.setProperty('--vh', `${vh}px`)
  }

  setViewportHeight()
  window.addEventListener('resize', setViewportHeight)
})
</script>

<style>
#app {
  width: 100%;
  min-height: 100vh;
  min-height: calc(var(--vh, 1vh) * 100);
}

.app-content {
  width: 100%;
  min-height: 100vh;
  min-height: calc(var(--vh, 1vh) * 100);
}

/* Add padding for bottom tabbar when navigation is visible */
.app-content.has-tabbar {
  padding-bottom: 50px;
}
</style>
