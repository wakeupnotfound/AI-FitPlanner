package service

import (
	"context"
	"time"

	"github.com/ai-fitness-planner/backend/internal/errors"
	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/repository"
)

// UpdateProfileRequest represents the profile update request data
type UpdateProfileRequest struct {
	Email  *string `json:"email" validate:"omitempty,email,max=100"`
	Nickname *string `json:"nickname" validate:"omitempty,min=1,max=50"`
	Phone  *string `json:"phone" validate:"omitempty,max=20"`
	Avatar *string `json:"avatar" validate:"omitempty,avatar"`
}

// BodyDataRequest represents the body data submission request
type BodyDataRequest struct {
	Age               int       `json:"age" validate:"required,min=1,max=150"`
	Gender            string    `json:"gender" validate:"required,oneof=male female other"`
	Height            float64   `json:"height" validate:"required,min=50,max=300"`
	Weight            float64   `json:"weight" validate:"required,min=20,max=500"`
	BodyFatPercentage *float64  `json:"body_fat_percentage" validate:"omitempty,min=0,max=100"`
	MusclePercentage  *float64  `json:"muscle_percentage" validate:"omitempty,min=0,max=100"`
	MeasurementDate   time.Time `json:"measurement_date" validate:"required"`
}

// FitnessGoalRequest represents the fitness goal submission request
type FitnessGoalRequest struct {
	GoalType        string     `json:"goal_type" validate:"required,max=100"`
	GoalDescription *string    `json:"goal_description"`
	TargetWeight    *float64   `json:"target_weight" validate:"omitempty,min=20,max=500"`
	Deadline        *time.Time `json:"deadline"`
	Priority        int        `json:"priority" validate:"min=1,max=10"`
}

// UserService interface defines methods for user profile operations
type UserService interface {
	GetProfile(ctx context.Context, userID int64) (*model.User, error)
	UpdateProfile(ctx context.Context, userID int64, req *UpdateProfileRequest) (*model.User, error)
	AddBodyData(ctx context.Context, userID int64, req *BodyDataRequest) (*model.UserBodyData, error)
	GetBodyDataHistory(ctx context.Context, userID int64) ([]*model.UserBodyData, error)
	SetFitnessGoals(ctx context.Context, userID int64, req *FitnessGoalRequest) (*model.FitnessGoal, error)
	GetFitnessGoals(ctx context.Context, userID int64) ([]*model.FitnessGoal, error)
	UpdateFitnessGoals(ctx context.Context, userID int64, goalID int64, req *FitnessGoalRequest) (*model.FitnessGoal, error)
}

// userService implements the UserService interface
type userService struct {
	userRepo        repository.UserRepository
	bodyDataRepo    repository.BodyDataRepository
	fitnessGoalRepo repository.FitnessGoalRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(
	userRepo repository.UserRepository,
	bodyDataRepo repository.BodyDataRepository,
	fitnessGoalRepo repository.FitnessGoalRepository,
) UserService {
	return &userService{
		userRepo:        userRepo,
		bodyDataRepo:    bodyDataRepo,
		fitnessGoalRepo: fitnessGoalRepo,
	}
}

// GetProfile retrieves a user's profile information
// Validates: Requirements 2.1
func (s *userService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get user profile")
	}
	if user == nil {
		return nil, errors.ErrResourceNotFound
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return user, nil
}

// UpdateProfile updates a user's profile information with validation
// Validates: Requirements 2.2
func (s *userService) UpdateProfile(ctx context.Context, userID int64, req *UpdateProfileRequest) (*model.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get user")
	}
	if user == nil {
		return nil, errors.ErrResourceNotFound
	}

	// Update fields if provided
	if req.Email != nil {
		// Check if email is already taken by another user
		existingUser, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrDatabase, "failed to check email")
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.ErrEmailExists
		}
		user.Email = *req.Email
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if req.Nickname != nil {
		user.Nickname = req.Nickname
	}

	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}

	user.UpdatedAt = time.Now()

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update user profile")
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return user, nil
}

