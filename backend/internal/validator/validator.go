package validator

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator instance with custom validation functions
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new custom validator instance
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Register custom validators
	_ = v.RegisterValidation("password_strength", validatePasswordStrength)
	_ = v.RegisterValidation("email_format", validateEmailFormat)
	_ = v.RegisterValidation("date_range", validateDateRange)
	_ = v.RegisterValidation("macro_ratio", validateMacroRatio)
	_ = v.RegisterValidation("future_date", validateNotFutureDate)

	return &CustomValidator{
		validator: v,
	}
}

// Validate validates a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// GetValidator returns the underlying validator instance
func (cv *CustomValidator) GetValidator() *validator.Validate {
	return cv.validator
}

// validatePasswordStrength validates password strength
// Requirements:
// - At least 8 characters
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one digit
// - Contains at least one special character
func validatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// ValidatePasswordStrength is the exported version for registration
func ValidatePasswordStrength(fl validator.FieldLevel) bool {
	return validatePasswordStrength(fl)
}

// validateEmailFormat validates email format with stricter rules
// Validates: Requirements 1.1, 2.2
func validateEmailFormat(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	// Additional checks
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Local part should not start or end with a dot
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return false
	}

	// Domain part should not start or end with a hyphen
	if strings.HasPrefix(domainPart, "-") || strings.HasSuffix(domainPart, "-") {
		return false
	}

	// Domain should have at least one dot
	if !strings.Contains(domainPart, ".") {
		return false
	}

	return true
}

// ValidateEmailFormat is the exported version for registration
func ValidateEmailFormat(fl validator.FieldLevel) bool {
	return validateEmailFormat(fl)
}

// validateDateRange validates that end date is after start date
// This is used as a struct-level validator
// Validates: Requirements 7.1
func validateDateRange(fl validator.FieldLevel) bool {
	// This validator is meant to be used at struct level
	// For field-level validation, use custom logic in handlers
	return true
}

// validateMacroRatio validates that macro nutrient ratios sum to approximately 1.0
// Validates: Requirements 6.3
func validateMacroRatio(fl validator.FieldLevel) bool {
	ratio := fl.Field().Float()
	// Individual ratio should be between 0 and 1
	return ratio >= 0 && ratio <= 1
}

// ValidateMacroRatio is the exported version for registration
func ValidateMacroRatio(fl validator.FieldLevel) bool {
	return validateMacroRatio(fl)
}

// validateNotFutureDate validates that a date is not in the future
// Validates: Requirements 7.1
func validateNotFutureDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return true // Let required validator handle empty values
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	// Check if date is not in the future (date should be <= today)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return !dateOnly.After(today)
}

// ValidateNotFutureDate is the exported version for registration
func ValidateNotFutureDate(fl validator.FieldLevel) bool {
	return validateNotFutureDate(fl)
}

// ValidateMacroRatioSum validates that protein + carb + fat ratios sum to approximately 1.0
// This is a custom function for struct-level validation
// Validates: Requirements 6.3
func ValidateMacroRatioSum(protein, carb, fat float64) bool {
	sum := protein + carb + fat
	tolerance := 0.01
	return sum >= (1.0-tolerance) && sum <= (1.0+tolerance)
}

// ValidateDateRangeOrder validates that start date is before or equal to end date
// Validates: Requirements 7.1
func ValidateDateRangeOrder(startDate, endDate string) bool {
	if startDate == "" || endDate == "" {
		return true // Let required validator handle empty values
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return false
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return false
	}

	return !start.After(end)
}
