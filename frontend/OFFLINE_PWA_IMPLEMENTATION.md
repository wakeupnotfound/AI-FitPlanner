# Offline Support and PWA Implementation

## Overview

This document describes the implementation of offline support and Progressive Web App (PWA) features for the AI Fitness Frontend application.

## Implemented Features

### 1. Service Worker Configuration (Task 17.1)

**Files Modified:**
- `frontend/vite.config.js` - Added VitePWA plugin configuration
- `frontend/src/main.js` - Added service worker registration
- `frontend/package.json` - Added vite-plugin-pwa dependency

**Key Features:**
- **Auto-update registration**: Service worker updates automatically when new version is available
- **PWA Manifest**: Configured app name, icons, theme colors, and display mode
- **Caching Strategies**:
  - **API Responses**: NetworkFirst strategy with 24-hour cache expiration
  - **Images**: CacheFirst strategy with 30-day cache expiration
  - **Fonts**: CacheFirst strategy with 1-year cache expiration
  - **CSS/JS**: StaleWhileRevalidate strategy with 7-day cache expiration
- **Cache Management**: Automatic cleanup of outdated caches
- **Skip Waiting**: New service worker activates immediately

**Caching Strategy Details:**

```javascript
// API Responses - NetworkFirst
// Tries network first, falls back to cache if offline
// Timeout: 10 seconds
// Cache: 24 hours, max 100 entries

// Images - CacheFirst
// Serves from cache if available, fetches if not
// Cache: 30 days, max 50 entries

// Fonts - CacheFirst
// Serves from cache if available, fetches if not
// Cache: 1 year, max 10 entries

// Static Resources (CSS/JS) - StaleWhileRevalidate
// Serves from cache immediately, updates in background
// Cache: 7 days, max 60 entries
```

### 2. Offline Queue Implementation (Task 17.3)

**Files Created:**
- `frontend/src/utils/offlineQueue.js` - Offline queue manager

**Files Modified:**
- `frontend/src/services/api.js` - Integrated offline queue with API client

**Key Features:**
- **Operation Queueing**: Automatically queues write operations (POST, PUT, DELETE, PATCH) when offline
- **Automatic Sync**: Syncs queued operations when connection is restored
- **Queue Persistence**: Stores queue in localStorage to survive page reloads
- **Queue Management**:
  - Maximum queue size: 100 operations
  - Oldest operations removed when queue is full
  - Invalid operations (4xx errors) automatically removed during sync
- **Network Event Listeners**: Listens for online/offline events
- **Sync Callbacks**: Notifies components when sync completes

**Usage Example:**

```javascript
// When offline, operations are automatically queued
await api.post('/training/record', workoutData)
// Returns: { queued: true, operationId: '...', message: '...' }

// When back online, operations sync automatically
// Or manually trigger sync:
await offlineQueue.syncQueue()
```

### 3. Network Status Indicator (Task 17.4)

**Files Created:**
- `frontend/src/composables/useNetworkStatus.js` - Network status composable
- `frontend/src/components/common/NetworkStatusIndicator.vue` - Network status UI component

**Files Modified:**
- `frontend/src/App.vue` - Added network status indicator to app layout
- `frontend/src/components/common/index.js` - Exported new component
- `frontend/src/locales/en.json` - Added English translations
- `frontend/src/locales/zh.json` - Added Chinese translations

**Key Features:**
- **Offline Banner**: Displays prominent banner when offline
- **Queue Status**: Shows number of queued operations
- **Syncing Indicator**: Shows when operations are being synced
- **Data Freshness**: Optional indicator showing when data was last updated
- **Reactive Updates**: Automatically updates based on network status changes

**UI Components:**

1. **Offline Banner** (Orange):
   - "You are offline"
   - Shows queued operation count
   - Closeable by user

2. **Syncing Banner** (Blue):
   - "Syncing queued operations..."
   - Shown during sync process

3. **Data Freshness Indicator**:
   - Shows relative time (e.g., "5m ago", "2h ago")
   - Optional, can be enabled per view
   - Useful for cached data

**Translations Added:**

```json
{
  "network": {
    "offline": "You are offline / 您已离线",
    "online": "Back online / 已恢复在线",
    "syncing": "Syncing queued operations... / 正在同步队列操作...",
    "queuedOperations": "{count} operations queued / {count} 个操作已排队",
    "justNow": "Just now / 刚刚",
    "minutesAgo": "{minutes}m ago / {minutes}分钟前",
    "hoursAgo": "{hours}h ago / {hours}小时前",
    "daysAgo": "{days}d ago / {days}天前"
  }
}
```

## Architecture

### Offline Queue Flow

```
User Action (POST/PUT/DELETE/PATCH)
    ↓
Check Network Status
    ↓
┌─────────────┬─────────────┐
│   Online    │   Offline   │
├─────────────┼─────────────┤
│ Execute API │ Queue       │
│ Request     │ Operation   │
└─────────────┴─────────────┘
                    ↓
            Save to localStorage
                    ↓
            Wait for Online Event
                    ↓
            Sync Queue Automatically
                    ↓
            Execute Queued Operations
                    ↓
            Update UI (Success/Failure)
```

### Service Worker Caching Flow

