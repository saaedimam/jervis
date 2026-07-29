package model

import (
	"context"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/contracts"
)

type Job struct {
	id       string
	name     string
	schedule contracts.Schedule
	handler  func(ctx context.Context) error
}

func NewJob(id, name string, schedule contracts.Schedule, handler func(ctx context.Context) error) contracts.Job {
	return &Job{
		id:       id,
		name:     name,
		schedule: schedule,
		handler:  handler,
	}
}

func (j *Job) ID() string {
	return j.id
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) Schedule() contracts.Schedule {
	return j.schedule
}

func (j *Job) Handle(ctx context.Context) error {
	if j.handler == nil {
		return nil
	}
	return j.handler(ctx)
}
