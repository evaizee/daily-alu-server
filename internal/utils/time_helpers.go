package utils

import (
    "context"
    "fmt"
    "time"
)

// TimeInputParsing parses a time input string into a time.Time
func TimeLocationParsing(ctx context.Context, timeInput string) (time.Time, error) {
    parsedTime, err := time.Parse("2006-01-02T15:04:05.999", timeInput)

    if err != nil {
        parsedTime, err = time.Parse(time.RFC3339, timeInput)
        if err != nil {
            return time.Time{}, fmt.Errorf("invalid happens_at format: %w", err)
        }
    }

    if parsedTime.Location() == time.UTC {
        loc, err := time.LoadLocation("Asia/Bangkok")
        if err != nil {
            return time.Time{}, fmt.Errorf("failed to load timezone: %w", err)
        }

        parsedTime = parsedTime.In(loc)
    }

    return parsedTime, nil
}