// AddBodyData stores a new body measurement with timestamp
// Validates: Requirements 2.3
func (s *userService) AddBodyData(ctx context.Context, userID int64, req *BodyDataRequest) (*model.UserBodyData, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get user")
	}
	if user == nil {
		return nil, errors.ErrResourceNotFound
	}

	// Create body data record
	bodyData := &model.UserBodyData{
		UserID:            userID,
		Age:               req.Age,
		Gender:            req.Gender,
		Height:            req.Height,
		Weight:            req.Weight,
		BodyFatPercentage: req.BodyFatPercentage,
		MusclePercentage:  req.MusclePercentage,
		MeasurementDate:   req.MeasurementDate,
		CreatedAt:         time.Now(),
	}

	if err := s.bodyDataRepo.Create(ctx, bodyData); err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to add body data")
	}

	return bodyData, nil
}

// GetBodyDataHistory retrieves body measurements ordered by date
// Validates: Requirements 2.4
func (s *userService) GetBodyDataHistory(ctx context.Context, userID int64) ([]*model.UserBodyData, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get user")
	}
	if user == nil {
		return nil, errors.ErrResourceNotFound
	}

	// Get body data history ordered by measurement date descending
	bodyDataList, err := s.bodyDataRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get body data history")
	}

	return bodyDataList, nil
}

// SetFitnessGoals validates and stores fitness goals with priority levels
// Validates: Requirements 2.5
func (s *userService) SetFitnessGoals(ctx context.Context, userID int64, req *FitnessGoalRequest) (*model.FitnessGoal, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get user")
	}
	if user == nil {
		return nil, errors.ErrResourceNotFound
	}

	var initialWeight *float64
	var initialBodyFat *float64
	var initialMuscle *float64
	if latestBodyData, err := s.bodyDataRepo.GetLatestByUserID(ctx, userID); err == nil && latestBodyData != nil {
		initialWeight = &latestBodyData.Weight
		if latestBodyData.BodyFatPercentage != nil {
			initialBodyFat = latestBodyData.BodyFatPercentage
		}
		if latestBodyData.MusclePercentage != nil {
			initialMuscle = latestBodyData.MusclePercentage
		}
	}

	// Create fitness goal
	goal := &model.FitnessGoal{
		UserID:          userID,
		GoalType:        req.GoalType,
		GoalDescription: req.GoalDescription,
		InitialWeight:   initialWeight,
		InitialBodyFat:  initialBodyFat,
		InitialMuscle:   initialMuscle,
		TargetWeight:    req.TargetWeight,
		Deadline:        req.Deadline,
		Priority:        req.Priority,
		Status:          string(model.GoalStatusActive),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.fitnessGoalRepo.Create(ctx, goal); err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to set fitness goal")
	}

	return goal, nil
}


// GetFitnessGoals retrieves all active fitness goals for a user
// Validates: Requirements 2.5
func (s *userService) GetFitnessGoals(ctx context.Context, userID int64) ([]*model.FitnessGoal, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get user")
	}
	if user == nil {
		return nil, errors.ErrResourceNotFound
	}

	// Get all fitness goals (empty status means all)
	goals, err := s.fitnessGoalRepo.GetByUserID(ctx, userID, "")
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get fitness goals")
	}

	return goals, nil
}

// UpdateFitnessGoals updates an existing fitness goal
// Validates: Requirements 2.5
func (s *userService) UpdateFitnessGoals(ctx context.Context, userID int64, goalID int64, req *FitnessGoalRequest) (*model.FitnessGoal, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get user")
	}
	if user == nil {
		return nil, errors.ErrResourceNotFound
	}

	// Get existing goals to find the one to update
	goals, err := s.fitnessGoalRepo.GetByUserID(ctx, userID, "")
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get fitness goals")
	}

	// Find the goal to update
	var goalToUpdate *model.FitnessGoal
	for _, g := range goals {
		if g.ID == goalID {
			goalToUpdate = g
			break
		}
	}

	if goalToUpdate == nil {
		return nil, errors.ErrResourceNotFound
	}

	// Update fields
	goalToUpdate.GoalType = req.GoalType
	goalToUpdate.GoalDescription = req.GoalDescription
	goalToUpdate.TargetWeight = req.TargetWeight
	goalToUpdate.Deadline = req.Deadline
	goalToUpdate.Priority = req.Priority
	goalToUpdate.UpdatedAt = time.Now()

	if err := s.fitnessGoalRepo.Update(ctx, goalToUpdate); err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update fitness goal")
	}

	return goalToUpdate, nil
}
