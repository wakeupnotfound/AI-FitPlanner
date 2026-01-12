-- 用户基础信息表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    nickname VARCHAR(50) COMMENT '昵称',
    email VARCHAR(100) UNIQUE NOT NULL COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    avatar MEDIUMTEXT COMMENT '头像URL/Base64',
    status TINYINT DEFAULT 1 COMMENT '1-正常, 0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户基础表';

-- AI API配置表
CREATE TABLE ai_apis (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '所属用户ID',
    provider VARCHAR(50) NOT NULL COMMENT '服务提供商',
    name VARCHAR(100) NOT NULL COMMENT '自定义名称',
    api_endpoint VARCHAR(500) NOT NULL COMMENT 'API地址',
    api_key_encrypted TEXT NOT NULL COMMENT '加密的API Key',
    model VARCHAR(100) COMMENT '使用的模型',
    max_tokens INT COMMENT '最大token数',
    temperature DECIMAL(3,2) DEFAULT 0.7 COMMENT '生成温度',
    is_default TINYINT DEFAULT 0 COMMENT '是否默认使用',
    status TINYINT DEFAULT 1 COMMENT '1-启用, 0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_provider (provider),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI API配置表';

-- 用户身体数据表
CREATE TABLE user_body_data (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    age INT NOT NULL COMMENT '年龄',
    gender ENUM('male', 'female', 'other') NOT NULL COMMENT '性别',
    height DECIMAL(5,2) NOT NULL COMMENT '身高(cm)',
    weight DECIMAL(5,2) NOT NULL COMMENT '体重(kg)',
    body_fat_percentage DECIMAL(4,2) COMMENT '体脂率',
    muscle_percentage DECIMAL(4,2) COMMENT '肌肉率',
    measurement_date DATE NOT NULL COMMENT '测量日期',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_date (user_id, measurement_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身体数据表';

-- 健身目标表
CREATE TABLE fitness_goals (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    goal_type VARCHAR(100) NOT NULL COMMENT '目标类型',
    goal_description TEXT COMMENT '目标描述',
    initial_weight DECIMAL(5,2) COMMENT '初始体重',
    initial_body_fat DECIMAL(4,2) COMMENT '初始体脂',
    initial_muscle_mass DECIMAL(4,2) COMMENT '初始肌肉量',
    target_weight DECIMAL(5,2) COMMENT '目标体重',
    deadline DATE COMMENT '截止日期',
    priority INT DEFAULT 1 COMMENT '优先级',
    status VARCHAR(20) DEFAULT 'active' COMMENT 'active/completed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_status (user_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='健身目标表';

-- 运动能力评估表
CREATE TABLE fitness_assessments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    experience_level ENUM('beginner', 'intermediate', 'advanced') NOT NULL COMMENT '经验水平',
    weekly_available_days INT NOT NULL COMMENT '每周可用天数',
    daily_available_minutes INT NOT NULL COMMENT '每日可用分钟数',
    activity_type VARCHAR(50) COMMENT '主要运动类型',
    injury_history TEXT COMMENT '伤病历史',
    health_conditions TEXT COMMENT '健康问题',
    preferred_days JSON COMMENT '偏好的训练日',
    equipment_available JSON COMMENT '可用的器材',
    assessment_date DATE NOT NULL COMMENT '评估日期',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_date (user_id, assessment_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='运动能力评估表';

-- 训练计划表
CREATE TABLE training_plans (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    plan_name VARCHAR(200) NOT NULL COMMENT '计划名称',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE NOT NULL COMMENT '结束日期',
    total_weeks INT NOT NULL COMMENT '总周数',
    difficulty_level ENUM('easy', 'medium', 'hard', 'extreme') COMMENT '难度等级',
    training_purpose VARCHAR(100) COMMENT '训练目的',
    ai_api_id BIGINT NOT NULL COMMENT '使用的AI API',
    plan_data JSON NOT NULL COMMENT '计划详细数据',
    status VARCHAR(20) DEFAULT 'active' COMMENT 'active/inactive/completed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (ai_api_id) REFERENCES ai_apis(id),
    INDEX idx_user_status (user_id, status),
    INDEX idx_start_date (start_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练计划表';

-- 饮食计划表
CREATE TABLE nutrition_plans (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    plan_name VARCHAR(200) NOT NULL COMMENT '计划名称',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE NOT NULL COMMENT '结束日期',
    daily_calories DECIMAL(7,2) COMMENT '每日卡路里',
    protein_ratio DECIMAL(3,2) COMMENT '蛋白质比例',
    carb_ratio DECIMAL(3,2) COMMENT '碳水化合物比例',
    fat_ratio DECIMAL(3,2) COMMENT '脂肪比例',
    dietary_restrictions JSON COMMENT '饮食限制',
    preferences JSON COMMENT '饮食偏好',
    plan_data JSON NOT NULL COMMENT '计划详细数据',
    ai_api_id BIGINT NOT NULL COMMENT '使用的AI API',
    status VARCHAR(20) DEFAULT 'active' COMMENT 'active/inactive/completed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (ai_api_id) REFERENCES ai_apis(id),
    INDEX idx_user_status (user_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='营养计划表';

-- 训练记录表
CREATE TABLE training_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    plan_id BIGINT COMMENT '所属计划ID',
    workout_date DATE NOT NULL COMMENT '训练日期',
    workout_type VARCHAR(100) NOT NULL COMMENT '训练类型',
    duration_minutes INT COMMENT '训练时长',
    exercises JSON COMMENT '训练项目',
    performance_data JSON COMMENT '表现数据',
    notes TEXT COMMENT '备注',
    rating INT COMMENT '自我评分1-5',
    injury_report TEXT COMMENT '伤病报告',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (plan_id) REFERENCES training_plans(id) ON DELETE SET NULL,
    INDEX idx_user_date (user_id, workout_date),
    INDEX idx_plan_id (plan_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练记录表';

-- 饮食记录表
CREATE TABLE nutrition_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    meal_date DATE NOT NULL COMMENT '用餐日期',
    meal_time ENUM('breakfast', 'lunch', 'dinner', 'snack') COMMENT '用餐时间',
    foods JSON NOT NULL COMMENT '食物详情',
    calories DECIMAL(7,2) COMMENT '卡路里',
    protein DECIMAL(6,2) COMMENT '蛋白质(g)',
    carbs DECIMAL(6,2) COMMENT '碳水化合物(g)',
    fat DECIMAL(6,2) COMMENT '脂肪(g)',
    fiber DECIMAL(6,2) COMMENT '纤维(g)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_date (user_id, meal_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='饮食记录表';

-- AI提示词模板表
CREATE TABLE prompt_templates (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    category VARCHAR(50) NOT NULL COMMENT '分类',
    subcategory VARCHAR(50) COMMENT '子分类',
    name VARCHAR(200) NOT NULL COMMENT '模板名称',
    template TEXT NOT NULL COMMENT '提示词模板',
    variables JSON COMMENT '变量列表',
    is_default TINYINT DEFAULT 0 COMMENT '是否默认模板',
    description TEXT COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_default (is_default)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI提示词模板表';

-- 反馈记录表
CREATE TABLE feedback_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    plan_type ENUM('training', 'nutrition') COMMENT '计划类型',
    plan_id BIGINT COMMENT '计划ID',
    feedback_type VARCHAR(50) COMMENT '反馈类型',
    feedback_data JSON COMMENT '反馈数据',
    satisfaction INT COMMENT '满意度1-5',
    ai_response TEXT COMMENT 'AI调整建议',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_date (user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户反馈表';
