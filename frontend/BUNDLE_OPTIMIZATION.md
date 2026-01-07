# Bundle Optimization Guide

This document describes the bundle optimization strategies implemented in the AI Fitness Frontend application.

## Overview

The application uses several techniques to minimize bundle size and improve load times:

1. **Code Splitting** - Route-level and component-level splitting
2. **Tree Shaking** - Removing unused code
3. **Lazy Loading** - Loading components and images on demand
4. **Chunk Optimization** - Strategic grouping of dependencies
5. **Minification** - Compressing JavaScript and CSS

## Bundle Analysis

### Running Bundle Analysis

To analyze the bundle size and composition:

```bash
npm run build:analyze
```

This will:
1. Build the production bundle
2. Generate a visual report at `dist/stats.html`
3. Open the report in your browser (macOS)

For other platforms, manually open `dist/stats.html` after running `npm run build`.

### Understanding the Report

The bundle analyzer shows:
- **Size by module** - Which dependencies take up the most space
- **Chunk composition** - How code is split across chunks
- **Gzip/Brotli sizes** - Compressed sizes (what users actually download)
- **Module relationships** - Dependencies between modules

## Optimization Strategies

### 1. Code Splitting

#### Route-Level Splitting

All routes use dynamic imports for automatic code splitting:

```javascript
{
  path: '/dashboard',
  component: () => import('@/views/DashboardView.vue')
}
```

**Benefits:**
- Users only download code for routes they visit
- Faster initial load time
- Better caching (route chunks change independently)

#### Component-Level Splitting

Heavy components use `defineAsyncComponent`:

```javascript
const ChartWidget = defineAsyncComponent(() => 
  import('@/components/statistics/ChartWidget.vue')
)
```

**When to use:**
- Large components (>50KB)
- Components with heavy dependencies
- Components not needed on initial render
- Components used conditionally

### 2. Manual Chunk Configuration

Dependencies are grouped into logical chunks:

```javascript
manualChunks: {
  'vendor-vue': ['vue', 'vue-router', 'pinia'],
  'vendor-vant': ['vant'],
  'vendor-i18n': ['vue-i18n'],
  'vendor-axios': ['axios']
}
```

**Benefits:**
- Better caching (vendor code changes less frequently)
- Parallel loading of chunks
- Smaller individual chunk sizes

### 3. Tree Shaking

Aggressive tree shaking configuration:

```javascript
treeshake: {
  moduleSideEffects: false,
  propertyReadSideEffects: false,
  tryCatchDeoptimization: false
}
```

**Best Practices:**
- Use ES6 imports/exports (not CommonJS)
- Import only what you need: `import { ref } from 'vue'`
- Avoid importing entire libraries: `import _ from 'lodash'` ❌
- Use specific imports: `import debounce from 'lodash/debounce'` ✅

### 4. Lazy Loading

#### Images

Use the `v-lazy` directive or `LazyImage` component:

```vue
<!-- Directive -->
<img v-lazy="imageUrl" alt="description" />

<!-- Component -->
<LazyImage 
  :src="imageUrl" 
  :width="300" 
  :height="200"
  aspect-ratio="16/9"
/>
```

#### Components

Use `defineAsyncComponent` for conditional components:

```javascript
const HeavyComponent = defineAsyncComponent(() => 
  import('./HeavyComponent.vue')
)
```

### 5. Minification

Terser configuration for optimal compression:

```javascript
terserOptions: {
  compress: {
    drop_console: true,      // Remove console.log
    drop_debugger: true,     // Remove debugger statements
    pure_funcs: [            // Remove specific function calls
      'console.log',
      'console.info',
      'console.debug'
    ]
  }
}
```

## Performance Targets

### Bundle Size Targets

- **Initial Bundle**: < 200KB (gzipped)
- **Vendor Chunks**: < 150KB each (gzipped)
- **Route Chunks**: < 50KB each (gzipped)
- **Component Chunks**: < 30KB each (gzipped)

### Load Time Targets

- **First Contentful Paint (FCP)**: < 1.5s on 3G
- **Time to Interactive (TTI)**: < 3.5s on 3G
- **Largest Contentful Paint (LCP)**: < 2.5s on 3G

## Monitoring Bundle Size

### During Development

Watch for warnings during build:

```bash
npm run build
```

Vite will warn about chunks exceeding the size limit (1000KB by default).

### In CI/CD

Add bundle size checks to your CI pipeline:

```yaml
- name: Check bundle size
  run: |
    npm run build
    # Add size check script here
```

## Common Issues and Solutions

### Issue: Large Vendor Chunks

**Problem:** `vendor-vant` chunk is too large

**Solutions:**
1. Import only needed Vant components
2. Use Vant's tree-shaking features
3. Consider alternative lighter UI libraries for specific components

### Issue: Duplicate Dependencies

**Problem:** Same dependency appears in multiple chunks

**Solutions:**
1. Check `manualChunks` configuration
2. Ensure consistent import paths
3. Use `optimizeDeps.include` to pre-bundle dependencies

### Issue: Large Route Chunks

**Problem:** A route chunk exceeds 100KB

**Solutions:**
1. Split large components within the route
2. Use lazy loading for heavy components
3. Move shared code to separate chunks
4. Consider virtual scrolling for large lists

## Best Practices

### DO:
✅ Use dynamic imports for routes
✅ Lazy load heavy components
✅ Import only what you need
✅ Use tree-shakeable libraries
✅ Optimize images (WebP, compression)
✅ Enable gzip/brotli compression on server
✅ Monitor bundle size regularly

### DON'T:
❌ Import entire libraries
❌ Include unused dependencies
❌ Ignore bundle size warnings
❌ Load all components eagerly
❌ Use large unoptimized images
❌ Disable tree shaking
❌ Skip bundle analysis

## Tools and Resources

### Analysis Tools
- **Rollup Plugin Visualizer** - Visual bundle analysis
- **Vite Bundle Analyzer** - Built-in size reporting
- **Lighthouse** - Performance auditing
- **WebPageTest** - Real-world performance testing

### Optimization Resources
- [Vite Performance Guide](https://vitejs.dev/guide/performance.html)
- [Vue Performance Guide](https://vuejs.org/guide/best-practices/performance.html)
- [Web.dev Performance](https://web.dev/performance/)

## Maintenance

### Regular Tasks

1. **Weekly**: Review bundle size after major changes
2. **Monthly**: Run full bundle analysis
3. **Quarterly**: Audit dependencies for updates/alternatives
4. **Before Release**: Verify all performance targets met

### Dependency Updates

When updating dependencies:
1. Check bundle size impact
2. Test tree shaking still works
3. Verify code splitting unchanged
4. Run performance tests

## Conclusion

Bundle optimization is an ongoing process. Regular monitoring and analysis ensure the application remains fast and efficient as it grows.

For questions or suggestions, please contact the development team.
