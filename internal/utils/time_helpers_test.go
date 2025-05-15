package utils

import (
	"context"
	"testing"
	"time"
)

func TestTimeLocationParsingSuccess(t *testing.T) {
	// Test cases with different time formats
	timeInput := "2025-03-28T07:40:00Z"
	expectedTime := time.Date(2025, 3, 28, 7, 40, 0, 0, time.UTC)
	parsedTime, err := TimeLocationParsing(context.Background(), timeInput)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Expected time: %v, but got: %v", expectedTime, parsedTime)
	}
}

func TestTimeLocationParsingFailedInvalidFormat(t *testing.T) {
	// Test cases with different time formats
	timeInput := "28-03-2025T07:40:00"
	_, err := TimeLocationParsing(context.Background(), timeInput)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	expectedError := "invalid happens_at format: parsing time \"28-03-2025T07:40:00\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"28-03-2025T07:40:00\" as \"2006\""
	if err.Error() != expectedError {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}