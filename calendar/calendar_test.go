package calendar

import (
	"testing"

	"google.golang.org/api/calendar/v3"
)

func TestFilterEvents(t *testing.T) {
	events := []*calendar.Event{
		{
			Summary:   "Normal Event",
			EventType: "event",
			Start:     &calendar.EventDateTime{DateTime: "2023-10-27T10:00:00Z"},
			End:       &calendar.EventDateTime{DateTime: "2023-10-27T11:00:00Z"},
		},
		{
			Summary:   "All-Day Event",
			EventType: "event",
			Start:     &calendar.EventDateTime{Date: "2023-10-27"},
			End:       &calendar.EventDateTime{Date: "2023-10-28"},
		},
		{
			Summary:   "Focus Time",
			EventType: "focusTime",
			Start:     &calendar.EventDateTime{DateTime: "2023-10-27T12:00:00Z"},
			End:       &calendar.EventDateTime{DateTime: "2023-10-27T13:00:00Z"},
		},
		{
			Summary:   "Ignored Event",
			EventType: "outOfOffice",
			Start:     &calendar.EventDateTime{DateTime: "2023-10-27T14:00:00Z"},
			End:       &calendar.EventDateTime{DateTime: "2023-10-27T15:00:00Z"},
		},
	}

	filtered := filterEvents(events)

	if len(filtered) != 3 {
		t.Errorf("Expected 3 filtered events, got %d", len(filtered))
	}

	for _, e := range filtered {
		if e.Name == "All-Day Event" {
			if e.Start == "" {
				t.Errorf("All-Day Event Start is empty")
			}
			if e.End == "" {
				t.Errorf("All-Day Event End is empty")
			}
			if e.Start != "2023-10-27" {
				t.Errorf("Expected Start '2023-10-27', got '%s'", e.Start)
			}
			if e.End != "2023-10-28" {
				t.Errorf("Expected End '2023-10-28', got '%s'", e.End)
			}
		}
	}
}
