package handler

import (
	"github.com/ai-fitness-planner/backend/internal/api/request"
	"github.com/ai-fitness-planner/backend/internal/api/response"
	"github.com/ai-fitness-planner/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// StatisticsHandler handles statistics-related HTTP requests
// Requirements: 10.1, 10.2, 10.3, 10.4
type StatisticsHandler struct {
	*BaseHandler
	statsService service.StatisticsService
}

// NewStatisticsHandler creates a new StatisticsHandler instance
func NewStatisticsHandler(statsService service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		BaseHandler:  NewBaseHandler(),
		statsService: statsService,
	}
}

// GetTrainingStatistics handles GET /api/v1/stats/training
// Requirements: 10.1
func (h *StatisticsHandler) GetTrainingStatistics(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var params request.TrainingStatsParams
	if !h.BindQuery(c, &params) {
		return
	}

	stats, err := h.statsService.GetTrainingStatistics(c.Request.Context(), userID, params.Period)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.TrainingStatsResponse{
		Period:            stats.Period,
		StartDate:         stats.StartDate.Format("2006-01-02"),
		EndDate:           stats.EndDate.Format("2006-01-02"),
		TotalWorkouts:     stats.TotalWorkouts,
		TotalDuration:     stats.TotalDuration,
		TotalCalories:     stats.TotalCalories,
		AverageRating:     stats.AverageRating,
		WorkoutsByType:    stats.WorkoutsByType,
		AverageDuration:   stats.AverageDuration,
		HasSufficientData: stats.HasSufficientData,
		Message:           stats.Message,
	}

	h.Success(c, resp)
}

// GetProgressReport handles GET /api/v1/stats/progress
// Requirements: 10.2
func (h *StatisticsHandler) GetProgressReport(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	report, err := h.statsService.GetProgressReport(c.Request.Context(), userID)
	if err != nil {
		h.Error(c, err)
		return
	}

	resp := response.ProgressReportResponse{
		CurrentPeriod: &response.PeriodSummaryInfo{
			StartDate:     report.CurrentPeriod.StartDate.Format("2006-01-02"),
			EndDate:       report.CurrentPeriod.EndDate.Format("2006-01-02"),
			TotalWorkouts: report.CurrentPeriod.TotalWorkouts,
			TotalDuration: report.CurrentPeriod.TotalDuration,
			TotalCalories: report.CurrentPeriod.TotalCalories,
			AverageRating: report.CurrentPeriod.AverageRating,
		},
		PreviousPeriod: &response.PeriodSummaryInfo{
			StartDate:     report.PreviousPeriod.StartDate.Format("2006-01-02"),
			EndDate:       report.PreviousPeriod.EndDate.Format("2006-01-02"),
			TotalWorkouts: report.PreviousPeriod.TotalWorkouts,
			TotalDuration: report.PreviousPeriod.TotalDuration,
			TotalCalories: report.PreviousPeriod.TotalCalories,
			AverageRating: report.PreviousPeriod.AverageRating,
		},
		WorkoutComparison: &response.WorkoutCompareInfo{
			WorkoutCountChange:  report.WorkoutComparison.WorkoutCountChange,
			DurationChange:      report.WorkoutComparison.DurationChange,
			CaloriesChange:      report.WorkoutComparison.CaloriesChange,
			WorkoutCountPercent: report.WorkoutComparison.WorkoutCountPercent,
			DurationPercent:     report.WorkoutComparison.DurationPercent,
			CaloriesPercent:     report.WorkoutComparison.CaloriesPercent,
		},
		HasSufficientData: report.HasSufficientData,
		Message:           report.Message,
	}

	// Add body progress if available
	if report.BodyProgress != nil {
		resp.BodyProgress = &response.BodyProgressInfo{
			CurrentWeight:   report.BodyProgress.CurrentWeight,
			PreviousWeight:  report.BodyProgress.PreviousWeight,
			WeightChange:    report.BodyProgress.WeightChange,
			CurrentBodyFat:  report.BodyProgress.CurrentBodyFat,
			PreviousBodyFat: report.BodyProgress.PreviousBodyFat,
			BodyFatChange:   report.BodyProgress.BodyFatChange,
		}
	}

	h.Success(c, resp)
}

// GetTrends handles GET /api/v1/stats/trends
// Requirements: 10.3
func (h *StatisticsHandler) GetTrends(c *gin.Context) {
	userID, ok := h.GetUserID(c)
	if !ok {
		return
	}

	var params request.TrendsParams
	if !h.BindQuery(c, &params) {
		return
	}

	// Default count to 12 if not provided
	count := params.Count
	if count == 0 {
		count = 12
	}

	trends, err := h.statsService.CalculateTrends(c.Request.Context(), userID, params.Period, count)
	if err != nil {
		h.Error(c, err)
		return
	}

	// Convert data points
	dataPoints := make([]response.TrendPointInfo, 0, len(trends.DataPoints))
	for _, dp := range trends.DataPoints {
		dataPoints = append(dataPoints, response.TrendPointInfo{
			PeriodLabel:   dp.PeriodLabel,
			StartDate:     dp.StartDate.Format("2006-01-02"),
			EndDate:       dp.EndDate.Format("2006-01-02"),
			TotalWorkouts: dp.TotalWorkouts,
			TotalDuration: dp.TotalDuration,
			TotalCalories: dp.TotalCalories,
			AverageRating: dp.AverageRating,
		})
	}

	resp := response.TrendsReportResponse{
		Period:            trends.Period,
		DataPoints:        dataPoints,
		HasSufficientData: trends.HasSufficientData,
		Message:           trends.Message,
	}

	h.Success(c, resp)
}
