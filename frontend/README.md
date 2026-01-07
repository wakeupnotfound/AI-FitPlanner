# AI Fitness Frontend

A mobile-first Progressive Web Application built with Vue 3, Vite, and Vant UI for AI-powered fitness planning.

## Tech Stack

- **Framework**: Vue 3 (Composition API)
- **Build Tool**: Vite
- **UI Library**: Vant 4.x
- **State Management**: Pinia
- **Router**: Vue Router
- **HTTP Client**: Axios
- **Internationalization**: Vue I18n
- **Testing**: Vitest + Vue Test Utils + fast-check
- **CSS Processing**: PostCSS + Autoprefixer + postcss-pxtorem

## Project Structure

```
src/
├── assets/          # Static assets (images, fonts)
├── components/      # Reusable components
│   ├── common/      # Common UI components
│   └── fitness/     # Domain-specific components
├── composables/     # Composition API composables
├── locales/         # i18n language files
├── router/          # Vue Router configuration
├── services/        # API service layer
├── stores/          # Pinia stores
├── tests/           # Test setup and utilities
├── utils/           # Utility functions
├── views/           # Page-level components
├── App.vue          # Root component
├── main.js          # Application entry point
└── style.css        # Global styles
```

## Getting Started

### Prerequisites

- Node.js >= 18.0.0
- npm >= 9.0.0

### Installation

```bash
# Install dependencies
npm install
```

### Development

```bash
# Start development server
npm run dev

# The app will be available at http://localhost:5173
```

### Build

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

### Testing

```bash
# Run unit tests
npm run test:unit

# Run tests in watch mode
npm run test:watch

# Run tests with UI
npm run test:ui
```

## Environment Variables

Create `.env.development` and `.env.production` files:

```env
VITE_APP_TITLE=AI Fitness Planner
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_ENV=development
VITE_ENABLE_MOCK=false
VITE_LOG_LEVEL=debug
```

## Mobile Optimization

- Responsive design with viewport meta tags
- Touch-optimized UI with minimum 44x44px touch targets
- px to rem conversion for consistent scaling
- Safe area support for iOS devices
- Pull-to-refresh functionality
- Virtual scrolling for long lists

## Browser Support

- Chrome >= 61
- Safari >= 11
- Firefox >= 60
- Edge >= 79
- iOS Safari >= 11
- Android Chrome >= 61

## License

Private - All rights reserved
