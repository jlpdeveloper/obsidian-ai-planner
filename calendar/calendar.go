package calendar

import (
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Event struct {
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
	Type  string `json:"type"`
}

func filterEvents(events []*calendar.Event) []Event {
	validEventTypes := make(map[string]struct{})
	validEventTypes["event"] = struct{}{}
	validEventTypes["focusTime"] = struct{}{}
	var filteredEvents []Event
	for _, e := range events {
		if _, ok := validEventTypes[e.EventType]; ok {
			filteredEvents = append(filteredEvents, Event{
				Name:  e.Summary,
				Start: e.Start.DateTime,
				End:   e.End.DateTime,
				Type:  e.EventType,
			})
		}
	}

	return filteredEvents
}

type GoogleCalendarIntegration struct {
	calendarService *calendar.Service
}

func (c GoogleCalendarIntegration) GetCalendarEvents(start time.Time) ([]Event, error) {
	minTime := start.Format(time.RFC3339)
	maxTime := start.Add(24 * time.Hour).Format(time.RFC3339)
	events, err := c.calendarService.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(minTime).TimeMax(maxTime).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	return filterEvents(events.Items), nil
}

func New(ctx context.Context) *GoogleCalendarIntegration {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	svr, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return &GoogleCalendarIntegration{
		calendarService: svr,
	}
}
