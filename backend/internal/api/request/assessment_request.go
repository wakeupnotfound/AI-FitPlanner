package request

// CreateAssessmentRequest represents the request to create a fitness assessment
type CreateAssessmentRequest struct {
	ExperienceLevel       string   `json:"experience_level" binding:"required,oneof=beginner intermediate advanced"`
	WeeklyAvailableDays   int      `json:"weekly_available_days" binding:"required,min=1,max=7"`
	DailyAvailableMinutes int      `json:"daily_available_minutes" binding:"required,min=10,max=480"`
	ActivityType          *string  `json:"activity_type" binding:"omitempty,min=1,max=50"`
	InjuryHistory         *string  `json:"injury_history" binding:"omitempty,max=1000"`
	HealthConditions      *string  `json:"health_conditions" binding:"omitempty,max=1000"`
	PreferredDays         []string `json:"preferred_days" binding:"omitempty,dive,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	EquipmentAvailable    []string `json:"equipment_available" binding:"omitempty,dive,min=1,max=100"`
	AssessmentDate        string   `json:"assessment_date" binding:"required,datetime=2006-01-02,future_date"`
}
