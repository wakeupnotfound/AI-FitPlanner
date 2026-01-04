# 前端架构设计文档

## 一、技术栈

```json
{
  "framework": "Vue 3.3+",
  "language": "JavaScript (ES2020)",
  "buildTool": "Vite 4.x",
  "stateManagement": "Pinia 2.x",
  "router": "Vue Router 4.x",
  "uiFramework": "Van 4.x",
  "httpClient": "Axios 1.x",
  "mobileUI": "Flexible Rem + Viewport",
  "icon": "Iconify",
  "dateTime": "Day.js",
  "validator": "VeeValidate 4.x",
  "linter": "ESLint",
  "formatter": "Prettier"
}
```

---

## 二、项目目录结构

```
src/
├── assets/                     # 静态资源
│   ├── images/                # 图片
│   ├── styles/                # 全局样式
│   │   ├── variables.css      # CSS变量
│   │   ├── mixins.css         # 混合宏
│   │   └── global.css         # 全局样式
│   └── icons/                 # SVG图标
│
├── components/                 # 通用组件
│   ├── common/                # 基础组件
│   │   ├── AppHeader.vue      # 页面头部
│   │   ├── AppFooter.vue      # 页面底部
│   │   ├── Loading.vue        # 加载动画
│   │   ├── EmptyState.vue     # 空状态
│   │   └── ErrorBoundary.vue  # 错误边界
│   │
│   ├── form/                  # 表单组件
│   │   ├── InputField.vue     # 输入框
│   │   ├── SelectField.vue    # 下拉框
│   │   ├── DatePicker.vue     # 日期选择器
│   │   ├── NumberStep.vue     # 数字步进器
│   │   └── FormSection.vue    # 表单区块
│   │
│   ├── fitness/               # 健身专用组件
│   │   ├── BodyMetricsCard.vue     # 身体数据卡片
│   │   ├── TrainingPlanCard.vue    # 训练计划卡片
│   │   ├── ExerciseItem.vue        # 训练项目
│   │   ├── NutritionPlanCard.vue   # 饮食计划卡片
│   │   ├── MealItem.vue            # 餐食项目
│   │   └── ProgressChart.vue       # 进度图表
│   │
│   └── ai/                    # AI相关组件
│       ├── ApiConfigForm.vue  # API配置表单
│       ├── ApiCard.vue        # API信息卡片
│       └── ModelSelector.vue  # 模型选择器
│
├── composables/                # 组合式函数
│   ├── useAuth.js             # 认证相关
│   ├── useApi.js              # API调用
│   ├── useForm.js             # 表单处理
│   ├── useTraining.js         # 训练数据
│   ├── useNutrition.js        # 营养数据
│   └── useAi.js               # AI功能
│
├── views/                      # 页面视图
│   ├── auth/                  # 认证相关
│   │   ├── Login.vue          # 登录页
│   │   └── Register.vue       # 注册页
│   │
   ├── dashboard/               # 首页仪表盘
│   │   └── Dashboard.vue      # 主仪表盘
│   │
│   ├── profile/               # 个人档案
│   │   ├── Profile.vue        # 个人主页
│   │   ├── BodyData.vue       # 身体数据
│   │   └── FitnessGoals.vue   # 健身目标
│   │
│   ├── assessment/            # 能力评估
│   │   └── FitnessAssessment.vue  # 运动能力评估
│   │
│   ├── api/                   # AI配置
│   │   ├── ApiConfig.vue      # API管理
│   │   └── ApiAdd.vue         # 添加API
│   │
│   ├── training/              # 训练计划
│   │   ├── TrainingPlans.vue  # 计划列表
│   │   ├── PlanDetail.vue     # 计划详情
│   │   └── TodayTraining.vue  # 今日训练
│   │
│   ├── nutrition/             # 饮食计划
│   │   ├── NutritionPlans.vue # 计划列表
│   │   ├── PlanDetail.vue     # 计划详情
│   │   └── TodayNutrition.vue # 今日饮食
│   │
│   └── record/                # 记录
│       ├── TrainingRecord.vue # 训练记录
│       └── NutritionRecord.vue# 饮食记录
│
├── stores/                     # 状态管理
│   ├── index.js               # Store入口
│   ├── auth.js                # 认证Store
│   ├── user.js                # 用户Store
│   ├── training.js            # 训练Store
│   ├── nutrition.js           # 营养Store
│   └── ai.js                  # AI配置Store
│
├── router/                     # 路由配置
│   ├── index.js               # 路由入口
│   ├── routes.js              # 路由定义
│   └── guards.js              # 路由守卫
│
├── plugins/                    # 插件
│   ├── vant.js                # Vant组件注册
│   ├── axios.js               # Axios配置
│   └── dayjs.js               # Day.js配置
│
├── utils/                      # 工具函数
│   ├── request.js             # 请求封装
│   ├── auth.js                # 认证工具
│   ├── date.js                # 日期处理
│   ├── validation.js          # 验证规则
│   ├── storage.js             # 本地存储
│   └── constants.js           # 常量定义
│
├── assets/                     # 静态资源
│
├── App.vue                     # 根组件
└── main.js                     # 入口文件
```

