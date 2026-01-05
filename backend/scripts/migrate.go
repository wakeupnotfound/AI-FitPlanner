package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ai-fitness-planner/backend/internal/config"
	"github.com/ai-fitness-planner/backend/internal/pkg/database"
	"github.com/ai-fitness-planner/backend/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Initialize configuration
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Logger.Sync()

	logger.Info("Starting database migration")

	// Initialize database connection
	if err := database.InitDatabase(); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.Close()

	db := database.GetDB()

	// Read schema file
	schemaPath := filepath.Join("database", "schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		logger.Fatal("Failed to read schema file", zap.Error(err), zap.String("path", schemaPath))
	}

	logger.Info("Executing schema migration...")

	// Execute schema SQL
	if err := db.Exec(string(schemaSQL)).Error; err != nil {
		logger.Fatal("Failed to execute schema migration", zap.Error(err))
	}

	logger.Info("Schema migration completed successfully")

	// Insert initial data (prompt templates)
	logger.Info("Inserting initial data...")

	if err := insertPromptTemplates(db); err != nil {
		logger.Fatal("Failed to insert initial data", zap.Error(err))
	}

	logger.Info("Initial data inserted successfully")
	logger.Info("Database migration completed")
}

// insertPromptTemplates inserts default AI prompt templates
func insertPromptTemplates(db *gorm.DB) error {
	templates := []map[string]interface{}{
		{
			"category":    "training",
			"subcategory": "plan_generation",
			"name":        "训练计划生成模板",
			"template": `你是一位专业的健身教练。请根据以下用户信息生成一个详细的训练计划：

用户信息：
- 年龄：{{.Age}}
- 性别：{{.Gender}}
- 身高：{{.Height}}cm
- 体重：{{.Weight}}kg
- 体脂率：{{.BodyFatPercentage}}%
- 经验水平：{{.ExperienceLevel}}
- 每周可训练天数：{{.WeeklyAvailableDays}}天
- 每日可训练时间：{{.DailyAvailableMinutes}}分钟
- 健身目标：{{.FitnessGoals}}
- 伤病历史：{{.InjuryHistory}}
- 可用器材：{{.EquipmentAvailable}}

请生成一个{{.TotalWeeks}}周的训练计划，包含以下内容：
1. 每周的训练安排（包括训练日和休息日）
2. 每天的具体训练项目（动作名称、组数、次数、重量建议、休息时间）
3. 训练强度的渐进安排
4. 安全注意事项

请以JSON格式返回，结构如下：
{
  "weeks": [
    {
      "week": 1,
      "days": [
        {
          "day": 1,
          "date": "2024-01-01",
          "type": "strength",
          "focus_area": "上肢",
          "exercises": [
            {
              "name": "卧推",
              "sets": 4,
              "reps": "8-10",
              "weight": "根据个人能力",
              "rest": "90秒",
              "difficulty": "medium",
              "safety_notes": "保持肩胛骨稳定"
            }
          ],
          "duration": 60,
          "estimated_calories": 300
        }
      ]
    }
  ]
}`,
			"variables":   `["Age","Gender","Height","Weight","BodyFatPercentage","ExperienceLevel","WeeklyAvailableDays","DailyAvailableMinutes","FitnessGoals","InjuryHistory","EquipmentAvailable","TotalWeeks"]`,
			"is_default":  1,
			"description": "用于生成个性化训练计划的默认模板",
		},
		{
			"category":    "nutrition",
			"subcategory": "plan_generation",
			"name":        "饮食计划生成模板",
			"template": `你是一位专业的营养师。请根据以下用户信息生成一个详细的饮食计划：

用户信息：
- 年龄：{{.Age}}
- 性别：{{.Gender}}
- 身高：{{.Height}}cm
- 体重：{{.Weight}}kg
- 活动水平：{{.ActivityLevel}}
- 健身目标：{{.FitnessGoals}}
- 饮食限制：{{.DietaryRestrictions}}
- 饮食偏好：{{.Preferences}}

目标营养摄入：
- 每日卡路里：{{.DailyCalories}}千卡
- 蛋白质比例：{{.ProteinRatio}}%
- 碳水化合物比例：{{.CarbRatio}}%
- 脂肪比例：{{.FatRatio}}%

请生成一个{{.TotalDays}}天的饮食计划，包含以下内容：
1. 每天的三餐和加餐安排
2. 每餐的具体食物、份量和营养成分
3. 烹饪建议
4. 营养补充建议

请以JSON格式返回，结构如下：
{
  "days": [
    {
      "day": 1,
      "date": "2024-01-01",
      "meals": [
        {
          "meal_time": "breakfast",
          "meal_name": "高蛋白早餐",
          "foods": [
            {
              "name": "鸡蛋",
              "amount": "2个",
              "calories": 140,
              "protein": 12,
              "carbs": 1,
              "fat": 10
            }
          ],
          "total_calories": 500,
          "cooking_notes": "水煮或煎蛋"
        }
      ],
      "daily_total": {
        "calories": 2000,
        "protein": 150,
        "carbs": 200,
        "fat": 67
      }
    }
  ]
}`,
			"variables":   `["Age","Gender","Height","Weight","ActivityLevel","FitnessGoals","DietaryRestrictions","Preferences","DailyCalories","ProteinRatio","CarbRatio","FatRatio","TotalDays"]`,
			"is_default":  1,
			"description": "用于生成个性化饮食计划的默认模板",
		},
		{
			"category":    "training",
			"subcategory": "adjustment",
			"name":        "训练计划调整模板",
			"template": `基于用户的反馈，请调整训练计划：

当前计划：
{{.CurrentPlan}}

用户反馈：
- 完成情况：{{.CompletionRate}}%
- 难度评价：{{.DifficultyRating}}/5
- 伤病报告：{{.InjuryReport}}
- 其他反馈：{{.Feedback}}

请提供调整建议，包括：
1. 训练强度调整
2. 动作替换建议
3. 休息时间调整
4. 其他优化建议`,
			"variables":   `["CurrentPlan","CompletionRate","DifficultyRating","InjuryReport","Feedback"]`,
			"is_default":  0,
			"description": "用于根据用户反馈调整训练计划",
		},
		{
			"category":    "nutrition",
			"subcategory": "adjustment",
			"name":        "饮食计划调整模板",
			"template": `基于用户的反馈，请调整饮食计划：

当前计划：
{{.CurrentPlan}}

用户反馈：
- 执行情况：{{.CompletionRate}}%
- 满意度：{{.SatisfactionRating}}/5
- 体重变化：{{.WeightChange}}kg
- 其他反馈：{{.Feedback}}

请提供调整建议，包括：
1. 卡路里摄入调整
2. 营养比例调整
3. 食物替换建议
4. 其他优化建议`,
			"variables":   `["CurrentPlan","CompletionRate","SatisfactionRating","WeightChange","Feedback"]`,
			"is_default":  0,
			"description": "用于根据用户反馈调整饮食计划",
		},
	}

	for _, template := range templates {
		result := db.Exec(`
			INSERT INTO prompt_templates (category, subcategory, name, template, variables, is_default, description)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				template = VALUES(template),
				variables = VALUES(variables),
				description = VALUES(description)
		`,
			template["category"],
			template["subcategory"],
			template["name"],
			template["template"],
			template["variables"],
			template["is_default"],
			template["description"],
		)

		if result.Error != nil {
			return fmt.Errorf("failed to insert template %s: %w", template["name"], result.Error)
		}
	}

	return nil
}
