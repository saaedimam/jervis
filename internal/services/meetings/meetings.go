package meetings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	storecontracts "github.com/ioriimasu/jervis/internal/memory/store/contracts"
	eventcontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/events"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

var (
	ErrInvalidMeeting   = errors.New("meetings: invalid meeting")
	ErrMeetingNotFound  = errors.New("meetings: meeting not found")
	ErrDuplicateMeeting = errors.New("meetings: duplicate meeting ID")
)

type MeetingStatus string

const (
	StatusScheduled MeetingStatus = "SCHEDULED"
	StatusCancelled MeetingStatus = "CANCELLED"
	StatusCompleted MeetingStatus = "COMPLETED"
)

type Meeting struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Location    string        `json:"location"`
	Status      MeetingStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type Service interface {
	CreateMeeting(ctx context.Context, id, title, description string, start, end time.Time, location string) (*Meeting, error)
	GetMeeting(ctx context.Context, id string) (*Meeting, error)
	UpdateMeetingStatus(ctx context.Context, id string, status MeetingStatus) (*Meeting, error)
	ListMeetings(ctx context.Context, from, to time.Time) ([]*Meeting, error)
}

type service struct {
	store     storecontracts.Store
	publisher eventcontracts.Publisher
}

func New(store storecontracts.Store, publisher eventcontracts.Publisher) Service {
	return &service{
		store:     store,
		publisher: publisher,
	}
}

func (s *service) CreateMeeting(ctx context.Context, id, title, description string, start, end time.Time, location string) (*Meeting, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(title) == "" {
		return nil, ErrInvalidMeeting
	}

	now := time.Now().UTC()
	m := &Meeting{
		ID:          id,
		Title:       title,
		Description: description,
		StartTime:   start.UTC(),
		EndTime:     end.UTC(),
		Location:    location,
		Status:      StatusScheduled,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := s.store.Exec(ctx,
		"INSERT INTO meetings (id, title, description, start_time, end_time, location, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		m.ID, m.Title, m.Description, m.StartTime, m.EndTime, m.Location, m.Status, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert meeting: %w", err)
	}

	s.publishEvent(ctx, "meetings.meeting.created", m)

	return m, nil
}

func (s *service) GetMeeting(ctx context.Context, id string) (*Meeting, error) {
	row := s.store.QueryRow(ctx,
		"SELECT id, title, description, start_time, end_time, location, status, created_at, updated_at FROM meetings WHERE id = ?",
		id,
	)

	var m Meeting
	err := row.Scan(&m.ID, &m.Title, &m.Description, &m.StartTime, &m.EndTime, &m.Location, &m.Status, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrMeetingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan meeting: %w", err)
	}

	return &m, nil
}

func (s *service) UpdateMeetingStatus(ctx context.Context, id string, status MeetingStatus) (*Meeting, error) {
	now := time.Now().UTC()

	result, err := s.store.Exec(ctx,
		"UPDATE meetings SET status = ?, updated_at = ? WHERE id = ?",
		status, now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update meeting: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrMeetingNotFound
	}

	m, err := s.GetMeeting(ctx, id)
	if err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "meetings.meeting.updated", m)

	return m, nil
}

func (s *service) ListMeetings(ctx context.Context, from, to time.Time) ([]*Meeting, error) {
	rows, err := s.store.Query(ctx,
		"SELECT id, title, description, start_time, end_time, location, status, created_at, updated_at FROM meetings WHERE start_time >= ? AND start_time <= ? ORDER BY start_time ASC",
		from.UTC(), to.UTC(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query meetings: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []*Meeting
	for rows.Next() {
		var m Meeting
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.StartTime, &m.EndTime, &m.Location, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan meeting row: %w", err)
		}
		results = append(results, &m)
	}

	return results, nil
}

func (s *service) publishEvent(ctx context.Context, eventType string, meeting *Meeting) {
	if s.publisher == nil {
		return
	}

	id, _ := types.NewEventID(fmt.Sprintf("%s-%d", meeting.ID, time.Now().UnixNano()))
	event, err := events.NewBuilder().
		SetID(id).
		SetType(events.EventType(eventType)).
		SetSource("meetings").
		SetPayload(meeting).
		Build()

	if err == nil {
		_ = s.publisher.Publish(event)
	}
}
