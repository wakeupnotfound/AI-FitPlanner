<template>
  <div class="chart-widget">
    <!-- Chart Header -->
    <div class="chart-header">
      <h3 class="chart-title">{{ title }}</h3>
      <div class="time-range-selector">
        <van-dropdown-menu>
          <van-dropdown-item 
            v-model="selectedRange" 
            :options="rangeOptions"
            @change="onRangeChange"
          />
        </van-dropdown-menu>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="chart-loading">
      <van-loading size="24px" />
    </div>

    <!-- No Data State -->
    <div v-else-if="!hasData" class="chart-empty">
      <van-empty 
        :description="t('statistics.noData')"
        image="search"
      />
    </div>

    <!-- Chart Content -->
    <div v-else class="chart-content">
      <!-- Line Chart for Trends -->
      <div v-if="type === 'line'" class="line-chart">
        <div class="chart-canvas" ref="chartCanvas">
          <svg :viewBox="`0 0 ${chartWidth} ${chartHeight}`" class="chart-svg">
            <!-- Grid Lines -->
            <g class="grid-lines">
              <line 
                v-for="(line, index) in gridLines" 
                :key="'grid-' + index"
                :x1="padding.left"
                :y1="line.y"
                :x2="chartWidth - padding.right"
                :y2="line.y"
                stroke="#ebedf0"
                stroke-dasharray="4,4"
              />
            </g>

            <!-- Y-Axis Labels -->
            <g class="y-axis-labels">
              <text 
                v-for="(line, index) in gridLines" 
                :key="'y-label-' + index"
                :x="padding.left - 8"
                :y="line.y + 4"
                text-anchor="end"
                class="axis-label"
              >
                {{ line.value }}
              </text>
            </g>

            <!-- Line Path -->
            <path 
              :d="linePath" 
              fill="none" 
              :stroke="lineColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />

            <!-- Area Fill -->
            <path 
              :d="areaPath" 
              :fill="areaFill"
              opacity="0.1"
            />

            <!-- Data Points -->
            <g class="data-points">
              <circle 
                v-for="(point, index) in dataPoints" 
                :key="'point-' + index"
                :cx="point.x"
                :cy="point.y"
                r="4"
                :fill="lineColor"
                @click="onPointClick(point, index)"
              />
            </g>

            <!-- X-Axis Labels -->
            <g class="x-axis-labels">
              <text 
                v-for="(label, index) in xAxisLabels" 
                :key="'x-label-' + index"
                :x="label.x"
                :y="chartHeight - 8"
                text-anchor="middle"
                class="axis-label"
              >
                {{ label.text }}
              </text>
            </g>
          </svg>
        </div>

        <!-- Legend -->
        <div v-if="showLegend" class="chart-legend">
          <span class="legend-item">
            <span class="legend-color" :style="{ backgroundColor: lineColor }"></span>
            <span class="legend-text">{{ legendLabel }}</span>
          </span>
        </div>
      </div>

      <!-- Bar Chart for Statistics -->
      <div v-else-if="type === 'bar'" class="bar-chart">
        <div class="chart-canvas" ref="chartCanvas">
          <svg :viewBox="`0 0 ${chartWidth} ${chartHeight}`" class="chart-svg">
            <!-- Grid Lines -->
            <g class="grid-lines">
              <line 
                v-for="(line, index) in gridLines" 
                :key="'grid-' + index"
                :x1="padding.left"
                :y1="line.y"
                :x2="chartWidth - padding.right"
                :y2="line.y"
                stroke="#ebedf0"
                stroke-dasharray="4,4"
              />
            </g>

            <!-- Y-Axis Labels -->
            <g class="y-axis-labels">
              <text 
                v-for="(line, index) in gridLines" 
                :key="'y-label-' + index"
                :x="padding.left - 8"
                :y="line.y + 4"
                text-anchor="end"
                class="axis-label"
              >
                {{ line.value }}
              </text>
            </g>

            <!-- Bars -->
            <g class="bars">
              <rect 
                v-for="(bar, index) in barData" 
                :key="'bar-' + index"
                :x="bar.x"
                :y="bar.y"
                :width="bar.width"
                :height="bar.height"
                :fill="getBarColor(index)"
                rx="4"
                @click="onBarClick(bar, index)"
              />
            </g>

            <!-- X-Axis Labels -->
            <g class="x-axis-labels">
              <text 
                v-for="(label, index) in xAxisLabels" 
                :key="'x-label-' + index"
                :x="label.x"
                :y="chartHeight - 8"
                text-anchor="middle"
                class="axis-label"
              >
                {{ label.text }}
              </text>
            </g>
          </svg>
        </div>
      </div>
    </div>

    <!-- Selected Point Info -->
    <div v-if="selectedPoint" class="selected-info">
      <span class="info-date">{{ selectedPoint.label }}</span>
      <span class="info-value">{{ selectedPoint.value }} {{ unit }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'line',
    validator: (value) => ['line', 'bar'].includes(value)
  },
  data: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  color: {
    type: String,
    default: '#1989fa'
  },
  unit: {
    type: String,
    default: ''
  },
  showLegend: {
    type: Boolean,
    default: false
  },
  legendLabel: {
    type: String,
    default: ''
  },
  defaultRange: {
    type: String,
    default: 'month'
  }
})

