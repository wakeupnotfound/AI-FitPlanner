-- 创建数据库和用户的SQL脚本
-- 使用方法: mysql -u root -p < database/setup_database.sql

-- 创建数据库
CREATE DATABASE IF NOT EXISTS fitness_planner CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（如果不存在）
CREATE USER IF NOT EXISTS 'fitness_user'@'localhost' IDENTIFIED BY 'fitness_password_2024';

-- 授予权限
GRANT ALL PRIVILEGES ON fitness_planner.* TO 'fitness_user'@'localhost';

-- 刷新权限
FLUSH PRIVILEGES;

-- 显示结果
SELECT 'Database and user created successfully!' AS Status;
SHOW DATABASES LIKE 'fitness_planner';
