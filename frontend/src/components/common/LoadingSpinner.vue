<template>
  <!-- Full-page overlay loading -->
  <van-overlay 
    v-if="overlay" 
    :show="loading" 
    :z-index="zIndex"
    class="loading-overlay"
  >
    <div class="loading-wrapper">
      <van-loading 
        :type="type" 
        :color="color" 
        :size="size"
        :text-color="textColor"
        vertical
      >
        {{ text || t('app.loading') }}
      </van-loading>
    </div>
  </van-overlay>

  <!-- Skeleton screen variant -->
  <div v-else-if="skeleton" class="skeleton-wrapper">
    <van-skeleton 
      :row="skeletonRows" 
      :row-width="skeletonRowWidth"
      :title="skeletonTitle"
      :avatar="skeletonAvatar"
      :avatar-size="skeletonAvatarSize"
      :avatar-shape="skeletonAvatarShape"
      :loading="loading"
      :animate="animate"
    >
      <slot />
    </van-skeleton>
  </div>

  <!-- Inline loading spinner -->
  <div v-else class="inline-loading" :class="{ centered }">
    <van-loading 
      :type="type" 
      :color="color" 
      :size="size"
      :text-color="textColor"
      :vertical="vertical"
    >
      <span v-if="text || showDefaultText">{{ text || t('app.loading') }}</span>
    </van-loading>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  // Common props
  loading: {
    type: Boolean,
    default: true
  },
  text: {
    type: String,
    default: ''
  },
  showDefaultText: {
    type: Boolean,
    default: true
  },
  type: {
    type: String,
    default: 'circular',
    validator: (value) => ['circular', 'spinner'].includes(value)
  },
  color: {
    type: String,
    default: '#1989fa'
  },
  size: {
    type: [String, Number],
    default: '30px'
  },
  textColor: {
    type: String,
    default: '#969799'
  },
  vertical: {
    type: Boolean,
    default: false
  },
  centered: {
    type: Boolean,
    default: false
  },

  // Overlay variant props
  overlay: {
    type: Boolean,
    default: false
  },
  zIndex: {
    type: Number,
    default: 1000
  },

  // Skeleton variant props
  skeleton: {
    type: Boolean,
    default: false
  },
  skeletonRows: {
    type: Number,
    default: 3
  },
  skeletonRowWidth: {
    type: [String, Array],
    default: () => ['100%', '100%', '60%']
  },
  skeletonTitle: {
    type: Boolean,
    default: true
  },
  skeletonAvatar: {
    type: Boolean,
    default: false
  },
  skeletonAvatarSize: {
    type: String,
    default: '32px'
  },
  skeletonAvatarShape: {
    type: String,
    default: 'round',
    validator: (value) => ['round', 'square'].includes(value)
  },
  animate: {
    type: Boolean,
    default: true
  }
})
</script>

<style scoped>
.loading-overlay {
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  min-width: 120px;
}

.inline-loading {
  display: inline-flex;
  align-items: center;
  padding: 8px;
}

.inline-loading.centered {
  display: flex;
  justify-content: center;
  width: 100%;
  padding: 20px;
}

.skeleton-wrapper {
  width: 100%;
}
</style>