---

## 三、状态管理设计 (Pinia)

### 1. auth.js - 认证状态
```javascript
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token'),
    user: null,
    isAuthenticated: false
  }),

  getters: {
    userInfo: (state) => state.user
  },

  actions: {
    async login(credentials) {
      // 登录逻辑
    },
    logout() {
      // 登出逻辑
    }
  }
})
```

### 2. user.js - 用户信息
```javascript
export const useUserStore = defineStore('user', {
  state: () => ({
    profile: null,
    bodyData: {},
    fitnessGoals: [],
    assessment: null
  }),

  actions: {
    async fetchProfile() {
      // 获取用户档案
    },
    async updateBodyData(data) {
      // 更新身体数据
    }
  }
})
```

### 3. training.js - 训练数据
```javascript
export const useTrainingStore = defineStore('training', {
  state: () => ({
    plans: [],
    currentPlan: null,
    todaySchedule: null,
    records: []
  }),

  actions: {
    async generatePlan(params) {
      // 生成训练计划
    },
    async fetchTodayTraining() {
      // 获取今日训练
    },
    async recordTraining(data) {
      // 记录训练
    }
  }
})
```

### 4. ai.js - AI配置
```javascript
export const useAiStore = defineStore('ai', {
  state: () => ({
    apis: [],
    defaultApi: null,
    templates: [],
    models: {
      openai: ['gpt-3.5-turbo', 'gpt-4'],
      wenxin: ['ERNIE-Bot'],
      tongyi: ['qwen-max', 'qwen-turbo']
    }
  }),

  actions: {
    async fetchApis() {
      // 获取API列表
    },
    async addApi(config) {
      // 添加API配置
    },
    async testApi(apiId) {
      // 测试API连通性
    }
  }
})
```

---

## 四、核心组件设计

### 1. 仪表盘组件 (Dashboard.vue)

**功能模块:**
- 今日训练概览（今天是否有训练）
- 本周出勤统计
- 身体数据变化趋势
- 目标完成进度
- 快速操作入口

**布局设计 (移动端):**
```
┌─────────────────────────┐
│  用户头像 + 今日状态     │
├─────────────────────────┤
│  今日训练/饮食          │
│  [开始训练] [记录饮食]   │
├─────────────────────────┤
│  本周统计               │
│  ○○○○○ ○○○○○ ○○○○○    │
├─────────────────────────┤
│  目标进度               │
│  ━━━━━━━━━ 70%         │
├─────────────────────────┤
│  快速访问               │
│  [计划] [记录] [分析]    │
└─────────────────────────┘
```

**关键代码示例:**
```vue
<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <UserInfoCard />
    </header>

    <section class="today-section">
      <TodaySchedule />
    </section>

    <section class="stats-section">
      <WeekProgress />
      <GoalProgress />
    </section>

    <QuickActions />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useTrainingStore } from '@/stores/training'
import { useUserStore } from '@/stores/user'

const trainingStore = useTrainingStore()
const userStore = useUserStore()

onMounted(async () => {
  await Promise.all([
    trainingStore.fetchTodayTraining(),
    userStore.fetchProfile()
  ])
})
</script>
```

