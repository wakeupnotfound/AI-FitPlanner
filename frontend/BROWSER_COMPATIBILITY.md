# Browser Compatibility Guide

## Overview

This document outlines the browser compatibility testing strategy and known issues for the AI Fitness Frontend application.

## Supported Browsers

### Desktop Browsers

| Browser | Minimum Version | Status | Notes |
|---------|----------------|--------|-------|
| Chrome | 90+ | ✅ Fully Supported | Primary development browser |
| Firefox | 88+ | ✅ Fully Supported | All features work |
| Safari | 14+ | ✅ Fully Supported | Some CSS features require prefixes |
| Edge | 90+ | ✅ Fully Supported | Chromium-based, same as Chrome |

### Mobile Browsers

| Browser | Minimum Version | Status | Notes |
|---------|----------------|--------|-------|
| Chrome (Android) | 90+ | ✅ Fully Supported | Primary mobile browser |
| Safari (iOS) | 14+ | ✅ Fully Supported | Requires viewport workarounds |
| Firefox (Android) | 88+ | ✅ Fully Supported | All features work |
| Samsung Internet | 14+ | ✅ Fully Supported | Chromium-based |

## Browser-Specific Issues and Workarounds

### Safari (Desktop & iOS)

#### Issue 1: Date Parsing
**Problem:** Safari has strict date parsing requirements.

**Solution:** Always use ISO 8601 format for dates:
```javascript
// ✅ Good - Works in all browsers
const date = new Date('2024-01-07T12:00:00.000Z')

// ❌ Bad - May fail in Safari
const date = new Date('2024-01-07 12:00:00')
```

#### Issue 2: Flexbox Gap Property
**Problem:** Safari 14.0 and earlier don't support `gap` in flexbox.

**Solution:** Use margin-based spacing as fallback:
```css
.flex-container {
  display: flex;
  gap: 1rem; /* Modern browsers */
}

.flex-container > * + * {
  margin-left: 1rem; /* Fallback for older Safari */
}
```

#### Issue 3: Backdrop Filter
**Problem:** Requires `-webkit-` prefix in Safari.

**Solution:** Include both prefixed and unprefixed versions:
```css
.blur-background {
  -webkit-backdrop-filter: blur(10px);
  backdrop-filter: blur(10px);
}
```

#### Issue 4: 100vh on Mobile Safari
**Problem:** iOS Safari's address bar affects viewport height calculations.

**Solution:** Use CSS custom properties with JavaScript:
```javascript
// Set actual viewport height
const setVH = () => {
  const vh = window.innerHeight * 0.01
  document.documentElement.style.setProperty('--vh', `${vh}px`)
}

window.addEventListener('resize', setVH)
setVH()
```

```css
.full-height {
  height: 100vh; /* Fallback */
  height: calc(var(--vh, 1vh) * 100); /* Actual height */
}
```

#### Issue 5: Safe Area Insets (iPhone X+)
**Problem:** Notch and home indicator require padding.

**Solution:** Use `env()` for safe areas:
```css
.app-container {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

#### Issue 6: Momentum Scrolling
**Problem:** iOS Safari doesn't have smooth momentum scrolling by default.

**Solution:** Enable webkit overflow scrolling:
```css
.scrollable {
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}
```

### Firefox

#### Issue 1: Scrollbar Styling
**Problem:** Firefox uses different scrollbar properties than Chrome/Safari.

**Solution:** Use Firefox-specific properties:
```css
/* Chrome/Safari */
.scrollable::-webkit-scrollbar {
  width: 8px;
}

/* Firefox */
.scrollable {
  scrollbar-width: thin;
  scrollbar-color: #888 #f1f1f1;
}
```

#### Issue 2: Input Type Date
**Problem:** Firefox date picker has different styling.

**Solution:** Provide consistent styling or use a custom date picker:
```css
input[type="date"] {
  appearance: none;
  -moz-appearance: none;
  -webkit-appearance: none;
}
```

### Chrome/Edge

#### Issue 1: Autofill Styling
**Problem:** Chrome applies yellow background to autofilled inputs.

**Solution:** Override autofill styles:
```css
input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus {
  -webkit-box-shadow: 0 0 0 1000px white inset;
  -webkit-text-fill-color: #333;
}
```

### Android Chrome

#### Issue 1: 300ms Tap Delay
**Problem:** Older Android browsers have 300ms delay on touch events.

**Solution:** Use `touch-action` CSS property:
```css
button, a, input {
  touch-action: manipulation;
}
```

#### Issue 2: Viewport Zoom on Input Focus
**Problem:** Android Chrome zooms in when input font-size is less than 16px.

**Solution:** Use minimum 16px font size for inputs:
```css
input, select, textarea {
  font-size: 16px; /* Prevents zoom */
}
```

## Testing Strategy

### Automated Testing

We use Vitest with happy-dom for automated cross-browser compatibility testing:

```bash
# Run all tests
npm run test:unit

