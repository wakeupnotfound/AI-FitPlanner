package validator

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
)

func TestPasswordStrength(t *testing.T) {
	v := NewCustomValidator()

	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{
			name:     "valid strong password",
			password: "Test@1234",
			valid:    true,
		},
		{
			name:     "valid with special chars",
			password: "MyP@ssw0rd!",
			valid:    true,
		},
		{
			name:     "too short",
			password: "Test@12",
			valid:    false,
		},
		{
			name:     "no uppercase",
			password: "test@1234",
			valid:    false,
		},
		{
			name:     "no lowercase",
			password: "TEST@1234",
			valid:    false,
		},
		{
			name:     "no digit",
			password: "Test@word",
			valid:    false,
		},
		{
			name:     "no special char",
			password: "Test1234",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				Password string `validate:"password_strength"`
			}

			ts := testStruct{Password: tt.password}
			err := v.Validate(ts)

			if tt.valid && err != nil {
				t.Errorf("expected valid password, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected invalid password, got no error")
			}
		})
	}
}

func TestEmailFormat(t *testing.T) {
	v := NewCustomValidator()

	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{
			name:  "valid email",
			email: "user@example.com",
			valid: true,
		},
		{
			name:  "valid with subdomain",
			email: "user@mail.example.com",
			valid: true,
		},
		{
			name:  "valid with plus",
			email: "user+tag@example.com",
			valid: true,
		},
		{
			name:  "invalid - no @",
			email: "userexample.com",
			valid: false,
		},
		{
			name:  "invalid - starts with dot",
			email: ".user@example.com",
			valid: false,
		},
		{
			name:  "invalid - ends with dot",
			email: "user.@example.com",
			valid: false,
		},
		{
			name:  "invalid - no domain extension",
			email: "user@example",
			valid: false,
		},
		{
			name:  "invalid - domain starts with hyphen",
			email: "user@-example.com",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				Email string `validate:"email_format"`
			}

			ts := testStruct{Email: tt.email}
			err := v.Validate(ts)

			if tt.valid && err != nil {
				t.Errorf("expected valid email, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected invalid email, got no error")
			}
		})
	}
}

func TestMacroRatio(t *testing.T) {
	v := NewCustomValidator()

	tests := []struct {
		name  string
		ratio float64
		valid bool
	}{
		{
			name:  "valid ratio 0.3",
			ratio: 0.3,
			valid: true,
		},
		{
			name:  "valid ratio 0",
			ratio: 0.0,
			valid: true,
		},
		{
			name:  "valid ratio 1",
			ratio: 1.0,
			valid: true,
		},
		{
			name:  "invalid negative",
			ratio: -0.1,
			valid: false,
		},
		{
			name:  "invalid over 1",
			ratio: 1.1,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				Ratio float64 `validate:"macro_ratio"`
			}

			ts := testStruct{Ratio: tt.ratio}
			err := v.Validate(ts)

			if tt.valid && err != nil {
				t.Errorf("expected valid ratio, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected invalid ratio, got no error")
			}
		})
	}
}

func TestNotFutureDate(t *testing.T) {
	v := NewCustomValidator()

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	tests := []struct {
		name  string
		date  string
		valid bool
	}{
		{
			name:  "today is valid",
			date:  today,
			valid: true,
		},
		{
			name:  "yesterday is valid",
			date:  yesterday,
			valid: true,
		},
		{
			name:  "tomorrow is invalid",
			date:  tomorrow,
			valid: false,
		},
		{
			name:  "past date is valid",
			date:  "2020-01-01",
			valid: true,
		},
		{
			name:  "invalid format",
			date:  "2020/01/01",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				Date string `validate:"future_date"`
			}

			ts := testStruct{Date: tt.date}
			err := v.Validate(ts)

			if tt.valid && err != nil {
				t.Errorf("expected valid date, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected invalid date, got no error")
			}
		})
	}
}

func TestValidateMacroRatioSum(t *testing.T) {
	tests := []struct {
		name    string
		protein float64
		carb    float64
		fat     float64
		valid   bool
	}{
		{
			name:    "valid sum exactly 1.0",
			protein: 0.3,
			carb:    0.4,
			fat:     0.3,
			valid:   true,
		},
		{
			name:    "valid sum within tolerance",
			protein: 0.33,
			carb:    0.33,
			fat:     0.34,
			valid:   true,
		},
		{
			name:    "invalid sum too low",
			protein: 0.2,
			carb:    0.3,
			fat:     0.3,
			valid:   false,
		},
		{
			name:    "invalid sum too high",
			protein: 0.4,
			carb:    0.4,
			fat:     0.4,
			valid:   false,
		},
		{
			name:    "valid edge case lower bound",
			protein: 0.33,
			carb:    0.33,
			fat:     0.33,
			valid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateMacroRatioSum(tt.protein, tt.carb, tt.fat)
			if result != tt.valid {
				t.Errorf("ValidateMacroRatioSum(%v, %v, %v) = %v, want %v",
					tt.protein, tt.carb, tt.fat, result, tt.valid)
			}
		})
	}
}

func TestValidateDateRangeOrder(t *testing.T) {
	tests := []struct {
		name      string
		startDate string
		endDate   string
		valid     bool
	}{
		{
			name:      "valid range",
			startDate: "2024-01-01",
			endDate:   "2024-12-31",
			valid:     true,
		},
		{
			name:      "valid same date",
			startDate: "2024-01-01",
			endDate:   "2024-01-01",
			valid:     true,
		},
		{
			name:      "invalid reversed",
			startDate: "2024-12-31",
			endDate:   "2024-01-01",
			valid:     false,
		},
		{
			name:      "invalid start format",
			startDate: "2024/01/01",
			endDate:   "2024-12-31",
			valid:     false,
		},
		{
			name:      "invalid end format",
			startDate: "2024-01-01",
			endDate:   "2024/12/31",
			valid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateDateRangeOrder(tt.startDate, tt.endDate)
			if result != tt.valid {
				t.Errorf("ValidateDateRangeOrder(%v, %v) = %v, want %v",
					tt.startDate, tt.endDate, result, tt.valid)
			}
		})
	}
}

// Test validator registration
func TestNewCustomValidator(t *testing.T) {
	v := NewCustomValidator()
	if v == nil {
		t.Fatal("NewCustomValidator returned nil")
	}

	if v.validator == nil {
		t.Fatal("validator instance is nil")
	}

	// Test that custom validators are registered
	type testStruct struct {
		Password string `validate:"password_strength"`
	}

	ts := testStruct{Password: "weak"}
	err := v.Validate(ts)
	if err == nil {
		t.Error("expected validation error for weak password")
	}

	// Check that it's a validation error
	if _, ok := err.(validator.ValidationErrors); !ok {
		t.Errorf("expected validator.ValidationErrors, got %T", err)
	}
}
