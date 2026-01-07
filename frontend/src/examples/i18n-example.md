# Internationalization (i18n) Implementation

## Overview

The AI Fitness Frontend now has full internationalization support with Vue I18n, supporting English (en) and Chinese (zh) languages.

## Features Implemented

### 1. Vue I18n Setup ✅
- Configured Vue I18n plugin with automatic language detection
- Browser language detection with fallback to English
- Language persistence in localStorage

### 2. Translation Files ✅
- Comprehensive English translations (`src/locales/en.json`)
- Comprehensive Chinese translations (`src/locales/zh.json`)
- Organized by feature/component (auth, profile, training, nutrition, etc.)

### 3. Language Switcher Component ✅
- Dropdown menu component for language selection
- Integrated into ProfileView
- Automatic language persistence
- Success toast notification on language change

### 4. Locale Formatting Utilities ✅
- Date formatting with locale awareness
- Number formatting with locale awareness
- Specialized formatters for:
  - Dates (full, short, relative time)
  - Numbers (decimal, percentage, currency)
  - Duration (minutes to hours/minutes)
  - Weight (kg/lbs)
  - Distance (meters/kilometers)

## Usage Examples

### Using Translations in Components

```vue
<template>
  <div>
    <h1>{{ t('dashboard.title') }}</h1>
    <p>{{ t('dashboard.welcome') }}</p>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>
```

### Using Locale Formatting

```vue
<script setup>
import { useLocale } from '@/composables/useLocale'

const { formatDate, formatNumber, formatWeight } = useLocale()

// Format a date
const formattedDate = formatDate(new Date()) // "January 7, 2026" or "2026年1月7日"

// Format a number
const formattedNumber = formatNumber(1234.56) // "1,234.56" or "1,234.56"

// Format weight
const weight = formatWeight(75.5, 'kg') // "75.5 kg" or "75.5 公斤"
</script>
```

### Language Switcher Component

```vue
<template>
  <LanguageSwitcher />
</template>

<script setup>
import { LanguageSwitcher } from '@/components/common'
</script>
```

## File Structure

```
frontend/src/
├── locales/
│   ├── index.js          # i18n configuration
│   ├── en.json           # English translations
│   └── zh.json           # Chinese translations
├── composables/
│   └── useLocale.js      # Locale formatting composable
├── utils/
│   └── localeFormatter.js # Locale formatting utilities
└── components/
    └── common/
        └── LanguageSwitcher.vue # Language switcher component
```

## Supported Languages

- **English (en)**: Default language
- **Chinese (zh)**: Full translation support

## Language Detection

The application automatically detects the user's browser language on first load:
1. Checks localStorage for saved language preference
2. Falls back to browser language (navigator.language)
3. Defaults to English if no match found

## Language Persistence

Language preference is automatically saved to localStorage when changed, ensuring the user's choice persists across sessions.

## Locale-Specific Formatting

All dates, numbers, and specialized values are formatted according to the current locale:

### Date Formatting
- English: "January 7, 2026"
- Chinese: "2026年1月7日"

### Number Formatting
- English: "1,234.56"
- Chinese: "1,234.56"

### Duration Formatting
- English: "2h 30m"
- Chinese: "2小时30分钟"

### Distance Formatting
- English: "5.2 km"
- Chinese: "5.2公里"

## Requirements Validated

✅ **Requirement 13.1**: Language detection implemented
✅ **Requirement 13.2**: Translation keys used throughout UI
✅ **Requirement 13.3**: Language switcher with persistence
✅ **Requirement 13.4**: Locale-specific date and number formatting
✅ **Requirement 13.5**: Language preference persistence

## Testing

To test the i18n implementation:

1. Open the application
2. Navigate to Profile page
3. Use the language switcher to change between English and Chinese
4. Observe that all UI text updates immediately
5. Refresh the page to verify language persistence
6. Check that dates and numbers are formatted according to the selected locale
