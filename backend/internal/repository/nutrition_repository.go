package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ai-fitness-planner/backend/internal/model"
	"gorm.io/gorm"
)

// NutritionPlanRepository defines the interface for nutrition plan operations
type NutritionPlanRepository interface {
	Create(ctx context.Context, plan *model.NutritionPlan) error
	GetByID(ctx context.Context, id int64) (*model.NutritionPlan, error)
	ListByUser(ctx context.Context, userID int64, status string) ([]*model.NutritionPlan, error)
	Update(ctx context.Context, plan *model.NutritionPlan) error
	Delete(ctx context.Context, id int64) error
	GetTodayMeals(ctx context.Context, userID int64, date time.Time) ([]model.NutritionPlanMeal, error)
}

// NutritionRecordRepository defines the interface for nutrition record operations
type NutritionRecordRepository interface {
	Create(ctx context.Context, record *model.NutritionRecord) error
	GetByID(ctx context.Context, id int64) (*model.NutritionRecord, error)
	ListByUser(ctx context.Context, userID int64, startDate, endDate *time.Time) ([]*model.NutritionRecord, error)
	GetDailySummary(ctx context.Context, userID int64, date time.Time) (*DailyNutritionSummary, error)
}

// DailyNutritionSummary represents aggregated nutrition data for a day
type DailyNutritionSummary struct {
	Date          time.Time
	TotalCalories float64
	TotalProtein  float64
	TotalCarbs    float64
	TotalFat      float64
	TotalFiber    float64
	MealCount     int64
}

// nutritionPlanRepository implements NutritionPlanRepository interface
type nutritionPlanRepository struct {
	db *gorm.DB
}

// NewNutritionPlanRepository creates a new instance of NutritionPlanRepository
func NewNutritionPlanRepository(db *gorm.DB) NutritionPlanRepository {
	return &nutritionPlanRepository{db: db}
}

