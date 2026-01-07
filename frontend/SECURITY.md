# Security Implementation

This document outlines the security measures implemented in the AI Fitness Frontend application.

## Overview

The application implements multiple layers of security to protect user data and prevent common web vulnerabilities.

## Security Features

### 1. Input Sanitization (Requirement 14.4)

**Location:** `src/utils/sanitizer.js`

All user inputs are sanitized before rendering to prevent XSS (Cross-Site Scripting) attacks.

**Features:**
- HTML escaping for special characters
- Script tag removal
- Event handler removal
- JavaScript protocol blocking
- URL sanitization
- Filename sanitization for path traversal prevention

**Usage:**
```javascript
import { sanitizeInput } from '@/utils/sanitizer'

// Sanitize user input
const cleanInput = sanitizeInput(userInput)

// Use as Vue directive
<input v-sanitize v-model="userInput" />
```

**Applied to:**
- Assessment form (injury history, health conditions)
- Meal record form (notes)
- All text inputs that accept free-form user content

### 2. Secure Token Storage (Requirement 14.1)

**Location:** `src/utils/secureStorage.js`

Tokens are stored securely with multiple protection layers:

**Features:**
- Prefers sessionStorage over localStorage (cleared on tab close)
- XOR encryption/obfuscation in production
- Browser fingerprint-based encryption key
- Token expiry tracking
- Automatic cleanup of expired tokens
- Fallback to in-memory storage if Web Storage unavailable

**Security Benefits:**
- Tokens are not stored in plain text in production
- SessionStorage provides better security than localStorage
- Automatic expiry prevents stale token usage
- In-memory fallback prevents storage errors

**Usage:**
```javascript
import secureStorage from '@/utils/secureStorage'

// Store tokens
secureStorage.setTokens(accessToken, refreshToken, expiresIn)

// Retrieve tokens
const token = secureStorage.getAccessToken()

// Check expiry
if (secureStorage.isTokenExpired()) {
  // Refresh token
}

// Clear tokens
secureStorage.clearTokens()
```

### 3. Content Security Policy (Requirement 14.4)

**Location:** `index.html`, `vite-plugin-csp.js`

CSP headers prevent unauthorized script execution and data exfiltration.

**Policy:**
```
default-src 'self';
script-src 'self' 'unsafe-inline' 'unsafe-eval';
style-src 'self' 'unsafe-inline';
img-src 'self' data: https: blob:;
font-src 'self' data:;
connect-src 'self' [API endpoints];
frame-src 'none';
object-src 'none';
base-uri 'self';
form-action 'self';
upgrade-insecure-requests;
```

**Additional Headers:**
- `X-Content-Type-Options: nosniff` - Prevents MIME type sniffing
- `X-Frame-Options: DENY` - Prevents clickjacking
- `X-XSS-Protection: 1; mode=block` - Browser XSS protection
- `Referrer-Policy: strict-origin-when-cross-origin` - Controls referrer information
- `Permissions-Policy` - Restricts browser features

**Note:** `unsafe-inline` and `unsafe-eval` are currently needed for Vue development. For production, consider:
- Using nonces for inline scripts
- Removing `unsafe-eval` if not needed
- Implementing stricter CSP

### 4. Error Log Sanitization (Requirement 14.6)

**Location:** `src/utils/errorHandler.js`

Error logs are sanitized to prevent sensitive data leakage.

**Features:**
- Removes tokens from error messages
- Redacts API keys
- Sanitizes passwords
- Removes email addresses from logs
- Cleans stack traces

**Patterns Removed:**
- Bearer tokens
- API keys
- Passwords
- Email addresses
- Credit card numbers

### 5. Authentication Security

**Token Management:**
- Access tokens stored securely
- Automatic token refresh on expiry
- Tokens cleared on logout
- Failed refresh redirects to login

**API Client Security:**
- Authorization header added automatically
- Token refresh with retry queue
- Failed requests queued during refresh
- Network error retry logic

### 6. HTTPS Enforcement

**Production:**
- `upgrade-insecure-requests` CSP directive
- All API calls use HTTPS in production
- Service worker only works over HTTPS

## Security Testing

### Development Testing

Security testing utilities are available in development mode:

```javascript
// Run all security tests
window.securityTest.runAll()

// Test XSS prevention
window.securityTest.testXSS('<script>alert("XSS")</script>')

// Test URL sanitization
window.securityTest.testURL('javascript:alert("XSS")')

// Test CSP compliance
window.securityTest.testCSP()
```

### XSS Test Vectors

The application is tested against common XSS attack vectors:
- Script injection
- Event handlers
- JavaScript protocol
- Data protocol
- HTML entities
- Mixed case
- Encoded attacks
- Unicode
- Null bytes
- Nested tags
- Style injection
- Meta refresh
- Form injection

## Best Practices

### For Developers

1. **Always sanitize user input** before rendering
2. **Never store sensitive data** in localStorage without encryption
3. **Use the secure storage utility** for tokens
4. **Test with XSS vectors** before deploying
5. **Review CSP violations** in browser console
6. **Keep dependencies updated** for security patches

### Input Handling

```javascript
// ✅ Good - Sanitized
const cleanInput = sanitizeInput(userInput)
element.textContent = cleanInput

// ❌ Bad - Unsanitized
element.innerHTML = userInput
```

### Token Storage

```javascript
// ✅ Good - Secure storage
secureStorage.setAccessToken(token)

// ❌ Bad - Plain localStorage
localStorage.setItem('token', token)
```

### URL Handling

```javascript
// ✅ Good - Sanitized URL
const cleanUrl = sanitizeUrl(userUrl)
window.location.href = cleanUrl

// ❌ Bad - Unsanitized URL
window.location.href = userUrl
```

## Security Checklist

- [x] Input sanitization implemented
- [x] XSS prevention tested
- [x] Secure token storage
- [x] CSP headers configured
- [x] Error log sanitization
- [x] HTTPS enforcement
- [x] Authentication security
- [x] Security testing utilities
- [ ] Regular security audits
- [ ] Dependency vulnerability scanning
- [ ] Penetration testing

## Known Limitations

1. **CSP unsafe-inline/unsafe-eval**: Required for Vue development. Should be removed in production with proper nonce implementation.

2. **XOR Encryption**: The token obfuscation uses simple XOR encryption. For production, consider:
   - Web Crypto API for stronger encryption
   - Server-side token encryption
   - HttpOnly cookies for tokens

3. **Client-side Storage**: Tokens are stored client-side. Consider:
   - HttpOnly cookies for better security
   - Server-side session management
   - Shorter token expiry times

## Reporting Security Issues

If you discover a security vulnerability, please:
1. Do not open a public issue
2. Contact the security team directly
3. Provide detailed reproduction steps
4. Allow time for patching before disclosure

## References

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Content Security Policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
- [Web Security Guidelines](https://infosec.mozilla.org/guidelines/web_security)
- [Vue.js Security Best Practices](https://vuejs.org/guide/best-practices/security.html)