---

### 2. AI配置管理组件 (ApiConfig.vue)

**功能模块:**
- API列表展示
- 添加/编辑API
- 测试API连通性
- 设置默认API
- 安全删除API

**交互设计:**
```
┌─────────────────────────┐
│  AI服务配置             │ ◄ 顶部标题
├─────────────────────────┤
│  + 添加新API [管理模板]  │
├─────────────────────────┤
│  ┌───────────────────┐ │ ◄ API卡片1
│  │  OpenAI GPT-4     │ │
│  │  • 模型: gpt-4    │ │
│  │  [测试] [编辑]    │ │
│  └───────────────────┘ │
│  ┌───────────────────┐ │ ◄ API卡片2
│  │  文心一言         │ │
│  │  • 默认使用       │ │
│  │  [测试] [编辑]    │ │
│  └───────────────────┘ │
└─────────────────────────┘
```

---

### 3. 训练计划生成组件 (TrainingPlans.vue)

**功能模块:**
- 计划列表展示
- 新建计划向导
- 计划详情查看
- 开始/暂停计划
- 计划执行进度

**新建计划向导流程:**
```
步骤1: 基本信息
┌─────────────────────────┐
│  新建训练计划            │
│                         │
│  计划名称: [          ]  │
│  目标: ○增肌 ○减脂 ○   │
│  周期: [4] 周           │
│                         │
│      [下一步]           │
└─────────────────────────┘

步骤2: 能力评估 (如未完成)
┌─────────────────────────┐
│  运动能力评估            │
│    ━━━━━━━━━ 60%        │
│  请先完成运动能力评估... │
│  [立即前往评估]         │
└─────────────────────────┘

步骤3: AI配置
┌─────────────────────────┐
│  选择AI服务              │
│                         │
│  ⊙ OpenAI GPT-4        │
│  ○ 文心一言            │
│  ○ 通义千问            │
│                         │
│     [开始生成] [取消]    │
└─────────────────────────┘
```

---

### 4. 今日训练执行组件 (TodayTraining.vue)

**功能模块:**
- 今日训练计划展示
- 训练计时器
- 记录训练完成度
- 受伤风险警示
- 训练后评价

**界面布局:**
```
┌─────────────────────────┐
│  今日训练：下肢训练       │
│  预计60分钟 · 350卡路里 │
├─────────────────────────┤
│  总体进度 ★★★☆☆ 60%     │
├─────────────────────────┤
│  □ 深蹲                 │
│     4组 × 12次 · 70kg   │
│     [开始] [跳过]       │
│                         │
│  ○ 硬拉                 │
│     3组 × 10次 · 80kg   │
│     [开始]              │
├─────────────────────────┤
│  当前项目计时: 00:03:45  │
│  [完成本轮] [暂停]       │
└─────────────────────────┘
```

---

## 五、移动端适配策略

### 1. 响应式布局
```css
/* 全局变量 */
:root {
  --px-1: 0.0267rem;  /* 1px = 0.0267rem */
  --primary-color: #1989fa;
  --text-color: #323233;
  --border-color: #ebedf0;

  /* 间距 */
  --padding-xs: 0.2133rem;  /* 8px */
  --padding-sm: 0.4267rem;  /* 16px */
  --padding-md: 0.64rem;    /* 24px */
  --padding-lg: 0.8533rem;  /* 32px */
}

/* 容器响应式 */
.page-container {
  width: 100%;
  max-width: 480px; /* 限制最大宽度 */
  margin: 0 auto;
  padding: var(--padding-sm);
}
```

### 2. Touch交互优化
```vue
<template>
  <!-- 触摸反馈 -->
  <van-button
    :loading="loading"
    :disabled="disabled"
    @touchstart="handleTouchStart"
    @click="handleClick"
  >
    {{ buttonText }}
  </van-button>
</template>

<script setup>
const handleTouchStart = () => {
  // 触摸反馈动画
  if ('vibrate' in navigator) {
    navigator.vibrate(50); // 轻微震动
  }
}
</script>
```

