import { ref, onMounted, onUnmounted } from 'vue'
import { offlineQueue } from '@/utils/offlineQueue'

/**
 * Composable for tracking network status
 */
export function useNetworkStatus() {
  const isOnline = ref(navigator.onLine)
  const queueSize = ref(offlineQueue.size())
  const lastSyncTime = ref(null)
  const isSyncing = ref(false)

  /**
   * Handle online event
   */
  const handleOnline = () => {
    isOnline.value = true
    syncOfflineQueue()
  }

  /**
   * Handle offline event
   */
  const handleOffline = () => {
    isOnline.value = false
  }

  /**
   * Sync offline queue
   */
  const syncOfflineQueue = async () => {
    if (offlineQueue.size() === 0) {
      return
    }

    isSyncing.value = true
    try {
      await offlineQueue.syncQueue()
      lastSyncTime.value = new Date()
      queueSize.value = offlineQueue.size()
    } catch (error) {
      console.error('Failed to sync offline queue:', error)
    } finally {
      isSyncing.value = false
    }
  }

  /**
   * Update queue size
   */
  const updateQueueSize = () => {
    queueSize.value = offlineQueue.size()
  }

  /**
   * Register sync callback
   */
  const setupSyncCallback = () => {
    offlineQueue.onSyncComplete((results) => {
      queueSize.value = offlineQueue.size()
      lastSyncTime.value = new Date()
    })
  }

  onMounted(() => {
    // Add event listeners
    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)
    
    // Setup sync callback
    setupSyncCallback()
    
    // Update initial queue size
    updateQueueSize()
  })

  onUnmounted(() => {
    // Remove event listeners
    window.removeEventListener('online', handleOnline)
    window.removeEventListener('offline', handleOffline)
  })

  return {
    isOnline,
    queueSize,
    lastSyncTime,
    isSyncing,
    syncOfflineQueue,
    updateQueueSize
  }
}
