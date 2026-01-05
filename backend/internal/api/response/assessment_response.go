package response

// AssessmentInfo represents a fitness assessment in responses
type AssessmentInfo struct {
	ID                    int64    `json:"id"`
	ExperienceLevel       string   `json:"experience_level"`
	WeeklyAvailableDays   int      `json:"weekly_available_days"`
	DailyAvailableMinutes int      `json:"daily_available_minutes"`
	ActivityType          string   `json:"activity_type,omitempty"`
	InjuryHistory         string   `json:"injury_history,omitempty"`
	HealthConditions      string   `json:"health_conditions,omitempty"`
	PreferredDays         []string `json:"preferred_days,omitempty"`
	EquipmentAvailable    []string `json:"equipment_available,omitempty"`
	AssessmentDate        string   `json:"assessment_date"`
	CreatedAt             string   `json:"created_at"`
}

// AssessmentDetailResponse represents a single assessment response
type AssessmentDetailResponse struct {
	Assessment AssessmentInfo `json:"assessment"`
}
