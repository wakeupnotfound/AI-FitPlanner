<template>
  <div class="search-bar">
    <van-search
      v-model="localQuery"
      :placeholder="placeholder"
      :show-action="showAction"
      :clearable="clearable"
      @update:model-value="handleSearch"
      @search="handleSearchSubmit"
      @clear="handleClear"
      @cancel="handleCancel"
    >
      <template v-if="$slots.action" #action>
        <slot name="action" />
      </template>
      <template v-if="$slots.leftIcon" #left-icon>
        <slot name="leftIcon" />
      </template>
      <template v-if="$slots.rightIcon" #right-icon>
        <slot name="rightIcon" />
      </template>
    </van-search>

    <!-- Search Results -->
    <div v-if="showResults && (isSearching || searchResults.length > 0)" class="search-results">
      <!-- Loading State -->
      <div v-if="isSearching" class="search-loading">
        <van-loading size="24" />
        <span class="loading-text">{{ loadingText }}</span>
      </div>

      <!-- Results List -->
      <div v-else-if="searchResults.length > 0" class="results-list">
        <slot name="results" :results="searchResults">
          <van-cell
            v-for="(result, index) in searchResults"
            :key="index"
            :title="getResultTitle(result)"
            :label="getResultLabel(result)"
            clickable
            @click="handleResultClick(result)"
          />
        </slot>
      </div>

      <!-- No Results -->
      <div v-else class="no-results">
        <van-empty :description="noResultsText" image="search" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useDebouncedSearch, DEBOUNCE_DELAYS } from '../../composables/useDebounce'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Search...'
  },
  searchFn: {
    type: Function,
    required: true
  },
  delay: {
    type: Number,
    default: DEBOUNCE_DELAYS.SEARCH
  },
  minLength: {
    type: Number,
    default: 1
  },
  showAction: {
    type: Boolean,
    default: false
  },
  clearable: {
    type: Boolean,
    default: true
  },
  showResults: {
    type: Boolean,
    default: true
  },
  resultTitleKey: {
    type: String,
    default: 'title'
  },
  resultLabelKey: {
    type: String,
    default: 'label'
  },
  loadingText: {
    type: String,
    default: 'Searching...'
  },
  noResultsText: {
    type: String,
    default: 'No results found'
  }
})

const emit = defineEmits([
  'update:modelValue',
  'search',
  'select',
  'clear',
  'cancel'
])

const localQuery = ref(props.modelValue)

// Use debounced search
const {
  searchResults,
  isSearching,
  error,
  search,
  clearSearch
} = useDebouncedSearch(props.searchFn, {
  delay: props.delay,
  minLength: props.minLength
})

// Handle search input
const handleSearch = (value) => {
  localQuery.value = value
  emit('update:modelValue', value)
  
  if (value.length >= props.minLength) {
    search(value)
  } else {
    clearSearch()
  }
}

// Handle search submit (Enter key)
const handleSearchSubmit = (value) => {
  emit('search', value)
}

// Handle clear
const handleClear = () => {
  localQuery.value = ''
  emit('update:modelValue', '')
  emit('clear')
  clearSearch()
}

// Handle cancel
const handleCancel = () => {
  emit('cancel')
  clearSearch()
}

// Handle result click
const handleResultClick = (result) => {
  emit('select', result)
}

// Get result title
const getResultTitle = (result) => {
  if (typeof result === 'string') return result
  return result[props.resultTitleKey] || ''
}

// Get result label
const getResultLabel = (result) => {
  if (typeof result === 'string') return ''
  return result[props.resultLabelKey] || ''
}

// Watch for external changes
watch(() => props.modelValue, (newValue) => {
  if (newValue !== localQuery.value) {
    localQuery.value = newValue
    if (newValue) {
      search(newValue)
    } else {
      clearSearch()
    }
  }
})
</script>

<style scoped>
.search-bar {
  position: relative;
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border-radius: 0 0 8px 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  max-height: 400px;
  overflow-y: auto;
  z-index: 100;
}

.search-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  gap: 12px;
}

.loading-text {
  font-size: 14px;
  color: #969799;
}

.results-list {
  padding: 8px 0;
}

.no-results {
  padding: 24px;
}
</style>