### 3. 手势滑动支持
```vue
<script setup>
import { ref } from 'vue'

const startX = ref(0)
const endX = ref(0)

const onTouchStart = (e) => {
  startX.value = e.touches[0].clientX
}

const onTouchEnd = (e) => {
  endX.value = e.changedTouches[0].clientX
  const diff = endX.value - startX.value

  if (Math.abs(diff) > 100) {
    if (diff > 0) {
      // 右滑 - 返回上一页
      router.back()
    } else {
      // 左滑 - 打开侧边栏
      showSidebar()
    }
  }
}
</script>
```

---

## 六、API请求封装

### 1. Axios配置 (utils/request.js)
```javascript
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { showNotify } from 'vant'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // Token过期，跳转到登录
      const authStore = useAuthStore()
      authStore.logout()
    }

    const message = error.response?.data?.message || '请求失败'
    showNotify({ type: 'danger', message })

    return Promise.reject(error)
  }
)

export default request
```

### 2. API调用示例 (composables/useTraining.js)
```javascript
import { ref } from 'vue'
import request from '@/utils/request'

export function useTraining() {
  const loading = ref(false)
  const error = ref(null)

  // 生成训练计划
  const generatePlan = async (params) => {
    loading.value = true
    error.value = null

    try {
      const response = await request.post('/api/training/generate', params)
      return response.data
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  // 获取今日训练
  const fetchTodayTraining = async () => {
    const response = await request.get('/api/training/today')
    return response.data
  }

  // 记录训练
  const recordTraining = async (data) => {
    const response = await request.post('/api/training/record', data)
    return response.data
  }

  return {
    loading,
    error,
    generatePlan,
    fetchTodayTraining,
    recordTraining
  }
}
```

---

## 七、性能优化

### 1. 懒加载路由
```javascript
const routes = [
  {
    path: '/training',
    component: () => import('@/views/training/TrainingPlans.vue'),
    meta: { requiresAuth: true }
  }
]
```

### 2. 组件缓存
```vue
<keep-alive :include="cacheViews">
  <router-view />
</keep-alive>

<script setup>
const cacheViews = ['Dashboard', 'TrainingPlans', 'NutritionPlans']
</script>
```

### 3. 图片懒加载
```vue
<img v-lazy="imageUrl" alt="描述" />
```

### 4. 数据缓存
```javascript
// 使用store缓存
store.$state = {
  trainingPlans: [],
  cachedAt: null
}

const CACHE_DURATION = 5 * 60 * 1000 // 5分钟

const fetchWithCache = async () => {
  if (store.trainingPlans.length &&
      Date.now() - store.cachedAt < CACHE_DURATION) {
    return store.trainingPlans
  }

  const data = await fetchTrainingPlans()
  store.trainingPlans = data
  store.cachedAt = Date.now()
  return data
}
```

---

## 八、开发规范

### 1. 命名规范
- 组件名：大驼峰 (Dashboard.vue)
- 文件名：小驼峰 (apiConfig.js)
- Store：use + 功能名 (useUserStore)
- 组合式函数：use + 功能名 (useTraining)

### 2. 代码风格
- 使用Composition API
- 组件小于300行
- 单一职责原则
- 优先使用依赖注入而非props drilling

### 3. Git提交规范
```
feat: 新增功能
fix: 修复bug
docs: 文档修改
style: 样式调整
refactor: 重构代码
test: 测试相关
chore: 构建或工具相关
```

---

## 九、构建配置

### Vite配置 (vite.config.js)
```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],

  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },

  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },

  build: {
    target: 'es2015',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    rollupOptions: {
      output: {
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: '[ext]/[name]-[hash].[ext]'
      }
    }
  }
})
```

### PostCSS配置
```javascript
module.exports = {
  plugins: {
    'postcss-pxtorem': {
      rootValue: 37.5, // 设计稿宽度375px
      propList: ['*'],
      selectorBlackList: ['.norem'] // 过滤.norem开头的class
    },
    autoprefixer: {}
  }
}
```
