/**
 * Offline Queue Manager
 * Handles queueing of write operations when offline and syncing when connection is restored
 */

const QUEUE_STORAGE_KEY = 'offline_queue'
const MAX_QUEUE_SIZE = 100

class OfflineQueue {
  constructor() {
    this.queue = this.loadQueue()
    this.isSyncing = false
    this.syncCallbacks = []
    
    // Listen for online/offline events
    this.setupNetworkListeners()
  }

  /**
   * Setup network event listeners
   */
  setupNetworkListeners() {
    window.addEventListener('online', () => {
      console.log('Network connection restored')
      this.syncQueue()
    })

    window.addEventListener('offline', () => {
      console.log('Network connection lost')
    })
  }

  /**
   * Load queue from localStorage
   */
  loadQueue() {
    try {
      const stored = localStorage.getItem(QUEUE_STORAGE_KEY)
      return stored ? JSON.parse(stored) : []
    } catch (error) {
      console.error('Failed to load offline queue:', error)
      return []
    }
  }

  /**
   * Save queue to localStorage
   */
  saveQueue() {
    try {
      localStorage.setItem(QUEUE_STORAGE_KEY, JSON.stringify(this.queue))
    } catch (error) {
      console.error('Failed to save offline queue:', error)
    }
  }

  /**
   * Add operation to queue
   * @param {Object} operation - Operation to queue
   * @param {string} operation.type - Operation type (e.g., 'POST', 'PUT', 'DELETE')
   * @param {string} operation.url - API endpoint URL
   * @param {Object} operation.data - Request data
   * @param {Object} operation.headers - Request headers
   */
  enqueue(operation) {
    // Check if queue is full
    if (this.queue.length >= MAX_QUEUE_SIZE) {
      console.warn('Offline queue is full, removing oldest operation')
      this.queue.shift()
    }

    // Add timestamp and unique ID
    const queuedOperation = {
      id: this.generateId(),
      timestamp: Date.now(),
      ...operation
    }

    this.queue.push(queuedOperation)
    this.saveQueue()

    console.log('Operation queued:', queuedOperation.id)
    return queuedOperation.id
  }

  /**
   * Remove operation from queue
   */
  dequeue(operationId) {
    const index = this.queue.findIndex(op => op.id === operationId)
    if (index !== -1) {
      this.queue.splice(index, 1)
      this.saveQueue()
    }
  }

  /**
   * Get all queued operations
   */
  getQueue() {
    return [...this.queue]
  }

  /**
   * Get queue size
   */
  size() {
    return this.queue.length
  }

  /**
   * Clear all queued operations
   */
  clear() {
    this.queue = []
    this.saveQueue()
  }

  /**
   * Check if network is online
   */
  isOnline() {
    return navigator.onLine
  }

  /**
   * Sync queued operations with server
   */
  async syncQueue() {
    if (this.isSyncing || !this.isOnline() || this.queue.length === 0) {
      return
    }

    this.isSyncing = true
    console.log(`Syncing ${this.queue.length} queued operations...`)

    const results = {
      success: [],
      failed: []
    }

    // Process queue in order
    const queueCopy = [...this.queue]
    
    for (const operation of queueCopy) {
      try {
        await this.executeOperation(operation)
        results.success.push(operation.id)
        this.dequeue(operation.id)
      } catch (error) {
        console.error('Failed to sync operation:', operation.id, error)
        results.failed.push({
          id: operation.id,
          error: error.message
        })
        
        // If it's a client error (4xx), remove from queue
        if (error.response && error.response.status >= 400 && error.response.status < 500) {
          console.warn('Removing invalid operation from queue:', operation.id)
          this.dequeue(operation.id)
        }
      }
    }

    this.isSyncing = false

    // Notify callbacks
    this.notifySyncComplete(results)

    console.log('Sync complete:', results)
    return results
  }

  /**
   * Execute a queued operation
   */
  async executeOperation(operation) {
    const { type, url, data, headers } = operation

    // Import axios dynamically to avoid circular dependencies
    const { default: axios } = await import('axios')

    const config = {
      method: type,
      url,
      data,
      headers: headers || {}
    }

    return axios(config)
  }

  /**
   * Register callback for sync completion
   */
  onSyncComplete(callback) {
    this.syncCallbacks.push(callback)
  }

  /**
   * Notify all registered callbacks
   */
  notifySyncComplete(results) {
    this.syncCallbacks.forEach(callback => {
      try {
        callback(results)
      } catch (error) {
        console.error('Sync callback error:', error)
      }
    })
  }

  /**
   * Generate unique ID for operation
   */
  generateId() {
    return `${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }
}

// Export singleton instance
export const offlineQueue = new OfflineQueue()

// Export class for testing
export { OfflineQueue }