```
Request
    ↓
Service Worker Intercepts
    ↓
Check Cache Strategy
    ↓
┌──────────────┬──────────────┬──────────────┐
│ NetworkFirst │ CacheFirst   │ StaleWhile   │
│ (API)        │ (Images)     │ Revalidate   │
│              │              │ (CSS/JS)     │
├──────────────┼──────────────┼──────────────┤
│ Try Network  │ Try Cache    │ Serve Cache  │
│ ↓            │ ↓            │ ↓            │
│ If Fail:     │ If Miss:     │ Update in    │
│ Use Cache    │ Fetch        │ Background   │
└──────────────┴──────────────┴──────────────┘
```

## Requirements Validation

### Requirement 15.1: Offline Data Access
✅ **Implemented**: Service worker caches API responses, images, and static resources. Cached data is served when offline.

### Requirement 15.2: Queue Write Operations
✅ **Implemented**: Offline queue automatically queues POST, PUT, DELETE, and PATCH operations when offline.

### Requirement 15.3: Sync on Reconnection
✅ **Implemented**: Offline queue automatically syncs when connection is restored via online event listener.

### Requirement 9.2: Network Status Indicator
✅ **Implemented**: NetworkStatusIndicator component shows offline banner and syncing status.

### Requirement 15.4: Data Freshness Indicators
✅ **Implemented**: NetworkStatusIndicator supports optional data freshness display with relative time formatting.

## Testing Recommendations

### Manual Testing

1. **Offline Mode**:
   - Open DevTools → Network tab → Set to "Offline"
   - Try to perform write operations (record workout, add meal)
   - Verify operations are queued
   - Check localStorage for queued operations

2. **Sync on Reconnection**:
   - While offline, queue several operations
   - Set network back to "Online"
   - Verify operations sync automatically
   - Check console for sync logs

3. **Cache Behavior**:
   - Load app while online
   - Go offline
   - Navigate between pages
   - Verify cached content loads

4. **Network Indicator**:
   - Go offline → Verify orange banner appears
   - Queue operations → Verify count shows
   - Go online → Verify syncing banner appears
   - Verify banner disappears after sync

### Automated Testing

Property-based tests should be written for:
- **Property 6**: Offline data access (cached data retrieval)
- Queue persistence round trip
- Sync operation ordering
- Cache expiration behavior

## Usage in Components

### Using Network Status

```vue
<script setup>
import { useNetworkStatus } from '@/composables/useNetworkStatus'

const { isOnline, queueSize, isSyncing } = useNetworkStatus()

// Show different UI based on network status
const canSubmit = computed(() => isOnline.value || queueSize.value < 100)
</script>
```

### Showing Data Freshness

```vue
<template>
  <NetworkStatusIndicator 
    :show-freshness-indicator="true"
    :last-update-time="lastFetchTime"
  />
</template>

<script setup>
import { ref } from 'vue'
import { NetworkStatusIndicator } from '@/components/common'

const lastFetchTime = ref(new Date())

// Update when data is fetched
const fetchData = async () => {
  const data = await api.get('/data')
  lastFetchTime.value = new Date()
  return data
}
</script>
```

## Configuration

### PWA Manifest

The PWA manifest can be customized in `vite.config.js`:

```javascript
manifest: {
  name: 'AI Fitness Planner',
  short_name: 'AI Fitness',
  description: 'AI-powered personal fitness and nutrition planning system',
  theme_color: '#1989fa',
  background_color: '#ffffff',
  display: 'standalone',
  orientation: 'portrait',
  // Add custom icons here
  icons: [...]
}
```

### Cache Strategies

Cache strategies can be adjusted in `vite.config.js` under `workbox.runtimeCaching`:

```javascript
{
  urlPattern: /pattern/,
  handler: 'NetworkFirst' | 'CacheFirst' | 'StaleWhileRevalidate',
  options: {
    cacheName: 'cache-name',
    expiration: {
      maxEntries: 100,
      maxAgeSeconds: 86400
    }
  }
}
```

### Offline Queue

Queue settings can be modified in `offlineQueue.js`:

```javascript
const QUEUE_STORAGE_KEY = 'offline_queue'
const MAX_QUEUE_SIZE = 100
```

## Browser Support

- **Service Workers**: Chrome 40+, Firefox 44+, Safari 11.1+, Edge 17+
- **PWA Features**: Chrome 40+, Firefox 44+, Safari 11.3+, Edge 17+
- **localStorage**: All modern browsers

## Known Limitations

1. **Queue Size**: Limited to 100 operations to prevent localStorage overflow
2. **Token Expiration**: Queued operations may fail if auth token expires before sync
3. **Network Detection**: `navigator.onLine` may not be 100% accurate in all scenarios
4. **Cache Storage**: Limited by browser storage quotas (typically 50MB+)

## Future Enhancements

1. **IndexedDB**: Use IndexedDB for larger queue storage
2. **Background Sync**: Use Background Sync API for more reliable syncing
3. **Conflict Resolution**: Handle conflicts when syncing stale data
4. **Selective Caching**: Allow users to control what gets cached
5. **Cache Inspection**: UI for viewing and managing cached data

## Conclusion

The offline support and PWA implementation provides a robust foundation for the AI Fitness Frontend to work reliably in offline scenarios. Users can continue to use the app and queue operations when offline, with automatic synchronization when connectivity is restored.
