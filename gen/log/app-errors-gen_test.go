package log

import (
	"errors"
	"testing"
)

// TestNewDetailedInternalServiceError tests the NewDetailedInternalServiceError function.
func TestNewDetailedInternalServiceError(t *testing.T) {
	t.Parallel()

	reason := "test reason"
	//nolint:err113
	errs := []error{errors.New("test error")}

	// Call the function
	result := NewDetailedInternalServiceError(reason, errs...)

	// Assertions
	if result == nil {
		t.Errorf("NewDetailedInternalServiceError() should not return nil")
	}

	if result.Code != "DetailedInternalServiceError" {
		t.Errorf("NewDetailedInternalServiceError().Code = %s, want %s", result.Code, "DetailedInternalServiceError")
	}

	if result.Text != "internal service error. Reason {{reason}}" {
		t.Errorf("NewDetailedInternalServiceError().Text = %s, want %s", result.Text, "internal service error. Reason {{reason}}")
	}

	if result.Error != "test error" {
		t.Errorf("NewDetailedInternalServiceError().Error = %s, want %s", result.Error, "test error")
	}
}

// TestNewInternalServiceError tests the NewInternalServiceError function.
func TestNewInternalServiceError(t *testing.T) {
	t.Parallel()

	errs := []error{errors.New("test error")}

	// Call the function
	result := NewInternalServiceError(errs...)

	// Assertions
	if result == nil {
		t.Errorf("NewInternalServiceError() should not return nil")
	}

	if result.Code != "InternalServiceError" {
		t.Errorf("NewInternalServiceError().Code = %s, want %s", result.Code, "InternalServiceError")
	}

	if result.Text != "internal service error" {
		t.Errorf("NewInternalServiceError().Text = %s, want %s", result.Text, "internal service error")
	}

	if result.Error != "test error" {
		t.Errorf("NewInternalServiceError().Error = %s, want %s", result.Error, "test error")
	}
}

// TestNewDetailedNumericInternalServiceError tests the NewDetailedNumericInternalServiceError function.
func TestNewDetailedNumericInternalServiceError(t *testing.T) {
	t.Parallel()

	// setup test data.
	errorCode := int64(123)
	errorMessage := errors.New("this is a detailed numeric internal service error")
	detailedError := NewDetailedNumericInternalServiceError(errorCode, errorMessage)

	// assertions
	if detailedError.Code != DetailedNumericInternalServiceError {
		t.Errorf("Expected error code %s, got %s", DetailedNumericInternalServiceError, detailedError.Code)
	}

	if detailedError.Text != "internal service error. Reason {{ reason }}" {
		t.Errorf("Expected error text %s, got %s", "internal service error. Reason {{ reason }}", detailedError.Text)
	}

	if detailedError.Error != "this is a detailed numeric internal service error" {
		t.Errorf("Expected error message %s, got %s", "This is a detailed numeric internal service error.", detailedError.Error)
	}
}