const emit = defineEmits(['range-change', 'point-click'])

// State
const selectedRange = ref(props.defaultRange)
const selectedPoint = ref(null)
const chartCanvas = ref(null)

// Chart dimensions
const chartWidth = 320
const chartHeight = 200
const padding = { top: 20, right: 20, bottom: 30, left: 45 }

// Range options
const rangeOptions = computed(() => [
  { text: t('statistics.week'), value: 'week' },
  { text: t('statistics.month'), value: 'month' },
  { text: t('statistics.threeMonths'), value: '3months' },
  { text: t('statistics.year'), value: 'year' }
])

// Computed
const hasData = computed(() => props.data && props.data.length > 0)

const lineColor = computed(() => props.color)
const areaFill = computed(() => props.color)

const dataValues = computed(() => {
  return props.data.map(item => item.value || 0)
})

const minValue = computed(() => {
  if (dataValues.value.length === 0) return 0
  return Math.min(...dataValues.value) * 0.9
})

const maxValue = computed(() => {
  if (dataValues.value.length === 0) return 100
  return Math.max(...dataValues.value) * 1.1
})

const valueRange = computed(() => maxValue.value - minValue.value || 1)

const chartAreaWidth = computed(() => chartWidth - padding.left - padding.right)
const chartAreaHeight = computed(() => chartHeight - padding.top - padding.bottom)

// Calculate data points for line chart
const dataPoints = computed(() => {
  if (!hasData.value) return []
  
  const points = []
  const step = chartAreaWidth.value / (props.data.length - 1 || 1)
  
  props.data.forEach((item, index) => {
    const x = padding.left + (index * step)
    const normalizedValue = (item.value - minValue.value) / valueRange.value
    const y = chartHeight - padding.bottom - (normalizedValue * chartAreaHeight.value)
    
    points.push({
      x,
      y,
      value: item.value,
      label: item.date || item.label || ''
    })
  })
  
  return points
})

// Generate line path
const linePath = computed(() => {
  if (dataPoints.value.length === 0) return ''
  
  return dataPoints.value.reduce((path, point, index) => {
    if (index === 0) {
      return `M ${point.x} ${point.y}`
    }
    return `${path} L ${point.x} ${point.y}`
  }, '')
})

// Generate area path
const areaPath = computed(() => {
  if (dataPoints.value.length === 0) return ''
  
  const baseline = chartHeight - padding.bottom
  const firstPoint = dataPoints.value[0]
  const lastPoint = dataPoints.value[dataPoints.value.length - 1]
  
  return `${linePath.value} L ${lastPoint.x} ${baseline} L ${firstPoint.x} ${baseline} Z`
})

