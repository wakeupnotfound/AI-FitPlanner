<template>
  <van-tabbar 
    v-model="activeTab" 
    route 
    :safe-area-inset-bottom="true"
    class="navigation-bar"
  >
    <van-tabbar-item 
      to="/dashboard" 
      icon="home-o"
      name="dashboard"
      class="touch-target-nav"
    >
      {{ t('nav.dashboard') }}
    </van-tabbar-item>
    <van-tabbar-item 
      to="/training" 
      icon="fire-o"
      name="training"
      class="touch-target-nav"
    >
      {{ t('nav.training') }}
    </van-tabbar-item>
    <van-tabbar-item 
      to="/nutrition" 
      icon="coupon-o"
      name="nutrition"
      class="touch-target-nav"
    >
      {{ t('nav.nutrition') }}
    </van-tabbar-item>
    <van-tabbar-item 
      to="/statistics" 
      icon="chart-trending-o"
      name="statistics"
      class="touch-target-nav"
    >
      {{ t('nav.statistics') }}
    </van-tabbar-item>
    <van-tabbar-item 
      to="/profile" 
      icon="user-o"
      name="profile"
      class="touch-target-nav"
    >
      {{ t('nav.profile') }}
    </van-tabbar-item>
  </van-tabbar>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()

// Active tab state - synced with current route
const activeTab = ref(getActiveTabFromRoute(route.name))

// Watch route changes to update active tab highlighting
watch(
  () => route.name,
  (newRouteName) => {
    activeTab.value = getActiveTabFromRoute(newRouteName)
  }
)

/**
 * Get the active tab name from the current route
 * @param {string} routeName - Current route name
 * @returns {string} Tab name to highlight
 */
function getActiveTabFromRoute(routeName) {
  const tabRoutes = ['dashboard', 'training', 'nutrition', 'statistics', 'profile']
  return tabRoutes.includes(routeName) ? routeName : 'dashboard'
}
</script>

<style scoped>
.navigation-bar {
  --van-tabbar-item-active-color: var(--van-primary-color, #1989fa);
}

/* Ensure minimum touch target size of 44x44px */
:deep(.van-tabbar-item) {
  min-height: 50px;
  min-width: 44px;
}

/* Touch target class for nav items */
.touch-target-nav {
  position: relative;
}

/* Touch feedback animation */
:deep(.van-tabbar-item)::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  background: rgba(25, 137, 250, 0.1);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  transition: width 0.3s ease, height 0.3s ease, opacity 0.3s ease;
  opacity: 0;
  pointer-events: none;
}

:deep(.van-tabbar-item:active)::after {
  width: 100%;
  height: 100%;
  opacity: 1;
  transition: width 0s, height 0s, opacity 0s;
}

:deep(.van-tabbar-item__icon) {
  font-size: 22px;
  margin-bottom: 4px;
}

:deep(.van-tabbar-item__text) {
  font-size: 11px;
}

/* Active state enhancement */
:deep(.van-tabbar-item--active) {
  transform: scale(1);
  transition: transform 0.2s ease;
}

:deep(.van-tabbar-item--active .van-tabbar-item__icon) {
  transform: scale(1.1);
  transition: transform 0.2s ease;
}
</style>
