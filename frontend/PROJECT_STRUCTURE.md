# Frontend Project Structure

## Directory Overview

```
frontend/
├── public/                 # Static assets served directly
│   └── vite.svg           # Favicon
├── src/
│   ├── assets/            # Static assets (images, fonts) processed by Vite
│   ├── components/        # Reusable Vue components
│   │   ├── common/        # Common UI components (buttons, modals, etc.)
│   │   └── fitness/       # Domain-specific components (training cards, meal cards)
│   ├── composables/       # Vue Composition API composables
│   ├── locales/           # i18n translation files
│   │   ├── en.json        # English translations
│   │   ├── zh.json        # Chinese translations
│   │   └── index.js       # i18n configuration
│   ├── router/            # Vue Router configuration
│   │   └── index.js       # Route definitions and navigation guards
│   ├── services/          # API service layer
│   ├── stores/            # Pinia state management stores
│   ├── tests/             # Test utilities and setup
│   │   └── setup.js       # Vitest global setup
│   ├── utils/             # Utility functions and helpers
│   │   └── constants.js   # Application constants
│   ├── views/             # Page-level components (routes)
│   │   ├── LoginView.vue
│   │   ├── RegisterView.vue
│   │   ├── DashboardView.vue
│   │   ├── ProfileView.vue
│   │   ├── AIConfigView.vue
│   │   ├── TrainingView.vue
│   │   ├── NutritionView.vue
│   │   └── StatisticsView.vue
│   ├── App.vue            # Root component
│   ├── main.js            # Application entry point
│   └── style.css          # Global styles
├── .env                   # Base environment variables
├── .env.development       # Development environment variables
├── .env.production        # Production environment variables
├── index.html             # HTML entry point
├── package.json           # Dependencies and scripts
├── postcss.config.js      # PostCSS configuration
├── vite.config.js         # Vite build configuration
└── README.md              # Project documentation
```

## Key Files

### Configuration Files

- **vite.config.js**: Vite build tool configuration with mobile optimization
- **postcss.config.js**: PostCSS plugins for autoprefixer and px-to-rem conversion
- **package.json**: Project dependencies and npm scripts

### Entry Points

- **index.html**: HTML template with mobile viewport meta tags
- **src/main.js**: JavaScript entry point, initializes Vue app with plugins
- **src/App.vue**: Root Vue component

### Core Modules

- **src/router/index.js**: Route definitions and navigation guards
- **src/locales/index.js**: i18n configuration with language detection
- **src/utils/constants.js**: Application-wide constants

## Naming Conventions

### Files
- Vue components: PascalCase (e.g., `LoginView.vue`, `TrainingCard.vue`)
- JavaScript modules: camelCase (e.g., `authService.js`, `errorHandler.js`)
- Test files: `*.test.js` or `*.spec.js`

### Components
- Views (pages): `*View.vue` (e.g., `DashboardView.vue`)
- Reusable components: Descriptive names (e.g., `NavigationBar.vue`, `LoadingSpinner.vue`)

### Stores
- Pinia stores: `use*Store` (e.g., `useAuthStore`, `useTrainingStore`)

### Services
- API services: `*Service` or `*.service.js` (e.g., `authService`, `training.service.js`)

## Development Workflow

1. **Start Development Server**: `npm run dev`
2. **Run Tests**: `npm run test:unit` or `npm run test:watch`
3. **Build for Production**: `npm run build`
4. **Preview Production Build**: `npm run preview`

## Mobile Optimization Features

- Viewport configuration for mobile devices
- px-to-rem conversion (base: 37.5px)
- Touch target minimum size: 44x44px
- Safe area support for iOS notch
- Pull-to-refresh functionality
- Virtual scrolling for performance

## Testing Strategy

- **Unit Tests**: Vitest + Vue Test Utils
- **Property-Based Tests**: fast-check library
- **Test Setup**: `src/tests/setup.js` with mocks for localStorage, sessionStorage, etc.

## Next Steps

The project structure is now set up. The following tasks will implement:
- Task 2: Core Infrastructure (API client, error handling, stores)
- Task 5: Authentication Services and Views
- Task 6: Common Components
- Task 7+: Feature-specific implementations