// Calculate bar data
const barData = computed(() => {
  if (!hasData.value) return []
  
  const bars = []
  const barWidth = (chartAreaWidth.value / props.data.length) * 0.7
  const gap = (chartAreaWidth.value / props.data.length) * 0.3
  
  props.data.forEach((item, index) => {
    const x = padding.left + (index * (barWidth + gap)) + (gap / 2)
    const normalizedValue = (item.value - minValue.value) / valueRange.value
    const height = normalizedValue * chartAreaHeight.value
    const y = chartHeight - padding.bottom - height
    
    bars.push({
      x,
      y,
      width: barWidth,
      height,
      value: item.value,
      label: item.date || item.label || ''
    })
  })
  
  return bars
})

// Grid lines
const gridLines = computed(() => {
  const lines = []
  const numLines = 5
  const step = valueRange.value / (numLines - 1)
  
  for (let i = 0; i < numLines; i++) {
    const value = minValue.value + (step * i)
    const y = chartHeight - padding.bottom - ((i / (numLines - 1)) * chartAreaHeight.value)
    lines.push({
      y,
      value: Math.round(value * 10) / 10
    })
  }
  
  return lines
})

// X-axis labels
const xAxisLabels = computed(() => {
  if (!hasData.value) return []
  
  const labels = []
  const maxLabels = 5
  const step = Math.ceil(props.data.length / maxLabels)
  
  props.data.forEach((item, index) => {
    if (index % step === 0 || index === props.data.length - 1) {
      const x = props.type === 'line' 
        ? dataPoints.value[index]?.x 
        : barData.value[index]?.x + (barData.value[index]?.width / 2)
      
      labels.push({
        x,
        text: formatDateLabel(item.date || item.label)
      })
    }
  })
  
  return labels
})

// Methods
const formatDateLabel = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

const getBarColor = (index) => {
  const colors = ['#1989fa', '#07c160', '#ffa940', '#ff6b6b', '#722ed1']
  return colors[index % colors.length]
}

const onRangeChange = (value) => {
  emit('range-change', value)
}

const onPointClick = (point, index) => {
  selectedPoint.value = point
  emit('point-click', { point, index })
}

const onBarClick = (bar, index) => {
  selectedPoint.value = bar
  emit('point-click', { point: bar, index })
}

// Watch for range changes
watch(selectedRange, (newRange) => {
  emit('range-change', newRange)
})
</script>

<style scoped>
.chart-widget {
  background: white;
  border-radius: 8px;
  padding: 16px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.chart-title {
  font-size: 16px;
  font-weight: 500;
  color: #323233;
  margin: 0;
}

.time-range-selector {
  min-width: 100px;
}

.time-range-selector :deep(.van-dropdown-menu__bar) {
  height: 32px;
  box-shadow: none;
  background: #f7f8fa;
  border-radius: 4px;
}

.time-range-selector :deep(.van-dropdown-menu__title) {
  font-size: 12px;
}

.chart-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

.chart-empty {
  padding: 20px 0;
}

.chart-content {
  position: relative;
}

.chart-canvas {
  width: 100%;
  overflow: hidden;
}

.chart-svg {
  width: 100%;
  height: auto;
}

.axis-label {
  font-size: 10px;
  fill: #969799;
}

.data-points circle {
  cursor: pointer;
  transition: r 0.2s;
}

.data-points circle:hover {
  r: 6;
}

.bars rect {
  cursor: pointer;
  transition: opacity 0.2s;
}

.bars rect:hover {
  opacity: 0.8;
}

.chart-legend {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #969799;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
  margin-right: 6px;
}

.selected-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
  padding: 8px 12px;
  background: #f7f8fa;
  border-radius: 4px;
}

.info-date {
  font-size: 12px;
  color: #969799;
}

.info-value {
  font-size: 14px;
  font-weight: 600;
  color: #1989fa;
}
</style>