// Create creates a new nutrition plan
func (r *nutritionPlanRepository) Create(ctx context.Context, plan *model.NutritionPlan) error {
	if err := r.db.WithContext(ctx).Create(plan).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a nutrition plan by ID
func (r *nutritionPlanRepository) GetByID(ctx context.Context, id int64) (*model.NutritionPlan, error) {
	var plan model.NutritionPlan
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// ListByUser retrieves all nutrition plans for a user, optionally filtered by status
func (r *nutritionPlanRepository) ListByUser(ctx context.Context, userID int64, status string) ([]*model.NutritionPlan, error) {
	var plans []*model.NutritionPlan
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

// Update updates an existing nutrition plan
func (r *nutritionPlanRepository) Update(ctx context.Context, plan *model.NutritionPlan) error {
	if err := r.db.WithContext(ctx).Save(plan).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a nutrition plan
func (r *nutritionPlanRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.NutritionPlan{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetTodayMeals retrieves the meal plan for a specific date
func (r *nutritionPlanRepository) GetTodayMeals(ctx context.Context, userID int64, date time.Time) ([]model.NutritionPlanMeal, error) {
	var plan model.NutritionPlan

	// Format date as string for comparison (YYYY-MM-DD)
	dateStr := date.Format("2006-01-02")

	// Find active plan that includes this date
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND start_date <= ? AND end_date >= ?",
			userID, "active", date, date).
		First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Parse plan data to find meals for the specific date
	mealsData, ok := plan.PlanData["meals"]
	if !ok {
		return nil, nil
	}

	mealsInterface, ok := mealsData.([]interface{})
	if !ok {
		return nil, nil
	}

	var meals []model.NutritionPlanMeal

	for _, mealInterface := range mealsInterface {
		mealMap, ok := mealInterface.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if this meal is for the requested date
		mealDate, ok := mealMap["date"].(string)
		if !ok || mealDate != dateStr {
			continue
		}

		meal := model.NutritionPlanMeal{}

		if time, ok := mealMap["time"].(string); ok {
			meal.Time = time
		}

		if totalCal, ok := mealMap["total_calories"].(float64); ok {
			meal.TotalCalories = totalCal
		}

		// Parse foods
		if foodsInterface, ok := mealMap["foods"].([]interface{}); ok {
			foods := make([]model.NutritionFoodItem, 0, len(foodsInterface))
			for _, foodInterface := range foodsInterface {
				foodMap, ok := foodInterface.(map[string]interface{})
				if !ok {
					continue
				}

				food := model.NutritionFoodItem{}
				if name, ok := foodMap["name"].(string); ok {
					food.Name = name
				}
				if amount, ok := foodMap["amount"].(string); ok {
					food.Amount = amount
				}
				if calories, ok := foodMap["calories"].(float64); ok {
					food.Calories = calories
				}
				if protein, ok := foodMap["protein"].(float64); ok {
					food.Protein = protein
				}
				if carbs, ok := foodMap["carbs"].(float64); ok {
					food.Carbs = carbs
				}
				if fat, ok := foodMap["fat"].(float64); ok {
					food.Fat = fat
				}
				if fiber, ok := foodMap["fiber"].(float64); ok {
					food.Fiber = fiber
				}

				foods = append(foods, food)
			}
			meal.Foods = foods
		}

		meals = append(meals, meal)
	}

	return meals, nil
}

// nutritionRecordRepository implements NutritionRecordRepository interface
type nutritionRecordRepository struct {
	db *gorm.DB
}

// NewNutritionRecordRepository creates a new instance of NutritionRecordRepository
func NewNutritionRecordRepository(db *gorm.DB) NutritionRecordRepository {
	return &nutritionRecordRepository{db: db}
}

// Create creates a new nutrition record
func (r *nutritionRecordRepository) Create(ctx context.Context, record *model.NutritionRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a nutrition record by ID
func (r *nutritionRecordRepository) GetByID(ctx context.Context, id int64) (*model.NutritionRecord, error) {
	var record model.NutritionRecord
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// ListByUser retrieves nutrition records for a user within an optional date range
func (r *nutritionRecordRepository) ListByUser(ctx context.Context, userID int64, startDate, endDate *time.Time) ([]*model.NutritionRecord, error) {
	var records []*model.NutritionRecord
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("meal_date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("meal_date <= ?", *endDate)
	}

	if err := query.Order("meal_date DESC, created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// GetDailySummary calculates aggregated nutrition data for a specific day
func (r *nutritionRecordRepository) GetDailySummary(ctx context.Context, userID int64, date time.Time) (*DailyNutritionSummary, error) {
	summary := &DailyNutritionSummary{
		Date: date,
	}

	// Get count of meals
	if err := r.db.WithContext(ctx).
		Model(&model.NutritionRecord{}).
		Where("user_id = ? AND meal_date = ?", userID, date).
		Count(&summary.MealCount).Error; err != nil {
		return nil, err
	}

	// Get sum of all nutrition values
	type AggregateResult struct {
		TotalCalories float64
		TotalProtein  float64
		TotalCarbs    float64
		TotalFat      float64
		TotalFiber    float64
	}

	var result AggregateResult
	if err := r.db.WithContext(ctx).
		Model(&model.NutritionRecord{}).
		Select(`
			COALESCE(SUM(calories), 0) as total_calories,
			COALESCE(SUM(protein), 0) as total_protein,
			COALESCE(SUM(carbs), 0) as total_carbs,
			COALESCE(SUM(fat), 0) as total_fat,
			COALESCE(SUM(fiber), 0) as total_fiber
		`).
		Where("user_id = ? AND meal_date = ?", userID, date).
		Scan(&result).Error; err != nil {
		return nil, err
	}

	summary.TotalCalories = result.TotalCalories
	summary.TotalProtein = result.TotalProtein
	summary.TotalCarbs = result.TotalCarbs
	summary.TotalFat = result.TotalFat
	summary.TotalFiber = result.TotalFiber

	return summary, nil
}
