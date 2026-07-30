package calendar

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ical "github.com/arran4/golang-ical"
	eventcontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/services/meetings"
)

// Service defines the Calendar Integration Service interface.
type Service interface {
	SyncEvents(ctx context.Context) error
	ImportICal(ctx context.Context, url string) error
	ExportICal(ctx context.Context) (string, error)
}

type service struct {
	publisher       eventcontracts.Publisher
	meetingsService meetings.Service
}

// New constructs a new Calendar Integration Service.
func New(publisher eventcontracts.Publisher, meetingsService meetings.Service) Service {
	return &service{
		publisher:       publisher,
		meetingsService: meetingsService,
	}
}

func (s *service) SyncEvents(ctx context.Context) error {
	// For now, SyncEvents is a no-op or can be mapped to a default set of URLs
	return nil
}

func (s *service) ImportICal(ctx context.Context, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch calendar: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse iCal: %w", err)
	}

	for _, event := range cal.Events() {
		summary := event.GetProperty(ical.ComponentPropertySummary).Value
		description := ""
		if prop := event.GetProperty(ical.ComponentPropertyDescription); prop != nil {
			description = prop.Value
		}

		start, _ := event.GetStartAt()
		end, _ := event.GetEndAt()
		uid := event.GetProperty(ical.ComponentPropertyUniqueId).Value

		location := ""
		if prop := event.GetProperty(ical.ComponentPropertyLocation); prop != nil {
			location = prop.Value
		}

		// Check if meeting already exists to avoid duplicates
		// For simplicity, we use the UID as the ID
		_, err := s.meetingsService.CreateMeeting(ctx, uid, summary, description, start, end, location)
		if err != nil {
			// If it already exists, skip or update.
			// meetings.CreateMeeting currently returns error on duplicate.
			// TODO: Implement UpdateMeeting or CheckExists
			continue
		}
	}

	return nil
}

func (s *service) ExportICal(ctx context.Context) (string, error) {
	now := time.Now()
	from := now.AddDate(0, -1, 0)
	to := now.AddDate(0, 1, 0)

	meetingsList, err := s.meetingsService.ListMeetings(ctx, from, to)
	if err != nil {
		return "", err
	}

	cal := ical.NewCalendar()
	cal.SetMethod(ical.MethodPublish)

	for _, m := range meetingsList {
		event := cal.AddEvent(m.ID)
		event.SetSummary(m.Title)
		event.SetDescription(m.Description)
		event.SetStartAt(m.StartTime)
		event.SetEndAt(m.EndTime)
		event.SetLocation(m.Location)
		event.SetDtStampTime(m.CreatedAt)
	}

	return cal.Serialize(), nil
}
