# Browser Issues Quick Reference

## Quick Lookup for Common Browser-Specific Issues

### Safari Issues

| Issue | Symptom | Solution | File |
|-------|---------|----------|------|
| Date parsing fails | Invalid Date errors | Use ISO 8601 format: `2024-01-07T12:00:00.000Z` | Any date handling code |
| 100vh includes address bar | Content cut off at bottom | Use CSS custom property `--vh` set by JS | `viewportFix.js`, `style.css` |
| Flexbox gap not working | Items not spaced correctly | Use margin fallback: `.item + .item { margin-left: 1rem; }` | `style-browser-fixes.css` |
| Backdrop filter not working | Blur effect missing | Add `-webkit-` prefix | `style-browser-fixes.css` |
| No momentum scrolling | Scrolling feels sluggish | Add `-webkit-overflow-scrolling: touch` | `style-browser-fixes.css` |
| Notch cuts off content | Content hidden by notch | Use `env(safe-area-inset-*)` | `style-browser-fixes.css` |

### Firefox Issues

| Issue | Symptom | Solution | File |
|-------|---------|----------|------|
| Scrollbar styling different | Scrollbars look wrong | Use `scrollbar-width` and `scrollbar-color` | `style-browser-fixes.css` |
| Date input styling | Date picker looks different | Use `appearance: none` or accept default | `style-browser-fixes.css` |
| Number input spinners | Spinners visible | Use `-moz-appearance: textfield` | `style-browser-fixes.css` |

### Chrome/Edge Issues

| Issue | Symptom | Solution | File |
|-------|---------|----------|------|
| Yellow autofill background | Inputs have yellow background | Override with `-webkit-box-shadow` | `style-browser-fixes.css` |
| Search input cancel button | X button appears | Hide with `-webkit-search-cancel-button` | `style-browser-fixes.css` |

### iOS Safari Issues

| Issue | Symptom | Solution | File |
|-------|---------|----------|------|
| Input zoom on focus | Page zooms when typing | Use `font-size: 16px` minimum | `style-browser-fixes.css` |
| Viewport height wrong | Content cut off | Use `--vh` custom property | `viewportFix.js` |
| No install prompt | Can't install PWA | iOS doesn't support install prompt | N/A - Expected behavior |
| Momentum scrolling missing | Scrolling not smooth | Add `-webkit-overflow-scrolling: touch` | `style-browser-fixes.css` |
| Safe area issues | Notch cuts content | Use `env(safe-area-inset-*)` | `style-browser-fixes.css` |

### Android Chrome Issues

| Issue | Symptom | Solution | File |
|-------|---------|----------|------|
| 300ms tap delay | Buttons feel slow | Use `touch-action: manipulation` | `style-browser-fixes.css` |
| Input zoom on focus | Page zooms when typing | Use `font-size: 16px` minimum | `style-browser-fixes.css` |
| Viewport zoom | Page zooms unexpectedly | Set `maximum-scale=1.0` in viewport meta | `viewportFix.js` |

## Browser Detection

Use the `viewportFix.js` utility to detect browsers:

```javascript
import { detectBrowser, isIOS, isAndroid } from '@/utils/viewportFix'

// Detect browser type
const browser = detectBrowser() // 'chrome', 'firefox', 'safari', 'edge'

// Detect platform
if (isIOS()) {
  // iOS-specific code
}

if (isAndroid()) {
  // Android-specific code
}
```

## Feature Detection

Use the `supportsFeature()` function:

```javascript
import { supportsFeature } from '@/utils/viewportFix'

if (supportsFeature('css-gap')) {
  // Use gap
} else {
  // Use margin fallback
}

if (supportsFeature('service-worker')) {
  // Register service worker
}
```

## CSS Classes Added by Browser Detection

The app automatically adds these classes to `<body>`:

- `browser-chrome` - Chrome browser
- `browser-firefox` - Firefox browser
- `browser-safari` - Safari browser
- `browser-edge` - Edge browser
- `is-ios` - iOS device
- `is-android` - Android device
- `is-touch` - Touch-capable device

Use these for browser-specific CSS:

```css
/* Safari-specific styles */
.browser-safari .my-element {
  -webkit-backdrop-filter: blur(10px);
}

/* iOS-specific styles */
.is-ios .full-height {
  height: calc(var(--vh, 1vh) * 100);
}

/* Touch device styles */
.is-touch .button {
  min-width: 44px;
  min-height: 44px;
}
```

## Common Fixes Checklist

When you encounter a browser issue:

1. **Check if it's a known issue** - Look in this document
2. **Check browser version** - Ensure it's a supported version
3. **Check console for errors** - Look for specific error messages
4. **Test in other browsers** - Confirm it's browser-specific
5. **Check CSS support** - Use caniuse.com
6. **Add fallback** - Provide alternative for unsupported features
7. **Test fix** - Verify fix works in affected browser
8. **Document fix** - Add to this document

## Testing Tools

### Desktop
- **Chrome DevTools** - F12, Device emulation
- **Firefox DevTools** - F12, Responsive design mode
- **Safari Web Inspector** - Develop menu → Show Web Inspector

### Mobile
- **iOS Safari** - Connect iPhone to Mac, use Safari Web Inspector
- **Android Chrome** - Enable USB debugging, use chrome://inspect
- **BrowserStack** - Cloud-based testing on real devices
- **Sauce Labs** - Cloud-based testing on real devices

### Automated
- **Vitest** - Unit and integration tests
- **Playwright** - E2E testing across browsers
- **Cypress** - E2E testing with visual debugging

## Resources

- [Can I Use](https://caniuse.com/) - Browser compatibility tables
- [MDN Web Docs](https://developer.mozilla.org/) - Browser compatibility data
- [Autoprefixer](https://autoprefixer.github.io/) - CSS vendor prefix tool
- [Browserslist](https://browsersl.ist/) - Target browser configuration
- [WebKit Blog](https://webkit.org/blog/) - Safari updates
- [Chrome Platform Status](https://chromestatus.com/) - Chrome feature status
- [Firefox Release Notes](https://www.mozilla.org/firefox/releases/) - Firefox updates

## Emergency Fixes

If you need a quick fix in production:

### 1. Add Browser-Specific CSS
```css
/* Add to style.css */
@supports (-webkit-backdrop-filter: blur(10px)) {
  .blur { -webkit-backdrop-filter: blur(10px); }
}
```

### 2. Add JavaScript Polyfill
```javascript
// Add to main.js
if (!Element.prototype.animate) {
  // Load polyfill
  import('web-animations-js')
}
```

### 3. Add User Agent Detection
```javascript
// Add to viewportFix.js
const isSafari = /^((?!chrome|android).)*safari/i.test(navigator.userAgent)
if (isSafari) {
  // Safari-specific fix
}
```

### 4. Add Feature Detection
```javascript
// Add to component
if ('IntersectionObserver' in window) {
  // Use IntersectionObserver
} else {
  // Fallback to scroll events
}
```

## Contact

For browser-specific issues not covered here:
1. Check the main `BROWSER_COMPATIBILITY.md` document
2. Search the issue tracker
3. Create a new issue with browser details
4. Tag with `browser-compatibility` label