# Run with UI
npm run test:ui

# Watch mode
npm run test:watch
```

### Manual Testing Checklist

#### Desktop Testing
- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Safari (latest)
- [ ] Edge (latest)

#### Mobile Testing
- [ ] iOS Safari (iPhone)
- [ ] iOS Safari (iPad)
- [ ] Chrome (Android phone)
- [ ] Chrome (Android tablet)
- [ ] Firefox (Android)

### Test Scenarios

For each browser, test the following:

1. **Authentication Flow**
   - [ ] Registration
   - [ ] Login
   - [ ] Logout
   - [ ] Token refresh

2. **Responsive Design**
   - [ ] Layout adapts to screen size
   - [ ] Navigation works on mobile
   - [ ] Touch targets are adequate (44x44px minimum)
   - [ ] No horizontal scrolling

3. **Forms**
   - [ ] Input validation
   - [ ] Date pickers work
   - [ ] Number inputs work
   - [ ] Form submission

4. **Charts and Visualizations**
   - [ ] Charts render correctly
   - [ ] Animations work smoothly
   - [ ] Touch interactions work

5. **Offline Support**
   - [ ] Service worker registers
   - [ ] Cached data loads offline
   - [ ] Offline indicator shows

6. **Performance**
   - [ ] Initial load < 2s on 3G
   - [ ] Smooth scrolling
   - [ ] No janky animations

7. **Internationalization**
   - [ ] Language switching works
   - [ ] Date/number formatting correct
   - [ ] RTL support (if applicable)

## Browser DevTools

### Chrome DevTools
- Device emulation: F12 → Toggle device toolbar
- Network throttling: Network tab → Throttling dropdown
- Lighthouse audits: Lighthouse tab

### Firefox DevTools
- Responsive design mode: Ctrl+Shift+M (Cmd+Opt+M on Mac)
- Network throttling: Network tab → Throttling dropdown

### Safari DevTools
- Enable: Safari → Preferences → Advanced → Show Develop menu
- Responsive design mode: Develop → Enter Responsive Design Mode
- iOS Simulator: Xcode → Open Developer Tool → Simulator

## Known Limitations

### Progressive Web App (PWA)
- iOS Safari: Limited PWA support, no install prompt
- Firefox: No install prompt on desktop

### Service Workers
- iOS Safari: Service workers work but with limitations
- Private/Incognito mode: Service workers may not work

### Web APIs
- Notification API: Requires user permission, not supported in iOS Safari
- Geolocation API: Requires HTTPS and user permission
- Clipboard API: Limited support in older browsers

## Polyfills and Fallbacks

We use modern JavaScript features that are supported in all target browsers. No polyfills are required for:

- ES2015+ features (arrow functions, classes, promises, async/await)
- Fetch API
- localStorage/sessionStorage
- CSS Grid and Flexbox
- CSS Custom Properties

## Continuous Integration

Our CI pipeline tests on:
- Node.js 18+ (simulates modern browser environment)
- happy-dom (lightweight DOM implementation)

For full browser testing, we recommend:
- BrowserStack or Sauce Labs for cross-browser testing
- Playwright or Cypress for E2E testing

## Resources

- [Can I Use](https://caniuse.com/) - Browser compatibility tables
- [MDN Web Docs](https://developer.mozilla.org/) - Browser compatibility data
- [Autoprefixer](https://autoprefixer.github.io/) - CSS vendor prefix tool
- [Browserslist](https://browsersl.ist/) - Target browser configuration

## Reporting Issues

If you encounter a browser-specific issue:

1. Document the browser and version
2. Provide steps to reproduce
3. Include screenshots or screen recordings
4. Check if the issue exists in other browsers
5. Create an issue in the project repository

## Version History

- **v1.0.0** (2024-01-07): Initial browser compatibility documentation
