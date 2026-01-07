/**
 * Vue directive for lazy loading images
 * Usage: <img v-lazy="imageUrl" alt="description" />
 */

import { getPlaceholderImage } from '../composables/useLazyLoad'

// Global observer instance
let observer = null

// Initialize observer
const initObserver = () => {
  if (!observer && 'IntersectionObserver' in window) {
    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const img = entry.target
            loadImage(img)
            observer.unobserve(img)
          }
        })
      },
      {
        rootMargin: '50px',
        threshold: 0.01
      }
    )
  }
}

// Load image
const loadImage = (img) => {
  const src = img.dataset.lazySrc
  if (!src) return

  img.classList.add('lazy-loading')

  const tempImg = new Image()
  tempImg.onload = () => {
    img.src = src
    img.classList.remove('lazy-loading')
    img.classList.add('lazy-loaded')
    img.removeAttribute('data-lazy-src')
  }
  tempImg.onerror = () => {
    img.classList.remove('lazy-loading')
    img.classList.add('lazy-error')
    // Set error placeholder
    img.src = getPlaceholderImage(img.width || 300, img.height || 200, '#ffebee')
  }
  tempImg.src = src
}

// Directive definition
export const lazyLoadDirective = {
  mounted(el, binding) {
    // Initialize observer if not already done
    initObserver()

    // Set placeholder
    if (!el.src || el.src === window.location.href) {
      el.src = getPlaceholderImage(el.width || 300, el.height || 200)
    }

    // Store the actual image URL
    el.dataset.lazySrc = binding.value

    // Add loading class
    el.classList.add('lazy-image')

    // If IntersectionObserver is not supported, load immediately
    if (!observer) {
      loadImage(el)
      return
    }

    // Start observing
    observer.observe(el)
  },

  updated(el, binding) {
    // If the image URL changes, update it
    if (binding.value !== binding.oldValue) {
      el.dataset.lazySrc = binding.value
      
      // If already loaded, load the new image immediately
      if (el.classList.contains('lazy-loaded')) {
        loadImage(el)
      }
    }
  },

  unmounted(el) {
    // Stop observing when element is unmounted
    if (observer) {
      observer.unobserve(el)
    }
  }
}

// CSS styles to add to your global styles
export const lazyLoadStyles = `
.lazy-image {
  transition: opacity 0.3s ease-in-out;
}

.lazy-loading {
  opacity: 0.6;
  filter: blur(5px);
}

.lazy-loaded {
  opacity: 1;
  filter: blur(0);
}

.lazy-error {
  opacity: 0.5;
  border: 1px dashed #ff5252;
}
`
